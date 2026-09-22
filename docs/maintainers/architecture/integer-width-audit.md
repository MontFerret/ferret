# Integer-width audit

Ferret Int follows the [signed 64-bit runtime rule](runtime.md#integer-width).
This audit covers handwritten production code, compatibility adapters, tool
modules, runtime/public boundaries, and the tests that exercise those boundaries.
Generated parser code and vendored dependencies are not edited.

## Corrected value paths

| Owner | Correction | Regression coverage |
| --- | --- | --- |
| Runtime integer conversion | `ParseInt` parses decimal strings as int64; `ToInt(DateTime)` preserves Unix seconds; `Int.Unwrap` returns Go int64. | `pkg/runtime/int_width_test.go`, compatibility unwrapping, SDK decoding, public parameter round trips |
| Runtime Float helpers | Preserve the existing integer-string grammar of `ParseFloat` with explicit int64 parsing; normalize `IsInf`'s sign without narrowing arbitrary Int values. | Runtime boundary and sign tests |
| JSON | Integer-form tokens within int64 remain Int, including nested values. Existing Float conversion outside that range remains. | JSON boundary round trips and public FQL tests |
| Compiler recovery | Retry count parsing, storage, and emitted constants use int64. | Success after one failure and oversized-literal rejection at None, Basic, and Full optimization |
| Math | Int inputs to floor/ceil/round retain their exact type and value; Float results must be finite and fit int64. | Numeric boundaries, non-finite values, and FQL error recovery |
| DateTime | Fixed units use checked runtime Duration and DateTime operations; calendar shifts validate native components, normalization, and epoch range before accepting AddDate results. | Beyond-int32 offsets, overflow, calendar boundaries, DST gaps/overlaps, month ends, and large/negative years |
| MessagePack | List, Map, and Range headers check nonnegative native-int and uint32 bounds before encoding. | Boundary checks and rejection before headers or traversal |
| Bytecode/VM | Aggregate selector constants must fit native int before becoming internal key indexes, including directly constructed VM programs. | Bytecode validation and VM execution boundaries |

`Unwrap() any` retains its signature, but the concrete Go representation of
Int is now `int64` on all hosts. No repository production consumer requires the
old `int` representation. SDK decoding to native int now reaches its existing
overflow validation rather than accepting a narrowed unwrapped value.

New rounding and datetime range failures are intentional FQL-visible changes.
The corresponding deprecated global adapters retain the same behavior and
error metadata as their canonical functions.

## Optional allocation hints

SDK `ToSlice`, array concat/union/flatten/legacy unshift, object list copying,
VM flatten/distinct, and object-comparison snapshots omit hints that cannot fit
native int or whose sizing arithmetic would overflow. Traversal, value order,
ownership, and source errors remain authoritative. Representable hints remain;
this is not a memory budget or protection against all allocation failures.

These paths share `runtime.CapacityHint`. Mandatory native-integer boundaries
introduced by this audit share `runtime.ToNativeInt`; MessagePack's format
limit, aggregate-selector diagnostics, and calendar arithmetic remain local.
The helpers do not change existing SDK decoding or allocation-constructor policy.

Existing checked SDK slice decoding and required MessagePack sizes retain
errors rather than adopting the optional-hint policy.

## Native integer uses retained

- Compiler registers, instruction operands, function bindings, aggregate slots,
  debug indexes, arities, and in-memory counters represent native structures.
  Aggregate sort indexes are constructed directly from the entries being sorted.
- Concrete Array/Object/FastObject lengths, stream-group descriptors, SDK view
  indexes, and debugger snapshots are bounded by their native backing storage.
- String find/slice bounds are clamped before conversion; replacement/split
  limits are bounded without changing observable results for an in-memory string.
  Repeat counts and crypto token lengths already have explicit limits.
- Formatter widths and placeholder indexes, stdlib assertion arities parsed by
  the API-reference tool, port numbers, and bytecode artifact metadata are host,
  format, or structurally bounded concepts rather than arbitrary FQL numbers.
- Existing literal parsing, debugger literal evaluation, bytecode value codecs,
  MessagePack integer decoding, host ValueOf/SDK Encode, and random integer
  generation already preserve signed 64-bit values.

## Separate design follow-ups

- `Array.CopyWithGrowth` still accepts an Int capacity and narrows it in
  `copyInternal`. Its production caller supplies the bounded constant 1.
  Supporting arbitrary host-supplied growth with returned errors requires a
  decision about this no-error public API.
- `NewArray64` and other no-error allocation constructors retain Go allocation
  failure behavior. Their signatures are unchanged; fallible runtime callers
  corrected here no longer send them unrepresentable optional hints.
- General Float-to-Int overflow behavior outside the corrected rounding
  functions remains unchanged. Defining a common conversion policy is separate.
- Extreme DateTime calendar accessors such as Go `Time.Year() int` have their own
  native-width representation limit. Extending those APIs beyond native calendar
  fields requires a DateTime design decision, not integer narrowing at their
  consumers. Calendar arithmetic now rejects a source year that cannot survive
  this boundary, including a zero-unit shift.

## Benchmark evidence (2026-09-22)

Compared the working implementation with an archive of `ef430a70`, copying the
new benchmark files unchanged into the baseline. Measurements used Go 1.26.5,
darwin/arm64, Apple M2 Max. Six alternating before/after pairs ran this command
in each source tree, followed by `benchstat` over the two output files:

```sh
go test ./pkg/encoding/json ./pkg/encoding/msgpack ./pkg/sdk \
  ./pkg/stdlib/datetime ./pkg/stdlib/arrays ./pkg/stdlib/objects \
  ./pkg/vm ./pkg/vm/internal/data ./test/benchmarks \
  -run '^$' \
  -bench 'Benchmark(JSONCodecDecode|MsgpackCodecEncode|DecodeNested|DateTime|ArraySets|IntegerWidth.*|GenericGroupedCollectAggregate_None)$' \
  -benchmem -count=1 -benchtime=200ms
```

Representative medians; `IntegerWidth` prefixes are omitted below:

| Path | Before | After | Allocation impact |
| --- | ---: | ---: | --- |
| SDK DecodeNested | 2.624 us | 2.508 us | Unchanged: 2160 B, 40 allocations |
| SDK ToSlice | 1.072 us | 1.047 us | Unchanged: 1048 B, 2 allocations |
| Calendar/hour | 67.05 ns | 129.35 ns | 64 to 104 B; 3 to 6 allocations |
| Calendar/day | 82.48 ns | 124.75 ns | Unchanged: 64 B, 3 allocations |
| Calendar/month | 83.80 ns | 123.15 ns | Unchanged: 64 B, 3 allocations |
| Calendar/year | 82.45 ns | 123.85 ns | Unchanged: 64 B, 3 allocations |
| ArrayHints/concat | 1.145 us | 1.105 us | Unchanged |
| ArrayHints/flatten | 1.037 us | 1.045 us | Unchanged |
| Object ListCopy | 989 ns | 962 ns | Unchanged |
| VM ArrayHints/flatten | 1.065 us | 1.072 us | Unchanged |
| VM ArrayHints/distinct | 4.328 us | 4.287 us | Unchanged |
| ObjectComparison | 28.99 us | 29.20 us | Unchanged |
| GenericGroupedCollectAggregate_None | 4.169 ms | 4.076 ms | Unchanged |

The fixed-unit overhead comes from the required checked runtime multiplication
and DateTime operation. Calendar validation preserves the ordinary path's
allocation count; extreme dates use exact arithmetic for range validation.
Calendar timings increased by approximately 47-51%, and direct hour shifting
by 93%. The legacy registered `date_add` hour benchmark increased from 53.68 ns
to 112.70 ns, with the same additional 40 B and 3 allocations.

SDK and collection timing differences were not statistically significant in
these samples. Several samples overlapped runner activity and showed high
variance; these results do not establish small performance improvements.

Codec measurements were repeated in six alternating pairs with longer samples
because the first JSON samples included large timing outliers:

```sh
go test ./pkg/encoding/json ./pkg/encoding/msgpack -run '^$' \
  -bench 'Benchmark(JSONCodecDecode|MsgpackCodecEncode)$/(flat_array_1024|flat_object_256|nested_array_10000|nested_object_5000)$' \
  -benchmem -count=1 -benchtime=300ms
```

| Codec path | Before | After | Allocation impact |
| --- | ---: | ---: | --- |
| JSON flat array, 1024 items | 78.87 us | 77.02 us | Unchanged |
| JSON flat object, 256 items | 45.66 us | 46.70 us | Unchanged |
| JSON nested array, depth 10000 | 1.097 ms | 1.087 ms | Unchanged |
| JSON nested object, depth 5000 | 1.197 ms | 1.197 ms | Unchanged |
| MessagePack flat object, 256 items | 14.09 us | 15.14 us | Unchanged |

JSON timing differences remained statistically insignificant. MessagePack flat
object encoding was 3.9% slower in the broad comparison and 7.5% slower in the
repeat (`p=0.041` and `p=0.026`, respectively). Its mandatory header validation is
retained. Timing variance limits precise estimates; no latency-neutrality claim
is made. MessagePack allocation counts stayed unchanged across all sampled paths.

## Validation and final review (2026-09-22)

All required local checks passed:

- Focused runtime, SDK, codec, compiler, stdlib, compatibility, bytecode, and VM
  tests, including public FQL regressions at None, Basic, and Full optimization.
- `make test` on darwin/arm64 with Go 1.26.5, including race-enabled unit and
  integration coverage and both tooling modules.
- `make compile-32bit` after the final production correction.
- Actual `make test-32bit` execution with Go 1.25.14 on Linux/386 via Docker
  emulation, covering the root module and both tooling modules. This was an
  executed test run, not just cross-compilation.
- `go test -race ./pkg/stdlib/datetime . -count=1` after the final extreme-source
  calendar correction; the final complete Linux/386 run also includes it.
- `make fmt`, `make lint`, `git diff --check`, new-file whitespace checks, and
  local links in the affected maintainer documentation.
- Website `mage build` after synchronizing host-value and serialization guides.

Go caches and temporary build files used task-scoped writable directories under
`$TMPDIR`. The final Linux runner command was:

```sh
docker run --rm --platform linux/386 \
  --mount "type=bind,src=$PWD,dst=/src,readonly" \
  --mount "type=bind,src=$TMPDIR,dst=/audit" -w /src \
  -e GOCACHE=/audit/ferret-int-audit-linux386-cache \
  -e GOMODCACHE=/audit/ferret-int-audit-linux386-mod \
  golang:1.25.14-bookworm sh -c \
  'apt-get update >/dev/null && apt-get install -y jq >/dev/null && setpriv --bounding-set=-dac_override,-dac_read_search make test-32bit'
```

The runner needed `jq` for script tests and ordinary filesystem permission
enforcement for the publishing tool's write-failure regression. Earlier runner
attempts without those prerequisites failed; the final complete run passed.
The existing CI matrix remains responsible for additional Go versions.

The mandatory full-diff self-review covered production changes, new tests and
benchmarks, documentation, ownership, and public behavior. Review corrections
included the aggregate-selector narrowing guard, rejection of truncated source
calendar years, synchronized legacy range-error metadata, and the JSON test's
stale native-integer overflow boundary. Relevant checks were rerun afterward.
Unwrapping and rounding regressions were also verified to fail against the
original source. No outstanding correctness finding remains within this scope;
the separate API design questions above remain deferred.

## Shared-helper follow-up (2026-09-22)

`runtime.ToNativeInt` now owns exact checked narrowing, and
`runtime.CapacityHint` owns optional capacity arithmetic. The duplicate array
helpers and equivalent SDK/object/VM inline checks were removed. Callers retain
their fallback capacities, domain restrictions, diagnostics, and ownership.
Calendar operations pass the checked native deltas directly to `AddDate`.
Existing SDK decoding checks, rounding, and allocation APIs are unchanged.

Direct helper tests now live in `pkg/runtime/native_int_test.go`, covering native
and int64 boundaries, exact conversions, rejected conversions returning zero,
invalid hint parameters, and multiplication/addition overflow without large
allocations. Subsystem regressions continue to cover traversal and mandatory
size failures.

The benchmark baseline was a snapshot of the completed, staged integer audit
immediately before this refactor, including its existing benchmark files. Both
trees used Go 1.26.5 on darwin/arm64, Apple M2 Max, with this command:

```sh
go test ./pkg/sdk ./pkg/stdlib/arrays ./pkg/stdlib/objects \
  ./pkg/stdlib/datetime ./pkg/vm ./pkg/vm/internal/data \
  ./pkg/encoding/msgpack ./test/benchmarks -run '^$' \
  -bench 'Benchmark(IntegerWidth.*|MsgpackCodecEncode|GenericGroupedCollectAggregate_None)$' \
  -benchmem -count=6 -benchtime=200ms
```

Allocation counts were unchanged across every measured path. Representative
medians, with the `IntegerWidth` prefix omitted:

| Path | Before | After |
| --- | ---: | ---: |
| SDK ToSlice | 1.007 us | 1.008 us |
| Object ListCopy | 916.3 ns | 922.4 ns |
| Calendar/day | 117.0 ns | 118.4 ns |
| Calendar/month | 117.6 ns | 119.5 ns |
| VM ArrayHints/flatten | 929.4 ns | 984.6 ns |
| VM ArrayHints/distinct | 4.145 us | 3.940 us |
| MessagePack flat object, 256 items | 12.85 us | 12.60 us |
| GenericGroupedCollectAggregate_None | 3.981 ms | 3.767 ms |

An initial approximately 3% array timing increase was investigated. Compiler
diagnostics confirmed that both helpers inline at the affected call sites.
Six alternating before/after pairs repeated the focused benchmark:

```sh
go test ./pkg/stdlib/arrays -run '^$' \
  -bench '^BenchmarkIntegerWidthArrayHints$' \
  -benchmem -count=1 -benchtime=500ms
```

Concat measured 1.073 us before and 1.077 us after (`p=0.390`); flatten measured
1.000 us before and 1.011 us after (`p=0.589`). Neither repeat showed a
statistically significant change, and allocation counts remained unchanged.
These samples establish no repeatable array slowdown and are not a guarantee
of identical timing on every host.

Follow-up validation passed: focused tests for all migrated packages,
`make test`, `make compile-32bit`, and actual `make test-32bit` execution across
the root and both tooling modules. The Linux/386 runner used Go 1.25.14 and the
same Docker command documented above. `make fmt`, `make lint`, `git diff --check`,
new-file whitespace checks, maintainer links, and the website's `mage build`
also passed.

The final self-review covered the complete refactor, the new public contracts,
boundary tests, error preservation, ownership, and benchmark evidence. No
outstanding finding remained. The existing staged audit was verified unchanged
byte-for-byte; unrelated website edits were preserved. This follow-up introduces
no additional FQL semantic change.

## Calendar boundary follow-up (2026-09-22)

The baseline for this correction is `dd5cfa07344e7a76482faa4cd9c326cfb46d9e9c`.
The preceding sections record the earlier audit and helper-refactor results;
their measurements are historical, not measurements of this correction.

At the supported upper epoch, `math.MaxInt64 - 62_135_596_800`, a UTC+01 wall
value is 3600 seconds greater than the actual instant. Validating that wall
value as a final DateTime prematurely rejected both source reconstruction and
destination calculation. Twelve positive-offset regression cases failed before
the production change, covering zero Add/Subtract shifts across calendar units,
adding a day into the boundary, and subtracting a day from it.

Calendar validation now keeps exact wall seconds separate from representable
`time.Time` intermediates. Bounded reference dates still handle extreme
Gregorian normalization, and source reconstruction applies the known source
offset before validating its epoch. Native field/delta checks, week scaling,
negation checks, and the ordinary-date fast path remain in place.

Go `AddDate` selects the offset and produces the candidate. When a wall value
or first offset guess exceeds the internal epoch range, validation reconstructs
the selected native offset from exact wall seconds and the candidate's Unix
seconds. A wrapped final subtraction would differ by 2^64 and cannot pass this
checked reconstruction. Only raw Unix seconds and nanoseconds are inspected
until the actual instant passes the runtime epoch converter. No out-of-range
wall or provisional `time.Time` is constructed for calendar extraction or
ordering. Zone-transition endpoints are compared by Unix seconds, since an
endpoint's internal epoch can itself wrap.

For an accepted result at the upper boundary, TZif offsets (signed 32-bit) leave
Go's initial offset subtraction within int64; the supported upper endpoint is
more than 62 billion seconds below `MaxInt64`. `FixedZone` can use a larger native
offset, but uses that same offset for both probes, so checked reconstruction
also detects its subtraction overflow. The representable-wall path retains an
explicit int64 check before the first probe.

Synthetic TZif regressions cover DST gaps with an out-of-range wall value or
provisional instant, including a selected offset that differs from the returned
date's `Zone()`. Fixed-zone tests cover UTC, positive and negative offsets, both
epoch endpoints, invalid final dates, and Unix-second wraparound. Existing
month-end, named-zone, and 386 source-year tests remain. Near `MinInt64`, source
calendar extraction can already wrap on 64-bit Go; those sources still fail
reconstruction, including zero shifts. The shared tests assert explicit 386
range failures instead of skipping the group.

Public parameter/FQL regressions run Add and Subtract at None, Basic, and Full
optimization. Successful results preserve exact seconds, nanoseconds, and
location; errors remain `None` plus `runtime.ErrRange` attributed to argument 1.
The runtime epoch limits, fixed-unit arithmetic, and other integer-audit
contracts are unchanged. The intentional FQL change is acceptance of these
representable boundary calendar results.

The website's existing Go embedding migration guide now explicitly replaces
`Unwrap().(int)` with `Unwrap().(int64)` and shows checked `ToNativeInt` use with
failure handling. It explains that small values and 64-bit hosts are affected,
and is linked from host-value guidance. Runtime and compatibility tests now
also exercise the small value 42. Existing website edits were preserved.

### Benchmark evidence

Both revisions used Go 1.26.5 on darwin/arm64, Apple M2 Max. The initial
before/after command was:

```sh
go test ./pkg/stdlib/datetime -run '^$' \
  -bench '^(BenchmarkDateTime|BenchmarkDateTimeSame|BenchmarkIntegerWidthCalendar)$' \
  -benchmem -count=6 -benchtime=200ms
```

All allocation counts and byte counts were unchanged. The first samples showed
calendar increases of roughly 7–8%, alongside 3–11% increases in several
untouched controls. To investigate the timing variation, the baseline commit
was extracted with `git archive` into a temporary directory, and each revision
was compiled once:

```sh
go -C "$TMPDIR/ferret-calendar-boundary-baseline" test -c \
  -o "$TMPDIR/ferret-calendar-boundary-before.test" ./pkg/stdlib/datetime
go test -c -o "$TMPDIR/ferret-calendar-boundary-after.test" ./pkg/stdlib/datetime
```

After native and Docker validation finished, six pairs alternated these
commands, reversing the order in every other pair:

```sh
"$TMPDIR/ferret-calendar-boundary-before.test" -test.run='^$' \
  -test.bench='^BenchmarkIntegerWidthCalendar$' \
  -test.benchmem -test.count=1 -test.benchtime=500ms
"$TMPDIR/ferret-calendar-boundary-after.test" -test.run='^$' \
  -test.bench='^BenchmarkIntegerWidthCalendar$' \
  -test.benchmem -test.count=1 -test.benchtime=500ms
benchstat "$TMPDIR/ferret-calendar-boundary-paired-before.txt" \
  "$TMPDIR/ferret-calendar-boundary-paired-after.txt"
```

The repeated comparison measured:

| IntegerWidthCalendar case | Before | After | Timing delta | B/op, unchanged | allocs/op, unchanged |
| --- | ---: | ---: | ---: | ---: | ---: |
| hour | 121.7 ns | 122.3 ns | no significant change (`p=0.102`) | 104 | 6 |
| day | 118.8 ns | 123.5 ns | +4.00% (`p=0.002`) | 64 | 3 |
| month | 119.6 ns | 123.7 ns | +3.34% (`p=0.002`) | 64 | 3 |
| year | 119.1 ns | 123.4 ns | +3.61% (`p=0.002`) | 64 | 3 |

The correction adds approximately 4 ns/op in these ordinary calendar samples,
with no additional allocation. Exact-arithmetic allocations remain outside the
ordinary-date fast path. Fixed-unit optimization remains separate; its
implementation was not changed to recover benchmark numbers.

### Validation and review

The final tree passed:

```sh
go test ./pkg/runtime ./pkg/stdlib/datetime ./compat/... . -count=1
make fmt
make lint
make test
make compile-32bit
```

Actual Linux/386 execution of `make test-32bit` passed across the root,
`tools/apiref`, and `tools/apipublish` modules, using the Go 1.25.14 Docker command
documented above. The installed Docker daemon was started for this run. This
was executed validation, separately from the Go 1.26.5 cross-compilation check.
Other Go versions remain covered by the existing CI matrix.

The website's `mage build`, the rendered migration anchor and link, affected
maintainer links, and `git diff --check` in both repositories passed. Native
checks used the task-scoped writable Go cache; lint was rerun with module-cache
access after the sandbox denied a metadata write.

The mandatory full-diff self-review covered arithmetic and `time.Time` safety,
native narrowing, DST normalization, error attribution, public/FQL behavior,
tests, performance, documentation, and scope. The transition regression exposed
premature validation of the first offset guess; this was corrected along with
unsafe ordering of extreme zone endpoints. Full validation also caught the
formatter reordering a test struct with positional literals; those cases now use
keyed fields. All affected checks were rerun successfully. No outstanding
finding remains in this follow-up's scope, and unrelated work was preserved.
