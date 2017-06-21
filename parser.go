package main

import(
    "fmt"
    "regexp"
    "runtime"
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
    handler func(ParserResult, string, []string)
}

func parse(argList []string) ParserResult {
    var argsExactMatches = map[string]FlagInfo{
        "-o": FlagInfo{1, outputFileCallback},
        "-c": FlagInfo{0, compileOnlyCallback},
        "-E": FlagInfo{0, preprocessOnlyCallback},
        "-S": FlagInfo{0, assembleOnlyCallback},

        "--verbose": FlagInfo{0, verboseFlagCallback},
        "--param": FlagInfo{1, defaultBinaryCallback},
        "-aux-info": FlagInfo{1, defaultBinaryCallback},

        "--version": FlagInfo{0, compileOnlyCallback},
        "-v": FlagInfo{0, compileOnlyCallback},

        "-w": FlagInfo{0, compileOnlyCallback},
        "-W": FlagInfo{0, compileOnlyCallback},

        "-emit-llvm": FlagInfo{0, emitLLVMCallback},

        "-pipe": FlagInfo{0, compileUnaryCallback},
        "-undef": FlagInfo{0, compileUnaryCallback},
        "-nostdinc": FlagInfo{0, compileUnaryCallback},
        "-nostdinc++": FlagInfo{0, compileUnaryCallback},
        "-Qunused-arguments": FlagInfo{0, compileUnaryCallback},
        "-no-integrated-as": FlagInfo{0, compileUnaryCallback},
        "-integrated-as": FlagInfo{0, compileUnaryCallback},

        "-pthread": FlagInfo{0, compileUnaryCallback},
        "-nostdlibinc": FlagInfo{0, compileUnaryCallback},

        "-mno-omit-leaf-frame-pointer": FlagInfo{0, compileUnaryCallback},
        "-maes": FlagInfo{0, compileUnaryCallback},
        "-mno-aes": FlagInfo{0, compileUnaryCallback},
        "-mavx": FlagInfo{0, compileUnaryCallback},
        "-mno-avx": FlagInfo{0, compileUnaryCallback},
        "-mcmodel=kernel": FlagInfo{0, compileUnaryCallback},
        "-mno-red-zone": FlagInfo{0, compileUnaryCallback},
        "-mmmx": FlagInfo{0, compileUnaryCallback},
        "-mno-mmx": FlagInfo{0, compileUnaryCallback},
        "-msse": FlagInfo{0, compileUnaryCallback},
        "-mno-sse2": FlagInfo{0, compileUnaryCallback},
        "-msse2": FlagInfo{0, compileUnaryCallback},
        "-mno-sse3": FlagInfo{0, compileUnaryCallback},
        "-msse3": FlagInfo{0, compileUnaryCallback},
        "-mno-sse": FlagInfo{0, compileUnaryCallback},
        "-msoft-float": FlagInfo{0, compileUnaryCallback},
        "-m3dnow": FlagInfo{0, compileUnaryCallback},
        "-mno-3dnow": FlagInfo{0, compileUnaryCallback},
        "-m32": FlagInfo{0, compileUnaryCallback},
        "-m64": FlagInfo{0, compileUnaryCallback},
        "-mstackrealign": FlagInfo{0, compileUnaryCallback},

        "-A": FlagInfo{1, compileBinaryCallback},
        "-D": FlagInfo{1, compileBinaryCallback},
        "-U": FlagInfo{1, compileBinaryCallback},

        "-M"  : FlagInfo{0, dependencyOnlyCallback},
        "-MM": FlagInfo{0, dependencyOnlyCallback},
        "-MF": FlagInfo{1, dependencyBinaryCallback},
        "-MG": FlagInfo{0, dependencyOnlyCallback},
        "-MP": FlagInfo{0, dependencyOnlyCallback},
        "-MT": FlagInfo{1, dependencyBinaryCallback},
        "-MQ": FlagInfo{1, dependencyBinaryCallback},
        "-MD": FlagInfo{0, dependencyOnlyCallback},
        "-MMD": FlagInfo{0, dependencyOnlyCallback},

        "-I": FlagInfo{1, compileBinaryCallback},
        "-idirafter": FlagInfo{1, compileBinaryCallback},
        "-include": FlagInfo{1, compileBinaryCallback},
        "-imacros": FlagInfo{1, compileBinaryCallback},
        "-iprefix": FlagInfo{1, compileBinaryCallback},
        "-iwithprefix": FlagInfo{1, compileBinaryCallback},
        "-iwithprefixbefore": FlagInfo{1, compileBinaryCallback},
        "-isystem": FlagInfo{1, compileBinaryCallback},
        "-isysroot": FlagInfo{1, compileBinaryCallback},
        "-iquote": FlagInfo{1, compileBinaryCallback},
        "-imultilib": FlagInfo{1, compileBinaryCallback},

        "-ansi": FlagInfo{0, compileUnaryCallback},
        "-pedantic": FlagInfo{0, compileUnaryCallback},
        "-x": FlagInfo{1, compileBinaryCallback},

        "-g": FlagInfo{0, compileUnaryCallback},
        "-g0": FlagInfo{0, compileUnaryCallback},
        "-ggdb": FlagInfo{0, compileUnaryCallback},
        "-ggdb3": FlagInfo{0, compileUnaryCallback},
        "-gdwarf-2": FlagInfo{0, compileUnaryCallback},
        "-gdwarf-3": FlagInfo{0, compileUnaryCallback},
        "-gline-tables-only": FlagInfo{0, compileUnaryCallback},

        "-p": FlagInfo{0, compileUnaryCallback},
        "-pg": FlagInfo{0, compileUnaryCallback},

        "-O": FlagInfo{0, compileUnaryCallback},
        "-O0": FlagInfo{0, compileUnaryCallback},
        "-O1": FlagInfo{0, compileUnaryCallback},
        "-O2": FlagInfo{0, compileUnaryCallback},
        "-O3": FlagInfo{0, compileUnaryCallback},
        "-Os": FlagInfo{0, compileUnaryCallback},
        "-Ofast": FlagInfo{0, compileUnaryCallback},
        "-Og": FlagInfo{0, compileUnaryCallback},

        "-Xclang": FlagInfo{1, compileBinaryCallback},
        "-Xpreprocessor": FlagInfo{1, defaultBinaryCallback},
        "-Xassembler": FlagInfo{1, defaultBinaryCallback},
        "-Xlinker": FlagInfo{1, defaultBinaryCallback},

        "-l": FlagInfo{1, linkBinaryCallback},
        "-L": FlagInfo{1, linkBinaryCallback},
        "-T": FlagInfo{1, linkBinaryCallback},
        "-u": FlagInfo{1, linkBinaryCallback},

        "-e": FlagInfo{1, linkBinaryCallback},
        "-rpath": FlagInfo{1, linkBinaryCallback},

        "-shared": FlagInfo{0, linkUnaryCallback},
        "-static": FlagInfo{0, linkUnaryCallback},
        "-pie": FlagInfo{0, linkUnaryCallback},
        "-nostdlib": FlagInfo{0, linkUnaryCallback},
        "-nodefaultlibs": FlagInfo{0, linkUnaryCallback},
        "-rdynamic": FlagInfo{0, linkUnaryCallback},

        "-dynamiclib": FlagInfo{0, linkUnaryCallback},
        "-current_version": FlagInfo{1, linkBinaryCallback},
        "-compatibility_version": FlagInfo{1, linkBinaryCallback},

        "-print-multi-directory": FlagInfo{0, compileUnaryCallback},
        "-print-multi-lib": FlagInfo{0, compileUnaryCallback},
        "-print-libgcc-file-name": FlagInfo{0, compileUnaryCallback},

        "-fprofile-arcs": FlagInfo{0, compileLinkUnaryCallback},
        "-coverage": FlagInfo{0, compileLinkUnaryCallback},
        "--coverage": FlagInfo{0, compileLinkUnaryCallback},

        "-Wl,-dead_strip": FlagInfo{0, darwinWarningLinkUnaryCallback},
    }

    var argPatterns = map[string]FlagInfo{
        `^.+\.(c|cc|cpp|C|cxx|i|s|S|bc)$`: FlagInfo{0, inputFileCallback},
        `^.+\.([fF](|[0-9][0-9]|or|OR|pp|PP))$`: FlagInfo{0, inputFileCallback},
        `^.+\.(o|lo|So|so|po|a|dylib)$`: FlagInfo{0, objectFileCallback},
        `^.+\.dylib(\.\d)+$`: FlagInfo{0, objectFileCallback},
        `^.+\.(So|so)(\.\d)+$`: FlagInfo{0, objectFileCallback},
        `^-(l|L).+$`: FlagInfo{0, linkUnaryCallback},
        `^-I.+$`: FlagInfo{0, compileUnaryCallback},
        `^-D.+$`: FlagInfo{0, compileUnaryCallback},
        `^-U.+$`: FlagInfo{0, compileUnaryCallback},
        `^-Wl,.+$`: FlagInfo{0, linkUnaryCallback},
        `^-W.*$`: FlagInfo{0, compileUnaryCallback},
        `^-f.+$`: FlagInfo{0, compileUnaryCallback},
        `^-rtlib=.+$`: FlagInfo{0, linkUnaryCallback},
        `^-std=.+$`: FlagInfo{0, compileUnaryCallback},
        `^-stdlib=.+$`: FlagInfo{0, compileLinkUnaryCallback},
        `^-mtune=.+$`: FlagInfo{0, compileUnaryCallback},
        `^--sysroot=.+$`: FlagInfo{0, compileUnaryCallback},
        `^-print-prog-name=.*$`: FlagInfo{0, compileUnaryCallback},
        `^-print-file-name=.*$`: FlagInfo{0, compileUnaryCallback},
    }

    var pr = ParserResult{}
    pr.InputList = argList

    for len(argList) > 0 && !(pr.IsAssembly || pr.IsAssembleOnly || pr.IsPreprocessOnly) {
        var elem = argList[0]

        // Try to match the flag exactly
        if fi, ok := argsExactMatches[elem]; ok {
            fi.handler(pr, elem, argList[1:1+fi.arity])
            argList = argList[1+fi.arity:]
        // Else try to match a pattern
        } else {
            var listShift = 0
            for pattern, fi := range argPatterns {
                var regExp = regexp.MustCompile(pattern)
                if regExp.MatchString(elem) {
                    fi.handler(pr, elem, argList[1:1+fi.arity])
                    listShift = fi.arity
                    break
                }
            }
            argList = argList[1+listShift:]
        }

    }

    return pr
}

func inputFileCallback(pr ParserResult, flag string, _ []string) {
    var regExp = regexp.MustCompile(`\\.(s|S)$`)
    pr.InputFiles = append(pr.InputFiles, flag)
    if regExp.MatchString(flag) {
        pr.IsAssembly = true
    }
}

func outputFileCallback(pr ParserResult, _ string, args []string) {
    pr.OutputFilename = args[0]
}

func objectFileCallback(pr ParserResult, flag string, _ []string) {
    pr.ObjectFiles = append(pr.ObjectFiles, flag)
}

func preprocessOnlyCallback(pr ParserResult, _ string, _ []string) {
    pr.IsPreprocessOnly = true
}

func dependencyOnlyCallback(pr ParserResult, flag string, _ []string) {
    pr.IsDependencyOnly = true
    pr.CompileArgs = append(pr.CompileArgs, flag)
}

func assembleOnlyCallback(pr ParserResult, _ string, _ []string) {
    pr.IsAssembleOnly = true
}

func verboseFlagCallback(pr ParserResult, _ string, _ []string) {
    pr.IsVerbose = true
}

func compileOnlyCallback(pr ParserResult, _ string, _ []string) {
    pr.IsCompileOnly = true
}

func emitLLVMCallback(pr ParserResult, _ string, _ []string) {
    pr.IsCompileOnly = true
    pr.IsEmitLLVM = true
}

func linkUnaryCallback(pr ParserResult, flag string, _ []string) {
    pr.LinkArgs = append(pr.LinkArgs, flag)
}

func compileUnaryCallback(pr ParserResult, flag string, _ []string) {
    pr.CompileArgs = append(pr.CompileArgs, flag)
}

func darwinWarningLinkUnaryCallback(pr ParserResult, flag string, _ []string) {
    if runtime.GOOS == "darwin" {
        fmt.Println("The flag", flag, "cannot be used with this tool.")
    } else {
        pr.LinkArgs = append(pr.LinkArgs, flag)
    }
}

func defaultBinaryCallback(_ ParserResult, _ string, _ []string) {
    // Do nothing
}

func dependencyBinaryCallback(pr ParserResult, flag string, args []string) {
    pr.CompileArgs = append(pr.CompileArgs, flag, args[0])
    pr.IsDependencyOnly = true
}

func compileBinaryCallback(pr ParserResult, flag string, args []string) {
    pr.CompileArgs = append(pr.CompileArgs, flag, args[0])
}

func linkBinaryCallback(pr ParserResult, flag string, args []string) {
    pr.LinkArgs = append(pr.LinkArgs, flag, args[0])
}

func compileLinkUnaryCallback(pr ParserResult, flag string, _ []string) {
    pr.LinkArgs = append(pr.LinkArgs, flag)
    pr.CompileArgs = append(pr.CompileArgs, flag)
}
