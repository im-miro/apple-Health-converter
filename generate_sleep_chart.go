package main

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/wcharczuk/go-chart/v2"
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
	months := 3
	if len(os.Args) >= 2 {
		m, err := strconv.Atoi(os.Args[1])
		if err == nil && m >= 1 && m <= 6 {
			months = m
		}
	}
	cutoff := time.Now().AddDate(0, -months, 0)

	xmlPath := filepath.Join("extracted", "apple_health_export", "export.xml")
	xmlFile, err := os.ReadFile(xmlPath)
	if err != nil {
		log.Fatalf("export.xml 読み込み失敗: %v", err)
	}

	var root Root
	if err := xml.Unmarshal(xmlFile, &root); err != nil {
		log.Fatalf("XML解析失敗: %v", err)
	}

	sleepData := map[string]float64{}
	for _, rec := range root.Records {
		if rec.Type != "HKCategoryTypeIdentifierSleepAnalysis" {
			continue
		}
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
		sleepData[date] += duration
	}

	var dates []string
	for d := range sleepData {
		if sleepData[d] < 1.0 {
			continue
		}
		dates = append(dates, d)
	}
	sort.Strings(dates)

	xValues := []time.Time{}
	yValues := []float64{}
	for _, d := range dates {
		parsed, _ := time.Parse("2006-01-02", d)
		xValues = append(xValues, parsed)
		yValues = append(yValues, sleepData[d])
	}

	graph := chart.Chart{
		Title:      "sleep Trend",
		TitleStyle: chart.Style{FontSize: 14, FontColor: chart.ColorBlack},
		Background: chart.Style{Padding: chart.Box{Top: 40}},
		XAxis: chart.XAxis{
			Name:           "date",
			NameStyle:      chart.Style{FontSize: 10, FontColor: chart.ColorBlack},
			Style:          chart.Style{FontSize: 10},
			ValueFormatter: chart.TimeDateValueFormatter,
		},
		YAxis: chart.YAxis{
			Name:      "sleep time (h)",
			NameStyle: chart.Style{FontSize: 10, FontColor: chart.ColorBlack},
			Style:     chart.Style{FontSize: 10},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: xValues,
				YValues: yValues,
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(0),
					FillColor:   chart.ColorAlternateGray,
				},
			},
		},
	}

	f, err := os.Create("sleep_chart.png")
	if err != nil {
		log.Fatalf("画像ファイル作成失敗: %v", err)
	}
	defer f.Close()

	err = graph.Render(chart.PNG, f)
	if err != nil {
		log.Fatalf("グラフ出力失敗: %v", err)
	}

	log.Println("✅ sleep_chart.png を生成しました")
}
