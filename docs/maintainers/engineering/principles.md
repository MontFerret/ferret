# Engineering Principles

These principles apply to everyone maintaining and extending Ferret v2.

## Ownership and boundaries

Implement behavior in the subsystem that owns its semantics. Consumer packages
must not duplicate an owning package's semantics. Use the
[architecture overview](../architecture/overview.md) for package ownership,
dependency direction, and embedding boundaries.

Prefer the smallest local change that fully solves the problem. New
abstractions, indirection, package moves, and refactors require a concrete
correctness, ownership, or maintainability need. Prefer an existing local
pattern over a new architectural pattern and leave already-correct code alone.
Avoid opportunistic refactoring unrelated to the requested change. Keep
necessary supporting cleanup narrow and explain why it is required.

## Correctness before optimization

Preserve correctness, lifecycle, ownership, and subsystem boundaries before
pursuing cleanup or performance. Optimizations must preserve correctness before
performance; evaluate significant changes under the
[testing and benchmark policy](testing.md#significant-changes-and-benchmarks).

## Public API and compatibility

Treat the top-level `ferret` package, `pkg/engine`, `uapi`, `pkg/module`,
`pkg/runtime`, and `pkg/sdk` as API-sensitive.

* Preserve existing public and language-visible behavior unless a behavior
  change is explicitly required.
* Do not export new symbols unless the change requires an external contract.
* Prefer unexported helpers in the owning package before expanding public API.
* Add contract-focused doc comments to necessary exported symbols.
* Do not move internals into `pkg/sdk` merely to make tests or cross-package
  access easier.
* Do not expose debugger-only behavior through the embedding surface unless the
  change explicitly requires it.
* Backward-incompatible behavior changes require coverage for the former edge
  case and the new expected behavior.

FQL semantics must not change as a side effect of refactoring. Do not infer
compatibility promises from obsolete design notes or the v1 branch. The
[embedding architecture](../architecture/overview.md#embedding-layer) describes
the public facade and native integration boundaries.

## Documentation as part of changes

Documentation is part of the change, not a follow-up activity. Evaluate whether
the implementation changes any documented architecture, ownership boundary,
invariant, workflow, API, behavior, example, or contributor guidance.

Update the relevant documentation in the same change:

* Update `docs/maintainers/**` when repository architecture, subsystem
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

Public documentation is maintained in the website repository. Make necessary
website documentation changes as part of the same change when that source is
available.
