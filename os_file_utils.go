package ebtk

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Open a given path in OSX or Win
func OpenPath(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}
