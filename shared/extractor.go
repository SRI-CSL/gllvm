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
	"debug/elf"
	"debug/macho"
	"flag"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type extractionArgs struct {
	InputFile           string
	InputType           int
	OutputFile          string
	LinkerName          string
	ArchiverName        string
	ArArgs              []string
	ObjectTypeInArchive int // Type of file that can be put into an archive
	Extractor           func(string) []string
	Verbose             bool
	WriteManifest       bool
	BuildBitcodeArchive bool
}

//Extract extracts the LLVM bitcode according to the arguments it is passed.
func Extract(args []string) {
	ea := parseSwitches()

	// Set arguments according to runtime OS
	switch platform := runtime.GOOS; platform {
	case "freebsd", "linux":
		ea.Extractor = extractSectionUnix
		if ea.Verbose {
			ea.ArArgs = append(ea.ArArgs, "xv")
		} else {
			ea.ArArgs = append(ea.ArArgs, "x")
		}
		ea.ObjectTypeInArchive = fileTypeELFOBJECT
	case "darwin":
		ea.Extractor = extractSectionDarwin
		ea.ArArgs = append(ea.ArArgs, "-x")
		if ea.Verbose {
			ea.ArArgs = append(ea.ArArgs, "-v")
		}
		ea.ObjectTypeInArchive = fileTypeMACHOBJECT
	default:
		LogFatal("Unsupported platform: %s.", platform)
	}

	// Create output filename if not given
	if ea.OutputFile == "" {
		if ea.InputType == fileTypeARCHIVE {
			var ext string
			if ea.BuildBitcodeArchive {
				ext = ".a.bc"
			} else {
				ext = ".bca"
			}
			ea.OutputFile = strings.TrimSuffix(ea.InputFile, ".a") + ext
		} else {
			ea.OutputFile = ea.InputFile + ".bc"
		}
	}

	switch ea.InputType {
	case fileTypeELFEXECUTABLE,
		fileTypeELFSHARED,
		fileTypeELFOBJECT,
		fileTypeMACHEXECUTABLE,
		fileTypeMACHSHARED,
		fileTypeMACHOBJECT:
		handleExecutable(ea)
	case fileTypeARCHIVE:
		handleArchive(ea)
	default:
		LogFatal("Incorrect input file type %v.", ea.InputType)
	}

}

func parseSwitches() (ea extractionArgs) {
	ea = extractionArgs{
		LinkerName:   "llvm-link",
		ArchiverName: "llvm-ar",
	}

	verbosePtr := flag.Bool("v", false, "verbose mode")

	writeManifestPtr := flag.Bool("m", false, "write the manifest")

	buildBitcodeArchive := flag.Bool("b", false, "build a bitcode module(FIXME? should this be archive)")

	outputFilePtr := flag.String("o", "", "the output file")

	archiverNamePtr := flag.String("a", "", "the llvm archiver")

	linkerNamePtr := flag.String("l", "", "the llvm linker")

	flag.Parse()

	ea.Verbose = *verbosePtr
	ea.WriteManifest = *writeManifestPtr
	ea.BuildBitcodeArchive = *buildBitcodeArchive

	if *archiverNamePtr != "" {
		ea.ArchiverName = *archiverNamePtr
	} else {
		if LLVMARName != "" {
			ea.ArchiverName = filepath.Join(LLVMToolChainBinDir, LLVMARName)
		}
	}

	if *linkerNamePtr != "" {
		ea.LinkerName = *linkerNamePtr
	} else {
		if LLVMLINKName != "" {
			ea.LinkerName = filepath.Join(LLVMToolChainBinDir, LLVMLINKName)
		}
	}

	ea.OutputFile = *outputFilePtr

	inputFiles := flag.Args()

	LogInfo("ea.Verbose: %v\n", ea.Verbose)
	LogInfo("ea.WriteManifest: %v\n", ea.WriteManifest)
	LogInfo("ea.BuildBitcodeArchive: %v\n", ea.BuildBitcodeArchive)
	LogInfo("ea.ArchiverName: %v\n", ea.ArchiverName)
	LogInfo("ea.LinkerName: %v\n", ea.LinkerName)
	LogInfo("ea.OutputFile: %v\n", ea.OutputFile)

	if len(inputFiles) != 1 {
		LogFatal("Can currently only deal with exactly one input file, sorry. You gave me %v\n.", len(inputFiles))
	}

	ea.InputFile = inputFiles[0]

	LogInfo("ea.InputFile: %v\n", ea.InputFile)

	if _, err := os.Stat(ea.InputFile); os.IsNotExist(err) {
		LogFatal("The input file %s  does not exist.", ea.InputFile)
	}
	realPath, err := filepath.EvalSymlinks(ea.InputFile)
	if err != nil {
		LogFatal("There was an error getting the real path of %s.", ea.InputFile)
	}
	ea.InputFile = realPath
	ea.InputType = getFileType(realPath)

	LogInfo("ea.InputFile real path: %v\n", ea.InputFile)

	return
}

func handleExecutable(ea extractionArgs) {
	artifactPaths := ea.Extractor(ea.InputFile)
	if len(artifactPaths) == 0 {
		return
	}
	filesToLink := make([]string, len(artifactPaths))
	for i, artPath := range artifactPaths {
		filesToLink[i] = resolveBitcodePath(artPath)
	}
	extractTimeLinkFiles(ea, filesToLink)

	// Write manifest
	if ea.WriteManifest {
		writeManifest(ea, filesToLink, artifactPaths)
	}
}

func handleArchive(ea extractionArgs) {
	// List bitcode files to link
	var bcFiles []string
	var artifactFiles []string

	// Create tmp dir
	tmpDirName, err := ioutil.TempDir("", "gllvm")
	if err != nil {
		LogFatal("The temporary directory in which to extract object files could not be created.")
	}
	defer os.RemoveAll(tmpDirName)

	// Extract objects to tmpDir
	arArgs := ea.ArArgs
	inputAbsPath, _ := filepath.Abs(ea.InputFile)
	arArgs = append(arArgs, inputAbsPath)
	success, err := execCmd("ar", arArgs, tmpDirName)
	if !success {
		LogFatal("Failed to extract object files from %s to %s because: %v.\n", ea.InputFile, tmpDirName, err)
	}

	// Define object file handling closure
	var walkHandlingFunc = func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			fileType := getFileType(path)
			if fileType == ea.ObjectTypeInArchive {
				artifactPaths := ea.Extractor(path)
				for _, artPath := range artifactPaths {
					bcPath := resolveBitcodePath(artPath)
					bcFiles = append(bcFiles, bcPath)
				}
				artifactFiles = append(artifactFiles, artifactPaths...)
			}
		}
		return nil
	}

	// Handle object files
	filepath.Walk(tmpDirName, walkHandlingFunc)

	// Build archive
	if ea.BuildBitcodeArchive {
		extractTimeLinkFiles(ea, bcFiles)
	} else {
		archiveBcFiles(ea, bcFiles)
	}

	// Write manifest
	if ea.WriteManifest {
		writeManifest(ea, bcFiles, artifactFiles)
	}
}

func archiveBcFiles(ea extractionArgs, bcFiles []string) {
	// We do not want full paths in the archive, so we need to chdir into each
	// bitcode's folder. Handle this by calling llvm-ar once for all bitcode
	// files in the same directory
	dirToBcMap := make(map[string][]string)
	for _, bcFile := range bcFiles {
		dirName, baseName := path.Split(bcFile)
		dirToBcMap[dirName] = append(dirToBcMap[dirName], baseName)
	}

	// Call llvm-ar from each directory
	absOutputFile, _ := filepath.Abs(ea.OutputFile)
	for dir, bcFilesInDir := range dirToBcMap {
		var args []string
		args = append(args, "rs", absOutputFile)
		args = append(args, bcFilesInDir...)
		success, err := execCmd(ea.ArchiverName, args, dir)
		LogInfo("ea.ArchiverName = %s, args = %v, dir = %s\n", ea.ArchiverName, args, dir)
		if !success {
			LogFatal("There was an error creating the bitcode archive: %v.\n", err)
		}
	}
	LogInfo("Built bitcode archive: %s.", ea.OutputFile)
}

func extractTimeLinkFiles(ea extractionArgs, filesToLink []string) {
	var linkArgs []string
	if ea.Verbose {
		linkArgs = append(linkArgs, "-v")
	}
	linkArgs = append(linkArgs, "-o", ea.OutputFile)
	linkArgs = append(linkArgs, filesToLink...)
	success, err := execCmd(ea.LinkerName, linkArgs, "")
	if !success {
		LogFatal("There was an error linking input files into %s because %v.\n", ea.OutputFile, err)
	}
	LogInfo("Bitcode file extracted to: %s.", ea.OutputFile)
}

func extractSectionDarwin(inputFile string) (contents []string) {
	machoFile, err := macho.Open(inputFile)
	if err != nil {
		LogFatal("Mach-O file %s could not be read.", inputFile)
	}
	section := machoFile.Section(DarwinSectionName)
	if section == nil {
		LogWarning("The %s section of %s is missing!\n", DarwinSectionName, inputFile)
		return
	}
	sectionContents, errContents := section.Data()
	if errContents != nil {
		LogFatal("Error reading the %s section of Mach-O file %s.", DarwinSectionName, inputFile)
	}
	contents = strings.Split(strings.TrimSuffix(string(sectionContents), "\n"), "\n")
	return
}

func extractSectionUnix(inputFile string) (contents []string) {
	elfFile, err := elf.Open(inputFile)
	if err != nil {
		LogFatal("ELF file %s could not be read.", inputFile)
		return
	}
	section := elfFile.Section(ELFSectionName)
	if section == nil {
		LogWarning("Error reading the %s section of ELF file %s.", ELFSectionName, inputFile)
		return
	}
	sectionContents, errContents := section.Data()
	if errContents != nil {
		LogWarning("Error reading the %s section of ELF file %s.", ELFSectionName, inputFile)
		return
	}
	contents = strings.Split(strings.TrimSuffix(string(sectionContents), "\n"), "\n")
	return
}

// Return the actual path to the bitcode file, or an empty string if it does not exist
func resolveBitcodePath(bcPath string) string {
	if _, err := os.Stat(bcPath); os.IsNotExist(err) {
		// If the bitcode file does not exist, try to find it in the store
		if LLVMBitcodeStorePath != "" {
			// Compute absolute path hash
			absBcPath, _ := filepath.Abs(bcPath)
			storeBcPath := path.Join(LLVMBitcodeStorePath, getHashedPath(absBcPath))
			if _, err := os.Stat(storeBcPath); os.IsNotExist(err) {
				return ""
			}
			return storeBcPath
		}
		return ""
	}
	return bcPath
}

func writeManifest(ea extractionArgs, bcFiles []string, artifactFiles []string) {
	section1 := "Physical location of extracted files:\n" + strings.Join(bcFiles, "\n") + "\n\n"
	section2 := "Build-time location of extracted files:\n" + strings.Join(artifactFiles, "\n")
	contents := []byte(section1 + section2)
	manifestFilename := ea.OutputFile + ".llvm.manifest"
	if err := ioutil.WriteFile(manifestFilename, contents, 0644); err != nil {
		LogFatal("There was an error while writing the manifest file: ", err)
	}
	LogInfo("Manifest file written to %s.", manifestFilename)
}
