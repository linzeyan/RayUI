//go:build linux

package sysproxy

import (
	"fmt"
	"os/exec"
	"strconv"
)

type linuxSysProxy struct{}

func newPlatformSysProxy() SysProxy {
	return &linuxSysProxy{}
}

func (l *linuxSysProxy) Set(httpAddr string, httpPort int, socksAddr string, socksPort int) error {
	commands := [][]string{
		{"gsettings", "set", "org.gnome.system.proxy", "mode", "manual"},
		{"gsettings", "set", "org.gnome.system.proxy.http", "host", httpAddr},
		{"gsettings", "set", "org.gnome.system.proxy.http", "port", strconv.Itoa(httpPort)},
		{"gsettings", "set", "org.gnome.system.proxy.https", "host", httpAddr},
		{"gsettings", "set", "org.gnome.system.proxy.https", "port", strconv.Itoa(httpPort)},
		{"gsettings", "set", "org.gnome.system.proxy.socks", "host", socksAddr},
		{"gsettings", "set", "org.gnome.system.proxy.socks", "port", strconv.Itoa(socksPort)},
	}

	for _, args := range commands {
		if err := exec.Command(args[0], args[1:]...).Run(); err != nil {
			return fmt.Errorf("sysproxy set %v: %w", args, err)
		}
	}
	return nil
}

func (l *linuxSysProxy) Clear() error {
	return exec.Command("gsettings", "set", "org.gnome.system.proxy", "mode", "none").Run()
}

func (l *linuxSysProxy) GetCurrent() (*ProxyState, error) {
	out, err := exec.Command("gsettings", "get", "org.gnome.system.proxy", "mode").Output()
	if err != nil {
		return &ProxyState{}, nil
	}
	mode := string(out)
	return &ProxyState{
		Enabled: mode != "'none'\n" && mode != "'none'",
	}, nil
}
