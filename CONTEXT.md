# Image Conversion

JPEG 画像および PDF 文書を AVIF 画像へ変換するコマンドラインツールのコンテキストです。

## Language

**Source JPEG**:
AVIF へ変換する入力 JPEG ファイル。
_Avoid_: input image, source image

**AVIF output**:
Source JPEG から生成される AVIF ファイル。
_Avoid_: converted image, result file

**Batch conversion**:
入力ディレクトリを再帰的に走査して、見つかった Source JPEG ごとに AVIF output を生成する処理。
_Avoid_: directory conversion, bulk conversion

**Source PDF**:
ページごとに AVIF output を生成する入力 PDF 文書。
_Avoid_: input document, source file

**PDF page output**:
Source PDF の 1 ページから生成される AVIF output。ファイル名は 3 桁以上のページ番号と `.avif` をこの順に連結する。
_Avoid_: PDF output, converted page
