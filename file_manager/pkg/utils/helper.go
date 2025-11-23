//go:build !windows
// +build !windows

package utils

import (
	"os/exec"
)

// ExecCommand executes a system command (Unix implementation)
func ExecCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Start()
}