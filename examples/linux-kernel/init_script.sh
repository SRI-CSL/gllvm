#!/usr/bin/env bash

### creating the folder architecture necessary for the kernel build

cd $HOME
mkdir standalone-build
sudo cp /vagrant/convert-thin-archive.sh /usr/bin/

cd standalone-build
cp /vagrant/copy-missing-o.sh .
cp /vagrant/copy-native-bi.sh .
cp /vagrant/handle-bi.sh .
cp /vagrant/copy.sh .

mkdir -p arch/x86/lib
mkdir -p arch/x86/kernel

mkdir -p built-ins/objects/lib_assembly_objects
mkdir -p built-ins/objects/arch_assembly_objects
mkdir -p built-ins/objects/xlib_assembly_objects
mkdir -p built-ins/objects/pow_assembly_objects
mkdir -p built-ins/objects/ker_objects
mkdir -p built-ins/objects/libx_objects

mkdir -p built-ins/fs/objects

mkdir lib/

bash /vagrant/build_linux_gllvm.sh

bash copy.sh

#bash /vagrant/bootable-kernel.sh
