build:
	go build

install: build
	go install
	chmod +x bin/*
	cp bin/* /usr/local/bin
