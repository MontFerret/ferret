# Debugger Architecture

Ferret's debugger is a source-level orchestration layer over retained VM
execution. It must preserve normal execution semantics and remain cheap to
bypass when debugging is disabled.

## Compilation and metadata

`Engine.CompileDebug` uses a compiler configured to emit debug information.
Debug compilation uses `OptimizationNone` so source-visible execution is
not rearranged by optimizer passes.

The compiler records source spans, logical debug points, function identities,
program counters, and visible bindings in `bytecode.Program.Metadata`. The
shared source and bytecode packages own those records; debugger code consumes
them rather than reconstructing compiler semantics.

`Plan.NewDebugSession` requires a program containing debug points. Programs
loaded or compiled without debug information cannot create a debug session.

## Layer boundaries

`pkg/vm` owns retained execution and paused-state access. Its debug execution
surface starts or resumes the VM in `DebugResumeContinue`, `DebugResumeStepIn`,
`DebugResumeStepOver`, and `DebugResumeStepOut` modes, reports stops, accepts pause
or termination requests, and exposes frames and runtime values without moving
source-level policy into the dispatch loop.

`pkg/debugger` owns source-level policy:

* binding requested source locations to compiler-emitted debug points;
* serializing commands and translating VM stops into debugger events;
* managing breakpoints, frames, locals, value references, evaluation, and
  presentation limits;
* running embedding lifecycle services and materializing final output.

`pkg/engine` validates admission, debug metadata, and native session options.
Its internal session constructor assembles a dedicated VM, host services, hooks,
source, and output configuration into `pkg/debugger.Session`, using the same
acquisition owner as ordinary execution. Root `ferret` aliases the supported
debugger types through the curated embedding API; source-level debugger policy
stays in `pkg/debugger`.

## Session state and concurrency

A debug session retains one execution across commands. `Start` stops at entry;
`Continue`, `StepIn`, `StepOver`, and `StepOut` resume according to the selected
VM mode. `StepIn` stops at the next debuggable source location and may enter
called functions. `StepOver` stops at the next debuggable location at the same
or shallower call depth. `StepOut` stops at the next debuggable location in a
caller after the current frame exits; at main it runs to completion.

Breakpoints and pause requests take precedence over stepping, including inside
deeper calls. All three stepping operations report the shared `ReasonStep` stop
reason when their stepping condition is reached. Commands are serialized.
`Pause` can safely request a stop while a command is running.

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

The embedding services receive their own host-resource manager from session
construction. It borrows Engine services and owns a session filesystem override,
if configured. Service closure runs close hooks, closes owned host resources,
and releases the limiter permit. Retained VM execution remains owned and closed
by the debugger session rather than by this manager. Native debug services bind
the concrete limiter and resource manager through a constructor and keep them
private. They never receive an ordinary `session.Execution`, pooled VM, or
VM-return callback. The internal acquisition owner retains rollback until
construction succeeds; the debugger serializes service use and invokes service
closure once.

When all before-run hooks succeed but context validation prevents VM entry,
`Start` settles that attempt's after-run hooks immediately. The session remains
new and accepts a later valid `Start`; `Close` does not repeat the aborted
attempt's hooks or consume the later run's hook obligation.

`Close` requests termination, waits for an active command to leave the retained
execution, invalidates value references, runs after-run handling when needed,
closes retained VM state, and releases embedding-owned resources. Cleanup errors
are aggregated without skipping later cleanup.

Termination events include both the execution cause and retained-state cleanup
failures. The debug execution caches cleanup failures for repeated and concurrent
`Close` calls. Cancellation itself is not cached as a cleanup failure, so a clean
cancellation does not make later `Close` calls fail. Ordinary runtime errors
remain inspectable until the next resume or close releases their retained state.
If `Close` wins while the VM returns a paused or runtime-error stop, the session
drains that state before reporting termination and preserves concurrent failures.

Debugger state must not change normal VM execution. Debugger-only work on a
normal execution path must be explicit, measurable, and immediately bypassable.

## Testing and performance

Keep VM retained-execution tests separate from source-level breakpoint,
stepping, evaluation, formatting, and lifecycle tests in `pkg/debugger`. Use
`pkg/engine` tests for native `DebugSession` composition and output behavior,
and root tests for supported façade usage.

Cover invalid state transitions, concurrent pause/close, cancellation, nested
frames, breakpoint resolution, stale value references, formatter bounds,
runtime errors, completion, and cleanup failures. Benchmark debug metadata hooks
and repeated inspection when a change could affect normal dispatch or paused
interaction cost.

## Related guides

* [Architecture](architecture.md)
* [Runtime and lifecycle](runtime.md)
* [Development workflow](workflow.md)

## Portable debugger boundary

Coordinates, values, variables, frames, breakpoints, reasons, and events alias
the Universal API types. Native source text/indexing and compiler debug tables
remain native. Table identifiers are validated natively and converted only when
constructing portable debugger values. The Universal API debugger package owns
`NoFunction` for the top-level body and the positive value-reference convention;
references are usable only in their paused state.

Native source positions use one-based lines and byte columns; spans are
zero-based half-open byte offsets. Source names can be anonymous or non-path
identities. The compiler converts ANTLR rune offsets into byte spans when
publishing diagnostics and program metadata, including debug points and
call-argument spans. Analysis already publishes byte spans. Breakpoint creation accepts a source location and optional binding
mode; unknown modes are rejected. Command cancellation reports termination while
preserving cancellation identity. Ordinary runtime failures remain inspectable
runtime-error stops. Completion retains an event and available output even when
an after-run hook or later result cleanup fails.
