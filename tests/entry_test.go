package test

import (
	"fmt"
	"testing"

	"github.com/SRI-CSL/gllvm/shared"
)

func Test_basic_functionality(t *testing.T) {
	args := []string{"../data/helloworld.c", "../data/testobj.c", "-o", "../data/hello"}

	exitCode := shared.Compile(args, "clang")

	if exitCode != 0 {
		t.Errorf("Compile of %v returned %v\n", args, exitCode)
	} else {
		fmt.Println("Compiled OK")
	}

	args = []string{"get-bc", "-v", "../data/hello"}

	exitCode = shared.Extract(args)

	if exitCode != 0 {
		t.Errorf("Extraction of %v returned %v\n", args, exitCode)
	} else {
		fmt.Println("Extraction OK")
	}

}
