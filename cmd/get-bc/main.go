package main

import (
	"os"
	"github.com/SRI-CSL/gllvm/shared"
)

func main() {
	// Parse command line
	var args = os.Args

	shared.Extract(args)

	shared.LogInfo("Calling %v DID NOT TELL US WHAT HAPPENED\n", os.Args)

	// could be more honest about our success here
	os.Exit(0)

}
