package test

import (
	"fmt"
	"github.com/SRI-CSL/gllvm/shared"
	"os"
	"testing"
)

const (
	verbose = false
)

func checkExecutables(t *testing.T, ea shared.ExtractionArgs,
	llvmLinker string, llvmArchiver string, archiver string,
	clang string, clangpp string) {
	if ea.LlvmLinkerName != llvmLinker {
		t.Errorf("ParseSwitches: LlvmLinkerName incorrect: %v\n", ea.LlvmLinkerName)
	}
	if ea.LlvmArchiverName != llvmArchiver {
		t.Errorf("ParseSwitches: LlvmArchiverName incorrect: %v\n", ea.LlvmArchiverName)
	}
	if ea.ArchiverName != archiver {
		t.Errorf("ParseSwitches: ArchiverName incorrect: %v\n", ea.ArchiverName)
	}

	eclang := shared.GetCompilerExecName("clang")
	if eclang != clang {
		t.Errorf("C compiler not correct: %v\n", eclang)
	}
	eclangpp := shared.GetCompilerExecName("clang++")
	if eclangpp != clangpp {
		t.Errorf("C++ compiler not correct: %v\n", eclangpp)
	}
	if verbose {
		fmt.Printf("\nParseSwitches: %v\nclang = %v\nclang++ = %v\n", ea, clang, clangpp)
	}
}

func Test_env_and_args(t *testing.T) {

	args := []string{"get-bc", "-v", "../data/hello"}

	if verbose {
		shared.PrintEnvironment()
	}


	shared.ResetEnvironment()

	ea := shared.ParseSwitches(args)
	if !ea.Verbose {
		t.Errorf("ParseSwitches: -v flag not working\n")
	}
	if ea.WriteManifest || ea.SortBitcodeFiles || ea.BuildBitcodeModule || ea.KeepTemp {
		t.Errorf("ParseSwitches: defaults not correct\n")
	}
	if ea.InputFile != "../data/hello" {
		t.Errorf("ParseSwitches: InputFile incorrect: %v\n", ea.InputFile)
	}

	//iam: this test assumes LLVMToolChainBinDir = ""
	checkExecutables(t, ea, "llvm-link", "llvm-ar", "ar", "clang", "clang++")

	os.Setenv("LLVM_COMPILER_PATH", "/the_future_is_here")
	os.Setenv("LLVM_CC_NAME", "clang-666")
	os.Setenv("LLVM_CXX_NAME", "clang++-666")
	os.Setenv("LLVM_LINK_NAME", "llvm-link-666")
	os.Setenv("LLVM_AR_NAME", "llvm-ar-666")

	shared.FetchEnvironment()

	ea = shared.ParseSwitches(args)

	checkExecutables(t, ea,
		"/the_future_is_here/llvm-link-666",
		"/the_future_is_here/llvm-ar-666",
		"ar",
		"/the_future_is_here/clang-666",
		"/the_future_is_here/clang++-666")

	args = []string{"get-bc", "-a", "llvm-ar-665", "-l", "llvm-link-665", "../data/hello"}

	ea = shared.ParseSwitches(args)

	checkExecutables(t, ea,
		"llvm-link-665",
		"llvm-ar-665",
		"ar",
		"/the_future_is_here/clang-666",
		"/the_future_is_here/clang++-666")
}
