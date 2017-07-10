package main

import (
	"bytes"
	"github.com/SRI-CSL/gllvm/shared"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const explainLLVMCCNAME = `

If your clang compiler is not called clang, but something else, then
you will need to set the environment variable LLVM_CC_NAME to the
appropriate string. For example if your clang is called clang-3.5 then
LLVM_CC_NAME should be set to clang-3.5.

`

const explainLLVMCXXNAME = `

If your clang++ compiler is not called clang++, but something else,
then you will need to set the environment variable LLVM_CXX_NAME to
the appropriate string. For example if your clang++ is called ++clang
then LLVM_CC_NAME should be set to ++clang.

`

const explainLLVMCOMPILERPATH = `

Your compiler should either be in your PATH, or else located where the
environment variable LLVM_COMPILER_PATH indicates. It can also be used
to indicate the directory that contains the other LLVM tools such as
llvm-link, and llvm-ar.

`

const explainLLVMLINKNAME = `

If your llvm linker is not called llvm-link, but something else, then
you will need to set the environment variable LLVM_LINK_NAME to the
appropriate string. For example if your llvm-link is called llvm-link-3.5 then
LLVM_LINK_NAME should be set to llvm-link-3.5.

`

const explainLLVMARNAME = `

If your llvm archiver is not called llvm-ar, but something else,
then you will need to set the environment variable LLVM_AR_NAME to
the appropriate string. For example if your llvm-ar is called llvm-ar-3.5
then LLVM_AR_NAME should be set to llvm-ar-3.5.

`

// Performs the environmental sanity check.
//
//        Performs the following checks in order:
//
//        1. Check that the OS is supported.
//        2. Checks that the compiler settings make sense.
//        3. Checks that the needed LLVM utilities exists.
//        4. Check that the store, if set, exists.
//
func sanityCheck() {

	checkOS()

	if !checkCompilers() {
		os.Exit(1)
	}

	if !checkAuxiliaries() {
		os.Exit(1)
	}

	checkStore()

}

func checkOS() {

	platform := runtime.GOOS

	if platform == "darwin" || platform == "linux" || platform == "freebsd" {
		return
	}

	shared.LogFatal("We do not support the OS %s", platform)
}

func checkCompilers() bool {

	cc := shared.GetCompilerExecName("clang")
	ccOK, _, ccVersion := checkExecutable(cc, "-v")
	if !ccOK {
		shared.LogError("The C compiler %s was not found or not executable.\nBetter not try using gclang!\n", cc)
	} else {
		shared.LogWrite("The C compiler %s is:\n\n\t%s\n\n", cc, extractLine(ccVersion, 0))
	}

	cxx := shared.GetCompilerExecName("clang++")
	cxxOK, _, cxxVersion := checkExecutable(cxx, "-v")
	if !ccOK {
		shared.LogError("The CXX compiler %s was not found or not executable.\nBetter not try using gclang++!\n", cxx)
	} else {
		shared.LogWrite("The CXX compiler %s is:\n\n\t%s\n\n", cxx, extractLine(cxxVersion, 0))
	}

	return ccOK || cxxOK
}

func extractLine(version string, n int) string {
	if len(version) == 0 {
		return version
	}
	lines := strings.Split(version, "\n")
	var line string
	lenLines := len(lines)
	if n < lenLines {
		line = lines[n]
	} else {
		line = lines[lenLines-1]
	}

	return strings.TrimSpace(line)

}

// Executes a command then returns true for success, false if there was an error, err is either nil or the error.
func checkExecutable(cmdExecName string, varg string) (success bool, err error, output string) {
	cmd := exec.Command(cmdExecName, varg)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	success = (err == nil)
	output = out.String()
	return
}

func checkAuxiliaries() bool {
	linkerName := shared.LLVMLINKName
	archiverName := shared.LLVMARName

	if linkerName == "" {
		linkerName = "llvm-link"
	}

	if archiverName == "" {
		archiverName = "llvm-ar"
	}

	linkerOK, _, linkerVersion := checkExecutable(linkerName, "-version")

	if !linkerOK {
		shared.LogError("The bitcode linker %s was not found or not executable.\nBetter not try using get-bc!\n", linkerName)
		shared.LogError(explainLLVMLINKNAME)
	} else {
		shared.LogWrite("The bitcode linker %s is:\n\n\t%s\n\n", linkerName, extractLine(linkerVersion, 1))
	}

	archiverOK, _, archiverVersion := checkExecutable(archiverName, "-version")

	if !archiverOK {
		shared.LogError("The bitcode archiver %s was not found or not executable.\nBetter not try using get-bc!\n", archiverName)
		shared.LogError(explainLLVMARNAME)
	} else {
		shared.LogWrite("The bitcode archiver %s is:\n\n\t%s\n\n", archiverName, extractLine(archiverVersion, 1))
	}

	return true
}

func checkStore() {
	storeDir := shared.BitcodeStorePath

	if storeDir != "" {
		finfo, err := os.Stat(storeDir)
		if err != nil && os.IsNotExist(err) {
			shared.LogError("The bitcode archive %s does not exist!\n\n", storeDir)
			return
		}
		if !finfo.Mode().IsDir() {
			shared.LogError("The bitcode archive %s is not a directory!\n\n", storeDir)
			return
		}
		shared.LogWrite("Using the bitcode archive %s\n\n", storeDir)
		return
	}
	shared.LogWrite("Not using a bitcode store.\n\n")
}
