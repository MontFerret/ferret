.PHONY: build install compile test e2e doc fmt lint vet release

export GOPATH
export GO111MODULE=on

VERSION ?= $(shell git describe --tags --always --dirty)
RELEASE_VERSION ?= $(version)
DIR_BIN = ./bin
DIR_PKG = ./pkg
DIR_CLI = ./cli
DIR_E2E = ./e2e

default: build

build: vet generate test compile

install:
	go get

compile:
	go build -race -v -o ${DIR_BIN}/ferret \
	-ldflags "-X main.version=${VERSION}" \
	./main.go

test:
	go test -race ${DIR_PKG}/...

cover:
	go test -race -coverprofile=coverage.txt -covermode=atomic ${DIR_PKG}/... && \
	curl -s https://codecov.io/bash | bash

e2e:
	lab --timeout=120 --times=5 --concurrency=1 --wait=http://127.0.0.1:9222/json/version --runtime=bin://./bin/ferret --files=./e2e/tests --cdn=./e2e/pages/dynamic --cdn=./e2e/pages/static

bench:
	go test -run=XXX -bench=. ${DIR_PKG}/...

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
	revive -config revive.toml -formatter stylish -exclude ./pkg/parser/fql/... -exclude ./vendor/... ./...

# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	go vet ${DIR_CLI}/... ${DIR_PKG}/...

release:
ifeq ($(RELEASE_VERSION), )
	$(error "Release version is required (version=x)")
else ifeq ($(GITHUB_TOKEN), )
	$(error "GitHub token is required (GITHUB_TOKEN)")
else
	rm -rf ./dist && \
	git tag -a v$(RELEASE_VERSION) -m "New $(RELEASE_VERSION) version" && \
	git push origin v$(RELEASE_VERSION) && \
	goreleaser
endif
