
# https://gist.github.com/ianamason/1c3a54697465436dce2822b80ee364d1

VERSION=v10.16.0
SRC=node-${VERSION}
TAR=${SRC}.tar.gz
URL=https://nodejs.org/dist/${VERSION}/${TAR}

all: node_from_bitcode


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

libnode.bc: ${SRC}/node
	get-bc -b -o libnode.bc ${SRC}/out/Release/obj.target/libnode.a

libzlib.bc: ${SRC}/node
	get-bc -b -o libzlib.bc ${SRC}/out/Release/obj.target/deps/zlib/libzlib.a

libuv.bc: ${SRC}/node
	get-bc -b -o libuv.bc ${SRC}/out/Release/obj.target/deps/uv/libuv.a

libv8_base.bc: ${SRC}/node
	get-bc -b -o libv8_base.bc ${SRC}/out/Release/obj.target/deps/v8/gypfiles/libv8_base.a

libopenssl.bc: ${SRC}/node
	get-bc -b -o libopenssl.bc ${SRC}/out/Release/obj.target/deps/openssl/libopenssl.a

node_main.bc: ${SRC}/node
	get-bc -b -o node_main.bc ${SRC}/out/Release/obj.target/node/src/node_main.o

libv8_libplatform.bc: ${SRC}/node
	get-bc -b -o libv8_libplatform.bc ${SRC}/out/Release/obj.target/deps/v8/gypfiles/libv8_libplatform.a

libicui18n.bc: ${SRC}/node
	get-bc -b -o libicui18n.bc ${SRC}/out/Release/obj.target/tools/icu/libicui18n.a

libhttp_parser.bc: ${SRC}/node
	get-bc -b -o libhttp_parser.bc ${SRC}/out/Release/obj.target/deps/http_parser/libhttp_parser.a

libcares.bc: ${SRC}/node
	get-bc -b -o libcares.bc ${SRC}/out/Release/obj.target/deps/cares/libcares.a

libnghttp2.bc: ${SRC}/node
	get-bc -b -o libnghttp2.bc ${SRC}/out/Release/obj.target/deps/nghttp2/libnghttp2.a

libbrotili.bc: ${SRC}/node
	get-bc -b -o libbrotli.bc ${SRC}/out/Release/obj.target/deps/brotli/libbrotli.a

libv8_libbase.bc: ${SRC}/node
	get-bc -b -o libv8_libbase.bc ${SRC}/out/Release/obj.target/deps/v8/gypfiles/libv8_libbase.a

libv8_sampler.bc: ${SRC}/node
	get-bc -b -o libv8_libsampler.bc ${SRC}/out/Release/obj.target/deps/v8/gypfiles/libv8_libsampler.a

libicuucx.bc: ${SRC}/node
	get-bc -b -o libicuucx.bc ${SRC}/out/Release/obj.target/tools/icu/libicuucx.a

libicudata.bc: ${SRC}/node
	get-bc -b -o libicudata.bc ${SRC}/out/Release/obj.target/tools/icu/libicudata.a

libicustubdata.bc: ${SRC}/node
	get-bc -b -o libicustubdata.bc ${SRC}/out/Release/obj.target/tools/icu/libicustubdata.a

libv8_snapshot.bc: ${SRC}/node
	get-bc -b -o libv8_snapshot.bc ${SRC}/out/Release/obj.target/deps/v8/gypfiles/libv8_snapshot.a

node_from_bitcode: bitcode
	clang++ -o node_from_bitcode -pthread -rdynamic -m64 node_main.bc libnode.bc libv8_libplatform.bc libicui18n.bc libzlib.bc libhttp_parser.bc libcares.bc libuv.bc libnghttp2.bc libbrotli.bc libopenssl.bc libv8_base.bc libv8_libbase.bc libv8_libsampler.bc libicuucx.bc libicudata.bc libicustubdata.bc libv8_snapshot.bc -ldl -lrt -lm

clean:
	make -C ${SRC} clean

spotless:
	rm -rf ${SRC}


.PHONY: bitcode clean spotless
