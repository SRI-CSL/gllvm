package shared

import (
	"os"
	"os/exec"
)

// Executes a command then returns true for success, false if there was an error, err is either nil or the error.
func execCmd(cmdExecName string, args []string, workingDir string) (success bool, err error) {
	cmd := exec.Command(cmdExecName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = workingDir
	err = cmd.Run()
	ecode := 0
	if err != nil {
		ecode = 1
	}
	LogDebug("execCmd: %v %v in %v had exitCode %v\n", cmdExecName, args, workingDir, ecode)
	if err != nil {
		LogDebug("execCmd: error was %v\n", err)
	}
	success = (ecode == 0)
	return
}
