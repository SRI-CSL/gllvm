GOROOT := $(shell go env GOPATH)

build:
	go build

install: build
	go install
	ln -f -s $(GOROOT)/bin/gllvm $(GOROOT)/bin/gclang
	ln -f -s $(GOROOT)/bin/gllvm $(GOROOT)/bin/gclang++
	ln -f -s $(GOROOT)/bin/gllvm $(GOROOT)/bin/get-bc

clean:
	go clean
	rm $(GOROOT)/bin/gclang
	rm $(GOROOT)/bin/gclang++
	rm $(GOROOT)/bin/get-bc
