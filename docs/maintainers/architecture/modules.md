# Modules, SDK, and Standard Library

Ferret extensions are divided between stable module contracts, supported
authoring helpers, built-in library registrations, and runtime-owned semantics.
Keeping those layers separate lets embedders control capabilities without
turning stdlib implementation details into module API.

## Module bootstrap

`pkg/module.Module` is the reusable engine extension contract. A module has a
stable name and registers against `module.Bootstrap` during `ferret.New`.
Registration completes before the engine host is finalized; a registration
error aborts construction and triggers cleanup of resources already created.

`Bootstrap.Host` exposes host-scoped registries and services:

* the runtime library and default parameters;
* encoding codec registration;
* logging;
* the configured filesystem and network service.

`Bootstrap.Hooks` exposes engine, plan, and session lifecycle registrars. Module
registration is the place to configure these shared services and callbacks, not
to retain mutable execution state that belongs to a session.

## Hook lifecycle

Engine hooks cover initialization and close. Plan hooks cover before/after
compilation and close. Session hooks cover before/after each run and close.

Ordering is part of the contract:

* init, before-compile, and before-run hooks execute in registration order;
* after-compile, after-run, and close hooks execute in reverse registration
  order;
* before hooks stop on the first error;
* after hooks receive the primary operation error;
* close paths aggregate errors and continue cleanup.

A before-run hook may return a derived context for later hooks and VM execution.
Once all before-run hooks succeed, after-run hooks run exactly once, including
when subsequent context validation prevents VM entry. They receive the returned
hook context, or the original caller context if the hook context is nil, and the
same execution or validation error. Hook failures are joined with that error.
Pre-canceled admission and failed before-run hooks do not invoke after-run hooks.
These rules apply to both normal and debug sessions.

An after-run failure does not discard successful encoded output. Normal sessions
run these hooks before encoding and result cleanup; the returned error retains
hook, encoding, and cleanup failures without invalidating available output.

## SDK authoring layer

`pkg/sdk` is the supported convenience layer for module and host-value authors.
It provides callback-backed modules, declarative function definitions, typed
runtime binders, context-aware codecs, host-value wrappers, collection views,
and the `sdktest` black-box harness.

SDK helpers should remain thin adapters over public module and runtime contracts.
Core implementation details do not move into `pkg/sdk` merely to make them
accessible across packages.

`sdk.RegisterFunctions` validates a complete registration set before mutating
the namespace. Function names are canonicalized to lowercase and resolved
case-insensitively in FQL. Fixed arities may overload one another, and a
variadic fallback may share the name; duplicate name/arity definitions are
rejected atomically.

Typed binders operate on `runtime.Value` constraints and delegate argument
conversion to runtime helpers. They do not reflect over arbitrary Go functions.
Host-boundary conversion must preserve context, optional configuration,
explicit `None`, unknown-field policy, and root-type validation as configured by
the SDK codec.

## Runtime ownership

Module and SDK code consume runtime value semantics; they do not redefine them.
Host values opt into equality, ordering, iteration, query, dispatch, resource,
or debugger behavior through runtime-owned interfaces.

Argument validation that is specific to an FQL function stays near the function
boundary. Semantics shared with the VM, encoding, debugger, or other functions
belong in `pkg/runtime`. See [Runtime and lifecycle](runtime.md), especially for
hash/equality and resource ownership requirements.

## Standard library

`pkg/stdlib` owns Ferret's built-in functions, namespaces, and immutable
capability-group selections. `stdlib.Full`, `Safe`, `Empty`, and selection
operations determine which groups are registered into a runtime namespace.
Filesystem and network functionality is exposed through the configured host
services rather than bypassing `pkg/fs` or `pkg/net` policy.

Built-in functions should remain small, validate Ferret-facing arguments at the
boundary, preserve argument context in errors, and delegate shared semantics to
runtime helpers. Reusable module contracts do not belong in stdlib, and
stdlib-specific behavior does not belong in `pkg/module`.

Standard-library functions propagate context and follow the canonical
[cancellation rule](../../../AGENTS.md#context-cancellation-in-the-standard-library).
The VM observes ordinary execution cancellation at its safepoints after control
returns; synchronous stdlib work may complete with a canceled context. Downstream
blocking capabilities own cancellation through the propagated context. WAIT owns
its timer/select handling. Contextless filesystem and entropy reads cannot be
interrupted by context; crypto token sampling does not poll between reads.

The Collections group registers global `count`, `count_distinct`, `includes`,
and `reverse`. Their minimum input capabilities, cancellation boundaries, and
ownership contracts are described in [Collection library contracts](../stdlib/collections.md).

The function definitions registered by `stdlib.Full()` are also the source for
the published Ferret Core API artifacts. Structured documentation requirements
and generation checks are described in the
[Core API artifact maintainer guide](../reference/core-api.md).

String operations stay global, while the independent Encoding and Crypto groups
register namespaced serialization and cryptographic operations. Their strict
argument, Unicode, formatting, and migration contracts, together with the full
`path::` namespace registered by the Path group, are described in
[String, encoding, and crypto contracts](../stdlib/strings.md).

The Arrays group registers immutable functions under `arrays::`, explicit
mutation under `arrays::mut::`, and frozen global migration aliases.
Vocabulary, mutation results, set ordering, copy ownership, and
legacy signatures are described in [Array library contracts](../stdlib/arrays.md).

The Objects group registers immutable functions under `object::`, explicit
mutation under `object::mut::`, and seven temporary deprecated global aliases.
The globals delegate to canonical implementations and advertise replacements
through structured API metadata, without compiler or runtime warnings. Shared map
transformations, copy ownership, iterable entry construction, and migration
behavior are described in [Object library contracts](../stdlib/objects.md).

The DateTime group registers canonical `datetime::` functions and deprecated
global migration adapters. Local-calendar precision equality, checked elapsed
differences, calendar arithmetic, and migration changes are described in
[DateTime library contracts](../stdlib/datetime.md).

## Math functions and compatibility

The Math group registers canonical `math::` functions and deprecated global
compatibility functions. Scalar signatures and calculations are shared. The
canonical names `mean`, `variance`, and `stddev` replace global `average`,
`variance_population`, and `stddev_population`; sample forms retain their suffix.
Global registration uses separate documented declarations where sharing a
function would also share its deprecation metadata. Deprecation is API metadata,
not a compiler or runtime warning. Existing exported Go entry points retain
their compatibility contracts.

Canonical collection math accepts native `Int` and `Float` elements, including
mixed lists and non-finite floats. Other elements fail with the argument position
and zero-based element index; there is no coercion. Deprecated globals retain
numeric filtering. Both policies share traversal and calculation helpers.

Shared traversal uses `runtime.List.ForEach`, propagates context, and returns
host errors without successful partial results. An operation's error is returned
directly even if cancellation occurs concurrently; successful traversal and
sorting do not inspect cancellation before returning. Counts come from
traversal rather than `Length`. Sum, mean, and extrema use constant additional
storage. Variance uses Welford's one-pass recurrence and divides by `N` or
`N - 1`; standard deviation takes the corresponding square root. Sources need
not support repeated traversal.

Median and percentile sort native numbers in private snapshots using runtime
comparison with the caller's context. They never copy, sort, index, or
mutate the source. Selected values retain their native type; even medians and
interpolated values are floats.

| Operation | Canonical empty list | Legacy empty list | Legacy nonempty list with no numbers |
| --- | --- | --- | --- |
| Sum | Integer `0` | Integer `0` | Float `0` |
| Mean / average | `NaN` | Float `0` | Float `0` |
| Min, max | `None` | `None` | `None` |
| Median | `NaN` | `None` | `None` |
| Variance, standard deviation, percentile | `NaN` | `NaN` | `NaN` |

Sample statistics need two numbers; population statistics return zero for one
finite number. Non-finite inputs produce `NaN` for variance and standard
deviation. Canonical lists containing non-numbers fail rather than becoming
numeric-empty inputs.

`math::percentile(values, p)` accepts exactly two arguments. `p` must be a finite
`Int` or `Float` in `0..100`, validated even for an empty list. Linear
interpolation uses ascending position `(p / 100) * (N - 1)`, preserving the
selected native value at exact positions. The legacy global keeps integer
percentiles in `1..100`, two/three arguments, nearest rank by default, and the
exact `"interpolation"` method string. Other strings retain nearest-rank fallback.
Its numeric-empty snapshot returns `NaN` before percentile validation; a supplied
method must still be a string.

Range generation belongs to Arrays: `arrays::range(start, end[, step])` shares
its implementation with the deprecated global and exported Go `math.Range`
forwarder. The global remains registered by Math for Math-only embeddings.
`math::range` is not registered. Global `rand` retains its existing coercion,
argument order, and rounded-result behavior as a deprecated adapter. Its draws
use the same Session source as the canonical `random::` namespace. There is no `math::rand`.

The registry supports functions, not namespaced constants, so pi is exposed as
`math::pi()` with a deprecated global `pi()`. Euler's number follows the same
model as `math::e()`, without a global alias or constant registry machinery.

The canonical-only scalar additions `clamp`, `sign`, `trunc`, `cbrt`, `hypot`,
`log1p`, and `expm1` accept native `Int` and `Float` arguments without coercion.
They use the existing argument validators and report the failing argument's
position. `clamp(value, min, max)` rejects NaN bounds and requires `min <= max`.
It uses runtime comparison to select the original `min` when the value is below
it, the original `max` when above it, or the original `value` otherwise. The
result keeps the selected Int or Float type and exact representation without
numeric conversion. Equality preserves `value`, including signed zero. Infinite
bounds are allowed; a NaN input value is returned unchanged after bounds are
validated, even for equal infinite bounds. `sign` returns Int -1, 0, or 1,
treats both signed zeros as zero, accepts infinities, and rejects NaN. Invalid
bounds and a NaN sign raise argument errors.

`trunc` returns Int inputs unchanged and uses Go's `math.Trunc` for Float
inputs, preserving NaN, infinities, and signed zero. `cbrt` and `hypot` use Go's
implementations, including real cube roots of negative values and distance
calculation that avoids unnecessary overflow and underflow. `log1p` computes
`ln(1+x)` and `expm1` computes `exp(x)-1` with improved numerical accuracy near
zero. These four operations return Float and retain Go's domain and non-finite
behavior. None introduces a deprecated global function.

Built-in `COLLECT AGGREGATE` reductions use VM-owned semantics, including in
generic finalization, independently of public math registrations. Only
unqualified `COUNT`, `SUM`, `MIN`, `MAX`, and `AVERAGE` identify those reductions.
Namespaced and other custom selectors retain ordinary function dispatch, so
explicit `math::` selectors remain strict. See [VM execution](runtime.md#vm-execution).

## Random value generation

The independent Random group registers `random::float()` and
`random::float(min, max)`, `random::int(min, max)`, `random::bool()`,
`random::choice(values)`, and `random::shuffle(values)` with fixed arities.
Full and Safe include this group. There are no new global aliases.
Math-only embeddings retain deprecated global `rand`; Random-only embeddings
expose only canonical names.

Float ranges use `[min, max)` and accept only native Int/Float bounds. Bounds
must be finite and ordered using runtime numeric comparison before conversion.
For unequal bounds, ceiling both endpoints to representable Floats preserves
the original numeric interval. Intervals with no representable Float fail.
Interpolation avoids overflow for extreme finite endpoints and corrects rounding
at the excluded upper bound. Equal bounds return ordinary Float conversion of
`min`, including its signed zero, without drawing.

Integer ranges use `[min, max]`, require native Int bounds, and cover the entire
int64 domain through unsigned width arithmetic and unbiased bounded sampling.
Equal bounds return the Int without drawing. Boolean draws consume the same
source. Argument validation does not consume randomness.

Collection operations accept `runtime.List`. Choice calls only `ForEach`, using
one-pass reservoir sampling with O(n) time and O(1) additional storage. It counts
visits independently of callback indices and returns one uniformly selected
original value, or `None` for an empty traversal.

Shuffle calls the source's `New(ctx)` factory and traverses the source once,
appending original value references into the independent destination. The result
preserves the source's implementation family and backend configuration; native
Array inputs produce Arrays. The source needs only `New` and `ForEach` during
this operation. The destination supplies `Append`, `Length`, and `Swap` for
materialization and descending Fisher-Yates. Negative destination lengths fail
with `ErrInvalidOperation`. This uses O(n) appends/swaps and constant algorithm
state; operation costs and destination storage depend on the backend. It does
not force a storage-backed List into an in-memory Array.

Both operations borrow the source and yielded values, preserving value identity
and leaving source contents unchanged. Traversal owns its iterator cleanup.
Factories own failed construction; after successful construction, shuffle owns
the destination until success. Any subsequent operation error, including host
cancellation, closes a closable destination exactly once and joins its cleanup error with the primary
failure, retaining source traversal and iterator cleanup causes. Successful
destinations remain open for normal caller/VM ownership. Borrowed sources and
elements are never explicitly closed. Context reaches each host operation;
stdlib traversal and shuffling do not poll it. Host errors, including cancellation,
are preserved. Successful work and random draws complete even if the context is
canceled while a host operation returns successfully.

Both algorithms use the existing unbiased `rnd.Source.Int64` primitive. Choice
draws for visits 2 through n; shuffle draws only after successful materialization,
for snapshot positions n-1 through 1. Empty and singleton inputs consume no
randomness. Failed choice traversals retain draws already consumed; failures
during shuffle materialization consume none.

Default sources acquire entropy on their first actual draw. Queries that never
draw, including calls with invalid arguments, equal canonical bounds, or empty
and singleton collection inputs, do not initialize the source. Explicitly
seeded sources are ready at construction.

`pkg/rnd` owns the non-cryptographic generator mechanics; the Session owns its
source and context only transports it. The VM's `OpRand` for WAITFOR jitter uses
the same source. Direct stdlib and VM integrations must provide one explicitly
with `rnd.WithContext`; missing sources fail with `runtime.ErrUnexpected`.
See [Session ownership and seeding](runtime.md#session-randomness).

Legacy `rand()` draws in `[0, 1)`. Legacy `rand(x)` uses `min=x/2`, `max=x*2`;
`rand(max, min)` keeps maximum-first ordering. Both ranged forms retain
`floor(u*(max-min+1))+min`, permissive `runtime.ToFloat` conversion, and historical
reversed/non-finite behavior. Each successful legacy call draws once, including
equal bounds. The Go `runtime.RandomDefault`, `Random`, and `Random2` helpers
remain deprecated standalone wrappers; each creates a fresh source and is never
used by Session execution.

Pseudo-random values support automation and reproducible tests. They are not
suitable for passwords, authentication tokens, secrets, keys, or other
security-sensitive use. Secure randomness remains in Crypto.

## Testing

Use package tests for module registration and hook ordering. Exercise SDK
authoring through its public surface and `pkg/sdk/sdktest`. Test built-in
function behavior through FQL whenever practical so argument validation,
registration, runtime semantics, and output are covered together.

Registration tests should include invalid definitions, duplicate arities,
case-insensitive names, atomic failure, and lifecycle cleanup. Host-value tests
should cover every capability they implement, especially comparison,
cancellation, ownership, conversion failures, and debugger inspection.

Changes to stdlib registrations or structured API documentation may require the
focused generator tests documented in the maintainer guide and consideration of
the release flow in [Release automation](../release/process.md).

## Related guides

* [Architecture](overview.md)
* [Runtime and lifecycle](runtime.md)
* [Development workflow](../engineering/workflow.md)
* [Release automation](../release/process.md)
