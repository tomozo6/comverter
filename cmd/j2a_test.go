package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestJ2AConvertsJPEGToAVIF(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "converted.avif")
	j2a(filepath.Join("..", "tests", "input.jpg"), outputPath, 80)

	output, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read AVIF output: %v", err)
	}
	if len(output) == 0 {
		t.Fatal("AVIF output is empty")
	}
	if !bytes.Contains(output[:min(len(output), 64)], []byte("ftypavif")) {
		t.Fatal("output is not an AVIF file")
	}
}

func TestGetOutputFileName(t *testing.T) {
	got := getOutputFileName("photos/cat.photo.jpg")
	if want := "cat.photo.avif"; got != want {
		t.Errorf("getOutputFileName() = %q, want %q", got, want)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
