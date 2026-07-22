package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestJ2ADirConvertsJPEGFilesRecursively(t *testing.T) {
	inputDir := t.TempDir()
	outputDir := filepath.Join(t.TempDir(), "output")
	if err := os.Mkdir(outputDir, 0o755); err != nil {
		t.Fatalf("create output directory: %v", err)
	}

	copyFixture(t, filepath.Join(inputDir, "first.jpg"))
	if err := os.WriteFile(filepath.Join(inputDir, "not-an-image.txt"), []byte("not a JPEG"), 0o644); err != nil {
		t.Fatalf("write non-JPEG input: %v", err)
	}
	nestedDir := filepath.Join(inputDir, "nested")
	if err := os.Mkdir(nestedDir, 0o755); err != nil {
		t.Fatalf("create nested input directory: %v", err)
	}
	copyFixture(t, filepath.Join(nestedDir, "second.jpg"))

	if err := j2adir(inputDir, outputDir, 30); err != nil {
		t.Fatalf("convert directory: %v", err)
	}

	for _, name := range []string{"first.avif", "second.avif"} {
		output, err := os.ReadFile(filepath.Join(outputDir, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if len(output) == 0 || !bytes.Contains(output[:min(len(output), 64)], []byte("ftypavif")) {
			t.Errorf("%s is not a valid AVIF output", name)
		}
	}
}

func TestIsDir(t *testing.T) {
	directory := t.TempDir()
	if got, err := isDir(directory); err != nil || !got {
		t.Errorf("isDir(directory) = (%v, %v), want (true, nil)", got, err)
	}

	file := filepath.Join(directory, "file.jpg")
	copyFixture(t, file)
	if got, err := isDir(file); err == nil || got {
		t.Errorf("isDir(file) = (%v, %v), want (false, error)", got, err)
	}
}

func TestIsJPEG(t *testing.T) {
	directory := t.TempDir()
	jpegPath := filepath.Join(directory, "input.jpg")
	copyFixture(t, jpegPath)
	if !isJPEG(jpegPath) {
		t.Fatal("isJPEG(JPEG fixture) = false, want true")
	}

	nonJPEGPath := filepath.Join(directory, "input.txt")
	if err := os.WriteFile(nonJPEGPath, []byte("not a JPEG"), 0o644); err != nil {
		t.Fatalf("write non-JPEG fixture: %v", err)
	}
	if isJPEG(nonJPEGPath) {
		t.Fatal("isJPEG(non-JPEG fixture) = true, want false")
	}
}

func copyFixture(t *testing.T, destination string) {
	t.Helper()

	input, err := os.ReadFile(filepath.Join("..", "tests", "input.jpg"))
	if err != nil {
		t.Fatalf("read JPEG fixture: %v", err)
	}
	if err := os.WriteFile(destination, input, 0o644); err != nil {
		t.Fatalf("write JPEG fixture: %v", err)
	}
}
