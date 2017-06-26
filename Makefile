GOROOT := $(shell go env GOPATH)

build:
	go build

install: build
	go install
	ln -f -s $(GOROOT)/bin/gowllvm $(GOROOT)/bin/gowclang
	ln -f -s $(GOROOT)/bin/gowllvm $(GOROOT)/bin/gowclang++
	ln -f -s $(GOROOT)/bin/gowllvm $(GOROOT)/bin/gowextract
