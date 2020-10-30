package test

import (
	"fmt"
	"github.com/SRI-CSL/gllvm/shared"
	"testing"
	"runtime"
)

func Test_basic_functionality(t *testing.T) {
	args := []string{"../data/helloworld.c", "-o", "../data/hello"}
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

func Test_more_functionality(t *testing.T) {
	objectFile := "../data/bhello.o"
	args := []string{"../data/helloworld.c", "-c", "-o", objectFile}
	exitCode := shared.Compile(args, "clang")
	if exitCode != 0 {
		t.Errorf("Compile of %v returned %v\n", args, exitCode)
	} else {
		fmt.Println("Compiled OK")
	}
	ok, err := shared.IsObjectFileForOS(objectFile, runtime.GOOS)
	if !ok {
		t.Errorf("isObjectFileForOS(%v, %v) = %v (err = %v)\n", objectFile, runtime.GOOS, ok, err)
	} else {
		fmt.Printf("isObjectFileForOS(%v, %v) = %v\n", objectFile, runtime.GOOS, ok)
	}
	args = []string{objectFile, "-o", "../data/bhello"}
	exitCode = shared.Compile(args, "clang")
	if exitCode != 0 {
		t.Errorf("Compile of %v returned %v\n", args, exitCode)
	} else {
		fmt.Println("Compiled OK")
	}
	args = []string{"get-bc", "-v", "../data/bhello"}
	exitCode = shared.Extract(args)
	if exitCode != 0 {
		t.Errorf("Extraction of %v returned %v\n", args, exitCode)
	} else {
		fmt.Println("Extraction OK")
	}
}

func Test_obscure_functionality(t *testing.T) {
	objectFile := "../data/bhello.notanextensionthatwerecognize"
	args := []string{"../data/helloworld.c", "-c", "-o", objectFile}
	exitCode := shared.Compile(args, "clang")
	if exitCode != 0 {
		t.Errorf("Compile of %v returned %v\n", args, exitCode)
	} else {
		fmt.Println("Compiled OK")
	}
	ok, err := shared.IsObjectFileForOS(objectFile, runtime.GOOS)
	if !ok {
		t.Errorf("isObjectFileForOS(%v, %v) = %v (err = %v)\n", objectFile, runtime.GOOS, ok, err)
	} else {
		fmt.Printf("isObjectFileForOS(%v, %v) = %v\n", objectFile, runtime.GOOS, ok)
	}
	args = []string{objectFile, "-o", "../data/bhello"}
	exitCode = shared.Compile(args, "clang")
	if exitCode != 0 {
		t.Errorf("Compile of %v returned %v\n", args, exitCode)
	} else {
		fmt.Println("Compiled OK")
	}
	args = []string{"get-bc", "-v", "../data/bhello"}
	exitCode = shared.Extract(args)
	if exitCode != 0 {
		t.Errorf("Extraction of %v returned %v\n", args, exitCode)
	} else {
		fmt.Println("Extraction OK")
	}
}
