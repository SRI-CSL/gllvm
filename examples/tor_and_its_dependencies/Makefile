#iam wrote this monstrosity

THISDIR= $(shell pwd)


ZLIB=zlib-1.2.11
ZLIB_TARBALL=${ZLIB}.tar.gz
ZLIB_URL=https://zlib.net/${ZLIB_TARBALL}
ZLIB_INSTALL=${THISDIR}/install/zlib

LIBEVENT=libevent-2.1.8-stable
LIBEVENT_TARBALL=${LIBEVENT}.tar.gz
LIBEVENT_URL=https://github.com/libevent/libevent/releases/download/release-2.1.8-stable/${LIBEVENT_TARBALL}
LIBEVENT_INSTALL=${THISDIR}/install/libevent

OPENSSL=openssl-1.1.0h
OPENSSL_TARBALL=${OPENSSL}.tar.gz
OPENSSL_URL=https://www.openssl.org/source/${OPENSSL_TARBALL}
OPENSSL_INSTALL=${THISDIR}/install/openssl


TOR=tor-0.3.3.6
TOR_TARBALL=${TOR}.tar.gz
TOR_URL=https://www.torproject.org/dist/${TOR_TARBALL}
TOR_INSTALL=${THISDIR}/install/tor


#LLVMCC=wllvm
#LLVMGET=extract-bc

LLVMCC=gclang
LLVMGET=get-bc

BITCODE = libz.a.bc libssl.a.bc libevent.a.bc tor.bc

all: ${BITCODE}

${TOR_TARBALL}:
	wget ${TOR_URL}

${TOR}: ${TOR_TARBALL}
	tar xvfz ${TOR_TARBALL}

${TOR_INSTALL}/bin/tor: ${TOR}
	cd ${TOR}; CC=${LLVMCC} ./configure --prefix=${TOR_INSTALL}; make && make install

tor.bc: ${TOR_INSTALL}/bin/tor
	cd ${TOR_INSTALL}/bin; ${LLVMGET} tor; mv tor.bc ${THISDIR}


${ZLIB_TARBALL}:
	wget ${ZLIB_URL}

${ZLIB}: ${ZLIB_TARBALL}
	tar xvfz ${ZLIB_TARBALL}

${ZLIB_INSTALL}/lib/libz.a: ${ZLIB}
	cd ${ZLIB}; CC=${LLVMCC} ./configure --prefix=${ZLIB_INSTALL}; make && make install

libz.a.bc: ${ZLIB_INSTALL}/lib/libz.a
	cd ${ZLIB_INSTALL}/lib; ${LLVMGET} -b libz.a; mv libz.a.bc ${THISDIR}


${OPENSSL_TARBALL}:
	wget ${OPENSSL_URL}

${OPENSSL}: ${OPENSSL_TARBALL}
	tar xvfz ${OPENSSL_TARBALL}

${OPENSSL_INSTALL}/lib/libssl.a: ${OPENSSL}
	cd ${OPENSSL}; CC=${LLVMCC} ./config --prefix=${OPENSSL_INSTALL} --openssldir=${OPENSSL_INSTALL}; make && make install

libssl.a.bc: ${OPENSSL_INSTALL}/lib/libssl.a
	cd ${OPENSSL_INSTALL}/lib; ${LLVMGET} -b libssl.a; mv libssl.a.bc ${THISDIR}


${LIBEVENT_TARBALL}:
	wget ${LIBEVENT_URL}

${LIBEVENT}: ${LIBEVENT_TARBALL}
	tar xvfz ${LIBEVENT_TARBALL}

${LIBEVENT_INSTALL}/lib/libevent.a: ${LIBEVENT}
	cd ${LIBEVENT}; CC=${LLVMCC} CFLAGS="-I${OPENSSL_INSTALL}/include" LDFLAGS="-L${OPENSSL_INSTALL}/lib" ./configure --prefix=${LIBEVENT_INSTALL}; make && make install

libevent.a.bc: ${LIBEVENT_INSTALL}/lib/libevent.a
	cd ${LIBEVENT_INSTALL}/lib; ${LLVMGET} -b libevent.a; mv libevent.a.bc ${THISDIR}



clean:
	rm ${BITCODE}
	make -C ${ZLIB} clean
	make -C ${LIBEVENT} clean
	make -C ${OPENSSL} clean
	make -C ${TOR} clean

very_clean: clean
	rm -rf install ${TOR} ${ZLIB} ${OPENSSL} ${LIBEVENT}

spotless: very_clean
	rm -rf ${TOR_TARBALL} ${ZLIB_TARBALL} ${OPENSSL_TARBALL} ${LIBEVENT_TARBALL}
