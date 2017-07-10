package main

import (
	"os"
)

func main() {
	// Parse command line
	var args = os.Args
	args = args[1:]

	sanityCheck()

}
