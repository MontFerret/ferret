# Development Workflow

This guide maps the repository's build, generation, test, lint, and benchmark
entry points. It describes the current workflows; `go.mod`, `Makefile`, and the
GitHub Actions definitions remain authoritative when versions or commands
change.

## Authorities and prerequisites

* `go.mod` declares the minimum Go version and root module dependencies.
* `Makefile` defines local commands and pins the auxiliary Go tools used by its
  targets.
* `.github/workflows/build.yml` defines the main CI matrix and validation path.
* `.github/workflows/benchmark.yml` defines automated benchmark collection and
  comparison.
* `.github/workflows/codeql.yml` defines CodeQL setup and its generation path.

Go is required. Make is the preferred entry point for repository-composed
workflows. Lint and format targets require the versions of `staticcheck`,
`fieldalignment`, `goimports`, `benchstat`, and `revive` installed by
`make install-tools`.

Parser regeneration additionally requires Java and an `antlr` wrapper backed by
ANTLR 4.13.2, matching the generated file headers and CI setup. Consult the
workflow files rather than duplicating their installation scripts elsewhere.

The API reference and publisher tools are independent Go modules under
`tools/apiref` and `tools/apipublish`. Root targets that cover them use `go -C`
and their own module files.

## Make targets

| Target | Purpose |
| --- | --- |
| `make build` | Run lint, generation/formatting, all tests, and compile the test CLI harness. |
| `make install-tools` | Install the exact auxiliary tool versions selected by the Makefile. |
| `make compile` | Build `test/cli.go` as `bin/ferret`; this is a repository test harness, not the separate MontFerret CLI product. |
| `make test` | Run unit/race, integration/race, and security suites. |
| `make test-unit` | Run race-enabled package, script, tool-module, compatibility, and root tests. |
| `make test-integration` | Run race-enabled tests under `test/integration`. |
| `make test-security` | Run `test/security` without the race flag. |
| `make cover` | Run package coverage and submit through the configured Codecov script. |
| `make lint` | Run `staticcheck` and `revive` with repository exclusions. |
| `make fmt` | Apply field alignment, Go formatting, and import formatting across configured roots. |
| `make generate` | Run package generation and then `make fmt`. |
| `make bench-unit` | Run package benchmarks with configurable filters and counts. |
| `make bench-integration` | Run benchmarks under `test/benchmarks`. |
| `make bench` | Run both benchmark groups. |

Several targets are intentionally mutating: `make fmt` rewrites Go files,
`make generate` rewrites generated artifacts and then formats, `make clean`
removes local build and coverage output, and `make cover` performs an external
upload. Choose targets according to the task rather than running the broadest
target automatically.

## Parser generation

The `go:generate` directives in `pkg/parser/parser.go` run:

```text
antlr -Xexact-output-dir -o fql -package fql -visitor -Dlanguage=Go antlr/FqlLexer.g4 antlr/FqlParser.g4
go run ./tools/patch_lexer.go
```

The grammar sources are `pkg/parser/antlr/FqlLexer.g4` and
`pkg/parser/antlr/FqlParser.g4`. Generated output includes the vocabulary at
`pkg/parser/antlr/FqlLexer.tokens` and the lexer, parser, listener, visitor,
token, and interpreter artifacts under `pkg/parser/fql`.

Run `make generate` only after changing grammar or generator inputs. Review both
the generated artifacts and any formatting changes produced by the target.
Generated Go files and token artifacts should never receive independent manual
fixes.

## Test layout

Package tests live beside their code. The root package tests the embedding
lifecycle and cross-package composition. Additional suites are grouped under:

* `test/integration/compiler`: language compilation and semantic behavior;
* `test/integration/optimization`: optimizer equivalence and lowering behavior;
* `test/integration/vm`: cross-layer execution behavior;
* `test/integration/formatter`: FQL formatter behavior;
* `test/security`: security-focused regression coverage;
* `test/spec`: shared test helpers and specification fixtures;
* `test/benchmarks`: cross-layer integration benchmarks.

Start with the package or focused integration suite that proves the change.
Broaden to the corresponding Make target when the impact crosses packages or
matches CI coverage. Race-enabled Make targets require CGO for the race detector.

`go test ./...` is a useful broad, non-race check for the root Go module. It does
not cover the independent tool modules or reproduce the full race and security
composition of `make test`.

## Benchmarks

The benchmark targets accept these Make variables:

* `BENCH_RUN` selects tests to run alongside benchmarks and defaults to none.
* `BENCH_FILTER` selects benchmark names and defaults to all.
* `BENCH_COUNT` controls repetitions.
* `BENCH_TIMEOUT` controls the Go test timeout.

Examples:

```sh
make bench-unit BENCH_FILTER='Comparison' BENCH_COUNT=5
make bench-integration BENCH_FILTER='Compiler' BENCH_COUNT=5 BENCH_TIMEOUT=15m
```

Use identical commands and environment for before/after comparisons. The
benchmark workflow records main baselines on `gh-pages`, runs unit comparisons
for pull requests, and gates integration comparisons behind its configured
label or manual inputs. The workflow file owns the current labels, thresholds,
and artifact names.

## Documentation-only changes

The repository currently has no Markdown formatter or link-check target. For a
documentation-only change, verify relative links and referenced paths directly,
run `git diff --check`, and inspect the complete diff. Do not run unrelated Go
build or integration suites unless the documentation change also touches code,
generated artifacts, or command behavior.

## Related guides

* [Architecture](architecture.md)
* [Runtime and lifecycle](runtime.md)
* [Debugger architecture](debugger.md)
* [Modules, SDK, and standard library](modules.md)
* [Release automation](release.md)
