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
for generic iterables; runtime.ForEach retains iterator cleanup and joined close
errors. Custom marshalers and generic Go-value fallbacks may buffer. Its standard
JSON tokenizer handles fragmented readers and strictly rejects extra roots,
malformed separators, and truncated containers; integer width is preserved.

MessagePack connects the dependency encoder/decoder to borrowed I/O through a
short-write and cancellation adapter. Instances are reset on every operation and
release I/O references on every exit. Generic iterables still materialize to
determine the length required by a MessagePack container header. Decoding in both
formats returns a fully materialized runtime.Value and may read ahead. This is
not a constant-memory guarantee or a lazy decoding API.

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
