cd $HOME
tar xf linux-4.14.39.tar.xz
mv linux-4.14.39 bootable-linux

cp /vagrant/make-script-clang.sh  bootable-linux/
cd bootable-linux
bash make-script-clang.sh
