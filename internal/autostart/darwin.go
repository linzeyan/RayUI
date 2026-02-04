//go:build darwin

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const plistLabel = "com.rayui.app"

type darwinAutoStart struct{}

func newPlatformAutoStart() AutoStart {
	return &darwinAutoStart{}
}

func (d *darwinAutoStart) plistPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library", "LaunchAgents", plistLabel+".plist")
}

func (d *darwinAutoStart) Enable(execPath string) error {
	plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>%s</string>
    <key>ProgramArguments</key>
    <array>
        <string>%s</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
</dict>
</plist>
`, plistLabel, execPath)

	path := d.plistPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(plist), 0o644)
}

func (d *darwinAutoStart) Disable() error {
	return os.Remove(d.plistPath())
}

func (d *darwinAutoStart) IsEnabled() (bool, error) {
	data, err := os.ReadFile(d.plistPath())
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return strings.Contains(string(data), plistLabel), nil
}
