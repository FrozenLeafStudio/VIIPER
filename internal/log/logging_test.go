package log

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenFileAppends(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.log")

	f, err := OpenFile(path)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	if _, err := f.WriteString("first\n"); err != nil {
		t.Fatalf("WriteString: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	f, err = OpenFile(path)
	if err != nil {
		t.Fatalf("OpenFile (reopen): %v", err)
	}
	if _, err := f.WriteString("second\n"); err != nil {
		t.Fatalf("WriteString: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if got, want := string(data), "first\nsecond\n"; got != want {
		t.Errorf("file content = %q, want %q", got, want)
	}
}

func TestOpenFileTruncatesWhenOversized(t *testing.T) {
	path := filepath.Join(t.TempDir(), "big.log")

	if err := os.WriteFile(path, make([]byte, maxAppendFileSize), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	f, err := OpenFile(path)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	if _, err := f.WriteString("fresh\n"); err != nil {
		t.Fatalf("WriteString: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if got, want := string(data), "fresh\n"; got != want {
		t.Errorf("file content = %q, want %q (oversized file was not truncated)", got, want)
	}
}

func TestOpenFileKeepsContentUnderCap(t *testing.T) {
	path := filepath.Join(t.TempDir(), "small.log")

	if err := os.WriteFile(path, []byte("existing\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	f, err := OpenFile(path)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	if _, err := f.WriteString("more\n"); err != nil {
		t.Fatalf("WriteString: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if got, want := string(data), "existing\nmore\n"; got != want {
		t.Errorf("file content = %q, want %q", got, want)
	}
}
