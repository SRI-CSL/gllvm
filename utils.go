package main

import (
	"os"
	"os/exec"
)

// Executes a command then returns true if there was an error
func execCmd(cmdExecName string, args []string, workingDir string) bool {
	cmd := exec.Command(cmdExecName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = workingDir
	return cmd.Run() != nil
}
