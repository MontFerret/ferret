# Maintainer Documentation

This directory contains repository-internal documentation for people maintaining
and extending Ferret. It covers architecture, engineering workflows, standard
library contracts, release automation, and maintainer references. Public
user-facing guides are maintained on the [Ferret documentation website](https://ferretlang.org/docs/introduction).

## Architecture

- [Overview](architecture/overview.md) — Execution pipeline, package ownership, and subsystem boundaries.
- [Runtime and lifecycle](architecture/runtime.md) — Runtime values, VM execution, and resource ownership.
- [Debugger](architecture/debugger.md) — Debugger layers, session state, and inspection boundaries.
- [Modules, SDK, and standard library](architecture/modules.md) — Module bootstrap, extension authoring, and built-in function contracts.
- [Universal API](architecture/universal-api.md) — Native-to-Universal adaptation, options, and lifecycle contracts.
- [Engine boundaries](architecture/engine-boundaries.md) — Native exports, internal ownership, and dependency direction.

## Engineering

- [Engineering principles](engineering/principles.md) — Ownership, compatibility, correctness, and documentation policy.
- [Go code style](engineering/code-style.md) — Handwritten Go organization, declarations, comments, and spacing.
- [Testing and performance](engineering/testing.md) — Contract coverage, validation requirements, and benchmark policy.
- [Development workflow](engineering/workflow.md) — Build, generation, testing, lint, and benchmark entry points.

## Standard Library

- [Shared contracts and overview](stdlib/README.md) — Standard-library cancellation policy and category navigation.
- [Arrays](stdlib/arrays.md) — Array operations, mutation, ownership, and legacy compatibility.
- [Collections](stdlib/collections.md) — Collection traversal, construction, and destination ownership.
- [DateTime](stdlib/datetime.md) — Date and time comparison, arithmetic, and migration compatibility.
- [Objects](stdlib/objects.md) — Object transformations, mutation, and shared-value ownership.
- [Strings, encoding, and crypto](stdlib/strings.md) — Text operations, formatting, encoding, and crypto contracts.

## Release

- [Release process](release/process.md) — Dependent notifications, website synchronization, and Core API publication.

## Reference

- [Core API artifacts](reference/core-api.md) — API reference authoring, generation, and versioned publication.
