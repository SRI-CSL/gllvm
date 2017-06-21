package main

import (
    "fmt"
    "os"
)

func compile(args []string) {
    if len(args) < 1 {
        fmt.Println("You must precise which compiler to use.")
        os.Exit(1)
    }
    var compilerName = args[0]
    args = args[1:]

    var pr = parse(args)
    var _ = pr
    var _ = compilerName
}
