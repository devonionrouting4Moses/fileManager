//go:build windows
// +build windows

package utils

import (
	"os/exec"
	"syscall"
)

// ExecCommand executes a system command (Windows implementation)
// Runs the command with the window hidden
func ExecCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Start()
}
