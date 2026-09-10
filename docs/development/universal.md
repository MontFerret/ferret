# Universal API adapter

`pkg/universal` exposes Ferret Native through `github.com/MontFerret/api` by
translating portable source, options, results, and errors into Native engine
operations. Use `New` to create and own a Native engine configured with Native
options:

```go
portable, err := universal.New(engine.WithParam("value", 42))
if err != nil {
    return err
}
defer portable.Close()

output, err := portable.Run(ctx, api.NewAnonymousSource("RETURN @value"))
// output can contain encoded data even when err reports a later cleanup failure.
```

Import Native and the adapter from `github.com/MontFerret/ferret/v2/pkg/engine`
and `github.com/MontFerret/ferret/v2/pkg/universal`. The root `ferret.Engine`
alias is also accepted. `New(opts ...engine.Option)` returns `(*Runtime, error)`;
Native owns construction and rollback, and the adapter projects construction
errors. Nil options are handled by Native and skipped.

Use `Wrap` to borrow a caller-supplied engine without acquiring resources:

```go
native, err := engine.New()
if err != nil {
    return err
}
defer native.Close()

var portable api.Runtime = universal.Wrap(native)
```

Configure Native modules, codecs, host services, and engine defaults before
wrapping. `Wrap(native *engine.Engine)` returns `*Runtime` and panics when native
is nil. The caller retains responsibility for closing the Native engine.

## Options and representation differences

| Universal option | Native translation |
| --- | --- |
| `api.WithOptimizationLevel` | Explicit None/Basic/Full mapping to `engine.WithPlanOptimizationLevel`; Aggressive and unknown levels fail. |
| `api.WithParam`, `api.WithParams` | `engine.WithSessionParam`, `engine.WithSessionParams`; Native owns host-value conversion. |
| `api.WithOutputContentType` | `engine.WithOutputContentType`; Native validates the value and resolves the codec. |
| `api.WithFSRoot` | `engine.WithSessionFSRoot`; Native validates, creates, and owns the override, preserving the engine's read-only policy. |

Omitted optimization inherits the engine default. Debug compilation accepts
omission or None. Plan setters reject unsupported levels immediately, without
changing prior settings or engine defaults. Native and Universal options remain
intentionally separate types.

Non-nil portable option callbacks run once in order. Session setters only queue
Native options and return nil. Returned callback failures are joined and stop
delegation. Otherwise Native applies the queued options and joins validation
failures before acquiring session resources. `Runtime.Run` delegates to
`Engine.Run`, so compilation and its hooks can precede Native session validation.
Native also closes the temporary plan if session validation fails.

Later values override earlier values, and parameter maps merge. Maps are
converted when Native applies the options, after all portable callbacks finish;
mutations during translation are visible to conversion. Even an initially empty
map is queued. Native rejects an empty single-parameter name or nil host value;
map values retain Native conversion semantics, including nil becoming `none`.
Nil and empty maps are no-ops when applied. Runtime-specific extension options
must reject incompatible targets; their errors are propagated.

Source conversion is `source.New(src.Name, src.Content)`. Coordinates and debugger
values already share portable types. Native output differs only by pointer versus
value: nil becomes zero output, and dereferencing does not copy encoded bytes.
Output belongs to the caller and survives cleanup. Debugger commands return Native
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

Caller contexts pass through unchanged. The adapter checks nil or already-canceled
contexts before portable callbacks and retains cancellation alongside callback
errors. Closed-engine and closed-plan rejection occurs in Native, after portable
option translation. The adapter adds no closing contexts, waiting, or tracking.

Use caller contexts to stop work, settle it, then close sessions, plans, and the
owning runtime or borrowed Native engine. Native engine close does not wait for
admitted work or close descendants. Native plan close rejects new sessions and
wakes capacity waiters without waiting for constructors. Ordinary sessions support sequential reuse;
settle Run before Close. Native debug Close terminates and settles commands.
See [Runtime and lifecycle](runtime.md).

## API release alignment

The pinned `api v1.0.0-alpha.14` documentation requires stronger parent-close
coordination than this adapter implements. The corresponding API contract update
permits borrowed runtime no-op close, Native parent-close behavior, and deferred
option validation. Publish that aligned contract and update Ferret's root and
API-reference-tool dependency pins before merging this refactor. Do not commit
local module replacements as a substitute.

ferretd migration remains separate. Its composition must retain Native engine
ownership when replacing its local adapter; daemon shutdown policy stays there.

## Validation

Contract and delegation tests live in `pkg/universal`, including external-package
embedding examples for both constructors. Native/Universal benchmarks cover
construction and close, compilation, ordinary/debug session creation, reusable
execution, and convenience Run. Run focused tests before the repository gates in
[Development workflow](workflow.md). Native APIs
and FQL semantics are unchanged.
