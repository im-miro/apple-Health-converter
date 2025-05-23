package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/mholt/archiver/v3"
)

type HealthRecord struct {
	Type      string `xml:"type,attr"`
	Value     string `xml:"value,attr"`
	StartDate string `xml:"startDate,attr"`
	EndDate   string `xml:"endDate,attr"`
}

type Root struct {
	Records []HealthRecord `xml:"Record"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("使用方法: go run main.go <zipファイル名> [1〜6の集計月数]")
	}
	zipPath := os.Args[1]
	unzipPath := "extracted"
	os.RemoveAll(unzipPath)

	months := 3
	if len(os.Args) >= 3 {
		m, err := strconv.Atoi(os.Args[2])
		if err != nil || m < 1 || m > 6 {
			log.Fatalf("集計月数は1〜6の整数で指定してください")
		}
		months = m
	}
	cutoff := time.Now().AddDate(0, -months, 0)

	fmt.Print("📄 レポートのタイトルを入力してください（例: 健康診断2025春）: ")
	var customTitle string
	fmt.Scanln(&customTitle)

	err := archiver.Unarchive(zipPath, unzipPath)
	if err != nil {
		log.Fatalf("ZIP解凍失敗: %v", err)
	}

	xmlPath := filepath.Join(unzipPath, "apple_health_export", "export.xml")
	xmlFile, err := os.ReadFile(xmlPath)
	if err != nil {
		log.Fatalf("export.xml 読み込み失敗: %v", err)
	}

	var root Root
	if err := xml.Unmarshal(xmlFile, &root); err != nil {
		log.Fatalf("XML解析失敗: %v", err)
	}

	typesToLabel := map[string]string{
		"HKQuantityTypeIdentifierHeartRate":         "HeartRate",
		"HKQuantityTypeIdentifierBodyMass":          "BodyMass",
		"HKQuantityTypeIdentifierBodyMassIndex":     "BMI",
		"HKQuantityTypeIdentifierBodyFatPercentage": "BodyFatPercent",
		"HKCategoryTypeIdentifierSleepAnalysis":     "SleepHours",
	}
	type stat struct {
		sum   float64
		count int
	}
	stats := map[string]map[string]stat{}

	sleepTotals := map[string]float64{}

	for _, rec := range root.Records {
		if rec.Type == "HKCategoryTypeIdentifierSleepAnalysis" {
			if rec.Value != "HKCategoryValueSleepAnalysisAsleepREM" &&
				rec.Value != "HKCategoryValueSleepAnalysisAsleepCore" &&
				rec.Value != "HKCategoryValueSleepAnalysisAsleepDeep" {
				continue
			}
			start, _ := time.Parse("2006-01-02 15:04:05 -0700", rec.StartDate)
			end, _ := time.Parse("2006-01-02 15:04:05 -0700", rec.EndDate)
			if start.Before(cutoff) || end.Before(cutoff) {
				continue
			}
			duration := end.Sub(start).Hours()
			date := start.Format("2006-01-02")
			sleepTotals[date] += duration
			continue
		}

		label, ok := typesToLabel[rec.Type]
		if !ok {
			continue
		}
		t, err := time.Parse("2006-01-02 15:04:05 -0700", rec.StartDate)
		if err != nil || t.Before(cutoff) {
			continue
		}
		val, err := strconv.ParseFloat(rec.Value, 64)
		if err != nil {
			continue
		}
		if rec.Type == "HKQuantityTypeIdentifierBodyFatPercentage" {
			val *= 100
		}
		date := t.Format("2006-01-02")
		if _, ok := stats[date]; !ok {
			stats[date] = map[string]stat{}
		}
		s := stats[date][label]
		s.sum += val
		s.count++
		stats[date][label] = s
	}

	for date, total := range sleepTotals {
		if total < 1.0 {
			continue
		}
		if _, ok := stats[date]; !ok {
			stats[date] = map[string]stat{}
		}
		stats[date]["SleepHours"] = stat{sum: total, count: 1}
	}

	var dates []string
	for d := range stats {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddUTF8Font("IPAexG", "", "ipaexg.ttf")
	pdf.AddPage()
	pdf.SetFont("IPAexG", "", 14)

	fullTitle := fmt.Sprintf("Health Report (過去%dヶ月分)", months)
	if customTitle != "" {
		fullTitle += " - " + customTitle
	}
	pdf.Cell(40, 10, fullTitle)
	pdf.Ln(12)

	headers := []string{"日付", "心拍数 (bpm)", "体重 (kg)", "BMI", "体脂肪率 (%)", "睡眠時間 (h)"}
	colWidth := 45.0

	pdf.SetFont("IPAexG", "", 10)
	for _, h := range headers {
		pdf.SetFillColor(220, 220, 220)
		pdf.CellFormat(colWidth, 10, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	fill := false
	for _, date := range dates {
		if fill {
			pdf.SetFillColor(245, 245, 245)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}
		pdf.CellFormat(colWidth, 10, date, "1", 0, "", true, 0, "")
		for _, key := range []string{"HeartRate", "BodyMass", "BMI", "BodyFatPercent", "SleepHours"} {
			val := "-"
			if s, ok := stats[date][key]; ok && s.count > 0 {
				val = fmt.Sprintf("%.1f", s.sum/float64(s.count))
			}
			pdf.CellFormat(colWidth, 10, val, "1", 0, "", true, 0, "")
		}
		pdf.Ln(-1)
		fill = !fill
	}

	if err := pdf.OutputFileAndClose("health_report.pdf"); err != nil {
		log.Fatalf("PDF保存失敗: %v", err)
	}

	fmt.Println("✅ PDF出力成功: health_report.pdf")
}
