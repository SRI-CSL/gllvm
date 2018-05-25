
all: test

#one:
#	${CC} test.c -dead_strip -fsanitize=address -o test


test.o:
	${CC} -c test.c -fsanitize=address


test:  test.o
	${CC} test.o -dead_strip  -fsanitize=address  -o test



clean:
	rm -f test.o test .test* *~
