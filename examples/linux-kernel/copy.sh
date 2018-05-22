#!/usr/bin/env bash

### Copy all necessary files to standalone-build and link the kernel from bitcode

export build_home=$HOME/standalone-build
export ker=$HOME/linux-stable

bash handle-bi.sh

cp $ker/arch/x86/kernel/vmlinux.lds ./arch/x86/kernel/vmlinux.lds
cp $ker/.tmp_kallsyms2.o .
#automated fs script building
cd $ker
python parse-bi.py fs/built-in.o fs/out.sh $build_home/instrfs -1

cd fs
bash out.sh

cd $build_home/built-ins

for bclib in ./*.bc; do
    clang -c -no-integrated-as -mcmodel=kernel $bclib -o ${bclib/%.o.bc/bc.o}
done

cd $build_home/built-ins/fs
for bclib in ./*.bc; do
    clang -c -no-integrated-as -mcmodel=kernel $bclib -o ${bclib/%.o.bc/bc.o}
done

cd objects
for bcobj in ./*.bc; do
    clang -c -no-integrated-as -mcmodel=kernel $bcobj -o ${bcobj/%.o.bc/bc.o}
done

cd $build_home/lib
clang -c -no-integrated-as -mcmodel=kernel lib.a.bc


cd $build_home/arch/x86/lib
clang -c -no-integrated-as -mcmodel=kernel lib.a.bc

cd $build_home
#linking command (full bc)
ld --build-id -T ./arch/x86/kernel/vmlinux.lds --whole-archive built-ins/objects/ker_objects/head_64.o built-ins/objects/ker_objects/head64.o built-ins/objects/ker_objects/ebda.o built-ins/objects/ker_objects/platform-quirks.o built-ins/inibibc.o built-ins/objects/ker_objects/initramfs_data.o built-ins/arcbibc.o built-ins/objects/arch_assembly_objects/* built-ins/kerbibc.o built-ins/mmbibc.o \@instrfs built-ins/ipcbibc.o built-ins/secbibc.o built-ins/cptbibc.o built-ins/blkbibc.o built-ins/libbibc.o built-ins/objects/lib_assembly_objects/* built-ins/xlibbibc.o built-ins/objects/xlib_assembly_objects/* built-ins/dribibc.o built-ins/sndbibc.o built-ins/pcibibc.o built-ins/powbibc.o built-ins/objects/pow_assembly_objects/* built-ins/vidbibc.o built-ins/netbibc.o --no-whole-archive --start-group lib/lib.a.o arch/x86/lib/lib.a.o built-ins/objects/libx_objects/*  .tmp_kallsyms2.o --end-group -o vmlinux

# #linking command (partial bc)
# ld --build-id -T ./arch/x86/kernel/vmlinux.lds --whole-archive \
# built-ins/objects/ker_objects/head_64.o built-ins/objects/ker_objects/head64.o built-ins/objects/ker_objects/ebda.o built-ins/objects/ker_objects/platform-quirks.o \
# built-ins/inibibc.o built-ins/objects/ker_objects/initramfs_data.o built-ins/arcbi.o built-ins/kerbibc.o built-ins/mmbibc.o built-ins/fsbi.o \
# built-ins/ipcbibc.o built-ins/secbibc.o built-ins/cptbibc.o built-ins/blkbibc.o built-ins/libbi.o built-ins/xlibbi.o built-ins/dribibc.o \
# built-ins/sndbibc.o built-ins/pcibibc.o built-ins/powbibc.o built-ins/objects/pow_assembly_objects/* built-ins/vidbibc.o built-ins/netbibc.o --no-whole-archive --start-group \
# lib/lib.a.o arch/x86/lib/lib.a.o built-ins/objects/libx_objects/* .tmp_kallsyms2.o --end-group -o vmlinux