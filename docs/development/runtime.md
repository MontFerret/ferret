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

The root embedding layer uses this mechanism to produce `encoding.Output`.
Encoders can discover nested resources during traversal, so encoding failures and
missing-codec failures must still close all resources already owned or adopted.

## Engine, plan, and session lifecycle

The root embedding lifecycle is hierarchical:

* `Engine` owns the configured host, compilers, loader, module hooks, session
  limiter, filesystem, and any network service it created.
* `Plan` owns one compiled program, plan/session hooks, and a bounded or
  unbounded pool of VMs for that program.
* `Session` borrows a VM, constructs an execution environment, injects logging,
  encoding, filesystem, and network services into the context, and materializes
  output. By default it borrows the engine filesystem.
  `ferret.WithSessionFSRoot` replaces only that session's root while retaining
  the engine's read-only policy; the session owns and closes the replacement
  filesystem.

`Engine.Run` is a convenience path that owns and closes its temporary plan and
session. A caller that creates a plan or session directly owns its `Close` call.
`Session.Close` is idempotent and concurrency-safe; it runs close hooks and
returns its VM to the pool. `Plan.Close` rejects new sessions, runs close hooks,
and closes retained VMs. Engine cleanup must continue even when a hook fails.

Before hooks run in registration order. After and close hooks unwind in reverse
order, with the error behavior defined by `pkg/module` and the root hook
implementation. See [Modules, SDK, and standard library](modules.md).

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
instruction execution and ownership transitions, and top-level tests for
Engine/Plan/Session lifecycle and output behavior. Cross-layer semantics should
also have coverage in `test/integration/vm` or the relevant compiler/optimization
suite.

Resource tests should cover aliases, borrowed versus owned values, failure,
cancellation, unwind, materialization, and idempotent cleanup. Performance work
should use the commands described in [Development workflow](workflow.md).

## Related guides

* [Architecture](architecture.md)
* [Debugger architecture](debugger.md)
* [Modules, SDK, and standard library](modules.md)
* [Development workflow](workflow.md)
