package assettest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// distDir returns the absolute path to frontend/dist relative to this test file.
func distDir(t *testing.T) string {
	t.Helper()
	_, filename, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(filename), "..", "..")
	dir := filepath.Join(root, "frontend", "dist")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Skip("frontend/dist not found — run 'pnpm run build' in frontend/ first")
	}
	return dir
}

// TestDistIndexHTMLExists verifies the built index.html exists and contains
// references to CSS and JS files.
func TestDistIndexHTMLExists(t *testing.T) {
	dir := distDir(t)
	data, err := os.ReadFile(filepath.Join(dir, "index.html"))
	if err != nil {
		t.Fatalf("read index.html: %v", err)
	}
	html := string(data)

	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("index.html missing DOCTYPE")
	}
	if !strings.Contains(html, `rel="stylesheet"`) {
		t.Error("index.html has no stylesheet link")
	}
	if !strings.Contains(html, `<script defer src="/static/js/`) {
		t.Error("index.html has no script tag")
	}
}

// TestDistCSSLinksResolvable verifies every CSS link in index.html
// points to an actual file on disk.
func TestDistCSSLinksResolvable(t *testing.T) {
	dir := distDir(t)
	data, err := os.ReadFile(filepath.Join(dir, "index.html"))
	if err != nil {
		t.Fatalf("read index.html: %v", err)
	}

	re := regexp.MustCompile(`href="(/static/css/[^"]+)"`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	if len(matches) == 0 {
		t.Fatal("no CSS links found in index.html")
	}

	for _, m := range matches {
		href := m[1]
		filePath := filepath.Join(dir, href)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("CSS link %q does not resolve to a file", href)
		}
	}
}

// TestDistCSSCompiled verifies CSS output is properly compiled by Tailwind —
// no raw @apply, @tailwind, @custom-variant, or @theme{} directives remain.
func TestDistCSSCompiled(t *testing.T) {
	dir := distDir(t)
	cssDir := filepath.Join(dir, "static", "css")

	entries, err := os.ReadDir(cssDir)
	if err != nil {
		t.Fatalf("read css dir: %v", err)
	}

	var totalSize int64
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".css") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(cssDir, e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		css := string(data)
		totalSize += int64(len(data))

		// These Tailwind directives should be compiled away
		if strings.Contains(css, "@apply ") {
			t.Errorf("%s contains uncompiled @apply directive — is @tailwindcss/postcss configured?", e.Name())
		}
		if strings.Contains(css, "@tailwind ") {
			t.Errorf("%s contains uncompiled @tailwind directive", e.Name())
		}
		if strings.Contains(css, "@custom-variant ") {
			t.Errorf("%s contains uncompiled @custom-variant directive", e.Name())
		}
		if matched, _ := regexp.MatchString(`@theme\s*\{`, css); matched {
			t.Errorf("%s contains uncompiled @theme block", e.Name())
		}
	}

	// Compiled Tailwind CSS with utility classes should be > 10KB
	if totalSize < 10000 {
		t.Errorf("total CSS size = %d bytes, want > 10000 (CSS appears uncompiled)", totalSize)
	}
}

// TestDistCSSContainsUtilities verifies the compiled CSS contains actual
// Tailwind utility class definitions.
func TestDistCSSContainsUtilities(t *testing.T) {
	dir := distDir(t)
	cssDir := filepath.Join(dir, "static", "css")

	entries, err := os.ReadDir(cssDir)
	if err != nil {
		t.Fatalf("read css dir: %v", err)
	}

	var allCSS strings.Builder
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".css") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(cssDir, e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		allCSS.Write(data)
	}

	css := allCSS.String()
	required := []string{".flex", ".hidden", ".items-center", "--background", "--foreground", "--primary"}
	for _, cls := range required {
		if !strings.Contains(css, cls) {
			t.Errorf("compiled CSS missing expected content: %q", cls)
		}
	}
}

// TestDistServeCSS verifies that serving frontend/dist via http.FileServer
// returns CSS files with the correct Content-Type and content.
func TestDistServeCSS(t *testing.T) {
	dir := distDir(t)
	handler := http.FileServer(http.Dir(dir))
	server := httptest.NewServer(handler)
	defer server.Close()

	// Read index.html to find CSS links
	data, err := os.ReadFile(filepath.Join(dir, "index.html"))
	if err != nil {
		t.Fatalf("read index.html: %v", err)
	}

	re := regexp.MustCompile(`href="(/static/css/[^"]+)"`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	if len(matches) == 0 {
		t.Fatal("no CSS links found")
	}

	for _, m := range matches {
		href := m[1]
		t.Run(href, func(t *testing.T) {
			resp, err := http.Get(server.URL + href)
			if err != nil {
				t.Fatalf("GET %s: %v", href, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				t.Errorf("status = %d, want 200", resp.StatusCode)
			}

			ct := resp.Header.Get("Content-Type")
			if !strings.Contains(ct, "text/css") {
				t.Errorf("Content-Type = %q, want to contain 'text/css'", ct)
			}

			body, _ := io.ReadAll(resp.Body)
			if len(body) < 1000 {
				t.Errorf("CSS body too small (%d bytes), expected compiled CSS", len(body))
			}
			if strings.Contains(string(body), "@tailwind ") {
				t.Error("served CSS contains uncompiled @tailwind directive")
			}
		})
	}
}

// TestDistServeJS verifies JS files are served with the correct Content-Type.
func TestDistServeJS(t *testing.T) {
	dir := distDir(t)
	handler := http.FileServer(http.Dir(dir))
	server := httptest.NewServer(handler)
	defer server.Close()

	data, err := os.ReadFile(filepath.Join(dir, "index.html"))
	if err != nil {
		t.Fatalf("read index.html: %v", err)
	}

	re := regexp.MustCompile(`src="(/static/js/[^"]+)"`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	if len(matches) == 0 {
		t.Fatal("no JS links found")
	}

	for _, m := range matches {
		src := m[1]
		t.Run(src, func(t *testing.T) {
			resp, err := http.Get(server.URL + src)
			if err != nil {
				t.Fatalf("GET %s: %v", src, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				t.Errorf("status = %d, want 200", resp.StatusCode)
			}

			ct := resp.Header.Get("Content-Type")
			if !strings.Contains(ct, "javascript") {
				t.Errorf("Content-Type = %q, want to contain 'javascript'", ct)
			}
		})
	}
}
