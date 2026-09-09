# Runtime and Lifecycle

Ferret's runtime semantics are shared by compilation consumers, VM execution,
the standard library, encoding, modules, and debugging. Keeping those semantics
in their owning packages prevents different execution paths from disagreeing.

## Runtime values

`pkg/runtime` owns `runtime.Value` and the optional capability interfaces layered
on it. A value supplies string representation, hashing, and shallow copying;
additional interfaces opt into equality, ordering, iteration, querying,
streaming, observability, dispatch, debugging, or resource behavior.

Consumers should use the runtime's shared operations rather than matching
concrete built-in types. Host values may implement capabilities without being a
built-in value, and consumers must preserve those contracts.

Map destination operations and immutable object-library ownership are described
in [Object library contracts](object-library.md). Runtime iterator traversal
closes acquired closable iterators and preserves both traversal and close errors;
the iterable itself remains borrowed.

## Equality, comparison, and hashing

A hash is an acceleration hint, not proof of equality. Hash collisions are
valid, so uniqueness, `DISTINCT`, set, grouping, and deduplication paths must
verify equality after selecting candidates by hash.

Shared runtime equality and comparison preserve Ferret's cross-type behavior,
type ordering, fallible host comparisons, and context propagation. A value that
implements both equality and ordering must keep those results consistent, and
semantically equal values must produce equal hashes.

Do not introduce a hash-only correctness path unless the contract is explicitly
probabilistic. FQL collection semantics are not probabilistic.

## VM execution

`pkg/vm` executes `bytecode.Program` in a register file with frames, cells,
catching, host-function dispatch, iterators, collectors, streams, and other
execution state. The compiler owns emission; the VM owns the runtime meaning of
the emitted instructions.

Normal execution errors are returned as user-facing runtime errors with source
context where possible. Internal invariant violations remain a distinct class.
The VM's panic policy may recover runtime failures or propagate invariants, so
new failure handling must preserve that distinction.

The VM observes cancellation at its own safepoints. Context-aware host
operations receive the execution context and are responsible for observing
cancellation while they retain control.

## Resource ownership

Runtime values may own closable resources. `runtime.Resource` adds a stable live
resource identity to `io.Closer`; other comparable closers can also participate
in VM tracking when their identity is safe to use.

The VM tracks resource ownership and aliases through register and frame changes.
Ownership must be transferred when a value escapes a frame or becomes a result,
and discarded values must be closed exactly once. Cleanup must remain correct on
normal return, runtime error, cancellation, recovery, unwind, and explicit close
paths.

Borrowed host parameters are not automatically owned by the VM. A change that
copies, stores, returns, materializes, or aliases a resource must make the
ownership transition explicit and preserve the caller's lifecycle.

Avoid eager materialization of lazy, iterable, streaming, or resource-backed
values unless the public contract requires it.

## Results and materialization

`vm.Result` holds the root runtime value and closers transferred from execution.
Callers must close it. `Root` supports low-level inspection while the result is
open; `vm.Materialize` performs a terminal ownership-aware conversion.

A result can be materialized at most once, including when the materializer
fails. Materializers may return additional closers or explicitly adopt values
discovered during traversal. The result remains responsible for those resources
until `Close`.

The native engine uses this mechanism to produce `encoding.Output`, also exposed
through the root embedding façade.
Encoders can discover nested resources during traversal, so encoding failures and
missing-codec failures must still close all resources already owned or adopted.

## Engine, plan, and session lifecycle

`pkg/engine` owns the native lifecycle implementation; root `ferret` aliases its
public embedding types. Native Ferret protects object integrity; callers
orchestrate object lifetime:

* `Engine` owns the configured host, compilers, loader, module hooks, session
  limiter, filesystem, and any network service it created.
* `Plan` owns one compiled program, plan/session hooks, and a bounded or
  unbounded pool of VMs for that program.
* `Session` owns native run admission, before/after-run hooks, and output/error
  decisions. Its internal `session.Execution` owns the borrowed VM, environment,
  logging, encoding, filesystem and network references, materialization, and
  cleanup. By default it borrows the engine filesystem.
  `ferret.WithSessionFSRoot` replaces only that session's root while retaining
  the engine's read-only policy; the session owns and closes the replacement
  filesystem.

Engine host-resource ownership is recorded in `pkg/engine/internal/resource.Manager`.
The filesystem and Ferret-created networks are owned; `WithNetwork` explicitly
borrows the caller's network. Each acquisition registers a callback that captures
that resource instance. Replacing a named resource retires its previous owned
callback immediately and records the replacement as the newest acquisition.
Borrowed resources are never closed. A retirement error leaves the replacement
registered, so construction rollback still cleans it up; retired callbacks are
never retried.

Failed option application or stdlib registration closes the configuration's
manager. The native constructor also retains cleanup responsibility if public
optimization-level translation fails. It then hands the manager to
`pkg/engine/internal/bootstrap.Build`, which owns rollback for all returned
construction errors and transfers ownership to Engine on success. Compiler and
host-service construction failures close resources without close hooks. After
host services are constructed, registration and host-build failures unwind the
mutable hooks; initialization failures unwind the finalized hook snapshot.
Construction panics propagate without error rollback, as before.
On rollback or shutdown, eligible hooks run first, then the manager attempts
every remaining owned callback in reverse acquisition/replacement order and
joins cleanup errors.
Consequently, default construction closes the network before the filesystem;
a network created by an option is acquired earlier and closes after the
filesystem. Bootstrap rollback and native Engine shutdown keep separate
orchestration with the same cleanup ordering and error aggregation. Duplicate
Engine closure retains the completed cleanup result.

Each ordinary or debug session has a separate host-resource manager, created
after acquiring its limiter permit. It explicitly borrows the Engine filesystem
and network. `WithSessionFSRoot` replaces only that manager's filesystem entry
with an owned instance; Engine and sibling managers remain independent.
The internal session acquisition owner constructs logging, acquires capacity,
creates host resources and the environment, and constructs the VM or debugger
session. Construction errors and panics close the untransferred manager without
close hooks, with a separate defer retaining permit release even if resource
cleanup panics. A successful build transfers both responsibilities to ordinary
`Execution` or debugger services before the native cancellation recheck.
Cancellation then closes the completed session through its normal cleanup path.

Session close hooks precede host-resource cleanup, which precedes permit release
and, for ordinary sessions, return of the borrowed VM. `Session.Close` marks the
native object closed; `Execution.Close` owns this complete once-only sequence and
caches its error for repeated and concurrent callers. Debugger closure still
settles retained execution before closing the embedding session services.
Those services keep their resource manager and acquired limiter permit private;
their constructor does not transfer construction rollback. Debug services return
the permit directly, while ordinary sessions also return a VM to the plan pool.
Manager calls are serialized by construction and the owner's existing once-only
shutdown. Managers have no independent synchronization or resource lookup API.
VM resources, query-value cleanup, and limiter permits retain their separate
lifecycle mechanisms.

`Engine.Run` owns its temporary session and plan, closes them in that order,
and joins execution and cleanup errors while retaining available encoded output.
Successful cleanup preserves the original runtime diagnostic. Compilation also
preserves the original diagnostic when hooks and cancellation checks add no
failure, so `FormatError` retains its source locations and hints.
Actual additional failures are joined without changing aggregate error rendering.
A caller that creates children directly closes sessions before plans, and plans
before the engine. Parents have no descendant registries. Ordinary execution
owners cancel and settle `Run` before closing their session; debug closure
terminates and settles active commands.

Engine, Plan, and Session closure is concurrency-safe and idempotent, retaining
the completed cleanup result. Engine uses an atomic closed flag; Plan uses one
close channel for its closed state and capacity-waiter notification. Each owns
its own once-only cleanup. Closed state becomes visible before close hooks run,
and immutable program and host references remain intact during cleanup.

Engine rejects new compilation and loading without waiting for admitted work.
Plan rejects new sessions and wakes pending capacity acquisition without waiting
for constructors. The session limiter observes the request context and that
Plan's close channel, never Engine shutdown. Once capacity is acquired, parent
closure does not roll back a constructed session. Ordinary construction may
still fail if the VM pool closes before acquisition finishes; the pool owns that
synchronization. A successfully returned session remains caller-owned.

Hosts must settle outstanding work before closing resources used by their hooks
or services. Native closure does not orchestrate graceful hierarchical shutdown.
No admission lock spans options, hooks, or cleanup; a close hook must not
recursively close its own object. `Session.Close` returns its VM to the pool;
`Plan.Close` closes idle VMs, leaving borrowed VMs with sessions until returned.
Engine and Plan cleanup continue even when a hook fails.

Execution, compilation, and debugger context boundaries require non-nil contexts.
Embedding/session admission rejects existing cancellation before options or hooks
and rechecks request cancellation before returning constructed resources, closing
them on cancellation. Raw VM execution retains its existing cancellation
safepoints within an admitted run.

`Compiler.Compile(ctx, src)` synchronously checks cancellation between parsing,
lowering, and program construction; phases are not individually preempted. The
context is required. Cancellation and deadline identities remain available
through `errors.Is`, including when cancellation accompanies an option or hook
failure.

Native `PlanOption` and `SessionOption` target private native configurations,
following the native `Option` pattern also aliased by root `ferret`. They are
distinct from Universal API functional options. Non-nil options run once in
order, joining validation errors before resource acquisition. The private
session option target contains internal
`session.Config`, which owns applied settings, defaults, and host-parameter
conversion; native option constructors retain validation and ordering. Portable
semantic data may be shared with the Universal API; runtime-specific option
adaptation belongs at a future adapter boundary.

`WithPlanOptimizationLevel` on `Engine.Compile` selects native None, Basic, or
Full for one compilation. Omitting it inherits the engine's optimization;
overrides do not mutate the engine default. Unsupported levels are rejected.
`CompileDebug` accepts omission or None. Compiler instances remain immutable and
safe for shared use.

Encoded `Output` is an alias of `api/result.Output`; native source indexing,
runtime values, compiler metadata, and diagnostic rendering remain native.
`Session.Run` retains successful encoded output even when an after-run hook or
result cleanup fails. Ordinary after-run hooks still precede encoding and receive
the execution or context-validation error. Ferret encodes the successful VM
result and closes it exactly once through `Execution.MaterializeAndClose`.
That operation returns encoding and cleanup errors separately so native
`Session.Run` preserves their ordering alongside hook failures. The caller owns
the returned encoded data, which has no `Close` method and remains available
after session cleanup. Failed execution or encoding returns no output. A lone
execution diagnostic is returned unchanged.

Diagnostic aggregates expose `Unwrap() []error`, and runtime errors unwrap to
their underlying diagnostic, which in turn retains its cause.

Before hooks run in registration order. After and close hooks unwind in reverse
order, with the error behavior defined by `pkg/module` and the engine's internal
host hook implementation. See [Modules, SDK, and standard library](modules.md).

## Performance-sensitive boundaries

The VM dispatch loop, register and frame operations, runtime comparison,
collection operations, resource tracking, result materialization, pooling, and
compiler-generated execution shape are performance-sensitive. Changes should be
measured with the relevant unit or integration benchmark before and after the
implementation.

Debugger hooks on normal execution paths must remain cheap to bypass. Encoding
and inspection should avoid repeated traversal or allocation without weakening
ownership correctness.

## Testing

Use `pkg/runtime` tests for value and capability contracts, `pkg/vm` tests for
instruction execution and ownership transitions, `pkg/engine` tests for native
Engine/Plan/Session lifecycle and output behavior, and root tests for the curated
public façade. Cross-layer semantics should also have coverage in
`test/integration/vm` or the relevant compiler/optimization suite.

Resource tests should cover aliases, borrowed versus owned values, failure,
cancellation, unwind, materialization, and idempotent cleanup. Performance work
should use the commands described in [Development workflow](workflow.md).

## Related guides

* [Architecture](architecture.md)
* [Debugger architecture](debugger.md)
* [Modules, SDK, and standard library](modules.md)
* [Development workflow](workflow.md)
