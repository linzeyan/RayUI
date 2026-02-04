//go:build windows

package autostart

import (
	"golang.org/x/sys/windows/registry"
)

const regRunPath = `Software\Microsoft\Windows\CurrentVersion\Run`
const appName = "RayUI"

type windowsAutoStart struct{}

func newPlatformAutoStart() AutoStart {
	return &windowsAutoStart{}
}

func (w *windowsAutoStart) Enable(execPath string) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, regRunPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()
	return key.SetStringValue(appName, execPath)
}

func (w *windowsAutoStart) Disable() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, regRunPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()
	return key.DeleteValue(appName)
}

func (w *windowsAutoStart) IsEnabled() (bool, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, regRunPath, registry.QUERY_VALUE)
	if err != nil {
		return false, nil
	}
	defer key.Close()
	_, _, err = key.GetStringValue(appName)
	return err == nil, nil
}
