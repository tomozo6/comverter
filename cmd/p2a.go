package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/spf13/cobra"
)

var p2aCmd = &cobra.Command{
	Use:   "p2a",
	Short: "Convert PDF pages to AVIF",
	Long:  "Convert every page in a PDF document to a separate AVIF file.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := cmd.Flags().GetString("input")
		if err != nil {
			return err
		}
		outputDir, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}
		quality, err := cmd.Flags().GetInt("quality")
		if err != nil {
			return err
		}
		density, err := cmd.Flags().GetInt("dpi")
		if err != nil {
			return err
		}

		if err := p2aWithDensity(input, outputDir, quality, density); err != nil {
			return err
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), "Successfully converted PDF pages to AVIF!")
		return err
	},
}

func init() {
	rootCmd.AddCommand(p2aCmd)
	p2aCmd.Flags().StringP("input", "i", "input.pdf", "Input PDF file.")
	p2aCmd.Flags().StringP("output", "o", ".", "Output directory for AVIF files.")
	p2aCmd.Flags().IntP("quality", "q", 30, "Quality of the output AVIF files (0-100).")
	p2aCmd.Flags().Int("dpi", 300, "PDF rasterization resolution in DPI.")
}

// p2a converts every page in a PDF document to a separate AVIF file.
func p2a(inputFile string, outputDir string, quality int) error {
	return p2aWithDensity(inputFile, outputDir, quality, 300)
}

// p2aWithDensity converts every PDF page at the supplied rasterization density.
func p2aWithDensity(inputFile string, outputDir string, quality int, density int) error {
	if err := validateQuality(quality); err != nil {
		return err
	}
	if density < 1 {
		return fmt.Errorf("dpi must be at least 1, got %d", density)
	}
	input, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("read PDF %q: %w", inputFile, err)
	}
	if !isPDF(input) {
		return fmt.Errorf("input file %q is not a PDF", inputFile)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory %q: %w", outputDir, err)
	}

	vipsStartupOnce.Do(func() {
		// PDF pages are processed one at a time. Do not retain completed page
		// operations in libvips' global cache during a long conversion.
		vips.Startup(&vips.Config{MaxCacheSize: 0, MaxCacheMem: 0})
	})
	// Loading every page merely to read metadata makes a high-DPI, long PDF
	// consume memory proportional to its page count. Loading the first page is
	// enough: pdfload still exposes the total page count in its metadata.
	params := vips.NewImportParams()
	params.NumPages.Set(1)
	document, err := vips.LoadImageFromBuffer(input, params)
	if err != nil {
		return fmt.Errorf("load PDF %q: %w", inputFile, err)
	}
	defer document.Close()

	pages := document.Pages()
	pageHeight := document.PageHeight()
	if pages < 1 || pageHeight < 1 {
		return fmt.Errorf("load PDF %q: no pages found", inputFile)
	}
	pageNumberWidth := max(3, len(strconv.Itoa(pages)))
	for page := 0; page < pages; page++ {
		pageParams := vips.NewImportParams()
		pageParams.Page.Set(page)
		pageParams.NumPages.Set(1)
		pageParams.Density.Set(density)
		image, err := vips.LoadImageFromBuffer(input, pageParams)
		if err != nil {
			return fmt.Errorf("load PDF page %d: %w", page+1, err)
		}
		options := vips.NewAvifExportParams()
		options.Quality = quality
		output, _, err := image.ExportAvif(options)
		image.Close()
		if err != nil {
			return fmt.Errorf("convert PDF page %d to AVIF: %w", page+1, err)
		}
		outputPath := filepath.Join(outputDir, getPDFPageOutputFileName(page+1, pageNumberWidth))
		if err := os.WriteFile(outputPath, output, 0o644); err != nil {
			return fmt.Errorf("write AVIF output %q: %w", outputPath, err)
		}
	}

	return nil
}

func isPDF(input []byte) bool {
	return len(input) >= len("%PDF-") && bytes.Equal(input[:len("%PDF-")], []byte("%PDF-"))
}

func getPDFPageOutputFileName(pageNumber int, pageNumberWidth int) string {
	return fmt.Sprintf("%0*d.avif", pageNumberWidth, pageNumber)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
