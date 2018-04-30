#!/usr/bin/env bash

sudo apt-get update

sudo apt-get install -y emacs24 dbus-x11 
sudo apt-get install -y git subversion wget 
sudo apt-get install -y llvm-5.0 libclang-5.0-dev clang-5.0
sudo apt-get install -y python-pip golang-go
sudo apt-get install -y flex bison bc libncurses5-dev
sudo apt-get install -y libelf-dev libssl-dev

sudo update-alternatives --install /usr/bin/clang++ clang++ /usr/bin/clang++-5.0 1000
sudo update-alternatives --install /usr/bin/clang clang /usr/bin/clang-5.0 1000
sudo update-alternatives --install /usr/bin/llvm-dis++ llvm-dis /usr/bin/llvm-dis-5.0 1000
sudo update-alternatives --install /usr/bin/llvm-dis llvm-dis /usr/bin/llvm-dis-5.0 1000
sudo update-alternatives --install /usr/bin/llvm-ar llvm-ar /usr/bin/llvm-ar-5.0 1000
sudo update-alternatives --install /usr/bin/llvm-link llvm-link /usr/bin/llvm-link-5.0 1000
sudo update-alternatives --install /usr/bin/llvm-config llvm-config /usr/bin/llvm-config-5.0 1000

cp /vagrant/bash_profile ~/.bash_profile

echo ". /vagrant/bash_profile" >> /home/vagrant/.bashrc
