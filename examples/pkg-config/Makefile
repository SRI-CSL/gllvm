# The latest is 0.29.2 as of 10/29/2019
#
# https://www.freedesktop.org/wiki/Software/pkg-config/
#
VERSION=0.29.2
SRC=pkg-config-${VERSION}
TARBALL=${SRC}.tar.gz
URL=https://pkg-config.freedesktop.org/releases/${TARBALL}

all: pkg-config.bc

${TARBALL}:
	wget ${URL}

${SRC}: ${TARBALL}
	tar xvfz ${TARBALL}
	touch ${SRC}

${SRC}/pkg-config:
	cd ${SRC}; CC=gclang ./configure --with-internal-glib; make

pkg-config.bc: ${SRC}/pkg-config
	get-bc -o pkg-config.bc ${SRC}/pkg-config

clean:
	make -C  ${SRC} clean


spotless:
	rm -rf ${SRC} ${TARBALL} pkg-config.bc
