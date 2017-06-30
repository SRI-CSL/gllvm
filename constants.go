package main

const (
	// Environment variables
	envCONFIGURE_ONLY    = "GLLVM_CONFIGURE_ONLY"
	envTOOLS_PATH        = "GLLVM_TOOLS_PATH"
	envC_COMPILER_NAME   = "GLLVM_CC_NAME"
	envCXX_COMPILER_NAME = "GLLVM_CXX_NAME"
	envLINKER_NAME       = "GLLVM_LINK_NAME"
	envAR_NAME           = "GLLVM_AR_NAME"
	envBC_STORE_PATH     = "GLLVM_BC_STORE"

	// Gllvm functioning  (once we have it working we can change the W to G; but for the time being leave it so that extract-bc works)
	elfSECTION_NAME    = ".llvm_bc"
	darwinSEGMENT_NAME = "__WLLVM"
	darwinSECTION_NAME = "__llvm_bc"
)

const (
	// File types
	ftUNDEFINED = iota
	ftELF_EXECUTABLE
	ftELF_OBJECT
	ftELF_SHARED
	ftMACH_EXECUTABLE
	ftMACH_OBJECT
	ftMACH_SHARED
	ftARCHIVE
)
