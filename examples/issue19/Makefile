



one:
	${CC} test.c -fsanitize=address -o test


test.o:
	${CC} -c test.c -fsanitize=address


test:  test.o
	${CC} test.o -fsanitize=address  -o test



clean:
	rm -f test.o test .test* *~

