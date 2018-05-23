export home=/home/pn/perso/bcfull
export ker=/home/pn/perso/linux-stable

cd $ker
python parse-bi.py fs/built-in.o fs/out.sh ../bcfull/instrfs 999

cd $home
bash copy.sh

ld --build-id -T ./arch/x86/kernel/vmlinux.lds --whole-archive built-ins/sep_objs/ker_objects/head_64.o \
built-ins/sep_objs/ker_objects/head64.o built-ins/sep_objs/ker_objects/ebda.o built-ins/sep_objs/ker_objects/platform-quirks.o\
 built-ins/inibibc.o built-ins/sep_objs/ker_objects/initramfs_data.o built-ins/arcbibc.o built-ins/sep_objs/arch_assembly_objects/* \
 built-ins/kerbibc.o built-ins/mmbibc.o \@instrfs built-ins/ipcbibc.o built-ins/secbibc.o built-ins/cptbibc.o built-ins/blkbibc.o \
 built-ins/libbibc.o built-ins/sep_objs/lib_assembly_objects/* built-ins/xlibbibc.o built-ins/sep_objs/xlib_assembly_objects/* \
 built-ins/dribi.o built-ins/sndbibc.o built-ins/pcibibc.o built-ins/powbibc.o built-ins/sep_objs/pow_assembly_objects/* \
 built-ins/vidbibc.o built-ins/netbibc.o --no-whole-archive --start-group lib/lib.a.o arch/x86/lib/lib.a.o \
 built-ins/sep_objs/libx_objects/* .tmp_kallsyms2.o --end-group -o vmlinux
 
cp vmlinux ../linux-stable-clang/
cd ../linux-stable-clang
./install.sh
sudo reboot
