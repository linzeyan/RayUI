// Package main provides tests for the embedded frontend assets.
package main

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestFrontendAssetStructure verifies the embedded frontend structure.
func TestFrontendAssetStructure(t *testing.T) {
	// List all files in the embedded FS
	var files []string
	err := fs.WalkDir(assets, ".", func(path string, d fs.DirEntry, err error) error {
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

	t.Logf("Total embedded files: %d", len(files))
	for _, f := range files {
		t.Logf("  %s", f)
	}

	// Verify expected files exist
	expectedPaths := []string{
		"frontend/dist/index.html",
		"frontend/dist/static/css",
		"frontend/dist/static/js",
	}
	for _, expected := range expectedPaths {
		found := false
		for _, f := range files {
			if strings.HasPrefix(f, expected) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected path %s not found", expected)
		}
	}
}

// TestFrontendAssetServing verifies assets are served correctly.
func TestFrontendAssetServing(t *testing.T) {
	// Create sub-filesystem for frontend/dist
	subFS, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		t.Fatalf("Failed to create sub FS: %v", err)
	}

	// Create file server
	handler := http.FileServer(http.FS(subFS))
	server := httptest.NewServer(handler)
	defer server.Close()

	// Test index.html
	t.Run("index.html", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/index.html")
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			t.Errorf("Status = %d, want 200", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			t.Errorf("Content-Type = %s, want text/html", contentType)
		}
	})

	// Find and test CSS files
	t.Run("css files", func(t *testing.T) {
		var cssFiles []string
		fs.WalkDir(subFS, "static/css", func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() && strings.HasSuffix(path, ".css") {
				cssFiles = append(cssFiles, path)
			}
			return nil
		})

		if len(cssFiles) == 0 {
			t.Error("No CSS files found in static/css")
			return
		}

		for _, cssFile := range cssFiles {
			resp, err := http.Get(server.URL + "/" + cssFile)
			if err != nil {
				t.Errorf("Request for %s failed: %v", cssFile, err)
				continue
			}

			if resp.StatusCode != 200 {
				t.Errorf("%s: Status = %d, want 200", cssFile, resp.StatusCode)
			}

			contentType := resp.Header.Get("Content-Type")
			if !strings.Contains(contentType, "text/css") {
				t.Errorf("%s: Content-Type = %s, want text/css", cssFile, contentType)
			}
			resp.Body.Close()
		}
	})

	// Find and test JS files
	t.Run("js files", func(t *testing.T) {
		var jsFiles []string
		fs.WalkDir(subFS, "static/js", func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() && strings.HasSuffix(path, ".js") {
				jsFiles = append(jsFiles, path)
			}
			return nil
		})

		if len(jsFiles) == 0 {
			t.Error("No JS files found in static/js")
			return
		}

		for _, jsFile := range jsFiles {
			resp, err := http.Get(server.URL + "/" + jsFile)
			if err != nil {
				t.Errorf("Request for %s failed: %v", jsFile, err)
				continue
			}

			if resp.StatusCode != 200 {
				t.Errorf("%s: Status = %d, want 200", jsFile, resp.StatusCode)
			}
			resp.Body.Close()
		}
	})
}
