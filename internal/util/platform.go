package util

import (
	"os"
	"os/exec"
	"runtime"
)

// GetOS returns the current operating system: "darwin", "windows", or "linux".
func GetOS() string {
	return runtime.GOOS
}

// GetArch returns the current architecture: "amd64", "arm64", etc.
func GetArch() string {
	return runtime.GOARCH
}

// IsAdmin checks whether the current process is running with elevated privileges.
func IsAdmin() bool {
	switch runtime.GOOS {
	case "windows":
		// net session requires admin; non-zero exit means not admin.
		err := exec.Command("net", "session").Run()
		return err == nil
	default:
		// Unix: UID 0 is root.
		return os.Getuid() == 0
	}
}
