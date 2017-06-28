# Go Whole Program LLVM


## Overview


This project, gllvm, provides tools for building whole-program (or
whole-library) LLVM bitcode files from an unmodified C or C++
source package. It currently runs on `*nix` platforms such as Linux,
FreeBSD, and Mac OS X. It is a Go port of the [wllvm](https://github.com/SRI-CSL/whole-program-llvm).

gllvm provides compiler wrappers that work in two
steps. The wrappers first invoke the compiler as normal. Then, for
each object file, they call a bitcode compiler to produce LLVM
bitcode. The wrappers also store the location of the generated bitcode
file in a dedicated section of the object file.  When object files are
linked together, the contents of the dedicated sections are
concatenated (so we don't lose the locations of any of the constituent
bitcode files). After the build completes, one can use a gllvm
utility to read the contents of the dedicated section and link all of
the bitcode into a single whole-program bitcode file. This utility
works for both executable and native libraries.

This two-phase build process is necessary to be a drop-in replacement
for gcc or g++ in any build system.  Using the LTO framework in gcc
and the gold linker plugin works in many cases, but fails in the
presence of static libraries in builds.  gllvm's approach has the
distinct advantage of generating working binaries, in case some part
of a build process requires that.

gllvm currently works with clang.

## Installation


#### Requirements

You need the Go compiler to compile gllvm, and both the clang/clang++
executables and the llvm tools -- llvm-link, llvm-ar -- to use gllvm. Follow
the instructions here to get started: https://golang.org/doc/install.

As for now, let us name `$GOROOT` your root Go path that you can obtain by
typing `go env GOPATH` in a terminal session -- it is usually `$HOME/go`
by default. It is worth noticing that a standard Go installation will install
the binaries generated for the project under `$GOROOT/bin`. Make sure that you
added the `$GOROOT/bin` directory to your `$PATH` variable.

#### Build

First, you must checkout the project under the directory `$GOROOT/src`:
```
cd $GOROOT/src
git clone https://github.com/SRI-CSL/gllvm
```

To build and install gllvm on your system, type:
```
make install
```

## Usage

gllvm includes three symlinks to the program's binary: `gclang` and
`gclang++`to compile C and C++, and an auxiliary tool `get-bc` for
extracting the bitcode from a build product (object file, executable, library
or archive).

Some useful environment variables are listed here:

 * `GLLVM_CC_NAME` can be set if your clang compiler is not called `clang` but
    something like `clang-3.7`. Similarly `GLLVM_CXX_NAME` can be used to
    describe what the C++ compiler is called. We also pay attention to the
    environment  variables `GLLVM_LINK_NAME` and `GLLVM_AR_NAME` in an
    analagous way, since they too get adorned with suffixes in various Linux
    distributions.

 * `GLLVM_TOOLS_PATH` can be set to the absolute path to the folder that
   contains the compiler and other LLVM tools such as `llvm-link` to be used.
   This prevents searching for the compiler in your PATH environment variable.
   This can be useful if you have different versions of clang on your system
   and you want to easily switch compilers without tinkering with your PATH
   variable.
   Example `GLLVM_TOOLS_PATH=/home/user/llvm_and_clang/Debug+Asserts/bin`.

* `GLLVM_CONFIGURE_ONLY` can be set to anything. If it is set, `gclang`
   and `gclang++` behave like a normal C or C++ compiler. They do not
   produce bitcode. Setting `GLLVM_CONFIGURE_ONLY` may prevent configuration
   errors caused by the unexpected production of hidden bitcode files. It is
   sometimes required when configuring a build.

## Examples

### Building a bitcode module with clang


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


### Building bitcode archive

```
tar -xvf bullet-2.81-rev2613.tgz
mkdir bullet-bin
cd bullet-bin
CC=gclang CXX=gclang++ cmake ../bullet-2.81-rev2613/
make

# Produces src/LinearMath/libLinearMath.bca
get-bc src/LinearMath/libLinearMath.a
```

Note that by default extracting bitcode from an archive produces an archive of
bitcode. You can also extract the bitcode directly into a module:
```
get-bc -b src/LinearMath/libLinearMath.a
```
produces `src/LinearMath/libLinearMath.a.bc`.


### Configuring without building bitcode

Sometimes it is necessary to disable the production of bitcode. Typically this
is during configuration, where the production of unexpected files can confuse
the configure script. For this we have a flag `GLLVM_CONFIGURE_ONLY` which
can be used as follows:
```
GLLVM_CONFIGURE_ONLY=1 CC=gclang ./configure
CC=gclang make
```


### Building a bitcode archive then extracting the bitcode

```
tar xvfz jansson-2.7.tar.gz
cd jansson-2.7
CC=gclang ./configure
make
mkdir bitcode
cp src/.libs/libjansson.a bitcode
cd bitcode
get-bc libjansson.a
llvm-ar x libjansson.bca
ls -la
```

## Miscellaneous Features

### Preserving bitcode files in a store

Sometimes it can be useful to preserve the bitcode files produced in a
build, either to prevent deletion or to retrieve them later. If the
environment variable `GLLVM_BC_STORE` is set to the absolute path of
an existing directory, then gllvm will copy the produced bitcode files
into that directory. The name of a copied bitcode file is the hash of the path
to the original bitcode file. For convenience, when using both the manifest
feature of `get-bc` and the store, the manifest will contain both the
original path, and the store path.


### Debugging


The GLLVM tools can show various levels of output to aid with debugging.
To show this output set the `GLLVM_OUTPUT_LEVEL` environment
variable to one of the following levels:

 * `ERROR`
 * `WARNING`
 * `INFO`
 * `DEBUG`

For example:
```
    export GLLVM_OUTPUT_LEVEL=DEBUG
```
Output will be directed to the standard error stream, unless you specify the
path of a logfile via the `GLLVM_OUTPUT_FILE` environment variable.

For example:
```
    export GLLVM_OUTPUT_FILE=/tmp/gllvm.log
```
