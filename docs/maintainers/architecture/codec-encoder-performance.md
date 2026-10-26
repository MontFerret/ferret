# Encoder performance recovery

Measurements taken on 2026-10-26 for the JSON and MessagePack encoder changes
following PR #1058. The current API and ownership contract remain in
[Codec I/O](codecs.md); the earlier migration evidence remains in the
[migration report](codec-streaming-migration.md).

## Revisions and method

- Pre-PR: `35693eb380d10bd80de76208b95e7598307e34d6`.
- Merged PR: `e45153dfd59c98b85313f4fa64fb9bf6fa6c854f`.
- Fixed: the working-tree implementation accompanying this report.

All performance measurements use Go 1.27.3, Darwin/arm64, Apple M2 Max,
`GOMAXPROCS=12`, and non-race binaries built with `CGO_ENABLED=0`. Archived
revision copies, binaries, profiles, and logs were kept outside the checkout.
No compilation or other validation ran concurrently with the final benchmarks.

The original encode benchmarks and their input fixtures were unchanged. Each
package binary ran this command once per round, for ten rounds, rotating the
three revisions' execution order:

```sh
export GOMAXPROCS=12 CGO_ENABLED=0
go test -c -o "$revision-$codec.test" ./pkg/encoding/$codec
"./$revision-$codec.test" -test.run '^$' \
  -test.bench 'Benchmark(JSON|Msgpack)CodecEncode$' \
  -test.benchmem -test.benchtime=250ms -test.count=1
benchstat verified-before.txt verified-after.txt verified-final.txt
```

The old API used an internal background context. The merged and fixed byte
benchmarks use `EncodeBytes(b.Context(), ...)`, so comparison with pre-PR includes
the new collecting helper, context validation and cancellation contract. Writer
and helper controls separate those costs. Additional writer, iterable, helper,
native session and adoption controls ran in six serial paired rounds, alternating
merged/fixed order, with `-test.benchtime=100ms -test.count=1`:

```sh
# Build/run package binaries for pkg/encoding, pkg/engine and
# pkg/engine/internal/session; append each round to its revision's result file.
"./$revision-$package.test" -test.run '^$' \
  -test.bench 'Benchmark(CodecWriter|CodecIterable|CodecEncodeBytes|SessionCodecIterable|CodecAdoption)$' \
  -test.benchmem -test.benchtime=100ms -test.count=1
benchstat controls-after.txt controls-final.txt
```

The added `BenchmarkCodecEncodeBytes` fixtures were copied unchanged into the
merged-PR benchmark tree. The former byte-only API cannot run the direct-writer
controls without an adapter, so those controls compare merged versus fixed.

## Common collecting workloads

Medians of ten samples. Times are ns/op; the last column is B/op / allocs/op.
Relative changes compare fixed with each baseline. The 5% target is assessed
per workload rather than from a geomean.

| Workload | Pre-PR | Merged PR | Fixed | Fixed vs pre | Fixed vs merged | Fixed memory |
| --- | ---: | ---: | ---: | ---: | ---: | ---: |
| JSON duration | 67.58 | 142.45 | 99.96 | +47.91% | -29.83% | 112 / 2 |
| JSON regexp | 102.25 | 164.50 | 121.20 | +18.53% | -26.32% | 112 / 2 |
| JSON range_1024 | 10,677.50 | 18,506.00 | 8,351.00 | -21.79% | -54.87% | 4,210 / 3 |
| JSON flat_array_1024 | 27,287.00 | 27,981.50 | 21,592.00 | -20.87% | -22.83% | 8,181 / 8 |
| JSON flat_object_256 | 17,581.00 | 18,988.50 | 17,111.50 | -2.67% | -9.88% | 12,328 / 266 |
| JSON nested_array_10000 | 319,064.00 | 448,243.50 | 295,463.50 | -7.40% | -34.08% | 65,545 / 11 |
| JSON nested_object_5000 | 801,207.50 | 868,075.50 | 689,491.50 | -13.94% | -20.57% | 390,686 / 15,011 |
| MessagePack regexp | 58.46 | 145.55 | 81.46 | +39.34% | -44.03% | 112 / 2 |
| MessagePack range_1024 | 17,858.00 | 23,292.50 | 9,195.50 | -48.51% | -60.52% | 8,181 / 8 |
| MessagePack flat_array_1024 | 13,146.50 | 18,863.00 | 13,581.00 | +3.31% | -28.00% | 8,181 / 8 |
| MessagePack flat_object_256 | 11,402.50 | 14,322.00 | 12,293.50 | +7.81% | -14.16% | 8,229 / 264 |
| MessagePack nested_array_10000 | 269,318.50 | 316,284.00 | 257,891.00 | -4.24% | -18.46% | 32,762 / 10 |
| MessagePack nested_object_5000 | 670,334.00 | 732,685.50 | 630,008.50 | -6.02% | -14.01% | 352,897 / 10,010 |

The 5% recovery target is met by nine of thirteen common workloads, including
nine of ten container workloads. It is **not fully met**: JSON duration and
regexp retain about 32 ns and 19 ns of overhead, MessagePack regexp about 23 ns,
and MessagePack's flat object is 7.81% slower than pre-PR. These paths retain
the new operation validation, context/discovery setup, checked writes and pool
lifecycle. Those costs are not waived by the improved aggregate result.

The fixed implementation is 10–61% faster than the merged PR on every common
workload. JSON's collecting integer array retains the migration's allocation
improvement (932 pre-PR, 9 merged, 8 fixed). Tiny JSON and MessagePack workloads
both return to 112 B/op and two allocations. MessagePack ranges fall from
776 pre-PR allocations (778 merged) to eight by avoiding generated-integer
interface boxes on the no-hook/no-adopter path. Smaller immutable encoder
captures reduce deep-object memory to approximately 391 KB for JSON and
353 KB for MessagePack, without removing ownership checks.

## Writer, helper and ownership controls

All 146 control cases completed six paired samples per revision. Benchstat found
no statistically significant timing regression against the merged PR. Selected
medians follow; these helper and writer cases use cancellable contexts. Memory
columns show B/op / allocs/op.

| Control | Merged ns/op | Fixed ns/op | Change | Merged memory | Fixed memory |
| --- | ---: | ---: | ---: | ---: | ---: |
| JSON integer, concrete `EncodeBytes` | 95.16 | 65.32 | -31.36% | 208 / 3 | 112 / 2 |
| MessagePack integer, concrete `EncodeBytes` | 106.75 | 61.44 | -42.44% | 288 / 4 | 112 / 2 |
| JSON boolean array, 1,024 values | 31,846.50 | 15,875.00 | -50.15% | 24,668 / 1,034 | 16,376 / 9 |
| JSON null array, 1,024 values | 29,157.50 | 13,779.50 | -52.74% | 24,668 / 1,034 | 16,376 / 9 |
| MessagePack 1 MiB string to Discard | 2,049.50 | 1,844.50 | -10.00% | 80 / 1 | 0 / 0 |
| MessagePack 1 MiB binary to Discard | 2,030.00 | 1,845.00 | -9.11% | 80 / 1 | 0 / 0 |

Concrete, preboxed-interface and pointer helper controls now all use 112 B and
two allocations for a scalar in both codecs. The JSON literal workloads lose
1,024 per-element allocations plus the concrete codec allocation; their nine
remaining allocations collect the result. The existing configured/background,
configured/cancellable, reusable-buffer and Discard controls remain in place.

Owned and borrowed iterable controls improve by 12–17%. Native session scalar
materialization improves by about 10% in both codecs. Nested resource-container
materialization medians decrease by 3.9–8.1%; the JSON depth-one decrease and
flat resource results are statistically unchanged. Adoption controls with
absent, no-op and real result callbacks have
no significant timing regression. These results preserve the migration's
streaming and ownership behavior while measuring their costs separately from
concrete helper boxing.

## Causes and implementation

The owning subsystem is `pkg/encoding`. JSON/MessagePack `codec.go` and
`config.go` own immutable configuration and operation setup; their `encoder.go`
files own traversal and serialization. `internal/codecutil/io.go` owns shared
checked writes and cancellation sampling. MessagePack's now-redundant custom
protocol predicate was removed from `helpers.go`. The registry test fixture was
realigned by the formatter after the embedded codec's representation changed.
New helper benchmarks and reuse fixtures live directly under `pkg/encoding`;
codec-specific configuration, numeric and payload tests stay with their codec.

The regression had several distinct costs; it was not an ownership graph scan
on every ordinary direct encode.

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

Root discovery, acquisition discovery, mutation/failure rescans, iterator
ownership, original hook contexts, error identities, and borrowed I/O lifetimes
are unchanged. No public method signatures, accepted formats, or FQL semantics
changed. The native ownership path is measured separately; it does not pay a
fresh concrete codec interface allocation on each operation.

## Validation and review

Validation used task-scoped writable `GOCACHE`, `GOTMPDIR` and
`STATICCHECK_CACHE`. The following passed:

```sh
go test -race ./pkg/encoding/... ./pkg/engine/... ./uapi/... -count=1
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
tests. New regressions cover configuration snapshots, concurrent/reentrant
encodes, retained output and error prefixes, reuse after failure/cancellation,
buffer-wrapper overrides, custom MessagePack encoder precedence, numeric
capacity boundaries, and multi-chunk JSON ranges. Review added explicit
second-chunk error/short-write and post-first-chunk cancellation tests, checking
exact prefixes, ancestor hook errors and no writes after failure. All affected
validation was repeated after the large-write correction.

Validation limitations were kept separate from codec failures:

- `make test` on Go 1.27.3 passed root package race tests but encountered a race
  in `go/types` through `golang.org/x/tools/go/packages` in the independent
  API-reference analyzer test. The same failing test reproduced on the archived
  merged-PR revision. The complete suite passed on Go 1.26.9.
- Installed formatting/static-analysis tools could not consume Go 1.27 export
  data. Formatting passed with Go 1.26.9. All three Staticcheck stages of
  `GOTOOLCHAIN=go1.26.9 make lint` passed; Revive then reported 216 pre-existing
  findings only in the ignored `bin/website-output-docs/core-snapshot` generated
  parser. Revive passed with that ignored directory excluded, as above.

The complete implementation and tests received independent second-pass review,
including the final bounded-range change, with no remaining actionable
correctness, ownership or API findings. Performance review drove the additional
hook-call, receiver-copy, range and payload-loop improvements; the missing
mid-payload failure coverage was corrected. Repository architecture and
performance documentation were updated; public documentation needs no change
because public behavior and contracts are preserved.
