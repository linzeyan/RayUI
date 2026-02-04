package netmgr

import (
	"net"
	"runtime"

	"github.com/RayUI/RayUI/internal/util"
)

// CheckTUNPermission checks if the process has permission to create TUN devices.
func CheckTUNPermission() error {
	if runtime.GOOS == "windows" {
		if !util.IsAdmin() {
			return &ErrPermission{Msg: "administrator privileges required for TUN mode"}
		}
	} else {
		if !util.IsAdmin() {
			return &ErrPermission{Msg: "root privileges required for TUN mode"}
		}
	}
	return nil
}

// GetNetworkInterfaces returns the names of all non-loopback network interfaces.
func GetNetworkInterfaces() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	var names []string
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		names = append(names, iface.Name)
	}
	return names
}

// ErrPermission is returned when elevated privileges are required.
type ErrPermission struct {
	Msg string
}

func (e *ErrPermission) Error() string {
	return e.Msg
}
