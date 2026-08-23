package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestP2ARejectsNonPDF(t *testing.T) {
	inputPath := filepath.Join(t.TempDir(), "input.pdf")
	if err := os.WriteFile(inputPath, []byte("not a PDF"), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	if err := p2a(inputPath, filepath.Join(t.TempDir(), "output"), 30); err == nil {
		t.Fatal("p2a(non-PDF) returned nil error")
	}
}

func TestP2ARejectsInvalidDPI(t *testing.T) {
	if err := p2aWithDensity("unused.pdf", t.TempDir(), 30, 0); err == nil {
		t.Fatal("p2aWithDensity(dpi 0) returned nil error")
	}
}

func TestGetPDFPageOutputFileName(t *testing.T) {
	if got, want := getPDFPageOutputFileName("documents/sample.pdf", 1, 3), "sample-001.avif"; got != want {
		t.Errorf("getPDFPageOutputFileName() = %q, want %q", got, want)
	}
	if got, want := getPDFPageOutputFileName("documents/book.pdf", 1000, 4), "book-1000.avif"; got != want {
		t.Errorf("getPDFPageOutputFileName() = %q, want %q", got, want)
	}
}
