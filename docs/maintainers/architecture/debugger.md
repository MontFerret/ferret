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

## Internal session composition

The public `Session` composes three private concrete components in
`pkg/debugger`:

* `sessionLifecycle` owns retained execution, embedding services, command and
  lifecycle locks, cancellation contexts, terminal state, hook pairing, output
  materialization, and once-only cleanup.
* `breakpointSet` owns binding, identities, immutable snapshots, the writer
  gate, the VM predicate, and captured hit IDs. It borrows lifecycle admission
  and publication guards so publication remains ordered with native termination.
* `inspector` borrows execution and value access and owns frames, locals,
  expression evaluation, bounded presentation, and expandable references.
  All inspection and reference invalidation run under the lifecycle command lock.

`Session` retains source-event projection and coordinates transitions between
these owners. Run preparation succeeds before it invalidates inspection
references and enters the VM. Close first rejects new work and cancels execution,
then acquires the command lock; Session clears inspection references and captured
hits before lifecycle cleanup releases execution and embedding resources.

Source metadata belongs to Session. Breakpoints and inspection borrow the same
source and debug-point index by pointer; the index is never copied after use.
Components do not retain Session or close each other's resources. Context
validation remains in the Session facade and lifecycle admission guards; the
inspector has no independent command lock or cancellation state.

## Session state and concurrency

A debug session retains one execution across commands. `Start` stops at entry;
`Continue`, `StepIn`, `StepOver`, and `StepOut` resume according to the selected
VM mode. `StepIn` stops at the next debuggable source location and may enter
called functions. `StepOver` stops at the next debuggable location at the same
or shallower call depth. `StepOut` stops at the next debuggable location in a
caller after the current frame exits; at main it runs to completion.

Breakpoints and pause requests take precedence over stepping, including inside
deeper calls. All three stepping operations report the shared `ReasonStep` stop
reason when their stepping condition is reached. Execution and inspection
commands are serialized. `Pause`, breakpoint mutation, and breakpoint listing
can proceed while a command is running.

All commands except Close take a non-nil context. Inspection, including Evaluate
and EvaluateFrame, checks cancellation before waiting for the command mutex and
after acquiring it, before reading frame bindings or evaluating; cancellation does
not interrupt that wait or introduce another inspector lock. A canceled Pause
returns before requesting a stop. Incremental breakpoint add/delete use the
request context for writer admission, construction, and publication checks.
Cancellation before publication prevents changes and ID consumption; cancellation
after publication preserves success and never cancels the debuggee.

Resume calls use the retained execution context unless a caller supplies an
additional context, in which case both lifetimes are observed. Starting or
resuming invalidates expandable value references from the previous pause.

Breakpoints retain both the requested location and the compiler-emitted point to
which they bind. Binding modes distinguish exact resolution from the supported
next-executable policies. Do not synthesize executable locations outside the
program's debug metadata.

### Live breakpoint replacement

`ReplaceBreakpoints(ctx, sourceName, requests)` replaces one source's complete
requested set before execution, while paused, or while running. Empty requests
clear that source. Empty source names select the launched source; other source
names retain unbound records without affecting launched-source breakpoints.
Invalid positions or binding modes fail the operation. Unresolved valid
locations return ordinary unbound results in request order.

The breakpoint set owns an immutable snapshot of requested/resolved records and a
PC-to-hit-IDs index. A cancellable writer gate serializes mutations independently
of `commandMu`. Resolution and construction precede atomic pointer publication;
no partially replaced set is visible. A brief `lifecycleMu` critical section
orders publication with native completion, termination, and closure. Publication
first succeeds; terminal commitment first rejects replacement. Cancellation
observed before publication aborts; cancellation afterward does not undo success.
Lifecycle owns the terminal notification observed by the gate, and Pause/Close
never acquire the writer gate. Lock ordering is command lock or writer gate,
then lifecycle lock, never the reverse.

VM resume receives a synchronous breakpoint predicate instead of a fixed PC
map. At existing debug source points, the predicate loads the current snapshot
and captures matching IDs when deciding a hit. Event conversion copies those
IDs, so removing a breakpoint cannot invalidate an already-decided stop or cause
automatic resume. Added breakpoints affect subsequent checks, never instructions
already passed. Normal VM dispatch retains its existing debugger bypass and
performs no breakpoint publication work when debugging is disabled.

Identity matching uses source name, requested position, and binding mode.
Unchanged requests retain IDs; duplicates match existing IDs in ascending order.
Removed IDs are never reused. Incremental add/delete methods use the same
publication mechanism and preserve their existing admission and validation
contracts. `Breakpoints(ctx)` returns a detached, ID-ordered snapshot and an error.
Nil and canceled contexts return errors; valid-context listing remains available
after close. Inspection reference lifetimes remain tied to execution commands.

Previously, mutation/listing shared the command mutex held throughout resume,
the VM referenced a PC map built at resume, and hit IDs were reconstructed from
the latest mutable maps. These dependencies are replaced together; merely
unlocking commands would not safely publish live changes or preserve hit IDs.

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

When retained execution returns an error directly from `Start` or `Resume` and
its status is terminal, lifecycle commits terminal state and immediately settles
`AfterRun` with the original execution error. Hook failures are joined with that
error. Later `Close` releases resources without repeating the hook or replacing
the run failure with cancellation. Direct admission errors that leave execution
nonterminal do not settle the run. A `DebugStopRuntimeError` event still settles
after-run hooks at the stop while retaining paused, inspectable state; it does
not commit terminal state.

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

* [Architecture](overview.md)
* [Runtime and lifecycle](runtime.md)
* [Development workflow](../engineering/workflow.md)

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
