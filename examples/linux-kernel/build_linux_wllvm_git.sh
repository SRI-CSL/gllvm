#!/usr/bin/env bash

sudo pip install wllvm

cd ${HOME}
git clone git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux-stable.git

cd linux-stable
git checkout tags/v4.14.34
cp /vagrant/tinyconfig64 .config


make CC=wllvm HOSTCC=wllvm

extract-bc -m -b built-in.o
extract-bc -m vmlinux
