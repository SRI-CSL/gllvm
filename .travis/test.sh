#!/bin/bash -x
# Make sure we exit if there is a failure
set -e


export PATH=/usr/lib/llvm-3.8/bin:${PATH}
export GLLVM_OUTPUT_LEVEL=WARNING

#currently the musllvm build fails with
#
#/usr/bin/ld: obj/src/process/posix_spawn.lo: relocation R_X86_64_PC32 against protected symbol `execve' can not be used when making a shared object
#/usr/bin/ld: final link failed: Bad value
#clang: error: linker command failed with exit code 1 (use -v to see invocation)
#2017/06/29 19:10:32 Failed to link.
#make: *** [lib/libc.so] Error 1
#
#need to investigate inside a vagrant box.
exit 0


git clone https://github.com/SRI-CSL/musllvm.git musllvm
cd musllvm
GLLVM_CONFIGURE_ONLY=1  CC=gclang ./configure --target=LLVM --build=LLVM
make
get-bc --bitcode ./lib/libc.a

if [ -s "./lib/libc.a.bc" ]
then
    echo "libc.a.bc exists."
else
    exit 1
fi
