package shared

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

//LLVMConfigureOnly is the user configured flag indicating a single pass mode is required.
var LLVMConfigureOnly string

//LLVMBitcodeStorePath is the user configured location of the bitcode archive.
var LLVMBitcodeStorePath string

//LLVMLoggingLevel is the user configured logging level: ERROR, WARNING, INFO, DEBUG.
var LLVMLoggingLevel string

//LLVMLoggingFile is the path to the optional logfile (useful when configuring)
var LLVMLoggingFile string

const (
	envpath = "GLLVM_TOOLS_PATH"     //"LLVM_COMPILER_PATH"
	envcc   = "GLLVM_CC_NAME"        //"LLVM_CC_NAME"
	envcxx  = "GLLVM_CXX_NAME"       //"LLVM_CXX_NAME"
	envar   = "GLLVM_AR_NAME"        //"LLVM_AR_NAME"
	envlnk  = "GLLVM_LINK_NAME"      //"LLVM_LINK_NAME"
	envcfg  = "GLLVM_CONFIGURE_ONLY" //"WLLVM_CONFIGURE_ONLY"
	envbc   = "GLLVM_BC_STORE"       //"WLLVM_BC_STORE"
	envlvl  = "GLLVM_OUTPUT_LEVEL"   //"WLLVM_OUTPUT_LEVEL
	envfile = "GLLVM_OUTPUT_FILE"    //"WLLVM_OUTPUT_FILE"
)

func init() {

	LLVMToolChainBinDir = os.Getenv(envpath)
	LLVMCCName = os.Getenv(envcc)
	LLVMCXXName = os.Getenv(envcxx)
	LLVMARName = os.Getenv(envar)
	LLVMLINKName = os.Getenv(envlnk)

	LLVMConfigureOnly = os.Getenv(envcfg)
	LLVMBitcodeStorePath = os.Getenv(envbc)

	LLVMLoggingLevel = os.Getenv(envlvl)
	LLVMLoggingFile = os.Getenv(envfile)

}
