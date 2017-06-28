package main

import (
    "os"
    "os/exec"
    "log"
    "runtime"
    "path"
    "path/filepath"
    "strings"
    "regexp"
    "encoding/hex"
    "io/ioutil"
    "fmt"
    "reflect"
)

type ExtractingArgs struct {
    InputFile string
    InputType int
    OutputFile string
    LinkerName string
    ArchiverName string
    ArArgs []string
    ObjectTypeInArchive int // Type of file that can be put into an archive
    Extractor func(string) string
    IsVerbose bool
    IsWriteManifest bool
    IsBuildBitcodeArchive bool
}

func extract(args []string) {
    ea := parseExtractingArgs(args)

    switch ea.InputType {
        case FT_ELF_EXECUTABLE,
            FT_ELF_SHARED,
            FT_ELF_OBJECT,
            FT_MACH_EXECUTABLE,
            FT_MACH_SHARED,
            FT_MACH_OBJECT:
            handleExecutable(ea)
        case FT_ARCHIVE:
            handleArchive(ea)
        default:
            log.Fatal("Incorrect input file type.")
    }

}



func parseExtractingArgs(args []string) ExtractingArgs {
    // Initializing args to defaults
    ea := ExtractingArgs{
        LinkerName: "llvm-link",
        ArchiverName: "llvm-ar",
    }

    // Checking environment variables
    if ln := os.Getenv(LINKER_NAME); ln != "" {
        if toolsPath := os.Getenv(TOOLS_PATH); toolsPath != "" {
            ea.LinkerName = toolsPath + ln
        } else {
            ea.LinkerName = ln
        }
    }
    if an := os.Getenv(AR_NAME); an != "" {
        if toolsPath := os.Getenv(TOOLS_PATH); toolsPath != "" {
            ea.ArchiverName = toolsPath + an
        } else {
            ea.ArchiverName = an
        }
    }

    // Parsing cli input
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
        ea.ObjectTypeInArchive = FT_ELF_OBJECT
    case "darwin":
        ea.Extractor = extractSectionDarwin
        ea.ArArgs = append(ea.ArArgs, "-x")
        if ea.IsVerbose {
            ea.ArArgs = append(ea.ArArgs, "-v")
        }
        ea.ObjectTypeInArchive = FT_MACH_OBJECT
    default:
        log.Fatal("Unsupported platform: ", platform)
    }

    // Create output filename if not given
    if ea.OutputFile == "" {
        if ea.InputType == FT_ARCHIVE {
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

func handleExecutable(ea ExtractingArgs) {
    artifactPath := ea.Extractor(ea.InputFile)
    filesToLink := []string{resolveBitcodePath(artifactPath)}
    extractTimeLinkFiles(ea, filesToLink)

    // Write manifest
    if ea.IsWriteManifest {
        writeManifest(ea, filesToLink, []string{artifactPath})
    }
}

func handleArchive(ea ExtractingArgs) {
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
    if execCmd("ar", arArgs, tmpDirName) {
        log.Fatal("Failed to extract object files from ", ea.InputFile, " to ", tmpDirName, ".")
    }

    // Define object file handling closure
    var walkHandlingFunc = func(path string, info os.FileInfo, err error) error {
        if err == nil && !info.IsDir() {
            ft := getFileType(path)
            if ft == ea.ObjectTypeInArchive {
                artifactPath := ea.Extractor(path)
                bcPath := resolveBitcodePath(artifactPath)
                bcFiles = append(bcFiles, bcPath)
                artifactFiles = append(artifactFiles, artifactPath)
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

func archiveBcFiles(ea ExtractingArgs, bcFiles []string) {
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
        if execCmd(ea.ArchiverName, args, dir) {
            log.Fatal("There was an error creating the bitcode archive.")
        }
    }
    fmt.Println("Built bitcode archive", ea.OutputFile)
}

func extractTimeLinkFiles(ea ExtractingArgs, filesToLink []string) {
    var linkArgs []string
    if ea.IsVerbose {
        linkArgs = append(linkArgs, "-v")
    }
    linkArgs = append(linkArgs, "-o", ea.OutputFile)
    linkArgs = append(linkArgs, filesToLink...)
    if execCmd(ea.LinkerName, linkArgs, "") {
        log.Fatal("There was an error linking input files into ", ea.OutputFile, ".")
    }
    fmt.Println("Bitcode file extracted to", ea.OutputFile)
}

func extractSectionDarwin(inputFile string) (contents string) {
    cmd := exec.Command("otool", "-X", "-s", DARWIN_SEGMENT_NAME, DARWIN_SECTION_NAME, inputFile)
    out, err := cmd.Output()
    if err != nil {
        log.Fatal("There was an error extracting the gllvm section from ", inputFile, ". Make sure that the 'otool' command is installed.")
    }
	//FIXME: ok this looks all wrong. need to deal with out properly and the go way, whatever that may be.
	logDebug("extractSectionDarwin: T(out) = %v \n", reflect.TypeOf(out))
	logDebug("extractSectionDarwin: out = %v \n", out)
    sectionLines := strings.Split(string(out), "\n")
	logDebug("extractSectionDarwin: sectionLines = %v\n", out)
    regExp := regexp.MustCompile(`^(?:[0-9a-f]{8,16}\t)?([0-9a-f\s]+)$`)
    var octets []byte

    for _, line := range sectionLines {
        if matches := regExp.FindStringSubmatch(line); matches != nil {
            hexline := []byte(strings.Join(strings.Split(matches[1], " "), ""))
            dst := make([]byte, hex.DecodedLen(len(hexline)))
            hex.Decode(dst, hexline)
            octets = append(octets, dst...)
        }
    }
    contents = strings.TrimSuffix(string(octets), "\n")
    return
}

func extractSectionUnix(inputFile string) (contents string) {
    cmd := exec.Command("objcopy", "--dump-section", ELF_SECTION_NAME + "=/dev/stdout", inputFile)
    out, err := cmd.Output()
    if err != nil {
        log.Fatal("There was an error reading the contents of ", inputFile, ". Make sure that the 'objcopy' command is installed.")
    }
    contents = strings.TrimSuffix(string(out), "\n")
    return
}

// Return the actual path to the bitcode file, or an empty string if it does not exist
func resolveBitcodePath(bcPath string) string {
    if _, err := os.Stat(bcPath); os.IsNotExist(err) {
        // If the bitcode file does not exist, try to find it in the store
        if bcStorePath := os.Getenv(BC_STORE_PATH); bcStorePath != "" {
            // Compute absolute path hash
            absBcPath, _ := filepath.Abs(bcPath)
            storeBcPath := path.Join(bcStorePath, getHashedPath(absBcPath))
            if _, err := os.Stat(storeBcPath); os.IsNotExist(err) {
                return ""
            } else {
                return storeBcPath
            }
        } else {
            return ""
        }
    } else {
        return bcPath
    }
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
        fileType = FT_ELF_EXECUTABLE
    } else if strings.Contains(fo, "Mach-O") && strings.Contains(fo, "executable") {
        fileType = FT_MACH_EXECUTABLE
    } else if strings.Contains(fo, "ELF") && strings.Contains(fo, "shared") {
        fileType = FT_ELF_SHARED
    } else if strings.Contains(fo, "Mach-O") && strings.Contains(fo, "dynamically linked shared") {
        fileType = FT_MACH_SHARED
    } else if strings.Contains(fo, "current ar archive") {
        fileType = FT_ARCHIVE
    } else if strings.Contains(fo, "ELF") && strings.Contains(fo, "relocatable") {
        fileType = FT_ELF_OBJECT
    } else if strings.Contains(fo, "Mach-O") && strings.Contains(fo, "object") {
        fileType = FT_MACH_OBJECT
    } else {
        fileType = FT_UNDEFINED
    }

    return
}

func writeManifest(ea ExtractingArgs, bcFiles []string, artifactFiles []string) {
    section1 := "Physical location of extracted files:\n" + strings.Join(bcFiles, "\n") + "\n\n"
    section2 := "Build-time location of extracted files:\n" + strings.Join(artifactFiles, "\n")
    contents := []byte(section1 + section2)
    manifestFilename := ea.OutputFile + ".llvm.manifest"
    if err := ioutil.WriteFile(manifestFilename, contents, 0644); err != nil {
        log.Fatal("There was an error while writing the manifest file: ", err)
    }
    fmt.Println("Manifest file written to", manifestFilename)
}
