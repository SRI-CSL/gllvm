GOROOT := $(shell go env GOPATH)

build:
	go build ./cmd/...


install: build
	go install ./cmd/...

clean:
	go clean

fmt:
	gofmt -s -w shared/*.go cmd/*/*.go

uninstall:
	rm -f $(GOROOT)/bin/gclang
	rm -f $(GOROOT)/bin/gclang++
	rm -f $(GOROOT)/bin/get-bc
	rm -f $(GOROOT)/bin/gsanity-check
