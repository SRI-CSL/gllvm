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
	"bytes"
	"os"
	"os/exec"
)

// Executes a command then returns true for success, false if there was an error, err is either nil or the error.
func execCmd(cmdExecName string, args []string, workingDir string) (success bool, err error) {
	cmd := exec.Command(cmdExecName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = workingDir
	err = cmd.Run()
	ecode := 0
	if err != nil {
		ecode = 1
	}
	LogDebug("execCmd: %v %v had exitCode %v\n", cmdExecName, args, ecode)
	if err != nil {
		LogDebug("execCmd: error was %v\n", err)
	}
	success = (ecode == 0)
	return
}

// Executes a command then returns the output as a string, err is either nil or the error.
func runCmd(cmdExecName string, args []string) (output string, err error) {
	var outb bytes.Buffer
	var errb bytes.Buffer
	cmd := exec.Command(cmdExecName, args...)
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	LogDebug("runCmd: %v %v\n", cmdExecName, args)
	if err != nil {
		LogDebug("runCmd: error was %v\n", err)
	}
	output = outb.String()
	return
}

// Deduplicate a potentially unsorted list of strings in-place without changing their order
func dedupeStrings(strings *[]string) {
	seen := make(map[string]bool)
	count := 0
	for _, s := range *strings {
		if _, exists := seen[s]; !exists {
			seen[s] = true
			(*strings)[count] = s
			count++
		}
	}
	*strings = (*strings)[:count]
}
