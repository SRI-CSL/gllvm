#!/usr/bin/env bash

### building from a tarball with gllvm

go get github.com/SRI-CSL/gllvm/cmd/...

cd ${HOME}
wget https://cdn.kernel.org/pub/linux/kernel/v4.x/linux-4.14.39.tar.xz
tar xf linux-4.14.39.tar.xz
mv linux-4.14.39 linux-stable
cd linux-stable

cp /vagrant/link-vmlinux.sh scripts/ #to retain a copy of kallsyms.o
cp /vagrant/parse-bi.py .
cp /vagrant/make-script.sh .

make defconfig
bash make-script.sh
