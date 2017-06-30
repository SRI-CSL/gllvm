package main

const (
	// Environment variables
	envCONFIGUREONLY    = "GLLVM_CONFIGURE_ONLY"
	envTOOLSPATH        = "GLLVM_TOOLS_PATH"
	envCCOMPILERNAME   = "GLLVM_CC_NAME"
	envCXXCOMPILERNAME = "GLLVM_CXX_NAME"
	envLINKERNAME       = "GLLVM_LINK_NAME"
	envARNAME           = "GLLVM_AR_NAME"
	envBCSTOREPATH     = "GLLVM_BC_STORE"

	// Gllvm functioning  (once we have it working we can change the W to G; but for the time being leave it so that extract-bc works)
	elfSECTIONNAME    = ".llvm_bc"
	darwinSEGMENTNAME = "__WLLVM"
	darwinSECTIONNAME = "__llvm_bc"
)

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
