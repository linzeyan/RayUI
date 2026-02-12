package core

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RayUI/RayUI/internal/model"
)

func TestSelectCore(t *testing.T) {
	tests := []struct {
		name    string
		profile model.ProfileItem
		want    model.ECoreType
	}{
		{
			name:    "default → xray",
			profile: model.ProfileItem{ConfigType: model.ConfigVMess},
			want:    model.CoreXray,
		},
		{
			name:    "hysteria2 → singbox",
			profile: model.ProfileItem{ConfigType: model.ConfigHysteria2},
			want:    model.CoreSingbox,
		},
		{
			name:    "tuic → singbox",
			profile: model.ProfileItem{ConfigType: model.ConfigTUIC},
			want:    model.CoreSingbox,
		},
		{
			name:    "wireguard → singbox",
			profile: model.ProfileItem{ConfigType: model.ConfigWireGuard},
			want:    model.CoreSingbox,
		},
		{
			name:    "grpc → singbox",
			profile: model.ProfileItem{ConfigType: model.ConfigVLESS, Network: "grpc"},
			want:    model.CoreSingbox,
		},
		{
			name:    "h2 → singbox",
			profile: model.ProfileItem{ConfigType: model.ConfigVLESS, Network: "h2"},
			want:    model.CoreSingbox,
		},
		{
			name:    "override xray",
			profile: model.ProfileItem{ConfigType: model.ConfigVLESS, CoreType: model.CoreXray, Network: "grpc"},
			want:    model.CoreXray,
		},
		{
			name:    "override singbox",
			profile: model.ProfileItem{ConfigType: model.ConfigVMess, CoreType: model.CoreSingbox},
			want:    model.CoreSingbox,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SelectCore(tt.profile)
			if got != tt.want {
				t.Errorf("SelectCore = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCoreManagerTypes(t *testing.T) {
	dir := t.TempDir()

	xrayMgr := NewCoreManager(model.CoreXray, dir)
	if xrayMgr.CoreType() != model.CoreXray {
		t.Errorf("expected CoreXray, got %v", xrayMgr.CoreType())
	}

	singboxMgr := NewCoreManager(model.CoreSingbox, dir)
	if singboxMgr.CoreType() != model.CoreSingbox {
		t.Errorf("expected CoreSingbox, got %v", singboxMgr.CoreType())
	}

	// CoreAuto defaults to xray.
	autoMgr := NewCoreManager(model.CoreAuto, dir)
	if autoMgr.CoreType() != model.CoreXray {
		t.Errorf("expected CoreXray for auto, got %v", autoMgr.CoreType())
	}
}

func TestXrayCoreInitialState(t *testing.T) {
	dir := t.TempDir()
	xc := NewXrayCore(dir)

	if xc.IsRunning() {
		t.Error("new XrayCore should not be running")
	}

	status := xc.GetStatus()
	if status.Running {
		t.Error("status.Running should be false")
	}
	if status.CoreType != model.CoreXray {
		t.Errorf("CoreType = %v, want CoreXray", status.CoreType)
	}
	if status.PID != 0 {
		t.Errorf("PID should be 0, got %d", status.PID)
	}

	// Stop on a non-running core should not error.
	if err := xc.Stop(); err != nil {
		t.Errorf("Stop on non-running: %v", err)
	}
}

func TestSingboxCoreInitialState(t *testing.T) {
	dir := t.TempDir()
	sc := NewSingboxCore(dir)

	if sc.IsRunning() {
		t.Error("new SingboxCore should not be running")
	}

	status := sc.GetStatus()
	if status.Running {
		t.Error("status.Running should be false")
	}
	if status.CoreType != model.CoreSingbox {
		t.Errorf("CoreType = %v, want CoreSingbox", status.CoreType)
	}

	if err := sc.Stop(); err != nil {
		t.Errorf("Stop on non-running: %v", err)
	}
}

func TestCoreBinaryPaths(t *testing.T) {
	dir := "/test/data"
	xc := NewXrayCore(dir)
	if xc.BinaryPath() != filepath.Join(dir, "cores", "xray") {
		t.Errorf("xray binary path = %q", xc.BinaryPath())
	}

	sc := NewSingboxCore(dir)
	if sc.BinaryPath() != filepath.Join(dir, "cores", "sing-box") {
		t.Errorf("singbox binary path = %q", sc.BinaryPath())
	}
}

func TestCoreVersionFromFile(t *testing.T) {
	dir := t.TempDir()
	coresDir := filepath.Join(dir, "cores")
	if err := os.MkdirAll(coresDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Write a version file for xray.
	if err := os.WriteFile(filepath.Join(coresDir, "xray.version"), []byte("1.8.4\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	xc := NewXrayCore(dir)
	v, err := xc.Version()
	if err != nil {
		t.Fatalf("Version: %v", err)
	}
	if v != "1.8.4" {
		t.Errorf("version = %q, want 1.8.4", v)
	}

	// Write a version file for singbox.
	if err := os.WriteFile(filepath.Join(coresDir, "sing-box.version"), []byte("1.9.0\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	sc := NewSingboxCore(dir)
	v, err = sc.Version()
	if err != nil {
		t.Fatalf("Version: %v", err)
	}
	if v != "1.9.0" {
		t.Errorf("version = %q, want 1.9.0", v)
	}
}

func TestCoreVersionMissingFile(t *testing.T) {
	dir := t.TempDir()
	xc := NewXrayCore(dir)
	_, err := xc.Version()
	if err == nil {
		t.Error("expected error for missing version file")
	}
}
