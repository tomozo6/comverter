/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// j2adirCmd represents the j2adir command
var j2adirCmd = &cobra.Command{
	Use:   "j2adir",
	Short: "j2a directory",
	Long:  `j2a directory.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("j2adir called")
		inputDirPath, _ := filepath.Abs(cmd.Flag("input").Value.String())
		outputDirPath := inputDirPath + "_avif"

		fmt.Println(inputDirPath)
		fmt.Println(outputDirPath)

		// inputがディレクトリかチェック
		if isDir, err := isDir(inputDirPath); !isDir {
			fmt.Println(err)
			return
		}

		// output用ディレクトリを作成
		// すでに存在している場合は何もしない
		if _, err := os.Stat(outputDirPath); os.IsNotExist(err) {
			os.Mkdir(outputDirPath, 0755)
		}

		// ディレクトリ内のjpgファイルをavifに変換してoutputディレクトリに保存
		if err := j2adir(inputDirPath, outputDirPath, quality); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(j2adirCmd)
	j2adirCmd.Flags().StringP("input", "i", "input", "Input Jpeg directory.")
	j2adirCmd.Flags().IntVarP(&quality, "quality", "q", 30, "Quality of the output Avif file.")
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
func isJPEG(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	// JPEGファイルのシグネチャを確認
	buf := make([]byte, 3)
	_, err = file.Read(buf)
	return err == nil
}

// ディレクトリ内のJPEGファイルをavifに変換してoutputディレクトリに保存する関数
func j2adir(inputDirPath string, outputDirPath string, quality int) error {
	return filepath.WalkDir(inputDirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if isJPEG(path) {
			outputFilePath := outputDirPath + "/" + getOutputFileName(path)

			j2a(path, outputFilePath, quality)
		}
		return nil
	})
}
