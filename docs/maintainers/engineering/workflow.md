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
| `make compile-32bit` | Cross-compile all root-module production and test packages for Linux/386 with CGO disabled; test binaries are not executed. |
| `make test` | Run unit/race, integration/race, and security suites. |
| `make test-32bit` | Run all root-module and independent API-tool-module tests on Linux/386 with CGO disabled and without the race detector. |
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
Never hand-edit generated parser artifacts or the token vocabulary; change the
authoritative grammar or generator inputs instead.

Syntax changes normally require coordinated grammar, generated parser,
diagnostic/parser integration, compiler lowering, formatter, and integration
coverage. Review all generated changes and commit them with their source change.

## Test layout

Package tests live beside their code. Native embedding lifecycle and composition
tests live in `pkg/engine`; component tests live in `pkg/engine/internal/bootstrap`,
`pkg/engine/internal/host`, `pkg/engine/internal/resource`, and
`pkg/engine/internal/session`. Bootstrap tests cover construction ownership
transfer, rollback ordering and errors, and initialization hook snapshots. Session
component tests cover acquisition failure and panic rollback, transfer to ordinary
execution or debugger services, exact permit/VM release, concurrent cleanup,
and materialization error separation. Native tests retain admission, cancellation,
option, hook, and output contracts without accessing execution internals. These
packages participate in the existing `pkg/...` unit/race, coverage, static
analysis, and formatting targets. Root tests guard the curated façade's exported
declarations, type aliases, function signatures, constants, and public embedding behavior.
Universal API adapter tests, external-package examples, and benchmarks live in
root-level `uapi`. It is explicitly included in unit/race tests, coverage,
unit benchmarks, static analysis, and import formatting alongside `pkg/...`.
Additional suites are grouped under:

* `test/integration/compiler`: language compilation and semantic behavior;
* `test/integration/optimization`: optimizer equivalence and lowering behavior;
* `test/integration/vm`: cross-layer execution behavior;
* `test/integration/formatter`: FQL formatter behavior;
* `test/security`: security-focused regression coverage;
* `test/spec`: shared test helpers and specification fixtures;
* `test/benchmarks`: cross-layer integration benchmarks.

Use the [testing and validation policy](testing.md) to choose coverage and checks.
Race-enabled Make targets require CGO for the race detector.

`go test ./...` is a useful broad, non-race check for the root Go module. It does
not cover the independent tool modules or reproduce the full race and security
composition of `make test`.

## Linux/386 validation

The main CI workflow has a separate Linux/386 matrix matching the Go versions
in the 64-bit build matrix. Each job runs `make compile-32bit`, followed by
`make test-32bit`. Matrix fail-fast is disabled so a failure on one Go version
does not cancel the others.

The compile target works from other supported Go build hosts. It uses
`go test -run '^$' -exec=true ./...` to compile root-module test packages without
executing their binaries. To execute the full test suite on a host that can run
Linux/386 binaries, use:

```sh
make test-32bit
```

The test target sets `GOOS=linux GOARCH=386 CGO_ENABLED=0` and runs
`go test -count=1 ./...` in the root module and both independent API tooling
modules, `tools/apiref` and `tools/apipublish`. This includes unit, integration,
security, compatibility, script, and tool tests. Test caching is disabled;
benchmarks and extended fuzzing remain separate.

The Linux/386 suite runs without the race detector. Existing 64-bit build,
coverage, security, and race jobs retain their coverage.

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

Baseline and comparison requirements live in the
[benchmark policy](testing.md#significant-changes-and-benchmarks). The benchmark
workflow records main baselines on `gh-pages`, runs unit comparisons
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

* [Engineering principles](principles.md)
* [Go code style](code-style.md)
* [Testing and performance](testing.md)
* [Architecture](../architecture/overview.md)
* [Runtime and lifecycle](../architecture/runtime.md)
* [Debugger architecture](../architecture/debugger.md)
* [Modules, SDK, and standard library](../architecture/modules.md)
* [Release automation](../release/process.md)
