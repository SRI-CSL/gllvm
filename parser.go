package main

import(
    "fmt"
    "regexp"
    "runtime"
    "path"
    "path/filepath"
    "strings"
    "crypto/sha256"
)

type ParserResult struct {
    InputList []string
    InputFiles []string
    ObjectFiles []string
    OutputFilename string
    CompileArgs []string
    LinkArgs []string
    IsVerbose bool
    IsDependencyOnly bool
    IsPreprocessOnly bool
    IsAssembleOnly bool
    IsAssembly bool
    IsCompileOnly bool
    IsEmitLLVM bool
}

type FlagInfo struct {
    arity int
    handler func(string, []string)
}

func parse(argList []string) ParserResult {
    var pr = ParserResult{}
    pr.InputList = argList

    var argsExactMatches = map[string]FlagInfo{
        "-o": FlagInfo{1, pr.outputFileCallback},
        "-c": FlagInfo{0, pr.compileOnlyCallback},
        "-E": FlagInfo{0, pr.preprocessOnlyCallback},
        "-S": FlagInfo{0, pr.assembleOnlyCallback},

        "--verbose": FlagInfo{0, pr.verboseFlagCallback},
        "--param": FlagInfo{1, pr.defaultBinaryCallback},
        "-aux-info": FlagInfo{1, pr.defaultBinaryCallback},

        "--version": FlagInfo{0, pr.compileOnlyCallback},
        "-v": FlagInfo{0, pr.compileOnlyCallback},

        "-w": FlagInfo{0, pr.compileOnlyCallback},
        "-W": FlagInfo{0, pr.compileOnlyCallback},

        "-emit-llvm": FlagInfo{0, pr.emitLLVMCallback},

        "-pipe": FlagInfo{0, pr.compileUnaryCallback},
        "-undef": FlagInfo{0, pr.compileUnaryCallback},
        "-nostdinc": FlagInfo{0, pr.compileUnaryCallback},
        "-nostdinc++": FlagInfo{0, pr.compileUnaryCallback},
        "-Qunused-arguments": FlagInfo{0, pr.compileUnaryCallback},
        "-no-integrated-as": FlagInfo{0, pr.compileUnaryCallback},
        "-integrated-as": FlagInfo{0, pr.compileUnaryCallback},

        "-pthread": FlagInfo{0, pr.compileUnaryCallback},
        "-nostdlibinc": FlagInfo{0, pr.compileUnaryCallback},

        "-mno-omit-leaf-frame-pointer": FlagInfo{0, pr.compileUnaryCallback},
        "-maes": FlagInfo{0, pr.compileUnaryCallback},
        "-mno-aes": FlagInfo{0, pr.compileUnaryCallback},
        "-mavx": FlagInfo{0, pr.compileUnaryCallback},
        "-mno-avx": FlagInfo{0, pr.compileUnaryCallback},
        "-mcmodel=kernel": FlagInfo{0, pr.compileUnaryCallback},
        "-mno-red-zone": FlagInfo{0, pr.compileUnaryCallback},
        "-mmmx": FlagInfo{0, pr.compileUnaryCallback},
        "-mno-mmx": FlagInfo{0, pr.compileUnaryCallback},
        "-msse": FlagInfo{0, pr.compileUnaryCallback},
        "-mno-sse2": FlagInfo{0, pr.compileUnaryCallback},
        "-msse2": FlagInfo{0, pr.compileUnaryCallback},
        "-mno-sse3": FlagInfo{0, pr.compileUnaryCallback},
        "-msse3": FlagInfo{0, pr.compileUnaryCallback},
        "-mno-sse": FlagInfo{0, pr.compileUnaryCallback},
        "-msoft-float": FlagInfo{0, pr.compileUnaryCallback},
        "-m3dnow": FlagInfo{0, pr.compileUnaryCallback},
        "-mno-3dnow": FlagInfo{0, pr.compileUnaryCallback},
        "-m32": FlagInfo{0, pr.compileUnaryCallback},
        "-m64": FlagInfo{0, pr.compileUnaryCallback},
        "-mstackrealign": FlagInfo{0, pr.compileUnaryCallback},

        "-A": FlagInfo{1, pr.compileBinaryCallback},
        "-D": FlagInfo{1, pr.compileBinaryCallback},
        "-U": FlagInfo{1, pr.compileBinaryCallback},

        "-M"  : FlagInfo{0, pr.dependencyOnlyCallback},
        "-MM": FlagInfo{0, pr.dependencyOnlyCallback},
        "-MF": FlagInfo{1, pr.dependencyBinaryCallback},
        "-MG": FlagInfo{0, pr.dependencyOnlyCallback},
        "-MP": FlagInfo{0, pr.dependencyOnlyCallback},
        "-MT": FlagInfo{1, pr.dependencyBinaryCallback},
        "-MQ": FlagInfo{1, pr.dependencyBinaryCallback},
        "-MD": FlagInfo{0, pr.dependencyOnlyCallback},
        "-MMD": FlagInfo{0, pr.dependencyOnlyCallback},

        "-I": FlagInfo{1, pr.compileBinaryCallback},
        "-idirafter": FlagInfo{1, pr.compileBinaryCallback},
        "-include": FlagInfo{1, pr.compileBinaryCallback},
        "-imacros": FlagInfo{1, pr.compileBinaryCallback},
        "-iprefix": FlagInfo{1, pr.compileBinaryCallback},
        "-iwithprefix": FlagInfo{1, pr.compileBinaryCallback},
        "-iwithprefixbefore": FlagInfo{1, pr.compileBinaryCallback},
        "-isystem": FlagInfo{1, pr.compileBinaryCallback},
        "-isysroot": FlagInfo{1, pr.compileBinaryCallback},
        "-iquote": FlagInfo{1, pr.compileBinaryCallback},
        "-imultilib": FlagInfo{1, pr.compileBinaryCallback},

        "-ansi": FlagInfo{0, pr.compileUnaryCallback},
        "-pedantic": FlagInfo{0, pr.compileUnaryCallback},
        "-x": FlagInfo{1, pr.compileBinaryCallback},

        "-g": FlagInfo{0, pr.compileUnaryCallback},
        "-g0": FlagInfo{0, pr.compileUnaryCallback},
        "-ggdb": FlagInfo{0, pr.compileUnaryCallback},
        "-ggdb3": FlagInfo{0, pr.compileUnaryCallback},
        "-gdwarf-2": FlagInfo{0, pr.compileUnaryCallback},
        "-gdwarf-3": FlagInfo{0, pr.compileUnaryCallback},
        "-gline-tables-only": FlagInfo{0, pr.compileUnaryCallback},

        "-p": FlagInfo{0, pr.compileUnaryCallback},
        "-pg": FlagInfo{0, pr.compileUnaryCallback},

        "-O": FlagInfo{0, pr.compileUnaryCallback},
        "-O0": FlagInfo{0, pr.compileUnaryCallback},
        "-O1": FlagInfo{0, pr.compileUnaryCallback},
        "-O2": FlagInfo{0, pr.compileUnaryCallback},
        "-O3": FlagInfo{0, pr.compileUnaryCallback},
        "-Os": FlagInfo{0, pr.compileUnaryCallback},
        "-Ofast": FlagInfo{0, pr.compileUnaryCallback},
        "-Og": FlagInfo{0, pr.compileUnaryCallback},

        "-Xclang": FlagInfo{1, pr.compileBinaryCallback},
        "-Xpreprocessor": FlagInfo{1, pr.defaultBinaryCallback},
        "-Xassembler": FlagInfo{1, pr.defaultBinaryCallback},
        "-Xlinker": FlagInfo{1, pr.defaultBinaryCallback},

        "-l": FlagInfo{1, pr.linkBinaryCallback},
        "-L": FlagInfo{1, pr.linkBinaryCallback},
        "-T": FlagInfo{1, pr.linkBinaryCallback},
        "-u": FlagInfo{1, pr.linkBinaryCallback},

        "-e": FlagInfo{1, pr.linkBinaryCallback},
        "-rpath": FlagInfo{1, pr.linkBinaryCallback},

        "-shared": FlagInfo{0, pr.linkUnaryCallback},
        "-static": FlagInfo{0, pr.linkUnaryCallback},
        "-pie": FlagInfo{0, pr.linkUnaryCallback},
        "-nostdlib": FlagInfo{0, pr.linkUnaryCallback},
        "-nodefaultlibs": FlagInfo{0, pr.linkUnaryCallback},
        "-rdynamic": FlagInfo{0, pr.linkUnaryCallback},

        "-dynamiclib": FlagInfo{0, pr.linkUnaryCallback},
        "-current_version": FlagInfo{1, pr.linkBinaryCallback},
        "-compatibility_version": FlagInfo{1, pr.linkBinaryCallback},

        "-print-multi-directory": FlagInfo{0, pr.compileUnaryCallback},
        "-print-multi-lib": FlagInfo{0, pr.compileUnaryCallback},
        "-print-libgcc-file-name": FlagInfo{0, pr.compileUnaryCallback},

        "-fprofile-arcs": FlagInfo{0, pr.compileLinkUnaryCallback},
        "-coverage": FlagInfo{0, pr.compileLinkUnaryCallback},
        "--coverage": FlagInfo{0, pr.compileLinkUnaryCallback},

        "-Wl,-dead_strip": FlagInfo{0, pr.darwinWarningLinkUnaryCallback},
    }

    var argPatterns = map[string]FlagInfo{
        `^.+\.(c|cc|cpp|C|cxx|i|s|S|bc)$`: FlagInfo{0, pr.inputFileCallback},
        `^.+\.([fF](|[0-9][0-9]|or|OR|pp|PP))$`: FlagInfo{0, pr.inputFileCallback},
        `^.+\.(o|lo|So|so|po|a|dylib)$`: FlagInfo{0, pr.objectFileCallback},
        `^.+\.dylib(\.\d)+$`: FlagInfo{0, pr.objectFileCallback},
        `^.+\.(So|so)(\.\d)+$`: FlagInfo{0, pr.objectFileCallback},
        `^-(l|L).+$`: FlagInfo{0, pr.linkUnaryCallback},
        `^-I.+$`: FlagInfo{0, pr.compileUnaryCallback},
        `^-D.+$`: FlagInfo{0, pr.compileUnaryCallback},
        `^-U.+$`: FlagInfo{0, pr.compileUnaryCallback},
        `^-Wl,.+$`: FlagInfo{0, pr.linkUnaryCallback},
        `^-W.*$`: FlagInfo{0, pr.compileUnaryCallback},
        `^-f.+$`: FlagInfo{0, pr.compileUnaryCallback},
        `^-rtlib=.+$`: FlagInfo{0, pr.linkUnaryCallback},
        `^-std=.+$`: FlagInfo{0, pr.compileUnaryCallback},
        `^-stdlib=.+$`: FlagInfo{0, pr.compileLinkUnaryCallback},
        `^-mtune=.+$`: FlagInfo{0, pr.compileUnaryCallback},
        `^--sysroot=.+$`: FlagInfo{0, pr.compileUnaryCallback},
        `^-print-prog-name=.*$`: FlagInfo{0, pr.compileUnaryCallback},
        `^-print-file-name=.*$`: FlagInfo{0, pr.compileUnaryCallback},
    }

    for len(argList) > 0 && !(pr.IsAssembly || pr.IsAssembleOnly || pr.IsPreprocessOnly) {
        var elem = argList[0]

        // Try to match the flag exactly
        if fi, ok := argsExactMatches[elem]; ok {
            fi.handler(elem, argList[1:1+fi.arity])
            argList = argList[1+fi.arity:]
        // Else try to match a pattern
        } else {
            var listShift = 0
            for pattern, fi := range argPatterns {
                var regExp = regexp.MustCompile(pattern)
                if regExp.MatchString(elem) {
                    fi.handler(elem, argList[1:1+fi.arity])
                    listShift = fi.arity
                    break
                }
            }
            argList = argList[1+listShift:]
        }

    }
    fmt.Println(pr)
    return pr
}

// Return the object and bc filenames that correspond to the i-th source file
func getArtifactNames(pr ParserResult, srcFileIndex int, hidden bool) (objBase string, bcBase string) {
    if len(pr.InputFiles) == 1 && pr.IsCompileOnly && len(pr.OutputFilename) > 0 {
        objBase = pr.OutputFilename
        dir, baseName := path.Split(objBase)
        bcBaseName := fmt.Sprintf(".%s.bc", baseName)
        bcBase = path.Join(dir, bcBaseName)
    } else {
        srcFile := pr.InputFiles[srcFileIndex]
        var dir, baseNameWithExt = path.Split(srcFile)
        var baseName = strings.TrimSuffix(baseNameWithExt, filepath.Ext(baseNameWithExt))
        bcBase = fmt.Sprintf(".%s.o.bc", baseName)
        if hidden {
            objBase = path.Join(dir, fmt.Sprintf(".%s.o", baseName))
        } else {
            objBase = path.Join(dir, fmt.Sprintf("%s.o", baseName))
        }
    }
    return
}

// Return a hash for the absolute object path
func getHashedPath(path string) string {
    inputBytes := []byte(path)
    hash := sha256.Sum256(inputBytes)
    return string(hash[:])
}

func (pr *ParserResult) inputFileCallback(flag string, _ []string) {
    var regExp = regexp.MustCompile(`\\.(s|S)$`)
    pr.InputFiles = append(pr.InputFiles, flag)
    if regExp.MatchString(flag) {
        pr.IsAssembly = true
    }
}

func (pr *ParserResult) outputFileCallback(_ string, args []string) {
    pr.OutputFilename = args[0]
}

func (pr *ParserResult) objectFileCallback(flag string, _ []string) {
    pr.ObjectFiles = append(pr.ObjectFiles, flag)
}

func (pr *ParserResult) preprocessOnlyCallback(_ string, _ []string) {
    pr.IsPreprocessOnly = true
}

func (pr *ParserResult) dependencyOnlyCallback(flag string, _ []string) {
    pr.IsDependencyOnly = true
    pr.CompileArgs = append(pr.CompileArgs, flag)
}

func (pr *ParserResult) assembleOnlyCallback(_ string, _ []string) {
    pr.IsAssembleOnly = true
}

func (pr *ParserResult) verboseFlagCallback(_ string, _ []string) {
    pr.IsVerbose = true
}

func (pr *ParserResult) compileOnlyCallback(_ string, _ []string) {
    pr.IsCompileOnly = true
}

func (pr *ParserResult) emitLLVMCallback(_ string, _ []string) {
    pr.IsCompileOnly = true
    pr.IsEmitLLVM = true
}

func (pr *ParserResult) linkUnaryCallback(flag string, _ []string) {
    pr.LinkArgs = append(pr.LinkArgs, flag)
}

func (pr *ParserResult) compileUnaryCallback(flag string, _ []string) {
    pr.CompileArgs = append(pr.CompileArgs, flag)
}

func (pr *ParserResult) darwinWarningLinkUnaryCallback(flag string, _ []string) {
    if runtime.GOOS == "darwin" {
        fmt.Println("The flag", flag, "cannot be used with this tool.")
    } else {
        pr.LinkArgs = append(pr.LinkArgs, flag)
    }
}

func (_ *ParserResult) defaultBinaryCallback(_ string, _ []string) {
    // Do nothing
}

func (pr *ParserResult) dependencyBinaryCallback(flag string, args []string) {
    pr.CompileArgs = append(pr.CompileArgs, flag, args[0])
    pr.IsDependencyOnly = true
}

func (pr *ParserResult) compileBinaryCallback(flag string, args []string) {
    pr.CompileArgs = append(pr.CompileArgs, flag, args[0])
}

func (pr *ParserResult) linkBinaryCallback(flag string, args []string) {
    pr.LinkArgs = append(pr.LinkArgs, flag, args[0])
}

func (pr *ParserResult) compileLinkUnaryCallback(flag string, _ []string) {
    pr.LinkArgs = append(pr.LinkArgs, flag)
    pr.CompileArgs = append(pr.CompileArgs, flag)
}
