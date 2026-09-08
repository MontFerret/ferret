.PHONY: build install compile test doc fmt lint vet release bench
export CGO_ENABLED=0

DIR_BIN = ./bin
DIR_PKG = ./pkg
DIR_TOOLS = ./tools
DIR_TOOL_APIREF = ${DIR_TOOLS}/apiref
DIR_TOOL_APIPUBLISH = ${DIR_TOOLS}/apipublish
DIR_COMPAT = ./compat
DIR_SCRIPTS = ./scripts
DIR_TEST = ./test
DIR_INTEG = ${DIR_TEST}/integration
DIR_BENCH = ${DIR_TEST}/benchmarks
DIR_SEC = ${DIR_TEST}/security
BENCH_RUN ?= '^$$'
BENCH_FILTER ?= .
BENCH_COUNT ?= 1
BENCH_TIMEOUT ?= 30m
STATICCHECK_VERSION = v0.7.0
GO_TOOLS_VERSION = v0.49.0
GO_PERF_VERSION = v0.0.0-20260813145340-fd4a688df892
REVIVE_VERSION = v1.15.0
STATICCHECK_FLAGS = -tests=false -checks=all,-U1000,-ST1000,-ST1001,-ST1020,-ST1022,-S1002

default: build

build: lint generate test compile

install-tools:
	go install honnef.co/go/tools/cmd/staticcheck@${STATICCHECK_VERSION} && \
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@${GO_TOOLS_VERSION} && \
	go install golang.org/x/tools/cmd/goimports@${GO_TOOLS_VERSION} && \
	go install golang.org/x/perf/cmd/benchstat@${GO_PERF_VERSION} && \
	go install github.com/mgechev/revive@${REVIVE_VERSION}

install:
	go get

compile:
	go build -v -o ${DIR_BIN}/ferret \
	${DIR_TEST}/cli.go

test: test-unit test-integration test-security

test-unit:
	CGO_ENABLED=1 go test -race ${DIR_PKG}/... ${DIR_SCRIPTS}/... && \
	CGO_ENABLED=1 go -C ${DIR_TOOL_APIREF} test -race ./... && \
	CGO_ENABLED=1 go -C ${DIR_TOOL_APIPUBLISH} test -race ./... && \
	CGO_ENABLED=1 go test -race ${DIR_COMPAT}/... .

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
	(cd ${DIR_TOOL_APIREF} && fieldalignment --fix ./...) && \
	(cd ${DIR_TOOL_APIPUBLISH} && fieldalignment --fix ./...) && \
	go fmt ./... && \
	go -C ${DIR_TOOL_APIREF} fmt ./... && \
	go -C ${DIR_TOOL_APIPUBLISH} fmt ./... && \
	goimports -w -local github.com/MontFerret ${DIR_PKG} ${DIR_TOOLS} ${DIR_INTEG} ${DIR_E2E}

# https://github.com/mgechev/revive
# go get github.com/mgechev/revive
lint:
	staticcheck ${STATICCHECK_FLAGS} $$(go list ${DIR_PKG}/... | grep -v /fql) && \
	(cd ${DIR_TOOL_APIREF} && staticcheck ${STATICCHECK_FLAGS} $$(go list ./...)) && \
	(cd ${DIR_TOOL_APIPUBLISH} && staticcheck ${STATICCHECK_FLAGS} $$(go list ./...)) && \
	revive -config revive.toml -formatter stylish -exclude ./pkg/parser/fql/... -exclude ./vendor/... -exclude ./*_test.go ./...
