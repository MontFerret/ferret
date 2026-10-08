# Codec I/O

`pkg/encoding` owns content-type registration and the runtime-value codec
contract. JSON and MessagePack implement context-aware synchronous I/O:

```go
Encode(ctx context.Context, dst io.Writer, value runtime.Value) error
Decode(ctx context.Context, src io.Reader) (runtime.Value, error)
```

The caller lends I/O only for the operation. Codecs never close it, retain it, or
use it after returning. Write buffers may be reused when Write returns. Success
completes codec-owned writes and flushing; callers flush, finalize, and close
their own wrappers. Failure may leave a written prefix or consumed input.
Non-nil I/O errors retain their identity. A short write without an error returns
`io.ErrShortWrite` and stops without retrying.

`EncodeBytes(ctx, enc, value)` uses a bytes.Buffer and returns its written prefix
with any encoding error. `DecodeBytes(ctx, dec, data)` uses bytes.NewReader.
Atomic publication requires buffering and checking encoding success before
publishing those bytes. There is no compatibility interface for the former
byte methods. SDK host-value codecs and bytecode serialization are independent.

## Hooks and contexts

Encoder hooks run recursively in depth-first order for visited runtime values:

```go
type PreEncoderHook func(context.Context, runtime.Value) error
type PostEncoderHook func(context.Context, runtime.Value, error) error
type PreDecoderHook func(context.Context) error
type PostDecoderHook func(context.Context, runtime.Value, error) error
```

Decoder hooks run once per operation. Their post-hook receives the materialized
value on decoding success and runtime.None on decoding failure. Same-kind hooks
run in registration order and stop at the first hook error. A pre-hook failure
skips the operation and its post-hooks; ancestor encoder post-hooks still see
descendant failures. Joined errors preserve operation and hook causes.

Contexts reach hooks and runtime capabilities. Built-ins reject nil contexts and
nil I/O with errors wrapping runtime.ErrInvalidArgument and check cancellation at
entry, bounded I/O boundaries, every 256 traversal values/tokens, and between
large payload chunks. They avoid ctx.Err calls on every item. Arbitrary blocking
I/O and context-free custom marshalers cannot be forcibly interrupted.
Context-free Go MarshalJSON protocols use background context only at that
unavoidable boundary; stdlib and HTTP adapters use their caller's context.

## Implementations and materialization

JSON writes numbers, containers, escaped strings, and Base64 progressively.
Numeric scratch and bounded 4 KiB payload scratch avoid container materialization
for generic iterables. Custom marshalers and generic Go-value fallbacks may
buffer. Its standard JSON tokenizer handles fragmented readers and strictly
rejects extra roots, malformed separators, and truncated containers; integer
width is preserved.

MessagePack connects the dependency encoder/decoder to borrowed I/O through a
short-write and cancellation adapter. Instances are reset on every operation and
release I/O references on every exit. Generic iterables still materialize to
determine the length required by a MessagePack container header. Decoding in both
formats returns a fully materialized runtime.Value and may read ahead. This is
not a constant-memory guarantee or a lazy decoding API.

Both encoders borrow the source value. Their internal iterator handle owns only
a distinct acquired iterator, traversed with runtime.ForEachIter. An iterator
that is the source itself, or shares its runtime.Resource identity, remains with
the source's owner. Comparable identities are checked safely; non-comparable
resource aliases require runtime.Resource for stable identity. The encoders do
not change runtime.ForEach's cleanup contract for other callers.

Distinct acquired iterators close once after traversal and child encoding,
including on failure and cancellation. MessagePack collects the items, encodes
the collected list after successful traversal, then closes the iterator. A close
failure therefore cannot skip the child hooks that register yielded resources
with the VM result. Native materialization also lends a private, synchronous
adoption callback through the operation context. MessagePack caches it once
in its codec operation, invoking it immediately after each successful yield,
before the collection cancellation checkpoint. It adopts the yielded resource
and resources already stored in built-in arrays and objects into the same VM
result as the normal ownership hook. Result resource-identity deduplication
prevents those routes from closing comparable resources or ResourceID aliases
twice. Direct codec calls without this callback keep caller-owned cleanup.

Early adoption uses iterative, cycle-safe traversal of concrete runtime.Array
and runtime.Object storage. It completes independently of cancellation, without
invoking host collection capabilities, lazy iterables, marshalers, or user
encoding hooks. Cycle protection resets for each yield so later mutations of a
repeated container do not hide newly stored resources. Failed collection still
skips child encoding and its hooks; it never reads more iterator values merely
to discover ownership. The callback and traversal scratch belong only to the
materialization operation and are not retained by codecs or their pools.

Iterator close errors join encoding errors before the parent's post-hooks. The
VM result remains responsible for closing adopted source
and child resources. Direct codec callers retain that responsibility themselves.

Native and UAPI output remain eager and byte-backed. Native materialization
collects into an owned buffer with the original invocation context; before-run
replacement contexts remain execution-only. Debugger materialization receives
the active execution or resume-command context. Successful zero-byte encoding
produces present empty Content; failed zero-byte encoding leaves Content absent.
Nonempty error prefixes have LengthKnown=false. Encoding and cleanup error
identities, resource adoption hooks, and result/session cleanup ordering remain
intact. Collected data survives codec scratch reuse and result/session closure.
Query execution, record iteration, and the output lifecycle are unchanged.

## Validation and performance

Shared reader/writer tests cover prefix failures, short and zero-progress writes,
fragmented and data-plus-error reads, cancellation, hook context, runtime
capabilities, and I/O reuse/closure. Codec-specific tests retain formats, deep
nesting, integer widths, custom marshaler precedence, escaping, and hook order.
Native/UAPI/debugger tests cover context forwarding and retained Content.
Engine regressions also cover self-returned iterators and yielded resources when
iterator closure fails, proving one cleanup owner through result and session
teardown. Cancellation regressions cover directly yielded resources, stored
container children, aliases, deep nesting, cycles, mutated containers, and a
resource yielded exactly at a sampled cancellation checkpoint. Child cleanup
errors remain separate from traversal and iterator cleanup errors. Shared codec
tests cover resource aliases, depth-first hooks, joined cleanup errors,
cancellation, and dynamically non-comparable identities.

Compare the existing byte-oriented workloads with:

```sh
go test ./pkg/encoding/json ./pkg/encoding/msgpack \
  -run '^$' -bench 'Benchmark(JSON|Msgpack)Codec(Encode|Decode)$' \
  -benchmem -benchtime=200ms -count=5
```

`BenchmarkCodecWriter` additionally measures large arrays, objects, ranges,
1 MiB strings/binary, default/configured encoders, background/cancellable
contexts, and reused bytes.Buffer/io.Discard destinations. Buffer workloads own
the accumulated output; Discard isolates codec overhead. Neither includes value
construction in timing. Decoder materialization and the buffering exceptions
above remain part of the design.

## Migration measurements (2026-10-07)

Baseline: Core 35693eb380d10bd80de76208b95e7598307e34d6. Migrated: the
reviewed reader/writer implementation in this change. Both use Go 1.27.1,
Darwin/arm64, Apple M2 Max, GOMAXPROCS=12, and the five-sample command above.
The migrated byte workloads use the byte helpers and b.Context; the former
codec had an internal background context. These comparisons include collecting
buffer allocation and cancellable-context handling.

| Byte workload | Baseline ns/op | Migrated ns/op | Time delta | B/op delta | allocs/op delta |
| --- | ---: | ---: | ---: | ---: | ---: |
| JSONCodecEncode/duration | 68.2 | 158.4 | +132.3% | +96 | +1 |
| JSONCodecEncode/regexp | 102.2 | 163.2 | +59.7% | +96 | +1 |
| JSONCodecEncode/range_1024 | 10,774.0 | 18,056.0 | +67.6% | +61 | -4 |
| JSONCodecEncode/flat_array_1024 | 26,862.0 | 25,125.0 | -6.5% | -2,875 | -923 |
| JSONCodecEncode/flat_object_256 | 17,061.0 | 18,110.0 | +6.1% | -410 | -155 |
| JSONCodecEncode/nested_array_10000 | 314,186.0 | 410,278.0 | +30.6% | +119 | +1 |
| JSONCodecEncode/nested_object_5000 | 780,898.0 | 805,101.0 | +3.1% | -79,695 | +1 |
| JSONCodecDecode/flat_array_1024 | 61,513.0 | 81,629.0 | +32.7% | +3,816 | +1,022 |
| JSONCodecDecode/flat_object_256 | 37,963.0 | 46,475.0 | +22.4% | -2,328 | +256 |
| JSONCodecDecode/nested_array_10000 | 784,485.0 | 1,019,629.0 | +30.0% | +301,360 | +24 |
| JSONCodecDecode/nested_object_5000 | 818,928.0 | 1,099,666.0 | +34.3% | +218,147 | +50 |
| MsgpackCodecEncode/regexp | 57.6 | 150.1 | +160.5% | +176 | +2 |
| MsgpackCodecEncode/range_1024 | 18,135.0 | 21,964.0 | +21.1% | +176 | +2 |
| MsgpackCodecEncode/flat_array_1024 | 13,152.0 | 17,425.0 | +32.5% | +176 | +2 |
| MsgpackCodecEncode/flat_object_256 | 11,474.0 | 13,265.0 | +15.6% | +176 | +2 |
| MsgpackCodecEncode/nested_array_10000 | 267,482.0 | 296,094.0 | +10.7% | +177 | +2 |
| MsgpackCodecEncode/nested_object_5000 | 671,063.0 | 687,322.0 | +2.4% | +177 | +2 |
| MsgpackCodecDecode/flat_array_1024 | 27,213.0 | 39,460.0 | +45.0% | +184 | +3 |
| MsgpackCodecDecode/flat_object_256 | 15,634.0 | 19,647.0 | +25.7% | +184 | +3 |
| MsgpackCodecDecode/nested_array_10000 | 516,408.0 | 554,153.0 | +7.3% | +179 | +3 |
| MsgpackCodecDecode/nested_object_5000 | 786,686.0 | 823,631.0 | +4.7% | +184 | +3 |

Flat JSON array encoding removes per-integer formatting allocations. Other
workloads expose reader/writer adaptation costs, context checkpoints, and the
standard JSON tokenizer's materialization overhead. This is not an across-the-board
speed improvement. Five samples do not establish benchstat's 95% confidence
interval; small timing shifts should be treated as inconclusive.

Writer measurements use:

```sh
go test ./pkg/encoding -run '^$' -bench '^BenchmarkCodecWriter$' \
  -benchmem -benchtime=100ms -count=3
```

The writer table reports three-sample medians. Arrays contain 16,384 integers,
objects 4,096 properties, and ranges 65,536 integers. The configured column
adds one no-op pre-hook and a cancellable context to io.Discard. Buffer is reused
and warmed before timing, so its allocation does not represent creating a fresh
collected payload. All context/configuration combinations exist in the benchmark.

| Codec / writer workload | Buffer ns/op | Discard ns/op | Discard B/op / allocs | Configured cancellable Discard ns/op | Configured B/op / allocs |
| --- | ---: | ---: | ---: | ---: | ---: |
| JSON / array | 364,633 | 346,635 | 0 / 0 | 366,840 | 0 / 0 |
| JSON / object | 238,875 | 232,430 | 65,652 / 4,098 | 235,280 | 65,652 / 4,098 |
| JSON / range | 1,029,103 | 954,243 | 0 / 0 | 1,824,670 | 522,524 / 65,280 |
| JSON / string_1MiB | 2,591,703 | 2,572,994 | 0 / 0 | 2,582,476 | 0 / 0 |
| JSON / binary_1MiB | 461,212 | 433,377 | 0 / 0 | 440,346 | 0 / 0 |
| MessagePack / array | 250,067 | 212,908 | 80 / 1 | 250,544 | 80 / 1 |
| MessagePack / object | 185,509 | 163,530 | 65,732 / 4,098 | 176,279 | 65,731 / 4,098 |
| MessagePack / range | 1,341,461 | 1,210,453 | 522,592 / 65,281 | 1,320,334 | 522,596 / 65,281 |
| MessagePack / string_1MiB | 23,814 | 998 | 80 / 1 | 1,973 | 80 / 1 |
| MessagePack / binary_1MiB | 24,824 | 999 | 80 / 1 | 1,964 | 80 / 1 |

Configured JSON ranges box values for recursive hooks; MessagePack ranges
already box values during encoding. Objects retain runtime traversal allocations.
JSON uses bounded scratch, but custom marshalers and generic fallbacks may buffer;
MessagePack materializes generic iterables for header lengths. Both decoders
materialize values. Native output still collects the entire encoded document.

## Review and validation notes

Second-pass review corrected a JSON closing-quote scratch boundary, ignored
MessagePack custom-encoder writer failures, and reader errors joined with EOF.
Regression tests cover these error causes and iterator cancellation settlement.
No parser or generator input changed. SDK host-value codecs and bytecode
serialization are outside this migration. The FQL-visible correction is that
encoding::json_parse now rejects malformed separators and truncated or trailing
documents that the former tokenizer could accept; query execution and output
lifecycle semantics remain intact.

Validation completed after those fixes, with writable task-scoped Go caches and
temporary directories:

```sh
go test ./pkg/encoding/... ./pkg/stdlib/encoding ./pkg/engine/... ./uapi/... -count=1
GOTOOLCHAIN=go1.26.8 make test
make compile compile-32bit
GOTOOLCHAIN=go1.26.8 make fmt
GOTOOLCHAIN=go1.26.8 make lint
git diff --check
```

The focused suites and compilation used Go 1.27.1. The full test target passed
on Go 1.26.8, including unit/integration race suites, independent tool modules,
compatibility, the facade, and security. Local HTTP tests required permitted
loopback access. A Go 1.27.1 full-target run reached the API-reference analyzer
and reported a race inside go/types during x/tools package loading; that
toolchain run remains unresolved. Linux/386 production and test compilation
passed; actual Linux/386 execution was not performed on this Darwin host.

Formatting and lint used the pinned tools with Go 1.26.8 because existing
source-analysis binaries could not process the Go 1.27 standard library. The
initial Revive run also included unrelated generated parser files in an ignored
bin/ snapshot. The successful rerun added a temporary Revive exclusion for
bin/ without modifying the snapshot or repository lint configuration.

Website codec/output guidance was synchronized while preserving its existing
edits. Its complete custom-codec example compiled and ran against this Core
checkout. The website's final mage build and git diff --check passed, including
the generated reference, Hugo, Pagefind/search, and Registry route pipeline.

### Iterator ownership fix measurements

The iterator cleanup correction was measured separately against
87684ac70684 on 2026-10-07, using Go 1.27.1 on Darwin/arm64, Apple M2 Max,
GOMAXPROCS=12. The generic-iterable benchmark was added before changing production
code, and the same workloads were run before and after the correction:

```sh
go test ./pkg/encoding/json ./pkg/encoding/msgpack \
  -run '^$' -bench 'Benchmark(JSON|Msgpack)CodecEncode$' \
  -benchmem -benchtime=200ms -count=5
go test ./pkg/encoding -run '^$' -bench '^BenchmarkCodecIterable$' \
  -benchmem -benchtime=200ms -count=5
```

The established encoding workloads retained every allocation count; median
bytes/op differed by at most 13 bytes. Final timing medians ranged from -4.4% to
+3.3% against the initial baseline. An earlier head
run was 2.5–10.2% slower; repeating the baseline from an isolated archive and
repeating the head did not reproduce that increase. Five samples are insufficient
for benchstat's 95% confidence interval. These timings do not establish a
general speedup or regression.

The generic workloads yield 1,024 values and use a warmed reusable bytes.Buffer
or io.Discard. These are median results; bytes/op differ by at most one byte, and
every allocation count is unchanged:

| Codec / iterator / writer | Before ns/op | Final ns/op | Time delta | Before → final B/op | Allocs/op |
| --- | ---: | ---: | ---: | ---: | ---: |
| JSON / owned / Buffer | 31,965 | 31,667 | -0.9% | 6,148 → 6,148 | 768 |
| JSON / owned / Discard | 30,625 | 30,387 | -0.8% | 6,147 → 6,147 | 768 |
| JSON / borrowed / Buffer | 32,275 | 31,585 | -2.1% | 6,148 → 6,148 | 768 |
| JSON / borrowed / Discard | 29,960 | 30,257 | +1.0% | 6,146 → 6,147 | 768 |
| MessagePack / owned / Buffer | 30,980 | 30,056 | -3.0% | 66,043 → 66,043 | 782 |
| MessagePack / owned / Discard | 28,979 | 28,245 | -2.5% | 66,042 → 66,042 | 782 |
| MessagePack / borrowed / Buffer | 31,153 | 30,060 | -3.5% | 66,043 → 66,043 | 782 |
| MessagePack / borrowed / Discard | 28,828 | 28,286 | -1.9% | 66,042 → 66,042 | 782 |

Profiling during the second-pass review found two allocations, approximately
eight bytes, in Go 1.27's general reflect.Value.Comparable check for pointer
sources. The final pointer-identity path avoids that check; other values retain
safe dynamic-comparability checks. Tests cover comparable values, pointers with
non-comparable acquired iterators, and non-comparable resource aliases. The
fixture's numeric iterator keys still allocate, and MessagePack still buffers
yielded values to determine its array header length.

The engine regressions failed before the fix with two closes of self-returned
iterators in both codecs and zero closes of yielded MessagePack resources on an
iterator-close failure. They pass after the fix, including preservation of
iterator and result-cleanup error causes through output consumption and teardown.
The complete second-pass review also corrected fixture organization, spacing,
and redundant writer declarations reported by Revive.

The fix passed focused codec/engine tests on Go 1.27.1, the full make test target
on Go 1.26.8 with permitted loopback access, make compile and make compile-32bit,
make fmt, make lint, and git diff --check. The added comparable-value case was
also rerun with the race detector. Formatting/lint used the previously described
Go 1.26.8 tools and temporary ignored-bin exclusion. Linux/386 was compiled but
not executed. The website encoder guide was updated without altering its
unrelated existing edits, and its mage build and git diff --check passed.

### Early MessagePack resource adoption measurements

The canceled-collection regression was reproduced against 852e08c76afc before
changing production code: yielded children received zero closes, including
children in materialized containers and a resource yielded at the collection
checkpoint. The fix registers ownership before that checkpoint without running
child encoding hooks. Tests cover borrowed and owned iterators, traversal,
deadline and cancellation errors, joined iterator and child cleanup failures,
aliases, deep and cyclic containers, mutated repeated containers, direct codec
cleanup responsibility, callback reuse, and successful depth-first hook order.

Measurements on 2026-10-08 used Go 1.27.1, Darwin/arm64, Apple M2 Max, and
GOMAXPROCS=12, with writable task-scoped caches and temporary directories. The
native benchmark and its fixtures were added before the production change.
Identical commands ran serially before and after; an isolated archive of the
baseline was also measured again with the same native fixtures:

```sh
go test ./pkg/encoding/json ./pkg/encoding/msgpack \
  -run '^$' -bench 'Benchmark(JSON|Msgpack)CodecEncode$' \
  -benchmem -benchtime=200ms -count=5
go test ./pkg/encoding -run '^$' -bench '^BenchmarkCodecIterable$' \
  -benchmem -benchtime=200ms -count=5
go test ./pkg/engine -run '^$' -bench '^BenchmarkSessionMsgpackIterable$' \
  -benchmem -benchtime=200ms -count=5
```

The established byte-helper workloads retain all allocation counts and differ
by at most two median bytes/op against the initial baseline. Timing medians
range from -4.9% to -2.8% for unchanged JSON and -2.4% to 0.0% for MessagePack.
Against the repeated baseline, they range from -3.0% to -0.5% and -3.1% to -0.4%,
respectively. The unchanged JSON control also moved, so these results do not
establish a general speedup. Five samples do not provide benchstat's 95%
confidence interval.

Generic iterables yield 1,024 integers to a warmed bytes.Buffer or io.Discard.
These are median results against the initial baseline; allocation counts are
unchanged:

| Codec / iterator / writer | Before ns/op | Final ns/op | Time delta | Before → final B/op | Allocs/op |
| --- | ---: | ---: | ---: | ---: | ---: |
| JSON / owned / Buffer | 31,272 | 31,054 | -0.7% | 6,148 → 6,147 | 768 |
| JSON / owned / Discard | 30,060 | 29,798 | -0.9% | 6,146 → 6,146 | 768 |
| JSON / borrowed / Buffer | 31,447 | 31,325 | -0.4% | 6,147 → 6,148 | 768 |
| JSON / borrowed / Discard | 30,313 | 30,087 | -0.7% | 6,147 → 6,146 | 768 |
| MessagePack / owned / Buffer | 29,497 | 30,002 | +1.7% | 66,043 → 66,043 | 782 |
| MessagePack / owned / Discard | 27,622 | 28,533 | +3.3% | 66,042 → 66,042 | 782 |
| MessagePack / borrowed / Buffer | 29,619 | 30,202 | +2.0% | 66,043 → 66,043 | 782 |
| MessagePack / borrowed / Discard | 27,772 | 28,346 | +2.1% | 66,042 → 66,042 | 782 |

The repeated baseline places MessagePack's generic-iterable increase at
1.0–2.3%. Review removed a prototype callback field that enlarged boxed codecs
by 16 bytes, then removed an extra callback argument from recursive encoding
calls that increased established MessagePack timings by approximately 6%.
The final implementation caches the callback in the operation context, preserving
codec and operation struct sizes and the existing recursive call signatures.

The native workload runs a reusable session returning 1,024 scalar values,
closable custom-marshaler values, or arrays containing objects with those
resources. Run, Collect, output closure, and iterator/resource reset are timed;
engine/plan/session construction is excluded. Median results against the initial
baseline are:

| Yielded values | Before ns/op | Final ns/op | Time delta | Before → final B/op | Before → final allocs/op |
| --- | ---: | ---: | ---: | ---: | ---: |
| Scalars | 45,352 | 50,136 | +10.5% | 75,765 → 75,920 | 816 → 820 |
| Direct resources | 249,739 | 280,179 | +12.2% | 435,480 → 435,626 | 3,924 → 3,928 |
| Resources inside containers | 394,889 | 518,225 | +31.2% | 550,230 → 550,619 | 5,973 → 5,980 |

Against the repeated baseline, the timing increases are 10.1%, 9.3%, and 28.8%.
Early adoption adds operation-scoped callback/context state and resource
registration in addition to the existing ownership hook. Stored-container
discovery also adds a cycle set and iterative traversal scratch. These costs
buy ownership settlement when collection or later encoding fails; they are not
a constant-memory guarantee. MessagePack still buffers all yielded values for
its array header, and the fixtures still allocate numeric keys and custom
marshaler buffers.

The mandatory second-pass review corrected per-yield cycle protection so a
mutated repeated container is rescanned, removed the measured codec/call-shape
overhead, and added a cancellation-checkpoint case whose cyclic container is
adopted after the context becomes canceled. Direct encoding tests cover both
success and failed collection without closing yielded resources. A host-map
regression verifies one callback lookup per operation even when predicates use
a different context. No public API, serialization, user-hook order, FQL
semantics, query execution, or output lifecycle changed.

Final validation completed after those corrections:

```sh
go test ./pkg/encoding/... ./pkg/internal/encodingownership \
  ./pkg/engine ./pkg/debugger/... ./uapi/... -count=1
GOTOOLCHAIN=go1.26.8 go test -race ./pkg/engine ./pkg/encoding/... \
  ./pkg/internal/encodingownership \
  -run 'Iterator|MsgpackOutputAdopts|MsgpackEarlyAdoption|CollectedResources|AdoptionCallback|ValueAdopter|DirectEncodingLeaves|CachesAdoption' \
  -count=1
GOTOOLCHAIN=go1.26.8 make test
make compile compile-32bit
GOTOOLCHAIN=go1.26.8 make fmt
GOTOOLCHAIN=go1.26.8 make lint
git diff --check
```

Focused tests, benchmarks, and compilation used Go 1.27.1; full tests, targeted
race tests, formatting, and lint used Go 1.26.8. The full target includes unit
and integration race suites, independent tools, compatibility, the facade,
and security tests with permitted loopback access. The targeted filter has no
matching JSON tests; the full target exercises JSON with the race detector.
Formatting/lint retained the previously described toolchain and temporary
ignored-bin exclusion. Read-only module-cache stat writes emitted non-fatal
warnings during compile/lint. Linux/386 compiled successfully but was not
executed. The website encoder guide was synchronized while preserving its
unrelated edits; its mage build and git diff --check passed. Parser generation
was not needed.
