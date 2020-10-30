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
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	// File types
	fileTypeUNDEFINED = iota
	fileTypeELFEXECUTABLE
	fileTypeELFOBJECT
	fileTypeELFSHARED
	fileTypeMACHEXECUTABLE
	fileTypeMACHOBJECT
	fileTypeMACHSHARED
	fileTypeARCHIVE
	fileTypeTHINARCHIVE

	fileTypeERROR
)

//iam:
// this is not that robust, because it depends on the file utility "file" which is
// often  missing on docker images (the klee doker file had this problem)
func getFileType(realPath string) (fileType int, err error) {

	// We need the file command to guess the file type
	fileType = fileTypeERROR
	err = nil
	cmd := exec.Command("file", realPath)
	out, err := cmd.Output()
	if err != nil {
		LogError("There was an error getting the type of %s. Make sure that the 'file' command is installed.", realPath)
		return
	}

	fo := string(out)

	if strings.Contains(fo, "ELF") {

		if strings.Contains(fo, "executable") {
			fileType = fileTypeELFEXECUTABLE
		} else if strings.Contains(fo, "shared") {
			fileType = fileTypeELFSHARED
		} else if strings.Contains(fo, "relocatable") {
			fileType = fileTypeELFOBJECT
		} else {
			fileType = fileTypeUNDEFINED
		}

	} else if strings.Contains(fo, "Mach-O") {

		if strings.Contains(fo, "executable") {
			fileType = fileTypeMACHEXECUTABLE
		} else if strings.Contains(fo, "dynamically linked shared") {
			fileType = fileTypeMACHSHARED
		} else if strings.Contains(fo, "object") {
			fileType = fileTypeMACHOBJECT
		} else {
			fileType = fileTypeUNDEFINED
		}

	} else if strings.Contains(fo, "current ar archive") {
		fileType = fileTypeARCHIVE
	} else if strings.Contains(fo, "thin archive") {
		fileType = fileTypeTHINARCHIVE
	} else {
		fileType = fileTypeUNDEFINED
	}
	return
}

// isPlainFile returns true if the file is stat-able (i.e. exists etc), and is not a directory, else it returns false.
func isPlainFile(objectFile string) (ok bool) {
	ok = false
	info, err := os.Stat(objectFile)
	if os.IsNotExist(err) || info.IsDir() {
		return
	}
	if err != nil {
		return
	}
	ok = true
	return
}

func injectableViaFileType(objectFile string) (ok bool, err error) {
	ok = false
	err = nil
	plain := isPlainFile(objectFile)
	if !plain {
		return
	}
	fileType, err := getFileType(objectFile)
	if err != nil {
		return
	}
	ok = (fileType == fileTypeELFOBJECT) || (fileType == fileTypeELFOBJECT)
	return
}

func injectableViaDebug(objectFile string) (ok bool, err error) {
	ok = false
	err = nil
	// I guess we are not doing cross compiling. Otherwise we are fucking up here.
	ok, err = IsObjectFileForOS(objectFile, runtime.GOOS)
	return
}

//IsObjectFileForOS returns true if the given file is an object file for the given OS, using the debug/elf and debug/macho packages.
func IsObjectFileForOS(objectFile string, operatingSys string) (ok bool, err error) {
	plain := isPlainFile(objectFile)
	if !plain {
		return
	}
	switch operatingSys {
	case "linux", "freebsd":
		var lbinFile *elf.File
		lbinFile, err = elf.Open(objectFile)
		if err != nil {
			return
		}
		dfileType := lbinFile.FileHeader.Type
		ok = (dfileType == elf.ET_REL)
		return
	case "darwin":
		var dbinFile *macho.File
		dbinFile, err = macho.Open(objectFile)
		if err != nil {
			return
		}
		dfileType := dbinFile.FileHeader.Type
		ok = (dfileType == macho.TypeObj)

		return
	}
	return
}
