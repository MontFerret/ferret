# Universal API adapter

Root-level `uapi` is Ferret's official adapter for `github.com/MontFerret/api`.
It translates portable source, options, results, and errors into Native engine
operations while keeping Native and Universal APIs independently evolvable.
Use `uapi.New` to create and own a Native engine configured with Native options:

```go
import (
    "fmt"

    "github.com/MontFerret/api"
    "github.com/MontFerret/ferret/v2"
    "github.com/MontFerret/ferret/v2/uapi"
)

portable, err := uapi.New(ferret.WithParam("value", 42))
if err != nil {
    return err
}
defer portable.Close()

output, err := portable.Run(ctx, api.NewAnonymousSource("RETURN @value"))
if output != nil {
    fmt.Printf("%s: %s\n", output.ContentType, output.Content)
}
if err != nil {
    return err
}
```

Use the root `ferret` package for ordinary Native embedding and configuration.
`uapi.New` accepts its Native options and returns `(*uapi.Runtime, error)`;
calling `uapi.New()` uses Native defaults. The adapter internally depends on
`pkg/engine`; root `ferret.Engine` and `ferret.Option` aliases preserve type
identity. Native owns construction and rollback, and the adapter projects
construction errors. Nil options are handled by Native and skipped.

Use `uapi.Wrap` to borrow a caller-supplied engine without acquiring resources:

```go
native, err := ferret.New()
if err != nil {
    return err
}
defer native.Close()

var portable api.Runtime = uapi.Wrap(native)
```

Configure Native modules, codecs, host services, and engine defaults before
wrapping. `uapi.Wrap(native)` returns `*uapi.Runtime` and panics when native is
nil. The caller retains responsibility for closing the Native engine. Both
constructors remain in `uapi`; the root façade exposes the Native API.

## Options and representation differences

| Universal option | Native translation |
| --- | --- |
| `api.WithOptimizationLevel` | Explicit None/Basic/Full mapping to `engine.WithPlanOptimizationLevel`; Aggressive and unknown levels fail. |
| `api.WithParam`, `api.WithParams` | `engine.WithSessionParam`, `engine.WithSessionParams`; Native owns host-value conversion. |
| `api.WithOutputContentType` | `engine.WithOutputContentType`; Native validates non-blank input when applying options and resolves the codec during result encoding. |
| `api.WithFSRoot` | `engine.WithSessionFSRoot`; Native validates, creates, and owns the override, preserving the engine's read-only policy. |

Omitted optimization inherits the engine default. Plan setters map supported
levels to queued Native options and reject unmappable levels immediately.
Native validates the queued levels against the compilation mode: debug compilation
accepts omission or None, and rejects Basic/Full even if followed by None.
Native and Universal options remain intentionally separate types.

Non-nil portable option callbacks run once in order. Session setters only queue
Native options and return nil. Returned callback failures are joined and stop
delegation. Otherwise Native applies the queued options and joins option-application
failures before acquiring session resources. `Runtime.Run` delegates to
`Engine.Run`, so compilation and its hooks can precede Native session validation.
Native also closes the temporary plan if session validation fails.

The portable contract permits runtime-specific validation at the point of use;
it does not require all configuration to be checked before compilation, resource
acquisition, or query execution. Native resolves output codecs when encoding the
result, after the query has run. An unavailable codec therefore returns an
operation error after callbacks and session creation succeeded and query side
effects may have occurred. Native still closes the result, and temporary or
caller-owned sessions retain their usual cleanup responsibilities.

Later values override earlier values, and parameter maps merge. Maps are
converted when Native applies the options, after all portable callbacks finish;
mutations during translation are visible to conversion. Even an initially empty
map is queued. Native rejects an empty single-parameter name or nil host value;
map values retain Native conversion semantics, including nil becoming `none`.
Nil and empty maps are no-ops when applied. Runtime-specific extension options
must reject incompatible targets; their errors are propagated.

Source conversion is `source.New(src.Name, src.Content)`. Coordinates and debugger
values already share portable types. Native and Universal execution return the
same `*api.Output` type, so the adapter preserves the Native pointer directly.
Nil output means no output was produced. Non-nil output with a nil error indicates
success, including empty output; a non-nil output may also accompany hook or
cleanup errors. A zero-valued `Output` is not an absence sentinel. Inspect output
independently of the error. Output belongs to the caller and survives cleanup.
Debugger commands return Native
snapshots with diagnostic projection applied to command and event errors.

Diagnostic projection preserves kind, message, hint, note, each diagnostic's source,
annotation order, and byte spans. Error-tree traversal preserves branch order and
distinct diagnostics while projecting repeated Native diagnostic pointers once.
`errors.As` exposes `api/diagnostics.Diagnostics`; standard Go error traversal
still reaches Native causes. Causes remain outside portable diagnostic fields.
Projection preserves Native error messages and available output or completion events.

## Ownership and delegation

For `New`, `Runtime.Close` delegates to the owned Native engine, which releases
its resources and rejects subsequent runtime operations. For `Wrap`, Close is a
no-op that leaves the adapter and borrowed engine usable. Multiple adapters may
borrow one engine independently; external engine closure is observed by all.

Plan and session close calls also delegate to Native. Native owns idempotence,
concurrent closure, and the completed cleanup result. Repeated projections need
not have identical pointers. Runtime stores only its Native pointer and an
immutable ownership flag.

Portable option callbacks run independently of the operation context. If they
fail, only their returned errors are projected; cancellation is included only
when a callback returns it. Otherwise the Native operation receives the original
caller context and owns context validation, cancellation, and closed-object
rejection. Even nil, canceled, or expired contexts do not skip portable callbacks.
The flow remains portable input, translation, Native operation, result projection.

Use caller contexts to stop work, settle it, then close sessions, plans, and the
owning runtime or borrowed Native engine. Native engine close does not wait for
operations already started or close descendants. Native plan close rejects new
sessions and wakes capacity waiters without waiting for constructors. These are
Native delegation guarantees. The portable contract requires cleanup of owned
resources and rejection of new descendants after plan closure, without requiring
waiting for constructors or coordinating descendant cleanup. Parent close does
not implicitly cancel caller-owned work; returned sessions and debug sessions
retain their own lifecycle. Callers coordinate cleanup when descendants use
parent-owned resources. Ordinary sessions support sequential reuse;
settle Run before Close. Native debug Close terminates and settles commands.
See [Runtime and lifecycle](runtime.md).

## API release alignment

The root module and API-reference tool pin `api v1.0.0-alpha.16`. This published
contract permits borrowed runtime no-op close, Native parent-close behavior,
deferred option validation at the point of use (including output encoding), and
portable translation before operation-context checks. `Runtime.Run` and
`Session.Run` return `(*Output, error)` so output presence survives adaptation.

Validate the root module and both tool modules with `GOWORK=off` to check their
independent module metadata against published dependencies. When root dependency
versions change, synchronize the API-reference tool's module files as well.

ferretd migration remains separate. Its composition must retain Native engine
ownership when replacing its local adapter; daemon shutdown policy stays there.

## Validation

Contract and delegation tests live in `uapi`, including external-package
embedding examples for both constructors. Native/Universal benchmarks cover
construction and close, compilation, ordinary/debug session creation, reusable
execution, and convenience Run. Run focused tests before the repository gates in
[Development workflow](workflow.md). Native APIs
and FQL semantics are unchanged.
