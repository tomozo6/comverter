# comverter

JPEG 画像と PDF の各ページを AVIF に変換する Go 製 CLI ツールです。単一ファイルの変換、ディレクトリ内の画像の一括変換、PDF ページの変換に対応しています。

## 必要要件

- Go 1.22 以降
- [libvips](https://www.libvips.org/)（AVIF エンコーダーを含むビルド）

macOS では Homebrew で libvips を導入できます。

```sh
brew install vips
```

## インストール

リポジトリを取得してビルドします。

```sh
git clone https://github.com/tomozo6/comverter.git
cd comverter
go build -o comverter .
```

または、Go から直接実行できます。

```sh
go run . --help
```

## 使い方

### JPEG ファイルを 1 枚変換する

`j2a` は JPEG を AVIF に変換します。出力ファイル名は入力ファイルのベース名に `.avif` を付けたものになります。

```sh
./comverter j2a --input path/to/photo.jpg
# photo.avif をカレントディレクトリに出力
```

画質は `--quality`（`-q`）で指定できます。既定値は `30` です。

```sh
./comverter j2a -i path/to/photo.jpg -q 80
```

### ディレクトリ内の画像を一括変換する

`j2adir` は指定ディレクトリを再帰的に走査し、読み込めるファイルを AVIF に変換します。出力先は、入力ディレクトリの絶対パスに `_avif` を付けたディレクトリです。

```sh
./comverter j2adir --input path/to/photos
# 例: /absolute/path/to/photos_avif に出力
```

画質を指定する場合:

```sh
./comverter j2adir -i path/to/photos -q 80
```

> 注意: 一括変換の出力ではサブディレクトリ構造を保持しません。同じファイル名が複数ある場合は、上書きを防ぐためエラーで終了します。

### PDF の各ページを変換する

`p2a` は PDF の全ページを個別の AVIF に変換します。既定ではカレントディレクトリへ、`sample-001.avif` のように出力します。

```sh
./comverter p2a --input sample.pdf
```

`--output`（`-o`）で出力ディレクトリを指定できます。指定先が存在しない場合は自動で作成されます。

```sh
./comverter p2a -i sample.pdf -o output/pages -q 80
```

## コマンド一覧

| コマンド | 説明 |
| --- | --- |
| `j2a -i <file> [-q <quality>]` | JPEG ファイルを AVIF に変換 |
| `j2adir -i <directory> [-q <quality>]` | ディレクトリを再帰的に変換 |
| `p2a -i <file> [-o <directory>] [-q <quality>]` | PDF の各ページを AVIF に変換 |
| `completion <shell>` | シェル補完スクリプトを生成 |

詳細は次で確認できます。

```sh
./comverter j2a --help
./comverter j2adir --help
./comverter p2a --help
```

## ライセンス

[MIT License](LICENSE)
