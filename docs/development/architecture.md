# Ferret v2 Architecture

Ferret v2 is a compiled, embeddable query language and runtime. This guide
describes the stable subsystem boundaries used by contributors. Current code and
tests remain authoritative for implementation details that evolve within those
boundaries.

## Execution pipeline

```text
source.Source
    -> ANTLR lexer and parser
    -> compiler frontend and diagnostics
    -> compiler lowering and optimization
    -> bytecode.Program
    -> VM execution
    -> runtime.Value
    -> encoding.Output
```

`pkg/source` owns indexed source text and identity. Portable positions, spans,
locations, and ranges are aliases of `github.com/MontFerret/api/source` types. Source-aware behavior begins
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

The root `ferret` package is the preferred public embedding façade. It contains
curated aliases, constants, and forwarding functions over `pkg/engine` and the
owning source, runtime, encoding, module, logging, diagnostics, artifact, and
debugger packages, without substantial engine implementation. Shared vocabulary
targets those owners directly, preserving type identity without an engine relay.
`pkg/engine` owns the native engine API and composes the compiler, bytecode
loader, runtime host, VM pool, modules,
filesystem, network, encoders, hooks, and logging:

```text
Engine -> Plan -> Session -> Output
              \-> DebugSession
```

An `Engine` owns host-scoped services and compilers. Compilation or artifact
loading creates a reusable `Plan`, which owns a VM pool for one program. A
`Session` borrows a VM and supplies per-execution environment and output
settings. A `DebugSession` uses retained VM state and debugger metadata. These
objects have explicit cleanup responsibilities described in
[Runtime and lifecycle](runtime.md) and [Debugger architecture](debugger.md).

The native engine owns composition and lifecycle policy, not the underlying
semantics of runtime values, modules, codecs, filesystems, networks, or debugger
inspection.

The façade fans out to `pkg/engine` for native engine semantics and directly to
lower-level packages for shared vocabulary and convenience helpers. Engine
implementation dependencies flow through `pkg/engine -> pkg/engine/internal/* ->
lower-level pkg/*`, with direct lower-level imports where needed.
`pkg/engine` exposes engine configuration, behavior, and semantic vocabulary; it
does not republish lower-level APIs for the façade. Engine internals are used
only within the engine subtree and must not import the native engine or root
façade. Lower-level production
packages must not import root `ferret`; compiler, runtime, VM, bytecode, and
stdlib remain independent of both embedding packages. Integrations such as
compatibility adapters and the SDK test harness may use `pkg/engine` because
they construct and own native executions. Public API tests and downstream-style
examples use root `ferret`.

`pkg/engine/internal/bootstrap` constructs compilers, orchestrates module
registration and host finalization, snapshots hooks, runs initialization, and
creates the session limiter. It owns rollback until a successful build transfers
its dependencies to the native Engine. `pkg/engine/internal/host` builds module
bootstrap services and snapshots host registries and lifecycle hooks.
`pkg/engine/internal/session` supplies session admission, permit release,
debugger services, and result materialization. The native engine retains public
configuration, option application, stdlib registration, optimization-level
translation, final Engine assembly, and Engine/Plan/Session orchestration.

`pkg/engine/internal/resource` owns lifecycle metadata for host services within
one owner.
The configuration creates one manager and hands it to bootstrap; successful
construction transfers that manager to the completed Engine. Session construction
creates a separate manager for each ordinary or debug session. Typed filesystem
and network references stay with their consumers.
Owned registrations carry cleanup callbacks, while borrowed registrations record
ownership outside the scope. The manager provides no service lookup or dependency
resolution.

See the [engine boundary audit](engine-boundary-audit.md) for the retained export
inventory, internal API rationale, consumer classifications, and deferred API
questions.

## Package ownership

### Language frontend and tooling

* `pkg/source`: indexed source identity and text, with API-owned coordinates.
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
| File or network policy | `pkg/fs` or `pkg/net` | native engine host/session context and stdlib adapters |
| Embedding API | root façade and `pkg/engine` | modules, VM, runtime, and public API tests |

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

The root package and `pkg/engine`, `pkg/module`, `pkg/runtime`, and `pkg/sdk` are
public, API-sensitive surfaces. `pkg/bytecode` artifacts also carry explicit versions
and validation. Do not infer compatibility promises from obsolete design notes
or the v1 branch.

## Related guides

* [Runtime and lifecycle](runtime.md)
* [Debugger architecture](debugger.md)
* [Modules, SDK, and standard library](modules.md)
* [Development workflow](workflow.md)
* [Release automation](release.md)
