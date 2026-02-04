package util

import (
	"os"
	"path/filepath"
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
