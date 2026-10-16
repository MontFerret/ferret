# Codec streaming migration report

Historical measurements and validation for [PR #1058](https://github.com/MontFerret/ferret/pull/1058).
The current contract is documented in [Codec I/O](codecs.md). The historical
sections below describe intermediate implementations; the final comparison and
validation section records the follow-up results.

## Cancellation initialization correction (2026-10-14)

Tested PR head `704fe1a603892a3de5de6bb61fa4209f0e9cb4d2` plus the
uncommitted correction to `codecutil.NewOperation`, the four JSON/MessagePack
entry points, their codec/session regressions and fixtures, and the lifecycle
documentation. The Core checkout was clean before this correction. Existing
website edits were preserved outside the narrow Value Encoders clarification.

Operation construction now validates arguments and initializes state without
checking cancellation. Encoders bind the adopter, report the stored root, and
check cancellation directly before acquiring serialization state or entering
hooks. Decoders check directly after validation, before hooks or reads. Invalid
arguments still take precedence over cancellation. The private callback cache,
stored-graph walker, iterator handle, result lifecycle, public API, successful
formats, and FQL semantics are unchanged.

Before the production correction, the strengthened root test failed independently
for JSON and MessagePack: canceled and expired contexts skipped discovery;
cancellation inside discovery allowed three writes (10 JSON bytes or 8
MessagePack bytes), and configured encoders ran two pre-hooks and two post-hooks.
The real materialization regression also failed for both codecs: result cleanup
closed neither stored descendant and lost the sentinel close error.

```sh
go test ./pkg/encoding -run '^TestCodecRootDiscoverySurvivesCancellation$' -count=1
go test ./pkg/encoding \
  -run '^TestCodec(RootDiscoverySurvivesCancellation|CanceledEncodeDoesNotEvaluateHostValues|RejectsTypedNilIOBeforeHooks|DecoderCancellationPreventsHooksAndReads)$' \
  -count=1
go test ./pkg/engine/internal/session \
  -run '^TestMaterializeCanceledRootRetainsDiscoveredResources$' -count=1
```

The new validation and decoder cases passed before the fix. Afterward, all new
regressions passed, including default/configured codecs, expired deadlines,
cancellation during discovery, callback/writer independence, typed-nil validation
precedence, ordinary and byte-scanner readers, and suppressed host capabilities
and custom marshalers. The session regression attaches nested stored descendants
only after live VM execution, uses the actual Materialize adopter, and proves
absent content, no premature closure, preserved cleanup errors, comparable-alias
deduplication, repeated cleanup, and no host/lazy traversal.

Validation used task-scoped writable GOCACHE, GOTMPDIR, TMPDIR, and
STATICCHECK_CACHE on Darwin/arm64. Go 1.27.1 passed the focused, full root-module,
and focused race suites:

```sh
go test ./pkg/encoding ./pkg/engine/internal/session \
  -run '^Test(Codec(RootDiscoverySurvivesCancellation|CanceledEncodeDoesNotEvaluateHostValues|RejectsTypedNilIOBeforeHooks|DecoderCancellationPreventsHooksAndReads)|MaterializeCanceledRootRetainsDiscoveredResources)$' \
  -count=1
go test ./pkg/encoding -run '^TestCodecRootDiscoverySurvivesCancellation$' -count=1
go test ./pkg/encoding/... ./pkg/engine/... ./pkg/vm/... ./pkg/debugger/... ./uapi/... -count=1
go test ./... -count=1
go test -race ./pkg/encoding/... ./pkg/engine/... ./pkg/vm/... -count=1
```

The first full-suite attempt failed only because sandbox restrictions denied
httptest loopback listeners in `pkg/net/http` and `test/security`. The unchanged
full-suite rerun passed with permitted loopback access.

The repository target also passed with permitted loopback access on Go 1.26.8,
including package and integration race suites, both independent API-tool modules,
compatibility, the public facade, and security tests. Test caching was disabled:

```sh
GOTOOLCHAIN=go1.26.8 GOFLAGS=-count=1 make test
```

Formatting and lint used Go 1.26.8 and the Makefile-pinned tools:

```sh
GOTOOLCHAIN=go1.26.8 make fmt
GOTOOLCHAIN=go1.26.8 make lint
GOTOOLCHAIN=go1.26.8 revive -config revive.toml -formatter stylish \
  -exclude ./pkg/parser/fql/... -exclude ./vendor/... -exclude './*_test.go' \
  -exclude ./bin/... ./...
git diff --check
```

Formatting passed without additional tracked-file changes. The unmodified lint
target passed all Staticcheck stages, then Revive reported 216 findings solely
in generated parser files under the pre-existing ignored
`bin/website-output-docs/core-snapshot`. Revive passed with that ignored `bin/`
tree excluded as shown above; no repository lint configuration or snapshot was
changed. The complete implementation, new scanner fixture, and documentation
received a second-pass review with no actionable findings.

The website Value Encoders guide received only the initialization clarification.
Its `mage build` passed with CI-pinned Hugo 0.159.1 extended on PATH, Go 1.27.1,
and a task-scoped MAGEFILE_CACHE. The rendered page contained the revised
validation/discovery/cancellation ordering; website `git diff --check` passed.
SHA-256 comparison confirmed all other 310 tracked/unignored paths unchanged,
and the edited page retained all pre-existing edits outside the correction.
The build emitted the existing `.Site.Data` deprecation warning. No performance
benchmarks were added or run for this correction; historical measurements below
remain unchanged.

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

## Follow-up implementation and review (2026-10-08)

Reviewed all three PR commits and their combined diff before implementation:
`87684ac70684`, `852e08c76afc`, and `d9b3ba820e99`, against
`35693eb380d1`. The follow-up retains the public codec API, accepted formats,
eager byte-backed output, caller I/O ownership, user-hook order, UAPI delegation,
and debugger content retention. No FQL syntax or successful serialization changes
are introduced. Joined non-EOF iterator failures now propagate instead of being
mistaken for exhaustion; wrapped and custom EOF markers remain accepted.
Because runtime.ForEach also uses this helper, FQL standard-library consumers of
host iterables now return such a joined failure instead of successful exhaustion.
This is the intentional FQL-visible error-propagation correction; normal EOF and
successful collection semantics are unchanged.

Changed responsibilities:

| Owner / files | Responsibility |
| --- | --- |
| `pkg/encoding/internal/codecutil/io.go` | Operation validation and private callback cache/registration |
| `pkg/internal/encodingownership/context.go` | Callback context key, installation and lookup |
| `pkg/encoding/{json,msgpack}/{codec,encoder}.go` | Root/acquisition discovery, hook/error rescans, iterator buffering and pool release |
| `pkg/engine/internal/session/output{,_resources}.go` | Exact codec type selection and VM-result-owned concrete-storage traversal |
| `pkg/vm/internal/mem/resource_key.go` | Safe value comparability while retaining identity fast paths |
| `pkg/runtime/iterator.go` | Joined non-EOF traversal error retention |
| Codec, session, engine, runtime and memory test/fixture files | Failure, cancellation, mutation, identity, I/O and compatibility regressions |
| `pkg/engine/output_iterator_bench_test.go`, session `output_bench*_test.go` | Both-format native depth benchmarks and adoption controls |
| `codecs.md`, this report and the measurements table | Durable contracts and reproducible migration evidence |
| Website `content/docs/embedding/go/value-encoders.md` | Public ownership and typed-nil validation guidance |

The private bridge now contains only its context key, installation, and lookup.
The operation retains a private cache wrapper under the approved performance-gate
exception: repeated common scalar measurements exceeded the 5% slowdown threshold
for a direct callback field. The wrapper lives in codecutil; encodingownership
contains no forwarding helpers. No public ownership interface, mutable handshake,
or runtime-value wrapper is needed. Session materialization selects concrete
JSON/MessagePack codec types,
including configured and pointer instances, for callback-only ownership. Custom
codecs and decorators keep their ownership pre-hook, irrespective of content type.

The callback distinguishes direct registration from stored discovery. Roots,
successful host/yield acquisitions, map keys, and abandoned containers report
ownership before checkpoints, writes, or hooks can lose them. The cancellation-
independent iterative walker inspects only concrete Array/Object storage, reuses
scratch, resets cycle protection per discovery, and rescans mutated containers.
Their ForEach methods can only return predicate errors; the walker's always-
successful predicate makes ignoring that return safe. Successful hook-free
ordinary visits register directly, avoiding repeated deep walks. User hooks and
failure exits can require repeated storage scans; those paths are not promised
linear discovery cost across the whole operation.

The small iterator handle remains necessary: runtime.ForEach closes every
acquired closer, while codecs borrow self-returned and ResourceID-alias iterators.
Value-level comparability prevents interface-bearing non-comparable identities
from panicking, with ResourceID and pointer fast paths preserved. Non-comparable
plain aliases still need ResourceID for deduplication. Typed-nil I/O is rejected
with runtime.ErrInvalidArgument using one allocation-free reflection check before
hooks or I/O. JSON pool state is cleared on return; MessagePack's PutEncoder and
PutDecoder release borrowed I/O references.

The mandatory second-pass review corrected additional abandonment paths: key or
separator writes can mutate an acquired container before failing; MessagePack's
internal collected list can hold another hook-mutated container; successful child
hooks can append beyond a captured list length. Those cases now rescan known
stored graphs before returning. This preserves the list snapshot and hook order
while settling resources. JSON now calls List.At before writing the next
separator, so a separator failure can acquire/adopt the next item before returning;
a failed At may leave a different partial prefix. This is the intentional
acquisition-before-abandonment change, while successful bytes remain unchanged.
The internal MessagePack buffer is scanned around user
hooks or failure; hook-free successful collection retains linear traversal.
An unnecessary extra MessagePack Reset was removed after checking the dependency's
pool-release implementation. Review also added a custom-codec decorator test with
the private callback disabled, proving ownership-hook fallback even under a
built-in content type. Callback reuse assertions now count every invocation so
a stale callback cannot go unnoticed on a different subsequent value. Positive
function-valued I/O checks ensure typed-nil validation still accepts valid values.
An unnecessary host-list fixture forwarding method was removed; embedding
already provides that capability. No broader runtime ownership unification was added.

Regressions cover both formats' acquisition and I/O failures, typed-nil pointer,
function, slice and map I/O, root/child hook failures, hook mutations, canceled
cyclic yields, 32,768-level stored graphs, host predicate contexts, callback reuse,
and exact cleanup through repeated teardown. Existing coverage retains iterator
aliases and identity policy, joined traversal/encoding/cleanup errors, partial and
empty-content presence, UAPI behavior, and debugger retention. The initial
regressions demonstrably failed on d9b3ba82 (including both identity and typed-nil
panics); targeted checks passed after the fixes.

### Validation

Task-scoped writable GOCACHE, GOTMPDIR, TMPDIR and STATICCHECK_CACHE were used.
Go 1.27.1 ran focused tests, benchmarks, and compilation; supported Go 1.26.8 ran
race/full tests, formatting, and lint. Commands actually run after review fixes:

```sh
go test ./pkg/encoding/... ./pkg/runtime/... ./pkg/vm/... ./pkg/engine/... \
  ./pkg/debugger/... ./pkg/stdlib/encoding ./uapi/... -count=1
GOTOOLCHAIN=go1.26.8 go test -race ./pkg/encoding/... \
  ./pkg/internal/encodingownership ./pkg/runtime/... ./pkg/vm/internal/mem \
  ./pkg/engine/... ./uapi/... \
  -run 'Iterator|AcquisitionAdopts|RootDiscovery|BuiltinAdopts|HookMutations|RejectsTypedNil|Dynamically|Adoption|ValueAdopter|Output|StoredMutation|Exhaustion|Adopts|Appended|CustomCodecRetains' \
  -count=1
GOTOOLCHAIN=go1.26.8 make test
GOTOOLCHAIN=go1.26.8 make fmt
GOTOOLCHAIN=go1.26.8 make lint
make compile compile-32bit
git diff --check
```

Full tests passed with permitted loopback access after the sandbox blocked a
local httptest listener. Lint passed with a temporary revive wrapper excluding
`./bin/...`: the initial unmodified target reported 216 generated-parser findings
inside an existing ignored website source snapshot there. No repository lint
configuration or snapshot was changed. Read-only module-cache stat writes emitted
non-fatal warnings. Formatting's complete rewrite surface was inspected; it
introduced no unrelated tracked-file edits. Linux/386 binaries compiled and tests
compiled with `-exec=true`; they were not executed.

The website's Value Encoders guide was synchronized, preserving every other dirty
file and unrelated text in the edited page. Its `mage build` passed with CI's
Hugo 0.159.1, rendered ownership/nil-validation text and Output links were checked,
and website `git diff --check` passed. The unchanged complete Value Encoders
example was extracted verbatim into a temporary module with a local replace to
this checkout; `GOWORK=off GOTOOLCHAIN=go1.27.1 go run -mod=mod .` printed
`text/plain` and `Hello, Ferret!`. Focused codec/engine tests validate the
short snippets' API and ownership contracts. No parser generation,
asynchronous codec, streaming VM output, or broader ownership framework is included.

### Final performance comparison

The full comparison covers 119 workloads at each of five revisions. All runs
used Go 1.27.1, Darwin/arm64 on an Apple M2 Max, GOMAXPROCS=12,
`-benchmem -benchtime=200ms -count=10`, without race instrumentation. Revisions:

| Label | Revision |
| --- | --- |
| Pre-migration | `35693eb380d10bd80de76208b95e7598307e34d6` |
| Initial reader/writer migration | `87684ac7068401e3fe10a81ba085e86264399c7a` |
| Iterator correction | `852e08c76afcb4a9f994926011de5e2a27fdf817` |
| Original follow-up head | `d9b3ba820e995fe18c0cae81df319af212010bd2` |
| Final | Uncommitted follow-up to d9b3ba82 described above |

Temporary baseline copies came from `git archive REV`, outside the checkout.
The shared Writer/Iterable and native iterable fixtures were identical across
copies. Only calls that crossed the old byte-based API were adapted: the former
Encode returned bytes, so its writer adapter necessarily collects before writing;
the former hooks had no context parameter. Its cancellation dimensions cannot
measure cancellation support that did not exist in that API. These old collecting
costs are part of the comparison, not a direct-writer implementation backported
to the baseline.

Commands run serially in each temporary copy:

```sh
export GOTOOLCHAIN=go1.27.1 GOMAXPROCS=12
# With task-scoped writable GOCACHE, TMPDIR and GOTMPDIR:
go test ./pkg/encoding/json ./pkg/encoding/msgpack ./pkg/encoding ./pkg/engine \
  -run '^$' -bench 'Benchmark(JSONCodec|MsgpackCodec|CodecWriter|CodecIterable|SessionCodecIterable)' \
  -benchmem -benchtime=200ms -count=10 -timeout=60m
# Separate common-workload repeat at all five revisions:
go test ./pkg/encoding/json ./pkg/encoding/msgpack -run '^$' \
  -bench 'Benchmark(JSON|Msgpack)Codec(Encode|Decode)' \
  -benchmem -benchtime=200ms -count=10
# Private adoption bridge exists at d9b3ba82 and final:
go test ./pkg/engine/internal/session -run '^$' -bench BenchmarkCodecAdoption \
  -benchmem -benchtime=200ms -count=10
benchstat before.txt after.txt
```

The initial broad runs overlapped brief validation/build activity. Several common
historical timings were noisy: d9b3ba82's MessagePack flat-array encode ranged
from 18,230 to 86,994 ns/op. The separate common repeats ran without concurrent
validation builds; their medians replace those 21 rows in the supporting table.
For the other 98 workloads, an interquartile timing spread above 10% of the
median selected 36 revision/workload rows for another serial repeat. The selection
was fixed before looking at the new results, including apparent improvements.
A final matched repeat isolated an apparent 17% MessagePack string/Buffer
regression, its binary/Buffer control, and an apparent 30% JSON object/Discard
improvement at d9b3ba82/final. Neither large timing difference reproduced; the
three matched rows also use their repeated medians. Those repeated medians replace the corresponding rows; other rows retain their
complete-pass medians. Statistical comparisons use benchstat; timing differences between unchanged paths are not attributed to
this implementation merely because they have different medians.

[Complete measurements](codec-streaming-measurements.md) record ns/op, B/op and
allocs/op for every workload and all five revisions, plus separate adoption and
callback transport controls. Writer dimensions distinguish warmed Buffer from
Discard, background from cancellable contexts, and absent from no-op user hooks.
Generic iterable dimensions distinguish owned from borrowed iterators. Native
materialization includes Run, Collect, output closure, iterator/resource reset,
and leaf marshaler buffers; construction is excluded. Containers alternate
Array/Object pairs at depths 1, 8 and 32, with 1,024 yielded items. Adoption
controls use Discard, fixed small leaf payloads, fresh vm.Result ownership, and
absent/no-op/real callbacks to separate invocation, discovery and cleanup costs.
At d9b3ba82, direct JSON encoding ignores this private callback, so its result
control does not adopt the resources. The final result control adds actual
registration/cleanup; that difference is a guarantee change, not pure overhead
for equivalent behavior.

#### Measured implementation choices

A direct callback field grew Operation from 32 to 40 bytes and moved MessagePack's
escaping writer to a larger allocation class. Two alternating wrapper/field
cohorts, each with ten samples per variant, reproduced the gate violation:

| Existing workload | Wrapper ns/op | Field ns/op | Change | B/op wrapper / field | Allocs/op |
| --- | ---: | ---: | ---: | ---: | ---: |
| JSON duration | 135.9 | 146.5 | +7.76% | 208 / 208 | 3 / 3 |
| JSON regexp | 155.7 | 165.2 | +6.10% | 208 / 208 | 3 / 3 |
| MessagePack regexp | 140.5 | 153.1 | +8.89% | 288 / 304 | 4 / 4 |

All three comparisons have p<0.001, n=20. A fresh full common-workload ablation
with the final ownership fixes reproduced MessagePack regexp's slowdown at
+10.37% (p<0.001, n=10); the JSON scalar deltas were +4.97% and +4.15% in that
pass. Under the approved >5% exception, production therefore retains the private
cache wrapper, its operation size, and recursive call signatures. Moving the
wrapper into codecutil still removes encodingownership's forwarding helpers.
Operations without an adopter allocate no wrapper; real callbacks incur one
private cache allocation. Ownership discovery is unchanged by this choice.

A local MessagePack improvement skips duplicate adoption for native Array.At
results, which immediately enter encodeValue's registration before any checkpoint.
Host At results retain deep discovery. A controlled before/after run with ten
samples improved native resource materialization from 326,689 to 294,535 ns/op
(-9.8%), and the real-adopter resource control from 198,160 to 172,584 ns/op
(-12.9%), with unchanged B/op and allocs/op. This improvement was retained.
Broader container/iterator ownership redesign and speculative pooling changes
were not added.

#### Profiles and remaining costs

CPU and allocation profiles were collected for d9b3ba82/final MessagePack native
container depth 8 and pre-migration/final JSON flat-array decoding:

```sh
go test ./pkg/engine -run '^$' \
  -bench 'BenchmarkSessionCodecIterable/application/vnd.msgpack/containers/depth_8$' \
  -benchtime=3s -count=1 -cpuprofile=codec.cpu -memprofile=codec.mem -o codec.test
go test ./pkg/encoding/json -run '^$' \
  -bench 'BenchmarkJSONCodecDecode/flat_array_1024$' \
  -benchmem -benchtime=3s -count=1 -cpuprofile=json.cpu -memprofile=json.mem -o json.test
pprof -top -focus='github.com/MontFerret' codec.cpu
pprof -top -alloc_space codec.mem
pprof -top -focus='github.com/MontFerret|encoding/json' json.cpu
pprof -top -alloc_space json.mem
```

The installed standalone pprof was usable even though `go tool pprof` was absent.
Embedded symbols sufficed; passing the Mach-O test binary to this installation
was unsupported. CPU profiles contain substantial runtime/idle samples, so their
application call trees identify costs rather than supplying precise wall-time
attribution. Final MessagePack allocation samples are dominated by map callback
closures (~44%), output buffer growth (~21%), resource registration (~15%),
Object.ForEach boxing (~9%), internal collection growth (~5%), and leaf custom
marshaler buffers (~4%). Those are existing buffering/serialization/ownership
costs, not a constant-memory claim. Its CPU tree confirms map iteration,
identity lookup, encoding, and stored discovery remain the relevant paths.

JSON's pre-migration decoder used goccy/go-json; the migration uses the standard
encoding/json tokenizer and validates complete documents, including trailing
roots, separators and truncation. Flat-array allocation counts rise from 1,812
to 2,834, with B/op rising from 90,120 to 93,936. Profiles place about 63% of final
allocation space in materialized Array growth and about 29% under Decoder.Token,
including token strings and tokenizer buffers. The pre-migration profile instead
shows goccy token/buffer work. The tokenizer replacement and stricter document
validation were introduced by the initial migration; this follow-up does not
change accepted serialization or decoder structure. Their cost is included in
the five-revision comparison, not attributed to adoption.

#### Native materialization against d9b3ba82

Medians from ten samples; each metric cell is ns/op / B/op / allocs/op.

| Format / yielded values | d9b3ba82 | Final | Time delta |
| --- | ---: | ---: | ---: |
| JSON/scalars | 39,656 / 15,982 / 805 | 51,329 / 15,718 / 801 | +29.44% |
| JSON/resources | 290,762 / 302,197 / 3,914 | 323,175 / 301,929 / 3,910 | +11.15% |
| JSON/containers/depth_1 | 475,634 / 441,617 / 6,987 | 652,570 / 441,618 / 6,986 | +37.20% |
| JSON/containers/depth_8 | 1,771,252 / 1,326,896 / 28,494 | 2,654,482 / 1,328,284 / 28,498 | +49.86% |
| JSON/containers/depth_32 | 7,056,552 / 4,497,724 / 102,226 | 10,193,488 / 4,504,510 / 102,234 | +44.45% |
| MessagePack/scalars | 52,030 / 75,920 / 820 | 56,092 / 75,633 / 815 | +7.81% |
| MessagePack/resources | 295,120 / 435,634 / 3,928 | 295,136 / 435,344 / 3,923 | +0.01% |
| MessagePack/containers/depth_1 | 551,682 / 550,624 / 5,980 | 583,118 / 550,333 / 5,975 | +5.70% |
| MessagePack/containers/depth_8 | 2,304,325 / 1,469,762 / 20,325 | 2,395,440 / 1,469,469 / 20,320 | +3.95% |
| MessagePack/containers/depth_32 | 8,867,428 / 4,621,040 / 69,484 | 9,106,469 / 4,620,781 / 69,480 | +2.70% |

JSON now discovers every yielded graph before hooks, cancellation and prefix writes;
d9b3ba82 performed early stored discovery only in MessagePack. That explains the
larger JSON increase, including direct scalar callback overhead. Root/acquisition
storage scans add one graph walk; ordinary visits remain direct and linear.
Removing the injected built-in ownership hook saves four JSON and five
MessagePack allocations on scalar/resource materialization. Container scratch
and cycle protection offset those savings at deeper JSON nesting. MessagePack
resource timing is unchanged against d9b3ba82; scalar and nested-container costs
remain about 3–8% higher for the stronger callback/error/mutation guarantees.
Hook-heavy graphs can incur repeated deep scans; this is explicitly retained
for cleanup correctness. The historical pass has timing noise (notably JSON
depth 32), so these medians describe this machine and fixture, not universal
percentages. No further unmeasured local improvement was adopted.

#### Common workloads and direct I/O

The common-workload repeats retain all allocation counts against d9b3ba82.
JSON encoding ranges from +0.38% (range) to +8.69% (flat array), with object/deep
cases about +5–8%. MessagePack encoding ranges from +1.11% (nested object) to
+9.52% (flat array). Every recursive visit now checks the cached callback even
when absent, and acquisition/key/hook/error boundaries report discovery. That
extra fixed CPU work remains for the stronger ownership guarantees; the cache
field ablation isolated and rejected its additional cost separately. Decoder
structure is unchanged: repeat timings range from -0.46% to +3.23% for JSON and
-0.54% to +1.51% for MessagePack, without an allocation-count increase.

Compared with the pre-migration API, fixed collecting/adaptation costs remain
visible in tiny encodes: JSON duration is 65.6 → 140.8 ns/op, 112 → 208 B/op,
2 → 3 allocs/op; MessagePack regexp is 59.2 → 144.1 ns/op, 112 → 288 B/op,
2 → 4 allocs/op. JSON flat-array encoding is roughly unchanged in time at
26,086 → 26,307 ns/op while dropping 932 → 9 allocations and
11,152 → 8,277 B/op. MessagePack's flat-array encode is
12,830 → 18,422 ns/op, 8,181 → 8,357 B/op, 8 → 10 allocs/op; its decode is
26,518 → 38,334 ns/op, 24,663 → 24,847 B/op, 771 → 774 allocs/op. Borrowed I/O
adapters, bounded checkpoints, writer-error tracking, strict complete-document
validation and the JSON tokenizer change are part of the migration costs.

Generic direct iterable timing is +0.9–3.2% against d9b3ba82, with unchanged
B/op and allocation counts. Large direct writers retain allocation counts too;
JSON's warmed flat array still uses 0 B/op and 0 allocs/op, while MessagePack
uses 80 B/op and 1 alloc/op. Hook/cancellation variants expose repeated callback
checks and the existing per-key allocations; Buffer variants also include output
copying. The matched 1 MiB Buffer repeat puts MessagePack string at
23,890 → 23,970 ns/op (p=0.315) and binary at 23,926 → 23,978 ns/op (p=0.436),
both with 80 B/op and 1 alloc/op. JSON object/Discard is
235,144 → 244,600 ns/op (+4.02%, p=0.002), with unchanged 4,098 allocations.
The earlier +17%/-30% outliers did not reproduce. The final direct-writer
comparisons range from -2.7% to +13.4%; repeated absent-callback checks around
values, keys and no-op hooks explain remaining CPU overhead without additional
allocations. Scalar payloads do not show a uniform slowdown from byte copying.
Timing changes on unchanged paths are not claimed as optimizations.
