
# https://gist.github.com/ianamason/1c3a54697465436dce2822b80ee364d1

VERSION=v10.16.0
SRC=node-${VERSION}
TAR=${SRC}.tar.gz
URL=https://nodejs.org/dist/${VERSION}/${TAR}

all: node.bc


bitcode: libnode.bc libzlib.bc libuv.bc libv8_base.bc libopenssl.bc node_main.bc \
         libv8_libplatform.bc libicui18n.bc libhttp_parser.bc libcares.bc libnghttp2.bc \
         libbrotili.bc libv8_libbase.bc libv8_sampler.bc libicuucx.bc libicudata.bc \
         libicustubdata.bc libv8_snapshot.bc


${TAR}:
	wget https://nodejs.org/dist/v10.16.0/node-v10.16.0.tar.gz

${SRC}: ${TAR}
	tar xvfz ${TAR}

${SRC}/Makefile: ${SRC}
	cd ${SRC}; CC=gclang CXX=gclang++ ./configure --openssl-no-asm
	#make sure ${SRC}/node isn't earlier than ${SRC}
	touch ${SRC}/Makefile

${SRC}/node: ${SRC}/Makefile
	cd ${SRC}; CC=gclang CXX=gclang++ make
	#make sure ${SRC}/node isn't earlier than ${SRC}
	touch ${SRC}/node

node.bc: ${SRC}/node
	get-bc -o node.bc ${SRC}/node

clean:
	make -C ${SRC} clean

spotless:
	rm -rf ${SRC}


.PHONY: bitcode clean spotless
