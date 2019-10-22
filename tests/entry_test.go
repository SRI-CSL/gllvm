package test

import (
	"fmt"
	"github.com/SRI-CSL/gllvm/shared"
	"testing"
)

func Test_basic(t *testing.T) {
	args := []string{"../data/helloworld.c", "-o", "../data/hello"}

	exitCode := shared.Compile(args, "clang")

	if exitCode != 0 {
		t.Errorf("Compile of %v returned %v\n", args, exitCode)
	} else {
		fmt.Println("Compiled OK")
	}
}
