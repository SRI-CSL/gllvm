package main

import (
    "os"
    "os/exec"
    "io/ioutil"
    "log"
    "strings"
    "path"
    "path/filepath"
    "runtime"
    "io"
)

func compile(args []string) {
    if len(args) < 1 {
        log.Fatal("You must precise which compiler to use.")
    }
    var compilerName = args[0]
    var compilerExecName = getCompilerExecName(compilerName)
    var configureOnly bool
    if os.Getenv(CONFIGURE_ONLY) != "" {
        configureOnly = true
    }
    args = args[1:]
    var pr = parse(args)

    // If configure only is set, try to execute normal compiling command then exit silently
    if configureOnly {
        execCompile(compilerExecName, pr)
        os.Exit(0)
    }
    // Else try to build objects and bitcode
    buildAndAttachBitcode(compilerExecName, pr)
}

// Compiles bitcode files and attach path to the object files
func buildAndAttachBitcode(compilerExecName string, pr ParserResult) {
    // If nothing to do, exit silently
    if len(pr.InputFiles) == 0 || pr.IsEmitLLVM || pr.IsAssembly || pr.IsAssembleOnly ||
        (pr.IsDependencyOnly && !pr.IsCompileOnly) || pr.IsPreprocessOnly {
        os.Exit(0)
    }

    var newObjectFiles []string
    var hidden = !pr.IsCompileOnly

    if len(pr.InputFiles) == 1 && pr.IsCompileOnly {
        var srcFile = pr.InputFiles[0]
        objFile, bcFile := getArtifactNames(pr, 0, hidden)
        buildObjectFile(compilerExecName, pr, srcFile, objFile)
        buildBitcodeFile(compilerExecName, pr, srcFile, bcFile)
        attachBitcodePathToObject(bcFile, objFile)
    } else {
        for i, srcFile := range pr.InputFiles {
            objFile, bcFile := getArtifactNames(pr, i, hidden)
            buildObjectFile(compilerExecName, pr, srcFile, objFile)
            if hidden {
                newObjectFiles = append(newObjectFiles, objFile)
            } else if strings.HasSuffix(srcFile, ".bc") {
                attachBitcodePathToObject(srcFile, objFile)
            } else {
                buildBitcodeFile(compilerExecName, pr, srcFile, bcFile)
                attachBitcodePathToObject(bcFile, objFile)
            }
        }
    }

    if !pr.IsCompileOnly {
        linkFiles(compilerExecName, pr, newObjectFiles)
    }
}

func attachBitcodePathToObject(bcFile, objFile string) {
    // We can only attach a bitcode path to certain file types
    switch filepath.Ext(objFile) {
    case
        ".o",
        ".lo",
        ".os",
        ".So",
        ".po":

        // Store bitcode path to temp file
        var absBcPath, _= filepath.Abs(bcFile)
        tmpContent := []byte(absBcPath+"\n")
        tmpFile, err := ioutil.TempFile("", "gowllvm")
        if err != nil {
            log.Fatal(err)
        }
        defer os.Remove(tmpFile.Name())
        if _, err := tmpFile.Write(tmpContent); err != nil {
            log.Fatal(err)
    	}
    	if err := tmpFile.Close(); err != nil {
    		log.Fatal(err)
    	}

        // Let's write the bitcode section
        var attachCmd string
        var attachCmdArgs []string
        if runtime.GOOS == "darwin" {
            attachCmd = "ld"
            attachCmdArgs = []string{"-r", "-keep_private_externs", objFile, "-sectcreate", DARWIN_SEGMENT_NAME, DARWIN_SECTION_NAME, tmpFile.Name(), "-o", objFile}
        } else {
            attachCmd = "objcopy"
            attachCmdArgs = []string{"objcopy", "--add-section", ELF_SECTION_NAME, "=", tmpFile.Name(), objFile}
        }

        // Run the attach command and ignore errors
        execCmd(attachCmd, attachCmdArgs)

        // Copy bitcode file to store, if necessary
        if bcStorePath := os.Getenv(BC_STORE_PATH); bcStorePath != "" {
            destFilePath := path.Join(bcStorePath, getHashedPath(absBcPath))
            in, _ := os.Open(absBcPath)
            defer in.Close()
            out, _ := os.Create(destFilePath)
            defer  out.Close()
            io.Copy(out, in)
            out.Sync()
        }
    }
}

func linkFiles(compilerExecName string, pr ParserResult, objFiles []string) {
    var outputFile = pr.OutputFilename
    if outputFile == "" {
        outputFile = "a.out"
    }
    args := append(pr.ObjectFiles, objFiles...)
    args = append(args, pr.LinkArgs...)
    args = append(args, "-o", outputFile)
    if execCmd(compilerExecName, args) {
        log.Fatal("Failed to link.")
    }
}

// Tries to build the specified source file to object
func buildObjectFile(compilerExecName string, pr ParserResult, srcFile string, objFile string) {
    args := pr.CompileArgs[:]
    args = append(args, srcFile, "-c", "-o", objFile)
    if execCmd(compilerExecName, args) {
        log.Fatal("Failed to build object file for ", srcFile)
    }
}

// Tries to build the specified source file to bitcode
func buildBitcodeFile(compilerExecName string, pr ParserResult, srcFile string, bcFile string) {
    args := pr.CompileArgs[:]
    args = append(args, "-emit-llvm", "-c", srcFile, "-o", bcFile)
    if execCmd(compilerExecName, args) {
        log.Fatal("Failed to build bitcode file for ", srcFile)
    }
}

// Tries to build object file
func execCompile(compilerExecName string, pr ParserResult) {
    if execCmd(compilerExecName, pr.InputList) {
        log.Fatal("Failed to execute compile command.")
    }
}

// Executes a command then returns true if there was an error
func execCmd(cmdExecName string, args []string) bool {
    cmd := exec.Command(cmdExecName, listToArgString(args))
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if cmd.Run() == nil {
        return false
    } else {
        return true
    }
}

// Joins a list of arguments to create a string
func listToArgString(argList []string) string {
    return strings.Join(argList, " ")
}

func getCompilerExecName(compilerName string) string {
    var compilerPath = os.Getenv(COMPILER_PATH)
    switch compilerName {
    case "clang":
        var clangName = os.Getenv(C_COMPILER_NAME)
        if clangName != "" {
            return compilerPath + clangName
        } else {
            return compilerPath + compilerName
        }
    case "clang++":
        var clangppName = os.Getenv(C_COMPILER_NAME)
        if clangppName != "" {
            return compilerPath + clangppName
        } else {
            return compilerPath + compilerName
        }
    default:
        log.Fatal("The compiler ", compilerName, " is not supported by this tool.")
        return ""
    }
}
