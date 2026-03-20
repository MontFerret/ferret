# AGENTS.md

This file is the canonical operating guide for coding agents working in this repository. It is written for Ferret v2 only. If you see documentation that conflicts with this file, prefer `Makefile`, `go.mod`, and `.github/workflows/build.yml`.

## Repo snapshot

- Module path: `github.com/MontFerret/ferret/v2`
- Go version: `1.23+`
- Toolchain in `go.mod`: `go1.24.5`
- This repository root is Ferret v2. Do not mix assumptions from the separate v1 branch.
- High-level flow: `Engine` -> `compiler` -> `bytecode.Program` -> `vm.VM`

## Architectural mental model

Ferret v2 is a compiled query language and runtime.

Primary pipeline:
`source -> parser -> diagnostics/AST -> compiler -> bytecode.Program -> vm.VM -> runtime values/results`

Subsystem responsibilities:
- `pkg/parser` handles syntax, parse-tree processing, and parser diagnostics.
- `pkg/compiler` performs lowering, semantic checks, bytecode emission, and optimization.
- `pkg/bytecode` defines executable program structures.
- `pkg/vm` executes bytecode programs.
- `pkg/runtime` defines value semantics and runtime-facing contracts.
- `pkg/stdlib` provides built-in host modules and functions.
- Top-level package exposes the embedding surface used by applications.

Agents should reason about changes by pipeline stage:
- Syntax changes usually begin in grammar/parser and continue into compiler lowering.
- Semantic/runtime changes usually live in compiler, runtime, or VM.
- Embedding/API changes usually affect the top-level package and integration boundaries.

## Canonical invariants

- Ferret v2 uses a register-based VM.
- `runtime.Value` is the common runtime/VM value abstraction.
- Parser-generated code is derived output, not the source of truth.
- Compiler changes must preserve program semantics expected by the VM.
- Optimizations must preserve correctness before performance.
- Runtime execution errors and internal invariant violations are different classes of failure and should not be collapsed conceptually.
- Do not assume behavior from old design notes or the v1 codebase unless it is reflected in the current v2 code.

## Package map

Agents should begin with the package whose responsibility owns the requested behavior. Do not infer ownership from file names alone when a package in this map already describes the intended boundary.

### Core execution pipeline

- `pkg/parser`
    - Grammar, parser pipeline, parse-tree processing, and parser-side diagnostics for FQL source code.
    - Do not hand-edit generated parser artifacts; edit grammar sources and regenerate.

- `pkg/compiler`
    - Lowers parsed FQL into `bytecode.Program`, performs semantic analysis, and runs optimization/code generation passes.

- `pkg/bytecode`
    - Core executable program model: instructions, operands, programs, and related structures consumed by the VM and produced by the compiler.
    - Changes here are cross-cutting and usually require corresponding compiler and VM updates.

- `pkg/runtime`
    - Core runtime semantics: values, function and module contracts, execution-facing interfaces, and shared language/runtime behavior.

- `pkg/vm`
    - Bytecode execution engine for Ferret v2, including program execution, runtime coordination, cleanup, and VM-facing result handling.
    - Performance-sensitive and semantics-sensitive; verify compiler/runtime assumptions before changing internals.

### Language and developer tooling

- `pkg/asm`
    - Assembly-layer support for Ferret bytecode programs, including parsing and encoding of assembly representations used for low-level tooling and debugging.

- `pkg/diagnostics`
    - Shared diagnostic primitives and formatting support for errors, warnings, spans, labels, and user-facing parser/compiler messages.

- `pkg/encoding`
    - Output encoding and materialization infrastructure for turning runtime values into external representations.

- `pkg/file`
    - Ferret file abstraction and loaders for bringing query sources into the parser/compiler pipeline, including path-aware source handling used by embedding and tooling flows.
    - Prefer this package when the behavior involves source origin, file-backed input, or file identity rather than general OS-level file utilities.

- `pkg/formatter`
    - Source formatting and pretty-printing support for Ferret code and related textual representations.

### Integration and extension surfaces

- `pkg/sdk`
    - Extension-oriented developer surface for building on top of Ferret internals, including helpers and contracts intended for external integrations, tools, and custom runtime/module implementations.
    - Prefer `pkg/sdk` when the goal is to support consumers building with Ferret rather than changing the core execution pipeline itself.

- `pkg/stdlib`
    - Built-in Ferret modules, namespaces, and host functions registered as the standard library.

## Primary surfaces

- Top-level package exposes the embedding API:
    - `Engine` compiles and runs FQL sources.
    - `Plan` wraps compiled bytecode plus environment state.
    - `Session` executes a plan.
    - `Module` and `ModuleRegistry` are the extension points for host modules.
- `pkg/asm` provides assembly-oriented support for working with Ferret bytecode at a lower level.
- `pkg/bytecode` defines instructions, operands, programs, and related executable structures.
- `pkg/compiler` lowers parsed FQL into bytecode and runs optimization passes.
- `pkg/diagnostics` provides shared diagnostics infrastructure used across parsing, compilation, and formatting of user-facing errors.
- `pkg/encoding` handles output encoding and materialization of runtime values.
- `pkg/file` provides file-backed source abstractions and related loading support used by parsing, compilation, and tooling flows.
- `pkg/formatter` provides source formatting and pretty-printing for Ferret code.
- `pkg/parser` parses FQL and assembles diagnostics.
- `pkg/parser/antlr` contains the grammar sources.
- `pkg/parser/fql` contains generated parser and lexer code.
- `pkg/runtime` defines runtime values, function registries, and core semantics.
- `pkg/sdk` contains extension-facing helpers and contracts for developers building on top of Ferret.
- `pkg/stdlib` registers the built-in namespaces and functions.
- `pkg/vm` executes bytecode programs.
- `test/integration/compiler`, `test/integration/optimization`, and `test/integration/vm` are the main regression suites.
- `test/e2e` covers CLI and browser-backed flows.

## Where to start by task

- Add or change syntax:
    - edit grammar under `pkg/parser/antlr`
    - regenerate parser artifacts
    - inspect parser diagnostics/code in `pkg/parser`
    - update compiler lowering in `pkg/compiler`
    - add or update integration coverage

- Add or change bytecode/opcodes:
    - inspect `pkg/bytecode`
    - inspect compiler emission sites
    - inspect VM execution in `pkg/vm`
    - validate with VM/integration tests

- Change runtime value semantics:
    - inspect `pkg/runtime`
    - inspect any relevant assumptions in `pkg/vm`
    - inspect result/materialization behavior if affected

- Change diagnostics behavior:
    - inspect `pkg/diagnostics`
    - inspect parser/compiler call sites that construct or transform diagnostics
    - validate both message content and span/label accuracy

- Change output/materialization behavior:
    - inspect `pkg/encoding`
    - inspect runtime and VM call sites that feed values into encoders/materializers
    - validate public-facing behavior and resource/cleanup interactions if relevant

- Change formatting behavior:
    - inspect `pkg/formatter`
    - validate formatting stability with targeted tests or fixtures

- Change file-backed source handling or source loading behavior:
    - inspect `pkg/file`
    - inspect parser/compiler call sites that consume source objects
    - validate path-aware diagnostics and any embedding/tooling behavior that depends on source identity

- Change embedding API:
    - inspect top-level package (`Engine`, `Plan`, `Session`)
    - inspect downstream compiler/runtime/VM interactions
    - validate public behavior with integration or e2e coverage as appropriate

- Change built-in functions/modules:
    - inspect `pkg/stdlib`
    - inspect host function/module registration and runtime contracts

- Change extension or integration support for external developers/tools:
    - inspect `pkg/sdk`
    - validate that the change belongs to an extension-facing surface rather than the core pipeline

- Change developer tooling or low-level program tooling:
    - inspect `pkg/asm`, `pkg/formatter`, `pkg/diagnostics`, or `pkg/sdk` depending on which surface owns the behavior

## Stability guide

Treat these as relatively stable unless the task explicitly targets them:
- the overall pipeline shape: parser -> compiler -> bytecode -> VM
- the parser generation workflow
- the top-level embedding entry points

Treat these as implementation-sensitive and verify current code before proposing changes:
- optimizer internals
- diagnostics plumbing
- VM execution internals
- runtime value behavior and cleanup/resource semantics
- encoding/materialization behavior

Do not treat historical discussion, stale comments, or old branches as authoritative.

## Go type and file structure rules

These rules are mandatory unless the task explicitly requires otherwise.

## Go type and file structure rules

These rules are mandatory unless the task explicitly requires otherwise.

- Do not define multiple method-bearing structs in the same `.go` file.
- Prefer declaring a method-bearing struct as a standalone `type Name struct { ... }`.
- A method-bearing struct should usually live in its own file, named after the primary type or responsibility whenever practical, for example:
    - `result.go` for `Result`
    - `register_file.go` for `RegisterFile`
    - `call_stack.go` for `CallStack`
- Grouped `type ( ... )` declarations are allowed for interfaces, passive data-only structs, and other small related helper/value types that belong to the same narrow concern.
- A grouped `type ( ... )` block may also contain exactly one method-bearing struct when:
    - it is the only behavioral type in the file, and
    - the other grouped types are passive helper/value types from the same narrow concern.
- Do not use grouped `type ( ... )` declarations to hide multiple substantial behavioral types.
- If a helper struct later gains methods and would create more than one method-bearing struct in the file, extract it into its own file immediately.
- Methods for a struct should live in the same file as the struct unless there is a strong, explicit reason to split by concern.
- Do not place a new method-bearing struct into an existing file just because the code compiles.

Allowed:

```go
type (
	PassResult struct {
		Metadata map[string]any
		Modified bool
	}

	PassContext struct {
		Program  *bytecode.Program
		CFG      *ControlFlowGraph
		Metadata map[string]any
	}

	Pass interface {
		Name() string
		Requires() []string
		Run(ctx *PassContext) (*PassResult, error)
	}
)
```

Avoid:

```go
type (
	Result struct {
		// ...
	}

	execState struct {
		// ...
	}
)
```

### Rationale:
- one method-bearing type per file keeps ownership of behavior obvious
- standalone method-bearing types make diffs and reviews clearer
- grouped type blocks are fine for passive, closely related types, but should not hide substantial behavioral types

Comment rules for functions and methods
- Do not add comments to every function or method by default.
- Exported functions and methods should usually have doc comments, especially in public, embedding-facing, or extension-facing packages.
- Unexported functions and methods should be commented only when they carry non-obvious behavior, invariants, side effects, ownership rules, cleanup expectations, or protocol/lifecycle constraints.
- Comments must explain intent, contract, invariants, side effects, or lifecycle behavior.
- Prefer comments that explain why the code exists, what must remain true, or how the method is meant to be used.
- Do not write comments that merely restate the method name or signature.
- For VM, runtime, compiler, encoding, and diagnostics internals, prefer comments on semantics and invariants over implementation narration.
- Avoid comment wallpaper. Dense, meaningful comments are preferred over mechanically documenting obvious code.

Preferred:

```go
// Close releases resources associated with the result.
// It is safe to call multiple times. Once closed, the result must not be reused.
func (r *Result) Close() error
```

Preferred for internal code:

```go
// promoteEscaped ensures a value that may outlive the current register write
// is no longer tied to the current ownership path.
func (s *execState) promoteEscaped(...)
```

Avoid:

```go
// Close closes the result.
func (r *Result) Close() error
```

## Development practice expectations

Agents must follow repository-specific engineering discipline rather than generic style preferences.

### Core principles
- Preserve correctness first.
- Preserve subsystem boundaries and invariants.
- Prefer the smallest local change that fully solves the task.
- Avoid introducing abstractions, indirection, or refactors unless they are necessary for correctness, maintainability, or an explicitly requested design change.
- Do not optimize by intuition alone; use measurements for performance-sensitive work.
- Keep behavioral ownership obvious in code structure, naming, and file layout.

### Mandatory expectations
- Identify the owning subsystem before making a non-trivial change.
- Preserve existing behavior unless the task explicitly requires changing it.
- Add or update tests for any behavior change.
- Add or update benchmarks for any significant change.
- Run the narrowest relevant validation first, then broaden as appropriate.
- Do not claim tests, benchmarks, or validation were completed unless they were actually run.
- Do not treat historical discussions, abandoned directions, or old branches as authoritative over current code and repository guidance.
- Do not perform opportunistic refactors unrelated to the requested task unless they are required for correctness.

### Required workflow for non-trivial changes

Before making a non-trivial change, agents must:
1.	Identify the owning subsystem.
2.	Identify the contract, invariant, or behavior being preserved or changed.
3.	Choose the smallest reasonable implementation that fits the existing design.
4.	Determine whether the change is significant.
5.	Add or update correctness tests.
6.	Add or update benchmarks if the change is significant.
7.	Run relevant validation and summarize the results accurately.

### Significant changes

A change is significant when it could reasonably affect:
- execution throughput
- compile-time performance
- latency on common paths
- allocation patterns
- memory reuse, pooling, or cleanup behavior
- result/materialization cost
- optimizer or code generation output relevant to performance

This includes, but is not limited to, changes in:
- pkg/vm
- pkg/runtime
- pkg/compiler
- pkg/bytecode
- pkg/encoding
- parser/compiler hot paths
- caching, pooling, register allocation, ownership tracking, or materialization logic

This usually does not include:
- comment-only, docs-only, or formatting-only edits
- pure renames with no behavior change
- test-only changes
- narrowly scoped refactors that do not affect behavior or hot paths

When in doubt, treat the change as significant and benchmark it.

### Benchmark workflow for significant changes

For significant changes, agents must:
- run relevant benchmarks before making the change and save the results as a baseline
- implement the change
- run the same benchmarks again after the change
- compare before/after results, preferably including ns/op, B/op, and allocs/op
- report the benchmark command used and summarize the performance delta

If no relevant benchmark exists for the changed hot path, add one.

If benchmark tooling or environment is unavailable, state that explicitly and do not claim benchmark validation was completed.

## Validation and evidence

When finishing a non-trivial change, agents should report:
- owning subsystem
- files changed
- tests added or updated
- benchmarks added or updated
- validation commands run
- benchmark commands run, if applicable
- notable invariants preserved or intentionally changed

For significant changes:
- tests alone are not sufficient
- both correctness tests and benchmarks are required
- benchmark results must be compared against a baseline when the environment allows it

### Change discipline
- Prefer adapting an existing local pattern over introducing a new architectural pattern.
- Do not add new helper layers, wrappers, interfaces, or abstractions only for aesthetic reasons.
- Do not move code across packages unless the ownership boundary is genuinely wrong.
- Keep diffs focused on the requested task.
- If a cleanup is necessary to make the requested change safe, keep it tightly scoped and explain why it was needed.

### Comment and documentation discipline
- Add comments where semantics, invariants, side effects, ownership, lifecycle, or recovery behavior are non-obvious.
- Do not add comment wallpaper.
- Prefer comments that explain why, contract, or invariants rather than implementation narration.
- Public and extension-facing behavior should be documented more carefully than local obvious helpers.

### Decision bias when uncertain

When uncertain:
- preserve existing behavior
- prefer the smaller local change
- add a focused test
- treat the change as significant if performance might be affected
- verify ownership before introducing a new abstraction or package-level dependency

## Tooling prerequisites
- Go must be installed.
- make is optional but is the preferred entrypoint for repo-defined workflows.
- Java plus ANTLR 4.13.2 are required when regenerating parser artifacts.
- lab plus a reachable Chromium instance are required for e2e coverage.
- staticcheck, goimports, and revive are needed for lint/format flows; install them with make install-tools.

## Command matrix
- Broad validation: go test ./...
- Race-heavy package and integration coverage: make test
- Lint: make lint
- Format: make fmt
- Regenerate parser/codegen artifacts: make generate
- Run this only when grammar or generator inputs change.
- Build the CLI binary: make compile
- Run e2e coverage: LAB_BIN=/absolute/path/to/lab make e2e
- Ensure Chromium is reachable at http://127.0.0.1:9222/json/version.
- CI uses docker run -d -p 9222:9222 ghcr.io/montferret/chromium:92.0.4512.0.

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

### Validation expectations
- After code changes, run the narrowest tests that prove the behavior you touched.
- Before finishing broader changes, run the relevant repo-level command from the matrix above.
- If you changed formatting-sensitive files, run make fmt.
- If you changed lint-sensitive code paths or public behavior, run make lint when the toolchain is available.
- If you changed parser grammar, generated lexer/parser output must be included and reviewed.

### Expectations for non-trivial changes

When proposing or implementing non-trivial changes:
- identify the owning subsystem first
- preserve invariants unless the task explicitly changes them
- prefer local, comprehensible changes before introducing new abstractions
- distinguish correctness work from performance work
- do not perform opportunistic refactors unrelated to the requested task unless they are necessary for correctness

## Secondary references
- README.md for product context and links to the broader Ferret ecosystem.
- CONTRIBUTING.md for human contributor process.
- .github/workflows/build.yml for the current CI validation path.