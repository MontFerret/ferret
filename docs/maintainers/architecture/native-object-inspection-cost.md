# Native object storage and bounded debugger inspection

The debugger can describe a native object in constant time and can expand a
small, eligible object in sorted key order. The native constructor and both
object decoders meet the allocation counts measured before debugger paging.
There is no separately allocated inspection record.

## Who owns the properties

PR #1068 stored a Go map and a pointer to `debugmap.Bounds` on each
`runtime.Object`. Copies of the Go struct shared both pointers. The extra
allocation made changes to the deletion count and long-key flag visible through
every alias. Putting the bounds on the original receiver and pointing aliases
back at it restored the allocation counts, but an old alias then kept that
receiver and all its replacement values alive after `Clear`.

An object is a two-word handle. Its capacity hint selects the storage.

| Initial hint | Property storage | Inspection history |
| --- | --- | --- |
| Up to 8, including no hint | First field inline in the detached owner, eight slots on the second distinct key | Bounds live on the owner |
| 9–128 | Detached owner, name index and 16, 32, 64 or 128 slots in one allocation | Bounds live on the owner |
| Above 128 | Direct Go map | Permanently summary-only for this backing |

In the eight-slot group, a 64-bit filter derived from the key length and final
byte can rule out a key before a linear scan. Larger slot groups index names
by a one-byte fingerprint. A unique fingerprint identifies one candidate. A
shared fingerprint requires scanning at most 128 bytes for candidates. All
matches require exact string equality.

When the slots fill, a new key moves the properties to a Go map in the same
owner. All existing handles see the move. Removing a field or moving a group
clears the old slots. Slots co-allocated with their header still occupy their
original bytes after promotion, but the cleared slots release their child values.

A Go value copy and an iterator can keep the old storage alive. Neither has
a pointer back to the receiver. Clearing either receiver gives it a fresh
empty owner. `Copy` creates independent storage for the live keys and preserves
child identities. `Clone` clones or copies each child under the existing value
contract. Both start with fresh inspection history. Iterators snapshot keys
and read values from the captured backing. Callback traversal also reads from
the original backing if user code clears the receiver or moves the properties
during the callback.

## Inspection eligibility

The existing 128-key and 64 KiB key-length limits apply to native slots,
promoted maps and VM dictionaries. The live count plus successful removals
must fit in 128. A missing-key deletion consumes no budget. Inserting a long
name disqualifies the backing even if that name is later deleted. A direct,
oversized reservation can reject inspection without a scan. The exact live
count remains available for every representation.

Explicit expansion collects and sorts **names** from eligible storage, then
looks up values only in the requested interval. The VM and session still own
page validation, atomic publication and stop-scoped references. The runtime
keeps no sorted-name cache or debugger session state. FQL object iteration
order remains unspecified.

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

- **Summary-only native objects.** A test overlay on #1068 removed the bounds
  and recovered every allocation baseline. The fresh singleton lost its
  expansion handle and the named-page test returned `errors.ErrUnsupported`.
- **Detached singleton with a Go map on the second key.** This passed the
  alias and GC checks, but growing to two keys needed four allocations and flat
  JSON still took 1,037. Keeping the first property in the detached owner and
  using an eight-slot group on the second key pays for storage as the object
  grows. An alternative that allocated all eight slots on the first write used
  320 bytes for a singleton and was slower on small object construction.
- **A Go map for all pre-sized objects.** Separate map and history owners kept
  the extra allocation for ordinary 16-, 64- and 128-field reservations.
  Using a flat slice of sorted strings instead made a 128-field build take
  about 21 µs in a prototype. The bounded fingerprint index avoids shifting
  properties or linearly comparing the full names on every insertion.
- **External hash table.** A prototype met the headline allocation counts, but
  building 16 pre-sized fields used 1,552 B and 128 fields used 11,024 B. Go's
  built-in map already handles unbounded growth. The native representation
  requires no additional dependency.

## Ownership and debugger tests

`TestObjectClearDetachesSurvivingCopy` keeps the old backing alive, forces a
GC and checks a weak pointer to a replacement child. It exercises both
directions of detachment with small, indexed and directly reserved storage.
`TestObjectBackingReleasesRemovedValues` checks that removing a child releases
it while the owner is still live. The alias and callback tests cover mutation,
promotion during iteration, multiple clear/reuse cycles and the empty key.
The `ForEach` adapter can inline into callers. A concrete callback that does
not retain its keys can traverse the object without allocating a boxed string
per property. A collision test exercises overwrite, removal, reinsertion and
sorted paging of distinct strings with the same fingerprint. A deterministic
mixed-write test compares every lookup and length against a Go map across
all four indexed size classes and the small and direct-map cases.

`TestObjectBackingSparseInspection` pre-reserves a million entries or inserts
a million and deletes all but one. The sparse original rejects its named page
before invoking the visitor. Its independent copy and clone both allow
inspection. VM and debugger session tests cover exact names, selected ranges,
atomic failure and invalidation of references when a stop ends.

## Reproducing the comparisons

The initial comparison between Task 2 and merged #1068 was recorded before
changing storage. The final comparison adds callback traversal and isolates the
insert/remove fixture so it cannot change a later benchmark's backing.

The measurements below use Go 1.27.4 on an Apple M2 Max (darwin/arm64). Each
revision used a clean source snapshot, identical keys and values, one CPU
and ten 200 ms samples. The two historical revisions ran in rotation with
the storage change. After specializing native string reads, the final code
ran the full matrix again. An additional paired comparison alternated the
merged and final implementations for 16- and 256-field lookups and 128-field
reserved construction. To reproduce, run both commands five times per
revision and compare the output files with `benchstat`.

```sh
export GOWORK=off
go test -mod=readonly -p=1 ./pkg/runtime ./pkg/vm/internal/data ./test/benchmarks -run '^$' \
  -bench '^Benchmark(NewObjectOf|ObjectOperations|FastObject|ObjectLiterals_None|DirectMutation_ObjectProperty_None|ObjectDictionaryLiterals_None|Match_ObjectPattern_None)$' \
  -benchmem -cpu=1 -benchtime=200ms -count=2
go test -mod=readonly -p=1 ./pkg/encoding/json ./pkg/encoding/msgpack -run '^$' \
  -bench '^Benchmark(JSONCodecDecode|MsgpackCodecDecode)$/^(flat_object_256|nested_object_5000)$' \
  -benchmem -cpu=1 -benchtime=200ms -count=2
benchstat task2.txt merged.txt candidate.txt
benchstat merged.txt candidate.txt
```

Copy `pkg/runtime/object_operations_bench_test.go` into both historical
snapshots. Task 2 also needs the `BenchmarkNewObjectOf` fixture from
`pkg/runtime/object_test.go`, `pkg/vm/internal/data/fast_object_bench_test.go`
and `test/benchmarks/bench_literals_test.go`. Its test-only `SetString`
adapter converts a nil value to `None` and assigns directly to `t.data[key]`.
This gives the historical native-write microbenchmark an equivalent Go-string
entry point. The codec and end-to-end fixtures execute each revision's own
production code.

## Measured performance

Empty construction and both decoder fixtures meet or improve on the Task 2
allocation counts. Empty objects take 64 B, compared with 56 B at Task 2
and 66 B in #1068. A one-field object also takes 64 B. Objects with two to
eight fields take 384 B, up from 354 B in #1068. At 5,000 nested objects,
JSON allocated 50% fewer bytes and took 20% less time than #1068. MessagePack
allocated 82% fewer bytes and took 52% less time.

Ordinary latency remains a tradeoff. In the alternating confirmation run,
16-field lookup went from 8.00 to 9.52 ns (+19%), and 256-field lookup from
7.65 to 8.24 ns (+7.7%). Pre-sized construction of 128 fields went from
2.05 to 2.31 µs (+13%), while using 4,880 B and two allocations in place
of 9,578 B and six. All three timing differences had p < 0.001 across ten
samples. Flat JSON also took 3.2% longer in the full matrix. Singleton reads,
writes, copies and clear/reuse were faster.

The four FQL workloads showed no significant timing regression against #1068.
FastObject retained the same allocation counts as #1068 across the matrix.
Dictionary transition, copy, clone and clear/reuse still allocate a 16-byte
wrapper relative to Task 2. The 65-field FQL dictionary fixture retains the
corresponding 16 extra allocations per execution.

Each cell gives median **ns/op / B/op / allocs/op**. Task 2 is
`3ac60f34` and merged is `ac2db0dc`. Construction benchmarks publish the
object to an escaping sink. Grow starts without a size hint and reserved
uses the final key count. Operation fixtures contain the indicated number
of properties. Insert/remove reserves one extra slot before timing.

### Decoding and FQL execution

| Benchmark | Task 2 | Merged #1068 | This change |
| --- | ---: | ---: | ---: |
| `DirectMutation_ObjectProperty` | 341 / 128 / 4 | 344 / 128 / 4 | 341 / 128 / 4 |
| `ObjectLiterals` | 412 / 288 / 6 | 409 / 288 / 6 | 410 / 288 / 6 |
| `ObjectDictionaryLiterals` | 369,042 / 663,244 / 1,123 | 367,945 / 663,505 / 1,139 | 369,196 / 663,504 / 1,139 |
| `Match_ObjectPattern` | 242 / 32 / 2 | 237 / 32 / 2 | 233 / 32 / 2 |
| `JSONCodecDecode/flat_object_256` | 55,585 / 51,640 / 1,036 | 55,730 / 51,648 / 1,037 | 57,514 / 51,120 / 1,035 |
| `JSONCodecDecode/nested_object_5000` | 1,497,094 / 2,828,288 / 20,075 | 1,515,921 / 2,878,297 / 25,075 | 1,217,642 / 1,428,290 / 15,075 |
| `MsgpackCodecDecode/flat_object_256` | 23,066 / 19,400 / 254 | 22,940 / 19,392 / 255 | 23,791 / 19,393 / 254 |
| `MsgpackCodecDecode/nested_object_5000` | 1,102,372 / 1,720,168 / 15,003 | 1,150,978 / 1,770,168 / 20,003 | 553,346 / 320,152 / 10,003 |

### Native construction and growth

| Benchmark | Task 2 | Merged #1068 | This change |
| --- | ---: | ---: | ---: |
| `NewObjectOf/0` | 31.9 / 56 / 2 | 36.4 / 66 / 3 | 24.3 / 64 / 2 |
| `NewObjectOf/1` | 29.4 / 56 / 2 | 36.3 / 66 / 3 | 24.4 / 64 / 2 |
| `NewObjectOf/256` | 974 / 18,528 / 5 | 1,059 / 18,538 / 6 | 993 / 18,536 / 5 |
| `grow_1` | 87.9 / 344 / 3 | 99.5 / 354 / 4 | 32.2 / 64 / 2 |
| `reserved_1` | 87.1 / 344 / 3 | 94.6 / 354 / 4 | 30.4 / 64 / 2 |
| `grow_2` | 96.3 / 344 / 3 | 102 / 354 / 4 | 90.0 / 384 / 3 |
| `reserved_2` | 95.1 / 344 / 3 | 102 / 354 / 4 | 89.0 / 384 / 3 |
| `grow_8` | 153 / 344 / 3 | 160 / 354 / 4 | 132 / 384 / 3 |
| `reserved_8` | 153 / 344 / 3 | 160 / 354 / 4 | 132 / 384 / 3 |
| `grow_16` | 926 / 2,144 / 8 | 936 / 2,154 / 9 | 549 / 1,624 / 7 |
| `reserved_16` | 337 / 1,248 / 5 | 349 / 1,258 / 6 | 306 / 1,040 / 2 |
| `grow_64` | 3,876 / 9,376 / 12 | 3,968 / 9,386 / 13 | 3,644 / 8,856 / 11 |
| `reserved_64` | 1,022 / 4,960 / 5 | 1,032 / 4,970 / 6 | 1,049 / 2,704 / 2 |
| `grow_128` | 7,743 / 18,880 / 14 | 7,690 / 18,890 / 15 | 7,771 / 18,360 / 13 |
| `reserved_128` | 1,938 / 9,568 / 5 | 1,959 / 9,578 / 6 | 2,108 / 4,880 / 2 |
| `grow_256` | 15,596 / 37,344 / 16 | 15,785 / 37,354 / 17 | 15,885 / 36,824 / 15 |
| `reserved_256` | 3,724 / 18,528 / 5 | 3,776 / 18,538 / 6 | 3,728 / 18,536 / 5 |

### Ordinary native operations

| Benchmark | Task 2 | Merged #1068 | This change |
| --- | ---: | ---: | ---: |
| `1/overwrite` | 7.91 / 0 / 0 | 7.87 / 0 / 0 | 3.17 / 0 / 0 |
| `1/set` | 7.92 / 0 / 0 | 9.60 / 0 / 0 | 4.64 / 0 / 0 |
| `1/lookup` | 7.36 / 0 / 0 | 7.37 / 0 / 0 | 3.53 / 0 / 0 |
| `1/insert_remove` | 20.6 / 0 / 0 | 22.7 / 0 / 0 | 11.8 / 0 / 0 |
| `1/iterate` | 88.8 / 80 / 3 | 89.1 / 80 / 3 | 48.2 / 80 / 3 |
| `1/for_each` | 29.9 / 0 / 0 | 29.9 / 0 / 0 | 4.72 / 0 / 0 |
| `1/copy` | 129 / 344 / 3 | 136 / 354 / 4 | 33.8 / 64 / 2 |
| `1/clone` | 141 / 344 / 3 | 150 / 354 / 4 | 36.8 / 64 / 2 |
| `16/overwrite` | 8.36 / 0 / 0 | 8.33 / 0 / 0 | 8.53 / 0 / 0 |
| `16/set` | 8.52 / 0 / 0 | 10.0 / 0 / 0 | 10.3 / 0 / 0 |
| `16/lookup` | 7.67 / 0 / 0 | 7.65 / 0 / 0 | 8.82 / 0 / 0 |
| `16/insert_remove` | 24.2 / 0 / 0 | 26.2 / 0 / 0 | 18.6 / 0 / 0 |
| `16/iterate` | 736 / 560 / 18 | 730 / 560 / 18 | 506 / 560 / 18 |
| `16/for_each` | 126 / 0 / 0 | 128 / 0 / 0 | 102 / 0 / 0 |
| `16/copy` | 1,284 / 2,144 / 8 | 1,310 / 2,154 / 9 | 337 / 1,040 / 2 |
| `16/clone` | 1,438 / 2,144 / 8 | 1,456 / 2,154 / 9 | 438 / 1,040 / 2 |
| `256/overwrite` | 8.50 / 0 / 0 | 8.56 / 0 / 0 | 8.85 / 0 / 0 |
| `256/set` | 8.59 / 0 / 0 | 10.2 / 0 / 0 | 10.1 / 0 / 0 |
| `256/lookup` | 7.79 / 0 / 0 | 7.66 / 0 / 0 | 8.24 / 0 / 0 |
| `256/insert_remove` | 23.6 / 0 / 0 | 26.2 / 0 / 0 | 26.2 / 0 / 0 |
| `256/iterate` | 10,760 / 9,008 / 258 | 10,507 / 9,008 / 258 | 10,754 / 9,008 / 258 |
| `256/for_each` | 1,650 / 0 / 0 | 1,678 / 0 / 0 | 1,672 / 0 / 0 |
| `256/copy` | 21,154 / 37,344 / 16 | 21,758 / 37,354 / 17 | 8,686 / 18,536 / 5 |
| `256/clone` | 23,432 / 37,344 / 16 | 24,284 / 37,354 / 17 | 9,299 / 18,536 / 5 |
| `clear_reuse` | 75.5 / 336 / 2 | 84.1 / 338 / 3 | 20.6 / 48 / 1 |

### VM FastObject

| Benchmark | Task 2 | Merged #1068 | This change |
| --- | ---: | ---: | ---: |
| `shape_build_4` | 317 / 224 / 5 | 305 / 224 / 5 | 320 / 224 / 5 |
| `shape_overwrite` | 8.54 / 0 / 0 | 8.28 / 0 / 0 | 8.31 / 0 / 0 |
| `transition` | 1,256 / 4,952 / 4 | 1,272 / 4,968 / 5 | 1,264 / 4,968 / 5 |
| `dictionary_build_65` | 15,615 / 40,280 / 69 | 15,670 / 40,296 / 70 | 15,670 / 40,296 / 70 |
| `dictionary_overwrite` | 15.3 / 0 / 0 | 9.62 / 0 / 0 | 9.72 / 0 / 0 |
| `dictionary_insert_remove` | 36.1 / 0 / 0 | 31.3 / 0 / 0 | 35.6 / 0 / 0 |
| `dictionary_copy` | 2,299 / 5,016 / 5 | 2,310 / 5,032 / 6 | 2,342 / 5,032 / 6 |
| `dictionary_clone` | 2,426 / 5,016 / 5 | 2,435 / 5,032 / 6 | 2,442 / 5,032 / 6 |
| `dictionary_clear_reuse` | 84.5 / 336 / 2 | 93.8 / 352 / 3 | 93.6 / 352 / 3 |
