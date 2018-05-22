#!/usr/bin/env bash

### copy the libraries compiled by clang to the build folder.
### currently we are only using the driver built-in.o, so this is mostly unnecessary
export build_home=$HOME/standalone-build
export ker=$HOME/linux-stable

convert-thin-archive.sh $ker/arch/x86/built-in.o
convert-thin-archive.sh $ker/arch/x86/lib/built-in.o
convert-thin-archive.sh $ker/drivers/built-in.o
convert-thin-archive.sh $ker/fs/built-in.o
convert-thin-archive.sh $ker/kernel/built-in.o
convert-thin-archive.sh $ker/lib/built-in.o
convert-thin-archive.sh $ker/mm/built-in.o
convert-thin-archive.sh $ker/security/built-in.o
convert-thin-archive.sh $ker/init/built-in.o
convert-thin-archive.sh $ker/sound/built-in.o
convert-thin-archive.sh $ker/net/built-in.o
convert-thin-archive.sh $ker/ipc/built-in.o
convert-thin-archive.sh $ker/crypto/built-in.o
convert-thin-archive.sh $ker/block/built-in.o
convert-thin-archive.sh $ker/lib/lib.a
convert-thin-archive.sh $ker/arch/x86/lib/lib.a
convert-thin-archive.sh $ker/arch/x86/pci/built-in.o
convert-thin-archive.sh $ker/arch/x86/video/built-in.o
convert-thin-archive.sh $ker/arch/x86/power/built-in.o



cp $ker/arch/x86/built-in.o.new ./built-ins/arcbi.o
cp $ker/arch/x86/lib/built-in.o.new ./built-ins/xlibbi.o
cp $ker/drivers/built-in.o.new ./built-ins/dribi.o
cp $ker/fs/built-in.o.new ./built-ins/fsbi.o
cp $ker/kernel/built-in.o.new ./built-ins/kerbi.o
cp $ker/lib/built-in.o.new ./built-ins/libbi.o
cp $ker/init/built-in.o.new ./built-ins/inibi.o
cp $ker/mm/built-in.o.new ./built-ins/mmbi.o
cp $ker/security/built-in.o.new ./built-ins/secbi.o
cp $ker/sound/built-in.o.new ./built-ins/sndbi.o
cp $ker/net/built-in.o.new ./built-ins/netbi.o
cp $ker/ipc/built-in.o.new ./built-ins/ipcbi.o
cp $ker/crypto/built-in.o.new ./built-ins/cptbi.o
cp $ker/block/built-in.o.new ./built-ins/blkbi.o
cp $ker/lib/lib.a.new ./lib/lib.a
cp $ker/arch/x86/lib/lib.a.new arch/x86/lib/lib.a
cp $ker/arch/x86/pci/built-in.o.new ./built-ins/pcibi.o
cp $ker/arch/x86/video/built-in.o.new ./built-ins/vidbi.o
cp $ker/arch/x86/power/built-in.o.new ./built-ins/powbi.o
