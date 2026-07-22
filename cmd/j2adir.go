/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// j2adirCmd represents the j2adir command
var j2adirCmd = &cobra.Command{
	Use:   "j2adir",
	Short: "Convert a directory of JPEG files to AVIF",
	Long:  "Convert JPEG files in a directory and its subdirectories to AVIF.",
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
		inputDirPath, err := filepath.Abs(input)
		if err != nil {
			return fmt.Errorf("resolve input directory: %w", err)
		}
		if err := j2adir(inputDirPath, inputDirPath+"_avif", quality); err != nil {
			return err
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), "Successfully converted JPEG files to AVIF!")
		return err
	},
}

func init() {
	rootCmd.AddCommand(j2adirCmd)
	j2adirCmd.Flags().StringP("input", "i", "input", "Input Jpeg directory.")
	j2adirCmd.Flags().IntP("quality", "q", 30, "Quality of the output Avif file (0-100).")
}

// ディレクトリかチェックする関数
func isDir(path string) (bool, error) {
	// ファイルorディレクトリが存在するか確認
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	// ディレクトリか確認
	if info.IsDir() {
		return true, nil
	}

	return false, fmt.Errorf("%s is not a directory", path)
}

// ファイルがJPEGかどうかを判定する関数
func isJPEG(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// JPEGファイルのシグネチャを確認
	buf := make([]byte, 3)
	if _, err := io.ReadFull(file, buf); err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return false, nil
		}
		return false, err
	}
	if !bytes.Equal(buf, []byte{0xff, 0xd8, 0xff}) {
		return false, nil
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return false, err
	}
	if _, err := jpeg.DecodeConfig(file); err != nil {
		return false, fmt.Errorf("decode JPEG: %w", err)
	}
	return true, nil
}

// ディレクトリ内のJPEGファイルをavifに変換してoutputディレクトリに保存する関数
func j2adir(inputDirPath string, outputDirPath string, quality int) error {
	if err := validateQuality(quality); err != nil {
		return err
	}
	if isDirectory, err := isDir(inputDirPath); err != nil {
		return err
	} else if !isDirectory {
		return fmt.Errorf("input path %q is not a directory", inputDirPath)
	}
	if err := os.MkdirAll(outputDirPath, 0o755); err != nil {
		return fmt.Errorf("create output directory %q: %w", outputDirPath, err)
	}

	targets, err := collectConversionTargets(inputDirPath, outputDirPath)
	if err != nil {
		return err
	}
	for _, target := range targets {
		if err := j2a(target.sourcePath, target.outputPath, quality); err != nil {
			return err
		}
	}
	return nil
}

type conversionTarget struct {
	sourcePath string
	outputPath string
}

func collectConversionTargets(inputDirPath string, outputDirPath string) ([]conversionTarget, error) {
	outputSources := make(map[string]string)
	var targets []conversionTarget
	err := filepath.WalkDir(inputDirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		jpeg, err := isJPEG(path)
		if err != nil {
			return fmt.Errorf("check input file %q: %w", path, err)
		}
		if jpeg {
			outputFilePath := filepath.Join(outputDirPath, getOutputFileName(path))
			if sourcePath, exists := outputSources[outputFilePath]; exists {
				return fmt.Errorf("output name collision: %q and %q both map to %q", sourcePath, path, outputFilePath)
			}
			outputSources[outputFilePath] = path
			targets = append(targets, conversionTarget{sourcePath: path, outputPath: outputFilePath})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return targets, nil
}
