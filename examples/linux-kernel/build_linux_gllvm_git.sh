#!/usr/bin/env bash

### building from a git clone with gllvm

go get github.com/SRI-CSL/gllvm/cmd/...

cd ${HOME}
git clone git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux-stable.git

cd linux-stable
git checkout tags/v4.14.39
cp /vagrant/tinyconfig64 .config

make CC=gclang HOSTCC=gclang

get-bc -m -b built-in.o
get-bc -m vmlinux
