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
}

type Root struct {
	Records []HealthRecord `xml:"Record"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("使用方法: go run main.go <zipファイル名>")
	}
	zipPath := os.Args[1]
	unzipPath := "extracted"

	// 💡 すでに存在する場合は削除
	os.RemoveAll(unzipPath)

	err := archiver.Unarchive(zipPath, unzipPath)
	if err != nil {
		log.Fatalf("ZIPファイルの解凍に失敗: %v", err)
	}

	// XML を読み込む
	xmlPath := filepath.Join(unzipPath, "apple_health_export", "export.xml")
	xmlFile, err := os.ReadFile(xmlPath)
	if err != nil {
		log.Fatalf("export.xml 読み込み失敗: %v", err)
	}

	var root Root
	if err := xml.Unmarshal(xmlFile, &root); err != nil {
		log.Fatalf("XML解析失敗: %v", err)
	}

	// 対象データ
	typesToLabel := map[string]string{
		"HKQuantityTypeIdentifierHeartRate":         "HeartRate",
		"HKQuantityTypeIdentifierBodyMass":          "BodyMass",
		"HKQuantityTypeIdentifierBodyMassIndex":     "BMI",
		"HKQuantityTypeIdentifierBodyFatPercentage": "BodyFatPercent",
	}

	// データ整形
	type stat struct {
		sum   float64
		count int
	}
	stats := map[string]map[string]stat{}

	for _, rec := range root.Records {
		label, ok := typesToLabel[rec.Type]
		if !ok {
			continue
		}
		t, err := time.Parse("2006-01-02 15:04:05 -0700", rec.StartDate)
		if err != nil {
			continue
		}
		val, err := strconv.ParseFloat(rec.Value, 64)
		if err != nil {
			continue
		}
		// 👇 BodyFatPercentage だけは 100倍して %
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

	// 日付並び替え
	var dates []string
	for d := range stats {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	// PDF出力（横向き）
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Health Report")
	pdf.Ln(12)

	headers := []string{"Date", "HeartRate (bpm)", "BodyMass (kg)", "BMI", "BodyFat (%)"}
	pdf.SetFont("Arial", "B", 10)
	for _, h := range headers {
		pdf.CellFormat(50, 10, h, "1", 0, "", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	for _, date := range dates {
		pdf.CellFormat(50, 10, date, "1", 0, "", false, 0, "")
		for _, key := range []string{"HeartRate", "BodyMass", "BMI", "BodyFatPercent"} {
			val := "-"
			if s, ok := stats[date][key]; ok && s.count > 0 {
				val = fmt.Sprintf("%.1f", s.sum/float64(s.count))
			}
			pdf.CellFormat(50, 10, val, "1", 0, "", false, 0, "")
		}
		pdf.Ln(-1)
	}

	if err := pdf.OutputFileAndClose("health_report.pdf"); err != nil {
		log.Fatalf("PDF出力失敗: %v", err)
	}

	fmt.Println("✅ PDF出力成功: health_report.pdf")
}
