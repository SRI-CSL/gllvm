package main

import (
       "os"
)

const (
	// Environment variables
	envCONFIGUREONLY   = "GLLVM_CONFIGURE_ONLY"
	envLINKERNAME      = "GLLVM_LINK_NAME"
	envARNAME          = "GLLVM_AR_NAME"
	envBCSTOREPATH     = "GLLVM_BC_STORE"

	// Gllvm functioning  (once we have it working we can change the W to G; but for the time being leave it so that extract-bc works)
	elfSECTIONNAME    = ".llvm_bc"
	darwinSEGMENTNAME = "__WLLVM"
	darwinSECTIONNAME = "__llvm_bc"
)

var LLVMToolChainBinDir = ""
var LLVMCCName          = ""
var LLVMCXXName         = ""

func init(){
     LLVMToolChainBinDir = os.Getenv("GLLVM_TOOLS_PATH")
     LLVMCCName = os.Getenv("GLLVM_CC_NAME")
     LLVMCXXName = os.Getenv("GLLVM_CXX_NAME")
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
