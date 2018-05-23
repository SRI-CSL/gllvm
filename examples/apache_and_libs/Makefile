INSTALL=${PWD}/apache-install

#gllvm
LLVMCC=gclang
LLVMCXX=gclang++
LLVMGET=get-bc

#wllvm
#LLVMCC=wllvm
#LLVMCXX=wllvm++
#LLVMGET=extract-bc


#current versions as of 5/16/2018
PCRE=pcre-8.42
APR=apr-1.6.3
APR_UTIL=apr-util-1.6.1
HTTPD=httpd-2.4.33

#the urls (if the versions change you should check these)
PCRE_URL=https://ftp.pcre.org/pub/pcre/${PCRE}.tar.gz
APR_URL=http://apache.mirrors.ionfish.org//apr/${APR}.tar.gz
APR_UTIL_URL=http://apache.mirrors.pair.com//apr/${APR_UTIL}.tar.gz
HTTPD_URL=http://apache.cs.utah.edu/httpd/${HTTPD}.tar.gz


all: build

build: ${PCRE}_build ${APR}_build ${APR_UTIL}_build ${HTTPD}_build

${PCRE}.tar.gz:
	wget ${PCRE_URL}

${PCRE}.tar: ${PCRE}.tar.gz
	gunzip -k ${PCRE}.tar.gz

${PCRE}: ${PCRE}.tar
	tar xvf ${PCRE}.tar

${PCRE}_build: ${PCRE}
	mkdir -p ${PCRE}_build                                                         && \
	cd ${PCRE}_build                                                               && \
	CXX=${LLVMCXX} CC=${LLVMCC} ../${PCRE}/configure --prefix=${INSTALL}/${PCRE}   && \
	CXX=${LLVMCXX} CC=${LLVMCC} make 	                                           && \
	cd  .libs                                                                      && \
	${LLVMGET} libpcre.so                                                          && \
	cp libpcre.so.*.bc ../../libpcre.so.bc                                         && \
	cd  ..                                                                         && \
	make install


${APR}.tar.gz:
	wget ${APR_URL}

${APR}.tar: ${APR}.tar.gz
	gunzip -k ${APR}.tar.gz

${APR}: ${APR}.tar
	tar xvf ${APR}.tar

${APR}_build: ${APR}
	mkdir -p ${APR}_build                                        && \
	cd ${APR}_build                                              && \
	CC=${LLVMCC} ../${APR}/configure --prefix=${INSTALL}/${APR}  && \
	CC=${LLVMCC} make                                            && \
	cd .libs                                                     && \
	${LLVMGET} libapr-1.so                                       && \
	cp libapr-1.so.*.bc ../../libapr-1.so.bc                     && \
	cd ..                                                        && \
	make install

${APR_UTIL}.tar.gz:
	wget ${APR_UTIL_URL}

${APR_UTIL}.tar: ${APR_UTIL}.tar.gz
	gunzip -k ${APR_UTIL}.tar.gz

${APR_UTIL}: ${APR_UTIL}.tar
	tar xvf ${APR_UTIL}.tar

${APR_UTIL}_build: ${APR_UTIL}
	mkdir -p ${APR_UTIL}_build                                             && \
	cd ${APR_UTIL}_build                                                   && \
	CC=${LLVMCC} ../${APR_UTIL}/configure --prefix=${INSTALL}/${APR_UTIL}     \
                                          --with-apr=${INSTALL}/${APR}     && \
	CC=${LLVMCC} make                                                      && \
	cd .libs                                                               && \
	${LLVMGET} libaprutil-1.so                                             && \
	cp libaprutil-1.so.*.bc ../../libaprutil-1.so.bc                       && \
	cd ..                                                                  && \
	make install

${HTTPD}.tar.gz:
	wget ${HTTPD_URL}

${HTTPD}.tar: ${HTTPD}.tar.gz
	gunzip -k ${HTTPD}.tar.gz

${HTTPD}: ${HTTPD}.tar
	tar xvf ${HTTPD}.tar

${HTTPD}_build: ${HTTPD}
	mkdir -p ${HTTPD}_build                                                           && \
	cd  ${HTTPD}_build                                                                && \
	CC=${LLVMCC} ../${HTTPD}/configure --prefix=${INSTALL}/${HTTPD}                      \
                                       --with-pcre=${INSTALL}/${PCRE}/bin/pcre-config    \
                                       --with-apr=${INSTALL}/${APR}                      \
				       --with-apr-util=${INSTALL}/${APR_UTIL}                         && \
	CC=${LLVMCC} make                                                                 && \
	${LLVMGET} httpd                                                                  && \
	cp httpd.bc ../                                                                   && \
	make install


clean:
	rm -rf *~  *_build  slash linked_slash httpd *.manifest httpd_*


very_clean: clean
	rm -rf ${HTTPD} ${APR_UTIL} ${APR} ${PCRE}

spotless: very_clean
	rm -rf *.tar.gz *.tar *.bc build ${INSTALL}
