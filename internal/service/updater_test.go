package service

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/RayUI/RayUI/internal/model"
)

func TestCoreRepo(t *testing.T) {
	tests := []struct {
		coreType model.ECoreType
		want     string
	}{
		{model.CoreSingbox, singboxRepo},
		{model.CoreXray, xrayRepo},
		{model.CoreAuto, ""},
		{model.ECoreType(99), ""},
	}
	for _, tt := range tests {
		got := coreRepo(tt.coreType)
		if got != tt.want {
			t.Errorf("coreRepo(%v) = %q, want %q", tt.coreType, got, tt.want)
		}
	}
}

func TestCoreBinaryName(t *testing.T) {
	xrayName := "xray"
	singboxName := "sing-box"
	if runtime.GOOS == "windows" {
		xrayName += ".exe"
		singboxName += ".exe"
	}

	if got := coreBinaryName(model.CoreXray); got != xrayName {
		t.Errorf("coreBinaryName(Xray) = %q, want %q", got, xrayName)
	}
	if got := coreBinaryName(model.CoreSingbox); got != singboxName {
		t.Errorf("coreBinaryName(Singbox) = %q, want %q", got, singboxName)
	}
}

func TestCoreVersionFileName(t *testing.T) {
	tests := []struct {
		coreType model.ECoreType
		want     string
	}{
		{model.CoreSingbox, "sing-box.version"},
		{model.CoreXray, "xray.version"},
		{model.CoreAuto, "core.version"},
	}
	for _, tt := range tests {
		got := coreVersionFileName(tt.coreType)
		if got != tt.want {
			t.Errorf("coreVersionFileName(%v) = %q, want %q", tt.coreType, got, tt.want)
		}
	}
}

func TestPlatformNames(t *testing.T) {
	osName, archName := platformNames()
	if osName != runtime.GOOS {
		t.Errorf("platformNames() osName = %q, want %q", osName, runtime.GOOS)
	}
	if archName != runtime.GOARCH {
		t.Errorf("platformNames() archName = %q, want %q", archName, runtime.GOARCH)
	}
}

func TestFindAsset(t *testing.T) {
	// Test with sing-box assets (uses "darwin" and standard arch names)
	singboxAssets := []ghAsset{
		{Name: "sing-box-1.12.19-darwin-arm64.tar.gz", BrowserDownloadURL: "https://example.com/sb-darwin-arm64.tar.gz"},
		{Name: "sing-box-1.12.19-darwin-amd64.tar.gz", BrowserDownloadURL: "https://example.com/sb-darwin-amd64.tar.gz"},
		{Name: "sing-box-1.12.19-linux-amd64.tar.gz", BrowserDownloadURL: "https://example.com/sb-linux-amd64.tar.gz"},
		{Name: "sing-box-1.12.19-windows-amd64.zip", BrowserDownloadURL: "https://example.com/sb-windows-amd64.zip"},
	}

	got := findAsset(singboxAssets, model.CoreSingbox, "darwin", "arm64")
	if got == nil {
		t.Fatal("findAsset(singbox, darwin/arm64) returned nil")
	}
	if got.Name != "sing-box-1.12.19-darwin-arm64.tar.gz" {
		t.Errorf("findAsset(singbox) = %q, want sing-box-1.12.19-darwin-arm64.tar.gz", got.Name)
	}

	got = findAsset(singboxAssets, model.CoreSingbox, "darwin", "amd64")
	if got == nil {
		t.Fatal("findAsset(singbox, darwin/amd64) returned nil")
	}

	// Test with real Xray-core assets (uses "macos" not "darwin", "64" not "amd64")
	xrayAssets := []ghAsset{
		{Name: "Xray-macos-arm64-v8a.zip", BrowserDownloadURL: "https://example.com/xray-macos-arm64.zip"},
		{Name: "Xray-macos-arm64-v8a.zip.dgst", BrowserDownloadURL: "https://example.com/xray-macos-arm64.dgst"},
		{Name: "Xray-macos-64.zip", BrowserDownloadURL: "https://example.com/xray-macos-64.zip"},
		{Name: "Xray-macos-64.zip.dgst", BrowserDownloadURL: "https://example.com/xray-macos-64.dgst"},
		{Name: "Xray-linux-64.zip", BrowserDownloadURL: "https://example.com/xray-linux-64.zip"},
		{Name: "Xray-linux-64.zip.dgst", BrowserDownloadURL: "https://example.com/xray-linux-64.dgst"},
		{Name: "Xray-linux-arm64-v8a.zip", BrowserDownloadURL: "https://example.com/xray-linux-arm64.zip"},
		{Name: "Xray-windows-64.zip", BrowserDownloadURL: "https://example.com/xray-windows-64.zip"},
		{Name: "Xray-windows-arm64-v8a.zip", BrowserDownloadURL: "https://example.com/xray-windows-arm64.zip"},
	}

	// Xray on macOS/arm64: "darwin" should match via "macos" alias, "arm64" via "arm64-v8a"
	got = findAsset(xrayAssets, model.CoreXray, "darwin", "arm64")
	if got == nil {
		t.Fatal("findAsset(xray, darwin/arm64) returned nil — should match Xray-macos-arm64-v8a.zip")
	}
	if got.Name != "Xray-macos-arm64-v8a.zip" {
		t.Errorf("findAsset(xray, darwin/arm64) = %q, want Xray-macos-arm64-v8a.zip", got.Name)
	}

	// Xray on macOS/amd64: "darwin" → "macos", "amd64" → "64"
	got = findAsset(xrayAssets, model.CoreXray, "darwin", "amd64")
	if got == nil {
		t.Fatal("findAsset(xray, darwin/amd64) returned nil — should match Xray-macos-64.zip")
	}
	if got.Name != "Xray-macos-64.zip" {
		t.Errorf("findAsset(xray, darwin/amd64) = %q, want Xray-macos-64.zip", got.Name)
	}

	// Xray on linux/amd64
	got = findAsset(xrayAssets, model.CoreXray, "linux", "amd64")
	if got == nil {
		t.Fatal("findAsset(xray, linux/amd64) returned nil")
	}
	if got.Name != "Xray-linux-64.zip" {
		t.Errorf("findAsset(xray, linux/amd64) = %q, want Xray-linux-64.zip", got.Name)
	}

	// Xray on linux/arm64
	got = findAsset(xrayAssets, model.CoreXray, "linux", "arm64")
	if got == nil {
		t.Fatal("findAsset(xray, linux/arm64) returned nil")
	}

	// Xray on windows/amd64
	got = findAsset(xrayAssets, model.CoreXray, "windows", "amd64")
	if got == nil {
		t.Fatal("findAsset(xray, windows/amd64) returned nil")
	}

	// No match returns nil
	got = findAsset(xrayAssets, model.CoreXray, "freebsd", "amd64")
	if got != nil {
		t.Errorf("findAsset() expected nil for freebsd, got %q", got.Name)
	}

	// Checksum files (.sha256, .txt, .dgst) should be skipped
	checksumOnly := []ghAsset{
		{Name: "Xray-macos-arm64-v8a.zip.dgst"},
		{Name: "checksums.txt"},
		{Name: "Xray-macos-arm64-v8a.sha256"},
	}
	got = findAsset(checksumOnly, model.CoreXray, "darwin", "arm64")
	if got != nil {
		t.Errorf("findAsset() should skip checksum files, got %q", got.Name)
	}
}

func TestCurrentVersion(t *testing.T) {
	dir := t.TempDir()
	u := &UpdaterService{DataDir: dir}

	// No version file returns empty string
	got := u.currentVersion(model.CoreXray)
	if got != "" {
		t.Errorf("currentVersion() = %q, want empty", got)
	}

	// Write a version file
	coresDir := filepath.Join(dir, "cores")
	os.MkdirAll(coresDir, 0o755)
	os.WriteFile(filepath.Join(coresDir, "xray.version"), []byte("1.8.4\n"), 0o644)

	got = u.currentVersion(model.CoreXray)
	if got != "1.8.4" {
		t.Errorf("currentVersion() = %q, want 1.8.4", got)
	}
}

func TestEmitProgress(t *testing.T) {
	// With nil callback - should not panic
	u := &UpdaterService{}
	u.emitProgress(UpdateProgress{Status: "test"})

	// With callback
	var received UpdateProgress
	u.OnProgress = func(p UpdateProgress) {
		received = p
	}
	u.emitProgress(UpdateProgress{Status: "downloading", CoreType: model.CoreXray})
	if received.Status != "downloading" {
		t.Errorf("emitProgress() status = %q, want downloading", received.Status)
	}
}

func TestCheckUpdate(t *testing.T) {
	// Create a mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := ghRelease{
			TagName: "v1.9.0",
			Assets: []ghAsset{
				{Name: "xray-darwin-arm64.zip", BrowserDownloadURL: "https://example.com/xray.zip", Size: 1000},
				{Name: "xray-linux-amd64.zip", BrowserDownloadURL: "https://example.com/xray-linux.zip", Size: 2000},
			},
		}
		json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	// Override fetchLatestRelease by testing via CheckUpdate with a custom server
	// Since fetchLatestRelease uses http.DefaultClient, we use the server URL directly
	// We'll test the pure logic parts instead

	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "cores"), 0o755)
	os.WriteFile(filepath.Join(dir, "cores", "xray.version"), []byte("1.8.0"), 0o644)

	u := &UpdaterService{DataDir: dir}

	// Test unknown core type
	_, err := u.CheckUpdate(model.CoreAuto)
	if err == nil {
		t.Error("CheckUpdate(CoreAuto) should return error")
	}
}

func TestExtractFromZip(t *testing.T) {
	dir := t.TempDir()

	// Create a test zip file containing a binary
	zipPath := filepath.Join(dir, "test.zip")
	destPath := filepath.Join(dir, "xray")

	f, err := os.Create(zipPath)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	fw, err := w.Create("xray")
	if err != nil {
		t.Fatal(err)
	}
	fw.Write([]byte("fake binary content"))
	w.Close()
	f.Close()

	// Extract
	err = extractFromZip(zipPath, "xray", destPath)
	if err != nil {
		t.Fatalf("extractFromZip() error: %v", err)
	}

	// Verify extracted file
	data, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("read extracted file: %v", err)
	}
	if string(data) != "fake binary content" {
		t.Errorf("extracted content = %q, want 'fake binary content'", string(data))
	}

	// Test missing binary in zip
	err = extractFromZip(zipPath, "nonexistent", filepath.Join(dir, "out"))
	if err == nil {
		t.Error("extractFromZip() should error when binary not found")
	}
}

func TestExtractFromTarGz(t *testing.T) {
	dir := t.TempDir()

	// Create a test tar.gz file
	tarGzPath := filepath.Join(dir, "test.tar.gz")
	destPath := filepath.Join(dir, "sing-box")

	f, err := os.Create(tarGzPath)
	if err != nil {
		t.Fatal(err)
	}
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)

	content := []byte("fake sing-box binary")
	hdr := &tar.Header{
		Name: "sing-box-1.0/sing-box",
		Mode: 0o755,
		Size: int64(len(content)),
	}
	tw.WriteHeader(hdr)
	tw.Write(content)
	tw.Close()
	gw.Close()
	f.Close()

	// Extract
	err = extractFromTarGz(tarGzPath, "sing-box", destPath)
	if err != nil {
		t.Fatalf("extractFromTarGz() error: %v", err)
	}

	data, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("read extracted file: %v", err)
	}
	if string(data) != "fake sing-box binary" {
		t.Errorf("extracted content = %q, want 'fake sing-box binary'", string(data))
	}

	// Test missing binary
	err = extractFromTarGz(tarGzPath, "nonexistent", filepath.Join(dir, "out"))
	if err == nil {
		t.Error("extractFromTarGz() should error when binary not found")
	}
}

func TestDownloadUpdate(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "cores"), 0o755)

	// Create a mock server that serves a zip file
	binaryContent := []byte("xray-binary-content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve a zip with the binary
		zw := zip.NewWriter(w)
		fw, _ := zw.Create("xray")
		fw.Write(binaryContent)
		zw.Close()
	}))
	defer server.Close()

	var progressStatuses []string
	u := &UpdaterService{
		DataDir: dir,
		OnProgress: func(p UpdateProgress) {
			progressStatuses = append(progressStatuses, p.Status)
		},
	}

	info := UpdateInfo{
		CoreType:       model.CoreXray,
		LatestVersion:  "1.9.0",
		DownloadURL:    server.URL + "/xray.zip",
		AssetName:      "xray.zip",
	}

	err := u.DownloadUpdate(info)
	if err != nil {
		t.Fatalf("DownloadUpdate() error: %v", err)
	}

	// Verify the binary was extracted
	extracted, err := os.ReadFile(filepath.Join(dir, "cores", "xray"))
	if err != nil {
		t.Fatalf("read extracted binary: %v", err)
	}
	if string(extracted) != "xray-binary-content" {
		t.Errorf("extracted binary = %q, want 'xray-binary-content'", string(extracted))
	}

	// Verify version file was written
	version, err := os.ReadFile(filepath.Join(dir, "cores", "xray.version"))
	if err != nil {
		t.Fatalf("read version file: %v", err)
	}
	if string(version) != "1.9.0" {
		t.Errorf("version = %q, want 1.9.0", string(version))
	}

	// Verify progress was emitted
	if len(progressStatuses) < 2 {
		t.Error("expected at least 2 progress events")
	}
	// Last should be "done"
	if progressStatuses[len(progressStatuses)-1] != "done" {
		t.Errorf("last progress status = %q, want 'done'", progressStatuses[len(progressStatuses)-1])
	}
}

func TestDownloadUpdateNoURL(t *testing.T) {
	u := &UpdaterService{}
	err := u.DownloadUpdate(UpdateInfo{CoreType: model.CoreXray})
	if err == nil {
		t.Error("DownloadUpdate() with empty URL should return error")
	}
}
