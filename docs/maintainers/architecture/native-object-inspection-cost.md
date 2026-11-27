# Native object storage and bounded debugger inspection

Native objects without a large capacity hint keep their first field inline.
They switch to a Go map on the ninth distinct field. A Go value copy shares the current backing. Clearing one receiver
detaches it, even when another copy or iterator survives. The debugger can
describe every representation in constant time.

## Who owns the properties

| Initial capacity hint | Property storage | Inspection history |
| --- | --- | --- |
| Up to 8, including no hint | First field inline on a detached owner. The second distinct field creates eight slots. The ninth promotes the slots to a Go map on the same owner. | On the owner from construction |
| 9–128 | A pre-sized Go map on the handle | One detached owner shared by value copies |
| Above 128 | A pre-sized Go map on the handle | Permanently summary-only for this backing |

Compact storage uses exact string equality. A 64-bit filter can rule out an
absent name before scanning the eight slots. A collision or a stale bit after
deletion only causes an exact scan. Replacement releases the old key and
value. Removal clears the vacated slot. Promotion copies the live entries and
clears all eight slots. Its map reservation is the constant 16, which removes
the allocation-size calculation reported in CodeQL alert 87.

A Go value copy or iterator can keep the old properties alive without a
reference back to a cleared receiver. An iterator snapshots names and reads
values from the captured backing. Callback traversal of compact slots keeps
up to eight original names, checks for removal or promotion, and delivers
current values. Map-backed traversal uses Go's map iteration contract. FQL
object iteration order is unspecified.

`Copy` creates independent storage and shares child identities. `Clone`
clones or copies children under the runtime value contract and propagates
errors. Both size their new backing from the live entries and start fresh
inspection history. A present raw nil from a host child survives native and
dictionary copies and clones. A shape clone converts such a child to `None`
because nil marks vacant slots.

## Inspection eligibility

The storage cutoff is eight. The debugger's separate limit is 128 keys.
The original reservation and the total of live keys plus successful removals
must fit in that limit. A missing-key deletion consumes no budget. Inserting
a key longer than 64 KiB disqualifies that backing, including after deletion.
A formerly large or churned map can therefore have a precise live count and
remain summary-only.

Explicit expansion validates the interval, sorts names, and reads values only
inside the requested page. A valid empty page skips name collection. The VM
and debugger session own page publication and stop-scoped references.

## Work before the first page

Every fixture below contains exactly `{"keep": true}` when inspection starts.
The first-key benchmark stops the Go range as soon as it yields a name. This
omits sorting, formatting, and handle creation. The known-key benchmark
directly looks up the same name to isolate the cost after discovery.

| Backing history | First key, Go 1.26.9 | First key, Go 1.27.3 | First key, Go 1.27.4 |
| --- | ---: | ---: | ---: |
| Fresh singleton | 29 ns | 27 ns | 29 ns |
| Reserved for 65,536 entries | 26.5 µs | 26.7 µs | 25.7 µs |
| Reserved for 1,000,000 entries | 1.02 ms | 1.04 ms | 2.01 ms |
| Grew to 1,000,000, then deleted 999,999 | 1.10 ms | 1.02 ms | 763 µs |

These are medians of five samples on an Apple M2 Max. Known-key lookup took
8–10 ns across the fixtures. Both operations allocated zero bytes per
iteration. Fixture construction and deletion were outside the timed loops.

```sh
GOWORK=off go test -mod=readonly ./pkg/runtime -run '^$' -bench '^BenchmarkObjectInspectionKeyDiscovery$' \
  -benchmem -cpu=1 -benchtime=150ms -count=5
```

The supported Go map API exposes the live count. Its iterator can cross
empty groups and tables before yielding control to the caller.
A limit on yielded keys or a context check between yields cannot bound that
work. For a lexicographic page of a multi-key map, the first yielded key
also need not be the first property in the requested order.

## Alternatives at a stop

- **Cache sorted names per stop.** Once a bounded scan has succeeded, its names
  can serve arbitrary offsets in stable order until resume or close. Building
  such an index still requires evidence that the first scan is bounded.
  The existing reader already produces the same order on repeated pages.
- **Start tracking an object on first sight or at debug-session creation.**
  A host can supply a map already reserved or thinned before the session.
  Writes and deletes through a Go value copy affect the same backing. A log
  started at the stop or keyed by receiver address misses those histories.
  Attaching state on the first exceptional mutation also needs a way for
  earlier aliases to find it. `Clear` must detach the receiver and leave
  earlier copies independent of its replacement values.
- **Use existing evidence when available.** A dense VM shape already supplies
  ordered names. An opt-in host provider can expose bounded children from
  storage it owns. Explicit member evaluation can look up a supplied name
  in a native object. Listing that object's children still requires discovering
  the names to establish their order. Checking key lengths before sorting
  could bound string comparisons after discovery, but leaves the first
  traversal unresolved.

## Storage alternatives

At `e68ab200`, reservations from 9 through 128 used a fingerprint vector
and a 256-byte table of candidate bounds. A four-byte miss on a fresh
128-field dictionary often still required an exact search. The historical
FQL miss result was 132.3 µs against 123.7 µs on merged #1068 (+6.9%,
p < 0.001). The previous profiles, rejected designs and four-way measurements
are preserved in the [investigation at e68ab200](https://github.com/MontFerret/ferret/blob/e68ab20049a84a99df97cee5f49664e06f0f99a6/docs/maintainers/architecture/native-object-inspection-cost.md).

Three interleaved 70 ms samples per variant screened the cutoff. Times below
are medians on the same Go 1.27.4 arm64 host. These pilots used the existing
tracked-map promotion path for medium reservations. The final candidate uses
the pre-sized map directly.

| Operation | Index through 128 | Index through 32 | Index through 16 | Go maps above 8 |
| --- | ---: | ---: | ---: | ---: |
| 16 fields, four-byte miss | 6.01 ns | 5.95 ns | 6.02 ns | 8.49 ns |
| 16 fields, long-key miss | 11.9 ns | 12.3 ns | 9.79 ns | 8.32 ns |
| 32 fields, long-key hit | 12.8 ns | 13.3 ns | 11.3 ns | 10.3 ns |
| 128 fields, four-byte miss | 17.4 ns | 8.54 ns | 8.53 ns | 9.04 ns |
| 128 fields, long-key miss | 16.3 ns | 9.25 ns | 8.68 ns | 9.12 ns |
| Reserved 16-field build, ns / B / allocations | 281 / 1,008 / 3 | 281 / 912 / 3 | 294 / 912 / 3 | 415 / 1,352 / 6 |
| Reserved 128-field build, ns / B / allocations | 1,608 / 5,360 / 3 | 2,456 / 9,672 / 6 | 2,429 / 9,672 / 6 | 2,380 / 9,672 / 6 |

The rotating fixtures also covered 7–9, 15–17 and 31–33 properties. The index
helped short names and construction at 16. Go maps handled the different name
patterns more consistently. Extending only the existing filtered linear scan
to 16 slots built a reserved object in 260 ns and two allocations, but four-byte
and long-key misses took 35.6 and 32.8 ns. A pre-sized Go map took 7.25 and
7.41 ns in that four-sample comparison. Removing the medium index also removes
the seeded fingerprint, short-key mixing, candidate bounds, collision probing
and deletion-time tag maintenance.

### The medium-map allocation

A profile of 1,000 escaping, reserved 128-field builds attributed 5,000
allocations to the Go map and native handle, and exactly 1,000 to the detached
owner that carries inspection history. The direct-map prototype used 9,624 B
and six allocations per build. Task 2 used 9,568 B and five. Storing the map
behind the compact group used 9,672 B and six. Putting mutable history on the
receiver would let a surviving alias retain the receiver's replacement children
after Clear.

The allocation exception covers one owner per independently created,
inspectable medium Go map, including dictionaries created or copied by
FastObject. The final read qualification meets the agreed condition for this exception. Singletons keep
their history in existing compact storage. Oversized reservations have no
history allocation.

## VM dictionaries

Shape IDs belong to the VM that allocated them. Its read instructions cache slots
for native string keys on its own shapes. Objects passed in from another VM use
named lookup. FastObject resolves a host key once before choosing its current
storage, since a host callback can promote the object. Hashing walks sorted
original names and contributes the hash of `None` for a sibling removed by an
earlier child's hash callback, as native objects do.

A named overwrite resolves its slot once and uses it to update occupancy and value. When a
distinct write crosses the configured threshold, the VM embeds a native
dictionary in the existing fast-object header and reserves the live count
including the incoming field. The default transition reserves 65 fields.
The previous 128-field reservation selected a larger Go map immediately. In
five alternating pilot samples the smaller hint reduced transition bytes
from 9,608 to 5,000 and the 16-literal FQL workload from about 720 KB to
646 KB. A dictionary that subsequently grows to 128 can pay for map growth.

In dictionary mode the VM already has a concrete FastObject at its property
cache entry point. Calling its getter directly avoids a redundant interface
dispatch on every read. Host key conversion still occurs before any
representation-sensitive lookup. Clearing a dictionary keeps dictionary mode
and detaches the receiver. Copy and Clone use the native ownership contract.

## Correctness and ownership coverage

- Hinted and grown 128-field dictionaries return the expected final page
  and accept an empty page without allocating.
  Deletion through an alias updates the shared budget, and replacing the
  deleted field exhausts it. A sparse original rejects expansion before a
  visitor runs. Copy and Clone of its single survivor regain eligibility.
- Weak-reference tests retain aliases and iterators across Clear and verify
  collection of overwritten values, deleted values, old equal-key substrings,
  and replacement children. Copies taken before and after promotion keep their
  independent keys. Clones own independent children.
- Exact-name and mixed-mutation tests cover filter collisions, moved slots,
  empty names and arbitrary string bytes. VM and session suites cover indexed
  and named views, selected child reads, atomic page failure and reference
  invalidation. The million-entry sparse fixtures exercise both initial
  reservation and deletion history.

## Reproducing the comparisons

Use Task 2 (`3ac60f34`), merged #1068 (`ac2db0dc`), current #1069 before
qualification (`e68ab200`), and the qualified implementation (`c51cc58f`).
Keep each revision's production code and modules. Use the final native
operations and rotating-key fixtures, FastObject fixture, FQL literals, and
JSON and MessagePack fixtures in all four snapshots. On Task 2 copy the
escaping `BenchmarkNewObjectOf` sink and supply a test-only `SetString` that
normalizes nil and writes to `t.data[key]`. On Task 2 and merged #1068 adapt
only the VM fixture constructors to `NewShapeCache(threshold)` and
`NewFastObject(cache, threshold)`.

All four binaries used Go 1.27.4, darwin/arm64, Apple M2 Max, `GOWORK=off`,
and `-mod=readonly`. Verify identical flat and nested wire bytes. The FQL
fixture verifies all 1,024 returned values before timing. Rotating keys use
distinct equal-content strings and wrap at the actual slice length.

Build each package with `go test -mod=readonly -p=2 -c`. Run one binary and
one benchmark row at a time with these flags, replacing the filter with the
anchored row name. Keep the samples for each revision in a separate file.

```sh
./runtime.test -test.run='^$' \
  -test.bench='^BenchmarkObjectLookupKeys$/^128$/^four_bytes$/^miss$' \
  -test.benchmem -test.cpu=1 -test.benchtime=200ms -test.count=1
./benchmarks.test -test.run='^$' \
  -test.bench='^BenchmarkObjectDictionaryReads_None$/^miss$' \
  -test.benchmem -test.cpu=1 -test.benchtime=1s -test.count=1
benchstat task2.txt merged.txt pr.txt final.txt
benchstat merged.txt final.txt
benchstat confirm-merged.txt confirm-final.txt
```

The four-way matrix uses ten 200 ms samples per row. Each row runs on all four
revisions consecutively in balanced order. A shuffle seeded with 1069 plus
the zero-based round number varies row order. The longer confirmation uses
twenty alternating one-second samples of merged #1068 and the final candidate
for the two FQL reads and the five disputed native reads. Compare the two
series separately. A further twenty one-second samples cover shape construction
and named overwrites. Benchstat's pairwise p-values use its default 0.05 level.
The final VM slot edit changed the VM and FQL binaries. Both suites were
remeasured against their baselines. The rebuilt native, JSON, and MessagePack
binaries were byte-identical to the ones already measured.

## Measured performance

Each cell in the four-way table gives median time and allocations per operation
from ten samples. The downloadable matrix includes bytes per operation, all 114
rows, and pairwise benchstat results against both merged #1068 and the previous
PR implementation. The archive also includes the individual samples, run order,
fixture checksums and cutoff pilots. [Benchmark evidence](native-object-qualification.tar.gz)

| Operation, time / allocations | Task 2 | Merged #1068 | Previous #1069 | Final |
| --- | ---: | ---: | ---: | ---: |
| FQL dictionary miss | 129.2 µs / 101 | 128.9 µs / 102 | 134.8 µs / 96 | 127.3 µs / 101 |
| FQL dictionary hit | 131.9 µs / 101 | 131.9 µs / 102 | 135.8 µs / 96 | 130.8 µs / 101 |
| Four-byte rotating miss, 128 fields | 9.08 ns / 0 | 9.42 ns / 0 | 17.8 ns / 0 | 8.37 ns / 0 |
| Build singleton | 95.9 ns / 3 | 101.5 ns / 4 | 31.1 ns / 2 | 32.7 ns / 2 |
| Overwrite singleton | 7.99 ns / 0 | 7.94 ns / 0 | 3.11 ns / 0 | 3.81 ns / 0 |
| Build 128 reserved fields | 1.958 µs / 5 | 1.988 µs / 6 | 1.710 µs / 3 | 2.303 µs / 6 |
| Copy 128 fields | 10.746 µs / 14 | 10.681 µs / 15 | 1.908 µs / 3 | 4.627 µs / 6 |
| Clone 128 fields | 11.564 µs / 14 | 11.698 µs / 15 | 2.557 µs / 3 | 4.987 µs / 6 |
| Build 256 reserved fields | 3.772 µs / 5 | 3.817 µs / 6 | 3.752 µs / 5 | 4.513 µs / 5 |
| Decode 5,000 JSON objects | 1.624 ms / 20,075 | 1.695 ms / 25,075 | 1.346 ms / 15,075 | 1.350 ms / 15,075 |
| Decode 5,000 MessagePack objects | 1.208 ms / 15,003 | 1.235 ms / 20,003 | 584.2 µs / 10,003 | 602.2 µs / 10,003 |
| Named shape overwrite | 8.85 ns / 0 | 8.65 ns / 0 | 8.35 ns / 0 | 8.35 ns / 0 |
| Shape to dictionary | 1.325 µs / 4 | 1.311 µs / 5 | 1.111 µs / 2 | 1.401 µs / 5 |
| Copy dictionary | 2.425 µs / 5 | 2.462 µs / 6 | 1.010 µs / 3 | 2.610 µs / 6 |
| Clone dictionary | 2.586 µs / 5 | 2.583 µs / 6 | 1.353 µs / 3 | 2.822 µs / 6 |
| Clear and reuse dictionary | 96.8 ns / 2 | 107.4 ns / 3 | 23.9 ns / 1 | 23.5 ns / 1 |
| FQL sixteen 65-field literals | 407.4 µs / 1,123 | 408.0 µs / 1,139 | 400.2 µs / 1,075 | 401.8 µs / 1,123 |

### Read and shape confirmation

Each FQL operation builds a dictionary and executes 1,024 reads. The fixture checks the results before timing.
The native lookups use a prebuilt 128-field object and allocate zero bytes.
The independent one-second comparison used twenty samples per revision.

| Confirmed workload | Merged #1068 | Final | Change | p |
| --- | ---: | ---: | ---: | ---: |
| FQL hit | 131.25 µs | 130.53 µs | -0.5% | 0.024 |
| FQL miss | 128.10 µs | 126.77 µs | -1.0% | < 0.001 |
| Four-byte hit | 9.99 ns | 9.47 ns | -5.2% | < 0.001 |
| Four-byte miss | 9.28 ns | 8.26 ns | -11.0% | < 0.001 |
| Eight-byte miss | 8.98 ns | 8.31 ns | -7.4% | 0.001 |
| Long common-suffix hit | 10.48 ns | 9.93 ns | -5.2% | < 0.001 |
| Long common-suffix miss | 9.04 ns | 8.46 ns | -6.4% | 0.046 |
| Build four-field shape | 321.30 ns | 319.00 ns | -0.7% | < 0.001 |
| Named shape overwrite | 8.60 ns | 8.34 ns | -3.0% | < 0.001 |

The median **paired** change in FQL miss time was -0.98%, with a
distribution-free 95.9% interval of -1.95% to -0.34%
(the sixth and fifteenth ordered changes). The named shape overwrite exercises
`SetStringCached`. VM bytecode with a warm write cache uses `SetSlotWithShape`.
On the first map-backed candidate, the named shape overwrite took 9.38 ns
versus 8.59 ns on #1068 (+9.2%, p < 0.001, twenty samples). Reusing the
resolved slot for both the occupancy check and assignment removed a second
bounds check on arm64. The final named overwrite took 8.34 ns, 3.0% faster than merged #1068.

### What the storage choice costs

- At seven and eight fields, names of equal length with the same final byte
  can force a full eight-slot scan. At eight fields the rotating long-name miss
  took 21.4 ns, against 10.8 ns on merged #1068 and 20.1 ns on the previous
  PR. Building eight reserved fields took 132 ns, against 166 ns on #1068.
- Building a reserved 128-field object took 2.30 µs and 9,624 B, against
  1.99 µs and 9,578 B on #1068. At 256 fields the difference was 4.51 versus
  3.82 µs (+18.3%, p < 0.001), even though both versions use Go maps. Three
  rearrangements of the setter gave no significant improvement in that build
  in five interleaved pilot samples. Unhinted growth to either 128 or 256 was
  statistically inconclusive against #1068.
- Fixed-key overwrite at 128 fields cost 9.67 ns, against 8.75 ns on #1068
  and 6.30 ns with the fingerprint index. Copy and Clone also cost more than
  the indexed representation, but took 57% less time than #1068. Singleton
  overwrite rose from 3.11 ns on the previous PR to 3.81 ns, still 52%
  faster than #1068. The complete matrix lists smaller regressions and
  inconclusive differences individually.
- FastObject dictionary transition took 6.9% longer than #1068, copying
  took 6.0% longer, and cloning took 9.2% longer (all p < 0.001). Transition
  used five allocations, and both Copy and Clone used six. Each operation
  spends the one additional medium-map owner accepted against Task 2.
  The separate FastObject dictionary wrapper is gone. Dictionary clear/reuse
  took 23.5 ns and one allocation, against 107 ns and three on #1068.
- End-to-end construction of sixteen 65-field FQL literals used 1,123
  allocations and 645,607 B, compared with 1,139 and 663,527 B on #1068.
  The shorter matrix found a 1.5% time improvement for this workload.
  Direct FQL property mutation was 1.4% slower than the previous PR but
  5.6% faster than #1068. None of the measured FQL workloads showed a
  confirmed slowdown against #1068.

The eight-field cutoff restores medium read performance across key patterns,
preserves the singleton and nested decode gains, and eliminates the medium
index. The remaining costs fall on specific compact-key patterns and
write-heavy construction and copying. The measurements support taking the
simpler representation with the scoped allocation exception.

## Qualification after the VM read fixes

The rerun compares commit `76aecf6a` with merged #1068 for dictionary reads and
the pre-review head `8a35deb0` for same-VM shape reads. Each query reads 1,024
properties. Each pair used twenty alternating one-second samples on the same
Apple M2 Max, Go 1.27.4, darwin/arm64. The dictionary fixture checks all
returned values before timing. The new shape fixture uses a parameterized
object literal and exercises the constant and dynamic instructions separately.

| FQL read | Baseline | Baseline µs/op | Candidate µs/op | Allocations/op | Time vs. baseline |
| --- | --- | ---: | ---: | ---: | --- |
| Dictionary hit | #1068 | 147.96 | 144.79 | 102 → 101 | inconclusive (p = 0.221) |
| Dictionary miss | #1068 | 146.09 | 138.92 | 102 → 101 | inconclusive (p = 0.149) |
| Dictionary repeat hit | #1068 | 137.45 | 135.10 | 102 → 101 | -1.7% (p = 0.002) |
| Dictionary repeat miss | #1068 | 134.35 | 131.39 | 102 → 101 | -2.2% (p < 0.001) |
| Shape constant | 8a35deb0 | 99.07 | 97.86 | 1,552 → 1,552 | inconclusive (p = 0.355) |
| Shape dynamic | 8a35deb0 | 110.60 | 109.04 | 1,552 → 1,552 | -1.4% (p = 0.014) |

Timing varied during the first series, so the dictionary comparison was repeated after the other workloads finished. The median paired change for misses and its distribution-free 95.9% interval were -4.35% (-8.22% to +2.31%) in the first series and -2.14% (-2.67% to -1.84%) in the repeat. Each interval uses the sixth and fifteenth ordered paired changes. The repeat supports the accepted one-allocation allowance for independently owned, inspectable medium maps. Both same-VM read fixtures retained their original allocations.

### Mutation on a tracked backing

The earlier `insert_remove` rows reserved one extra field. At 128 fields that
omitted the inspection-history owner. The corrected fixture reserves the
starting size and checks that the candidate's 128-field object is initially
pageable. Inserting the extra key crosses the 128-key limit. Removing it
consumes deletion history, so this measures repeated churn on the same
history-tracked backing after page eligibility is lost.

Each cell is the median in ns/op from ten balanced, interleaved 200 ms samples.
All four sizes used zero bytes and zero allocations per timed operation on
all four revisions.

| Starting fields | Task 2 | Merged #1068 | Previous #1069 | Candidate | Candidate vs. #1068 |
| ---: | ---: | ---: | ---: | ---: | --- |
| 1 | 21.41 | 23.78 | 12.39 | 10.73 | -54.9% (p < 0.001) |
| 16 | 24.84 | 27.62 | 29.84 | 29.05 | +5.2% (p < 0.001) |
| 128 | 24.91 | 27.47 | 29.95 | 29.52 | +7.4% (p < 0.001) |
| 256 | 25.02 | 27.50 | 27.69 | 28.82 | +4.8% (p = 0.004) |

These four mutation rows supersede the corresponding rows in the original
matrix. The other storage and codec measurements retain their original commit
labels. The `review/` directory in the archive includes the raw samples,
run order, binary checksums, commands and benchstat output.
