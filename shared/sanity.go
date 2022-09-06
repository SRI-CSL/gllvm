//
// OCCAM
//
// Copyright (c) 2017, SRI International
//
//  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of SRI International nor the names of its contributors may
//   be used to endorse or promote products derived from this software without
//   specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package shared

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
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
const explainLLVMFNAME = `

If your flang compiler is not called flang, but something else,
then you will need to set the environment variable LLVM_F_NAME to
the appropriate string. For example if your flang is called flang-7
then LLVM_F_NAME should be set to flang-7.

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

type sanityArgs struct {
	Environment bool
}

// SanityCheck performs the environmental sanity check.
//
//	Performs the following checks in order:
//	0. Check the logging
//	1. Check that the OS is supported.
//	2. Checks that the compiler settings make sense.
//	3. Checks that the needed LLVM utilities exists.
//	4. Check that the store, if set, exists.
func SanityCheck() {

	sa := parseSanitySwitches()

	informUser("\nVersion info: gsanity-check version %v\nReleased: %v\n", gllvmVersion, gllvmReleaseDate)

	if sa.Environment {
		PrintEnvironment()
	}

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

func parseSanitySwitches() (sa sanityArgs) {
	sa = sanityArgs{
		Environment: false,
	}

	environmentPtr := flag.Bool("e", false, "show environment")

	flag.Parse()

	sa.Environment = *environmentPtr

	return
}

func checkOS() {

	platform := runtime.GOOS

	if platform == osDARWIN || platform == osLINUX || platform == osFREEBSD {
		informUser("Happily sitting atop \"%s\" operating system.\n\n", platform)
		return
	}

	informUser("We do not support the OS %s", platform)
	os.Exit(1)

}

func checkCompilers() bool {

	cc := GetCompilerExecName("clang")
	ccOK, ccVersion, _ := checkExecutable(cc, "-v")
	if !ccOK {
		informUser("The C compiler %s was not found or not executable.\nBetter not try using gclang!\n", cc)
		informUser(explainLLVMCOMPILERPATH)
		informUser(explainLLVMCCNAME)

	} else {
		informUser("The C compiler %s is:\n\n\t%s\n\n", cc, extractLine(ccVersion, 0))
	}

	cxx := GetCompilerExecName("clang++")
	cxxOK, cxxVersion, _ := checkExecutable(cxx, "-v")
	if !cxxOK {
		informUser("The CXX compiler %s was not found or not executable.\nBetter not try using gclang++!\n", cxx)
		informUser(explainLLVMCOMPILERPATH)
		informUser(explainLLVMCXXNAME)
	} else {
		informUser("The CXX compiler %s is:\n\n\t%s\n\n", cxx, extractLine(cxxVersion, 0))
	}
	f := GetCompilerExecName("flang")
	fOK, fVersion, _ := checkExecutable(f, "-v")
	if !fOK {
		informUser("The Fortran compiler %s was not found or not executable.\nBetter not try using gflang!\n", f)
		informUser(explainLLVMCOMPILERPATH)
		informUser(explainLLVMFNAME)
	} else {
		informUser("The Fortran compiler %s is:\n\n\t%s\n\n", f, extractLine(fVersion, 0))
	}

	//FIXME: why "or" rather than "and"? BECAUSE: if you only need CC, not having CXX is not an error.
	return ccOK || cxxOK || fOK
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
	LogInfo("checkExecutable: %s %s returned %s\n", cmdExecName, varg, output)
	LogInfo("checkExecutable: returning (%s %s %s)\n", success, output, err)
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

	linkerName = filepath.Join(LLVMToolChainBinDir, linkerName)

	linkerOK, linkerVersion, _ := checkExecutable(linkerName, "-version")

	// iam: 5/8/2018 3.4 llvm-link and llvm-ar return exit status 1 for -version. GO FIGURE.
	if !linkerOK && !strings.Contains(linkerVersion, "LLVM") {
		informUser("The bitcode linker %s was not found or not executable.\nBetter not try using get-bc!\n", linkerName)
		informUser(explainLLVMLINKNAME)
	} else {
		informUser("The bitcode linker %s is:\n\n\t%s\n\n", linkerName, extractLine(linkerVersion, 1))
	}

	archiverName = filepath.Join(LLVMToolChainBinDir, archiverName)
	archiverOK, archiverVersion, _ := checkExecutable(archiverName, "-version")

	// iam: 5/8/2018 3.4 llvm-link and llvm-ar return exit status 1 for -version. GO FIGURE.
	if !archiverOK && !strings.Contains(linkerVersion, "LLVM") {
		informUser("The bitcode archiver %s was not found or not executable.\nBetter not try using get-bc!\n", archiverName)
		informUser(explainLLVMARNAME)
	} else {
		informUser("The bitcode archiver %s is:\n\n\t%s\n\n", archiverName, extractLine(archiverVersion, 1))
	}

	return linkerOK && archiverOK
}

func checkStore() {
	storeDir := LLVMBitcodeStorePath

	if storeDir != "" {
		finfo, err := os.Stat(storeDir)
		if err != nil && os.IsNotExist(err) {
			informUser("The bitcode archive %s does not exist!\n\n", storeDir)
			return
		}
		if !finfo.Mode().IsDir() {
			informUser("The bitcode archive %s is not a directory!\n\n", storeDir)
			return
		}
		informUser("Using the bitcode archive %s\n\n", storeDir)
		return
	}
	informUser("Not using a bitcode store.\n\n")
}

func checkLogging() {

	if LLVMLoggingFile != "" {
		informUser("\nLogging output directed to %s.\n", LLVMLoggingFile)
	} else {
		informUser("\nLogging output to standard error.\n")
	}
	if LLVMLoggingLevel != "" {
		if _, ok := loggingLevels[LLVMLoggingLevel]; ok {
			informUser("Logging level is set to %s.\n\n", LLVMLoggingLevel)
		} else {
			informUser("Logging level is set to UNKNOWN level %s, using default of ERROR.\n\n", LLVMLoggingLevel)
		}
	} else {
		informUser("Logging level not set, using default of WARNING.\n\n")
	}

	informUser("Logging configuration uses the environment variables:\n\n\tWLLVM_OUTPUT_LEVEL and WLLVM_OUTPUT_FILE.\n\n")
}
