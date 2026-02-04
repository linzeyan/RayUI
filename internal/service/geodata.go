package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/RayUI/RayUI/internal/model"
)

// GeoDataService manages GeoIP/GeoSite data files.
type GeoDataService struct {
	DataDir    string
	OnProgress func(GeoProgress)
	client     *http.Client
}

// GeoProgress reports download progress.
type GeoProgress struct {
	FileName    string `json:"fileName"`
	Downloaded  int64  `json:"downloaded"`
	Total       int64  `json:"total"`
	Status      string `json:"status"` // "downloading", "extracting", "done", "error"
	Description string `json:"description"`
}

// GeoDataInfo contains information about geo data files.
type GeoDataInfo struct {
	GeoIPVersion   string `json:"geoipVersion"`
	GeoSiteVersion string `json:"geositeVersion"`
	GeoIPPath      string `json:"geoipPath"`
	GeoSitePath    string `json:"geositePath"`
	LastUpdated    int64  `json:"lastUpdated"`
}

// NewGeoDataService creates a new GeoDataService.
func NewGeoDataService(dataDir string) *GeoDataService {
	return &GeoDataService{
		DataDir: dataDir,
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

// GitHub release URLs for geo data
const (
	// sing-box uses .srs format from sagernet/sing-geoip and sagernet/sing-geosite
	SingboxGeoIPURL   = "https://github.com/SagerNet/sing-geoip/releases/latest/download/geoip.db"
	SingboxGeoSiteURL = "https://github.com/SagerNet/sing-geosite/releases/latest/download/geosite.db"

	// xray uses .dat format from Loyalsoldier/v2ray-rules-dat
	XrayGeoIPURL   = "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat"
	XrayGeoSiteURL = "https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat"
)

// EnsureGeoData ensures geo data files exist, downloading if necessary.
func (s *GeoDataService) EnsureGeoData(coreType model.ECoreType) error {
	dataPath := filepath.Join(s.DataDir, "data")
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}

	var geoIPURL, geoSiteURL, geoIPFile, geoSiteFile string

	switch coreType {
	case model.CoreSingbox:
		geoIPURL = SingboxGeoIPURL
		geoSiteURL = SingboxGeoSiteURL
		geoIPFile = filepath.Join(dataPath, "geoip.db")
		geoSiteFile = filepath.Join(dataPath, "geosite.db")
	case model.CoreXray:
		geoIPURL = XrayGeoIPURL
		geoSiteURL = XrayGeoSiteURL
		geoIPFile = filepath.Join(dataPath, "geoip.dat")
		geoSiteFile = filepath.Join(dataPath, "geosite.dat")
	default:
		return fmt.Errorf("unsupported core type: %d", coreType)
	}

	// Check if files exist, download if not
	if _, err := os.Stat(geoIPFile); os.IsNotExist(err) {
		if err := s.downloadFile(geoIPURL, geoIPFile, "geoip"); err != nil {
			return fmt.Errorf("download geoip: %w", err)
		}
	}

	if _, err := os.Stat(geoSiteFile); os.IsNotExist(err) {
		if err := s.downloadFile(geoSiteURL, geoSiteFile, "geosite"); err != nil {
			return fmt.Errorf("download geosite: %w", err)
		}
	}

	return nil
}

// UpdateGeoData downloads the latest geo data files.
func (s *GeoDataService) UpdateGeoData(coreType model.ECoreType) error {
	dataPath := filepath.Join(s.DataDir, "data")
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}

	var geoIPURL, geoSiteURL, geoIPFile, geoSiteFile string

	switch coreType {
	case model.CoreSingbox:
		geoIPURL = SingboxGeoIPURL
		geoSiteURL = SingboxGeoSiteURL
		geoIPFile = filepath.Join(dataPath, "geoip.db")
		geoSiteFile = filepath.Join(dataPath, "geosite.db")
	case model.CoreXray:
		geoIPURL = XrayGeoIPURL
		geoSiteURL = XrayGeoSiteURL
		geoIPFile = filepath.Join(dataPath, "geoip.dat")
		geoSiteFile = filepath.Join(dataPath, "geosite.dat")
	default:
		return fmt.Errorf("unsupported core type: %d", coreType)
	}

	// Download both files
	if err := s.downloadFile(geoIPURL, geoIPFile, "geoip"); err != nil {
		return fmt.Errorf("download geoip: %w", err)
	}

	if err := s.downloadFile(geoSiteURL, geoSiteFile, "geosite"); err != nil {
		return fmt.Errorf("download geosite: %w", err)
	}

	// Update version file
	versionFile := filepath.Join(dataPath, "geo.version")
	versionData := fmt.Sprintf("%d", time.Now().Unix())
	_ = os.WriteFile(versionFile, []byte(versionData), 0644)

	return nil
}

// GetGeoDataInfo returns information about installed geo data files.
func (s *GeoDataService) GetGeoDataInfo(coreType model.ECoreType) GeoDataInfo {
	dataPath := filepath.Join(s.DataDir, "data")
	info := GeoDataInfo{}

	var geoIPFile, geoSiteFile string
	switch coreType {
	case model.CoreSingbox:
		geoIPFile = filepath.Join(dataPath, "geoip.db")
		geoSiteFile = filepath.Join(dataPath, "geosite.db")
	case model.CoreXray:
		geoIPFile = filepath.Join(dataPath, "geoip.dat")
		geoSiteFile = filepath.Join(dataPath, "geosite.dat")
	}

	if fi, err := os.Stat(geoIPFile); err == nil {
		info.GeoIPPath = geoIPFile
		info.GeoIPVersion = fi.ModTime().Format("2006-01-02")
	}

	if fi, err := os.Stat(geoSiteFile); err == nil {
		info.GeoSitePath = geoSiteFile
		info.GeoSiteVersion = fi.ModTime().Format("2006-01-02")
	}

	// Read last update time
	versionFile := filepath.Join(dataPath, "geo.version")
	if data, err := os.ReadFile(versionFile); err == nil {
		var ts int64
		fmt.Sscanf(string(data), "%d", &ts)
		info.LastUpdated = ts
	}

	return info
}

// CheckGeoDataUpdate checks if newer geo data is available via GitHub API.
func (s *GeoDataService) CheckGeoDataUpdate(coreType model.ECoreType) (bool, string, error) {
	var repo string
	switch coreType {
	case model.CoreSingbox:
		repo = "SagerNet/sing-geoip"
	case model.CoreXray:
		repo = "Loyalsoldier/v2ray-rules-dat"
	default:
		return false, "", fmt.Errorf("unsupported core type")
	}

	// Fetch latest release from GitHub API
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	resp, err := s.client.Get(apiURL)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("github api returned %d", resp.StatusCode)
	}

	var release struct {
		TagName     string `json:"tag_name"`
		PublishedAt string `json:"published_at"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false, "", err
	}

	// Compare with local version
	info := s.GetGeoDataInfo(coreType)
	if info.LastUpdated == 0 {
		// No local data, update needed
		return true, release.TagName, nil
	}

	// Parse published time
	pubTime, err := time.Parse(time.RFC3339, release.PublishedAt)
	if err != nil {
		return false, release.TagName, nil
	}

	// If remote is newer than local, update needed
	if pubTime.Unix() > info.LastUpdated {
		return true, release.TagName, nil
	}

	return false, release.TagName, nil
}

func (s *GeoDataService) downloadFile(url, destPath, name string) error {
	s.emitProgress(name, 0, 0, "downloading", fmt.Sprintf("Downloading %s...", name))

	resp, err := s.client.Get(url)
	if err != nil {
		s.emitProgress(name, 0, 0, "error", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("http %d", resp.StatusCode)
		s.emitProgress(name, 0, 0, "error", err.Error())
		return err
	}

	// Handle redirects for GitHub releases
	if strings.Contains(resp.Request.URL.Host, "objects.githubusercontent.com") ||
		strings.Contains(resp.Request.URL.Host, "github.com") {
		// Followed redirect, continue
	}

	total := resp.ContentLength

	// Create temp file
	tmpFile := destPath + ".tmp"
	f, err := os.Create(tmpFile)
	if err != nil {
		s.emitProgress(name, 0, 0, "error", err.Error())
		return err
	}

	// Download with progress
	var downloaded int64
	buf := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := f.Write(buf[:n]); writeErr != nil {
				f.Close()
				os.Remove(tmpFile)
				s.emitProgress(name, 0, 0, "error", writeErr.Error())
				return writeErr
			}
			downloaded += int64(n)
			s.emitProgress(name, downloaded, total, "downloading", "")
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			f.Close()
			os.Remove(tmpFile)
			s.emitProgress(name, 0, 0, "error", err.Error())
			return err
		}
	}
	f.Close()

	// Move temp file to final destination
	if err := os.Rename(tmpFile, destPath); err != nil {
		os.Remove(tmpFile)
		s.emitProgress(name, 0, 0, "error", err.Error())
		return err
	}

	s.emitProgress(name, downloaded, total, "done", "")
	return nil
}

func (s *GeoDataService) emitProgress(name string, downloaded, total int64, status, desc string) {
	if s.OnProgress != nil {
		s.OnProgress(GeoProgress{
			FileName:    name,
			Downloaded:  downloaded,
			Total:       total,
			Status:      status,
			Description: desc,
		})
	}
}
