package util

import (
	"math"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestEnsureDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "a", "b", "c")
	if err := EnsureDir(dir); err != nil {
		t.Fatalf("EnsureDir: %v", err)
	}
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected directory")
	}
}

func TestAtomicWriteAndReadJSON(t *testing.T) {
	type sample struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	path := filepath.Join(t.TempDir(), "test.json")

	want := sample{Name: "hello", Count: 42}
	if err := AtomicWriteJSON(path, want); err != nil {
		t.Fatalf("AtomicWriteJSON: %v", err)
	}

	// Verify the temp file was cleaned up.
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Error("expected temp file to be removed")
	}

	var got sample
	if err := ReadJSON(path, &got); err != nil {
		t.Fatalf("ReadJSON: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestReadJSONNotExist(t *testing.T) {
	var v struct{}
	err := ReadJSON(filepath.Join(t.TempDir(), "noexist.json"), &v)
	if !os.IsNotExist(err) {
		t.Fatalf("expected not-exist error, got %v", err)
	}
}

func TestGenerateUUIDUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := GenerateUUID()
		if id == "" {
			t.Fatal("GenerateUUID returned empty string")
		}
		if seen[id] {
			t.Fatalf("duplicate UUID: %s", id)
		}
		seen[id] = true
	}
}

func TestFormatSpeed(t *testing.T) {
	tests := []struct {
		input int64
		want  string
	}{
		{0, "0 B/s"},
		{1024, "1.00 KB/s"},
		{1048576, "1.00 MB/s"},
		{1073741824, "1.00 GB/s"},
	}
	for _, tt := range tests {
		got := FormatSpeed(tt.input)
		if got != tt.want {
			t.Errorf("FormatSpeed(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatBytesNegative(t *testing.T) {
	got := FormatBytes(-1)
	if got != "-1 B" {
		t.Errorf("FormatBytes(-1) = %q, want \"-1 B\"", got)
	}
}

func TestFormatBytesMaxInt64(t *testing.T) {
	got := FormatBytes(math.MaxInt64)
	if got == "" {
		t.Error("FormatBytes(MaxInt64) should not be empty")
	}
	// Should be in GB range.
	if len(got) < 3 {
		t.Errorf("FormatBytes(MaxInt64) = %q, too short", got)
	}
}

func TestEnsureDirNestedCreation(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "a", "b", "c", "d", "e")
	if err := EnsureDir(dir); err != nil {
		t.Fatal(err)
	}
	// Creating again should be a no-op.
	if err := EnsureDir(dir); err != nil {
		t.Fatal(err)
	}
}

func TestReadJSONInvalidContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte("not valid json{{{"), 0o644); err != nil {
		t.Fatal(err)
	}
	var v map[string]any
	if err := ReadJSON(path, &v); err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestGetOSAndArch(t *testing.T) {
	os := GetOS()
	if os != runtime.GOOS {
		t.Errorf("GetOS() = %q, want %q", os, runtime.GOOS)
	}
	arch := GetArch()
	if arch != runtime.GOARCH {
		t.Errorf("GetArch() = %q, want %q", arch, runtime.GOARCH)
	}
}

func TestAppDataDir(t *testing.T) {
	dir := AppDataDir()
	if dir == "" {
		t.Error("AppDataDir should not be empty")
	}
}
