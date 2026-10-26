# Encoder performance recovery

Measurements taken on 2026-10-26 for the JSON and MessagePack encoder changes
following PR #1058. The current API and ownership contract remain in
[Codec I/O](codecs.md); the earlier migration evidence remains in the
[migration report](codec-streaming-migration.md).

## Revisions and method

- Pre-PR: `35693eb380d10bd80de76208b95e7598307e34d6`.
- Merged PR: `e45153dfd59c98b85313f4fa64fb9bf6fa6c854f`.
- First recovery: `a3241fc5`.
- Candidate: the working-tree implementation accompanying this report, closing
  the remaining scalar and MessagePack flat-object gaps.

All performance measurements use Go 1.27.3, Darwin/arm64, Apple M2 Max,
`GOMAXPROCS=12`, and non-race binaries built with `CGO_ENABLED=0`. Archived
revision copies, binaries, profiles, and logs were kept outside the checkout.
Compilation and other validation run separately from timed benchmark rounds.

The original encode benchmarks and their input fixtures were unchanged. Each
package binary runs this command once per round, for ten serial rounds, rotating
the four revisions' execution order:

```sh
export GOMAXPROCS=12 CGO_ENABLED=0
go test -c -o "$revision-$codec.test" ./pkg/encoding/$codec
"./$revision-$codec.test" -test.run '^$' \
  -test.bench 'Benchmark(JSON|Msgpack)CodecEncode$' \
  -test.benchmem -test.benchtime=250ms -test.count=1
benchstat common-before.txt common-after.txt common-prior.txt common-candidate.txt
```

The old API used an internal background context. The merged and recovery byte
benchmarks use `EncodeBytes(b.Context(), ...)`, so comparison with pre-PR includes
the new collecting helper, context validation and cancellation contract. Writer
and helper controls separate those costs. The 146 writer, iterable, helper,
native session and adoption controls compare the candidate against the first
recovery in six serial paired rounds, alternating execution order, with
`-test.benchtime=100ms -test.count=1`:

```sh
# Build/run package binaries for pkg/encoding, pkg/engine and
# pkg/engine/internal/session; append each round to its revision's result file.
"./$revision-$package.test" -test.run '^$' \
  -test.bench 'Benchmark(CodecWriter|CodecIterable|CodecEncodeBytes|SessionCodecIterable|CodecAdoption)$' \
  -test.benchmem -test.benchtime=100ms -test.count=1
benchstat control-prior.txt control-candidate.txt
```

Five additional JSON root-string controls compare the first recovery and
candidate in six serial paired rounds at 150 ms per sample:

```sh
"./$revision-json.test" -test.run '^$' \
  -test.bench '^BenchmarkJSONRootStringFallbacks$' \
  -test.benchmem -test.benchtime=150ms -test.count=1
benchstat final-fallback-prior.txt final-fallback-candidate.txt
```

These fixtures cover ASCII and microsecond durations, short escaped and Unicode
strings, and an ASCII string just beyond the direct-output limit. The new
benchmark file is copied unchanged into the first-recovery tree. Common
benchmarks remain unchanged. The former byte-only API cannot run direct-writer
controls without an adapter, so those controls use the first recovery as their
baseline.

Acceptance evaluates every common workload against pre-PR +5% and checks the
controls against the first recovery; an aggregate speedup does not excuse an
individual regression. `ns/op`, `B/op` and `allocs/op` are compared together.

## Common collecting workloads

All 13 candidate medians meet the per-workload target of at most 1.05 times
the freshly measured pre-PR median. JSON duration is closest at **+4.29%**;
MessagePack’s flat array is **+3.17%**. Every other workload is faster than
pre-PR. The four previously open gaps are JSON duration **+4.29%**, JSON regexp
**−39.18%**, MessagePack regexp **−2.28%**, and MessagePack flat object **−0.88%**.

Each timing is the median of ten samples, followed by benchstat’s rounded 95%
confidence range. A displayed ±0% is a range smaller than its rounding unit.
Acceptance compares medians; confidence ranges describe measurement uncertainty.

| Workload | Pre-PR | Merged PR | First recovery | Candidate | Candidate vs pre |
| --- | ---: | ---: | ---: | ---: | ---: |
| JSON duration | 65.96 ns ±1% | 138.00 ns ±1% | 96.81 ns ±1% | 68.78 ns ±1% | +4.29% |
| JSON regexp | 99.88 ns ±1% | 157.80 ns ±1% | 117.20 ns ±1% | 60.75 ns ±1% | -39.18% |
| JSON range (1,024) | 10.425 µs ±1% | 17.973 µs ±0% | 8.097 µs ±1% | 8.081 µs ±1% | -22.49% |
| JSON flat array (1,024) | 26.51 µs ±2% | 27.02 µs ±1% | 20.83 µs ±1% | 19.92 µs ±0% | -24.84% |
| JSON flat object (256) | 17.02 µs ±1% | 18.35 µs ±1% | 16.60 µs ±1% | 16.55 µs ±1% | -2.71% |
| JSON nested array (depth 10,000) | 309.9 µs ±1% | 434.7 µs ±1% | 284.4 µs ±1% | 282.8 µs ±2% | -8.76% |
| JSON nested object (depth 5,000) | 771.2 µs ±1% | 834.0 µs ±2% | 666.8 µs ±1% | 667.0 µs ±1% | -13.51% |
| MessagePack regexp | 55.77 ns ±1% | 141.85 ns ±1% | 79.18 ns ±3% | 54.50 ns ±1% | -2.28% |
| MessagePack range (1,024) | 17.267 µs ±0% | 22.650 µs ±1% | 8.929 µs ±1% | 8.816 µs ±1% | -48.94% |
| MessagePack flat array (1,024) | 12.79 µs ±0% | 18.32 µs ±1% | 13.18 µs ±0% | 13.20 µs ±0% | +3.17% |
| MessagePack flat object (256) | 11.12 µs ±1% | 13.83 µs ±1% | 11.84 µs ±1% | 11.02 µs ±1% | -0.88% |
| MessagePack nested array (depth 10,000) | 260.5 µs ±1% | 307.5 µs ±1% | 247.5 µs ±0% | 243.2 µs ±1% | -6.62% |
| MessagePack nested object (depth 5,000) | 646.5 µs ±1% | 702.6 µs ±1% | 607.7 µs ±1% | 611.2 µs ±1% | -5.47% |

Memory medians are **B/op / allocs/op**, in the same revision order. Allocation
counts are unchanged from the first recovery on all common workloads. An even
sample count can yield fractional byte medians; small B/op variation includes
amortized allocation effects.

| Workload | Pre-PR memory | Merged PR memory | First recovery memory | Candidate memory |
| --- | ---: | ---: | ---: | ---: |
| JSON duration | 112 / 2 | 208 / 3 | 112 / 2 | 112 / 2 |
| JSON regexp | 112 / 2 | 208 / 3 | 112 / 2 | 112 / 2 |
| JSON range (1,024) | 8,216 / 13 | 8,277 / 9 | 4,210 / 3 | 4,210 / 3 |
| JSON flat array (1,024) | 11,152 / 932 | 8,277 / 9 | 8,181 / 8 | 8,181 / 8 |
| JSON flat object (256) | 12,867 / 422 | 12,457 / 267 | 12,328.5 / 266 | 12,328.5 / 266 |
| JSON nested array (depth 10,000) | 65,520 / 11 | 65,640 / 12 | 65,543 / 11 | 65,543.5 / 11 |
| JSON nested object (depth 5,000) | 630,524.5 / 15,011 | 550,831 / 15,012 | 390,689 / 15,011 | 390,682 / 15,011 |
| MessagePack regexp | 112 / 2 | 288 / 4 | 112 / 2 | 112 / 2 |
| MessagePack range (1,024) | 14,329 / 776 | 14,505 / 778 | 8,181 / 8 | 8,181 / 8 |
| MessagePack flat array (1,024) | 8,181 / 8 | 8,357 / 10 | 8,181 / 8 | 8,181 / 8 |
| MessagePack flat object (256) | 8,261 / 264 | 8,437 / 266 | 8,229 / 264 | 8,229 / 264 |
| MessagePack nested array (depth 10,000) | 32,763 / 10 | 32,939 / 12 | 32,762 / 10 | 32,763 / 10 |
| MessagePack nested object (depth 5,000) | 512,962 / 10,010 | 513,129.5 / 10,012 | 352,901 / 10,010 | 352,897.5 / 10,010 |

Tiny encodes retain **112 B and two allocations**. JSON’s flat integer array
retains the drop from 932 pre-PR allocations to eight, and MessagePack ranges
retain the drop from 776 to eight. The candidate has no significant common
workload slowdown above 5% against `a3241fc5`.

## Writer, helper and ownership controls

All **146 controls** completed six alternating paired rounds against `a3241fc5`.
No timing median increased more than 5%; the largest increase was **+3.62%** for
the JSON object/default/background/Buffer case. Benchstat found no statistically
significant slowdown above 5% (default significance level 0.05). Every control’s
allocation-count median matches the first recovery. B/op increases were all
below **0.025%**; the largest absolute increase was 54 B on a 3.46 MB native
resource workload. These small byte differences are retained in the evidence.

Selected controls include the largest slowdown and the helper, literal, payload
and native paths. Times show 95% confidence ranges; memory is B/op / allocs/op.

| Control | First recovery | Candidate | Change | First recovery memory | Candidate memory |
| --- | ---: | ---: | ---: | ---: | ---: |
| JSON object, default/background/Buffer | 227.8 µs ±1% | 236.0 µs ±0% | +3.62% | 65,620.5 / 4,098 | 65,621 / 4,098 |
| JSON integer, concrete/cancellable helper | 63.12 ns ±1% | 64.44 ns ±1% | +2.08% | 112 / 2 | 112 / 2 |
| MessagePack integer, concrete/cancellable helper | 59.84 ns ±1% | 60.99 ns ±1% | +1.90% | 112 / 2 | 112 / 2 |
| JSON booleans (1,024), concrete/cancellable helper | 15.46 µs ±1% | 15.42 µs ±2% | -0.25% | 16,376 / 9 | 16,376 / 9 |
| JSON nulls (1,024), concrete/cancellable helper | 13.45 µs ±1% | 13.45 µs ±2% | +0.00% | 16,376 / 9 | 16,376 / 9 |
| JSON 1 MiB ASCII to Discard, cancellable | 2611.7 µs ±3% | 744.2 µs ±2% | -71.51% | 0 / 0 | 0 / 0 |
| MessagePack 1 MiB string to Discard, cancellable | 1.798 µs ±0% | 1.781 µs ±1% | -0.95% | 0 / 0 | 0 / 0 |
| Native JSON scalar iterable | 44.92 µs ±0% | 45.88 µs ±1% | +2.13% | 15,718 / 801 | 15,719 / 801 |
| Native MessagePack scalar iterable | 50.51 µs ±1% | 50.42 µs ±1% | -0.16% | 75,553.5 / 814 | 75,554 / 814 |

Concrete, preboxed and pointer scalar helpers retain 112 B and two allocations.
JSON boolean/null collecting arrays retain nine allocations, without the former
per-element allocations. Owned and borrowed iterable
timing changes range from −3.47% to +0.19%. Native session controls range from
−0.39% to +2.13%; adoption controls range from −2.93% to +0.21%, with no
statistically significant adoption slowdown.

The small B/op differences were investigated separately. The MessagePack pool
state grows from 88 to 96 bytes, both within the same 96-byte allocation class;
neither adapter adds a per-operation allocation. JSON's pooled state and map
callback captures are unchanged. Three extra alternating pairs with fixed
iteration counts produced 4,082 B / 7 allocations on both revisions for the
MessagePack boolean array (`100000x`), and 65,622 versus 65,623 B / 4,098
allocations for JSON object-to-Discard (`5000x`). This supports amortized
measurement effects rather than a new fixed allocation. These three-pair runs
are allocation diagnostics, not additional throughput acceptance samples or
confidence estimates:

```sh
"./$revision-encoding.test" -test.run '^$' \
  -test.bench '^BenchmarkCodecWriter$/^application$/^json$/^object$/^default$/^background$/^Discard$' \
  -test.benchmem -test.benchtime=5000x -test.count=1
"./$revision-encoding.test" -test.run '^$' \
  -test.bench '^BenchmarkCodecEncodeBytes$/^msgpack$/^boolean_array_1024$/^concrete_value$/^background$' \
  -test.benchmem -test.benchtime=100000x -test.count=1
```

The five added JSON root-string controls also meet the 5% limit against the
first recovery, with unchanged allocation counts and byte medians. Short Unicode
strings retain a measured +3.36% cost; this is reported rather than hidden by the
large ASCII improvement.

| JSON root string | First recovery | Candidate | Change | Both revisions’ memory |
| --- | ---: | ---: | ---: | ---: |
| Microsecond duration | 85.31 ns ±1% | 86.91 ns ±1% | +1.88% | 112 / 2 |
| ASCII duration | 83.01 ns ±1% | 59.55 ns ±1% | -28.27% | 112 / 2 |
| Short escaped string | 71.89 ns ±1% | 71.72 ns ±1% | -0.24% | 112 / 2 |
| Short Unicode string | 75.24 ns ±1% | 77.77 ns ±1% | +3.36% | 112 / 2 |
| 4,090-byte ASCII string | 10.477 µs ±2% | 3.263 µs ±1% | -68.85% | 4,146 / 2 |

The serial rounds produced 520 common-workload samples, 1,752 original-control
samples and 60 additional fallback-control samples. Raw outputs, comparisons,
validation logs and the runner are retained locally under
`/private/tmp/ferret-encoder-perf.secSu3/gap-closure/verified/`.

## Causes and implementation

The owning subsystem is `pkg/encoding`. JSON/MessagePack `codec.go` and
`config.go` own immutable configuration and operation setup; their `encoder.go`
files own traversal and serialization. `internal/codecutil/io.go` owns shared
checked writes and cancellation sampling. New helper benchmarks and reuse fixtures live directly under `pkg/encoding`;
codec-specific configuration, numeric, string and payload tests stay with their
codec. Shared buffer-writer accounting tests live with `codecutil`.

The first recovery addressed several distinct costs; the regression was not an
ownership graph scan on every ordinary direct encode.

- Passing the former 96-byte concrete codec to the non-inlined `EncodeBytes`
  helper allocated an interface copy. Pointer and preboxed controls removed
  exactly 96 bytes and one allocation. Each codec now holds one pointer to
  immutable hook configuration; configuration builders publish independent
  snapshots and zero-value codecs remain valid.
- MessagePack allocated an escaping 80-byte checked writer on each operation.
  Its codec-local pool now owns both writer and dependency encoder. Every exit
  resets dependency flags/dictionaries and releases borrowed writer, context,
  callback and error references. Returned output buffers are never pooled.
- CPU profiles placed substantial time in the added checked writes and JSON's
  extra recursive dispatch frame. JSON now uses one recursive dispatcher; both
  codecs borrow immutable encoder configuration and skip empty hook calls.
  The sampled traversal check remains out of line so the counter increment can
  inline into hot loops. Cancellation still checks every 256 traversal values
  and at bounded I/O/payload boundaries.
- JSON's `[]byte(text)` literal conversion allocated once per boolean/null.
  Literals now use reusable scratch. Numeric output can append into an exact
  `bytes.Buffer`'s available capacity. Hook-free native ranges use bounded
  numeric chunks only for exact buffers, without changing arbitrary writer
  failure timing or configured hook traversal.
- MessagePack dispatches exact native types without redundant custom-protocol
  assertions. Hook-free ranges without an adopter avoid boxing each generated
  integer. Host custom encoders retain precedence over collection capabilities.
- The first complete controls pass caught a 6% slowdown for cancellable
  MessagePack payloads sent to Discard. Large writes now loop directly over
  checked chunks instead of recursively re-entering the small-write dispatcher.
  A paired six-sample ablation confirmed the affected payload cases became
  approximately 9% faster than the merged PR, preserving buffer-writer results.

### Closing the remaining gaps

This follow-up changes JSON's `codec.go`, `encoder.go` and `string_ascii.go`,
MessagePack's `codec.go` and `encoder.go`, and shared
`internal/codecutil/io.go`. New scalar, string, map and checked-writer tests stay
with those owners; shared reuse tests and the new JSON fallback benchmarks cover
the operation boundaries.

The remaining scalar cases paid operation and pool setup for very small outputs.
The MessagePack object profile also showed repeated checked-writer dispatch and
native-value traversal costs. The candidate adds these bounded optimizations:

- Exact root String, Duration and non-nil Regexp values can write directly to an
  exact `bytes.Buffer` only without encoder hooks or an adopter. Argument
  validation, root discovery and the initial cancellation check still precede
  this path; it performs a traversal step and the final cancellation check.
  JSON accepts unescaped ASCII up to 4,089 bytes; MessagePack accepts strings
  below 32 bytes using its fixstr representation. Other writers retain ordinary
  synchronous error and partial-prefix behavior.
- JSON renders an eligible native root value once. If ASCII or size eligibility
  fails, the already-rendered text goes through the ordinary bounded string
  writer. This avoids formatting a duration twice when it contains `µ`.
  The generic writer copies contiguous safe ASCII runs for strings of at least
  64 bytes while preserving escaping, invalid-UTF-8 replacement and its existing
  bounded flush/checkpoint behavior.
- MessagePack chooses its concrete-buffer adapter once per operation. The
  adapter embeds the checked writer by value, so both generic and concrete
  paths share exactly one operation, error and byte counter. Small writes call
  the buffer directly; large writes delegate to the existing chunk loop.
  Configured and custom encodes can use the adapter. Ignored custom writer
  errors are captured from either codec-owned writer type before post-hooks.
- MessagePack keeps dependency string encoding for ordinary values and keys,
  preserving dictionaries and interning changed by custom encoders. Acquisition
  uses `ResetWriter` because the private pool performs full `Reset(nil)` and
  clears the complete writer state on every release, including failure or panic.
- Exact native string keys avoid a dynamic `String()` call; host keys still
  invoke their own method. Exact integer map values use sampled traversal and
  the dependency's `EncodeInt` when hooks and adoption are absent. Key writes
  retain their order, and all other values use the ordinary visitor. Acquisition
  reporting remains unconditional so the callback closure stays compact.

Root discovery, acquisition discovery, mutation/failure rescans, iterator
ownership, original hook contexts, error identities, and borrowed I/O lifetimes
are unchanged. No public method signatures, accepted formats, or FQL semantics
changed. The native ownership path is measured separately; it does not pay a
fresh concrete codec interface allocation on each operation.

## Validation and review

Validation used task-scoped writable `GOCACHE`, `GOTMPDIR` and
`STATICCHECK_CACHE`. These commands passed for the candidate:

```sh
CGO_ENABLED=1 go test -race ./pkg/encoding/... ./pkg/engine/... ./uapi/... -count=1
GOTOOLCHAIN=go1.25.14 go test ./pkg/encoding/... -count=1
GOTOOLCHAIN=go1.26.9 GOFLAGS=-count=1 make test
GOTOOLCHAIN=go1.26.9 make fmt
GOTOOLCHAIN=go1.26.9 revive -config revive.toml -formatter stylish \
  -exclude ./pkg/parser/fql/... -exclude ./vendor/... \
  -exclude './*_test.go' -exclude ./bin/... ./...
git diff --check
```

The full Go 1.26 suite includes package and integration race tests, both
independent API-tool modules, compatibility, the public facade and security
tests. The Go 1.27.3 race run covers codecs, engine and UAPI. The minimum-version
codec suite uses Go 1.25.14. No decoder or FQL behavior changed.

Added coverage checks scalar eligibility and fallback equivalence, zero-value,
pointer and decoder-only codecs, borrowed buffer prefixes and consumed offsets,
retained output, typed-nil validation, discovery-triggered cancellation, and
original hook contexts. JSON cases include duration extremes and microseconds,
escaping, invalid UTF-8, limit boundaries, large ASCII runs, short writes and
multi-chunk cancellation. MessagePack cases cover fixstr boundaries, mixed
base/adapter writes, shared sticky errors and chunk remainder accounting,
same-operation custom string dictionaries/interning, and dependency state after
pool reuse. Map tests preserve host-key methods, custom scalar encoding,
signed-64-bit limits, integer hook/adoption visits and cancellation before the
larger I/O checkpoint. Concurrent encodes alternate pooled containers and scalar
shortcuts. Existing reentrant and resource-ownership tests remain in the suite.

The final self-review and independent complete-diff review checked the production
paths, all added tests, and documentation. Review strengthened the ignored-writer
cancellation test to verify the ancestor array's post-hook and exact partial
header, then reran the MessagePack package on all three toolchains:

```sh
CGO_ENABLED=1 go test -race ./pkg/encoding/msgpack -count=1
GOTOOLCHAIN=go1.25.14 go test ./pkg/encoding/msgpack -count=1
GOTOOLCHAIN=go1.26.9 CGO_ENABLED=1 go test -race ./pkg/encoding/msgpack -count=1
```

Review and measurement corrected supporting regressions before acceptance:
a cached map-adopter boolean enlarged each callback allocation by 16 bytes and
was removed; JSON fallback formatted Unicode durations twice and now reuses the
rendered text; scalar eligibility probing slowed integer roots and is now gated
by exact root type. The concrete-buffer adapter embeds its sole writer state,
avoiding another pointer load without duplicating ownership or error state.
JSON performs its bounded root write directly in the operation entry point;
its buffer slice is sized once after `Grow`.

The previously documented tooling limitations were rechecked separately from
codec results:

- `GOFLAGS=-count=1 make test` on Go 1.27.3 passed root package race tests but
  reproduced a race in `go/types` through `golang.org/x/tools/go/packages` in the
  independent API-reference analyzer test. The earlier investigation reproduced
  the same failure on the archived merged-PR revision. The complete candidate
  suite passed on Go 1.26.9.
- `make fmt` and `make lint` with Go 1.27 could not consume export-data version 5
  with the installed tools. Formatting passed with Go 1.26.9. All three
  Staticcheck stages of `GOTOOLCHAIN=go1.26.9 make lint` passed; Revive then
  reported 216 pre-existing findings only in the ignored
  `bin/website-output-docs/core-snapshot` generated parser. Revive passed with
  that ignored directory excluded, as above.

Repository codec and performance documentation were updated. Public
documentation remains applicable because signatures, accepted encodings,
ownership, cancellation and partial-output contracts are unchanged.
