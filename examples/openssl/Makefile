OPENSSL=openssl-1.1.1c
SRC="https://www.openssl.org/source/${OPENSSL}.tar.gz"
all: libcrypto.bca libssl.bca

${OPENSSL}.tar.gz:
	wget ${SRC}

${OPENSSL}.tar: ${OPENSSL}.tar.gz
	gunzip -k ${OPENSSL}.tar.gz

${OPENSSL}: ${OPENSSL}.tar
	tar xvf ${OPENSSL}.tar

openssl: ${OPENSSL}
	cd ${OPENSSL}; CC=gclang ./config no-asm no-shared; make

libcrypto.bca: openssl
	get-bc -o libcrypto.bc -m -b ${OPENSSL}/libcrypto.a

libssl.bca: openssl
	get-bc -o libssl.bc -m -b ${OPENSSL}/libssl.a

clean:
	rm -f *.bc *.llvm.manifest
	make -C ${OPENSSL} clean


spotless: clean
	rm -rf ${OPENSSL} ${OPENSSL}.tar.gz ${OPENSSL}.tar
