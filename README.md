<p align="center">
<img align="center" src="data/dragon128x128.png?raw_true">
</p>

# Whole Program LLVM in Go

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
|  wfortran                   | gflang                      |
|  extract-bc                 | get-bc                      |
|  wllvm-sanity-checker       | gsanity-check               |
|  LLVM_COMPILER_PATH         | LLVM_COMPILER_PATH          |
|  LLVM_CC_NAME      ...      | LLVM_CC_NAME          ...   |
|                             | LLVM_F_NAME                 |
|  WLLVM_CONFIGURE_ONLY       | WLLVM_CONFIGURE_ONLY        |
|  WLLVM_OUTPUT_LEVEL         | WLLVM_OUTPUT_LEVEL          |
|  WLLVM_OUTPUT_FILE          | WLLVM_OUTPUT_FILE           |
|  LLVM_COMPILER              | *not supported* (clang only)|
|  LLVM_GCC_PREFIX            | *not supported* (clang only)|
|  LLVM_DRAGONEGG_PLUGIN      | *not supported* (clang only)|
|  LLVM_LINK_FLAGS            | LLVM_LINK_FLAGS             |


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

To use `gllvm` you need clang/clang++/flang and the llvm tools llvm-link and llvm-ar.
`gllvm` is agnostic to the actual llvm version. `gllvm` also relies on standard build
tools such as `objcopy` and `ld`.


## Installation

To install, simply do (making sure to include those `...`)
```
go install github.com/SRI-CSL/gllvm/cmd/...@latest
```
This should install six binaries: `gclang`, `gclang++`, `gflang`, `get-bc`, `gparse`, and `gsanity-check`
in the `$GOPATH/bin` directory. 

## Usage

`gclang` and
`gclang++` are the wrappers used to compile C and C++.  
`gflang` is the wrapper used to compile Fortran.
`get-bc` is used for
extracting the bitcode from a build product (either an object file, executable, library
or archive). `gsanity-check` can be used for detecting configuration errors. `gparse` can be used to examine how `gllvm` parses compiler/linker lines.

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

which will produce the bitcode module `pkg-config.bc`. For more on this example
see [here](https://github.com/SRI-CSL/gllvm/tree/master/examples/pkg-config).

## Advanced Configuration

If clang and the llvm tools are not in your `PATH`, you will need to set some
environment variables.


 * `LLVM_COMPILER_PATH` can be set to the absolute path of the directory that
   contains the compiler and the other LLVM tools to be used.

 * `LLVM_CC_NAME` can be set if your clang compiler is not called `clang` but
    something like `clang-3.7`. Similarly `LLVM_CXX_NAME` and `LLVM_F_NAME` can be used to
    describe what the C++ and Fortran compilers are called, respectively. We also pay attention to the
    environment variables `LLVM_LINK_NAME` and `LLVM_AR_NAME` in an
    analogous way.

Another useful, and sometimes necessary, environment variable is `WLLVM_CONFIGURE_ONLY`.

* `WLLVM_CONFIGURE_ONLY` can be set to anything. If it is set, `gclang`
   and `gclang++` behave like a normal C or C++ compiler. They do not
   produce bitcode.  Setting `WLLVM_CONFIGURE_ONLY` may prevent
   configuration errors caused by the unexpected production of hidden
   bitcode files. It is sometimes required when configuring a build.
   For example:
   ```
   WLLVM_CONFIGURE_ONLY=1 CC=gclang ./configure
   make
   ```

## Extracting the Bitcode

The `get-bc` tool is used to extract the bitcode from a build artifact, such as an executable, object file, thin archive, archive, or library. In the simplest use case, as seen above,
one simply does:

```
get-bc -o <name of bitcode file> <path to executable>
```
This will produce the desired bitcode file. The situation is similar for an object file.
For an archive or library, there is a choice as to whether you produce a bitcode module
or a bitcode archive. This choice is made by using the `-b` switch.

Another useful switch is the `-m` switch which will, in addition to producing the
bitcode, will also produce a manifest of the bitcode files
that made up the final product. As is typical

```
get-bc -h
```
will list all the commandline switches. Since we use the `golang` `flag` module,
the switches must precede the artifact path.



## Preserving bitcode files in a store

Sometimes, because of pathological build systems, it can be useful
to preserve the bitcode files produced in a
build, either to prevent deletion or to retrieve it later. If the
environment variable `WLLVM_BC_STORE` is set to the absolute path of
an existing directory,
then WLLVM will copy the produced bitcode file into that directory.
The name of the copied bitcode file is the hash of the path to the
original bitcode file.  For convenience, when using both the manifest
feature of `get-bc` and the store, the manifest will contain both
the original path, and the store path.

## Debugging


The gllvm tools can show various levels of output to aid with debugging.
To show this output set the `WLLVM_OUTPUT_LEVEL` environment
variable to one of the following levels:

 * `ERROR`
 * `WARNING`
 * `AUDIT`
 * `INFO`
 * `DEBUG`

For example:
```
    export WLLVM_OUTPUT_LEVEL=DEBUG
```
Output will be directed to the standard error stream, unless you specify the
path of a logfile via the `WLLVM_OUTPUT_FILE` environment variable.
The `AUDIT` level, new in 2022, logs only the calls to the compiler, and indicates 
whether each call is *compiling* or *linking*, the compiler used, and the arguments provided.

For example:
```
    export WLLVM_OUTPUT_FILE=/tmp/gllvm.log
```

## Dragons Begone

`gllvm` does not support the dragonegg plugin.


## Sanity Checking

Too many environment variables? Try doing a sanity check:

```
gsanity-check
```
it might point out what is wrong.



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

## Customizing the BitCode Generation (e.g. LTO)

In some situations it is desirable to pass certain flags to `clang` in the step that
produces the bitcode. This can be fulfilled by setting the
`LLVM_BITCODE_GENERATION_FLAGS` environment variable to the desired
flags, for example `"-flto -fwhole-program-vtables"`.

In other situations it is desirable to pass certain flags to `llvm-link` in the step
that merges multiple individual bitcode files together (i.e., within `get-bc`).
This can be fulfilled by setting the `LLVM_LINK_FLAGS` environment variable to
the desired flags, for example `"-internalize -only-needed"`.

## Beware of link time optimization.

If the package you are building happens to take advantage of recent `clang` developments 
such as *link time optimization* (indicated by the presence of compiler flag `-flto`), then
your build is unlikely to produce anything that `get-bc` will work on. This is to be
expected. When working under these flags, the compiler actually produces object files that are bitcode,
your only recourse here is to try and save these object files, and retrieve them yourself.
This can be done by setting the `LTO_LINKING_FLAGS` to be something like
`"-g -Wl,-plugin-opt=save-temps"` which will be appended to the flags at link time.
This will at least preserve the bitcode files, even if `get-bc` will not be able to retrieve them for you.

## Cross-compilation notes

When cross-compiling a project (i.e. you pass the `--target=` or `-target` flag to the compiler), 
you'll need to set the `GLLVM_OBJCOPY` variable to either 
* `llvm-objcopy` to use LLVM's objcopy, which naturally supports all targets that clang does.
* `YOUR-TARGET-TRIPLE-objcopy` to use GNU's objcopy, since `objcopy` only supports the native architecture.

Example:
```sh
# test program
echo 'int main() { return 0; }' > a.c 
clang --target=aarch64-linux-gnu a.c # works
gclang --target=aarch64-linux-gnu a.c # breaks
GLLVM_OBJCOPY=llvm-objcopy gclang --target=aarch64-linux-gnu a.c # works
GLLVM_OBJCOPY=aarch64-linux-gnu-objcopy gclang --target=aarch64-linux-gnu a.c # works if you have GNU's arm64 toolchain
```

## Developer tools

Debugging usually boils down to looking in the logs, maybe adding a print statement or two.
There is an additional executable, not mentioned above, called `gparse` that gets installed 
along with `gclang`, `gclang++`, `gflang`, `get-bc` and `gsanity-check`. `gparse` takes the command line
arguments to the compiler, and outputs how it parsed them. This can sometimes be helpful.

## License

`gllvm` is released under a BSD license. See the file `LICENSE` for [details.](LICENSE)

---

This material is based upon work supported by the National Science
Foundation under Grant
[ACI-1440800](http://www.nsf.gov/awardsearch/showAward?AWD_ID=1440800). Any
opinions, findings, and conclusions or recommendations expressed in
this material are those of the author(s) and do not necessarily
reflect the views of the National Science Foundation.
