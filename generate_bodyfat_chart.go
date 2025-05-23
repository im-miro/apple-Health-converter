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

	bodyFatData := map[string]float64{}
	counts := map[string]int{}

	for _, rec := range root.Records {
		if rec.Type != "HKQuantityTypeIdentifierBodyFatPercentage" {
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
		val *= 100 // 体脂肪率は0.0〜1.0の小数として記録されているため %
		date := t.Format("2006-01-02")
		bodyFatData[date] += val
		counts[date]++
	}

	var dates []string
	for d := range bodyFatData {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	xValues := []time.Time{}
	yValues := []float64{}
	for _, d := range dates {
		parsed, _ := time.Parse("2006-01-02", d)
		xValues = append(xValues, parsed)
		avg := bodyFatData[d] / float64(counts[d])
		yValues = append(yValues, avg)
	}

	graph := chart.Chart{
		Title:      "Body Fat Percentage Trend",
		TitleStyle: chart.Style{FontSize: 14, FontColor: chart.ColorBlack},
		Background: chart.Style{Padding: chart.Box{Top: 40}},
		XAxis: chart.XAxis{
			Name:           "Date",
			NameStyle:      chart.Style{FontSize: 10, FontColor: chart.ColorBlack},
			Style:          chart.Style{FontSize: 10},
			ValueFormatter: chart.TimeDateValueFormatter,
		},
		YAxis: chart.YAxis{
			Name:      "Body Fat (%)",
			NameStyle: chart.Style{FontSize: 10, FontColor: chart.ColorBlack},
			Style:     chart.Style{FontSize: 10},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: xValues,
				YValues: yValues,
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(3),
					FillColor:   chart.ColorAlternateGray,
				},
			},
		},
	}

	f, err := os.Create("bodyfat_chart.png")
	if err != nil {
		log.Fatalf("画像ファイル作成失敗: %v", err)
	}
	defer f.Close()

	err = graph.Render(chart.PNG, f)
	if err != nil {
		log.Fatalf("グラフ出力失敗: %v", err)
	}

	log.Println("✅ bodyfat_chart.png を生成しました")
}
