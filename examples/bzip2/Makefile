all: bzip2.bc bzip2_from_bitcode test

CLANG = $(shell which clang)
ifeq ($(CLANG),)
	CLANG = gclang
endif

./bzip2:
	git clone git://sourceware.org/git/bzip2.git
# the patch is to remove the hard coded "CC=gcc". Bad Julian. No biscuit.
	patch bzip2/Makefile < patch.txt

./bzip2/bzip2: bzip2
	CC=gclang make -C bzip2

./bzip2.bc: bzip2/bzip2
	get-bc -o bzip2.bc bzip2/bzip2

# to build an bitcode archive we do not use any additional flags; the '-m' flag writes the manifest (i.e. the ingredients)
libbz2.bca:
	get-bc -m -o libbz2.bca bzip2/libbz2.a

# to build a module we just use the '-b' flag; the '-m' flag writes the manifest (i.e. the ingredients)
libbz2.a.bc:
	get-bc -m -b -o libbz2.a.bc bzip2/libbz2.a

# groked this list from the manifests
modules: bzip2/bzip2
	cp ./bzip2/.bzip2.c.o.bc      bzip2.c.o.bc
	cp ./bzip2/.blocksort.c.o.bc  blocksort.c.o.bc
	cp ./bzip2/.huffman.c.o.bc    huffman.c.o.bc
	cp ./bzip2/.crctable.c.o.bc   crctable.c.o.bc
	cp ./bzip2/.randtable.c.o.bc  randtable.c.o.bc
	cp ./bzip2/.compress.c.o.bc   compress.c.o.bc
	cp ./bzip2/.decompress.c.o.bc decompress.c.o.bc
	cp ./bzip2/.bzlib.c.o.bc      bzlib.c.o.bc


# the override flag is necessary and is there to prevent any complaints about multiple definitions.
bzip2_from_bitcode: bzip2.bc modules
	${CLANG} bzip2.c.o.bc blocksort.c.o.bc huffman.c.o.bc crctable.c.o.bc randtable.c.o.bc compress.c.o.bc decompress.c.o.bc bzlib.c.o.bc -o bzip2_from_bitcode

test: bzip2_from_bitcode
	./bzip2_from_bitcode -1  < ./bzip2/sample1.ref > sample1.rb2
	./bzip2_from_bitcode -2  < ./bzip2/sample2.ref > sample2.rb2
	./bzip2_from_bitcode -3  < ./bzip2/sample3.ref > sample3.rb2
	./bzip2_from_bitcode -d  < ./bzip2/sample1.bz2 > sample1.tst
	./bzip2_from_bitcode -d  < ./bzip2/sample2.bz2 > sample2.tst
	./bzip2_from_bitcode -ds < ./bzip2/sample3.bz2 > sample3.tst
	cmp ./bzip2/sample1.bz2 sample1.rb2
	cmp ./bzip2/sample2.bz2 sample2.rb2
	cmp ./bzip2/sample3.bz2 sample3.rb2
	cmp sample1.tst ./bzip2/sample1.ref
	cmp sample2.tst ./bzip2/sample2.ref
	cmp sample3.tst ./bzip2/sample3.ref

clean:
	rm -f *.bc *.bca *.manifest bzip2_from_bitcode *.rb2 *.tst bzip2_slashed
	rm -rf slash_specialized
	make -C bzip2 clean

spotless: clean
	rm -rf bzip2 bzip2_tests



.PHONY: clean spotless modules
