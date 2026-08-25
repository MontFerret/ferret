# Ferret v2 Architecture

Ferret v2 is a compiled, embeddable query language and runtime. This guide
describes the stable subsystem boundaries used by contributors. Current code and
tests remain authoritative for implementation details that evolve within those
boundaries.

## Execution pipeline

```text
ferret.Source alias of api/source.File (embedding boundary)
    -> source.Source
    -> ANTLR lexer and parser
    -> compiler frontend and diagnostics
    -> compiler lowering and optimization
    -> bytecode.Program
    -> VM execution
    -> runtime.Value
    -> encoding.Output
```

`pkg/source` owns source text, identity, and spans. Source-aware behavior begins
there so parser, compiler, diagnostics, formatter, and debugger integrations use
the same coordinates instead of translating them independently.

`pkg/parser` creates the FQL token stream and parse tree. Grammar sources live in
`pkg/parser/antlr`; generated Go, token, and interpreter artifacts live primarily
in `pkg/parser/fql`. Parser-specific diagnostic collection lives in
`pkg/parser/diagnostics`, while `pkg/diagnostics` provides shared diagnostic and
rendering concepts.

`pkg/compiler` owns semantic analysis, lowering, optimization, and program
construction. It consumes parsed FQL and produces `bytecode.Program`; runtime-only
behavior does not belong in the compiler unless it is an explicit compile-time
semantic or validation rule.

`pkg/bytecode` owns instructions, operands, functions, executable metadata, and
the program model consumed by the VM. Its `artifact` and `format` subpackages own
serialized program loading and formats. Bytecode shape changes normally affect
compiler emission, VM execution, validation, debugger metadata, and possibly
`pkg/asm`.

`pkg/vm` owns register-based execution, frame and register state, dispatch,
cancellation safepoints, cleanup, result ownership, VM pooling, and retained
debug execution. `pkg/runtime` owns the values and semantic contracts used by the
VM. See [Runtime and lifecycle](runtime.md) for those boundaries.

`pkg/encoding` turns runtime values into external representations. Encoding is a
consumer of runtime semantics and VM result ownership; it does not redefine
value behavior.

## Embedding layer

The root `ferret` package composes the compiler, bytecode loader, runtime host,
VM pool, modules, filesystem, network, encoders, hooks, and logging into the
embedding API. Its concrete `Engine`, `Plan`, and `Session` types directly
implement `api.Runtime`, `api.Plan`, and `api.Session`:

```text
Engine -> Plan -> Session -> Output
              \-> DebugSession
```

An `Engine` owns host-scoped services and immutable compilers. Its public
compile and run boundary accepts `ferret.Source`, an exact alias of
`api/source.File`, then converts it to Core `pkg/source.Source`; parser,
compiler, diagnostics, bytecode, VM, and encoding packages do not depend on the
universal API. Per-compilation optimization options are passed explicitly
without mutating the shared compiler.

Compilation or artifact loading creates a reusable concrete `Plan`, which owns
a VM pool for one program, while the universal compile methods return it as
`api.Plan`. A `Session` borrows a VM and supplies per-execution environment and
output settings. A debug plan exposes `api/debugger.Session` through a private
translation bridge to the Core debugger. These objects have explicit cleanup
responsibilities described in
[Runtime and lifecycle](runtime.md) and [Debugger architecture](debugger.md).

The root package owns composition and lifecycle policy, not the underlying
semantics of runtime values, modules, codecs, filesystems, networks, or debugger
inspection. It aliases the Universal source, option, optimization, and output
types used by this boundary and forwards their common constructors, while
portable implementation-independent consumers may import `api` directly.

## Package ownership

### Language frontend and tooling

* `pkg/source`: source identity, text, spans, and locations.
* `pkg/parser`: grammar integration, token transformation, parse trees, and
  parser diagnostics.
* `pkg/diagnostics`: shared diagnostics, labels, notes, hints, and rendering.
* `pkg/compiler`: semantic analysis, lowering, optimization, debug metadata
  emission, and program construction.
* `pkg/formatter`: FQL pretty-printing and canonical source layout.
* `pkg/asm`: assembly and disassembly support for bytecode tooling.

### Execution and external boundaries

* `pkg/bytecode`: executable program and artifact contracts.
* `pkg/vm`: execution machinery, VM state, results, and pooling.
* `pkg/runtime`: value types, comparison, function registries, iteration,
  streaming, and host-facing runtime contracts.
* `pkg/encoding`: codecs, output encoding, and materialization integration.
* `pkg/fs`: controlled filesystem construction and access policy.
* `pkg/net`: controlled network services and HTTP policy.
* `pkg/logging`: observational logging support.

### Extensions and developer tools

* `pkg/module`: reusable engine module and lifecycle-hook contracts.
* `pkg/sdk`: supported helpers for module and host-value authors.
* `pkg/stdlib`: built-in functions, namespaces, and capability groups.
* `pkg/debugger`: source-level debugger policy over retained VM execution.
* `pkg/internal`: implementation details that are not extension contracts.

See [Modules, SDK, and standard library](modules.md) and
[Debugger architecture](debugger.md) for the extension-specific boundaries.

## Change routing

| Change | Begin with | Coordinate with |
| --- | --- | --- |
| FQL syntax | `pkg/parser/antlr` | parser diagnostics, compiler, formatter, integration tests |
| Semantic analysis | `pkg/compiler` | source spans, diagnostics, compiler tests |
| Optimization | `pkg/compiler/internal/optimization` | bytecode behavior, VM/integration tests, benchmarks |
| Opcode or program shape | `pkg/bytecode` | compiler emission, VM, artifact validation, debugger, asm |
| Runtime value semantics | `pkg/runtime` | VM, encoding, stdlib, debugger consumers |
| Execution or cleanup | `pkg/vm` | runtime ownership, embedding lifecycle, benchmarks |
| Output materialization | `pkg/encoding` | VM results and runtime resource ownership |
| Source formatting | `pkg/formatter` | parser grammar and formatter fixtures |
| File or network policy | `pkg/fs` or `pkg/net` | root host/session context and stdlib adapters |
| Embedding API | root package | compiler, modules, VM, runtime, and integration tests |

Start with the primary owner even when a behavior has several consumers. Shared
semantics should flow outward from their owner rather than being recreated at
each call site.

## Generated and source artifacts

The `.g4` files under `pkg/parser/antlr` are parser sources. The generated
`pkg/parser/antlr/FqlLexer.tokens` file and artifacts under `pkg/parser/fql` are
derived by the `go:generate` directives in `pkg/parser/parser.go`. Do not edit
derived parser output directly.

Run `make generate` after grammar changes. That target runs Go generation and
then the repository formatter, so inspect the full diff for unintended changes.
Generated changes must be committed with the grammar source that produced them.

## Diagnostics and formatting

User-facing diagnostics should retain the most accurate source span available.
Parser and compiler failures distinguish syntax errors, semantic errors, runtime
type errors, and internal invariants instead of replacing them with a generic
message. Labels, notes, and hints should make common misuse actionable. Tests for
diagnostic changes should assert both the category or message contract and span
accuracy.

`pkg/formatter` consumes current parser semantics and owns canonical FQL layout.
Formatter changes must preserve executable meaning and should be covered by
focused fixtures or integration tests. A syntax change should update formatter
handling in the same change when the new construct can be emitted or reformatted.

Low-level textual program representations belong in `pkg/asm`. They follow the
bytecode model and debugger metadata rather than defining independent opcode or
source semantics.

## Stability and compatibility

The pipeline, parser-generation workflow, and root embedding lifecycle are
stable architectural boundaries. Optimizer internals, diagnostics plumbing, VM
state, cleanup, encoding/materialization, and debugger integration are
implementation-sensitive and should be verified in current code before being
changed.

The root package and `pkg/module`, `pkg/runtime`, and `pkg/sdk` are public,
API-sensitive surfaces. `pkg/bytecode` artifacts also carry explicit versions
and validation. Do not infer compatibility promises from obsolete design notes
or the v1 branch.

## Related guides

* [Runtime and lifecycle](runtime.md)
* [Debugger architecture](debugger.md)
* [Modules, SDK, and standard library](modules.md)
* [Development workflow](workflow.md)
* [Release automation](release.md)
