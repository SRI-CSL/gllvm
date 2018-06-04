OS   =  $(shell uname)
ARCH =  $(shell uname -p)

LLVMCC=gclang
LLVMGET=get-bc

TECLA=libtecla
TECLA_TARBALL=${TECLA}.tar.gz

ifeq (Darwin, $(findstring Darwin, ${OS}))

LIBEXT=dylib

else

LIBEXT=so

endif


TECLA_BITCODE=enhance.bc libtecla.a.bc libtecla.${LIBEXT}.bc

all: ${TECLA_BITCODE}

${TECLA}: ${TECLA_TARBALL}
	tar xvfz ${TECLA_TARBALL}

${TECLA}/enhance: ${TECLA}
	cd ${TECLA}; CC=${LLVMCC} make

enhance.bc: ${TECLA}/enhance
	cd ${TECLA}; ${LLVMGET} enhance; mv enhance.bc ../


libtecla.a.bc: ${TECLA}/enhance
	cd ${TECLA}; ${LLVMGET} -b libtecla.a; mv libtecla.a.bc ../

libtecla.${LIBEXT}.bc: ${TECLA}/enhance
	cd ${TECLA}; ${LLVMGET} libtecla.${LIBEXT}; mv libtecla*${LIBEXT}*.bc ../libtecla.${LIBEXT}.bc



very_clean:
	rm -rf ${TECLA} ${TECLA_BITCODE}

spotless: very_clean
	rm -rf ${TECLA_TARBALL}
