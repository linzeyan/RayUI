package service

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/RayUI/RayUI/internal/model"
)

// UpdateInfo contains version information for a core.
type UpdateInfo struct {
	CoreType       model.ECoreType `json:"coreType"`
	CurrentVersion string          `json:"currentVersion"`
	LatestVersion  string          `json:"latestVersion"`
	HasUpdate      bool            `json:"hasUpdate"`
	DownloadURL    string          `json:"downloadUrl"`
	AssetName      string          `json:"assetName"`
}

// UpdateProgress is emitted during downloads.
type UpdateProgress struct {
	CoreType    model.ECoreType `json:"coreType"`
	Downloaded  int64           `json:"downloaded"`
	Total       int64           `json:"total"`
	Status      string          `json:"status"` // "downloading", "extracting", "done", "error"
	Description string          `json:"description,omitempty"`
}

// UpdaterService handles checking and downloading core updates from GitHub.
type UpdaterService struct {
	DataDir    string
	OnProgress func(UpdateProgress)
}

// GitHub API release structure (minimal fields).
type ghRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []ghAsset `json:"assets"`
}

type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int    `json:"size"`
}

const (
	singboxRepo = "SagerNet/sing-box"
	xrayRepo    = "XTLS/Xray-core"
)

// CheckUpdate checks GitHub for the latest release of the given core.
func (u *UpdaterService) CheckUpdate(coreType model.ECoreType) (*UpdateInfo, error) {
	repo := coreRepo(coreType)
	if repo == "" {
		return nil, fmt.Errorf("unknown core type: %v", coreType)
	}

	release, err := fetchLatestRelease(repo)
	if err != nil {
		return nil, fmt.Errorf("fetch release: %w", err)
	}

	osName, archName := platformNames()
	asset := findAsset(release.Assets, coreType, osName, archName)

	current := u.currentVersion(coreType)
	latest := strings.TrimPrefix(release.TagName, "v")

	info := &UpdateInfo{
		CoreType:       coreType,
		CurrentVersion: current,
		LatestVersion:  latest,
		HasUpdate:      current != latest && latest != "",
	}
	if asset != nil {
		info.DownloadURL = asset.BrowserDownloadURL
		info.AssetName = asset.Name
	}
	return info, nil
}

// DownloadUpdate downloads and installs the given core update.
func (u *UpdaterService) DownloadUpdate(info UpdateInfo) error {
	if info.DownloadURL == "" {
		return fmt.Errorf("no download URL for %v", info.CoreType)
	}

	u.emitProgress(UpdateProgress{
		CoreType: info.CoreType,
		Status:   "downloading",
	})

	// Download to temp file.
	resp, err := http.Get(info.DownloadURL)
	if err != nil {
		return fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned %d", resp.StatusCode)
	}

	tmpFile := filepath.Join(u.DataDir, "cores", "download.tmp")
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}

	total := resp.ContentLength
	written := int64(0)
	buf := make([]byte, 32*1024)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, wErr := f.Write(buf[:n]); wErr != nil {
				f.Close()
				os.Remove(tmpFile)
				return wErr
			}
			written += int64(n)
			u.emitProgress(UpdateProgress{
				CoreType:   info.CoreType,
				Downloaded: written,
				Total:      total,
				Status:     "downloading",
			})
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			f.Close()
			os.Remove(tmpFile)
			return readErr
		}
	}
	f.Close()

	// Extract binary.
	u.emitProgress(UpdateProgress{
		CoreType: info.CoreType,
		Status:   "extracting",
	})

	binaryName := coreBinaryName(info.CoreType)
	coresDir := filepath.Join(u.DataDir, "cores")
	destPath := filepath.Join(coresDir, binaryName)

	var extractErr error
	if strings.HasSuffix(info.AssetName, ".zip") {
		extractErr = extractFromZip(tmpFile, binaryName, destPath)
	} else if strings.HasSuffix(info.AssetName, ".tar.gz") || strings.HasSuffix(info.AssetName, ".tgz") {
		extractErr = extractFromTarGz(tmpFile, binaryName, destPath)
	} else {
		// Assume raw binary.
		extractErr = os.Rename(tmpFile, destPath)
	}
	os.Remove(tmpFile)

	if extractErr != nil {
		return fmt.Errorf("extract: %w", extractErr)
	}

	// Make executable.
	_ = os.Chmod(destPath, 0o755)

	// Write version file.
	versionFile := filepath.Join(coresDir, coreVersionFileName(info.CoreType))
	_ = os.WriteFile(versionFile, []byte(info.LatestVersion), 0o644)

	u.emitProgress(UpdateProgress{
		CoreType:    info.CoreType,
		Status:      "done",
		Description: info.LatestVersion,
	})
	return nil
}

func (u *UpdaterService) emitProgress(p UpdateProgress) {
	if u.OnProgress != nil {
		u.OnProgress(p)
	}
}

func (u *UpdaterService) currentVersion(coreType model.ECoreType) string {
	vFile := filepath.Join(u.DataDir, "cores", coreVersionFileName(coreType))
	data, err := os.ReadFile(vFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func coreRepo(ct model.ECoreType) string {
	switch ct {
	case model.CoreSingbox:
		return singboxRepo
	case model.CoreXray:
		return xrayRepo
	default:
		return ""
	}
}

func coreBinaryName(ct model.ECoreType) string {
	name := "xray"
	if ct == model.CoreSingbox {
		name = "sing-box"
	}
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return name
}

func coreVersionFileName(ct model.ECoreType) string {
	switch ct {
	case model.CoreSingbox:
		return "sing-box.version"
	case model.CoreXray:
		return "xray.version"
	default:
		return "core.version"
	}
}

func fetchLatestRelease(repo string) (*ghRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}
	return &release, nil
}

func platformNames() (string, string) {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Map Go arch names to common release asset names.
	switch arch {
	case "amd64":
		arch = "amd64"
	case "arm64":
		arch = "arm64"
	case "386":
		arch = "386"
	}

	return osName, arch
}

func findAsset(assets []ghAsset, coreType model.ECoreType, osName, archName string) *ghAsset {
	// Build list of OS name variants to try.
	// Xray uses "macos" instead of "darwin", most others use "darwin".
	osAliases := []string{osName}
	if osName == "darwin" {
		osAliases = append(osAliases, "macos")
	}

	// Build list of arch matchers to try.
	// Xray uses "64" for amd64, "arm64-v8a" for arm64, "32" for 386.
	// Short aliases like "64" could match "arm64", so we use delimited matching.
	type archMatcher struct {
		pattern   string
		delimited bool // if true, require "-" before pattern to avoid partial matches
	}
	archMatchers := []archMatcher{{pattern: archName}}
	switch archName {
	case "amd64":
		archMatchers = append(archMatchers, archMatcher{pattern: "64", delimited: true})
	case "arm64":
		archMatchers = append(archMatchers, archMatcher{pattern: "arm64-v8a"})
	case "386":
		archMatchers = append(archMatchers, archMatcher{pattern: "32", delimited: true})
	}

	for _, oa := range osAliases {
		for _, am := range archMatchers {
			for i := range assets {
				name := strings.ToLower(assets[i].Name)
				if !strings.Contains(name, oa) {
					continue
				}
				if am.delimited {
					// Match "-64." or "-64-" or end with "-64" to avoid matching "arm64"
					if !strings.Contains(name, "-"+am.pattern+".") &&
						!strings.Contains(name, "-"+am.pattern+"-") &&
						!strings.HasSuffix(name, "-"+am.pattern) {
						continue
					}
				} else {
					if !strings.Contains(name, am.pattern) {
						continue
					}
				}
				// Skip checksum / digest files.
				if strings.HasSuffix(name, ".sha256") || strings.HasSuffix(name, ".txt") || strings.HasSuffix(name, ".dgst") {
					continue
				}
				return &assets[i]
			}
		}
	}
	return nil
}

func extractFromZip(zipPath, binaryName, destPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if filepath.Base(f.Name) == binaryName {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			out, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer out.Close()

			_, err = io.Copy(out, rc)
			return err
		}
	}
	return fmt.Errorf("binary %q not found in zip", binaryName)
}

func extractFromTarGz(tarGzPath, binaryName, destPath string) error {
	f, err := os.Open(tarGzPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if filepath.Base(header.Name) == binaryName {
			out, err := os.Create(destPath)
			if err != nil {
				return err
			}
			defer out.Close()

			_, err = io.Copy(out, tr)
			return err
		}
	}
	return fmt.Errorf("binary %q not found in tar.gz", binaryName)
}
