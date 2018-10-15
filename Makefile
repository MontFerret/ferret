.PHONY: build compile install test e2e doc fmt lint vet release

export GOPATH

VERSION ?= $(shell git describe --tags --always --dirty)
DIR_BIN = ./bin
DIR_PKG = ./pkg
DIR_CLI = ./cli
DIR_E2E = ./e2e

default: build

build: install vet generate test compile

compile:
	go build -v -o ${DIR_BIN}/ferret \
	-ldflags "-X main.version=${VERSION}" \
	./main.go

install:
	dep ensure

test:
	go test ${DIR_PKG}/...

e2e:
	go run ${DIR_E2E}/main.go --tests ${DIR_E2E}/tests --pages ${DIR_E2E}/pages

generate:
	go generate ${DIR_PKG}/...

doc:
	godoc -http=:6060 -index

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ${DIR_CLI}/... ${DIR_PKG}/...

# https://github.com/mgechev/revive
# go get github.com/mgechev/revive
lint:
	revive -config revive.toml -formatter friendly -exclude ./pkg/parser/fql/... -exclude ./vendor/... ./...

# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	go vet ${DIR_CLI}/... ${DIR_PKG}/...

release:
	goreleaser