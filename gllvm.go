package main

import (
	"os"
	"path"
)

func main() {
	// Parse command line
	var args = os.Args
	_, callerName := path.Split(args[0])
	args = args[1:]
	
	var exitCode int

	switch callerName {
	case "gclang":
		exitCode = compile(args, "clang")
	case "gclang++":
		exitCode = compile(args, "clang++")
	case "get-bc":
		extract(args)
	default:
		logError("You should call %s with a valid name.", callerName)
	}

	logInfo("Calling %v returned %v\n",  os.Args, exitCode)

	//important to pretend to look like the actual wrapped command
	os.Exit(exitCode)

}
