GOROOT := $(shell go env GOPATH)

build:
	go build ./shared ./cmd/gclang ./cmd/gclang++ ./cmd/get-bc



install: build
	go install ./cmd/gclang ./cmd/gclang++ ./cmd/get-bc

clean:
	go clean
	rm -f gclang gclang++ get-bc

uninstall:
	rm -f $(GOROOT)/bin/gclang
	rm -f $(GOROOT)/bin/gclang++
	rm -f $(GOROOT)/bin/get-bc
