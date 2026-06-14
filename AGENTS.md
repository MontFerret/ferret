# AGENTS.md

This file is the canonical operating guide for coding agents working in this repository. It is written for Ferret v2 only. If repository documentation conflicts with this file, prefer `Makefile`, `go.mod`, and `.github/workflows/build.yml` for commands, toolchain, and CI behavior.

## Repo snapshot

- Module path: `github.com/MontFerret/ferret/v2`
- Go version: `1.23+`
- Toolchain in `go.mod`: `go1.24.5`
- This repository root is Ferret v2. Do not mix assumptions from the separate v1 branch.
- High-level flow: `Engine` -> `compiler` -> `bytecode.Program` -> `vm.VM`

## Architectural mental model

Ferret v2 is a compiled query language and runtime.

Primary pipeline:

```text
source -> parser -> diagnostics/AST -> compiler -> bytecode.Program -> vm.VM -> runtime values/results
```

Agents should reason about changes by pipeline stage and ownership boundary:
- Source identity, source ranges, and source-origin behavior usually begin in `pkg/source` and affect parser, diagnostics, compiler, formatter, or tooling call sites.
- Syntax changes usually begin in grammar/parser and continue into compiler lowering.
- Diagnostic changes usually involve `pkg/diagnostics` plus the parser/compiler/runtime call sites that create or wrap the diagnostic.
- Semantic/runtime changes usually live in compiler, runtime, or VM.
- Bytecode changes usually require coordinated updates in `pkg/bytecode`, compiler emission, VM execution, and low-level tooling such as `pkg/asm` or debugger metadata.
- Runtime value behavior usually belongs in `pkg/runtime`; VM, stdlib, encoding, and debugger should consume those semantics rather than redefine them.
- Output and materialization changes usually involve `pkg/encoding`, runtime values, and VM result handling.
- Debugging changes usually involve `pkg/debugger`, VM execution hooks/state, compiler or bytecode metadata, and runtime value inspection contracts.
- Built-in module/function changes usually belong in `pkg/stdlib`, while reusable module contracts belong in `pkg/module` or `pkg/sdk`.
- Embedding/API changes usually affect the top-level package and integration boundaries.
- File system access, sandboxing, or path-policy behavior usually belongs in `pkg/fs`, not parser/compiler logic directly.

## Canonical invariants

- Ferret v2 uses a register-based VM.
- `runtime.Value` is the common runtime/VM value abstraction.
- Parser-generated code is derived output, not the source of truth.
- Compiler changes must preserve program semantics expected by the VM.
- Optimizations must preserve correctness before performance.
- Runtime execution errors and internal invariant violations are different classes of failure and should not be collapsed conceptually.
- Do not assume behavior from old design notes or the v1 codebase unless it is reflected in the current v2 code.
- Do not change FQL language semantics as a side effect of refactoring.

## Package map

Agents should begin with the package whose responsibility owns the requested behavior. Do not infer ownership from file names alone when a package in this map already describes the intended boundary.

### Core execution pipeline

* `pkg/source`
    * Owns source text, source identity, source ranges, and source-origin metadata.
    * Prefer this package when behavior depends on where code came from.
    * Parser, diagnostics, compiler, formatter, and tooling may depend on it, but source identity should not be reimplemented in those packages.
* `pkg/parser`
    * Owns FQL syntax parsing, parse-tree processing, parser diagnostics, and parser-generated code integration.
    * Grammar changes should begin under pkg/parser/antlr.
    * pkg/parser/antlr contains grammar sources.
    * pkg/parser/fql contains generated parser and lexer code.
    * Do not hand-edit generated parser artifacts; edit grammar sources and regenerate.
* `pkg/diagnostics`
    * Owns shared diagnostic primitives and formatting support for errors, warnings, spans, labels, notes, hints, and user-facing messages.
    * Parser, compiler, runtime, formatter, and tooling should use shared diagnostic concepts rather than inventing local diagnostic formats.
    * Changes here should preserve diagnostic category, span, label, note, and hint quality.
* `pkg/compiler`
    * Owns semantic analysis, lowering from parsed FQL into bytecode.Program, bytecode emission, and optimization/code generation passes.
    * Compiler changes must preserve the runtime semantics expected by the VM.
    * Do not move runtime-only behavior into the compiler unless the behavior is explicitly compile-time validation or compile-time semantics.
* `pkg/bytecode`
    * Owns the executable program model consumed by the VM and produced by the compiler.
    * Includes instructions, operands, programs, and related executable metadata.
    * Changes here are cross-cutting and usually require coordinated updates in compiler emission, VM execution, debugger metadata, and low-level tooling such as pkg/asm.
* `pkg/vm`
    * Owns bytecode execution, VM state, instruction dispatch, runtime coordination, cleanup, and VM-facing result handling.
    * This package is performance-sensitive and semantics-sensitive.
    * Verify compiler, bytecode, and runtime assumptions before changing VM internals.
* `pkg/runtime`
    * Owns core runtime semantics: values, value comparison/equality behavior, type behavior, function registries, runtime contracts, and execution-facing interfaces.
    * VM, stdlib, encoding, and debugger should consume runtime semantics rather than redefine them locally.
    * Prefer this package for shared value behavior that must remain consistent across execution, materialization, stdlib functions, and debugging.

Output, formatting, and low-level tooling

* `pkg/encoding`
    * Owns output encoding and materialization infrastructure for turning runtime values into external representations.
    * Changes here must account for runtime value semantics, resource ownership, cleanup behavior, and result lifetimes.
* `pkg/formatter`
    * Owns Ferret source formatting and pretty-printing behavior.
    * Formatting changes should preserve source semantics and should be validated with targeted tests or fixtures.
* `pkg/asm`
    * Owns assembly-layer support for Ferret bytecode programs.
    * Includes parsing, encoding, and low-level textual representations used by bytecode tooling and debugging.
    * Opcode or bytecode shape changes may require corresponding updates here.

Debugging and developer tooling

* `pkg/debugger`
    * Owns Ferret debugging support: debugger state, breakpoints, stepping coordination, execution inspection, and debugger-facing protocol adaptation.
    * It should consume VM execution hooks, compiler/bytecode metadata, and runtime value inspection contracts.
    * It should not own core runtime value semantics.
* `pkg/logging`
    * Owns internal logging support.
    * Logging should remain observational.
    * Do not make language semantics, diagnostics, execution behavior, or control flow depend on log output.

Integration and extension surfaces

* `pkg/module`
    * Owns Ferret module API contracts, host module interfaces, and module registration boundaries.
    * Use this package for reusable module integration concepts.
    * Do not place stdlib-specific behavior here.
* `pkg/sdk`
    * Owns extension-facing helpers and contracts for developers building on top of Ferret internals.
    * Prefer `pkg/sdk` when the goal is to support external integrations, tools, or custom runtime/module implementations.
    * Do not move internals into `pkg/sdk` only to make cross-package access easier.
* `pkg/stdlib`
    * Owns built-in Ferret namespaces, modules, and host functions registered as the standard library.
    * Built-in functions should delegate shared semantics to runtime-owned helpers when behavior must be consistent outside stdlib.
    * Do not duplicate runtime value semantics inside stdlib functions.

File systems, internals, and support packages

* `pkg/fs`
    * Owns file system abstractions and security-aware file access.
    * Use this package for controlled file access, virtual file systems, path policy, and sandbox-like behavior.
    * Do not confuse it with pkg/source; source identity and file system access are related but separate concerns.
* `pkg/internal`
    * Owns implementation-only packages that are not intended for use outside of the Ferret project.
    * Prefer `pkg/internal` for shared implementation details that should not become public or extension contracts.
    * Do not expose APIs from here to external users.

Top-level package and regression suites

* The top-level package owns the embedding surface used by applications:
    * Engine compiles and runs FQL sources.
    * Plan wraps compiled bytecode plus environment state.
    * Session executes a plan.
    * Module and ModuleRegistry expose extension points for host modules.
* Changes here are API-sensitive and should be treated as embedding-facing behavior changes.
* `test/integration/compiler`, `test/integration/optimization`, and `test/integration/vm` are the main regression suites.

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
    - inspect relevant assumptions in `pkg/vm`
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
    - inspect `pkg/source`
    - inspect parser/compiler call sites that consume source objects
    - validate path-aware diagnostics and any embedding/tooling behavior that depends on source identity

- Change embedding API:
    - inspect top-level package (`Engine`, `Plan`, `Session`)
    - inspect downstream compiler/runtime/VM interactions
    - validate public behavior with integration coverage as appropriate

- Change built-in functions/modules:
    - inspect `pkg/stdlib`
    - inspect host function/module registration and runtime contracts

- Change extension or integration support for external developers/tools:
    - inspect `pkg/sdk`
    - validate that the change belongs to an extension-facing surface rather than the core pipeline

- Change developer tooling or low-level program tooling:
    - inspect `pkg/asm`, `pkg/debugger`, `pkg/formatter`, `pkg/diagnostics`, or `pkg/sdk` depending on which surface owns the behavior

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
- debugger integration points

Do not treat historical discussion, stale comments, or old branches as authoritative.

## Public API and package boundary rules

- Treat the top-level package, `pkg/module`, `pkg/runtime`, and `pkg/sdk` as API-sensitive.
- Do not export new symbols from API-sensitive packages unless the task explicitly requires an external contract.
- Prefer unexported helpers inside the owning package before adding exported APIs.
- If a new exported symbol is necessary, add a doc comment that explains the external contract and stability expectation.
- Do not move internals into `pkg/sdk` only to make tests or cross-package access easier.
- Do not expose debugger-only APIs through the public embedding surface unless explicitly requested.

## Language behavior change rules

- Do not change FQL language semantics as a side effect of refactoring.
- Any intentional language behavior change must be called out explicitly in the final summary.
- Backward-incompatible behavior changes require tests showing both the old edge case and the new expected behavior.
- When behavior differs from v1 or from old docs, prefer current v2 tests and this file.
- Preserve existing behavior unless the task explicitly requires changing it.

## Runtime value correctness rules

- Hashes are acceleration hints, not proof of equality.
- Any uniqueness, `DISTINCT`, set, grouping, or deduplication behavior must verify equality after hash comparison.
- Reuse shared runtime equality/comparison semantics instead of inventing local equality rules.
- Do not introduce hash-only correctness paths unless the behavior is explicitly probabilistic, which language semantics normally are not.
- Preserve Ferret type ordering and comparison semantics consistently across compiler, VM, stdlib, and encoding behavior.

## Resource and lifecycle rules

- Values or results that own resources must document ownership and cleanup behavior.
- Cleanup must be deterministic where the API exposes `Close` or equivalent lifecycle methods.
- VM execution must preserve cleanup behavior on normal return, runtime error, and cancellation paths.
- Do not materialize lazy or streaming values eagerly unless the task explicitly requires it.
- Encoding/materialization changes must consider resource ownership and whether values can outlive the current execution frame.

## Debugging architecture rules

- Debugger support must not change normal execution semantics.
- Debugger-disabled execution paths must remain allocation-conscious and should avoid debugger-specific work on hot VM paths.
- Runtime values may expose debugger-facing metadata or display information through runtime-owned contracts, not debugger-owned runtime type checks.
- Debugger contracts should live near the value/runtime boundary when they describe value behavior.
- `pkg/debugger` should consume those contracts and translate them into debugger state, protocol state, or user-facing inspection data; it should not own core runtime semantics.
- Prefer optional interfaces over modifying every runtime value type.
- Debugger integration in the VM must be explicit, measurable, and easy to bypass when debugging is disabled.

## Diagnostic quality rules

- User-facing errors should include accurate source spans whenever source information is available.
- Prefer actionable hints when an error is likely caused by a common misuse.
- Do not replace specific diagnostics with generic errors.
- Parser/compiler diagnostics should distinguish syntax errors, semantic errors, runtime type errors, and internal invariants.
- Tests for diagnostics should verify message category and span accuracy for behavior changes.

## Standard library rules

- Built-in functions should keep Ferret-facing argument validation close to the function boundary.
- Prefer small host functions that delegate core behavior to runtime-owned helpers when behavior is shared.
- Do not duplicate runtime semantics inside stdlib functions.
- Stdlib errors should be user-facing and should preserve argument context where practical.
- Test stdlib behavior at the Ferret-language level whenever practical.

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

Rationale:

- one method-bearing type per file keeps ownership of behavior obvious
- standalone method-bearing types make diffs and reviews clearer
- grouped type blocks are fine for passive, closely related types, but should not hide substantial behavioral types

## Function and method ownership rules

These rules are mandatory unless the task explicitly requires otherwise.

- A file centered on a method-bearing type should contain the type, its methods, and its constructors only.
- Do not mix package-level helper functions into a file that already contains methods for a primary type.
- In type-centered files, constructor functions are the only normally allowed package-level functions.
- If logic conceptually belongs to the primary type, implement it as a method.
- If logic does not belong to the type and must remain a package-level function, place it in a separate helper-focused file.
- Package-level functions are preferred only when there is no natural owning type or when the behavior is genuinely package-level.
- If a file contains both methods and non-constructor package-level functions, that is usually a structure violation and should be refactored.

## Comment rules for functions and methods

- Do not add comments to every function or method by default.
- Exported functions and methods should usually have doc comments, especially in public, embedding-facing, or extension-facing packages.
- Unexported functions and methods should be commented only when they carry non-obvious behavior, invariants, side effects, ownership rules, cleanup expectations, or protocol/lifecycle constraints.
- Comments must explain intent, contract, invariants, side effects, or lifecycle behavior.
- Prefer comments that explain why the code exists, what must remain true, or how the method is meant to be used.
- Do not write comments that merely restate the method name or signature.
- For VM, runtime, compiler, encoding, diagnostics, and debugger internals, prefer comments on semantics and invariants over implementation narration.
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

## Response and code style

When assisting with this repository, avoid large unstructured blocks of prose or code.

Prefer responses that are easy to scan:

- Use short sections with clear headings.
- Use bullet points for decisions, trade-offs, and follow-up work.
- Use code blocks only for actual code, commands, or configuration.
- Prefer focused snippets or diffs over full-file dumps.
- Explain why a change is needed before showing how to implement it.
- Keep comments in code useful and minimal.
- Avoid repeating the same context in multiple places.
- When the change touches multiple files, summarize the role of each file first.

The expected tone is practical, concise, and engineering-focused.

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

1. Identify the owning subsystem.
2. Identify the contract, invariant, or behavior being preserved or changed.
3. Choose the smallest reasonable implementation that fits the existing design.
4. Determine whether the change is significant.
5. Add or update correctness tests.
6. Add or update benchmarks if the change is significant.
7. Run relevant validation and summarize the results accurately.

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

- `pkg/vm`
- `pkg/runtime`
- `pkg/compiler`
- `pkg/bytecode`
- `pkg/encoding`
- parser/compiler hot paths
- caching, pooling, register allocation, ownership tracking, or materialization logic
- debugger hooks on execution hot paths

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
- compare before/after results, preferably including `ns/op`, `B/op`, and `allocs/op`
- report the benchmark command used and summarize the performance delta

If no relevant benchmark exists for the changed hot path, add one.

If benchmark tooling or environment is unavailable, state that explicitly and do not claim benchmark validation was completed.

## Test placement rules

- Parser syntax behavior should have parser-focused tests or fixtures.
- Compiler semantic behavior should have compiler tests and diagnostics/span assertions when relevant.
- Bytecode emission changes should include compiler or integration tests that verify emitted behavior, not just VM behavior.
- VM opcode behavior should have VM-level tests plus integration coverage when user-visible.
- Stdlib behavior should be tested at the Ferret-language level whenever practical.
- Public embedding behavior should have top-level API tests, not only package-internal tests.
- Debugger behavior should test protocol/inspection output separately from VM execution semantics when possible.

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
- `make` is optional but is the preferred entrypoint for repo-defined workflows.
- Java plus ANTLR 4.13.2 are required when regenerating parser artifacts.
- `staticcheck`, `goimports`, and `revive` are needed for lint/format flows; install them with `make install-tools`.

## Command matrix

- Broad validation: `go test ./...`
- Race-heavy package and integration coverage: `make test`
- Lint: `make lint`
- Format: `make fmt`
- Regenerate parser/codegen artifacts: `make generate`
- Build the CLI binary: `make compile`

Run `make generate` only when grammar or generator inputs change.

## Editing rules

- Never hand-edit generated files under `pkg/parser/fql` or `pkg/parser/antlr/gen`.
- Parser generation is driven by `pkg/parser/parser.go`:
    - `antlr -Xexact-output-dir -o fql -package fql -visitor -Dlanguage=Go antlr/FqlLexer.g4 antlr/FqlParser.g4`
    - `go run ./tools/patch_lexer.go`
- If you change grammar files in `pkg/parser/antlr`, run `make generate` and commit the generated output in the same change.
- Treat `Makefile` and `.github/workflows/build.yml` as the source of truth for validation commands.
- Prefer narrow validation first, then broaden:
    - Package-local changes: run the affected `go test` package or packages.
    - Compiler, optimizer, or VM changes: run the relevant integration suites.
    - Cross-cutting changes: finish with `go test ./...` or `make test`.

### Validation expectations

- After code changes, run the narrowest tests that prove the behavior you touched.
- Before finishing broader changes, run the relevant repo-level command from the matrix above.
- If you changed formatting-sensitive files, run `make fmt`.
- If you changed lint-sensitive code paths or public behavior, run `make lint` when the toolchain is available.
- If you changed parser grammar, generated lexer/parser output must be included and reviewed.

### Expectations for non-trivial changes

When proposing or implementing non-trivial changes:

- identify the owning subsystem first
- preserve invariants unless the task explicitly changes them
- prefer local, comprehensible changes before introducing new abstractions
- distinguish correctness work from performance work
- do not perform opportunistic refactors unrelated to the requested task unless they are necessary for correctness

## Secondary references

- `README.md` for product context and links to the broader Ferret ecosystem.
- `CONTRIBUTING.md` for human contributor process.
- `.github/workflows/build.yml` for the current CI validation path.
