# AGENT.md

This file is the canonical operating guide for coding agents working in this repository. It is written for Ferret v2 only. If you see documentation that conflicts with this file, prefer `Makefile`, `go.mod`, and `.github/workflows/build.yml`.

## Repo snapshot

- Module path: `github.com/MontFerret/ferret/v2`
- Go version: `1.23+`
- Toolchain in `go.mod`: `go1.24.5`
- This repository root is Ferret v2. Do not mix assumptions from the separate v1 branch.
- High-level flow: `Engine` -> `compiler` -> `bytecode.Program` -> `vm.VM`

## Primary surfaces

- Top-level package exposes the embedding API:
  - `Engine` compiles and runs FQL sources.
  - `Plan` wraps compiled bytecode plus environment state.
  - `Session` executes a plan.
  - `Module` and `ModuleRegistry` are the extension points for host modules.
- `pkg/parser` parses FQL and assembles diagnostics.
- `pkg/parser/antlr` contains the grammar sources.
- `pkg/parser/fql` contains generated parser and lexer code.
- `pkg/compiler` lowers parsed FQL into bytecode and runs optimization passes.
- `pkg/vm` executes bytecode programs.
- `pkg/runtime` defines runtime values, function registries, and core semantics.
- `pkg/stdlib` registers the built-in namespaces and functions.
- `test/integration/compiler`, `test/integration/optimization`, and `test/integration/vm` are the main regression suites.
- `test/e2e` covers CLI and browser-backed flows.

## Tooling prerequisites

- Go must be installed.
- `make` is optional but is the preferred entrypoint for repo-defined workflows.
- Java plus ANTLR `4.13.2` are required when regenerating parser artifacts.
- `lab` plus a reachable Chromium instance are required for e2e coverage.
- `staticcheck`, `goimports`, and `revive` are needed for lint/format flows; install them with `make install-tools`.

## Command matrix

- Broad validation: `go test ./...`
- Race-heavy package and integration coverage: `make test`
- Lint: `make lint`
- Format: `make fmt`
- Regenerate parser/codegen artifacts: `make generate`
  - Run this only when grammar or generator inputs change.
- Build the CLI binary: `make compile`
- Run e2e coverage: `LAB_BIN=/absolute/path/to/lab make e2e`
  - Ensure Chromium is reachable at `http://127.0.0.1:9222/json/version`.
  - CI uses `docker run -d -p 9222:9222 ghcr.io/montferret/chromium:92.0.4512.0`.

## Editing rules

- Never hand-edit generated files under `pkg/parser/fql` or `pkg/parser/antlr/gen`.
- Parser generation is driven by `pkg/parser/parser.go`:
  - `antlr -Xexact-output-dir -o fql -package fql -visitor -Dlanguage=Go antlr/FqlLexer.g4 antlr/FqlParser.g4`
  - `go run ./tools/patch_lexer.go`
- If you change grammar files in `pkg/parser/antlr`, run `make generate` and commit the generated output in the same change.
- Treat `Makefile` and `.github/workflows/build.yml` as the source of truth for validation commands.
- Prefer narrow validation first, then broaden:
  - Package-local changes: run the affected `go test` package(s).
  - Compiler, optimizer, or VM changes: run the relevant integration suite(s).
  - Cross-cutting changes: finish with `go test ./...` or `make test`.
- Do not assume e2e is available locally. If Chromium or `lab` is missing, state that explicitly.

## Validation expectations

- After code changes, run the narrowest tests that prove the behavior you touched.
- Before finishing broader changes, run the relevant repo-level command from the matrix above.
- If you changed formatting-sensitive files, run `make fmt`.
- If you changed lint-sensitive code paths or public behavior, run `make lint` when the toolchain is available.
- If you changed parser grammar, generated lexer/parser output must be included and reviewed.

## Secondary references

- `README.md` for product context and links to the broader Ferret ecosystem.
- `CONTRIBUTING.md` for human contributor process.
- `.github/workflows/build.yml` for the current CI validation path.
