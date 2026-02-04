//go:build linux

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
)

type linuxAutoStart struct{}

func newPlatformAutoStart() AutoStart {
	return &linuxAutoStart{}
}

func (l *linuxAutoStart) desktopPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "autostart", "rayui.desktop")
}

func (l *linuxAutoStart) Enable(execPath string) error {
	content := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=RayUI
Exec=%s
X-GNOME-Autostart-enabled=true
`, execPath)
	path := l.desktopPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func (l *linuxAutoStart) Disable() error {
	return os.Remove(l.desktopPath())
}

func (l *linuxAutoStart) IsEnabled() (bool, error) {
	_, err := os.Stat(l.desktopPath())
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
