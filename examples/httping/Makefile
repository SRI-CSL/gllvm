# sudo apt-get install gettext
SRC=httping-2.5
TAR=${SRC}.tgz
URL=https://www.vanheusden.com/httping/${TAR}


all: httping.bc

${TAR}:
	wget ${URL}

${SRC}: ${TAR}
	tar xvfz ${TAR}

${SRC}/httping: ${SRC}
	cd ${SRC}; CC=gclang ./configure; CC=gclang make

httping.bc: ${SRC}/httping
	get-bc -m -o httping.bc ${SRC}/httping


clean:
	rm -rf *.bc *.bca *.manifest
	rm -rf slash_specialized
	make -C ${SRC} clean


spotless: clean
	rm -rf ${SRC} ${TAR}
