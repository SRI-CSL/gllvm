#!/usr/bin/env bash

### Copy vmlinux into the bootable linux folder and install the new kernel
cp $HOME/standalone-build/vmlinux $HOME/linux-stable/

cd $HOME/linux-stable

scripts/sortextable vmlinux
nm -n vmlinux | grep -v '\( [aNUw] \)\|\(__crc_\)\|\( \$[adt]\)\|\( .L\)' > System.map
make CC=gclang HOSTCC=gclang
sudo make modules_install install
