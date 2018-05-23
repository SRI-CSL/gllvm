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
