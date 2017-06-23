package main

import (
    "os"
    "os/exec"
    "log"
    "runtime"
    "path/filepath"
    "strings"
    "regexp"
    "encoding/hex"
)

type ExtractingArgs struct {
    InputFile string
    InputType int
    OutputFile string
    LinkerName string
    ArchiverName string
    ArArgs []string
    ObjectTypeInArchive int // Type of file that can be put into an archive
    Extractor func(ExtractingArgs) []string
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
        ea.LinkerName = ln
    }
    if an := os.Getenv(AR_NAME); an != "" {
        ea.ArchiverName = an
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
        ea.Extractor = extractSectionUnix
        //ea.Extractor = extractSectionDarwin
        ea.ArArgs = append(ea.ArArgs, "-x")
        if ea.IsVerbose {
            ea.ArArgs = append(ea.ArArgs, "-v")
        }
        ea.ObjectTypeInArchive = FT_MACH_OBJECT

    // Create output filename if not given
    if ea.OutputFile == "" {
        if ea.InputType == FT_ARCHIVE {
            if ea.IsBuildBitcodeArchive {
                ea.OutputFile = ea.InputFile + ".a.bc"
            } else {
                ea.OutputFile = ea.InputFile + ".bca"
            }
        } else {
            ea.OutputFile = ea.InputFile + ".bc"
        }
    }

    default:
        log.Fatal("Unsupported platform: ", platform)
    }

    return ea
}

func handleExecutable(ea ExtractingArgs) {
    filesToLink := ea.Extractor(ea)
    var _ = filesToLink

}

func handleArchive(_ ExtractingArgs) {
    // TODO
}

func extractSectionDarwin(ea ExtractingArgs) (contents []string) {
    cmd := exec.Command("otool", "-X", "-s", DARWIN_SEGMENT_NAME, DARWIN_SECTION_NAME, ea.InputFile)
    out, err := cmd.Output()
    if err != nil {
        log.Fatal("There was an error extracting the Gowllvm section from ", ea.InputFile, ". Make sure that the 'otool' command is installed.")
    }
    sectionLines := strings.Split(string(out), "\n")
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
    contents = strings.Split(strings.TrimSuffix(string(octets), "\n"), "\n")
    return
}

func extractSectionUnix(ea ExtractingArgs) (contents []string) {
    // TODO CORRECT -D to -w
    cmd := exec.Command("objcopy", "--dump-section", ELF_SECTION_NAME + "=/dev/stdout", ea.InputFile)
    out, err := cmd.Output()
    if err != nil {
        log.Fatal("There was an error reading the contents of ", ea.InputFile, ". Make sure that the 'objcopy' command is installed.")
    }
    contents = strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
    return
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
        log.Fatal("The type of the input file is not handled.")
    }

    return
}
