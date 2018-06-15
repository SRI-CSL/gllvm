# Building a recent Linux Kernel.

In this directory we include all the necessary files needed to
build the kernel in a Ubuntu 16.04 vagrant box. We will guide the reader through
the relatively simple task. We assume familiarity with [Vagrant.](https://www.vagrantup.com/)

We begin with a warm up exercise that just builds a version of the kernel (based on tinyconfig) that does not boot.
We then beef up the process somewhat so that we produce a bootable kernel.

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
    vb.customize ["modifyvm", :id, "--cpus", "2"]
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
by `make tinyconfig` and then using `make menuconfig` to specialize the build to 64 bits.

## The Tarball Build with gllvm

The build process is carried out by running the `build_linux_gllvm_tarball.sh`
script within the vagrant box, configured as described above.

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


## Comparing the two


`gclang` build:

```
real	2m55.689s
user	4m10.036s
sys     0m34.780s
```

`wllvm` build:
```
real	6m52.443s
user	4m32.124s
sys  	0m44.072s

```


## Building from a git clone

You can also build from a [git clone using gllvm,](https://github.com/SRI-CSL/gllvm/blob/master/examples/linux-kernel/build_linux_gllvm_git.sh)
or build from a [git clone using wllvm.](https://github.com/SRI-CSL/gllvm/blob/master/examples/linux-kernel/build_linux_wllvm_git.sh)
Though using a tarball is faster, and seemingly more reliable.

# Building a Bootable Kernel from the Bitcode


In this section we will describe how to build a bootable kernel from LLVM bitcode.
The [init_script.sh](init_script.sh) script will build a bootable kernel that is constructed from mostly bitcode (drivers and ext4 file system are currently not translated).

The init script first builds the required folder architecture for the build,  and then calls build_linux_gllvm,
only this time with a default configuration instead of tinyconfig.

The copy.sh script will then extract the bitcode from the archives in the linux build folder, and copy them along with necessary object files (the files compiled straight from assembly will not emit a bitcode file).
It will then call the link command on those files and generate a vmlinux executable containing the kernel.

The gclang build of the kernel adds llvm_bc headers to most files, and those mess with the generation of a compressed bootable kernel.
We need to have a separate folder built form clang or gcc on which to finish the kernel build and install.
Finally, calling the install-kernel script will copy the new kernel into the clang generated folder and finish the build and install. Rebooting will be on the bitcode kernel.

NB: I was not able to boot on any custom kernel via Vagrant with a defconfig build.

NB2: On a dedicated VirtualBox machine, the generated kernel boots properly but it may be buggy. Most notably, I have experienced issues when shutting down and booting the machine a second time.

NB3: Some default kernel modules loaded with olddefconfig cannot be compiled with clang due to VLAIS


## Using built-in-parsing.py

Another possibility after building the linux with gclang is running [built-in-parsing.py](built-in-parsing.py) in order to write a script that will do the extracting, copying and linking of bitcode.
This script automates the script-writing process for other configs than defconfig.
Running "python built-in-parsing.py BUILD_PATH drivers fs/ext4" from whithin the kernel folder writes a new build_script.sh with the right instructions to build the kernel in BUILD_PATH.
NB: You will have to set the gclang output file to /vagrant/wrapper-logs/wrapper.log before running the python script.
