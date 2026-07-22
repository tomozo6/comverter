/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/spf13/cobra"
)

var vipsStartupOnce sync.Once

// j2aCmd represents the j2a command
var j2aCmd = &cobra.Command{
	Use:   "j2a",
	Short: "Convert JPG to AVIF",
	Long:  "Convert JPG to AVIF.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := cmd.Flags().GetString("input")
		if err != nil {
			return err
		}
		quality, err := cmd.Flags().GetInt("quality")
		if err != nil {
			return err
		}
		if err := j2a(input, getOutputFileName(input), quality); err != nil {
			return err
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), "Successfully converted JPG to AVIF!")
		return err
	},
}

func init() {
	rootCmd.AddCommand(j2aCmd)
	j2aCmd.Flags().StringP("input", "i", "input.jpg", "Input Jpeg file.")
	j2aCmd.Flags().IntP("quality", "q", 30, "Quality of the output Avif file (0-100).")
}

// j2a converts a JPG image to an AVIF image.
// The input file is specified by inputFile, and the output file is specified by outputFile.
// The quality parameter specifies the quality of the output AVIF image.
func j2a(inputFile string, outputFile string, quality int) error {
	if err := validateQuality(quality); err != nil {
		return err
	}
	if err := validateJPEG(inputFile); err != nil {
		return err
	}

	// libvips can only be started once in a process. It is released when the
	// process exits, allowing batch conversion to reuse the same instance.
	vipsStartupOnce.Do(func() {
		vips.Startup(nil)
	})

	// 入力ファイルを読み込む
	img, err := vips.NewImageFromFile(inputFile)
	if err != nil {
		return fmt.Errorf("load JPEG %q: %w", inputFile, err)
	}
	defer img.Close()

	// 画像をAVIF形式に変換
	options := vips.NewAvifExportParams()
	options.Quality = quality

	output, _, err := img.ExportAvif(options)
	if err != nil {
		return fmt.Errorf("convert JPEG %q to AVIF: %w", inputFile, err)
	}

	// 出力ファイルに書き込み
	err = os.WriteFile(outputFile, output, 0o644)
	if err != nil {
		return fmt.Errorf("write AVIF output %q: %w", outputFile, err)
	}

	return nil
}

func validateQuality(quality int) error {
	if quality < 0 || quality > 100 {
		return fmt.Errorf("quality must be between 0 and 100, got %d", quality)
	}
	return nil
}

func validateJPEG(path string) error {
	isJPEG, err := isJPEG(path)
	if err != nil {
		return fmt.Errorf("validate JPEG %q: %w", path, err)
	}
	if !isJPEG {
		return fmt.Errorf("input file %q is not a JPEG", path)
	}
	return nil
}

// getOutputFileName returns the output file name.
// The output file name is the same as the input file name, but with the extension changed to ".avif".
// For example, "input.jpg" -> "input.avif".
func getOutputFileName(inputFileName string) string {
	return filepath.Base(inputFileName[:len(inputFileName)-len(filepath.Ext(inputFileName))]) + ".avif"
}
