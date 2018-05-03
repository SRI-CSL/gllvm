# Building a recent Linux Kernel.

In this directory we include all the necessary files needed to
build the kernel in a Ubuntu 16.04 vagrant box. We will guide the reader through
the relatively simple task. We assume familiarity with [Vagrant.](https://www.vagrantup.com/)

## Vagrantfile

```ruby
# -*- mode: ruby -*-
# vi: set ft=ruby :


Vagrant.configure("2") do |config|
  
  config.vm.box = "ubuntu/xenial64"
  config.vm.provision :shell, path: "bootstrap.sh"

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "4096"
    vb.customize ["modifyvm", :id, "--ioapic", "on"]
    vb.customize ["modifyvm", :id, "--memory", "4096"]
    vb.customize ["modifyvm", :id, "--cpus", "4"]
   end

end
```

## Bootstrapping 

```bash
#!/usr/bin/env bash

sudo apt-get update

sudo apt-get install -y emacs24 dbus-x11 
sudo apt-get install -y git
sudo apt-get install -y llvm-5.0 libclang-5.0-dev clang-5.0
sudo apt-get install -y python-pip golang-go
sudo apt-get install -y flex bison bc libncurses5-dev
sudo apt-get install -y libelf-dev libssl-dev

echo ". /vagrant/bash_profile" >> /home/vagrant/.bashrc
```

## Shell Settings

```bash
####  llvm
export LLVM_HOME=/usr/lib/llvm-5.0
export GOPATH=/vagrant/go

######## gllvm/wllvm configuration #############

export LLVM_COMPILER=clang
export WLLVM_OUTPUT_LEVEL=WARNING
export WLLVM_OUTPUT_FILE=/vagrant/wrapper.log
export PATH=${GOPATH}/bin:${LLVM_HOME}/bin:${PATH}
```



## Configuration stuff.

The file `tinyconfig64` is generated ...

## The Build with gllvm

The build process is carried out by running the `build_linux_gllvm.sh`
script.

```bash
#!/usr/bin/env bash

mkdir -p ${GOPATH}
go get github.com/SRI-CSL/gllvm/cmd/...

mkdir ${HOME}/linux_kernel
cd ${HOME}/linux_kernel
git clone git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux-stable.git

cd linux-stable
git checkout tags/v4.14.34
cp /vagrant/tinyconfig64 .config

make CC=gclang HOSTCC=gclang

get-bc -m -b built-in.o
get-bc -m vmlinux
```

## The Build with wllvm

The build process is carried out by running the `build_linux_wllvm.sh`
script.

```bash
#!/usr/bin/env bash

sudo pip install wllvm

mkdir ${HOME}/linux_kernel
cd ${HOME}/linux_kernel
git clone git://git.kernel.org/pub/scm/linux/kernel/git/stable/linux-stable.git

cd linux-stable
git checkout tags/v4.14.34
cp /vagrant/tinyconfig64 .config


make CC=wllvm HOSTCC=wllvm

extract-bc -m -b built-in.o
extract-bc -m vmlinux
```

## Extracting the bitcode

