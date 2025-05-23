# 🍎 Apple Health Converter

Apple 端末からエクスポートした `export.zip`（Health データ）を読み込み、
健康指標（心拍数・体重・BMI・体脂肪率・睡眠時間）を PDF レポートとして出力します。
また、各指標のトレンドをグラフ化し PDF に統合できます。

---

## 📦 構成ファイル一覧

| ファイル名                   | 説明                                        |
| ---------------------------- | ------------------------------------------- |
| `generate_health_report.go`  | データ集計と PDF 表形式出力を行うメイン処理 |
| `generate_sleep_chart.go`    | 睡眠時間の推移グラフを出力（PNG）           |
| `generate_bodymass_chart.go` | 体重の推移グラフを出力（PNG）               |
| `generate_bmi_chart.go`      | BMI の推移グラフを出力（PNG）               |
| `generate_bodyfat_chart.go`  | 体脂肪率の推移グラフを出力（PNG）           |
| `append_graphs_to_pdf.go`    | 生成されたグラフ画像を PDF に追加統合       |
| `run_all.sh`                 | 上記処理を一括実行するシェルスクリプト      |
| `ipaexg.ttf`                 | PDF 内で日本語表示を行うためのフォント      |

---

## 🚀 使い方

### 1. 必要な準備

- Go 1.20 以上がインストールされていること
- Apple ヘルスケアアプリから `export.zip` をエクスポートしておく（例: `書き出したデータ.zip or export.zip等`）

### 2. 一括実行コマンド

```bash
sh run_all.sh <zipファイル> <集計月数 (1〜6)>
```

### 例:

```bash
sh run_all.sh export.zip 3
```

処理が完了すると、以下のファイルが出力されます：

- `health_report_with_graphs.pdf`: 表形式＋グラフを統合した PDF レポート
- `sleep_chart.png`、`bodymass_chart.png` などの一時画像（実行後に削除されます）

---

## 📂 出力内容

### PDF レポート

- 日付別に以下の平均値を出力
  - 心拍数 (bpm)
  - 体重 (kg)
  - BMI
  - 体脂肪率 (%)
  - 睡眠時間 (h)

### グラフ画像（PDF 末尾に追加）

- 睡眠時間推移
- 体重推移
- BMI 推移
- 体脂肪率推移

---

## 🤝 コントリビューター募集

このプロジェクトはオープンソースであり、貢献を歓迎します！  
バグ修正、機能追加、ドキュメント改善など、どんな形でも構いません。お気軽に参加してください！

### ✅ 貢献の手順

1. **このリポジトリをフォーク**  
   GitHub 上の [Fork] ボタンを押して自分のアカウントにコピーします。

2. **ローカルにクローン**

```bash
git clone https://github.com/<あなたのユーザー名>/apple-Health-converter.git
cd apple-Health-converter
```

3. **新しいブランチを作成**

```bash
git checkout -b feature/あなたの機能名
```

4. **変更を加えてコミット**

```bash
git commit -m "Add: 新機能の説明や修正内容を記載"
```

5. **プッシュして Pull Request を作成**

```bash
git push origin feature/あなたの機能名
```

GitHub 上で Pull Request を作成し、修正内容の説明を記入してください。

### 📌 注意点

- Go 1.20 以上での動作確認を推奨しています。
- `run_all.sh` で一括実行できることを確認してください。
- `README.md` やコメントの日本語対応も歓迎です。
- メジャーな変更（大幅な仕様変更など）は事前に Issue を立ててご相談ください。

---

## 📝 ライセンス

考え中ですが、個人利用で再配布しない場合は連絡いりません
