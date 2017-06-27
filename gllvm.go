package main

import(
  "os"
  "log"
  )

func main() {
    // Parse command line
    var args = os.Args
    var callerName = args[0]
    args = args[1:]

    switch callerName {
    case "gclang":
        compile(args, "clang")
    case "gclang++":
        compile(args, "clang++")
    case "get-bc":
        extract(args)
    default:
        log.Fatal("You should call gowllvm with a valid mode.")
    }
}
