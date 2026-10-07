# Consumable output migration

Core uses Universal API `v1.0.0-alpha.21`, published at
`ceb51144c51034ede3b766985086990b76059274`. Native Engine.Run and Session.Run
return `(encoding.Output, error)`; the root facade and UAPI share the exact
portable Output, Content, Metadata, Consumer, and output-error identities.

## Implementation ownership

- `pkg/engine/output.go` and `output_errors.go` own the single native output
  lifecycle, context coordination, ownership transfer, and terminal errors.
- `pkg/engine/session.go`, `engine.go`, and `internal/session/output.go` retain
  native orchestration and VM/result responsibilities, with separate operation
  and cleanup outcomes and early convenience cleanup.
- `pkg/encoding/codec.go`, `pkg/encoding/output.go`, and root `output.go` define
  the encoder ownership contract and exact portable aliases and sentinels.
- `uapi/output.go`, `runtime.go`, and `session.go` delegate native output methods
  and project diagnostics. Debugger services and lifecycle now return detached
  Content while preserving debugger command and event behavior.
- Compatibility adapters, SDK harnesses/examples, the internal CLI, native/root/
  UAPI/security tests, and benchmarks explicitly settle execution outputs.
  Root and API-reference module files pin alpha.21; the other tool module is
  unchanged. README, GoDoc, examples, and architecture/lifecycle/UAPI/debugger
  guides document the new boundary.

## Execution and errors

Run executes and encodes eagerly. Compilation, option application, session
construction, before-run failure, and context validation before VM entry are
immediate errors with no handle. Successful before-run hooks are still paired
when their replacement context prevents entry. Once VM entry is admitted, Run
returns a usable handle and nil error. Consume or Collect reports execution,
encoding, after-run-hook, consumer, cancellation, and cleanup failures.

Ordinary after-run hooks remain before encoding and receive the primary execution
error. Materialization closes the VM result immediately. Convenience Run closes
its temporary session, then its plan, before publishing the output, retaining
cleanup failures separately. The byte-backed handle owns no VM, plan, session,
or permit. A caller-created session remains open and supports sequential reuse
after its preceding output settles. Existing parent-close behavior is unchanged.

Close returns the recorded cleanup outcome, not execution success. A deferred
Close after successful consumption does not create an error. Error traversal
retains native diagnostic, context, consumer, and joined cleanup causes; UAPI
projects diagnostics from all three output methods without changing payloads.

## Consumption, contexts, and closure

Consume and Collect are alternative one-shot operations. Admission validates a
non-nil context, then a non-nil consumer for Consume, then consumption and
invocation cancellation, before claiming the output. Rejections leave active
consumption unaffected. Valid competitors match ErrOutputInUse; once stopping or
finalization begins they match ErrOutputClosed. Use errors.Is with these shared
sentinels. An irrevocably canceled invocation also preserves already-known eager
operation and cleanup failures without claiming the output, whether consumption
uses the original Run context, a live context, or another canceled context. When
both contexts are canceled, their causes are preserved. Cancellation of only the
consumption context leaves the output available for retry within the invocation's
lifetime.

Keep the original Run context alive through settlement. The consumption context
may shorten that lifetime, never extend it. Consumers receive an effective
context observing both cancellation sources and the earliest deadline. Hook
replacement contexts remain execution-only. Coordination uses cancellation
registrations only where the consumption parent does not already propagate the
invocation; terminal paths stop registrations and cancel derived contexts.

Native Consume supplies the existing encoded buffer synchronously as borrowed,
read-only bytes. Empty present content invokes an empty callback; absent content
invokes none. Copy bytes retained beyond the callback. Public chunk boundaries
remain unspecified. Consumer errors stop delivery, and panic unwinding finalizes
without swallowing the panic. Collect directly transfers owned bytes.

Close abandons unread content without draining it. During consumption it marks
stopping, requests cancellation, then waits for the callback and finalization.
Concurrent closers converge on one settlement. Borrowed storage remains valid
until the callback returns. Explicit closure winning before outcome commitment
adds ErrOutputClosed; later closure only waits. Automatic finalization does not
relabel a failure as explicit closure. Arbitrary callback code cannot be forcibly
interrupted. Synchronous Close from the output's own consumer is unsupported;
return an error to stop instead. Metadata reads and reentrant consumption checks
remain safe from callbacks.

## Content, encoding, and debugger retention

Nil *Content means absent; any non-nil pointer means present, including nil or
empty Data. Successful encoding records the selected codec's actual media type,
the exact byte length, and LengthKnown=true. Non-nil encoder bytes returned with
an error remain available with unknown total length; nil bytes with that error
mean absence. Collected metadata remains the descriptor established before Run
returns. Output.Metadata is immutable and available after closure.

Encoder.Encode transfers ownership of every returned slice, including on error.
External custom encoders must detach reusable storage and borrowed runtime bytes
inside Encode while access remains exclusive, before unlocking or recycling.
Configured encoders have the same requirement. Built-in JSON and MessagePack
already return buffers independent of their encoder and VM state. Core does not
blanket-clone data in materialization, consumption, collection, or adaptation.

Debugger services materialize *Content directly and retain their independent
result cleanup and hook ordering. Completed events can carry content alongside
cleanup errors. They never contain a live, one-shot output. Returned bytes remain
valid after result cleanup, session closure, and subsequent executions.

Serialize detached Content, not Output. JSON now has a nested `metadata` object
and a `data` field using standard byte-slice encoding. Nil and empty data encode
as null and an empty string; the containing non-nil Content remains present.

## Validation and performance evidence

Lifecycle tests use the native implementation, with real facade and UAPI
execution coverage. Channel-controlled concurrent-close tests prove cancellation,
callback settlement, borrowed-buffer validity, and exactly-once native cleanup.
Tests also cover admission, consumer panic, context causes and deadlines, joined
diagnostics, codec detachment, content serialization, and debugger completion.

Benchmark workloads compare the original buffered API against terminal collection,
discard consumption, and consuming writes for 32-byte, 64 KiB, and 1 MiB strings.
They cover convenience runs, reused sessions, Native and UAPI, plus output-only
handle overhead. Query execution and output-owned cleanup are included. The
writer reuses a bytes.Buffer; its payload copy belongs to that destination.

The baseline is Core `77e403e667a96338250944c86836a5f9d4b2a315` with API alpha.20,
plus the comparable benchmark fixtures. Environment: Go 1.26.5, Darwin/arm64,
Apple M2 Max, GOMAXPROCS=12. Both revisions use:

```sh
GOCACHE=/private/tmp/ferret-output-plan-go-cache go test ./pkg/engine ./uapi \
  -run '^$' -bench 'Benchmark(BufferedOutput|OutputExecution)$' \
  -benchmem -benchtime=200ms -count=5
benchstat baseline.txt head.txt
```

The final five-sample collection medians are below. Byte and allocation deltas
are per operation; large-buffer measurements include small amortization noise.
`Reusable` retains its caller-owned session between iterations; `Convenience`
includes compilation and temporary session/plan cleanup.

| Payload / path | Baseline ns/op | Migrated ns/op | Time delta | B/op delta | allocs/op delta |
| --- | ---: | ---: | ---: | ---: | ---: |
| 32 B / reusable Native | 747.0 | 1,050.0 | +40.56% | +352 | +4 |
| 32 B / reusable UAPI | 753.9 | 1,073.0 | +42.33% | +368 | +5 |
| 64 KiB / reusable Native | 270,815 | 268,766 | -0.76% | +352 | +4 |
| 64 KiB / reusable UAPI | 269,792 | 269,653 | -0.05% | +368 | +5 |
| 1 MiB / reusable Native | 4,153,763 | 4,005,139 | -3.58% | +358 | +4 |
| 1 MiB / reusable UAPI | 4,129,114 | 4,025,406 | -2.51% | +372 | +5 |
| 32 B / convenience Native | 347,079 | 348,552 | +0.42% | +352 | +4 |
| 32 B / convenience UAPI | 346,086 | 342,016 | -1.18% | +367 | +5 |
| 1 MiB / convenience Native | 4,510,468 | 4,376,051 | -2.98% | +355 | +4 |
| 1 MiB / convenience UAPI | 4,477,911 | 4,373,627 | -2.33% | +375 | +5 |

Discard and writer workloads have the same allocation counts as collection:
no additional allocation comes from invoking the callback. The benchmark's
consumer is created outside the timed loop so a caller-created escaping closure
does not distort that comparison. Small reusable discard/writer medians increase
32.24–42.67%; 64 KiB results are close to baseline. Output-only collection/discard
costs about 229–286 ns, 400 B, and five allocations regardless of payload size,
against the baseline's optimized direct struct access at about 2 ns and no
allocations. UAPI's extra wrapper accounts for 16 B and one allocation per run.
The 1 MiB writer's amortized buffer allocation depends on iteration count.

The allocation increase is fixed, rather than an additional payload-sized buffer.
Tests establish encoder-buffer identity through native and UAPI collection and
consumption, and borrowed-byte validity during closure. Unread outputs release
temporary permits and VMs, and settled handles discard payload and cancellation
references. Large-payload timing gains vary between runs and are not attributed
to this migration.

Five samples do not establish a 95% confidence interval. Output-only comparisons
against an optimized materialized struct measure the added lifecycle cost, not a
throughput regression in encoding. Native encoding remains byte-backed and is
not constant-memory streaming. Core results establish neither Wire streaming nor
the separate CLI allocation regression.

Validation completed after the implementation review:

```sh
GOCACHE=/private/tmp/ferret-output-plan-go-cache CGO_ENABLED=1 go test -race \
  ./pkg/engine/... ./pkg/debugger ./uapi . ./compat/... \
  -run 'Test(Output|RealOutput|RunPanic|SessionAfterRunFailure|EngineRunPreserves|Debug|Session|Native|Execution|Runtime)' -count=1
GOCACHE=/private/tmp/ferret-output-plan-go-cache make fmt
GOCACHE=/private/tmp/ferret-output-plan-go-cache \
  STATICCHECK_CACHE=/private/tmp/ferret-output-staticcheck-cache make lint
GOCACHE=/private/tmp/ferret-output-plan-go-cache make test
GOCACHE=/private/tmp/ferret-output-plan-go-cache make compile compile-32bit
GOCACHE=/private/tmp/ferret-output-plan-go-cache GOWORK=off go mod tidy -diff
GOCACHE=/private/tmp/ferret-output-plan-go-cache GOWORK=off go -C tools/apiref mod tidy -diff
GOCACHE=/private/tmp/ferret-output-plan-go-cache GOWORK=off go -C tools/apipublish mod tidy -diff
git diff --check
```

The full test target covers unit and integration race suites, compatibility, both
independent tool modules, the root facade, and security tests. HTTP security tests
required permitted local loopback access. Some sandboxed tools reported denied
module stat-cache writes but completed successfully. CLI smoke checks verify
successful forwarding, exit status for execution errors, preserved source
diagnostics, and no extra blank line for absent error-only content.

The mandatory complete-diff review corrected preservation of both cancellation
causes, lone runtime diagnostic formatting, result rollback when an after-run
hook or encoder panics, and missing proof of unread-output capacity release and
portable session reuse. Affected checks were rerun. Local raw benchmark, benchstat,
validation, and CLI evidence is retained in ignored `bin/consumable-output/`.

Linux/386 builds and test compilation passed. Actual Linux/386 test execution
remains unverified: the local host is Darwin/arm64 and its installed Docker client
has no running daemon. Parser generation was unnecessary because no parser or
generator input changed. Dependency updates preserve the tool module's relative
root replacement; ignored vendor contents were refreshed to API alpha.21.

## Downstream requirements

This migration changes Core only. Follow-up integrations must:

- **Wire:** adopt consumable handles, cancellation/closure, portable terminal
  diagnostics, metadata, content presence, and detached Content serialization.
- **FerretD:** migrate ordinary execution and debugger/DAP completion content
  while retaining daemon/session ownership and event-plus-error handling.
- **CLI:** explicitly collect or forward through Consume, including side-effect
  runs; keep invocation contexts alive and observe terminal errors.
- **Lab, editors, and tooling:** migrate live execution consumers and materialized
  event payloads to the new signatures and JSON shape.
- **Custom codecs:** honor returned-byte ownership in ordinary and configured
  encoder implementations, including error returns.
- **Website documentation:** update embedding, configuration, SDK, codec, and
  remote-runtime examples to consumption, new error timing, and detached Content.

No FQL semantics, parser inputs, VM-result contract, transport implementation,
downstream repositories, or release publication are changed by this migration.
