package main

import (
	"os"
)

const (
	//ELFSectionName is the name of our ELF section of "bitcode paths".
	ELFSectionName = ".llvm_bc"
	//DarwinSegmentName is the name of our MACH-O segment of "bitcode paths".
	DarwinSegmentName = "__WLLVM"
	//DarwinSectionName is the name of our MACH-O section of "bitcode paths".
	DarwinSectionName = "__llvm_bc"
)

//LLVMToolChainBinDir is the user configured directory holding the LLVM binary tools.
var LLVMToolChainBinDir string

//LLVMCCName is the user configured name of the clang compiler.
var LLVMCCName string

//LLVMCXXName is the user configured name of the clang++ compiler.
var LLVMCXXName string

//LLVMARName is the user configured name of the llvm-ar.
var LLVMARName string

//LLVMLINKName is the user configured name of the llvm-link.
var LLVMLINKName string

//ConfigureOnly is the user configured flag indicating a single pass mode is required.
var ConfigureOnly string

//BitcodeStorePath is the user configured location of the bitcode archive.
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
