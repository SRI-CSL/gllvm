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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
)

type bitcodeToObjectLink struct {
	bcPath  string
	objPath string
}

//Compile wraps a call to the compiler with the given args.
func Compile(args []string, compiler string) (exitCode int) {

	exitCode = 0
	// in the configureOnly case we have to know the exit code of the compile
	// because that is how configure figures out what it can and cannot do.

	compileOk := true
	attachOk := true

	compilerExecName := GetCompilerExecName(compiler)

	pr := Parse(args)

	var wg sync.WaitGroup

	LogDebug("Compile using parsed arguments:%v\n", &pr)

	// If configure only, emit-llvm, flto, or print only are set, just execute the compiler
	if pr.SkipBitcodeGeneration() {
		wg.Add(1)
		go execCompile(compilerExecName, pr, &wg, &compileOk)
		wg.Wait()

		if !compileOk {
			exitCode = 1
		}

		// Else try to build bitcode as well
	} else {
		var bcObjLinks []bitcodeToObjectLink
		var newObjectFiles []string

		wg.Add(2)
		go execCompile(compilerExecName, pr, &wg, &compileOk)
		go buildAndAttachBitcode(compilerExecName, pr, &bcObjLinks, &newObjectFiles, &wg, &attachOk)
		wg.Wait()

		//grok the exit code
		if !compileOk || !attachOk {
			exitCode = 1
		} else {
			// When objects and bitcode are built we can attach bitcode paths
			// to object files and link
			for _, link := range bcObjLinks {
				attachBitcodePathToObject(link.bcPath, link.objPath)
			}
			if !pr.IsCompileOnly {
				compileTimeLinkFiles(compilerExecName, pr, newObjectFiles)
			}
		}
	}
	return
}

// Compiles bitcode files and mutates the list of bc->obj links to perform + the list of
// new object files to link
func buildAndAttachBitcode(compilerExecName string, pr ParserResult, bcObjLinks *[]bitcodeToObjectLink, newObjectFiles *[]string, wg *sync.WaitGroup, ok *bool) {
	defer (*wg).Done()

	// var hidden = !pr.IsCompileOnly
	var success bool
	var objFile string
	var bcFile string

	for _, srcFile := range pr.InputFiles {
		objFile, success = buildObjectFileContentAddressed(compilerExecName, pr, srcFile)
		if !success {
			*ok = false
			break
		}

		bcFile, success = buildBitcodeFileContentAddressed(compilerExecName, pr, srcFile)
		if !success {
			*ok = false
			break
		}

		*newObjectFiles = append(*newObjectFiles, objFile)
		*bcObjLinks = append(*bcObjLinks, bitcodeToObjectLink{bcPath: bcFile, objPath: objFile})
	}
}

func attachBitcodePathToObject(bcFile, objFile string) (success bool) {
	success = false
	// We can only attach a bitcode path to certain file types
	// this is too fragile, we need to look into a better way to do this.
	// We probably should be using debug/macho and debug/elf according to the OS we are atop of.
	extension := filepath.Ext(objFile)
	switch extension {
	case
		".o",
		".lo",
		".os",
		".So",
		".pico",      //iam: pico is FreeBSD, ".pico" denotes a position-independent relocatable object.
		".nossppico", //iam: also FreeBSD, ".nossppico" denotes a position-independent relocatable object without stack smashing protection.
		".po":        //iam: profiled object
		LogDebug("attachBitcodePathToObject recognized %v as something it can inject into.\n", extension)
		success = injectPath(extension, bcFile, objFile)
		return
	default:
		//OK we have to work harder here
		ok, err := injectableViaFileType(objFile)
		LogDebug("attachBitcodePathToObject: injectableViaFileType returned  ok=%v  err=%v", ok, err)
		if ok {
			success = injectPath(extension, bcFile, objFile)
			return
		}
		if err != nil {
			// OK we have to work EVEN harder here (the file utility is not installed - probably)
			// N.B. this will probably fail if we are cross compiling.
			ok, err = injectableViaDebug(objFile)
			LogDebug("attachBitcodePathToObject: injectableViaDebug returned  ok=%v  err=%v", ok, err)
			if ok {
				success = injectPath(extension, bcFile, objFile)
				return
			}
			if err != nil {
				LogWarning("attachBitcodePathToObject: injectableViaDebug failed %v", err)
			}
		}
		LogWarning("attachBitcodePathToObject ignoring unrecognized extension %v of file %v of unknown type\n", extension, objFile)
	}
	return
}

// move this out to concentrate on the object path analysis above.
func injectPath(extension, bcFile, objFile string) (success bool) {
	success = false
	// Store bitcode path to temp file
	var absBcPath, _ = filepath.Abs(bcFile)
	tmpContent := []byte(absBcPath + "\n")
	tmpFile, err := ioutil.TempFile("", "gllvm")
	if err != nil {
		LogError("attachBitcodePathToObject: %v\n", err)
		return
	}
	defer CheckDefer(func() error { return os.Remove(tmpFile.Name()) })
	if _, err := tmpFile.Write(tmpContent); err != nil {
		LogError("attachBitcodePathToObject: %v\n", err)
		return
	}
	if err := tmpFile.Close(); err != nil {
		LogError("attachBitcodePathToObject: %v\n", err)
		return
	}

	// Let's write the bitcode section
	var attachCmd string
	var attachCmdArgs []string
	if runtime.GOOS == osDARWIN {
		if len(LLVMLd) > 0 {
			attachCmd = LLVMLd
		} else {
			attachCmd = "ld"
		}
		attachCmdArgs = []string{"-r", "-keep_private_externs", objFile, "-sectcreate", DarwinSegmentName, DarwinSectionName, tmpFile.Name(), "-o", objFile}
	} else {
		if len(LLVMObjcopy) > 0 {
			attachCmd = LLVMObjcopy
		} else {
			attachCmd = "objcopy"
		}
		attachCmdArgs = []string{"--add-section", ELFSectionName + "=" + tmpFile.Name(), objFile}
	}

	// Run the attach command and ignore errors
	_, nerr := execCmd(attachCmd, attachCmdArgs, "")
	if nerr != nil {
		LogWarning("attachBitcodePathToObject: %v %v failed because %v\n", attachCmd, attachCmdArgs, nerr)
		return
	}

	// Copy bitcode file to store, if necessary
	if bcStorePath := LLVMBitcodeStorePath; bcStorePath != "" {
		destFilePath := path.Join(bcStorePath, getHashedPath(absBcPath))
		in, _ := os.Open(absBcPath)
		defer CheckDefer(func() error { return in.Close() })
		out, _ := os.Create(destFilePath)
		defer CheckDefer(func() error { return out.Close() })
		_, err := io.Copy(out, in)
		if err != nil {
			LogWarning("Copying bc to bitcode archive %v failed because %v\n", destFilePath, err)
			return
		}
		err = out.Sync()
		if err != nil {
			LogWarning("Syncing bitcode archive %v failed because %v\n", destFilePath, err)
			return
		}

	}
	success = true
	return
}

func compileTimeLinkFiles(compilerExecName string, pr ParserResult, objFiles []string) {
	var outputFile = pr.OutputFilename
	if outputFile == "" {
		outputFile = "a.out"
	}
	args := []string{}
	//iam: unclear if this is necessary here
	if pr.IsLTO {
		args = append(args, LLVMLtoLDFLAGS...)
	}
	args = append(args, objFiles...)
	args = append(args, pr.LinkArgs...)
	args = append(args, "-o", outputFile)
	success, err := execCmd(compilerExecName, args, "")
	if !success {
		LogError("%v %v failed to link: %v.", compilerExecName, args, err)
	} else {
		LogInfo("LINKING: %v %v", compilerExecName, args)
	}
}

// Tries to build the specified source file to an object, returning a content-addressed path
func buildObjectFileContentAddressed(compilerExecName string, pr ParserResult, srcFile string) (objFile string, success bool) {
	objFile = ""
	success = false

	// NOTE: We don't remove the temporary file because we rename it.
	tempObjFile, err := ioutil.TempFile("", ".*.o")
	if err != nil {
		return
	}

	args := pr.CompileArgs[:]
	args = append(args, "-c", srcFile, "-o", tempObjFile.Name())
	LogDebug("buildObjectFileContentAddressed: %v", args)

	success, err = execCmd(compilerExecName, args, "")
	if !success {
		LogError("Failed to build object file for %s because: %v\n", srcFile, err)
		return
	}

	bcContents := []byte{}
	bcContents, err = ioutil.ReadFile(tempObjFile.Name())
	if err != nil {
		LogError("Failed to read temporary object for %s: %v\n", srcFile, err)
		return
	}

	objFile = fmt.Sprintf(".%s.o", sha256Hash(bcContents))
	err = os.Rename(tempObjFile.Name(), objFile)
	if err != nil {
		LogError("Failed to rename object for %s (%s -> %s): %v\n", srcFile, tempObjFile.Name(), objFile, err)
		return
	}

	success = true
	return
}

// Tries to build the specified source file to bitcode, returning a content-addressed path
func buildBitcodeFileContentAddressed(compilerExecName string, pr ParserResult, srcFile string) (bcFile string, success bool) {
	bcFile = ""
	success = false

	// NOTE: We don't remove the temporary file because we rename it.
	tempBcFile, err := ioutil.TempFile("", ".*.bc")
	if err != nil {
		return
	}

	args := pr.CompileArgs[:]
	args = append(args, LLVMbcGen...)
	args = append(args, "-emit-llvm", "-c", srcFile, "-o", tempBcFile.Name())

	success, err = execCmd(compilerExecName, args, "")
	if !success {
		LogError("Failed to build bitcode file for %s because: %v\n", srcFile, err)
		return
	}

	bcContents := []byte{}
	bcContents, err = ioutil.ReadFile(tempBcFile.Name())
	if err != nil {
		LogError("Failed to read temporary bitcode for %s: %v\n", srcFile, err)
		return
	}

	bcFile = fmt.Sprintf(".%s.bc", sha256Hash(bcContents))
	err = os.Rename(tempBcFile.Name(), bcFile)
	if err != nil {
		LogError("Failed to rename bitcode for %s (%s -> %s): %v\n", srcFile, tempBcFile.Name(), bcFile, err)
		return
	}

	success = true
	return
}

// Tries to build object file or link...
func execCompile(compilerExecName string, pr ParserResult, wg *sync.WaitGroup, ok *bool) {
	defer (*wg).Done()
	//iam: strickly speaking we should do more work here depending on whether this is
	//     a compile only, a link only, or ...
	//     But for the now, we just remove forbidden arguments
	var success bool
	var err error
	// start afresh
	arguments := []string{}
	// we are linking rather than compiling
	if len(pr.InputFiles) == 0 && len(pr.LinkArgs) > 0 {
		if pr.IsLTO {
			arguments = append(arguments, LLVMLtoLDFLAGS...)
		}
	}
	//iam: this is clunky. is there a better way?
	if len(pr.ForbiddenFlags) > 0 {
		for _, arg := range pr.InputList {
			found := false
			for _, bad := range pr.ForbiddenFlags {
				if bad == arg {
					found = true
					break
				}
			}
			if !found {
				arguments = append(arguments, arg)
			}
		}
	} else {
		arguments = append(arguments, pr.InputList...)
	}
	LogDebug("Calling execCmd(%v, %v)", compilerExecName, arguments)
	success, err = execCmd(compilerExecName, arguments, "")
	if !success {
		LogError("Failed to compile using given arguments:\n%v %v\nexit status: %v\n", compilerExecName, arguments, err)
		*ok = false
	}
}

// Tries to build object file

func oldExecCompile(compilerExecName string, pr ParserResult, wg *sync.WaitGroup, ok *bool) {
	defer (*wg).Done()
	//iam: strickly speaking we should do more work here depending on whether this is
	//     a compile only, a link only, or ...
	//     But for the now, we just remove forbidden arguments
	var success bool
	var err error

	if len(pr.ForbiddenFlags) > 0 {
		filteredArgs := pr.InputList[:0]
		for _, arg := range pr.InputList {
			found := false
			for _, bad := range pr.ForbiddenFlags {
				if bad == arg {
					found = true
					break
				}
			}
			if !found {
				filteredArgs = append(filteredArgs, arg)
			}
		}
		success, err = execCmd(compilerExecName, filteredArgs, "")
	} else {
		success, err = execCmd(compilerExecName, pr.InputList, "")
	}

	if !success {
		LogError("Failed to compile using given arguments:\n%v %v\nexit status: %v\n", compilerExecName, pr.InputList, err)
		*ok = false
	}
}

// GetCompilerExecName returns the full path of the executable
func GetCompilerExecName(compiler string) string {
	switch compiler {
	case "clang":
		if LLVMCCName != "" {
			return filepath.Join(LLVMToolChainBinDir, LLVMCCName)
		}
		return filepath.Join(LLVMToolChainBinDir, compiler)
	case "clang++":
		if LLVMCXXName != "" {
			return filepath.Join(LLVMToolChainBinDir, LLVMCXXName)
		}
		return filepath.Join(LLVMToolChainBinDir, compiler)
	case "flang":
		if LLVMFName != "" {
			return filepath.Join(LLVMToolChainBinDir, LLVMCCName)
		}
		return filepath.Join(LLVMToolChainBinDir, compiler)
	default:
		LogError("The compiler %s is not supported by this tool.", compiler)
		return ""
	}
}

//CheckDefer is used to check the return values of defers
func CheckDefer(f func() error) {
	if err := f(); err != nil {
		LogWarning("CheckDefer received error: %v\n", err)
	}
}
