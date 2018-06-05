#!/bin/bash -x
# Make sure we exit if there is a failure
set -e


export PATH=/usr/lib/llvm-3.5/bin:${PATH}
export WLLVM_OUTPUT=WARNING

gsanity-check

#setup the store so we test that feature as well
export WLLVM_BC_STORE=/tmp/bc
mkdir -p /tmp/bc

git clone https://github.com/SRI-CSL/musllvm.git musllvm
cd musllvm
WLLVM_CONFIGURE_ONLY=1  CC=gclang ./configure --target=LLVM --build=LLVM
make
get-bc -b ./lib/libc.a

if [ -s "./lib/libc.a.bc" ]
then
    echo "libc.a.bc exists (built from build artifacts)."
else
    exit 1
fi

#now lets makes sure the store has the bitcode too.
mv ./lib/libc.a .
make clean
get-bc -b ./libc.a

if [ -s "./libc.a.bc" ]
then
    echo "libc.a.bc exists (built from store)."
else
    exit 1
fi
