package main

import (
	"debug/elf"
	"debug/macho"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type extractionArgs struct {
	InputFile             string
	InputType             int
	OutputFile            string
	LinkerName            string
	ArchiverName          string
	ArArgs                []string
	ObjectTypeInArchive   int // Type of file that can be put into an archive
	Extractor             func(string) []string
	IsVerbose             bool
	IsWriteManifest       bool
	IsBuildBitcodeArchive bool
}

func extract(args []string) {
	ea := parseExtractionArgs(args)

	switch ea.InputType {
	case ftELFEXECUTABLE,
		ftELFSHARED,
		ftELFOBJECT,
		ftMACHEXECUTABLE,
		ftMACHSHARED,
		ftMACHOBJECT:
		handleExecutable(ea)
	case ftARCHIVE:
		handleArchive(ea)
	default:
		log.Fatal("Incorrect input file type.")
	}

}

func parseExtractionArgs(args []string) extractionArgs {
	// Initializing args to defaults
	ea := extractionArgs{
		LinkerName:   "llvm-link",
		ArchiverName: "llvm-ar",
	}

	// Checking environment variables
	if LLVMLINKName != "" {
		//FIXME: this check can be eliminated because filepath.Join("", "baz") = "baz"
		if LLVMToolChainBinDir != "" {
			ea.LinkerName = filepath.Join(LLVMToolChainBinDir, LLVMLINKName)
		} else {
			ea.LinkerName = LLVMLINKName
		}
	}
	if LLVMARName != "" {
		//FIXME: this check can be eliminated because filepath.Join("", "baz") = "baz"
		if LLVMToolChainBinDir != "" {
			ea.ArchiverName = filepath.Join(LLVMToolChainBinDir, LLVMARName)
		} else {
			ea.ArchiverName = LLVMARName
		}
	}

	// Parsing cli input. FIXME:  "get-bc -mb libfoo.a" should work just like "get-bc -m -b libfoo.a"
	for len(args) > 0 {
		switch arg := args[0]; arg {
		case "-b":
			ea.IsBuildBitcodeArchive = true
			args = args[1:]
		case "-v":
			ea.IsVerbose = true
			args = args[1:]
		case "-m":
			ea.IsWriteManifest = true
			args = args[1:]
		case "-o":
			if len(args) < 2 {
				log.Fatal("There was an error parsing the arguments.")
			}
			ea.OutputFile = args[1]
			args = args[2:]
		default:
			ea.InputFile = arg
			args = args[1:]
		}
	}

	// Sanity-check the parsed arguments
	if len(ea.InputFile) == 0 {
		log.Fatal("No input file was given.")
	}
	if _, err := os.Stat(ea.InputFile); os.IsNotExist(err) {
		log.Fatal("The input file ", ea.InputFile, " does not exist.")
	}
	realPath, err := filepath.EvalSymlinks(ea.InputFile)
	if err != nil {
		log.Fatal("There was an error getting the real path of ", ea.InputFile, ".")
	}
	ea.InputFile = realPath
	ea.InputType = getFileType(realPath)

	// Set arguments according to runtime OS
	switch platform := runtime.GOOS; platform {
	case "freebsd", "linux":
		ea.Extractor = extractSectionUnix
		if ea.IsVerbose {
			ea.ArArgs = append(ea.ArArgs, "xv")
		} else {
			ea.ArArgs = append(ea.ArArgs, "x")
		}
		ea.ObjectTypeInArchive = ftELFOBJECT
	case "darwin":
		ea.Extractor = extractSectionDarwin
		ea.ArArgs = append(ea.ArArgs, "-x")
		if ea.IsVerbose {
			ea.ArArgs = append(ea.ArArgs, "-v")
		}
		ea.ObjectTypeInArchive = ftMACHOBJECT
	default:
		log.Fatal("Unsupported platform: ", platform)
	}

	// Create output filename if not given
	if ea.OutputFile == "" {
		if ea.InputType == ftARCHIVE {
			var ext string
			if ea.IsBuildBitcodeArchive {
				ext = ".a.bc"
			} else {
				ext = ".bca"
			}
			ea.OutputFile = strings.TrimSuffix(ea.InputFile, ".a") + ext
		} else {
			ea.OutputFile = ea.InputFile + ".bc"
		}
	}

	return ea
}

func handleExecutable(ea extractionArgs) {
	artifactPaths := ea.Extractor(ea.InputFile)
	filesToLink := make([]string, len(artifactPaths))
	for i, artPath := range artifactPaths {
		filesToLink[i] = resolveBitcodePath(artPath)
	}
	extractTimeLinkFiles(ea, filesToLink)

	// Write manifest
	if ea.IsWriteManifest {
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
		log.Fatal("The temporary directory in which to extract object files could not be created.")
	}
	defer os.RemoveAll(tmpDirName)

	// Extract objects to tmpDir
	arArgs := ea.ArArgs
	inputAbsPath, _ := filepath.Abs(ea.InputFile)
	arArgs = append(arArgs, inputAbsPath)
	success, err := execCmd("ar", arArgs, tmpDirName)
	if !success {
		logFatal("Failed to extract object files from %s to %s because: %v.\n", ea.InputFile, tmpDirName, err)
	}

	// Define object file handling closure
	var walkHandlingFunc = func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			ft := getFileType(path)
			if ft == ea.ObjectTypeInArchive {
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
	if ea.IsBuildBitcodeArchive {
		extractTimeLinkFiles(ea, bcFiles)
	} else {
		archiveBcFiles(ea, bcFiles)
	}

	// Write manifest
	if ea.IsWriteManifest {
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
		if !success {
			logFatal("There was an error creating the bitcode archive: %v.\n", err)
		}
	}
	fmt.Println("Built bitcode archive", ea.OutputFile)
}

func extractTimeLinkFiles(ea extractionArgs, filesToLink []string) {
	var linkArgs []string
	if ea.IsVerbose {
		linkArgs = append(linkArgs, "-v")
	}
	linkArgs = append(linkArgs, "-o", ea.OutputFile)
	linkArgs = append(linkArgs, filesToLink...)
	success, err := execCmd(ea.LinkerName, linkArgs, "")
	if !success {
		log.Fatal("There was an error linking input files into %s because %v.\n", ea.OutputFile, err)
	}
	fmt.Println("Bitcode file extracted to", ea.OutputFile)
}

func extractSectionDarwin(inputFile string) (contents []string) {
	machoFile, err := macho.Open(inputFile)
	if err != nil {
		log.Fatal("Mach-O file ", inputFile, " could not be read.")
	}
	section := machoFile.Section(darwinSECTIONNAME)
	sectionContents, errContents := section.Data()
	if errContents != nil {
		log.Fatal("Error reading the ", darwinSECTIONNAME, " section of Mach-O file ", inputFile, ".")
	}
	contents = strings.Split(strings.TrimSuffix(string(sectionContents), "\n"), "\n")
	return
}

func extractSectionUnix(inputFile string) (contents []string) {
	elfFile, err := elf.Open(inputFile)
	if err != nil {
		log.Fatal("ELF file ", inputFile, " could not be read.")
	}
	section := elfFile.Section(elfSECTIONNAME)
	sectionContents, errContents := section.Data()
	if errContents != nil {
		log.Fatal("Error reading the ", elfSECTIONNAME, " section of ELF file ", inputFile, ".")
	}
	contents = strings.Split(strings.TrimSuffix(string(sectionContents), "\n"), "\n")
	return
}

// Return the actual path to the bitcode file, or an empty string if it does not exist
func resolveBitcodePath(bcPath string) string {
	if _, err := os.Stat(bcPath); os.IsNotExist(err) {
		// If the bitcode file does not exist, try to find it in the store
		if bcStorePath := os.Getenv(BitcodeStorePath); bcStorePath != "" {
			// Compute absolute path hash
			absBcPath, _ := filepath.Abs(bcPath)
			storeBcPath := path.Join(bcStorePath, getHashedPath(absBcPath))
			if _, err := os.Stat(storeBcPath); os.IsNotExist(err) {
				return ""
			}
			return storeBcPath
		}
		return ""
	}
	return bcPath
}

func getFileType(realPath string) (fileType int) {
	// We need the file command to guess the file type
	cmd := exec.Command("file", realPath)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal("There was an error getting the type of ", realPath, ". Make sure that the 'file' command is installed.")
	}

	// Test the output
	if fo := string(out); strings.Contains(fo, "ELF") && strings.Contains(fo, "executable") {
		fileType = ftELFEXECUTABLE
	} else if strings.Contains(fo, "Mach-O") && strings.Contains(fo, "executable") {
		fileType = ftMACHEXECUTABLE
	} else if strings.Contains(fo, "ELF") && strings.Contains(fo, "shared") {
		fileType = ftELFSHARED
	} else if strings.Contains(fo, "Mach-O") && strings.Contains(fo, "dynamically linked shared") {
		fileType = ftMACHSHARED
	} else if strings.Contains(fo, "current ar archive") {
		fileType = ftARCHIVE
	} else if strings.Contains(fo, "ELF") && strings.Contains(fo, "relocatable") {
		fileType = ftELFOBJECT
	} else if strings.Contains(fo, "Mach-O") && strings.Contains(fo, "object") {
		fileType = ftMACHOBJECT
	} else {
		fileType = ftUNDEFINED
	}

	return
}

func writeManifest(ea extractionArgs, bcFiles []string, artifactFiles []string) {
	section1 := "Physical location of extracted files:\n" + strings.Join(bcFiles, "\n") + "\n\n"
	section2 := "Build-time location of extracted files:\n" + strings.Join(artifactFiles, "\n")
	contents := []byte(section1 + section2)
	manifestFilename := ea.OutputFile + ".llvm.manifest"
	if err := ioutil.WriteFile(manifestFilename, contents, 0644); err != nil {
		log.Fatal("There was an error while writing the manifest file: ", err)
	}
	fmt.Println("Manifest file written to", manifestFilename)
}
