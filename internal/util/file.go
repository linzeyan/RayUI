package util

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AppDataDir returns the path to ~/.RayUI, creating it if necessary.
func AppDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	dir := filepath.Join(home, ".RayUI")
	_ = os.MkdirAll(dir, 0o755)
	return dir
}

// EnsureDir creates the directory (and parents) if it does not exist.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// AtomicWriteJSON writes data as indented JSON to path using a temp file + rename.
func AtomicWriteJSON(path string, data interface{}) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}

	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}

	return os.Rename(tmp, path)
}

// ReadJSON reads a JSON file into target.
func ReadJSON(path string, target interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(target)
}
