
install:
	go install -race github.com/SRI-CSL/gllvm/cmd/...


race: install
	 go test -v -race ./tests


clean:
	rm -f tests/.helloworld.c.o tests/.helloworld.c.o.bc
