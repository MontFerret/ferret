# Debugger Architecture

Ferret's debugger is a source-level orchestration layer over retained VM
execution. It must preserve normal execution semantics and remain cheap to
bypass when debugging is disabled.

## Compilation and metadata

`Engine.CompileDebug` uses a compiler configured to emit debug information.
Debug compilation uses effective O0 optimization so source-visible execution is
not rearranged by optimizer passes.

The compiler records source spans, logical debug points, function identities,
program counters, and visible bindings in `bytecode.Program.Metadata`. The
shared source and bytecode packages own those records; debugger code consumes
them rather than reconstructing compiler semantics.

`Plan.NewDebugSession` requires a program containing debug points. Programs
loaded or compiled without debug information cannot create a debug session.
Plan optimization options may still be supplied to `CompileDebug`, but the
compiler forces their effective level to O0 while producing debug metadata.

## Layer boundaries

`pkg/vm` owns retained execution and paused-state access. Its debug execution
surface starts or resumes the VM in continue, step, next, and out modes, reports
stops, accepts pause or termination requests, and exposes frames and runtime
values without moving source-level policy into the dispatch loop.

`pkg/debugger` owns source-level policy:

* binding requested source locations to compiler-emitted debug points;
* serializing commands and translating VM stops into debugger events;
* managing breakpoints, frames, locals, value references, evaluation, and
  presentation limits;
* running embedding lifecycle services and materializing final output.

The root package wires a plan's VM, host services, hooks, source, and output
configuration into `pkg/debugger.Session`. The public plan contract returns
`api/debugger.Session`; a private compatibility bridge translates universal
debugger calls and DTOs to the existing Core debugger while preserving Core
errors and execution behavior. This bridge is temporary until the debugger
types themselves migrate to the universal API, and debugger policy remains in
`pkg/debugger`.

## Session state and concurrency

A debug session retains one execution across commands. `Start` stops at entry;
`Continue`, `Step`, `Next`, and `Out` resume according to the selected VM mode.
Commands are serialized. `Pause` can safely request a stop while a command is
running.

Resume calls use the retained execution context unless a caller supplies an
additional context, in which case both lifetimes are observed. Starting or
resuming invalidates expandable value references from the previous pause.

Breakpoints retain both the requested location and the compiler-emitted point to
which they bind. Binding modes distinguish exact resolution from the supported
next-executable policies. Do not synthesize executable locations outside the
program's debug metadata.

## Values and expression evaluation

Runtime values may opt into debugger inspection through runtime-owned contracts.
The VM exposes a generic inspection view, and `pkg/debugger` applies bounded
formatting and translates it into debugger values and child references.
Debugger code should not accumulate concrete type switches for runtime-owned
semantics.

Evaluation is deliberately conservative and side-effect-free. It reads paused
frame bindings and supports the documented expression subset; calls, mutation,
queries, async/event behavior, and full collection execution remain outside the
debug evaluator.

Formatting limits bound depth, item count, and rendered bytes. Expansion should
be deterministic and must not mutate the live paused value.

## Completion and cleanup

Debug session services apply the same before-run, after-run, encoding,
filesystem, network, logging, and close-hook behavior as normal sessions.
Completion materializes output through the embedding layer.

`Close` requests termination, waits for an active command to leave the retained
execution, invalidates value references, runs after-run handling when needed,
closes retained VM state, and releases embedding-owned resources. Cleanup errors
are aggregated without skipping later cleanup.

Debugger state must not change normal VM execution. Debugger-only work on a
normal execution path must be explicit, measurable, and immediately bypassable.

## Testing and performance

Keep VM retained-execution tests separate from source-level breakpoint,
stepping, evaluation, formatting, and lifecycle tests in `pkg/debugger`. Use
top-level tests for public `DebugSession` composition and output behavior.

Cover invalid state transitions, concurrent pause/close, cancellation, nested
frames, breakpoint resolution, stale value references, formatter bounds,
runtime errors, completion, and cleanup failures. Benchmark debug metadata hooks
and repeated inspection when a change could affect normal dispatch or paused
interaction cost.

## Related guides

* [Architecture](architecture.md)
* [Runtime and lifecycle](runtime.md)
* [Development workflow](workflow.md)
