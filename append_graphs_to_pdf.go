package main

import (
	"log"
	"os"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	outputPDF := "health_report_with_graphs.pdf"

	graphImages := []string{
		"sleep_chart.png",
		"bodymass_chart.png",
		"bmi_chart.png",
		"bodyfat_chart.png",
	}

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddUTF8Font("IPAexG", "", "ipaexg.ttf")

	for _, img := range graphImages {
		if _, err := os.Stat(img); err == nil {
			pdf.AddPage()
			pdf.SetFont("IPAexG", "", 14)
			pdf.CellFormat(40, 10, img+" グラフ", "", 0, "", false, 0, "")
			pdf.ImageOptions(
				img,
				10, 20, 270, 0,
				false,
				gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true},
				0,
				"",
			)
		} else {
			log.Printf("⚠️ 画像が見つかりません: %s", img)
		}
	}

	if err := pdf.OutputFileAndClose(outputPDF); err != nil {
		log.Fatalf("❌ PDF出力失敗: %v", err)
	}

	log.Printf("✅ %s を生成しました", outputPDF)

	// 不要になった画像ファイルを削除
	for _, img := range graphImages {
		if err := os.Remove(img); err != nil {
			log.Printf("⚠️ %s の削除に失敗しました: %v", img, err)
		} else {
			log.Printf("🗑 削除しました: %s", img)
		}
	}
}
