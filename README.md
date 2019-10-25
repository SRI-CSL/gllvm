![GLLVM](data/dragon128x128.png?raw_true)
# Concurrent Whole Program LLVM in Go

[![License](https://img.shields.io/badge/License-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![Build Status](https://travis-ci.org/SRI-CSL/gllvm.svg?branch=master)](https://travis-ci.org/SRI-CSL/gllvm)
[![Go Report Card](https://goreportcard.com/badge/github.com/SRI-CSL/gllvm)](https://goreportcard.com/report/github.com/SRI-CSL/gllvm)

**TL; DR:**  A drop-in replacement for [wllvm](https://github.com/SRI-CSL/whole-program-llvm), that builds the
bitcode in parallel, and is faster. A comparison between the two tools can be gleaned from building the [Linux kernel.](https://github.com/SRI-CSL/gllvm/tree/master/examples/linux-kernel)

## Quick Start Comparison Table

| wllvm command/env variable  | gllvm command/env variable  |
|-----------------------------|-----------------------------|
|  wllvm                      | gclang                      |
|  wllvm++                    | gclang++                    |
|  extract-bc                 | get-bc                      |
|  wllvm-sanity-checker       | gsanity-check               |
|  LLVM_COMPILER_PATH         | LLVM_COMPILER_PATH          |
|  LLVM_CC_NAME      ...      | LLVM_CC_NAME          ...   |
|  WLLVM_CONFIGURE_ONLY       | WLLVM_CONFIGURE_ONLY        |
|  WLLVM_OUTPUT_LEVEL         | WLLVM_OUTPUT_LEVEL          |
|  WLLVM_OUTPUT_FILE          | WLLVM_OUTPUT_FILE           |
|  LLVM_COMPILER              | *not supported* (clang only)|
|  LLVM_GCC_PREFIX            | *not supported* (clang only)|
|  LLVM_DRAGONEGG_PLUGIN      | *not supported* (clang only)|


This project, `gllvm`, provides tools for building whole-program (or
whole-library) LLVM bitcode files from an unmodified C or C++
source package. It currently runs on `*nix` platforms such as Linux,
FreeBSD, and Mac OS X. It is a Go port of [wllvm](https://github.com/SRI-CSL/whole-program-llvm).

`gllvm` provides compiler wrappers that work in two
phases. The wrappers first invoke the compiler as normal. Then, for
each object file, they call a bitcode compiler to produce LLVM
bitcode. The wrappers then store the location of the generated bitcode
file in a dedicated section of the object file.  When object files are
linked together, the contents of the dedicated sections are
concatenated (so we don't lose the locations of any of the constituent
bitcode files). After the build completes, one can use a `gllvm`
utility to read the contents of the dedicated section and link all of
the bitcode into a single whole-program bitcode file. This utility
works for both executable and native libraries.

For more details see [wllvm](https://github.com/SRI-CSL/whole-program-llvm).

## Prerequisites

To install `gllvm` you need the go language [tool](https://golang.org/doc/install).

To use `gllvm` you need clang/clang++ and the llvm tools llvm-link and llvm-ar.
`gllvm` is agnostic to the actual llvm version. `gllvm` also relies on standard build
tools such as `objcopy` and `ld`.


## Installation

To install, simply do
```
go get github.com/SRI-CSL/gllvm/cmd/...
```
This should install four binaries: `gclang`, `gclang++`, `get-bc`, and `gsanity-check`
in the `$GOPATH/bin` directory.

## Usage

`gclang` and
`gclang++` are the wrappers used to compile C and C++.  `get-bc` is used for
extracting the bitcode from a build product (either an object file, executable, library
or archive). `gsanity-check` can be used for detecting configuration errors.

Here is a simple example. Assuming that clang is in your `PATH`, you can build
bitcode for `pkg-config` as follows:

```
tar xf pkg-config-0.26.tar.gz
cd pkg-config-0.26
CC=gclang ./configure
make
```

This should produce the executable `pkg-config`. To extract the bitcode:
```
get-bc pkg-config
```

which will produce the bitcode module `pkg-config.bc`.


If clang and the llvm tools are not in your `PATH`, you will need to set some
environment variables.


 * `LLVM_COMPILER_PATH` can be set to the absolute path of the directory that
   contains the compiler and the other LLVM tools to be used.

 * `LLVM_CC_NAME` can be set if your clang compiler is not called `clang` but
    something like `clang-3.7`. Similarly `LLVM_CXX_NAME` can be used to
    describe what the C++ compiler is called. We also pay attention to the
    environment variables `LLVM_LINK_NAME` and `LLVM_AR_NAME` in an
    analogous way.

Another useful environment variable is `WLLVM_CONFIGURE_ONLY`. Its use is explained in the
README of  [wllvm](https://github.com/SRI-CSL/whole-program-llvm).

`gllvm` does not support the dragonegg plugin. All other features of `wllvm`, such as logging, and the bitcode store,
are supported in exactly the same fashion as documented [here](https://github.com/SRI-CSL/whole-program-llvm).


## Under the hoods


Both `wllvm` and `gllvm` toolsets do much the same thing, but the way
they do it is slightly different. The `gllvm` toolset's code base is
written in `golang`, and is largely derived from the `wllvm`'s python
codebase.

Both generate object files and bitcode files using the
compiler. `wllvm` can use `gcc` and `dragonegg`, `gllvm` can only use
`clang`. The `gllvm` toolset does these two tasks in parallel, while
`wllvm` does them sequentially.  This together with the slowness of
python's `fork exec`-ing, and it's interpreted nature accounts for the
large efficiency gap between the two toolsets.

Both inject the path of the bitcode version of the `.o` file into a
dedicated segment of the `.o` file itself. This segment is the same
across toolsets, so extracting the bitcode can be done by the
appropriate tool in either toolset. On `*nix` both toolsets use
`objcopy` to add the segment, while on OS X they use `ld`.

When the object files are linked into the resulting library or
executable, the bitcode path segments are appended, so the resulting
binary contains the paths of all the bitcode files that constitute the
binary.  To extract the sections the `gllvm` toolset uses the golang
packages `"debug/elf"` and `"debug/macho"`, while the `wllvm` toolset
uses `objdump` on `*nix`, and `otool` on OS X.

Both tools then use `llvm-link` or `llvm-ar` to combine the bitcode
files into the desired form.

## Customization under the hood.

You can specify the exact version of `objcopy` and `ld` that `gllvm` uses
to manipulate the artifacts by setting the `GLLVM_OBJCOPY` and `GLLVM_LD`
environment variables. For more details of what's under the `gllvm` hood, try
```
gsanity-check -e
```

## License

`gllvm` is released under a BSD license. See the file `LICENSE` for [details.](LICENSE)

---

This material is based upon work supported by the National Science
Foundation under Grant
[ACI-1440800](http://www.nsf.gov/awardsearch/showAward?AWD_ID=1440800). Any
opinions, findings, and conclusions or recommendations expressed in
this material are those of the author(s) and do not necessarily
reflect the views of the National Science Foundation.
