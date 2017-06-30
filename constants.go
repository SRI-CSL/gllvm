package main

import (
	"os"
)

const (
	//The name of our ELF section of "bitcode paths".
	ELFSectionName    = ".llvm_bc"
	//The name of our MACH-O segment of "bitcode paths".
	DarwinSegmentName = "__WLLVM"
	//The name of our MACH-O section of "bitcode paths".
	DarwinSectionName = "__llvm_bc"
)

//The user configured directory holding the LLVM binary tools.
var LLVMToolChainBinDir string
//The user configured name of the clang compiler.
var LLVMCCName string
//The user configured name of the clang++ compiler.
var LLVMCXXName string
//The user configured name of the llvm-ar.
var LLVMARName string
//The user configured name of the llvm-link.
var LLVMLINKName string

var ConfigureOnly string
var BitcodeStorePath string

func init() {

	LLVMToolChainBinDir = os.Getenv("GLLVM_TOOLS_PATH")
	LLVMCCName = os.Getenv("GLLVM_CC_NAME")
	LLVMCXXName = os.Getenv("GLLVM_CXX_NAME")
	LLVMARName = os.Getenv("GLLVM_AR_NAME")
	LLVMLINKName = os.Getenv("GLLVM_LINK_NAME")

	ConfigureOnly = os.Getenv("GLLVM_CONFIGURE_ONLY")
	BitcodeStorePath = os.Getenv("GLLVM_BC_STORE")
}
