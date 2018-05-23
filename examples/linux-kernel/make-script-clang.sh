#to avoid random missing files errors when building with clang
while [ ! -e "vmlinux" ]; do
    make vmlinux CC=clang HOSTCC=clang
done