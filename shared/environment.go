//
// OCCAM
//
// Copyright (c) 2017, SRI International
//
//  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of SRI International nor the names of its contributors may
//   be used to endorse or promote products derived from this software without
//   specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package shared

import (
	"os"
	"strconv"
	"strings"
)

const (

	// SectionNameBitCode section of "bitcode paths"
	SectionNameBitCode = "llvm_bc"

	// SectionNameFlags is the name of our ELF section of frontend flags for each translation unit
	SectionNameFlags = "gllvm_flags"

	//SegmentNameDarwin is the name of our MACH-O segment of "bitcode paths".
	SegmentNameDarwin = "__WLLVM"

	// //DarwinSectionName is the name of our MACH-O section of "bitcode paths".
	// DarwinSectionName = "__llvm_bc"

	// // DarwinFrontendFlagsSectionName is the name of our MACH-O section of frontend flags for each translation unit
	// DarwinFrontendFlagsSectionName = "__gllvm_flags"
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

//LLVMObjcopy is the path to the objcopy executable used to attach the bitcode on *nix.
var LLVMObjcopy string

//LLVMLd is the path to the ld executable used to attach the bitcode on OSX.
var LLVMLd string

//LLVMbcGen is the list of args to pass to clang during the bitcode generation step.
var LLVMbcGen []string

//LLVMEmbedFrontendArgs flags whether or not we embed frontend arguments in each compiled object.
var LLVMEmbedFrontendArgs bool

const (
	envpath    = "LLVM_COMPILER_PATH"
	envcc      = "LLVM_CC_NAME"
	envcxx     = "LLVM_CXX_NAME"
	envar      = "LLVM_AR_NAME"
	envlnk     = "LLVM_LINK_NAME"
	envcfg     = "WLLVM_CONFIGURE_ONLY"
	envbc      = "WLLVM_BC_STORE"
	envlvl     = "WLLVM_OUTPUT_LEVEL"
	envfile    = "WLLVM_OUTPUT_FILE"
	envld      = "GLLVM_LD"      //iam: we are deviating from wllvm here.
	envobjcopy = "GLLVM_OBJCOPY" //iam: we are deviating from wllvm here.
	//wllvm uses a BINUTILS_TARGET_PREFIX, which seems less general.
	//iam: 03/24/2020 new feature to pass things like "-flto -fwhole-program-vtables"
	// to clang during the bitcode generation step
	envbcgen = "LLVM_BITCODE_GENERATION_FLAGS"

	// 4/9/2020 new feature to embed the frontend compiler and linker flags
	// passed to clang within each object
	envembedargs = "GLLVM_EMBED_FRONTEND_ARGS"
)

func init() {
	FetchEnvironment()
}

// PrintEnvironment is used for printing the aspects of the environment that concern us
func PrintEnvironment() {
	vars := []string{envpath, envcc, envcxx, envar, envlnk, envcfg, envbc, envlvl, envfile, envobjcopy, envld, envbcgen, envembedargs}

	informUser("\nLiving in this environment:\n\n")
	for _, v := range vars {
		val, defined := os.LookupEnv(v)
		if defined {
			informUser("%v = \"%v\"\n", v, val)
		} else {
			informUser("%v is NOT defined\n", v)
		}
	}

}

// ResetEnvironment resets the globals, it is only used in testing
func ResetEnvironment() {
	LLVMToolChainBinDir = ""
	LLVMCCName = ""
	LLVMCXXName = ""
	LLVMARName = ""
	LLVMLINKName = ""
	LLVMConfigureOnly = ""
	LLVMBitcodeStorePath = ""
	LLVMLoggingLevel = ""
	LLVMLoggingFile = ""
	LLVMObjcopy = ""
	LLVMLd = ""
	LLVMbcGen = []string{}
	LLVMEmbedFrontendArgs = false
}

// FetchEnvironment is used in initializing our globals, it is also used in testing
func FetchEnvironment() {
	LLVMToolChainBinDir = os.Getenv(envpath)
	LLVMCCName = os.Getenv(envcc)
	LLVMCXXName = os.Getenv(envcxx)
	LLVMARName = os.Getenv(envar)
	LLVMLINKName = os.Getenv(envlnk)

	LLVMConfigureOnly = os.Getenv(envcfg)
	LLVMBitcodeStorePath = os.Getenv(envbc)

	LLVMLoggingLevel = os.Getenv(envlvl)
	LLVMLoggingFile = os.Getenv(envfile)

	LLVMObjcopy = os.Getenv(envobjcopy)
	LLVMLd = os.Getenv(envld)

	LLVMbcGen = strings.Fields(os.Getenv(envbcgen))

	var err error
	LLVMEmbedFrontendArgs, err = strconv.ParseBool(os.Getenv(envembedargs))
	if err != nil {
		LLVMEmbedFrontendArgs = false
	}
}
