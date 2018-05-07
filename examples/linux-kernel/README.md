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

# vagrant bootstrapping file

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
#### /vagrant/bash_profile

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

The file [`tinyconfig64`](https://github.com/SRI-CSL/gllvm/blob/master/examples/linux-kernel/tinyconfig64) is generated
by `make tinyconfig` and the using `make menuconfig` to specialize the build to 64 bits. 

## The Tarball Build with gllvm

The build process is carried out by running the `build_linux_gllvm.sh`
script.

```bash
#!/usr/bin/env bash

### building from a tarball with gllvm

go get github.com/SRI-CSL/gllvm/cmd/...

cd ${HOME}
wget https://cdn.kernel.org/pub/linux/kernel/v4.x/linux-4.14.39.tar.xz
tar xvf linux-4.14.39.tar.xz
cd linux-4.14.39

cp /vagrant/tinyconfig64 .config

make CC=gclang HOSTCC=gclang

get-bc -m -b built-in.o
get-bc -m vmlinux
```

## The Tarball Build with wllvm

The build process is carried out by running the `build_linux_wllvm.sh`
script.

```bash
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
```

## Building from a git clone

You can also build from a [git clone using gllvm,](https://github.com/SRI-CSL/gllvm/blob/master/examples/linux-kernel/build_linux_gllvm_git.sh)
or build from a [git clone using wllvm.](https://github.com/SRI-CSL/gllvm/blob/master/examples/linux-kernel/build_linux_wllvm_git.sh)
Though using a tarball is faster, and seemingly more reliable.

