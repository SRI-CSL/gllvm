package main

import (
	"os"
)

const (
	ELFSectionName    = ".llvm_bc"
	DarwinSegmentName = "__WLLVM"
	DarwinSectionName = "__llvm_bc"
)

var LLVMToolChainBinDir string
var LLVMCCName string
var LLVMCXXName string
var LLVMARName string
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

