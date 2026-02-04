package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/RayUI/RayUI/internal/model"
)

func TestNewGeoDataService(t *testing.T) {
	svc := NewGeoDataService("/tmp/test")
	if svc.DataDir != "/tmp/test" {
		t.Errorf("DataDir = %q, want /tmp/test", svc.DataDir)
	}
	if svc.client == nil {
		t.Error("client should not be nil")
	}
}

func TestGetGeoDataInfoNoFiles(t *testing.T) {
	dir := t.TempDir()
	svc := NewGeoDataService(dir)

	info := svc.GetGeoDataInfo(model.CoreXray)
	if info.GeoIPPath != "" {
		t.Errorf("GeoIPPath = %q, want empty", info.GeoIPPath)
	}
	if info.GeoSitePath != "" {
		t.Errorf("GeoSitePath = %q, want empty", info.GeoSitePath)
	}
	if info.LastUpdated != 0 {
		t.Errorf("LastUpdated = %d, want 0", info.LastUpdated)
	}
}

func TestGetGeoDataInfoWithFiles(t *testing.T) {
	dir := t.TempDir()
	dataPath := filepath.Join(dir, "data")
	os.MkdirAll(dataPath, 0o755)

	// Create fake geo data files for Xray
	os.WriteFile(filepath.Join(dataPath, "geoip.dat"), []byte("geoip data"), 0o644)
	os.WriteFile(filepath.Join(dataPath, "geosite.dat"), []byte("geosite data"), 0o644)

	// Create version file
	ts := time.Now().Unix()
	os.WriteFile(filepath.Join(dataPath, "geo.version"), []byte(fmt.Sprintf("%d", ts)), 0o644)

	svc := NewGeoDataService(dir)
	info := svc.GetGeoDataInfo(model.CoreXray)

	if info.GeoIPPath == "" {
		t.Error("GeoIPPath should not be empty")
	}
	if info.GeoSitePath == "" {
		t.Error("GeoSitePath should not be empty")
	}
	if info.GeoIPVersion == "" {
		t.Error("GeoIPVersion should not be empty")
	}
	if info.LastUpdated != ts {
		t.Errorf("LastUpdated = %d, want %d", info.LastUpdated, ts)
	}
}

func TestGetGeoDataInfoSingbox(t *testing.T) {
	dir := t.TempDir()
	dataPath := filepath.Join(dir, "data")
	os.MkdirAll(dataPath, 0o755)

	// Create sing-box format files (.db)
	os.WriteFile(filepath.Join(dataPath, "geoip.db"), []byte("geoip"), 0o644)
	os.WriteFile(filepath.Join(dataPath, "geosite.db"), []byte("geosite"), 0o644)

	svc := NewGeoDataService(dir)
	info := svc.GetGeoDataInfo(model.CoreSingbox)

	if info.GeoIPPath == "" {
		t.Error("GeoIPPath should not be empty for sing-box")
	}
	if info.GeoSitePath == "" {
		t.Error("GeoSitePath should not be empty for sing-box")
	}
}

func TestCheckGeoDataUpdateUnsupportedCore(t *testing.T) {
	svc := NewGeoDataService(t.TempDir())
	_, _, err := svc.CheckGeoDataUpdate(model.CoreAuto)
	if err == nil {
		t.Error("CheckGeoDataUpdate(CoreAuto) should return error")
	}
}

func TestCheckGeoDataUpdateNoLocalData(t *testing.T) {
	// Mock GitHub API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := struct {
			TagName     string `json:"tag_name"`
			PublishedAt string `json:"published_at"`
		}{
			TagName:     "v20240101",
			PublishedAt: time.Now().Format(time.RFC3339),
		}
		json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	dir := t.TempDir()
	svc := &GeoDataService{
		DataDir: dir,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	// Override by directly calling the API with our test server
	// Since CheckGeoDataUpdate builds the URL internally, we test the helper logic
	// by testing GetGeoDataInfo returns LastUpdated=0 for empty dir
	info := svc.GetGeoDataInfo(model.CoreXray)
	if info.LastUpdated != 0 {
		t.Errorf("LastUpdated should be 0 for empty dir, got %d", info.LastUpdated)
	}
}

func TestEnsureGeoDataUnsupportedCore(t *testing.T) {
	svc := NewGeoDataService(t.TempDir())
	err := svc.EnsureGeoData(model.CoreAuto)
	if err == nil {
		t.Error("EnsureGeoData(CoreAuto) should return error")
	}
}

func TestEnsureGeoDataExistingFiles(t *testing.T) {
	dir := t.TempDir()
	dataPath := filepath.Join(dir, "data")
	os.MkdirAll(dataPath, 0o755)

	// Pre-create the files so EnsureGeoData doesn't try to download
	os.WriteFile(filepath.Join(dataPath, "geoip.dat"), []byte("existing"), 0o644)
	os.WriteFile(filepath.Join(dataPath, "geosite.dat"), []byte("existing"), 0o644)

	svc := NewGeoDataService(dir)
	err := svc.EnsureGeoData(model.CoreXray)
	if err != nil {
		t.Fatalf("EnsureGeoData() with existing files should not error: %v", err)
	}
}

func TestUpdateGeoDataUnsupportedCore(t *testing.T) {
	svc := NewGeoDataService(t.TempDir())
	err := svc.UpdateGeoData(model.CoreAuto)
	if err == nil {
		t.Error("UpdateGeoData(CoreAuto) should return error")
	}
}

func TestUpdateGeoDataWithServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("fake geodata content"))
	}))
	defer server.Close()

	dir := t.TempDir()
	var progresses []GeoProgress
	svc := &GeoDataService{
		DataDir: dir,
		client:  &http.Client{Timeout: 5 * time.Second},
		OnProgress: func(p GeoProgress) {
			progresses = append(progresses, p)
		},
	}

	// We can't easily override the hardcoded URLs, but we can test downloadFile directly
	dataPath := filepath.Join(dir, "data")
	os.MkdirAll(dataPath, 0o755)

	destPath := filepath.Join(dataPath, "test.dat")
	err := svc.downloadFile(server.URL+"/test.dat", destPath, "test")
	if err != nil {
		t.Fatalf("downloadFile() error: %v", err)
	}

	data, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("read downloaded file: %v", err)
	}
	if string(data) != "fake geodata content" {
		t.Errorf("downloaded content = %q, want 'fake geodata content'", string(data))
	}

	// Verify progress was emitted
	if len(progresses) < 2 {
		t.Error("expected at least 2 progress events")
	}
	// Last should be "done"
	last := progresses[len(progresses)-1]
	if last.Status != "done" {
		t.Errorf("last progress status = %q, want 'done'", last.Status)
	}
}

func TestDownloadFileServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	dir := t.TempDir()
	svc := &GeoDataService{
		DataDir: dir,
		client:  &http.Client{Timeout: 5 * time.Second},
	}

	err := svc.downloadFile(server.URL, filepath.Join(dir, "out"), "test")
	if err == nil {
		t.Error("downloadFile() should return error on 404")
	}
}

func TestEmitGeoProgress(t *testing.T) {
	// nil callback should not panic
	svc := &GeoDataService{}
	svc.emitProgress("test", 0, 0, "downloading", "")

	// With callback
	var received GeoProgress
	svc.OnProgress = func(p GeoProgress) {
		received = p
	}
	svc.emitProgress("geoip", 100, 200, "downloading", "Downloading geoip...")
	if received.FileName != "geoip" {
		t.Errorf("FileName = %q, want geoip", received.FileName)
	}
	if received.Downloaded != 100 {
		t.Errorf("Downloaded = %d, want 100", received.Downloaded)
	}
	if received.Total != 200 {
		t.Errorf("Total = %d, want 200", received.Total)
	}
}

func TestGeoDataURLConstants(t *testing.T) {
	// Verify constants are not empty
	urls := []string{SingboxGeoIPURL, SingboxGeoSiteURL, XrayGeoIPURL, XrayGeoSiteURL}
	for _, u := range urls {
		if u == "" {
			t.Error("geo data URL constant should not be empty")
		}
	}
}
