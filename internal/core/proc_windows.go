//go:build windows

package core

import (
	"os"
	"syscall"
)

func coreSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}

func gracefulStop(p *os.Process) error {
	return p.Kill()
}
