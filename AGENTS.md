# AGENTS.md

## Scope

This file is the canonical operating guide for coding agents working in this
repository. It applies to Ferret v2 only. Do not import assumptions from the
separate v1 branch unless current v2 code and tests establish the same behavior.
Shared technical policy belongs in the maintainer documentation linked below.

## Sources of truth

Use the most direct repository authority for facts that can change:

* `go.mod` owns the module path and minimum Go version.
* `Makefile` owns local targets, tool versions, and command composition.
* `.github/workflows/build.yml` owns the primary CI validation path and tested Go
  versions; other workflow files own their named automation.
* Grammar sources under `pkg/parser/antlr` own FQL syntax. Generated parser files
  are derived output.
* Current code and tests own architecture and behavior. Historical notes, old
  branches, and stale comments are not authoritative.
* `README.md` provides product context, while `CONTRIBUTING.md` describes the
  human contribution process.

When these sources disagree with descriptive documentation, verify the current
implementation and correct the documentation rather than copying stale values.

## Maintainer documentation

Use the [maintainer index](docs/maintainers/README.md) to find the owning guide.
Follow the repository engineering guidance:

* [Engineering principles](docs/maintainers/engineering/principles.md)
* [Go code style](docs/maintainers/engineering/code-style.md)
* [Testing and performance](docs/maintainers/engineering/testing.md)
* [Development workflow](docs/maintainers/engineering/workflow.md)

Before a substantial change, read the relevant subsystem guide; do not load
every guide for unrelated work. Start with
[Architecture and package ownership](docs/maintainers/architecture/overview.md)
to locate the owner. For stdlib work, also read the
[shared standard-library contracts](docs/maintainers/stdlib/README.md) and the
relevant category guide.

## Agent engineering discipline

For every non-trivial change:

1. Identify the owning subsystem and read its maintainer guide.
2. Identify the contract, invariant, and existing behavior being preserved or
   intentionally changed.
3. Choose the smallest implementation that fits the current architecture.
4. Add or update correctness tests for behavior changes.
5. Decide whether the change is significant for performance and benchmark it
   when required.
6. Run the narrowest relevant validation, then broaden in proportion to risk.
7. Evaluate documentation impact and update affected repository and public
   documentation.
8. Review the complete resulting implementation and diff using the mandatory
   self-review below.
9. Fix actual findings and rerun every affected test, check, and benchmark.
10. Report changed behavior, documentation impact, evidence, review results,
    and limitations accurately.

A task is not complete merely because the first implementation compiles or its
tests pass. Do not perform opportunistic refactors unrelated to the requested
change.

## Working with generated sources

Before changing grammar or generator inputs, read and follow the
[parser generation workflow](docs/maintainers/engineering/workflow.md#parser-generation).
Identify the authoritative source inputs and include their derived output in
the complete diff review.

## Validation and evidence

Select validation using the
[testing and performance policy](docs/maintainers/engineering/testing.md) and
[development workflow](docs/maintainers/engineering/workflow.md). Rerun checks
affected by every review-driven fix before reporting completion.

Do not claim that tests, lint, generation, benchmarks, or review passed unless
they actually ran successfully. Report tooling or environment limitations
explicitly, explain skipped validation, and distinguish blocked checks from
product failures. If a required benchmark cannot run, report the limitation
without claiming benchmark validation. For benchmarks that ran, report exact
commands and summarized deltas.

## Mandatory final self-review

After implementation and initial validation, review the complete resulting
change before considering any non-trivial task finished. This is a second-pass
evaluation of the implementation, not a confirmation that tests passed.

### Correctness and lifecycle

Verify the task is completely satisfied. Look for missing cases, incorrect
assumptions, boundary conditions, invalid states, cancellation and concurrency
errors, cleanup failures, resource leaks, ownership mistakes, and incorrect
error propagation. Confirm public and FQL-visible behavior matches the intended
contract. Check idiomatic Go error wrapping, context propagation, synchronization,
and lifecycle management. For bug fixes, prefer a regression test that fails
without the fix.

### Architecture and API

Verify responsibilities remain in the correct package, type, and layer. Check
dependency direction, compile-time/runtime separation, runtime-owned semantics,
public API necessity, and compatibility against the
[architecture guides](docs/maintainers/README.md#architecture) and
[engineering principles](docs/maintainers/engineering/principles.md). Reject
duplicated semantics, leaked implementation details, misplaced behavior, and
abstractions at the wrong level.

### Clarity and organization

Look for unnecessary complexity, duplication, nesting, misleading names, dead
branches, debugging artifacts, and comments about abandoned approaches. Check
the [Go code style guide](docs/maintainers/engineering/code-style.md). Avoid both
overloaded files and unnecessary fragmentation into excessive helpers or
abstractions.

### Tests and performance

Review whether tests cover meaningful positive, negative, boundary, error,
cleanup, and cancellation cases without merely mirroring implementation. Check
assertion strength and unnecessary brittleness. Use the
[testing and performance guide](docs/maintainers/engineering/testing.md) to assess
coverage and validation. For significant changes, inspect allocations, repeated
work, materialization, synchronization, and benchmark comparability.

### Scope and complete diff

Inspect the complete final diff, not only individual files. Verify that:

* every changed line belongs to the request or a necessary supporting change;
* no temporary code, accidental API or behavior change, or unrelated refactor
  remains;
* generated files changed only because their source inputs changed;
* tests describe intended behavior;
* comments describe current contracts and invariants;
* file, type, function, and package boundaries remain coherent;
* resource ownership and lifecycle behavior remain correct;
* affected repository and public documentation has been updated, or any
  unavailable external documentation dependency is explicitly reported;
* the final implementation is the smallest coherent solution.

When review finds a correctness, architecture, ownership, lifecycle, API,
organization, performance, documentation, or meaningful coverage problem, fix
it and repeat affected validation. Minor optional style preferences do not
justify churn.

Do not use self-review to expand scope through speculative refactoring,
unrelated cleanup, unrelated API redesign, broad package moves, or FQL semantic
changes outside the task.

## Change and reporting discipline

Apply the [engineering principles](docs/maintainers/engineering/principles.md)
within the requested scope. Explain why any necessary supporting cleanup is
required. Call out every intentional FQL semantic change in the final report.

The final report for a non-trivial change must state:

* owning subsystem and files changed;
* behavior and invariants preserved or intentionally changed;
* tests and benchmarks added or updated;
* validation and benchmark commands actually run;
* documentation updated, or documentation impact explicitly evaluated as none;
* completion of the mandatory final self-review;
* material findings corrected during review;
* remaining limitations or skipped validation.

## Documentation synchronization

Before completing every non-trivial task, evaluate documentation impact under
the [shared documentation policy](docs/maintainers/engineering/principles.md#documentation-as-part-of-changes)
and complete the required updates in the same task. Report which documentation
changed, or explicitly state why no documentation update was needed.

When the website repository or another required documentation source is
available, make the necessary changes as part of the task. If it is
unavailable, identify the exact required follow-up in the final report rather
than silently leaving known documentation stale.

## Response style

Keep responses practical and easy to scan. Use short sections, focused bullets,
and code blocks only for code, commands, or configuration. Explain why a change
is needed before how it works, summarize each changed file's responsibility, and
avoid repeating the same context.
