VERSION=8.31
SRC=coreutils-${VERSION}
TAR=${SRC}.tar.xz
URL=https://ftp.gnu.org/gnu/coreutils/${TAR}


all: touch.bc

${TAR}:
	wget ${URL}

${SRC}: ${TAR}
	tar xvf ${TAR}
	touch ${SRC}


${SRC}/Makefile: ${SRC}
	cd ${SRC}; CC=gclang ./configure

${SRC}/src/touch: ${SRC}/Makefile
	cd ${SRC}; make
	touch ${SRC}/src/touch

touch.bc: ${SRC}/src/touch
	get-bc -o touch.bc -m ${SRC}/src/touch

touch_from_bitcode: touch.bc
	clang++ touch.bc -o touch

clean:
	rm -rf touch  touch.bc  touch.bc.llvm.manifest
	make -C ${SRC}

spotless:
	rm -rf ${SRC} ${TAR}

