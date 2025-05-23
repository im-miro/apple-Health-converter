#!/bin/bash

# 使用法チェック
if [ $# -lt 2 ]; then
  echo "使い方: $0 <zipファイル> <集計月数 (1〜6)>"
  exit 1
fi

ZIP_FILE="$1"
MONTHS="$2"

# ファイル存在確認
if [ ! -f "$ZIP_FILE" ]; then
  echo "❌ ZIPファイルが存在しません: $ZIP_FILE"
  exit 1
fi

# 集計月数バリデーション（1〜6）
case "$MONTHS" in
  [1-6]) ;;
  *) echo "❌ 集計月数は 1〜6 の整数で指定してください"; exit 1 ;;
esac

# ステップ1: generate_report.go 実行（PDF + データ表作成）
echo "📄 generate_report.go を実行中..."
go run generate_report.go "$ZIP_FILE" "$MONTHS"

# ステップ2: グラフ生成
for file in generate_sleep_chart.go generate_bodymass_chart.go generate_bmi_chart.go generate_bodyfat_chart.go; do
  echo "📊 $file を実行中..."
  if ! go run "$file" "$MONTHS"; then
    echo "❌ $file の実行に失敗しました"
    exit 1
  fi
done

# ステップ3: 画像をPDFに統合
echo "📎 グラフ画像をPDFに統合中..."
if ! go run append_graphs_to_pdf.go; then
  echo "❌ PDF統合処理に失敗しました"
  exit 1
fi

echo "✅ すべての処理が完了しました！ -> health_report_with_graphs.pdf"
