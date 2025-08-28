.PHONY: build install compile test e2e doc fmt lint vet release
export CGO_ENABLED=0

LAB_BIN ?= lab
DIR_BIN = ./bin
DIR_PKG = ./pkg
DIR_INTEG = ./test/integration
DIR_E2E = ./test/e2e

default: build

build: vet generate test compile

install-tools:
	go install honnef.co/go/tools/cmd/staticcheck@latest && \
	go install golang.org/x/tools/cmd/goimports@latest && \
	go install github.com/mgechev/revive@latest

install:
	go get

compile:
	go build -v -o ${DIR_BIN}/ferret \
	${DIR_E2E}/cli.go

test:
	go test ${DIR_PKG}/... && go test ${DIR_INTEG}/...

cover:
	go test -coverprofile=coverage.txt -covermode=atomic ${DIR_PKG}/... && \
	curl -s https://codecov.io/bash | bash

e2e:
	${LAB_BIN} --timeout=120 --attempts=5 --concurrency=1 --wait=http://127.0.0.1:9222/json/version --runtime=bin://./bin/ferret --files=./test/e2e/tests --cdn=./test/e2e/pages/dynamic --cdn=./test/e2e/pages/static

bench:
	go test -run=XXX -bench=. ${DIR_PKG}/...

generate:
	go generate ${DIR_PKG}/...

doc:
	godoc -http=:6060 -index

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ${DIR_PKG}/... ${DIR_INTEG}/... ${DIR_E2E}/... && \
	goimports -w -local github.com/MontFerret ${DIR_PKG} ${DIR_INTEG} ${DIR_E2E}

# https://github.com/mgechev/revive
# go get github.com/mgechev/revive
lint:
	staticcheck -tests=false ./pkg/... && \
	echo 'revive -config revive.toml -formatter stylish -exclude ./pkg/parser/fql/... -exclude ./vendor/... -exclude ./*_test.go ./...'
