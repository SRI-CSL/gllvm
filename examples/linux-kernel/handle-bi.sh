#!/usr/bin/env bash

### converts the built-in.o files from the different folders into bitcode and copies them to the build folder

export build_home=$HOME/standalone-build
export ker=$HOME/linux-stable

bash copy-native-bi.sh

cd $ker
get-bc -b arch/x86/built-in.o
get-bc -b arch/x86/lib/built-in.o
get-bc -b drivers/built-in.o
get-bc -b fs/built-in.o
get-bc -b kernel/built-in.o
get-bc -b lib/built-in.o
get-bc -b mm/built-in.o
get-bc -b security/built-in.o
get-bc -b init/built-in.o
get-bc -b sound/built-in.o
get-bc -b net/built-in.o
get-bc -b ipc/built-in.o
get-bc -b crypto/built-in.o
get-bc -b block/built-in.o
get-bc -b lib/lib.a
get-bc -b arch/x86/lib/lib.a
get-bc -b arch/x86/pci/built-in.o
get-bc -b arch/x86/video/built-in.o
get-bc -b arch/x86/power/built-in.o


cd $build_home
cp $ker/arch/x86/built-in.o.a.bc ./built-ins/arcbi.o.bc
cp $ker/arch/x86/lib/built-in.o.a.bc ./built-ins/xlibbi.o.bc
cp $ker/drivers/built-in.o.a.bc ./built-ins/dribi.o.bc
cp $ker/fs/built-in.o.a.bc ./built-ins/fsbi.o.bc
cp $ker/kernel/built-in.o.a.bc ./built-ins/kerbi.o.bc
cp $ker/lib/built-in.o.a.bc ./built-ins/libbi.o.bc
cp $ker/init/built-in.o.a.bc ./built-ins/inibi.o.bc
cp $ker/mm/built-in.o.a.bc ./built-ins/mmbi.o.bc
cp $ker/security/built-in.o.a.bc ./built-ins/secbi.o.bc
cp $ker/sound/built-in.o.a.bc ./built-ins/sndbi.o.bc
cp $ker/net/built-in.o.a.bc ./built-ins/netbi.o.bc
cp $ker/ipc/built-in.o.a.bc ./built-ins/ipcbi.o.bc
cp $ker/crypto/built-in.o.a.bc ./built-ins/cptbi.o.bc
cp $ker/block/built-in.o.a.bc ./built-ins/blkbi.o.bc
cp $ker/lib/lib.a.bc ./lib/lib.a.bc
cp $ker/arch/x86/lib/lib.a.bc arch/x86/lib/lib.a.bc
cp $ker/arch/x86/pci/built-in.o.a.bc ./built-ins/pcibi.o.bc
cp $ker/arch/x86/video/built-in.o.a.bc ./built-ins/vidbi.o.bc
cp $ker/arch/x86/power/built-in.o.a.bc ./built-ins/powbi.o.bc

bash copy-missing-o.sh