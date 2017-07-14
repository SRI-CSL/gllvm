#!/bin/bash -x
# Make sure we exit if there is a failure
set -e


export PATH=/usr/lib/llvm-3.8/bin:${PATH}
export WLLVM_OUTPUT_LEVEL=WARNING


git clone https://github.com/SRI-CSL/musllvm.git musllvm

cd musllvm

WLLVM_CONFIGURE_ONLY=1  CC=gclang ./configure --target=LLVM --build=LLVM

make

exit $?


get-bc -b ./lib/libc.a

if [ -s "./lib/libc.a.bc" ]
then
    echo "libc.a.bc exists."
else
    exit 1
fi

exit 0
