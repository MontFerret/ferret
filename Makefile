.PHONY: build install compile test e2e doc fmt lint vet release bench
export CGO_ENABLED=0

LAB_BIN ?= lab
DIR_BIN = ./bin
DIR_PKG = ./pkg
DIR_COMPAT = ./compat
DIR_TEST = ./test
DIR_INTEG = ${DIR_TEST}/integration
DIR_BENCH = ${DIR_TEST}/benchmarks
DIR_SEC = ${DIR_TEST}/security
BENCH_RUN ?= '^$$'
BENCH_FILTER ?= .
BENCH_COUNT ?= 1
BENCH_TIMEOUT ?= 30m

default: build

build: lint generate test compile

install-tools:
	go install honnef.co/go/tools/cmd/staticcheck@latest && \
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest && \
	go install golang.org/x/tools/cmd/goimports@latest && \
	go install golang.org/x/perf/cmd/benchstat@latest && \
	go install github.com/mgechev/revive@latest

install:
	go get

compile:
	go build -v -o ${DIR_BIN}/ferret \
	${DIR_TEST}/cli.go

test: test-unit test-integration test-security

test-unit:
	CGO_ENABLED=1 go test -race ${DIR_PKG}/... && CGO_ENABLED=1 go test -race ${DIR_COMPAT}/... .

test-integration:
	CGO_ENABLED=1 go test -race ${DIR_INTEG}/...

test-security:
	go test ${DIR_SEC}/...

clean:
	rm -rf ${DIR_BIN}/* && \
	rm -rf coverage.txt && \
	go clean -testcache

cover:
	go test -coverprofile=coverage.txt -covermode=atomic ${DIR_PKG}/... && \
	curl -s https://codecov.io/bash | bash

e2e:
	${LAB_BIN} --timeout=120 --attempts=5 --concurrency=1 --wait=http://127.0.0.1:9222/json/version --runtime=bin://./bin/ferret --files=./test/e2e/tests --cdn=./test/e2e/pages/dynamic --cdn=./test/e2e/pages/static

bench-unit:
	go test ${DIR_PKG}/... -run ${BENCH_RUN} -bench ${BENCH_FILTER} -benchmem -count=${BENCH_COUNT} -timeout ${BENCH_TIMEOUT}

bench-integration:
	go test ${DIR_BENCH}/... -run ${BENCH_RUN} -bench ${BENCH_FILTER} -benchmem -count=${BENCH_COUNT} -timeout ${BENCH_TIMEOUT}

bench: bench-unit bench-integration

generate:
	go generate ${DIR_PKG}/... && \
	make fmt

doc:
	godoc -http=:6060 -index

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	fieldalignment --fix  ./... && \
	go fmt ./... && \
	goimports -w -local github.com/MontFerret ${DIR_PKG} ${DIR_INTEG} ${DIR_E2E}

# https://github.com/mgechev/revive
# go get github.com/mgechev/revive
lint:
	staticcheck -tests=false -checks=all,-U1000,-ST1000,-ST1001,-ST1020,-ST1022,-S1002 $$(go list ./pkg/... | grep -v /fql) && \
	revive -config revive.toml -formatter stylish -exclude ./pkg/parser/fql/... -exclude ./vendor/... -exclude ./*_test.go ./...
