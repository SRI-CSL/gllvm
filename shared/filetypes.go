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
	"os/exec"
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
)

func getFileType(realPath string) (fileType int) {
	// We need the file command to guess the file type
	cmd := exec.Command("file", realPath)
	out, err := cmd.Output()
	if err != nil {
		LogFatal("There was an error getting the type of %s. Make sure that the 'file' command is installed.", realPath)
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
	}  else {
		fileType = fileTypeUNDEFINED
	}



	// Test the output
	//	if fo := string(out); strings.Contains(fo, "ELF") && strings.Contains(fo, "executable") {
	//		fileType = fileTypeELFEXECUTABLE
	//	} else if strings.Contains(fo, "Mach-O") && strings.Contains(fo, "executable") {
	//		fileType = fileTypeMACHEXECUTABLE
	//	} else if strings.Contains(fo, "ELF") && strings.Contains(fo, "shared") {
	//		fileType = fileTypeELFSHARED
	//	} else if strings.Contains(fo, "Mach-O") && strings.Contains(fo, "dynamically linked shared") {
	//		fileType = fileTypeMACHSHARED
	//	} else if strings.Contains(fo, "current ar archive") {
	//		fileType = fileTypeARCHIVE
	//	} else if strings.Contains(fo, "thin archive") {
	//		fileType = fileTypeTHINARCHIVE
	//	} else if strings.Contains(fo, "ELF") && strings.Contains(fo, "relocatable") {
	//		fileType = fileTypeELFOBJECT
	//	} else if strings.Contains(fo, "Mach-O") && strings.Contains(fo, "object") {
	//		fileType = fileTypeMACHOBJECT
	//	} else {
	//		fileType = fileTypeUNDEFINED
	//	}

	return
}
