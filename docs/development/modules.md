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
Context propagation, error joining, and cleanup behavior must remain consistent
between normal and debug sessions.

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

The function definitions registered by `stdlib.Full()` are also the source for
the published Ferret Core API artifacts. Structured documentation requirements
and generation checks are described in the
[Core API artifact maintainer guide](../maintainers/core-api-reference.md).

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
the release flow in [Release automation](release.md).

## Related guides

* [Architecture](architecture.md)
* [Runtime and lifecycle](runtime.md)
* [Development workflow](workflow.md)
* [Release automation](release.md)
