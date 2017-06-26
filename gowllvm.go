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
    case "gowclang":
        compile(args, "clang")
    case "gowclang++":
        compile(args, "clang++")
    case "gowextract":
        extract(args)
    default:
        log.Fatal("You should call gowllvm with a valid mode.")
    }
}
