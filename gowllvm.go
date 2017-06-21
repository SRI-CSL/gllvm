package main

import(
  "fmt"
  "os"
  )

func main() {
    // Parse command line
    var args = os.Args
    if len(args) < 2 {
        fmt.Println("Not enough arguments.")
        os.Exit(1)
    }
    var modeFlag = args[1]
    args = args[2:]

    switch modeFlag {
    case "compile":
        // Call main compile function with args
        compile(args)
    case "extract":
        fmt.Println("The extract feature is not implemented yet.")
        os.Exit(1)
    default:
        fmt.Println("You should call gowllvm with a valid mode.")
        os.Exit(1)
    }
}
