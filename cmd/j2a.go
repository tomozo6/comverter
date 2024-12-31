/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/spf13/cobra"
)

var (
	quality int
)

// j2aCmd represents the j2a command
var j2aCmd = &cobra.Command{
	Use:   "j2a",
	Short: "Convert JPG to AVIF",
	Long:  "Convert JPG to AVIF.",

	Run: func(cmd *cobra.Command, args []string) {
		input := cmd.Flag("input").Value.String()
		output := getOutputFileName(cmd.Flag("input").Value.String())

		// Debug
		// fmt.Println("j2a called")
		// fmt.Println("input:", input)
		// fmt.Println("output:", output)
		// fmt.Println(quality)

		// JPG -> AVIF変換
		j2a(input, output, quality)
	},
}

func init() {
	rootCmd.AddCommand(j2aCmd)
	j2aCmd.Flags().StringP("input", "i", "input.jpg", "Input Jpeg file.")
	j2aCmd.Flags().IntVarP(&quality, "quality", "q", 30, "Quality of the output Avif file.")
}

// j2a converts a JPG image to an AVIF image.
// The input file is specified by inputFile, and the output file is specified by outputFile.
// The quality parameter specifies the quality of the output AVIF image.
func j2a(inputFile string, outputFile string, quality int) {
	// VIPS初期化
	vips.Startup(nil)
	defer vips.Shutdown()

	// 入力ファイルを読み込む
	img, err := vips.NewImageFromFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to load image: %v", err)
	}

	// 画像をAVIF形式に変換
	options := vips.NewAvifExportParams()
	options.Quality = quality

	output, _, err := img.ExportAvif(options)
	if err != nil {
		log.Fatalf("Failed to convert to AVIF: %v", err)
	}

	// 出力ファイルに書き込み
	err = os.WriteFile(outputFile, output, 0644)
	if err != nil {
		log.Fatalf("Failed to write output file: %v", err)
	}

	fmt.Println("Successfully converted JPG to AVIF!")
}

// getOutputFileName returns the output file name.
// The output file name is the same as the input file name, but with the extension changed to ".avif".
// For example, "input.jpg" -> "input.avif".
func getOutputFileName(inputFileName string) string {
	return filepath.Base(inputFileName[:len(inputFileName)-len(filepath.Ext(inputFileName))]) + ".avif"
}
