# AGENTS.md

This file is the canonical operating guide for coding agents working in this
repository. It applies to Ferret v2 only. Do not import assumptions from the
separate v1 branch unless current v2 code and tests establish the same behavior.

## Sources of truth

Use the most direct repository authority for facts that can change:

* `go.mod` owns the module path and minimum Go version.
* `Makefile` owns local targets, tool versions, and command composition.
* `.github/workflows/build.yml` owns the primary CI validation path and tested Go
  versions; other workflow files own their named automation.
* Grammar sources under `pkg/parser/antlr` own FQL syntax. Generated parser files
  are derived output.
* Current code and tests own architecture and behavior. Historical notes, old
  branches, and stale comments are not authoritative.
* `README.md` provides product context, while `CONTRIBUTING.md` describes the
  human contribution process.

When these sources disagree with descriptive documentation, verify the current
implementation and correct the documentation rather than copying stale values.

## Development documentation

Detailed contributor documentation lives under `docs/development`. Before making
a substantial change, read the guide for the affected subsystem; do not load
every guide for unrelated work.

* [Architecture and package ownership](docs/development/architecture.md)
* [Runtime, VM, values, and resource lifecycle](docs/development/runtime.md)
* [Debugger](docs/development/debugger.md)
* [Modules, SDK, and standard library](docs/development/modules.md)
* [Tooling, generation, tests, and benchmarks](docs/development/workflow.md)
* [Release automation](docs/development/release.md)

## Universal architecture and ownership

Ferret is a compiled query language and runtime. Its primary pipeline is:

```text
source -> parser -> parse tree/diagnostics -> compiler -> bytecode.Program -> VM -> runtime values -> output
```

The following invariants apply across the repository:

* Ferret v2 uses a register-based VM.
* `runtime.Value` is the shared runtime and VM value abstraction.
* Compiler output must preserve the program semantics expected by the VM.
* Optimizations must preserve correctness before performance.
* Runtime execution errors and internal invariant violations are different
  failure classes and must not be collapsed.
* Logging is observational; language semantics and control flow must not depend
  on log output.
* FQL semantics must not change as a side effect of refactoring.

Begin in the package that owns the requested behavior:

| Concern | Primary owner |
| --- | --- |
| Source identity and spans | `pkg/source` |
| Syntax and parse-tree integration | `pkg/parser` |
| Shared diagnostics | `pkg/diagnostics` |
| Semantic analysis, lowering, and optimization | `pkg/compiler` |
| Executable program model and artifacts | `pkg/bytecode` |
| Execution and VM state | `pkg/vm` |
| Value semantics and runtime contracts | `pkg/runtime` |
| Output encoding and materialization | `pkg/encoding` |
| Source formatting | `pkg/formatter` |
| Assembly-layer bytecode tooling | `pkg/asm` |
| Source-level debugging | `pkg/debugger` |
| Module contracts and authoring support | `pkg/module`, `pkg/sdk` |
| Built-in functions and namespaces | `pkg/stdlib` |
| Controlled filesystem and network access | `pkg/fs`, `pkg/net` |
| Embedding lifecycle and composition | top-level `ferret` package |

Do not duplicate an owning package's semantics in a consumer. In particular,
runtime value behavior belongs in `pkg/runtime`, not in VM, stdlib, encoding, or
debugger-specific type switches.

## Public API and compatibility

Treat the top-level package, `pkg/module`, `pkg/runtime`, and `pkg/sdk` as
API-sensitive.

* Preserve existing public and language-visible behavior unless the task
  explicitly changes it.
* Do not export new symbols unless the task requires an external contract.
* Prefer unexported helpers in the owning package before expanding public API.
* Add contract-focused doc comments to necessary exported symbols.
* Do not move internals into `pkg/sdk` merely to make tests or cross-package
  access easier.
* Do not expose debugger-only behavior through the embedding surface unless the
  task explicitly requires it.
* Call out every intentional FQL semantic change in the final report.
* Backward-incompatible behavior changes require coverage for the former edge
  case and the new expected behavior.

The top-level embedding surface centers on `Engine`, `Plan`, `Session`,
`DebugSession`, and `Output`. Reusable module contracts live in `pkg/module`.

## Generated files and language changes

Never hand-edit generated parser artifacts under `pkg/parser/fql` or the
generated `pkg/parser/antlr/FqlLexer.tokens` vocabulary. Edit the `.g4` grammar
sources and regenerate with `make generate`.

Syntax changes normally require coordinated grammar, generated parser,
diagnostic/parser integration, compiler lowering, formatter, and integration
coverage. Review all generated changes and commit them with their source change.

## Go type and file structure

These rules are mandatory unless the task explicitly requires otherwise.

* Prefer grouped `type ( ... )` declarations for package-level types.
* Types declared in the same file should normally be placed in a single grouped
  `type` declaration rather than written as independent `type` declarations.
* This applies equally to structs, interfaces, aliases, named primitive types,
  and method-bearing types.
* Do not split types into independent declarations merely because one or more of
  them have methods.
* Keep related types together when they belong to the same narrow responsibility
  and their proximity makes the implementation easier to understand.
* A file may contain multiple related behavioral types when they form one
  cohesive concern.
* Split types into separate files based on responsibility and ownership, not
  simply because multiple types have methods.
* When a file contains only one package-level type, a standalone declaration is
  acceptable; do not create an artificial group containing a single type.
* When adding a package-level type to a file that already contains type
  declarations, incorporate it into the existing type group when it belongs to
  the same concern.
* Avoid scattering a cohesive family of small types across multiple files.
* Do not create `helpers.go`, `utils.go`, or similarly generic files as dumping
  grounds. Organize files around predictable responsibilities.

Preferred:

```go
type (
	PassResult struct {
		Metadata map[string]any
		Modified bool
	}

	PassContext struct {
		Program *bytecode.Program
	}

	Pass interface {
		Run(*PassContext) (*PassResult, error)
	}
)
```

Avoid independent declarations when the types belong to the same concern:

```go
type PassResult struct {
	Metadata map[string]any
	Modified bool
}

type PassContext struct {
	Program *bytecode.Program
}

type Pass interface {
	Run(*PassContext) (*PassResult, error)
}
```

The grouped declaration expresses that these types form one related family.

## Function and method ownership

These rules are mandatory unless the task explicitly requires otherwise.

* Organize files around cohesive responsibilities rather than individual types.
* A file may contain multiple related types and their methods when they
  participate in the same narrow concern.
* Keep methods close to the types they belong to.
* A file containing methods must not also contain unrelated package-level
  functions.
* Package-level functions may coexist with methods in a type-centered file only
  when they are constructors for types owned by that file.
* Constructors include conventional `New...` functions and other explicit
  construction functions whose primary responsibility is creating or
  initializing one of the file's types.
* If package-level behavior is not a constructor and has no natural receiver,
  place it in a separate responsibility-focused file.
* If behavior conceptually belongs to a type's state, invariants, lifecycle, or
  resources, implement it as a method rather than a package-level function.
* Do not keep a regular function beside methods merely because that function is
  used only by those methods.
* Split a file when it contains distinct responsibilities, not merely because it
  contains multiple behavioral or method-bearing types.
* Do not split cohesive behavior across files merely to enforce one type or one
  method-bearing type per file.

Preferred:

```go
type (
	Registry struct {
		entries map[string]*Entry
	}

	Entry struct {
		value Value
	}
)

func NewRegistry() *Registry {
	return &Registry{
		entries: make(map[string]*Entry),
	}
}

func (r *Registry) Add(name string, value Value) {
	// ...
}

func (r *Registry) Get(name string) (*Entry, bool) {
	// ...
}
```

Avoid mixing regular package functions with methods:

```go
func (r *Registry) Add(name string, value Value) {
	// ...
}

func normalizeName(name string) string {
	// ...
}

func (r *Registry) Get(name string) (*Entry, bool) {
	// ...
}
```

If `normalizeName` is intrinsic to `Registry`, make it a method. If it is
genuinely package-level behavior, move it to an appropriately named
responsibility-focused file.

## Comment conventions

* Do not comment every function or method mechanically.
* Exported functions and methods should normally have doc comments, especially
  in embedding-facing and extension-facing packages.
* Comment unexported code only when it carries non-obvious semantics,
  invariants, side effects, ownership, cleanup, recovery, or protocol behavior.
* Explain why the code exists, what must remain true, or how it must be used.
* Do not merely restate the symbol name or signature.
* Prefer semantic and lifecycle comments in compiler, VM, runtime, encoding,
  diagnostics, and debugger internals.
* Avoid comment wallpaper.

Preferred:

```go
// Close releases resources associated with the result.
// It is safe to call multiple times. Once closed, the result must not be reused.
func (r *Result) Close() error
```

Avoid comments such as `// Close closes the result.`

## Go control-flow spacing

These rules are mandatory for handwritten Go code. Blank lines separate logical
units and make control transfer visible.

### Producer and immediate check

A declaration, assignment, call, assertion, lookup, or parse operation stays
adjacent to an `if` that immediately validates or consumes its result:

```go
res, err := doSome()
if err != nil {
	return err
}

named, ok := value.(*types.Named)
if !ok {
	return ErrUnsupported
}
```

Do not insert a blank line between the producer and its immediate check.

If this producer/check unit follows separate logic, add a blank line before it:

```go
prepareState()

value := lookup(name)
if value == nil {
	return ErrNotFound
}
```

No leading blank line is needed when the producer starts the enclosing block.

### Independent control flow

Separate independent `if` blocks with a blank line:

```go
if foo != nil {
	useFoo(foo)
}

if bar != nil {
	useBar(bar)
}
```

After a completed control-flow block, add a blank line before a separate
statement or logical unit.

### Return and break

`return` and `break` begin a new logical group when another statement precedes
them in the same block:

```go
result := buildResult()

return result
```

```go
if is(curr, "WAITFOR") {
	found = true

	break
}
```

No blank line is required when `return` is already the first statement in its
block. Do not add an artificial leading blank line at the start of a function or
block.

## Local type declarations

A function-local type is appropriate when it is small, passive, method-free,
used only by that function, and makes the local algorithm easier to understand.

Prefer a package-level unexported type when it represents a meaningful domain or
algorithm concept, spans a substantial part of a complex function, benefits
from being visible outside control flow, may gain methods, or is likely to be
reused by nearby helpers.

Do not promote a tiny throwaway struct merely for consistency, and do not hide a
meaningful concept inside a function merely to avoid a package-level type.

## Engineering discipline

Preserve correctness and subsystem boundaries before pursuing cleanup or
performance. Prefer the smallest local change that fully solves the task. New
abstractions, indirection, package moves, and refactors require a concrete
correctness, ownership, or maintainability need.

For every non-trivial change:

1. Identify the owning subsystem and read its development guide.
2. Identify the contract, invariant, and existing behavior being preserved or
   intentionally changed.
3. Choose the smallest implementation that fits the current architecture.
4. Add or update correctness tests for behavior changes.
5. Decide whether the change is significant for performance and benchmark it
   when required.
6. Run the narrowest relevant validation, then broaden in proportion to risk.
7. Evaluate documentation impact and update affected repository and public
   documentation.
8. Review the complete resulting implementation and diff using the mandatory
   self-review below.
9. Fix actual findings and rerun every affected test, check, and benchmark.
10. Report changed behavior, documentation impact, evidence, review results,
    and limitations accurately.

A task is not complete merely because the first implementation compiles or its
tests pass. Do not perform opportunistic refactors unrelated to the requested
change.

## Tests and validation

Add or update tests for every behavior change. Test at the layer that owns the
contract and add integration coverage when behavior crosses package boundaries:

* Parser syntax belongs in parser tests or fixtures.
* Compiler semantics belong in compiler tests, including diagnostic category and
  span assertions when relevant.
* Bytecode emission requires compiler or integration coverage, not only VM
  behavior tests.
* VM opcode behavior requires VM tests and integration coverage when visible to
  users.
* Stdlib behavior should be exercised through FQL whenever practical.
* Embedding behavior requires top-level API coverage.
* Debugger protocol and inspection behavior should be separable from normal VM
  execution tests.

Run the narrowest relevant command first. Use `Makefile` and
`.github/workflows/build.yml` to select current repository-level commands. Run
`make generate` only when generator inputs change, and include the reviewed
generated output. Re-run affected validation after every review-driven change.

Run `make fmt` when handwritten Go formatting is affected and inspect its broad
rewrite surface. Run `make lint` for lint-sensitive code or public API changes
when the required tools are available. Broader, cross-package changes should
finish with the relevant repository-level test target.

Do not claim that tests, lint, generation, benchmarks, or review passed unless
they actually ran successfully. Report tooling or environment limitations
explicitly.

## Significant changes and benchmarks

A change is significant when it could reasonably affect execution throughput,
compile time, common-path latency, allocations, memory reuse, pooling, cleanup,
materialization, optimizer output, or code generation.

This commonly includes changes in `pkg/vm`, `pkg/runtime`, `pkg/compiler`,
`pkg/bytecode`, `pkg/encoding`, parser/compiler hot paths, ownership tracking,
caching, pooling, register allocation, or debugger hooks on execution paths.
Documentation-only, test-only, formatting-only, and behavior-neutral rename
changes are normally not significant.

For a significant change:

1. Run the relevant benchmark before implementation and save the baseline.
2. Run the same benchmark after implementation.
3. Compare `ns/op`, `B/op`, and `allocs/op` where available.
4. Investigate meaningful regressions before completing the task.
5. Add a benchmark when no relevant one covers the changed hot path.
6. Report the exact commands and summarized delta.

If the environment cannot run a required benchmark, say so and do not claim
benchmark validation. Never trade clear correctness or maintainability for a
speculative micro-optimization.

## Mandatory final self-review

After implementation and initial validation, review the complete resulting
change before considering any non-trivial task finished. This is a second-pass
evaluation of the implementation, not a confirmation that tests passed.

### Correctness and lifecycle

Verify the task is completely satisfied. Look for missing cases, incorrect
assumptions, boundary conditions, invalid states, cancellation and concurrency
errors, cleanup failures, resource leaks, ownership mistakes, and incorrect
error propagation. Confirm public and FQL-visible behavior matches the intended
contract. Check idiomatic Go error wrapping, context propagation, synchronization,
and lifecycle management. For bug fixes, prefer a regression test that fails
without the fix.

### Architecture and API

Verify responsibilities remain in the correct package, type, and layer. Check
dependency direction, compile-time/runtime separation, runtime-owned semantics,
public API necessity, and compatibility. Reject duplicated semantics, leaked
implementation details, misplaced behavior, and abstractions at the wrong level.

### Clarity and organization

Look for unnecessary complexity, duplication, nesting, misleading names, dead
branches, debugging artifacts, and comments about abandoned approaches. Check
the type/file/function ownership rules above. Avoid both overloaded files and
unnecessary fragmentation into excessive helpers or abstractions.

### Tests and performance

Review whether tests cover meaningful positive, negative, boundary, error,
cleanup, and cancellation cases without merely mirroring implementation. Check
assertion strength and unnecessary brittleness. For significant changes, inspect
allocations, repeated work, materialization, synchronization, and benchmark
comparability.

### Scope and complete diff

Inspect the complete final diff, not only individual files. Verify that:

* every changed line belongs to the request or a necessary supporting change;
* no temporary code, accidental API or behavior change, or unrelated refactor
  remains;
* generated files changed only because their source inputs changed;
* tests describe intended behavior;
* comments describe current contracts and invariants;
* file, type, function, and package boundaries remain coherent;
* resource ownership and lifecycle behavior remain correct;
* affected repository and public documentation has been updated, or any
  unavailable external documentation dependency is explicitly reported;
* the final implementation is the smallest coherent solution.

When review finds a correctness, architecture, ownership, lifecycle, API,
organization, performance, documentation, or meaningful coverage problem, fix
it and repeat affected validation. Minor optional style preferences do not
justify churn.

Do not use self-review to expand scope through speculative refactoring,
unrelated cleanup, unrelated API redesign, broad package moves, or FQL semantic
changes outside the task.

## Change and reporting discipline

Prefer an existing local pattern over a new architectural pattern. Leave
already-correct code alone. If requested work exposes a necessary supporting
cleanup, keep it narrow and explain why it is required.

The final report for a non-trivial change must state:

* owning subsystem and files changed;
* behavior and invariants preserved or intentionally changed;
* tests and benchmarks added or updated;
* validation and benchmark commands actually run;
* documentation updated, or documentation impact explicitly evaluated as none;
* completion of the mandatory final self-review;
* material findings corrected during review;
* remaining limitations or skipped validation.

## Documentation synchronization

Documentation is part of the change, not a follow-up activity. Before completing
every non-trivial task, evaluate whether the implementation changes any
documented architecture, ownership boundary, invariant, workflow, API, behavior,
example, or contributor guidance.

Update the relevant documentation in the same task:

* Update `docs/development/*` when repository architecture, subsystem
  responsibilities, internal contracts, lifecycle behavior, development
  workflows, tooling, testing, benchmarking, or release behavior changes.
* Update repository-facing documentation such as `README.md`,
  `CONTRIBUTING.md`, or other local documentation when their documented behavior,
  setup instructions, workflows, or examples are affected.
* Update the corresponding website documentation when changing FQL syntax or
  semantics, the embedding API, SDK contracts, extension points, standard
  library behavior, or other documented public behavior.
* Update both repository and website documentation when a change affects both
  internal development guidance and public behavior.

Do not update documentation mechanically when the implementation does not affect
its contract, behavior, examples, or guidance. Documentation-only churn is not a
substitute for evaluating documentation impact.

Public documentation is maintained in the website repository. When that
repository or another required documentation source is available, make the
necessary documentation changes as part of the same task. If it is unavailable,
identify the exact required follow-up explicitly in the final report rather than
silently leaving known documentation stale.

## Response style

Keep responses practical and easy to scan. Use short sections, focused bullets,
and code blocks only for code, commands, or configuration. Explain why a change
is needed before how it works, summarize each changed file's responsibility, and
avoid repeating the same context.