#!/usr/bin/env bash

### building from a git clone with wllvm

sudo pip install wllvm

cd ${HOME}
git clone git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux-stable.git

cd linux-stable
git checkout tags/v4.14.39
cp /vagrant/tinyconfig64 .config


make CC=wllvm HOSTCC=wllvm

extract-bc -m -b built-in.o
extract-bc -m vmlinux
