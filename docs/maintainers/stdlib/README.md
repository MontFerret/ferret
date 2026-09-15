# Standard Library Maintenance

This guide defines contracts shared by Ferret's standard-library functions.
See [Modules, SDK, and standard library](../architecture/modules.md) for
registration, capability groups, and extension boundaries.

## Context cancellation

The VM owns normal execution cancellation. Standard-library functions must not
poll `ctx.Err()` or `ctx.Done()` for ordinary synchronous computation, including
collection traversal and CPU loops. Do not add entry, exit, per-item, phase, or
periodic cancellation checks, or bookkeeping solely to support them.

Propagate context to downstream context-aware APIs. Observe cancellation directly
only when stdlib itself owns blocking, waiting, polling, retry, timer, or similar
work that must be interrupted locally. CPU-bound cancellation is an exception
requiring a demonstrated cancellation-latency problem and benchmark evidence.

## Category guides

* [Arrays](arrays.md): array operations, mutation, ownership, and compatibility.
* [Collections](collections.md): traversal, construction, and destination ownership.
* [DateTime](datetime.md): comparison, arithmetic, and migration compatibility.
* [Objects](objects.md): transformations, mutation, and shared-value ownership.
* [Strings, encoding, and crypto](strings.md): text, formatting, and encoding contracts.
