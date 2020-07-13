all: help


help:
	@echo ''
	@echo 'Here are the targets:'
	@echo ''
	@echo 'To test                :    "make check"'
	@echo 'To develop             :    "make develop"'
	@echo 'To install             :    "make install"'
	@echo 'To format              :    "make format"'
	@echo 'To lint                :    "make lint"'
	@echo 'To clean               :    "make clean"'
	@echo ''



develop:
	go install github.com/SRI-CSL/gllvm/cmd/...


check: develop
	 go test -v ./tests

format:
	gofmt -s -w shared/*.go tests/*.go cmd/*/*.go

lint:
	golint ./shared/ ./tests/ ./cmd/...

clean:
	rm -f data/hello data/hello.bc [td]*/.helloworld.c.o [td]*/.helloworld.c.o.bc
