.PHONY: build install test doc fmt lint vet

export GOPATH

VERSION ?= $(shell git describe --tags --always --dirty)
DIR_BIN = ./bin
DIR_PKG = ./pkg
DIR_CLI = ./cli

default: build

build: install vet generate test compile

compile:
	go build -v -o ${DIR_BIN}/ferret \
	-ldflags "-X main.Version=${VERSION}" \
	./main.go

install:
	dep ensure

test:
	go test ${DIR_PKG}/...

generate:
	go generate ${DIR_PKG}/...

doc:
	godoc -http=:6060 -index

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ${DIR_CLI}/... ${DIR_PKG}/...

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	golint ${DIR_CLI}/... ${DIR_PKG}/...

# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	go vet ${DIR_CLI}/... ${DIR_PKG}/...