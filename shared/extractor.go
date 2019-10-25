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
	"bytes"
	"debug/elf"
	"debug/macho"
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

type extractionArgs struct {
	Verbose             bool // inform the user of what is going on
	WriteManifest       bool // write a manifest of bitcode files used
	SortBitcodeFiles    bool // sort the arguments to linking and archiving (debugging too)
	BuildBitcodeModule  bool // buld an archive rather than a module
	KeepTemp            bool // keep temporary linking folder
	LinkArgSize         int  // maximum size of a llvm-link command line
	InputType           int
	ObjectTypeInArchive int // Type of file that can be put into an archive
	InputFile           string
	OutputFile          string
	LlvmLinkerName      string
	LlvmArchiverName    string
	ArchiverName        string
	ArArgs              []string
	Extractor           func(string) []string
}

//Extract extracts the LLVM bitcode according to the arguments it is passed.
func Extract(args []string)(exitCode int) {
	ea := parseSwitches(args)

	// Set arguments according to runtime OS
	switch platform := runtime.GOOS; platform {
	case osFREEBSD, osLINUX:
		ea.Extractor = extractSectionUnix
		if ea.Verbose {
			ea.ArArgs = append(ea.ArArgs, "xv")
		} else {
			ea.ArArgs = append(ea.ArArgs, "x")
		}
		ea.ObjectTypeInArchive = fileTypeELFOBJECT
	case osDARWIN:
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
		if ea.InputType == fileTypeARCHIVE || ea.InputType == fileTypeTHINARCHIVE {
			var ext string
			if ea.BuildBitcodeModule {
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
	case fileTypeTHINARCHIVE:
		handleThinArchive(ea)
	default:
		LogFatal("Incorrect input file type %v.", ea.InputType)
	}

	//need to actually get an exitCode eventually.
	return 
}

func resolveTool(defaultPath string, envPath string, usrPath string) (path string) {
	if usrPath != defaultPath {
		path = usrPath
	} else {
		if LLVMToolChainBinDir != "" {
			if envPath != "" {
				path = filepath.Join(LLVMToolChainBinDir, envPath)
			} else {
				path = filepath.Join(LLVMToolChainBinDir, defaultPath)
			}
		} else {
			if envPath != "" {
				path = envPath
			} else {
				path = defaultPath
			}
		}
	}
	LogDebug("defaultPath = %s", defaultPath)
	LogDebug("envPath = %s", envPath)
	LogDebug("usrPath = %s", usrPath)
	LogDebug("path = %s", path)
	return
}

func parseSwitches(args []string) (ea extractionArgs) {
	var flagSet *flag.FlagSet = flag.NewFlagSet(args[0], flag.ExitOnError)
	flagSet.BoolVar(&ea.Verbose, "v", false, "verbose mode")
	flagSet.BoolVar(&ea.WriteManifest, "m", false, "write the manifest")
	flagSet.BoolVar(&ea.SortBitcodeFiles, "s", false, "sort the bitcode files")
	flagSet.BoolVar(&ea.BuildBitcodeModule, "b", false, "build a bitcode module")
	flagSet.StringVar(&ea.OutputFile, "o", "", "the output file")
	flagSet.StringVar(&ea.LlvmArchiverName, "a", "llvm-ar", "the llvm archiver (i.e. llvm-ar)")
	flagSet.StringVar(&ea.ArchiverName, "r", "ar", "the system archiver (i.e. ar)")
	flagSet.StringVar(&ea.LlvmLinkerName, "l", "llvm-link", "the llvm linker (i.e. llvm-link)")
	flagSet.IntVar(&ea.LinkArgSize, "n", 0, "maximum llvm-link command line size (in bytes)")
	flagSet.BoolVar(&ea.KeepTemp, "t", false, "keep temporary linking folder")
	flagSet.Parse(args[1:])
	ea.LlvmArchiverName = resolveTool("llvm-ar", LLVMARName, ea.LlvmArchiverName)
	ea.LlvmLinkerName = resolveTool("llvm-link", LLVMLINKName, ea.LlvmLinkerName)
	inputFiles := flagSet.Args()
	if len(inputFiles) != 1 {
		LogFatal("Can currently only deal with exactly one input file, sorry. You gave me %v input files.\n", len(inputFiles))
	}
	ea.InputFile = inputFiles[0]
	if _, err := os.Stat(ea.InputFile); os.IsNotExist(err) {
		LogFatal("The input file %s  does not exist.", ea.InputFile)
	}
	realPath, err := filepath.EvalSymlinks(ea.InputFile)
	if err != nil {
		LogFatal("There was an error getting the real path of %s.", ea.InputFile)
	}
	ea.InputFile = realPath
	ea.InputType = getFileType(realPath)
	//<this should be a method of extractionArgs>
	LogInfo("ea.Verbose: %v\n", ea.Verbose)
	LogInfo("ea.WriteManifest: %v\n", ea.WriteManifest)
	LogInfo("ea.BuildBitcodeModule: %v\n", ea.BuildBitcodeModule)
	LogInfo("ea.LlvmArchiverName: %v\n", ea.LlvmArchiverName)
	LogInfo("ea.LlvmLinkerName: %v\n", ea.LlvmLinkerName)
	LogInfo("ea.ArchiverName: %v\n", ea.ArchiverName)
	LogInfo("ea.OutputFile: %v\n", ea.OutputFile)
	LogInfo("ea.InputFile: %v\n", ea.InputFile)
	LogInfo("ea.InputFile real path: %v\n", ea.InputFile)
	LogInfo("ea.LinkArgSize %d", ea.LinkArgSize)
	LogInfo("ea.KeepTemp %v", ea.KeepTemp)
	//</this should be a method of extractionArgs>
	return
}

func handleExecutable(ea extractionArgs) {
	artifactPaths := ea.Extractor(ea.InputFile)

	if len(artifactPaths) < 20 {
		// naert: to avoid saturating the log when dealing with big file lists
		LogInfo("handleExecutable: artifactPaths = %v\n", artifactPaths)
	}

	if len(artifactPaths) == 0 {
		return
	}
	filesToLink := make([]string, len(artifactPaths))
	for i, artPath := range artifactPaths {
		filesToLink[i] = resolveBitcodePath(artPath)
	}

	// Sort the bitcode files
	if ea.SortBitcodeFiles {
		LogWarning("Sorting bitcode files.")
		sort.Strings(filesToLink)
		sort.Strings(artifactPaths)
	}

	// Write manifest
	if ea.WriteManifest {
		writeManifest(ea, filesToLink, artifactPaths)
	}

	linkBitcodeFiles(ea, filesToLink)
}

func handleThinArchive(ea extractionArgs) {
	// List bitcode files to link
	var artifactFiles []string

	var objectFiles []string
	var bcFiles []string

	objectFiles = listArchiveFiles(ea, ea.InputFile)

	LogInfo("handleThinArchive: extractionArgs = %v\nobjectFiles = %v\n", ea, objectFiles)

	for index, obj := range objectFiles {
		LogInfo("obj = '%v'\n", obj)
		if len(obj) > 0 {
			artifacts := ea.Extractor(obj)
			LogInfo("\t%v\n", artifacts)
			artifactFiles = append(artifactFiles, artifacts...)
			for _, bc := range artifacts {
				bcPath := resolveBitcodePath(bc)
				if bcPath != "" {
					bcFiles = append(bcFiles, bcPath)
				}
			}
		} else {
			LogDebug("\tskipping empty entry at index %v\n", index)
		}
	}

	LogInfo("bcFiles: %v\n", bcFiles)
	LogInfo("len(bcFiles) = %v\n", len(bcFiles))

	if len(bcFiles) > 0 {

		// Sort the bitcode files
		if ea.SortBitcodeFiles {
			LogWarning("Sorting bitcode files.")
			sort.Strings(bcFiles)
			sort.Strings(artifactFiles)
		}

		// Build archive
		if ea.BuildBitcodeModule {
			linkBitcodeFiles(ea, bcFiles)
		} else {
			archiveBcFiles(ea, bcFiles)
		}

		// Write manifest
		if ea.WriteManifest {
			writeManifest(ea, bcFiles, artifactFiles)
		}
	} else {
		LogError("No bitcode files found\n")
	}

}

func listArchiveFiles(ea extractionArgs, inputFile string) (contents []string) {
	var arArgs []string
	arArgs = append(arArgs, "-t")
	arArgs = append(arArgs, inputFile)
	output, err := runCmd(ea.ArchiverName, arArgs)
	if err != nil {
		LogWarning("ar command: %v %v", ea.ArchiverName, arArgs)
		LogFatal("Failed to extract contents from archive %s because: %v.\n", inputFile, err)
	}
	contents = strings.Split(output, "\n")
	return
}

func extractFile(ea extractionArgs, archive string, filename string, instance int) bool {
	var arArgs []string
	if runtime.GOOS != osDARWIN {
		arArgs = append(arArgs, "xN")
		arArgs = append(arArgs, strconv.Itoa(instance))
	} else {
		if instance > 1 {
			LogWarning("Cannot extract instance %v of %v from archive %s for instance > 1.\n", instance, filename, archive)
			return false
		}
		arArgs = append(arArgs, "x")
	}
	arArgs = append(arArgs, archive)
	arArgs = append(arArgs, filename)
	_, err := runCmd(ea.ArchiverName, arArgs)
	if err != nil {
		LogWarning("The archiver %v failed to extract instance %v of %v from archive %s because: %v.\n", ea.ArchiverName, instance, filename, archive, err)
		return false
	}
	return true
}

func fetchTOC(ea extractionArgs, inputFile string) map[string]int {
	toc := make(map[string]int)

	contents := listArchiveFiles(ea, inputFile)

	for _, item := range contents {
		//iam: this is a hack to make get-bc work on libcurl.a
		if item != "" && !strings.HasPrefix(item, "__.SYMDEF") {
			toc[item]++
		}
	}
	return toc
}

//handleArchive processes an archive, and creates either a bitcode archive, or a module, depending on the flags used.
//
//    Archives are strange beasts. handleArchive processes the archive by:
//
//      1. first creating a table of contents of the archive, which maps file names (in the archive) to the number of
//    times a file with that name is stored in the archive.
//
//      2. for each OCCURRENCE of a file (name and count) it extracts the section from the object file, and adds the
//    bitcode paths to the bitcode list.
//
//      3. it then either links all these bitcode files together using llvm-link,  or else is creates a bitcode
//    archive using llvm-ar
//
//iam: 5/1/2018
func handleArchive(ea extractionArgs) {
	// List bitcode files to link
	var bcFiles []string
	var artifactFiles []string

	inputFile, _ := filepath.Abs(ea.InputFile)

	LogInfo("handleArchive: extractionArgs = %v\n", ea)

	// Create tmp dir
	tmpDirName, err := ioutil.TempDir("", "gllvm")
	if err != nil {
		LogFatal("The temporary directory in which to extract object files could not be created.")
	}
	defer CheckDefer(func() error { return os.RemoveAll(tmpDirName) })

	homeDir, err := os.Getwd()
	if err != nil {
		LogFatal("Could not ascertain our whereabouts: %v", err)
	}

	err = os.Chdir(tmpDirName)
	if err != nil {
		LogFatal("Could not cd to %v because: %v", tmpDirName, err)
	}

	//1. fetch the Table of Contents
	toc := fetchTOC(ea, inputFile)

	LogDebug("Table of Contents of %v:\n%v\n", inputFile, toc)

	for obj, instance := range toc {
		for i := 1; i <= instance; i++ {

			if obj != "" && extractFile(ea, inputFile, obj, i) {

				artifacts := ea.Extractor(obj)
				LogInfo("\t%v\n", artifacts)
				artifactFiles = append(artifactFiles, artifacts...)
				for _, bc := range artifacts {
					bcPath := resolveBitcodePath(bc)
					if bcPath != "" {
						bcFiles = append(bcFiles, bcPath)
					}
				}
			}
		}
	}

	err = os.Chdir(homeDir)
	if err != nil {
		LogFatal("Could not cd to %v because: %v", homeDir, err)
	}

	LogDebug("handleArchive: walked %v\nartifactFiles:\n%v\nbcFiles:\n%v\n", tmpDirName, artifactFiles, bcFiles)

	if len(bcFiles) > 0 {

		// Sort the bitcode files
		if ea.SortBitcodeFiles {
			LogWarning("Sorting bitcode files.")
			sort.Strings(bcFiles)
			sort.Strings(artifactFiles)
		}

		// Build archive
		if ea.BuildBitcodeModule {
			linkBitcodeFiles(ea, bcFiles)
		} else {
			archiveBcFiles(ea, bcFiles)
		}

		// Write manifest
		if ea.WriteManifest {
			writeManifest(ea, bcFiles, artifactFiles)
		}
	} else {
		LogError("No bitcode files found\n")
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
		success, err := execCmd(ea.LlvmArchiverName, args, dir)
		LogInfo("ea.LlvmArchiverName = %s, args = %v, dir = %s\n", ea.LlvmArchiverName, args, dir)
		if !success {
			LogFatal("There was an error creating the bitcode archive: %v.\n", err)
		}
	}
	informUser("Built bitcode archive: %s.\n", ea.OutputFile)
}

func getsize(stringslice []string) (totalLength int) {
	totalLength = 0
	for _, s := range stringslice {
		totalLength += len(s)
	}
	return totalLength
}

func formatStdOut(stdout bytes.Buffer, usefulIndex int) string {
	infoArr := strings.Split(stdout.String(), "\n")[usefulIndex]
	ret := strings.Fields(infoArr)
	return ret[0]
}

func fetchArgMax(ea extractionArgs) (argMax int) {
	if ea.LinkArgSize == 0 {
		getArgMax := exec.Command("getconf", "ARG_MAX")
		var argMaxStr bytes.Buffer
		getArgMax.Stdout = &argMaxStr
		err := getArgMax.Run()
		if err != nil {
			LogError("getconf ARG_MAX failed with %s\n", err)
		}
		argMax, err = strconv.Atoi(formatStdOut(argMaxStr, 0))
		if err != nil {
			LogError("string conversion for argMax failed with %s\n", err)
		}
		argMax = int(0.9 * float32(argMax)) // keeping a comfort margin
	} else {
		argMax = ea.LinkArgSize
	}
	LogInfo("argMax = %v\n", argMax)
	return
}

func linkBitcodeFilesIncrementally(ea extractionArgs, filesToLink []string, argMax int, linkArgs []string) {
	var tmpFileList []string
	// Create tmp dir
	tmpDirName, err := ioutil.TempDir(".", "glinking")
	if err != nil {
		LogFatal("The temporary directory in which to put temporary linking files could not be created.")
	}
	if !ea.KeepTemp { // delete temporary folder after used unless told otherwise
		LogInfo("Temporary folder will be deleted")
		defer CheckDefer(func() error { return os.RemoveAll(tmpDirName) })
	} else {
		LogInfo("Keeping the temporary folder")
	}

	tmpFile, err := ioutil.TempFile(tmpDirName, "tmp")
	if err != nil {
		LogFatal("The temporary linking file could not be created.")
	}
	tmpFileList = append(tmpFileList, tmpFile.Name())
	linkArgs = append(linkArgs, "-o", tmpFile.Name())

	LogInfo("llvm-link argument size : %d", getsize(filesToLink))
	for _, file := range filesToLink {
		linkArgs = append(linkArgs, file)
		if getsize(linkArgs) > argMax {
			LogInfo("Linking command size exceeding system capacity : splitting the command")
			var success bool
			success, err = execCmd(ea.LlvmLinkerName, linkArgs, "")
			if !success || err != nil {
				LogFatal("There was an error linking input files into %s because %v, on file %s.\n", ea.OutputFile, err, file)
			}
			linkArgs = nil

			if ea.Verbose {
				linkArgs = append(linkArgs, "-v")
			}
			tmpFile, err = ioutil.TempFile(tmpDirName, "tmp")
			if err != nil {
				LogFatal("Could not generate a temp file in %s because %v.\n", tmpDirName, err)
			}
			tmpFileList = append(tmpFileList, tmpFile.Name())
			linkArgs = append(linkArgs, "-o", tmpFile.Name())
		}

	}
	success, err := execCmd(ea.LlvmLinkerName, linkArgs, "")
	if !success {
		LogFatal("There was an error linking input files into %s because %v.\n", tmpFile.Name(), err)
	}
	linkArgs = nil
	if ea.Verbose {
		linkArgs = append(linkArgs, "-v")
	}
	linkArgs = append(linkArgs, tmpFileList...)

	linkArgs = append(linkArgs, "-o", ea.OutputFile)

	success, err = execCmd(ea.LlvmLinkerName, linkArgs, "")
	if !success {
		LogFatal("There was an error linking input files into %s because %v.\n", ea.OutputFile, err)
	}
	LogInfo("Bitcode file extracted to: %s, from files %v \n", ea.OutputFile, tmpFileList)
}

func linkBitcodeFiles(ea extractionArgs, filesToLink []string) {
	var linkArgs []string
	// Extracting the command line max size from the environment if it is not specified
	argMax := fetchArgMax(ea)
	if ea.Verbose {
		linkArgs = append(linkArgs, "-v")
	}
	if getsize(filesToLink) > argMax { //command line size too large for the OS (necessitated by chromium)
		linkBitcodeFilesIncrementally(ea, filesToLink, argMax, linkArgs)
	} else {
		linkArgs = append(linkArgs, "-o", ea.OutputFile)
		linkArgs = append(linkArgs, filesToLink...)
		success, err := execCmd(ea.LlvmLinkerName, linkArgs, "")
		if !success {
			LogFatal("There was an error linking input files into %s because %v.\n", ea.OutputFile, err)
		}
		informUser("Bitcode file extracted to: %s.\n", ea.OutputFile)
	}
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
		LogWarning("Error reading the %s section of Mach-O file %s.", DarwinSectionName, inputFile)
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
		LogWarning("Failed to find the file %v\n", bcPath)
		return ""
	}
	return bcPath
}

func writeManifest(ea extractionArgs, bcFiles []string, artifactFiles []string) {
	manifestFilename := ea.OutputFile + ".llvm.manifest"
	//only go into the gory details if we have a store around.
	if LLVMBitcodeStorePath != "" {
		section1 := "Physical location of extracted files:\n" + strings.Join(bcFiles, "\n") + "\n\n"
		section2 := "Build-time location of extracted files:\n" + strings.Join(artifactFiles, "\n")
		contents := []byte(section1 + section2)
		if err := ioutil.WriteFile(manifestFilename, contents, 0644); err != nil {
			LogFatal("There was an error while writing the manifest file: ", err)
		}
	} else {
		contents := []byte("\n" + strings.Join(bcFiles, "\n") + "\n")
		if err := ioutil.WriteFile(manifestFilename, contents, 0644); err != nil {
			LogFatal("There was an error while writing the manifest file: ", err)
		}
	}
	informUser("Manifest file written to %s.\n", manifestFilename)
}
