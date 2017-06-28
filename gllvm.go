package main

import(
  "os"
  )

func main() {
    // Parse command line
    var args = os.Args
    var callerName = args[0]
    args = args[1:]

	logDebug("DEBUG: this is two: %v\n", 2)
	logInfo("INFO: this is two: %v\n", 2)

    switch callerName {
    case "gclang":
        compile(args, "clang")
    case "gclang++":
        compile(args, "clang++")
    case "get-bc":
        extract(args)
    default:
        logError("You should call gowllvm with a valid mode.")
    }
}
