//go:build !windows

package core

import (
	"os"
	"syscall"
)

func coreSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}

func gracefulStop(p *os.Process) error {
	return p.Signal(syscall.SIGTERM)
}
