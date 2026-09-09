# Universal API adapter

`pkg/universal` exposes a Native Ferret engine through
`github.com/MontFerret/api`. Use it when an integration consumes the portable
Runtime, Plan, Session, or debugger contracts. Configure Native modules, codecs,
host services, and engine defaults before wrapping the engine.

```go
native, err := engine.New()
if err != nil {
    return err
}
defer native.Close()

var portable api.Runtime = universal.New(native)
defer portable.Close()

output, err := portable.Run(ctx, api.NewAnonymousSource("RETURN @value"),
    api.WithParam("value", 42))
// output can contain encoded data even when err reports a later cleanup failure.
```

Import the Native engine and adapter from
`github.com/MontFerret/ferret/v2/pkg/engine` and
`github.com/MontFerret/ferret/v2/pkg/universal`. The root `ferret.Engine` alias is
also accepted. `New` returns `*universal.Runtime` and panics on a nil engine.
It does not acquire or take ownership of the supplied engine's resources.

## Options and conversions

| Universal option | Native translation |
| --- | --- |
| `api.WithOptimizationLevel` | Explicit None/Basic/Full mapping to `engine.WithPlanOptimizationLevel`; Aggressive and unknown levels are rejected. |
| `api.WithParam`, `api.WithParams` | Immediate runtime-owned host-value conversion, followed by Native runtime-parameter options. |
| `api.WithOutputContentType` | `engine.WithOutputContentType`; Native resolves the codec. |
| `api.WithFSRoot` | `engine.WithSessionFSRoot`; Native creates and owns the session override and preserves the engine's read-only policy. |

Omitted optimization inherits the engine default. Debug compilation accepts
omission or None. Per-plan overrides never change engine defaults.
Non-nil options run once in order; validation failures are joined before Native
calls acquire resources. Setters validate immediately, and rejected setters
leave earlier settings intact. Later values override earlier values; parameter
maps merge and are converted when applied. Single parameters reject empty names
and nil host values; map parameters retain Native conversion semantics, including
nil values becoming the runtime's none value. Nil and empty maps are no-ops.
Unsupported runtime extension options must reject incompatible targets; their
errors are propagated. Native options are intentionally distinct types.

Portable sources become Native indexed sources. Coordinates and debugger values
already share portable types. Native encoded output belongs to the caller, so
the adapter converts its pointer to a value without copying its bytes. Outputs
survive session and parent cleanup. Debugger commands return Native snapshots
with diagnostic projection applied to command errors and event errors.

The diagnostic adapter walks each immediate error-tree branch in order. It
deduplicates repeated Native diagnostic pointers, retains distinct diagnostics,
and projects each diagnostic's own source and byte spans. `errors.As` can obtain
`api/diagnostics.Diagnostics`; `errors.Is` and `errors.As` still reach the original
Native causes. Execution, cancellation, hook, and cleanup errors remain distinct
causes. Available output and completion events are retained alongside errors.

## Lifetime and concurrency

Callers close sessions, plans, the adapter, and finally the Native engine.
Multiple adapters can borrow one engine independently. Parent closure does not
close published descendants; the adapter has no descendant registry.

`Runtime.Close` rejects new calls and waits for admitted Compile, CompileDebug,
and Run calls. It does not cancel them or close the engine. Callers use their
contexts to stop active work. `Runtime.Run` delegates execution and temporary
session/plan cleanup to Native.

Adapter plan close rejects new constructors and cancels outstanding constructor
contexts so capacity waiters can finish. It waits before running Native plan
cleanup. Constructor contexts are detached on return and never become the
lifetime of published ordinary or debug sessions. Successfully returned sessions
remain caller-owned even when plan close overlaps their creation.

Ordinary sessions support sequential reuse. Callers settle Run before Close;
debug Close terminates and settles active commands. Close calls are idempotent
and concurrency-safe, retaining the completed cleanup error projection. No
admission bookkeeping lock spans user options, hooks, Native calls, or cleanup.
Synchronous close from an option or hook inside an admitted operation on that
same parent, and recursive close from a close hook, are unsupported.

## ferretd follow-up

ferretd migration remains separate. Its existing `internal/ferretapi.New`
transfers engine cleanup responsibility to its runtime; `universal.New` borrows
instead, so daemon composition must explicitly retain and close the Native
engine. The local adapter's direct Universal-to-Native option forwarding is
obsolete for current Native option types. The official adapter also supplies
the portable admission/close coordination and avoids redundant copies of
Native-owned output and debugger snapshots.

Migration must align the published Ferret dependency and generated standard
library catalog, replace the local adapter at the composition boundary, and
retain daemon workspace, transport, debug leases, and shutdown orchestration
inside ferretd. Neither runtime policy nor that migration belongs in this package.

## Validation

Contract and race tests live in `pkg/universal`, including an external-package
embedding example. Native/Universal benchmarks cover compilation, session
creation with and without translated options, and reusable execution. Run
focused tests before the repository-level gates described in
[Development workflow](workflow.md). Native APIs and FQL semantics are unchanged.
