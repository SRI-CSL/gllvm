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

func init() {

	LLVMToolChainBinDir = os.Getenv("GLLVM_TOOLS_PATH") //os.Getenv("LLVM_COMPILER_PATH")
	LLVMCCName = os.Getenv("GLLVM_CC_NAME")             //os.Getenv("LLVM_CC_NAME")
	LLVMCXXName = os.Getenv("GLLVM_CXX_NAME")           //os.Getenv("LLVM_CXX_NAME")
	LLVMARName = os.Getenv("GLLVM_AR_NAME")             //os.Getenv("LLVM_AR_NAME")
	LLVMLINKName = os.Getenv("GLLVM_LINK_NAME")         //os.Getenv("LLVM_LINK_NAME")

	LLVMConfigureOnly = os.Getenv("GLLVM_CONFIGURE_ONLY") //os.Getenv("WLLVM_CONFIGURE_ONLY")
	LLVMBitcodeStorePath = os.Getenv("GLLVM_BC_STORE")    //os.Getenv("WLLVM_BC_STORE")

	LLVMLoggingLevel = os.Getenv("GLLVM_OUTPUT_LEVEL") //os.Getenv("WLLVM_OUTPUT_LEVEL")
	LLVMLoggingFile = os.Getenv("GLLVM_OUTPUT_FILE")   //os.Getenv("WLLVM_OUTPUT_FILE")

}
