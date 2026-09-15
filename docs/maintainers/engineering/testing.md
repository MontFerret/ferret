# Testing and Performance

This guide defines Ferret's shared testing, validation, and benchmark policy.
See [Development workflow](workflow.md) for commands, tool requirements, test
locations, and CI mechanics.

## Contract ownership

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

## Validation

Run the narrowest relevant command first, then broaden validation in proportion
to risk. Use `Makefile` and `.github/workflows/build.yml` to select current
repository-level commands. Changes that cross packages or match CI coverage
should finish with the corresponding repository-level test target.

For generator-input changes, follow the
[parser generation workflow](workflow.md#parser-generation). Re-run affected
validation after every review-driven change.

Run `make fmt` when handwritten Go formatting is affected and inspect its broad
rewrite surface. Run `make lint` for lint-sensitive code or public API changes
when the required tools are available.

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
2. Run the same benchmark command after implementation with the same environment.
3. Compare `ns/op`, `B/op`, and `allocs/op` where available.
4. Investigate meaningful regressions before completing the task.
5. Add a benchmark when no relevant one covers the changed hot path.
6. Record the exact commands and summarized delta.

Never trade clear correctness or maintainability for a speculative
micro-optimization. See [Engineering principles](principles.md#correctness-before-optimization)
for the broader priority of correctness, lifecycle, and ownership.
