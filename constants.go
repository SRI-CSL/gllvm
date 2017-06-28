package main

const(
	// Environment variables
	CONFIGURE_ONLY    = "GLLVM_CONFIGURE_ONLY"
	TOOLS_PATH        = "GLLVM_TOOLS_PATH"
	C_COMPILER_NAME   = "GLLVM_CC_NAME"
	CXX_COMPILER_NAME = "GLLVM_CXX_NAME"
	LINKER_NAME       = "GLLVM_LINK_NAME"
	AR_NAME           = "GLLVM_AR_NAME"
	BC_STORE_PATH     = "GLLVM_BC_STORE"

	// Gllvm functioning  (once we have it working we can change the W to G; but for the time being leave it so that extract-bc works)
	ELF_SECTION_NAME    = ".llvm_bc"
	DARWIN_SEGMENT_NAME = "__WLLVM"
	DARWIN_SECTION_NAME = "__llvm_bc"

)

const (
	// File types
	FT_UNDEFINED = iota
	FT_ELF_EXECUTABLE
	FT_ELF_OBJECT
	FT_ELF_SHARED
	FT_MACH_EXECUTABLE
	FT_MACH_OBJECT
	FT_MACH_SHARED
	FT_ARCHIVE
)
