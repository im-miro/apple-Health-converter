# Apple Health データ PDFレポート生成ツール

Apple のヘルスケアアプリからエクスポートした ZIP ファイルをもとに、過去の睡眠・体重・BMI・体脂肪率・心拍数の推移を集計し、PDF 形式でレポートとして出力するツールです。

## 🔧 構成ファイル一覧

| ファイル名                     | 説明 |
|------------------------------|------|
| `main.go`                    | ZIP を解凍し、表形式の PDF (`health_report.pdf`) を生成 |
| `generate_sleep_chart.go`    | 睡眠時間の推移グラフ (`sleep_chart.png`) を生成 |
| `generate_bodymass_chart.go` | 体重の推移グラフ (`bodymass_chart.png`) を生成 |
| `generate_bmi_chart.go`      | BMI の推移グラフ (`bmi_chart.png`) を生成 |
| `generate_bodyfat_chart.go`  | 体脂肪率の推移グラフ (`bodyfat_chart.png`) を生成 |
| `append_graphs_to_pdf.go`    | 上記 PNG 画像を `health_report.pdf` に結合し `health_report_with_graphs.pdf` を生成 |
| `run_all.sh`                 | 上記処理を一括で実行するシェルスクリプト |
| `ipaexg.ttf`                 | 日本語対応フォントファイル（PDF内で使用） |

## 🚀 使い方

### 1. 必要パッケージのインストール

```bash
go mod tidy
```

### 2. Appleヘルスケアのデータをエクスポート

1. iPhone の「ヘルスケア」アプリで「データを書き出す」
2. `export.zip` 形式で保存し、Mac や Linux 環境に転送

### 3. 実行

```bash
bash run_all.sh export.zip 3
```

> `3` は集計対象とする過去月数（1〜6）です。

### 4. 出力ファイル

- `health_report_with_graphs.pdf`：グラフ統合済みのレポート
- `health_report.pdf`：表のみのPDF
- `*_chart.png`：各種グラフ画像（自動で削除されます）

## 📄 ライセンス

未定、連絡してください。
