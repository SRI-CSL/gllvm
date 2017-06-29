package main

import(
	"os"
	"path"
)

func main() {
	// Parse command line
	var args = os.Args
	_, callerName := path.Split(args[0])
	args = args[1:]

	switch callerName {
	case "gclang":
		compile(args, "clang")
	case "gclang++":
		compile(args, "clang++")
	case "get-bc":
		extract(args)
	default:
		logError("You should call %s with a valid name.", callerName)
	}
}
