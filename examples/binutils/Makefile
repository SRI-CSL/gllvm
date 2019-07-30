#http://ftp.gnu.org/gnu/binutils/binutils-2.32.tar.gz

SRC=binutils-2.32
TAR=${SRC}.tar.gz
URL=http://ftp.gnu.org/gnu/binutils/${TAR}


all: readelf.bc objdump.bc readelf_from_bitcode objdump_from_bitcode


${TAR}:
	wget ${URL}


${SRC}: ${TAR}
	tar xvfz ${TAR}

${SRC}/binutils/readelf: ${SRC}
	cd ${SRC}; CC=gclang ./configure; CC=gclang make;

${SRC}/binutils/objdump: ${SRC}
	cd ${SRC}; CC=gclang ./configure; CC=gclang make;

readelf.bc: ${SRC}/binutils/readelf
	get-bc -m -o readelf.bc ${SRC}/binutils/readelf

objdump.bc: ${SRC}/binutils/objdump
	get-bc -m -o objdump.bc ${SRC}/binutils/objdump


readelf_modules:
	mkdir readelf_modules
	cp ${SRC}/binutils/.readelf.o.bc readelf_modules/readelf.o.bc
	cp ${SRC}/binutils/.version.o.bc readelf_modules/version.o.bc
	cp ${SRC}/binutils/.unwind-ia64.o.bc readelf_modules/unwind-ia64.o.bc
	cp ${SRC}/binutils/.dwarf.o.bc readelf_modules/dwarf.o.bc
	cp ${SRC}/binutils/.elfcomm.o.bc readelf_modules/elfcomm.o.bc
	cp ${SRC}/libiberty/.argv.o.bc readelf_modules/argv.o.bc
	cp ${SRC}/libiberty/.concat.o.bc readelf_modules/concat.o.bc
	cp ${SRC}/libiberty/.dwarfnames.o.bc readelf_modules/dwarfnames.o.bc
	cp ${SRC}/libiberty/.lbasename.o.bc readelf_modules/lbasename.o.bc
	cp ${SRC}/libiberty/.lrealpath.o.bc readelf_modules/lrealpath.o.bc
	cp ${SRC}/libiberty/.safe-ctype.o.bc readelf_modules/safe-ctype.o.bc
	cp ${SRC}/libiberty/.xexit.o.bc readelf_modules/xexit.o.bc
	cp ${SRC}/libiberty/.xmalloc.o.bc readelf_modules/xmalloc.o.bc
	cp ${SRC}/libiberty/.xstrdup.o.bc readelf_modules/xstrdup.o.bc
	cp ${SRC}/zlib/.libz_a-inflate.o.bc readelf_modules/libz_a-inflate.o.bc
	cp ${SRC}/zlib/.libz_a-inftrees.o.bc readelf_modules/libz_a-inftrees.o.bc
	cp ${SRC}/zlib/.libz_a-zutil.o.bc readelf_modules/libz_a-zutil.o.bc
	cp ${SRC}/zlib/.libz_a-adler32.o.bc readelf_modules/libz_a-adler32.o.bc
	cp ${SRC}/zlib/.libz_a-crc32.o.bc readelf_modules/libz_a-crc32.o.bc
	cp ${SRC}/zlib/.libz_a-inffast.o.bc readelf_modules/libz_a-inffast.o.bc


objdump_modules:
	mkdir objdump_modules
	cp ${SRC}/binutils/.objdump.o.bc objdump_modules/objdump.o.bc
	cp ${SRC}/binutils/.dwarf.o.bc objdump_modules/dwarf.o.bc
	cp ${SRC}/binutils/.prdbg.o.bc objdump_modules/prdbg.o.bc
	cp ${SRC}/binutils/.rddbg.o.bc objdump_modules/rddbg.o.bc
	cp ${SRC}/binutils/.debug.o.bc objdump_modules/debug.o.bc
	cp ${SRC}/binutils/.stabs.o.bc objdump_modules/binutils_stabs.o.bc
	cp ${SRC}/binutils/.rdcoff.o.bc objdump_modules/rdcoff.o.bc
	cp ${SRC}/binutils/.bucomm.o.bc objdump_modules/bucomm.o.bc
	cp ${SRC}/binutils/.version.o.bc objdump_modules/version.o.bc
	cp ${SRC}/binutils/.filemode.o.bc objdump_modules/filemode.o.bc
	cp ${SRC}/binutils/.elfcomm.o.bc objdump_modules/elfcomm.o.bc
	cp ${SRC}/opcodes/.disassemble.o.bc objdump_modules/disassemble.o.bc
	cp ${SRC}/opcodes/.dis-init.o.bc objdump_modules/dis-init.o.bc
	cp ${SRC}/opcodes/.i386-dis.o.bc objdump_modules/i386-dis.o.bc
	cp ${SRC}/opcodes/.dis-buf.o.bc objdump_modules/dis-buf.o.bc
	cp ${SRC}/bfd/.archive.o.bc objdump_modules/archive.o.bc
	cp ${SRC}/bfd/.archures.o.bc objdump_modules/archures.o.bc
	cp ${SRC}/bfd/.bfd.o.bc objdump_modules/bfd.o.bc
	cp ${SRC}/bfd/.bfdio.o.bc objdump_modules/bfdio.o.bc
	cp ${SRC}/bfd/.coff-bfd.o.bc objdump_modules/coff-bfd.o.bc
	cp ${SRC}/bfd/.compress.o.bc objdump_modules/compress.o.bc
	cp ${SRC}/bfd/.elf-properties.o.bc objdump_modules/elf-properties.o.bc
	cp ${SRC}/bfd/.format.o.bc objdump_modules/format.o.bc
	cp ${SRC}/bfd/.hash.o.bc objdump_modules/hash.o.bc
	cp ${SRC}/bfd/.init.o.bc objdump_modules/init.o.bc
	cp ${SRC}/bfd/.libbfd.o.bc objdump_modules/libbfd.o.bc
	cp ${SRC}/bfd/.opncls.o.bc objdump_modules/opncls.o.bc
	cp ${SRC}/bfd/.section.o.bc objdump_modules/section.o.bc
	cp ${SRC}/bfd/.simple.o.bc objdump_modules/simple.o.bc
	cp ${SRC}/bfd/.stab-syms.o.bc objdump_modules/stab-syms.o.bc
	cp ${SRC}/bfd/.syms.o.bc objdump_modules/syms.o.bc
	cp ${SRC}/bfd/.targets.o.bc objdump_modules/targets.o.bc
	cp ${SRC}/bfd/.binary.o.bc objdump_modules/binary.o.bc
	cp ${SRC}/bfd/.ihex.o.bc objdump_modules/ihex.o.bc
	cp ${SRC}/bfd/.srec.o.bc objdump_modules/srec.o.bc
	cp ${SRC}/bfd/.tekhex.o.bc objdump_modules/tekhex.o.bc
	cp ${SRC}/bfd/.verilog.o.bc objdump_modules/verilog.o.bc
	cp ${SRC}/bfd/.elf64-x86-64.o.bc objdump_modules/elf64-x86-64.o.bc
	cp ${SRC}/bfd/.elfxx-x86.o.bc objdump_modules/elfxx-x86.o.bc
	cp ${SRC}/bfd/.elf-ifunc.o.bc objdump_modules/elf-ifunc.o.bc
	cp ${SRC}/bfd/.elf-nacl.o.bc objdump_modules/elf-nacl.o.bc
	cp ${SRC}/bfd/.elf64.o.bc objdump_modules/elf64.o.bc
	cp ${SRC}/bfd/.elf.o.bc objdump_modules/elf.o.bc
	cp ${SRC}/bfd/.elflink.o.bc objdump_modules/elflink.o.bc
	cp ${SRC}/bfd/.elf-attrs.o.bc objdump_modules/elf-attrs.o.bc
	cp ${SRC}/bfd/.elf-strtab.o.bc objdump_modules/elf-strtab.o.bc
	cp ${SRC}/bfd/.elf-eh-frame.o.bc objdump_modules/elf-eh-frame.o.bc
	cp ${SRC}/bfd/.dwarf1.o.bc objdump_modules/dwarf1.o.bc
	cp ${SRC}/bfd/.dwarf2.o.bc objdump_modules/dwarf2.o.bc
	cp ${SRC}/bfd/.elf32-i386.o.bc objdump_modules/elf32-i386.o.bc
	cp ${SRC}/bfd/.elf-vxworks.o.bc objdump_modules/elf-vxworks.o.bc
	cp ${SRC}/bfd/.elf32.o.bc objdump_modules/elf32.o.bc
	cp ${SRC}/bfd/.pei-i386.o.bc objdump_modules/pei-i386.o.bc
	cp ${SRC}/bfd/.peigen.o.bc objdump_modules/peigen.o.bc
	cp ${SRC}/bfd/.cofflink.o.bc objdump_modules/cofflink.o.bc
	cp ${SRC}/bfd/.coffgen.o.bc objdump_modules/coffgen.o.bc
	cp ${SRC}/bfd/.pei-x86_64.o.bc objdump_modules/pei-x86_64.o.bc
	cp ${SRC}/bfd/.pex64igen.o.bc objdump_modules/pex64igen.o.bc
	cp ${SRC}/bfd/.elf64-gen.o.bc objdump_modules/elf64-gen.o.bc
	cp ${SRC}/bfd/.elf32-gen.o.bc objdump_modules/elf32-gen.o.bc
	cp ${SRC}/bfd/.plugin.o.bc objdump_modules/plugin.o.bc
	cp ${SRC}/bfd/.cpu-i386.o.bc objdump_modules/cpu-i386.o.bc
	cp ${SRC}/bfd/.cpu-iamcu.o.bc objdump_modules/cpu-iamcu.o.bc
	cp ${SRC}/bfd/.cpu-l1om.o.bc objdump_modules/cpu-l1om.o.bc
	cp ${SRC}/bfd/.cpu-k1om.o.bc objdump_modules/cpu-k1om.o.bc
	cp ${SRC}/bfd/.cpu-plugin.o.bc objdump_modules/cpu-plugin.o.bc
	cp ${SRC}/bfd/.archive64.o.bc objdump_modules/archive64.o.bc
	cp ${SRC}/bfd/.cache.o.bc objdump_modules/cache.o.bc
	cp ${SRC}/bfd/.corefile.o.bc objdump_modules/corefile.o.bc
	cp ${SRC}/bfd/.linker.o.bc objdump_modules/linker.o.bc
	cp ${SRC}/bfd/.merge.o.bc objdump_modules/merge.o.bc
	cp ${SRC}/bfd/.reloc.o.bc objdump_modules/reloc.o.bc
	cp ${SRC}/bfd/.stabs.o.bc objdump_modules/bfd_stabs.o.bc
	cp ${SRC}/zlib/.libz_a-compress.o.bc objdump_modules/libz_a-compress.o.bc
	cp ${SRC}/zlib/.libz_a-deflate.o.bc objdump_modules/libz_a-deflate.o.bc
	cp ${SRC}/zlib/.libz_a-inflate.o.bc objdump_modules/libz_a-inflate.o.bc
	cp ${SRC}/zlib/.libz_a-inftrees.o.bc objdump_modules/libz_a-inftrees.o.bc
	cp ${SRC}/zlib/.libz_a-trees.o.bc objdump_modules/libz_a-trees.o.bc
	cp ${SRC}/zlib/.libz_a-zutil.o.bc objdump_modules/libz_a-zutil.o.bc
	cp ${SRC}/zlib/.libz_a-adler32.o.bc objdump_modules/libz_a-adler32.o.bc
	cp ${SRC}/zlib/.libz_a-crc32.o.bc objdump_modules/libz_a-crc32.o.bc
	cp ${SRC}/zlib/.libz_a-inffast.o.bc objdump_modules/libz_a-inffast.o.bc
	cp ${SRC}/libiberty/.cplus-dem.o.bc objdump_modules/cplus-dem.o.bc
	cp ${SRC}/libiberty/.cp-demangle.o.bc objdump_modules/cp-demangle.o.bc
	cp ${SRC}/libiberty/.argv.o.bc objdump_modules/argv.o.bc
	cp ${SRC}/libiberty/.concat.o.bc objdump_modules/concat.o.bc
	cp ${SRC}/libiberty/.cp-demint.o.bc objdump_modules/cp-demint.o.bc
	cp ${SRC}/libiberty/.d-demangle.o.bc objdump_modules/d-demangle.o.bc
	cp ${SRC}/libiberty/.dwarfnames.o.bc objdump_modules/dwarfnames.o.bc
	cp ${SRC}/libiberty/.filename_cmp.o.bc objdump_modules/filename_cmp.o.bc
	cp ${SRC}/libiberty/.getpwd.o.bc objdump_modules/getpwd.o.bc
	cp ${SRC}/libiberty/.hashtab.o.bc objdump_modules/hashtab.o.bc
	cp ${SRC}/libiberty/.hex.o.bc objdump_modules/hex.o.bc
	cp ${SRC}/libiberty/.lbasename.o.bc objdump_modules/lbasename.o.bc
	cp ${SRC}/libiberty/.lrealpath.o.bc objdump_modules/lrealpath.o.bc
	cp ${SRC}/libiberty/.make-relative-prefix.o.bc objdump_modules/make-relative-prefix.o.bc
	cp ${SRC}/libiberty/.make-temp-file.o.bc objdump_modules/make-temp-file.o.bc
	cp ${SRC}/libiberty/.objalloc.o.bc objdump_modules/objalloc.o.bc
	cp ${SRC}/libiberty/.rust-demangle.o.bc objdump_modules/rust-demangle.o.bc
	cp ${SRC}/libiberty/.safe-ctype.o.bc objdump_modules/safe-ctype.o.bc
	cp ${SRC}/libiberty/.unlink-if-ordinary.o.bc objdump_modules/unlink-if-ordinary.o.bc
	cp ${SRC}/libiberty/.xexit.o.bc objdump_modules/xexit.o.bc
	cp ${SRC}/libiberty/.xmalloc.o.bc objdump_modules/xmalloc.o.bc
	cp ${SRC}/libiberty/.xstrdup.o.bc objdump_modules/xstrdup.o.bc
	cp ${SRC}/libiberty/.xstrerror.o.bc objdump_modules/xstrerror.o.bc


readelf_from_bitcode: readelf_modules
	clang readelf_modules/*.bc -o readelf_from_bitcode

objdump_from_bitcode: objdump_modules
	clang objdump_modules/*.bc -ldl -o objdump_from_bitcode


clean:
	rm -rf *~ *.bc *.llvm.manifest objdump_modules readelf_modules
	make -C ${SRC} clean


spotless: clean
	rm -rf ${TAR} ${SRC}
