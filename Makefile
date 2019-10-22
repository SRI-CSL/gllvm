
install:
	go install -race github.com/SRI-CSL/gllvm/cmd/...


race: install
	 go test -v -race -timeout 24h ./tests


clean:
	rm -f tests/.helloworld.c.o tests/.helloworld.c.o.bc
