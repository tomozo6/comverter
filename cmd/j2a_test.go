package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestJ2AConvertsJPEGToAVIF(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "converted.avif")
	if err := j2a(filepath.Join("..", "tests", "input.jpg"), outputPath, 80); err != nil {
		t.Fatalf("convert JPEG: %v", err)
	}

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

func TestJ2ARejectsNonJPEGAndInvalidQuality(t *testing.T) {
	directory := t.TempDir()
	nonJPEGPath := filepath.Join(directory, "input.txt")
	if err := os.WriteFile(nonJPEGPath, []byte("not a JPEG"), 0o644); err != nil {
		t.Fatalf("write non-JPEG fixture: %v", err)
	}

	if err := j2a(nonJPEGPath, filepath.Join(directory, "output.avif"), 30); err == nil {
		t.Fatal("j2a(non-JPEG) returned nil error")
	}
	if err := j2a(filepath.Join("..", "tests", "input.jpg"), filepath.Join(directory, "output.avif"), 101); err == nil {
		t.Fatal("j2a(quality 101) returned nil error")
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
