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
| 9–128 | Detached owner and name index, plus a slice with the requested capacity | Bounds live on the owner |
| Above 128 | Direct Go map | Permanently summary-only for this backing |

For two to eight properties, a 64-bit filter derived from the key length and
last byte can rule out a key before a linear scan. For a medium reservation,
each property has a one-byte fingerprint. The index records the last possible
slot for each fingerprint. A read checks that slot first, then searches only
earlier fingerprint bytes for other candidates. Every match requires exact
string equality. Appending a new name records its slot. Removing a name moves
the last property left, so the previous upper bounds remain valid. After
repeated removals an absent tag can retain a stale upper bound. Such a miss
may scan the remaining fingerprint bytes, bounded by 128 live slots. The
rotating lookup fixtures use freshly built objects.

Filling the reserved slots and then inserting another distinct key moves the
properties into a Go map on the same owner. Existing handles see the move.
Overwriting a slot replaces both the key and the value. An equal short name
can therefore release the large string it originally sliced. Removal clears the
vacated slot. Promotion clears child references and drops
the slot slice, allowing the separate medium property array to be collected.
The owner retains its inlined name-index bytes.

A Go value copy and an iterator can keep the old storage alive. Neither has
a pointer back to the receiver. Clearing either receiver gives it a fresh
empty owner. `Copy` creates independent storage for the live keys and preserves
child identities. `Clone` clones or copies each child under the existing value
contract. Native objects and embedded dictionaries preserve a present raw nil
from a host child across subsequent copies and clones. Its field hashes like
`None`, consistent with structural equality. Shapes use nil for vacant slots,
so a live child that clones to nil becomes `None` in a shape clone. Both copy
operations start with fresh inspection history. Iterators snapshot keys
and read values from the captured backing. Callback traversal also reads from
the original backing if user code clears the receiver or moves the properties
during the callback. Multi-property slots snapshot their names before the first
callback and look up each surviving value when delivering it. The bounded
name array lives in a separate call frame used only for that case.

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

## Where the time went

In #1069, the public native lookup inlined into its caller, but the string
adapter, representation dispatch and indexed search remained separate calls.
In a CPU profile of the fixed 16- and 256-field reads, those three functions
accounted for 1.19 seconds of flat samples. The generic seeded string hash
alone accounted for 0.71 seconds including its callees. The 256-field fixture
already had a direct Go map. It still paid the common lookup dispatch on the
way to the map. The follow-up takes a native string directly to that map.
For indexed storage it uses the fingerprint and bounded search described
below. Other key types still use their existing string conversion.

In the arm64 binaries, the direct path drops one call to
`(*Object).lookup` before `runtime.mapaccess2_faststr`. The indexed path
still calls `find` and `objectNameTag`. For four- to eight-byte names the
latter uses two word loads, an add, a multiply and extraction of the top byte.
The string-to-byte conversions allocate nothing.

Reserved construction also paid for fingerprinting and probing on each
insertion. The old index had to scan on collisions. The allocation profile
of 100 escaping 128-field builds showed one handle and one co-allocated
backing per build. The follow-up's allocation profile found 100 handles,
100 owners with an inline name index, and 100 separately allocated slot
arrays for the same 100 builds. The prototype measurements below separate
the cost of finding a slot from the benefit of co-allocation.

The JSON decoder has a different path. It creates an unhinted object, then
spills its eight slots to a Go map on the ninth distinct key. The 256-field
MessagePack input supplies a size hint and goes straight to the map. The
JSON CPU profile shows token reading, map assignment and map growth. It
cannot attribute the previously reported 3.2% to a single instruction or
storage choice. The repeatable wire order and longer confirmation runs below
test that end-to-end difference directly.

Two-sample, 65 ms prototype runs used 256 unhinted writes to isolate the
promoted-map path. Moving its map pointer onto the owner changed construction
from 16.15 to 16.18 µs and added 16 measured bytes to two- and eight-field
objects. Caching the map on the receiver changed 16.13 to 16.15 µs. With a
different fingerprint the cache changed 16.17 to 15.95 µs. Neither variant
established a consistent gain. A receiver cache would also need to recognize
promotion through older aliases. These measurements cover native writes.
The decoder timings below include tokenization and value construction.

## Storage alternatives

Short, isolated screening runs on the same Go 1.27.4 arm64 host compared the
following designs. These medians come from two samples per variant. The
longer four-revision qualification below measures the selected design.

| Medium property storage | 16-field lookup | Reserved 128-field build | B/op | Allocations |
| --- | ---: | ---: | ---: | ---: |
| #1069 before this follow-up | 9.24 ns | 2.40 µs | 4,880 | 2 |
| Detached, pre-sized Go map | 8.53 ns | 2.55 µs | 9,672 | 6 |
| Linear unsorted slots | 23.26 ns | 12.59 µs | 5,104 | 3 |
| Fingerprint vector, seeded Go hash | 10.94 ns | 2.81 µs | 5,104 | 3 |
| Fingerprint vector, last four bytes | 7.60 ns | 2.10 µs | 5,104 | 3 |
| Fingerprint vector and last possible slot, full key | 6.53 ns | 1.84 µs | 5,360 | 3 |

The suffix-only variant took 552 ns for rotating hits and 1.11 µs for misses
among 128 distinct names ending in the same suffix. The selected discriminator
uses all bytes. It mixes the first and last 32-bit words of four- to eight-byte
names, folds shorter names byte by byte, and uses Go's seeded string hash for
longer names. The 256 one-byte upper bounds save a full fingerprint scan on
a new tag and give repeated reads a direct candidate. The previous collision
sentinel and deletion-time index repairs are gone.

The four co-allocated slot capacities are also gone. A medium reservation
pays one allocation for the exact requested capacity, and can release those
bytes after promotion. The standalone map exceeded Task 2's five-allocation
budget for a reserved 128-field build. Keeping the eight-slot filter made
an unhinted eight-field build take 142 ns in the filter pilot, compared with
199 ns for exact scans without it. The inline first property keeps singleton
objects, including the nested decoder fixtures, out of both larger forms.

A prior prototype that put the first property directly in a Go map needed four
allocations to grow to two fields. An external hash table used 11,024 bytes
for 128 fields. The built-in map handles unbounded growth. Skipping native
expansion entirely would also save bookkeeping, but would lose the named page
even for a fresh singleton.

## VM dictionary transition

A `FastObject` keeps its shape and slots until a distinct write crosses the
configured live-key threshold. The shape cache owns that per-VM setting.
Moving the threshold out of each fast object makes room for a two-word native
dictionary handle in the same 64-byte header on 64-bit systems. A nil shape
identifies dictionary mode. Cached shape reads and writes continue to use VM
slots. A cleared dictionary stays in dictionary mode and gets fresh storage.

On arm64 a warm cached read checks the shape and loads its slot. A cached
write still calls `SetSlotWithShape`. Neither path loads the dictionary
threshold. A distinct uncached insertion loads the threshold through the
shared cache. The constructor starts with empty slots. Each shape transition
that increases the width allocates a fresh slice. The previous reservation for
a literal's field count was discarded on its first write. In the 65-field
case, the first write on an escaping object now takes 80 B and two allocations,
down from 1,232 B and three, through either the cached or uncached write path.
Existing bytecode still carries a validated size operand. Dictionary reservation
continues to use the live field count at promotion.

Host-supplied keys can mutate the receiver while computing their string name.
Get, Lookup, ContainsKey and RemoveKey observe the storage after this conversion.
Forwarding a built-in string to the native dictionary reuses the caller's
interface value.

The first prototype held a pointer to a native object. It removed the
`fastDictionary` wrapper at the default threshold, but still allocated a
separate native handle at thresholds 1, 128 and 255. Embedding the handle by
value eliminates that allocation. Native Copy and Clone let the compiler
place the destination handle directly in its caller. Both continue to own
their child-copy semantics, including errors and a host child returning nil.

At the default threshold, reserving just 65 slots favored exactly 65 fields
but immediately spilled to a Go map at 66. Two alternating 110 ms pilot samples
from a prebuilt 64-field shape measured 2.98 µs and 12,408 B to reach 66 fields,
and 3.96 µs and 12,408 B to reach 128. Reserving 128 slots gave 1.18 µs and
5,408 B at 66, and 1.95 µs and 5,408 B at 128. At exactly 65 the larger
reservation used 376 B more than #1069's map and wrapper. The transition
doubles the current live count when it can reserve within the 128-key
inspection limit. A transition that already requires more than 128 keys
uses the direct map. The timing matrix below also includes end-to-end FQL.

The size of a shape cache grew from 32 to 40 bytes, once per VM. The following
isolated check ran 200 operations at each threshold. Copy and Clone gave the
same counts. Each cell is **B/op / allocs/op**.

| Threshold and operation | Task 2 | Pre-follow-up #1069 | Follow-up |
| --- | ---: | ---: | ---: |
| 1, transition | 336 / 2 | 352 / 3 | 368 / 2 |
| 1, copy or clone | 400 / 3 | 416 / 4 | 432 / 3 |
| 64, transition | 4,952 / 4 | 4,968 / 5 | 5,344 / 2 |
| 64, copy or clone | 5,016 / 5 | 5,032 / 6 | 2,848 / 3 |
| 128, transition | 9,560 / 4 | 9,576 / 5 | 9,560 / 4 |
| 128, copy or clone | 9,624 / 5 | 9,640 / 6 | 9,624 / 5 |
| 255, transition | 18,520 / 4 | 18,536 / 5 | 18,520 / 4 |
| 255, copy or clone | 18,584 / 5 | 18,600 / 6 | 18,584 / 5 |

At threshold 64, allocation profiles of 100 transitions attributed all
200 follow-up allocations to the native property owner and slots. The
pre-follow-up transition allocated 500 objects across those 100 operations.

## Ownership and debugger tests

`TestObjectClearDetachesSurvivingCopy` keeps the old alias and iterator alive,
forces a GC and checks a weak pointer to a replacement child. It exercises both
directions of detachment with small, indexed and directly reserved storage.
`TestObjectBackingReleasesRemovedValues` checks that removing a child releases
it while the owner is still live. A separate weak-pointer test overwrites a
substring key with an equal clone and checks reclamation through a live alias. The alias and callback tests cover mutation,
promotion during iteration, multiple clear/reuse cycles and the empty key.
The `ForEach` adapter can inline into callers. A concrete callback that does
not retain its keys can traverse the object without allocating a boxed string
per property. A collision test exercises overwrite, removal, reinsertion and
sorted paging of distinct strings with the same fingerprint. A deterministic
mixed-write test compares every lookup and length against a Go map across
uneven indexed reservations, the small group and direct-map storage.

`TestObjectBackingSparseInspection` pre-reserves a million entries or inserts
a million and deletes all but one. The sparse original rejects its named page
before invoking the visitor. Its independent copy and clone both allow
inspection. VM and debugger session tests cover exact names, selected ranges,
atomic failure and invalidation of references when a stop ends.

The VM tests make a dictionary alias before an oversized-key mutation and
verify that both aliases reject expansion. Weak-reference tests keep old
dictionaries and iterators alive across Clear. Shape traversal tests cover
growth, removal, overwrite, and clear followed by reuse. A shape callback
visits the original names. If a later slot still holds the same name it reads
the current value, including a replacement written by the callback. For a
detached or unrelated slot it reads the captured value. Debug paging reads
the captured shape and slots throughout its visit. The host-key test promotes
the object inside `String` and checks each operation against the new dictionary.
The nil-child tests check field presence, key enumeration, hash compatibility
and a second clone across native, shape and dictionary storage.

## Reproducing the comparisons

Use separate source snapshots for Task 2 (`3ac60f34`), merged #1068
(`ac2db0dc`), #1069 before this follow-up (`ccbda571`), and the follow-up.
Leave production sources and module versions at their respective revisions.
Use the follow-up's benchmark files for native object operations and rotating
lookups, FastObject, FQL literals, and both codecs in all four snapshots.
For Task 2, also copy the escaping `BenchmarkNewObjectOf` fixture and its sink.
Its test-only `SetString` adapter normalizes nil to `None` and writes directly
to `t.data[key]`, giving the native-write fixture an equivalent entry point.
Adapt the VM fixture's two constructors in test code. On historical revisions
use `NewShapeCache(threshold)` and `NewFastObject(cache, threshold)`. On the
follow-up use `NewShapeCache(threshold, threshold)` and
`cache.NewFastObject()`. Both tiny adapters inline. Keep the same benchmark
body, including the escaping sink and warmed cache.

The flat decoder fixtures use identical bytes in every process. Each encodes
the same 256 decimal property names and integers in lexical key order before
timing. The older fixture encoded a Go map whose traversal order varied.
The nested inputs are chains of 5,000 single-property objects. The VM fixtures
warm the shape cache and the FQL workloads compile before timing. Dictionary
insert/remove measures the public `RemoveKey` entry point on all revisions.
The read-heavy FQL fixture builds a 128-field dictionary and reads each of
128 parameter keys eight times. It verifies each of the 1,024 results before
starting the benchmark timer. Hits have key-specific values and misses return
`None`. Rotating native lookups exercise short decimal names, four-byte and
eight-byte names, and long names with a common suffix.

Run on the same machine with the same Go toolchain and one benchmark CPU.
Rotate the four revisions between samples, running only one benchmark binary
at a time. Run ten samples for each row, keeping one output file per revision.
In the reported matrix each operation ran on all four revisions consecutively.
A four-sequence Williams design balanced their order, and a shuffle seeded
with 1069 plus the round number determined benchmark order within each phase.
Timing starts after builds and validation finish. These commands reproduce
the workloads and parameters for one sample per row.

```sh
export GOWORK=off
go test -mod=readonly -p=1 ./pkg/runtime ./pkg/vm/internal/data ./test/benchmarks -run '^$' \
  -bench '^Benchmark(NewObjectOf|ObjectOperations|ObjectLookupKeys|FastObject|ObjectLiterals_None|DirectMutation_ObjectProperty_None|ObjectDictionaryLiterals_None|ObjectDictionaryReads_None|MemberAccess_Short_None|Match_ObjectPattern_None)$' \
  -benchmem -cpu=1 -benchtime=200ms -count=1
go test -mod=readonly -p=1 ./pkg/encoding/json ./pkg/encoding/msgpack -run '^$' \
  -bench '^Benchmark(JSONCodecDecode|MsgpackCodecDecode)$/^(flat_object_256|nested_object_5000)$' \
  -benchmem -cpu=1 -benchtime=200ms -count=1
go test -mod=readonly -p=1 ./pkg/encoding/json ./pkg/encoding/msgpack -run '^$' \
  -bench '^Benchmark(JSONCodecEncode|MsgpackCodecEncode)$/^nested_object_5000$' \
  -benchmem -cpu=1 -benchtime=200ms -count=1
benchstat task2.txt merged.txt pr.txt final.txt
benchstat merged.txt final.txt
benchstat pr.txt final.txt
benchstat task2.txt final.txt
```

Confirm the four disputed operations in longer alternating runs using
`-benchtime=1s` and the filters
`^BenchmarkObjectOperations$/^(16|256)$/^lookup$`,
`^BenchmarkObjectOperations$/^reserved_128$`, and
`^BenchmarkJSONCodecDecode$/^flat_object_256$`. Retain the raw samples when
reporting differences below a nanosecond or a few percent.

For a separate confirmation of flat JSON and the 256-field fixed-key lookup,
alternate merged #1068 and the follow-up for twenty one-second samples of
each operation. Compare those samples separately from the four-way run.

## Measured performance

All four revisions ran on the same Apple M2 Max with Go 1.27.4 on darwin/arm64.
The first table gives medians from ten interleaved one-second samples. The full
tables below give **ns/op / B/op / allocs/op** from ten interleaved 200 ms
samples. Construction escapes through a sink. Rotating hits use separate strings
with equal contents. The flat decode fixtures consume the same wire bytes in each
revision.

| One-second benchmark | Task 2 | Merged #1068 | Pre-follow-up #1069 | Revised | Revised vs merged |
| --- | ---: | ---: | ---: | ---: | ---: |
| 16-field fixed lookup | 7.83 ns | 7.87 ns | 9.11 ns | 6.65 ns | −15.5%, p < 0.001 |
| 256-field fixed lookup | 7.88 ns | 7.86 ns | 8.49 ns | 7.56 ns | −3.9%, p = 0.001 |
| Reserved 128-field build | 1.95 µs | 1.97 µs | 2.29 µs | 1.65 µs | −16.2%, p < 0.001 |
| Flat JSON, 256 fields | 56.10 µs | 56.73 µs | 57.70 µs | 58.48 µs | +3.1%, p = 0.009 |

The 128-field build uses 5,360 B and three allocations. Task 2 used 9,568 B and
five, while pre-follow-up #1069 used 4,880 B and two. The independently allocated
slots can be released on promotion. Fixed reads and rotating reads allocate zero
bytes in every revision. Flat JSON retains 1,035 allocations, two fewer than
merged #1068.

An independent twenty-sample comparison against merged #1068 confirmed the
256-field lookup improvement (7.93 to 7.64 ns, −3.7%, p = 0.008). Flat JSON moved
from 58.51 to 59.21 µs (+1.2%, p = 0.174). The shorter 200 ms matrix also failed
to resolve the flat JSON difference. The one-second four-way result supports a
small regression, but the separate pairwise run leaves its size uncertain. Keep
the unhinted singleton path for the measured nested-decoding and allocation gains.

### Native and VM results

The escaping empty object takes 64 B and two allocations. Two through eight
populated properties take 384 B and three allocations. Singleton copy and clone
each take two allocations. All native construction and operation rows met or
improved on Task 2's allocation count. In the 256-field fixture, differences from
merged #1068 in overwrite, insert/remove and iteration did not reach statistical
significance in the 200 ms run. Against Task 2, generic 256-field Set took 10.50
vs 8.75 ns and insert/remove took 27.32 vs 24.48 ns. The two-word handle
dispatches each direct-map operation according to its representation.

At a depth of 5,000 objects, nested JSON decoding took 24% less time and used half
the heap bytes of merged #1068. MessagePack took 54% less time and used 82% fewer
bytes. Each saved 10,000 allocations per decode. Flat MessagePack retained Task
2's 254 allocations.

The four-field shape build retained 224 B and five allocations. Neither its time
nor warm shape overwrite time differed significantly from merged #1068. At the
default threshold, dictionary transition fell from five allocations to two, and
copy and clone from six to three. The 65-field FQL workload constructs 16 objects
per execution. It took 1,075 allocations compared with 1,091 before the review
fixes, 1,139 on merged #1068 and 1,123 on Task 2. It used about 12 KiB less than
merged per execution. Its 4.3% lower median time was inconclusive (p = 0.529).

### Callback traversal and object construction

The 128-name callback snapshot now occupies a frame reached only for
multi-property slots. In the Go 1.27.4 arm64 executable the dispatcher uses 160
bytes of stack, down from 2,256 before the review fixes. A separate cold probe
walks a chain of 5,000 singleton objects on a fresh goroutine. A host child
samples Go's process-wide stack-in-use counter at the leaf, with a single
scheduler processor and automatic GC disabled. Three fresh processes gave
identical increments.

| Operation | Merged #1068 | Before review fixes | Revised |
| --- | ---: | ---: | ---: |
| Clone | 2,016 KiB | 16,352 KiB | 4,064 KiB |
| JSON encode | 4,064 KiB | 16,352 KiB | 8,160 KiB |

The revised chain still crosses one more Go stack growth boundary than merged
#1068. The warm JSON and MessagePack encoding fixtures ran 11.8% and 13.1% faster
than merged, respectively, with equal allocations per operation. In a separate
ten-pair comparison against the pre-review worktree, warm JSON and MessagePack
encode times fell 14.4% and 21.2%, respectively. Singleton clone, 16-field clone
and 16-field overwrite timings showed no significant change in that comparison.

### Rotating dictionary reads

With 16 entries, every four- and eight-byte hit name has a distinct fingerprint
and every miss can be ruled out by its tag. With 128 entries, the four-byte hit
names occupy 68 distinct tags, and 90 of the 128 absent names require an exact-key
search. The eight-byte names occupy 96 tags, with 78 misses requiring a search.
Compared with merged #1068, misses rose from 9.04 to 17.70 ns for four-byte names
and from 9.04 to 12.54 ns for eight-byte names. The corresponding 16-field misses
took 5.54 ns in both cases. At 128 fields the long-name hit and miss costs rose by
6.6 and 7.2 ns. All native lookup cases allocated zero bytes.

The FQL fixture exposes this cost beyond an isolated lookup. Across 1,024 long-key
reads of a 128-field dictionary, hits took 127.6 µs on merged #1068 and 130.4 µs
on the revision (+2.2%, p = 0.052). Misses took 123.7 and 132.3 µs (+6.9%, p <
0.001). Both queries used 96 allocations and 98.6 KiB, down from 102 allocations
and 109.8 KiB. Tests verified each returned value before timing.

### Merge decision

The reviewed correctness and ownership failures are fixed. The allocation target
supports keeping the detached native slots and the VM's embedded dictionary. A
pre-sized Go map used six allocations for a reserved 128-field build in the
prototype, exceeding Task 2's five. The revised implementation uses three.

Hold #1069 on read-heavy dictionary latency. The 6.9% FQL miss regression is
significant in a query with verified output. The construction, clone and decode
gains are substantial, but the fixtures supply no production workload mix that
would establish a net win for readers of medium-sized dictionaries. The
performance decision applies specifically to those indexed dictionaries.

### Codec and FQL benchmarks

| Benchmark | Task 2 | Merged #1068 | Pre-follow-up #1069 | Follow-up |
| --- | ---: | ---: | ---: | ---: |
| `ObjectLiterals` | 448 / 288 / 6 | 451 / 288 / 6 | 447 / 288 / 6 | 409 / 224 / 5 |
| `DirectMutation_ObjectProperty` | 376 / 128 / 4 | 380 / 128 / 4 | 381 / 128 / 4 | 349 / 112 / 3 |
| `ObjectDictionaryLiterals` | 415,654 / 663,276 / 1,123 | 416,308 / 663,524 / 1,139 | 415,764 / 663,528 / 1,139 | 398,454 / 651,106 / 1,075 |
| `Match_ObjectPattern` | 275 / 32 / 2 | 273 / 32 / 2 | 271 / 32 / 2 | 269 / 32 / 2 |
| `MemberAccess_Short` | 240 / 96 / 3 | 241 / 96 / 3 | 241 / 96 / 3 | 218 / 80 / 2 |
| `JSONCodecDecode/flat_object_256` | 57,758 / 51,640 / 1,036 | 60,630 / 51,648 / 1,037 | 55,324 / 51,120 / 1,035 | 59,020 / 51,120 / 1,035 |
| `JSONCodecDecode/nested_object_5000` | 1,767,463 / 2,828,288 / 20,075 | 1,798,559 / 2,878,297 / 25,075 | 1,367,466 / 1,428,291 / 15,075 | 1,363,940 / 1,428,291 / 15,075 |
| `MsgpackCodecDecode/flat_object_256` | 24,126 / 19,362 / 254 | 23,738 / 19,373 / 255 | 21,961 / 19,370 / 254 | 22,196 / 19,370 / 254 |
| `MsgpackCodecDecode/nested_object_5000` | 1,320,968 / 1,720,170 / 15,003 | 1,355,414 / 1,770,170 / 20,003 | 621,006 / 320,153 / 10,003 | 623,478 / 320,153 / 10,003 |
| `ObjectDictionaryReads_None/hit` | 127,238 / 112,420 / 101 | 127,566 / 112,438 / 102 | 126,554 / 112,438 / 102 | 130,390 / 101,012 / 96 |
| `ObjectDictionaryReads_None/miss` | 124,779 / 112,420 / 101 | 123,737 / 112,434 / 102 | 124,086 / 112,435 / 102 | 132,255 / 101,012 / 96 |
| `JSONCodecEncode/nested_object_5000` | 803,075 / 390,532 / 15,011 | 792,710 / 390,532 / 15,011 | 719,518 / 390,528 / 15,011 | 699,094 / 390,536 / 15,011 |
| `MsgpackCodecEncode/nested_object_5000` | 746,718 / 352,763 / 10,010 | 738,564 / 352,763 / 10,010 | 667,452 / 352,758 / 10,010 | 642,066 / 352,766 / 10,010 |

### Native construction and growth

| Benchmark | Task 2 | Merged #1068 | Pre-follow-up #1069 | Follow-up |
| --- | ---: | ---: | ---: | ---: |
| `NewObjectOf/0` | 32.7 / 56 / 2 | 40.5 / 66 / 3 | 27.8 / 64 / 2 | 26.6 / 64 / 2 |
| `NewObjectOf/1` | 32.5 / 56 / 2 | 40.8 / 66 / 3 | 28.0 / 64 / 2 | 27.1 / 64 / 2 |
| `NewObjectOf/256` | 873 / 18,528 / 5 | 831 / 18,538 / 6 | 894 / 18,536 / 5 | 772 / 18,536 / 5 |
| `grow_1` | 103 / 344 / 3 | 108 / 354 / 4 | 34.7 / 64 / 2 | 34.8 / 64 / 2 |
| `reserved_1` | 99.5 / 344 / 3 | 106 / 354 / 4 | 33.4 / 64 / 2 | 33.2 / 64 / 2 |
| `grow_2` | 104 / 344 / 3 | 110 / 354 / 4 | 109 / 384 / 3 | 107 / 384 / 3 |
| `reserved_2` | 103 / 344 / 3 | 108 / 354 / 4 | 106 / 384 / 3 | 106 / 384 / 3 |
| `grow_8` | 163 / 344 / 3 | 169 / 354 / 4 | 143 / 384 / 3 | 142 / 384 / 3 |
| `reserved_8` | 164 / 344 / 3 | 175 / 354 / 4 | 140 / 384 / 3 | 142 / 384 / 3 |
| `grow_16` | 993 / 2,144 / 8 | 1,012 / 2,154 / 9 | 594 / 1,624 / 7 | 582 / 1,624 / 7 |
| `reserved_16` | 369 / 1,248 / 5 | 378 / 1,258 / 6 | 330 / 1,040 / 2 | 302 / 1,008 / 3 |
| `grow_64` | 4,225 / 9,376 / 12 | 4,106 / 9,386 / 13 | 3,884 / 8,856 / 11 | 3,932 / 8,856 / 11 |
| `reserved_64` | 1,080 / 4,960 / 5 | 1,081 / 4,970 / 6 | 1,130 / 2,704 / 2 | 834 / 2,800 / 3 |
| `grow_128` | 8,217 / 18,880 / 14 | 8,180 / 18,890 / 15 | 8,066 / 18,360 / 13 | 8,092 / 18,360 / 13 |
| `reserved_128` | 2,068 / 9,568 / 5 | 2,006 / 9,578 / 6 | 2,384 / 4,880 / 2 | 1,732 / 5,360 / 3 |
| `grow_256` | 16,951 / 37,344 / 16 | 16,794 / 37,354 / 17 | 17,326 / 36,824 / 15 | 17,288 / 36,824 / 15 |
| `reserved_256` | 3,972 / 18,528 / 5 | 3,968 / 18,538 / 6 | 3,884 / 18,536 / 5 | 3,904 / 18,536 / 5 |

### Ordinary native operations

| Benchmark | Task 2 | Merged #1068 | Pre-follow-up #1069 | Follow-up |
| --- | ---: | ---: | ---: | ---: |
| `1/overwrite` | 8.10 / 0 / 0 | 8.10 / 0 / 0 | 3.35 / 0 / 0 | 3.30 / 0 / 0 |
| `1/set` | 8.24 / 0 / 0 | 9.96 / 0 / 0 | 5.80 / 0 / 0 | 4.88 / 0 / 0 |
| `1/lookup` | 7.61 / 0 / 0 | 7.60 / 0 / 0 | 3.66 / 0 / 0 | 3.98 / 0 / 0 |
| `1/insert_remove` | 21.0 / 0 / 0 | 23.5 / 0 / 0 | 12.3 / 0 / 0 | 12.3 / 0 / 0 |
| `1/iterate` | 93.7 / 80 / 3 | 94.4 / 80 / 3 | 52.1 / 80 / 3 | 51.3 / 80 / 3 |
| `1/for_each` | 30.9 / 0 / 0 | 30.5 / 0 / 0 | 4.87 / 0 / 0 | 4.92 / 0 / 0 |
| `1/copy` | 144 / 344 / 3 | 155 / 354 / 4 | 38.5 / 64 / 2 | 37.6 / 64 / 2 |
| `1/clone` | 157 / 344 / 3 | 164 / 354 / 4 | 40.3 / 64 / 2 | 40.8 / 64 / 2 |
| `16/overwrite` | 8.54 / 0 / 0 | 8.64 / 0 / 0 | 8.85 / 0 / 0 | 6.39 / 0 / 0 |
| `16/set` | 8.61 / 0 / 0 | 10.4 / 0 / 0 | 11.2 / 0 / 0 | 7.89 / 0 / 0 |
| `16/lookup` | 7.88 / 0 / 0 | 7.87 / 0 / 0 | 9.08 / 0 / 0 | 6.73 / 0 / 0 |
| `16/insert_remove` | 24.5 / 0 / 0 | 27.1 / 0 / 0 | 19.3 / 0 / 0 | 17.1 / 0 / 0 |
| `16/iterate` | 769 / 560 / 18 | 775 / 560 / 18 | 522 / 560 / 18 | 455 / 560 / 18 |
| `16/for_each` | 132 / 0 / 0 | 126 / 0 / 0 | 106 / 0 / 0 | 106 / 0 / 0 |
| `16/copy` | 1,371 / 2,144 / 8 | 1,419 / 2,154 / 9 | 356 / 1,040 / 2 | 333 / 1,008 / 3 |
| `16/clone` | 1,508 / 2,144 / 8 | 1,550 / 2,154 / 9 | 462 / 1,040 / 2 | 440 / 1,008 / 3 |
| `256/overwrite` | 8.65 / 0 / 0 | 8.90 / 0 / 0 | 9.21 / 0 / 0 | 9.14 / 0 / 0 |
| `256/set` | 8.75 / 0 / 0 | 10.5 / 0 / 0 | 10.4 / 0 / 0 | 10.5 / 0 / 0 |
| `256/lookup` | 7.95 / 0 / 0 | 7.97 / 0 / 0 | 8.57 / 0 / 0 | 7.66 / 0 / 0 |
| `256/insert_remove` | 24.5 / 0 / 0 | 27.2 / 0 / 0 | 27.1 / 0 / 0 | 27.3 / 0 / 0 |
| `256/iterate` | 11,201 / 9,008 / 258 | 11,366 / 9,008 / 258 | 11,510 / 9,008 / 258 | 11,414 / 9,008 / 258 |
| `256/for_each` | 1,724 / 0 / 0 | 1,700 / 0 / 0 | 1,784 / 0 / 0 | 1,728 / 0 / 0 |
| `256/copy` | 22,982 / 37,344 / 16 | 23,378 / 37,354 / 17 | 9,270 / 18,536 / 5 | 9,230 / 18,536 / 5 |
| `256/clone` | 25,104 / 37,344 / 16 | 25,426 / 37,354 / 17 | 9,854 / 18,536 / 5 | 10,112 / 18,536 / 5 |
| `clear_reuse` | 86.3 / 336 / 2 | 93.4 / 338 / 3 | 22.3 / 48 / 1 | 21.2 / 48 / 1 |

### VM FastObject

| Benchmark | Task 2 | Merged #1068 | Pre-follow-up #1069 | Follow-up |
| --- | ---: | ---: | ---: | ---: |
| `shape_build_4` | 323 / 224 / 5 | 325 / 224 / 5 | 320 / 224 / 5 | 321 / 224 / 5 |
| `shape_overwrite` | 8.93 / 0 / 0 | 8.64 / 0 / 0 | 8.66 / 0 / 0 | 8.25 / 0 / 0 |
| `transition` | 1,312 / 4,952 / 4 | 1,331 / 4,968 / 5 | 1,357 / 4,968 / 5 | 1,110 / 5,344 / 2 |
| `dictionary_build_65` | 17,267 / 40,280 / 69 | 17,442 / 40,296 / 70 | 17,248 / 40,296 / 70 | 17,108 / 40,672 / 67 |
| `dictionary_overwrite` | 15.9 / 0 / 0 | 9.93 / 0 / 0 | 9.99 / 0 / 0 | 7.91 / 0 / 0 |
| `dictionary_insert_remove` | 38.4 / 0 / 0 | 33.9 / 0 / 0 | 34.3 / 0 / 0 | 24.7 / 0 / 0 |
| `dictionary_copy` | 2,449 / 5,016 / 5 | 2,481 / 5,032 / 6 | 2,454 / 5,032 / 6 | 1,012 / 2,848 / 3 |
| `dictionary_clone` | 2,542 / 5,016 / 5 | 2,532 / 5,032 / 6 | 2,562 / 5,032 / 6 | 1,342 / 2,848 / 3 |
| `dictionary_clear_reuse` | 101 / 336 / 2 | 111 / 352 / 3 | 109 / 352 / 3 | 24.4 / 48 / 1 |

### Rotating lookup keys

| Benchmark | Task 2 | Merged #1068 | Pre-follow-up #1069 | Follow-up |
| --- | ---: | ---: | ---: | ---: |
| `ObjectLookupKeys/16/short/hit` | 10.0 / 0 / 0 | 9.99 / 0 / 0 | 11.5 / 0 / 0 | 8.89 / 0 / 0 |
| `ObjectLookupKeys/16/short/miss` | 7.83 / 0 / 0 | 8.92 / 0 / 0 | 9.24 / 0 / 0 | 6.59 / 0 / 0 |
| `ObjectLookupKeys/16/four_bytes/hit` | 9.21 / 0 / 0 | 9.27 / 0 / 0 | 10.4 / 0 / 0 | 8.25 / 0 / 0 |
| `ObjectLookupKeys/16/four_bytes/miss` | 8.85 / 0 / 0 | 9.00 / 0 / 0 | 9.30 / 0 / 0 | 5.54 / 0 / 0 |
| `ObjectLookupKeys/16/eight_bytes/hit` | 9.14 / 0 / 0 | 8.95 / 0 / 0 | 11.1 / 0 / 0 | 7.97 / 0 / 0 |
| `ObjectLookupKeys/16/eight_bytes/miss` | 7.92 / 0 / 0 | 7.82 / 0 / 0 | 9.26 / 0 / 0 | 5.54 / 0 / 0 |
| `ObjectLookupKeys/16/common_suffix/hit` | 10.0 / 0 / 0 | 10.1 / 0 / 0 | 12.4 / 0 / 0 | 11.8 / 0 / 0 |
| `ObjectLookupKeys/16/common_suffix/miss` | 8.78 / 0 / 0 | 8.05 / 0 / 0 | 9.37 / 0 / 0 | 10.1 / 0 / 0 |
| `ObjectLookupKeys/128/short/hit` | 10.2 / 0 / 0 | 9.78 / 0 / 0 | 18.2 / 0 / 0 | 10.2 / 0 / 0 |
| `ObjectLookupKeys/128/short/miss` | 8.89 / 0 / 0 | 8.99 / 0 / 0 | 15.1 / 0 / 0 | 10.2 / 0 / 0 |
| `ObjectLookupKeys/128/four_bytes/hit` | 9.66 / 0 / 0 | 9.78 / 0 / 0 | 21.3 / 0 / 0 | 13.1 / 0 / 0 |
| `ObjectLookupKeys/128/four_bytes/miss` | 8.94 / 0 / 0 | 9.04 / 0 / 0 | 16.1 / 0 / 0 | 17.7 / 0 / 0 |
| `ObjectLookupKeys/128/eight_bytes/hit` | 9.55 / 0 / 0 | 9.60 / 0 / 0 | 19.7 / 0 / 0 | 10.7 / 0 / 0 |
| `ObjectLookupKeys/128/eight_bytes/miss` | 8.91 / 0 / 0 | 9.04 / 0 / 0 | 15.1 / 0 / 0 | 12.5 / 0 / 0 |
| `ObjectLookupKeys/128/common_suffix/hit` | 10.4 / 0 / 0 | 10.5 / 0 / 0 | 20.1 / 0 / 0 | 17.1 / 0 / 0 |
| `ObjectLookupKeys/128/common_suffix/miss` | 9.08 / 0 / 0 | 9.20 / 0 / 0 | 15.2 / 0 / 0 | 16.4 / 0 / 0 |
