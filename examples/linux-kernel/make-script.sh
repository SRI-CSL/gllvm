#to avoid random missing files errors when building
while [ ! -e "vmlinux" ]; do
    make vmlinux CC=gclang HOSTCC=gclang
done