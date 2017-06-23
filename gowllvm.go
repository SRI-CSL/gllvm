package main

import(
  "os"
  "os/exec"
  "log"
  )

func main() {
    // Parse command line
    var args = os.Args
    if len(args) < 2 {
        log.Fatal("Not enough arguments.")
    }
    var modeFlag = args[1]
    args = args[2:]

    switch modeFlag {
    case "compile":
        // Call main compiling function with args
        compile(args)
    case "extract":
        // Call main extracting function with args
        extract(args)
    default:
        log.Fatal("You should call gowllvm with a valid mode.")
    }
}

// Executes a command then returns true if there was an error
func execCmd(cmdExecName string, args []string) bool {
    cmd := exec.Command(cmdExecName, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if cmd.Run() == nil {
        return false
    } else {
        return true
    }
}
