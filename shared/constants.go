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

//  IAN! Remember to tag the repo, and publish a release on GitHub.
//
//  version history:
//
//  1.0.0
//  1.0.1                  various bug fixes
//
//  1.2.0 April 28 2018    linux kernel work, sorting bitcode files, etc.
//        May 2 2018       handleArchives rewritten to handle multiple occurrences of files with the same name.
//                         corresponds with wllvm 1.2.0. Gonna try and keep them in synch.
//
//  1.2.1 May 13th 2018    -fsanitize= needs to be compile AND link.
//
//  1.2.2 June 1st 2018    Fix extracting from archives on darwin, plus travis build for both linux and darwin,
//                         a few ittle fixes from building tor and it's dependencies.
//
//  1.2.4 June  21 2019    Random fixes (basically the same as wllvm 1.2.7)
//
//  1.2.5 October 23 2019  Fixes for issues #30 and #31. Plus a new branch where I can futz hygenically.
//
//  1.2.6 March 24 2020    Added the support for the LLVM_BITCODE_GENERATION_FLAGS environment variable.
//                         See https://github.com/travitch/whole-program-llvm/issues/96 for details.
//
//  1.2.7 August 4 2020    William Woodruff's (@woodruffw) tweaks to the get-bc command (a strict mode).
//
//  1.2.8 October 13 2020  Wensheng Tang's (@legendtang) requests and fixes to the -flto issues (issue #39).
//                         Added support for the LTO_LINKING_FLAGS environment variable.
//
//  1.2.9 October 31 2020  Gleb Popov's FreeBSD bug fixes (@arrowd) (issue #41), plus some more improvements
//                         to the command line parsing (a never ending saga), including the golang packages
//                         debug/elf and debug/macho to spot object files...
//
//  1.3.0 February 18 2021 William Woodruff's (@woodruffw) fix to the (-w -W fiasco).

const gllvmVersion = "1.3.0"
const gllvmReleaseDate = "February 18 2021"

const osDARWIN = "darwin"
const osLINUX = "linux"
const osFREEBSD = "freebsd"
