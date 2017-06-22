package main

const(
    // Environment variables
    CONFIGURE_ONLY = "GOWLLVM_CONFIGURE_ONLY"
    COMPILER_PATH = "GOWLLVM_COMPILER_PATH"
    C_COMPILER_NAME = "GOWLLVM_CC_NAME"
    CXX_COMPILER_NAME = "GOWLLVM_CXX_NAME"
    BC_STORE_PATH = "GOWLLVM_BC_STORE"

    // Gowllvm functioning
    ELF_SECTION_NAME = ".llvm_bc"
    DARWIN_SEGMENT_NAME = "__WLLVM"
    DARWIN_SECTION_NAME = "__llvm_bc"
)
