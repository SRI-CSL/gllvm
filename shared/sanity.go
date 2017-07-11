package shared

import (
	"bytes"
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

// SanityCheck performs the environmental sanity check.
//
//        Performs the following checks in order:
//        0. Check the logging
//        1. Check that the OS is supported.
//        2. Checks that the compiler settings make sense.
//        3. Checks that the needed LLVM utilities exists.
//        4. Check that the store, if set, exists.
//
func SanityCheck() {

	checkLogging()

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
		LogWrite("Happily sitting atop \"%s\" operating system.\n\n", platform)
		return
	}

	LogFatal("We do not support the OS %s", platform)
}

func checkCompilers() bool {

	cc := GetCompilerExecName("clang")
	ccOK, ccVersion, _ := checkExecutable(cc, "-v")
	if !ccOK {
		LogError("The C compiler %s was not found or not executable.\nBetter not try using gclang!\n", cc)
	} else {
		LogWrite("The C compiler %s is:\n\n\t%s\n\n", cc, extractLine(ccVersion, 0))
	}

	cxx := GetCompilerExecName("clang++")
	cxxOK, cxxVersion, _ := checkExecutable(cxx, "-v")
	if !ccOK {
		LogError("The CXX compiler %s was not found or not executable.\nBetter not try using gclang++!\n", cxx)
	} else {
		LogWrite("The CXX compiler %s is:\n\n\t%s\n\n", cxx, extractLine(cxxVersion, 0))
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
func checkExecutable(cmdExecName string, varg string) (success bool, output string, err error) {
	cmd := exec.Command(cmdExecName, varg)
	var out bytes.Buffer
	//strangely clang writes it's version out on stderr
	//so we conflate the two to be tolerant.
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	success = (err == nil)
	output = out.String()
	return
}

func checkAuxiliaries() bool {
	linkerName := LLVMLINKName
	archiverName := LLVMARName

	if linkerName == "" {
		linkerName = "llvm-link"
	}

	if archiverName == "" {
		archiverName = "llvm-ar"
	}

	linkerOK, linkerVersion, _ := checkExecutable(linkerName, "-version")

	if !linkerOK {
		LogError("The bitcode linker %s was not found or not executable.\nBetter not try using get-bc!\n", linkerName)
		LogError(explainLLVMLINKNAME)
	} else {
		LogWrite("The bitcode linker %s is:\n\n\t%s\n\n", linkerName, extractLine(linkerVersion, 1))
	}

	archiverOK, archiverVersion, _ := checkExecutable(archiverName, "-version")

	if !archiverOK {
		LogError("The bitcode archiver %s was not found or not executable.\nBetter not try using get-bc!\n", archiverName)
		LogError(explainLLVMARNAME)
	} else {
		LogWrite("The bitcode archiver %s is:\n\n\t%s\n\n", archiverName, extractLine(archiverVersion, 1))
	}

	return true
}

func checkStore() {
	storeDir := LLVMBitcodeStorePath

	if storeDir != "" {
		finfo, err := os.Stat(storeDir)
		if err != nil && os.IsNotExist(err) {
			LogError("The bitcode archive %s does not exist!\n\n", storeDir)
			return
		}
		if !finfo.Mode().IsDir() {
			LogError("The bitcode archive %s is not a directory!\n\n", storeDir)
			return
		}
		LogWrite("Using the bitcode archive %s\n\n", storeDir)
		return
	}
	LogWrite("Not using a bitcode store.\n\n")
}

func checkLogging() {

	if LLVMLoggingFile != "" {
		// override the redirection so we output to the terminal (would be unnecessary
		// if we multiplexed in logging.go)
		loggingFilePointer = os.Stderr
		LogWrite("\nLogging output directed to %s.\n", LLVMLoggingFile)
	} else {
		LogWrite("\nLogging output to standard error.\n")
	}
	if LLVMLoggingLevel != "" {
		if _, ok := loggingLevels[LLVMLoggingLevel]; ok {
			LogWrite("Logging level is set to %s.\n\n", LLVMLoggingLevel)
		} else {
			LogWrite("Logging level is set to UNKNOWN level %s, using default of ERROR.\n\n", LLVMLoggingLevel)
		}
	} else {
		LogWrite("Logging level not set, using default of ERROR.\n\n")
	}
}
