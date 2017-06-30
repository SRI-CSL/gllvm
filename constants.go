package main

const (
	// Environment variables
	env_CONFIGURE_ONLY    = "GLLVM_CONFIGURE_ONLY"
	env_TOOLS_PATH        = "GLLVM_TOOLS_PATH"
	env_C_COMPILER_NAME   = "GLLVM_CC_NAME"
	env_CXX_COMPILER_NAME = "GLLVM_CXX_NAME"
	env_LINKER_NAME       = "GLLVM_LINK_NAME"
	env_AR_NAME           = "GLLVM_AR_NAME"
	env_BC_STORE_PATH     = "GLLVM_BC_STORE"

	// Gllvm functioning  (once we have it working we can change the W to G; but for the time being leave it so that extract-bc works)
	elf_SECTION_NAME    = ".llvm_bc"
	darwin_SEGMENT_NAME = "__WLLVM"
	darwin_SECTION_NAME = "__llvm_bc"
)

const (
	// File types
	ft_UNDEFINED = iota
	ft_ELF_EXECUTABLE
	ft_ELF_OBJECT
	ft_ELF_SHARED
	ft_MACH_EXECUTABLE
	ft_MACH_OBJECT
	ft_MACH_SHARED
	ft_ARCHIVE
)
