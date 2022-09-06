all: help


help:
	@echo ''
	@echo 'Here are the targets:'
	@echo ''
	@echo 'To test                :    "make test"'
	@echo 'To develop             :    "make develop"'
	@echo 'To install             :    "make install"'
	@echo 'To format              :    "make format"'
	@echo 'To vet                 :    "make vet"'
	@echo 'To staticcheck         :    "make check"'
	@echo 'To clean               :    "make clean"'
	@echo ''


install: develop

develop:
	go install github.com/SRI-CSL/gllvm/cmd/...


test: develop
	 go test -v ./tests

format:
	gofmt -s -w shared/*.go tests/*.go cmd/*/*.go

check:
	staticcheck ./...

vet:
	go vet ./...

clean:
	rm -f data/*hello data/*.bc [td]*/.*.c.o [td]*/*.o [td]*/.*.c.o.bc data/*.notanextensionthatwerecognize
