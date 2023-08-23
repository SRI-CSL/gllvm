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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// ParserResult is the result of parsing and partioning the command line arguments.
type ParserResult struct {
	InputList        []string
	InputFiles       []string
	ObjectFiles      []string
	OutputFilename   string
	CompileArgs      []string
	LinkArgs         []string
	ForbiddenFlags   []string
	IsVerbose        bool
	IsDependencyOnly bool
	IsPreprocessOnly bool
	IsAssembleOnly   bool
	IsAssembly       bool
	IsCompileOnly    bool
	IsEmitLLVM       bool
	IsLTO            bool
	IsPrintOnly      bool
}

const parserResultFormat = `
InputList:         %v
InputFiles:        %v
ObjectFiles:       %v
OutputFilename:    %v
CompileArgs:       %v
LinkArgs:          %v
ForbiddenFlags:    %v
IsVerbose:         %v
IsDependencyOnly:  %v
IsPreprocessOnly:  %v
IsAssembleOnly:    %v
IsAssembly:        %v
IsCompileOnly:     %v
IsEmitLLVM:        %v
IsLTO:             %v
IsPrintOnly:       %v
`

func (pr *ParserResult) String() string {
	return fmt.Sprintf(parserResultFormat,
		pr.InputList,
		pr.InputFiles,
		pr.ObjectFiles,
		pr.OutputFilename,
		pr.CompileArgs,
		pr.LinkArgs,
		pr.ForbiddenFlags,
		pr.IsVerbose,
		pr.IsDependencyOnly,
		pr.IsPreprocessOnly,
		pr.IsAssembleOnly,
		pr.IsAssembly,
		pr.IsCompileOnly,
		pr.IsEmitLLVM,
		pr.IsLTO,
		pr.IsPrintOnly)
}

type flagInfo struct {
	arity   int
	handler func(string, []string)
}

type argPattern struct {
	pattern string
	finfo   flagInfo
}

// SkipBitcodeGeneration indicates whether or not we should generate bitcode for these command line options.
func (pr *ParserResult) SkipBitcodeGeneration() bool {
	reason := "No particular reason"
	retval := false
	squark := LogDebug
	if LLVMConfigureOnly != "" {
		reason = "we are in configure only mode"
		retval = true
	} else if len(pr.InputFiles) == 0 {
		reason = "we did not see any input files"
		retval = true
	} else if pr.IsEmitLLVM {
		squark = LogWarning
		reason = "we are in emit-llvm mode, and so the compiler is doing the job for us"
		retval = true
	} else if pr.IsLTO {
		squark = LogWarning
		reason = "we are doing link time optimization, and so the compiler is doing the job for us"
		retval = true
	} else if pr.IsAssembly {
		reason = "the input file(s) are written in assembly"
		squark = LogWarning
		retval = true
	} else if pr.IsAssembleOnly {
		reason = "we are assembling only, and so have nowhere to embed the path of the bitcode"
		retval = true
	} else if pr.IsDependencyOnly && !pr.IsCompileOnly {
		reason = "we are only computing dependencies at this stage"
		retval = true
	} else if pr.IsPreprocessOnly {
		reason = "we are in preprocess only mode"
		retval = true
	} else if pr.IsPrintOnly {
		reason = "we are in print only mode, and so have nowhere to embed the path of the bitcode"
		retval = true
	}
	if retval {
		squark(" We are skipping bitcode generation because %v.\n", reason)
	}
	return retval
}

// Parse analyzes the command line aruguments and returns the result of that analysis.
func Parse(argList []string) ParserResult {
	var pr = ParserResult{}
	pr.InputList = argList

	var argsExactMatches = map[string]flagInfo{

		"/dev/null": {0, pr.inputFileCallback}, //iam: linux kernel

		"-":  {0, pr.printOnlyCallback},
		"-o": {1, pr.outputFileCallback},
		"-c": {0, pr.compileOnlyCallback},
		"-E": {0, pr.preprocessOnlyCallback},
		"-S": {0, pr.assembleOnlyCallback},

		"--verbose": {0, pr.verboseFlagCallback},
		"--param":   {1, pr.defaultBinaryCallback},
		"-aux-info": {1, pr.defaultBinaryCallback},

		"-target": {1, pr.compileLinkBinaryCallback},

		"--version": {0, pr.compileOnlyCallback},
		"-v":        {0, pr.compileOnlyCallback},

		"-w": {0, pr.compileUnaryCallback},
		"-W": {0, pr.compileUnaryCallback},

		"-emit-llvm": {0, pr.emitLLVMCallback},
		"-flto":      {0, pr.linkTimeOptimizationCallback},

		"-pipe":                  {0, pr.compileUnaryCallback},
		"-undef":                 {0, pr.compileUnaryCallback},
		"-nostdinc":              {0, pr.compileUnaryCallback},
		"-nostdinc++":            {0, pr.compileUnaryCallback},
		"-Qunused-arguments":     {0, pr.compileUnaryCallback},
		"-no-integrated-as":      {0, pr.compileUnaryCallback},
		"-integrated-as":         {0, pr.compileUnaryCallback},
		"-no-canonical-prefixes": {0, pr.compileLinkUnaryCallback},

		"--sysroot": {1, pr.compileLinkBinaryCallback}, //iam: musl stuff

		//<archaic flags>
		"-no-cpp-precomp": {0, pr.compileUnaryCallback},
		//</archaic flags>

		"-pthread":     {0, pr.linkUnaryCallback},
		"-nostdlibinc": {0, pr.compileUnaryCallback},

		"-mno-omit-leaf-frame-pointer": {0, pr.compileUnaryCallback},
		"-maes":                        {0, pr.compileUnaryCallback},
		"-mno-aes":                     {0, pr.compileUnaryCallback},
		"-mavx":                        {0, pr.compileUnaryCallback},
		"-mno-avx":                     {0, pr.compileUnaryCallback},
		"-mavx2":                       {0, pr.compileUnaryCallback},
		"-mno-avx2":                    {0, pr.compileUnaryCallback},
		"-mno-red-zone":                {0, pr.compileUnaryCallback},
		"-mmmx":                        {0, pr.compileUnaryCallback},
		"-mbmi":                        {0, pr.compileUnaryCallback},
		"-mbmi2":                       {0, pr.compileUnaryCallback},
		"-mf161c":                      {0, pr.compileUnaryCallback},
		"-mfma":                        {0, pr.compileUnaryCallback},
		"-mno-mmx":                     {0, pr.compileUnaryCallback},
		"-mno-global-merge":            {0, pr.compileUnaryCallback}, //iam: linux kernel stuff
		"-mno-80387":                   {0, pr.compileUnaryCallback}, //iam: linux kernel stuff
		"-msse":                        {0, pr.compileUnaryCallback},
		"-mno-sse":                     {0, pr.compileUnaryCallback},
		"-msse2":                       {0, pr.compileUnaryCallback},
		"-mno-sse2":                    {0, pr.compileUnaryCallback},
		"-msse3":                       {0, pr.compileUnaryCallback},
		"-mno-sse3":                    {0, pr.compileUnaryCallback},
		"-mssse3":                      {0, pr.compileUnaryCallback},
		"-mno-ssse3":                   {0, pr.compileUnaryCallback},
		"-msse4":                       {0, pr.compileUnaryCallback},
		"-mno-sse4":                    {0, pr.compileUnaryCallback},
		"-msse4.1":                     {0, pr.compileUnaryCallback},
		"-mno-sse4.1":                  {0, pr.compileUnaryCallback},
		"-msse4.2":                     {0, pr.compileUnaryCallback},
		"-mno-sse4.2":                  {0, pr.compileUnaryCallback},
		"-msoft-float":                 {0, pr.compileUnaryCallback},
		"-m3dnow":                      {0, pr.compileUnaryCallback},
		"-mno-3dnow":                   {0, pr.compileUnaryCallback},
		"-m16":                         {0, pr.compileLinkUnaryCallback}, //iam: linux kernel stuff
		"-m32":                         {0, pr.compileLinkUnaryCallback},
		"-m64":                         {0, pr.compileLinkUnaryCallback},
		"-mstackrealign":               {0, pr.compileUnaryCallback},
		"-mretpoline-external-thunk":   {0, pr.compileUnaryCallback}, //iam: linux kernel stuff
		"-mno-fp-ret-in-387":           {0, pr.compileUnaryCallback}, //iam: linux kernel stuff
		"-mskip-rax-setup":             {0, pr.compileUnaryCallback}, //iam: linux kernel stuff
		"-mindirect-branch-register":   {0, pr.compileUnaryCallback}, //iam: linux kernel stuff

		"-mllvm": {1, pr.compileBinaryCallback}, //iam: chromium

		"-A": {1, pr.compileBinaryCallback},
		"-D": {1, pr.compileBinaryCallback},
		"-U": {1, pr.compileBinaryCallback},

		"-arch": {1, pr.compileBinaryCallback}, //iam: openssl

		"-P": {1, pr.compileUnaryCallback}, //iam: linux kernel stuff (linker script stuff)
		"-C": {1, pr.compileUnaryCallback}, //iam: linux kernel stuff (linker script stuff)

		"-M":   {0, pr.dependencyOnlyCallback},
		"-MM":  {0, pr.dependencyOnlyCallback},
		"-MF":  {1, pr.dependencyBinaryCallback},
		"-MJ":  {1, pr.dependencyBinaryCallback},
		"-MG":  {0, pr.dependencyOnlyCallback},
		"-MP":  {0, pr.dependencyOnlyCallback},
		"-MT":  {1, pr.dependencyBinaryCallback},
		"-MQ":  {1, pr.dependencyBinaryCallback},
		"-MD":  {0, pr.dependencyOnlyCallback},
		"-MV":  {0, pr.dependencyOnlyCallback},
		"-MMD": {0, pr.dependencyOnlyCallback},

		"-I":                 {1, pr.compileBinaryCallback},
		"-idirafter":         {1, pr.compileBinaryCallback},
		"-include":           {1, pr.compileBinaryCallback},
		"-imacros":           {1, pr.compileBinaryCallback},
		"-iprefix":           {1, pr.compileBinaryCallback},
		"-iwithprefix":       {1, pr.compileBinaryCallback},
		"-iwithprefixbefore": {1, pr.compileBinaryCallback},
		"-isystem":           {1, pr.compileBinaryCallback},
		"-isysroot":          {1, pr.compileBinaryCallback},
		"-iquote":            {1, pr.compileBinaryCallback},
		"-imultilib":         {1, pr.compileBinaryCallback},

		"-ansi":     {0, pr.compileUnaryCallback},
		"-pedantic": {0, pr.compileUnaryCallback},
		"-x":        {1, pr.compileBinaryCallback},

		"-g":                    {0, pr.compileUnaryCallback},
		"-g0":                   {0, pr.compileUnaryCallback},
		"-g1":                   {0, pr.compileUnaryCallback},
		"-g2":                   {0, pr.compileUnaryCallback},
		"-g3":                   {0, pr.compileUnaryCallback},
		"-ggdb":                 {0, pr.compileUnaryCallback},
		"-ggdb0":                {0, pr.compileUnaryCallback},
		"-ggdb1":                {0, pr.compileUnaryCallback},
		"-ggdb2":                {0, pr.compileUnaryCallback},
		"-ggdb3":                {0, pr.compileUnaryCallback},
		"-gdwarf":               {0, pr.compileUnaryCallback},
		"-gdwarf-2":             {0, pr.compileUnaryCallback},
		"-gdwarf-3":             {0, pr.compileUnaryCallback},
		"-gdwarf-4":             {0, pr.compileUnaryCallback},
		"-gline-tables-only":    {0, pr.compileUnaryCallback},
		"-grecord-gcc-switches": {0, pr.compileUnaryCallback},
		"-ggnu-pubnames":        {0, pr.compileUnaryCallback},

		"-p":  {0, pr.compileUnaryCallback},
		"-pg": {0, pr.compileUnaryCallback},

		"-O":     {0, pr.compileUnaryCallback},
		"-O0":    {0, pr.compileUnaryCallback},
		"-O1":    {0, pr.compileUnaryCallback},
		"-O2":    {0, pr.compileUnaryCallback},
		"-O3":    {0, pr.compileUnaryCallback},
		"-Os":    {0, pr.compileUnaryCallback},
		"-Ofast": {0, pr.compileUnaryCallback},
		"-Og":    {0, pr.compileUnaryCallback},
		"-Oz":    {0, pr.compileUnaryCallback}, //iam: linux kernel

		"-Xclang":        {1, pr.compileBinaryCallback},
		"-Xpreprocessor": {1, pr.defaultBinaryCallback},
		"-Xassembler":    {1, pr.defaultBinaryCallback},
		"-Xlinker":       {1, pr.defaultBinaryCallback},

		"-l":            {1, pr.linkBinaryCallback},
		"-L":            {1, pr.linkBinaryCallback},
		"-T":            {1, pr.linkBinaryCallback},
		"-u":            {1, pr.linkBinaryCallback},
		"-install_name": {1, pr.linkBinaryCallback},

		"-e":     {1, pr.linkBinaryCallback},
		"-rpath": {1, pr.linkBinaryCallback},

		"-shared":        {0, pr.linkUnaryCallback},
		"-static":        {0, pr.linkUnaryCallback},
		"-static-libgcc": {0, pr.linkUnaryCallback}, //iam: musl stuff
		"-pie":           {0, pr.linkUnaryCallback},
		"-nostdlib":      {0, pr.linkUnaryCallback},
		"-nodefaultlibs": {0, pr.linkUnaryCallback},
		"-rdynamic":      {0, pr.linkUnaryCallback},

		"-dynamiclib":            {0, pr.linkUnaryCallback},
		"-current_version":       {1, pr.linkBinaryCallback},
		"-compatibility_version": {1, pr.linkBinaryCallback},

		"-print-multi-directory":  {0, pr.compileUnaryCallback},
		"-print-multi-lib":        {0, pr.compileUnaryCallback},
		"-print-libgcc-file-name": {0, pr.compileUnaryCallback},
		"-print-search-dirs":      {0, pr.compileUnaryCallback},

		"-fprofile-arcs": {0, pr.compileLinkUnaryCallback},
		"-coverage":      {0, pr.compileLinkUnaryCallback},
		"--coverage":     {0, pr.compileLinkUnaryCallback},
		"-fopenmp":       {0, pr.compileLinkUnaryCallback},

		"-Wl,-dead_strip": {0, pr.warningLinkUnaryCallback},
		"-dead_strip":     {0, pr.warningLinkUnaryCallback}, //iam: tor does this. We lose the bitcode :-(
	}

	// iam: this is a list because matching needs to be done in order.
	// if you add a NEW pattern; make sure it is before any existing pattern that also
	// matches and has a conflicting flagInfo value. Also be careful with flags that can contain filenames, like
	// linker info flags or dependency flags
	var argPatterns = [...]argPattern{
		// dependency file generation (comes first because it can conatin file names)
		{`^-MF.*$`, flagInfo{0, pr.compileUnaryCallback}}, //Write depfile output from -MMD, -MD, -MM, or -M to <file>
		{`^-MJ.*$`, flagInfo{0, pr.compileUnaryCallback}}, //Write a compilation database entry per input
		{`^-MQ.*$`, flagInfo{0, pr.compileUnaryCallback}}, //Specify name of main file output to quote in depfile
		{`^-MT.*$`, flagInfo{0, pr.compileUnaryCallback}}, //Specify name of main file output in depfile
		//iam, need to be careful here, not mix up linker and warning flags.
		{`^-Wl,.+$`, flagInfo{0, pr.linkUnaryCallback}},
		{`^-W[^l].*$`, flagInfo{0, pr.compileUnaryCallback}},
		{`^-W[l][^,].*$`, flagInfo{0, pr.compileUnaryCallback}}, //iam: tor has a few -Wl...
		{`^-(l|L).+$`, flagInfo{0, pr.linkUnaryCallback}},
		{`^-I.+$`, flagInfo{0, pr.compileUnaryCallback}},
		{`^-D.+$`, flagInfo{0, pr.compileUnaryCallback}},
		{`^-B.+$`, flagInfo{0, pr.compileLinkUnaryCallback}},
		{`^-isystem.+$`, flagInfo{0, pr.compileLinkUnaryCallback}},
		{`^-U.+$`, flagInfo{0, pr.compileUnaryCallback}},
		{`^-fsanitize=.+$`, flagInfo{0, pr.compileLinkUnaryCallback}},
		{`^-fuse-ld=.+$`, flagInfo{0, pr.linkUnaryCallback}},         //iam:  musl stuff
		{`^-flto=.+$`, flagInfo{0, pr.linkTimeOptimizationCallback}}, //iam: new lto stuff
		{`^-f.+$`, flagInfo{0, pr.compileUnaryCallback}},
		{`^-rtlib=.+$`, flagInfo{0, pr.linkUnaryCallback}},
		{`^-std=.+$`, flagInfo{0, pr.compileUnaryCallback}},
		{`^-stdlib=.+$`, flagInfo{0, pr.compileLinkUnaryCallback}},
		{`^-mtune=.+$`, flagInfo{0, pr.compileUnaryCallback}},
		{`^--sysroot=.+$`, flagInfo{0, pr.compileLinkUnaryCallback}}, //both compile and link time
		{`^-print-.*$`, flagInfo{0, pr.compileUnaryCallback}},        // generic catch all for the print commands
		{`^-mmacosx-version-min=.+$`, flagInfo{0, pr.compileLinkUnaryCallback}},
		{`^-mstack-alignment=.+$`, flagInfo{0, pr.compileUnaryCallback}},          //iam, linux kernel stuff
		{`^-march=.+$`, flagInfo{0, pr.compileUnaryCallback}},                     //iam: linux kernel stuff
		{`^-mregparm=.+$`, flagInfo{0, pr.compileUnaryCallback}},                  //iam: linux kernel stuff
		{`^-mcmodel=.+$`, flagInfo{0, pr.compileUnaryCallback}},                   //iam: linux kernel stuff
		{`^-mpreferred-stack-boundary=.+$`, flagInfo{0, pr.compileUnaryCallback}}, //iam: linux kernel stuff
		{`^-mindirect-branch=.+$`, flagInfo{0, pr.compileUnaryCallback}},          //iam: linux kernel stuff
		{`^--param=.+$`, flagInfo{0, pr.compileUnaryCallback}},                    //iam: linux kernel stuff
		//the above come first because they are anchored at the start, and some can contain filenames.
		{`^.+\.(c|cc|cpp|C|cxx|i|s|S|bc)$`, flagInfo{0, pr.inputFileCallback}},
		{`^.+\.([fF](|[0-9][0-9]|or|OR|pp|PP))$`, flagInfo{0, pr.inputFileCallback}},
		//iam: it's a bit fragile as to what we recognize as an object file.
		// this also shows up in the compile function attachBitcodePathToObject, so additions
		// here, should also be additions there.
		{`^.+\.(o|lo|So|so|po|a|dylib|pico|nossppico)$`, flagInfo{0, pr.objectFileCallback}}, //iam: pico and nossppico are FreeBSD
		{`^.+\.dylib(\.\d)+$`, flagInfo{0, pr.objectFileCallback}},
		{`^.+\.(So|so)(\.\d)+$`, flagInfo{0, pr.objectFileCallback}},
	}

	for len(argList) > 0 {
		var elem = argList[0]

		// Try to match the flag exactly
		if fi, ok := argsExactMatches[elem]; ok {
			fi.handler(elem, argList[1:1+fi.arity])
			argList = argList[1+fi.arity:]
			// else it is more complicated, either a pattern or a group
		} else {
			var listShift = 0
			//need to handle the N-ary grouping flag
			if elem == "-Wl,--start-group" {
				endgroup := indexOf("-Wl,--end-group", argList)
				if endgroup > 0 {
					pr.linkerGroupCallback(elem, endgroup+1, argList)
					listShift = endgroup
				} else {
					LogWarning("Failed to find '-Wl,--end-group' matching '-Wl,--start-group'\n")
					pr.compileUnaryCallback(elem, argList[1:1])
				}
				//else try to match a pattern
			} else {
				var matched = false
				for _, argPat := range argPatterns {
					pattern := argPat.pattern
					fi := argPat.finfo
					var regExp = regexp.MustCompile(pattern)
					if regExp.MatchString(elem) {
						fi.handler(elem, argList[1:1+fi.arity])
						listShift = fi.arity
						matched = true
						break
					}
				}
				if !matched {
					ok, _ := IsObjectFileForOS(elem, runtime.GOOS)
					if ok {
						pr.objectFileCallback(elem, argList[1:1])
					} else {
						LogWarning("Did not recognize the compiler flag: %v\n", elem)
						pr.compileUnaryCallback(elem, argList[1:1])
					}
				}
			}
			argList = argList[1+listShift:]
		}
	}
	return pr
}

func indexOf(value string, slice []string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
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
		var _, baseNameWithExt = path.Split(srcFile)
		// issue #30:  main.cpp and main.c cause conflicts.
		var baseName = strings.TrimSuffix(baseNameWithExt, filepath.Ext(baseNameWithExt))
		bcBase = fmt.Sprintf(".%s.o.bc", baseNameWithExt)
		if hidden {
			objBase = fmt.Sprintf(".%s.o", baseNameWithExt)
		} else {
			objBase = fmt.Sprintf("%s.o", baseName)
		}
	}
	return
}

// Return a hash for the absolute object path
func getHashedPath(path string) string {
	inputBytes := []byte(path)
	hasher := sha256.New()
	//Hash interface claims this never returns an error
	hasher.Write(inputBytes)
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash
}

func (pr *ParserResult) inputFileCallback(flag string, _ []string) {
	var regExp = regexp.MustCompile(`\.(s|S)$`)
	pr.InputFiles = append(pr.InputFiles, flag)
	if regExp.MatchString(flag) {
		pr.IsAssembly = true
	}
}

func (pr *ParserResult) outputFileCallback(_ string, args []string) {
	pr.OutputFilename = args[0]
}

func (pr *ParserResult) objectFileCallback(flag string, _ []string) {
	// FIXME: the object file is appended to ObjectFiles that
	// is used nowhere else in the code
	pr.ObjectFiles = append(pr.ObjectFiles, flag)
	// We append the object files to link args to handle the
	// -Wl,--start-group obj_1.o ... obj_n.o -Wl,--end-group case
	pr.LinkArgs = append(pr.LinkArgs, flag)
}

func (pr *ParserResult) linkerGroupCallback(start string, count int, args []string) {
	group := args[0:count]
	pr.LinkArgs = append(pr.LinkArgs, group...)
}

func (pr *ParserResult) preprocessOnlyCallback(_ string, _ []string) {
	pr.IsPreprocessOnly = true
}

func (pr *ParserResult) dependencyOnlyCallback(flag string, _ []string) {
	pr.IsDependencyOnly = true
	pr.CompileArgs = append(pr.CompileArgs, flag)
}

func (pr *ParserResult) printOnlyCallback(flag string, _ []string) {
	pr.IsPrintOnly = true
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

func (pr *ParserResult) linkTimeOptimizationCallback(_ string, _ []string) {
	pr.IsLTO = true
}

func (pr *ParserResult) linkUnaryCallback(flag string, _ []string) {
	pr.LinkArgs = append(pr.LinkArgs, flag)
}

func (pr *ParserResult) compileUnaryCallback(flag string, _ []string) {
	pr.CompileArgs = append(pr.CompileArgs, flag)
}

func (pr *ParserResult) warningLinkUnaryCallback(flag string, _ []string) {
	LogWarning("The flag %v cannot be used with this tool, we ignore it, else we lose the bitcode section.\n", flag)
	pr.ForbiddenFlags = append(pr.ForbiddenFlags, flag)
}

func (pr *ParserResult) defaultBinaryCallback(_ string, _ []string) {
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

func (pr *ParserResult) compileLinkBinaryCallback(flag string, args []string) {
	pr.LinkArgs = append(pr.LinkArgs, flag, args[0])
	pr.CompileArgs = append(pr.CompileArgs, flag, args[0])
}
