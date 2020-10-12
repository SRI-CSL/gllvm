

all: main

main: main.o
	${CC} -flto a.o main.o -o main

main.o: main.c a.o
	${CC} -c main.c -o main.o

a.o: a.c a.h
	${CC} -flto -c a.c -o a.o

plain: main.c a.c a.h
	${CC} main.c a.c -o plain


clean:
	rm -f *.o main plain


