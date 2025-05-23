#!/bin/bash

# 使用法チェック
if [ $# -lt 2 ]; then
  echo "使い方: $0 <zipファイル> <集計月数 (1〜6)>"
  exit 1
fi

ZIP_FILE="$1"
MONTHS="$2"

# 集計月数バリデーション
if ! [[ "$MONTHS" =~ ^[1-6]$ ]]; then
  echo "❌ 集計月数は 1〜6 の整数で指定してください"
  exit 1
fi

# ステップ1: main.go 実行（PDF + データ表作成）
echo "📄 main.go を実行中..."
go run main.go "$ZIP_FILE" "$MONTHS"
if [ $? -ne 0 ]; then
  echo "❌ main.go の実行に失敗しました"
  exit 1
fi

# ステップ2: グラフ生成
for file in generate_sleep_chart.go generate_bodymass_chart.go generate_bmi_chart.go generate_bodyfat_chart.go; do
  echo "📊 $file を実行中..."
  go run "$file" "$MONTHS"
  if [ $? -ne 0 ]; then
    echo "❌ $file の実行に失敗しました"
    exit 1
  fi
done

# ステップ3: 画像をPDFに統合
echo "📎 グラフ画像をPDFに統合中..."
go run append_graphs_to_pdf.go
if [ $? -ne 0 ]; then
  echo "❌ PDF統合処理に失敗しました"
  exit 1
fi

echo "✅ すべての処理が完了しました！"
