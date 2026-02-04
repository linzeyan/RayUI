//go:build darwin

package sysproxy

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type darwinSysProxy struct{}

func newPlatformSysProxy() SysProxy {
	return &darwinSysProxy{}
}

func (d *darwinSysProxy) Set(httpAddr string, httpPort int, socksAddr string, socksPort int) error {
	service, err := getActiveNetworkService()
	if err != nil {
		return err
	}

	commands := [][]string{
		{"networksetup", "-setwebproxy", service, httpAddr, strconv.Itoa(httpPort)},
		{"networksetup", "-setsecurewebproxy", service, httpAddr, strconv.Itoa(httpPort)},
		{"networksetup", "-setsocksfirewallproxy", service, socksAddr, strconv.Itoa(socksPort)},
	}

	for _, args := range commands {
		if err := exec.Command(args[0], args[1:]...).Run(); err != nil {
			return fmt.Errorf("sysproxy set %v: %w", args, err)
		}
	}
	return nil
}

func (d *darwinSysProxy) Clear() error {
	service, err := getActiveNetworkService()
	if err != nil {
		return err
	}

	commands := [][]string{
		{"networksetup", "-setwebproxystate", service, "off"},
		{"networksetup", "-setsecurewebproxystate", service, "off"},
		{"networksetup", "-setsocksfirewallproxystate", service, "off"},
	}

	for _, args := range commands {
		if err := exec.Command(args[0], args[1:]...).Run(); err != nil {
			return fmt.Errorf("sysproxy clear %v: %w", args, err)
		}
	}
	return nil
}

func (d *darwinSysProxy) GetCurrent() (*ProxyState, error) {
	service, err := getActiveNetworkService()
	if err != nil {
		return nil, err
	}

	out, err := exec.Command("networksetup", "-getwebproxy", service).Output()
	if err != nil {
		return nil, err
	}

	state := &ProxyState{}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		switch key {
		case "Enabled":
			state.Enabled = val == "Yes"
		case "Server":
			state.HTTPHost = val
		case "Port":
			state.HTTPPort, _ = strconv.Atoi(val)
		}
	}
	return state, nil
}

// getActiveNetworkService detects the primary active network service.
func getActiveNetworkService() (string, error) {
	// Get the default route interface.
	out, err := exec.Command("route", "-n", "get", "default").Output()
	if err != nil {
		return "Wi-Fi", nil // fallback
	}

	var iface string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "interface:") {
			iface = strings.TrimSpace(strings.TrimPrefix(line, "interface:"))
			break
		}
	}
	if iface == "" {
		return "Wi-Fi", nil
	}

	// Map interface to network service name.
	listOut, err := exec.Command("networksetup", "-listallhardwareports").Output()
	if err != nil {
		return "Wi-Fi", nil
	}

	var currentService string
	for _, line := range strings.Split(string(listOut), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Hardware Port:") {
			currentService = strings.TrimPrefix(line, "Hardware Port: ")
		}
		if strings.HasPrefix(line, "Device:") {
			device := strings.TrimSpace(strings.TrimPrefix(line, "Device:"))
			if device == iface {
				return currentService, nil
			}
		}
	}

	return "Wi-Fi", nil
}
