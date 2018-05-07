#!/usr/bin/env bash

### building from a tarball with wllvm

sudo pip install wllvm

cd ${HOME}
wget https://cdn.kernel.org/pub/linux/kernel/v4.x/linux-4.14.39.tar.xz
tar xvf linux-4.14.39.tar.xz
cd linux-4.14.39

cp /vagrant/tinyconfig64 .config


make CC=wllvm HOSTCC=wllvm

extract-bc -m -b built-in.o
extract-bc -m vmlinux
