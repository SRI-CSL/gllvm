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

// BinaryType is the 'intersection' of elf.Type and macho.Type and partitions
// the binary world into categories we are most interested in. Missing is
// ARCHIVE but that is because it is not an elf format, so we cannot entirely
// eliminate the use of the 'file' utility (cf getFileType below).
type BinaryType uint32

const (
	//BinaryUnknown signals that the file does not fit into our three simple minded categories
	BinaryUnknown BinaryType = 0
	//BinaryObject is the type of an object file, the output unit of compilation
	BinaryObject BinaryType = 1
	//BinaryExecutable is the type of an executable file
	BinaryExecutable BinaryType = 2
	//BinaryShared is the type of a shared or dynamic library
	BinaryShared BinaryType = 3
)

func (bt BinaryType) String() string {
	switch bt {
	case BinaryUnknown:
		return "Unknown"
	case BinaryObject:
		return "Object"
	case BinaryExecutable:
		return "Executable"
	case BinaryShared:
		return "Library"
	default:
		return "Error"
	}
}

// GetBinaryType gets the binary type of the given path
func GetBinaryType(path string) (bt BinaryType) {
	bt = BinaryUnknown
	plain := IsPlainFile(path)
	if !plain {
		return
	}
	// try the format that suits the platform first
	operatingSys := runtime.GOOS
	switch operatingSys {
	case "linux", "freebsd":
		bt, _ = ElfFileType(path)
	case "darwin":
		bt, _ = MachoFileType(path)
	}
	if bt != BinaryUnknown {
		return
	}
	// try the other format instead
	switch operatingSys {
	case "linux", "freebsd":
		bt, _ = MachoFileType(path)
	case "darwin":
		bt, _ = ElfFileType(path)
	}
	return
}

func elfType2BinaryType(et elf.Type) (bt BinaryType) {
	bt = BinaryUnknown
	switch et {
	case elf.ET_NONE:
		bt = BinaryUnknown
	case elf.ET_REL:
		bt = BinaryObject
	case elf.ET_EXEC:
		bt = BinaryExecutable
	case elf.ET_DYN:
		bt = BinaryShared
	case elf.ET_CORE, elf.ET_LOOS, elf.ET_HIOS, elf.ET_LOPROC, elf.ET_HIPROC:
		bt = BinaryUnknown
	default:
		bt = BinaryUnknown
	}
	return
}

func machoType2BinaryType(mt macho.Type) (bt BinaryType) {
	bt = BinaryUnknown
	switch mt {
	case macho.TypeObj:
		bt = BinaryObject
	case macho.TypeExec:
		bt = BinaryExecutable
	case macho.TypeDylib:
		bt = BinaryShared
	case macho.TypeBundle:
		bt = BinaryUnknown
	default:
		bt = BinaryUnknown
	}
	return
}

// IsPlainFile returns true if the file is stat-able (i.e. exists etc), and is not a directory, else it returns false.
func IsPlainFile(objectFile string) (ok bool) {
	info, err := os.Stat(objectFile)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		return
	}
	if info.IsDir() {
		return
	}
	ok = true
	return
}

func injectableViaFileType(objectFile string) (ok bool, err error) {
	plain := IsPlainFile(objectFile)
	if !plain {
		return
	}
	fileType, err := getFileType(objectFile)
	if err != nil {
		return
	}
	ok = (fileType == fileTypeELFOBJECT) || (fileType == fileTypeMACHOBJECT)
	return
}

func injectableViaDebug(objectFile string) (ok bool, err error) {
	// I guess we are not doing cross compiling. Otherwise we are fucking up here.
	ok, err = IsObjectFileForOS(objectFile, runtime.GOOS)
	return
}

// ElfFileType returns the elf.Type of the given file name
func ElfFileType(objectFile string) (code BinaryType, err error) {
	var lbinFile *elf.File
	lbinFile, err = elf.Open(objectFile)
	if err != nil {
		return
	}
	code = elfType2BinaryType(lbinFile.FileHeader.Type)
	return
}

// MachoFileType returns the macho.Type of the given file name
func MachoFileType(objectFile string) (code BinaryType, err error) {
	var dbinFile *macho.File
	dbinFile, err = macho.Open(objectFile)
	if err != nil {
		return
	}
	code = machoType2BinaryType(dbinFile.FileHeader.Type)
	return
}

// IsObjectFileForOS returns true if the given file is an object file for the given OS, using the debug/elf and debug/macho packages.
func IsObjectFileForOS(objectFile string, operatingSys string) (ok bool, err error) {
	plain := IsPlainFile(objectFile)
	if !plain {
		return
	}
	var binaryType BinaryType
	switch operatingSys {
	case "linux", "freebsd":
		binaryType, err = ElfFileType(objectFile)
	case "darwin":
		binaryType, err = MachoFileType(objectFile)
	}
	if err != nil {
		return
	}
	ok = (binaryType == BinaryObject)
	return
}

// file types via the unix 'file' utility
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

// iam: this is not that robust, because it depends on the file utility "file" which is
// often missing on docker images (the klee docker file had this problem)
// this is only used in extraction, not in compilation.
func getFileType(realPath string) (fileType int, err error) {
	// We need the file command to guess the file type
	fileType = fileTypeERROR
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
