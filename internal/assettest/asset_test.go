// Package assettest provides tests for verifying embedded asset serving.
package assettest

import (
	"embed"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//go:embed testdata
var testAssets embed.FS

// TestEmbedFSStructure verifies the embedded filesystem structure.
func TestEmbedFSStructure(t *testing.T) {
	// List all files in the embedded FS
	var files []string
	err := fs.WalkDir(testAssets, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to walk embedded FS: %v", err)
	}
	t.Logf("Embedded files: %v", files)
}

// TestHTTPFileServer verifies http.FileServer serves files correctly.
func TestHTTPFileServer(t *testing.T) {
	// Create a sub-filesystem
	subFS, err := fs.Sub(testAssets, "testdata")
	if err != nil {
		t.Fatalf("Failed to create sub FS: %v", err)
	}

	// Create file server
	handler := http.FileServer(http.FS(subFS))
	server := httptest.NewServer(handler)
	defer server.Close()

	// Test requests
	tests := []struct {
		path        string
		wantStatus  int
		wantContain string
		wantType    string
	}{
		{"/index.html", 200, "<!DOCTYPE html>", "text/html"},
		{"/style.css", 200, "body", "text/css"},
		{"/script.js", 200, "console", "javascript"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			resp, err := http.Get(server.URL + tt.path)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			contentType := resp.Header.Get("Content-Type")
			if !strings.Contains(contentType, tt.wantType) {
				t.Errorf("Content-Type = %s, want to contain %s", contentType, tt.wantType)
			}
		})
	}
}
