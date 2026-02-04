//go:build windows

package sysproxy

import (
	"fmt"
	"strconv"

	"golang.org/x/sys/windows/registry"
)

type windowsSysProxy struct{}

func newPlatformSysProxy() SysProxy {
	return &windowsSysProxy{}
}

const regPath = `Software\Microsoft\Windows\CurrentVersion\Internet Settings`

func (w *windowsSysProxy) Set(httpAddr string, httpPort int, socksAddr string, socksPort int) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, regPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("open registry: %w", err)
	}
	defer key.Close()

	proxyServer := fmt.Sprintf("socks=%s:%d;http=%s:%d;https=%s:%d",
		socksAddr, socksPort, httpAddr, httpPort, httpAddr, httpPort)
	proxyOverride := "localhost;127.*;10.*;192.168.*"

	if err := key.SetDWordValue("ProxyEnable", 1); err != nil {
		return err
	}
	if err := key.SetStringValue("ProxyServer", proxyServer); err != nil {
		return err
	}
	if err := key.SetStringValue("ProxyOverride", proxyOverride); err != nil {
		return err
	}
	return nil
}

func (w *windowsSysProxy) Clear() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, regPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("open registry: %w", err)
	}
	defer key.Close()

	return key.SetDWordValue("ProxyEnable", 0)
}

func (w *windowsSysProxy) GetCurrent() (*ProxyState, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, regPath, registry.QUERY_VALUE)
	if err != nil {
		return nil, err
	}
	defer key.Close()

	enabled, _, _ := key.GetIntegerValue("ProxyEnable")
	server, _, _ := key.GetStringValue("ProxyServer")

	state := &ProxyState{
		Enabled:  enabled == 1,
		HTTPHost: server,
	}

	// Parse port from server string if possible.
	_ = strconv.Itoa(0) // import used
	return state, nil
}
