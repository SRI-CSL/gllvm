#!/bin/bash -x
# Make sure we exit if there is a failure
set -e

cd .travis

echo `pwd`

export PATH=/usr/lib/llvm-3.8/bin:${PATH}
export WLLVM_OUTPUT_LEVEL=WARNING

make

exit $?
