package main

import (
       "os"
)

const (
	elfSECTIONNAME    = ".llvm_bc"
	darwinSEGMENTNAME = "__WLLVM"
	darwinSECTIONNAME = "__llvm_bc"
)

var LLVMToolChainBinDir = ""
var LLVMCCName          = ""
var LLVMCXXName         = ""
var LLVMARName          = ""
var LLVMLINKName        = ""

var ConfigureOnly       = ""
var BitcodeStorePath    = ""

func init(){

     LLVMToolChainBinDir = os.Getenv("GLLVM_TOOLS_PATH")
     LLVMCCName          = os.Getenv("GLLVM_CC_NAME")
     LLVMCXXName         = os.Getenv("GLLVM_CXX_NAME")
     LLVMARName          = os.Getenv("GLLVM_AR_NAME")
     LLVMLINKName        = os.Getenv("GLLVM_LINK_NAME")
 
     ConfigureOnly        = os.Getenv("GLLVM_CONFIGURE_ONLY")
     BitcodeStorePath     = os.Getenv("GLLVM_BC_STORE")
}



const (
	// File types
	ftUNDEFINED = iota
	ftELFEXECUTABLE
	ftELFOBJECT
	ftELFSHARED
	ftMACHEXECUTABLE
	ftMACHOBJECT
	ftMACHSHARED
	ftARCHIVE
)
