# Native object inspection allocation cost

The target is zero additional allocations in ordinary execution when the
debugger is unused. Native `runtime.Object` currently allocates one shared
record per backing map to bound debugger enumeration. Moving the record to
the first inspection would leave the debugger unable to tell whether it can
safely discover even one property name.

## Cost without a debugger

The table gives allocations per operation on Go 1.27.3, darwin/arm64.

| Operation | Task 2 (`3ac60f34`) | Current | Summary-only control |
| --- | ---: | ---: | ---: |
| Construct an empty native object | 2 | 3 | 2 |
| JSON decode, 256 properties | 1,036 | 1,037 | 1,036 |
| MessagePack decode, 256 properties | 254 | 255 | 254 |
| JSON decode, 5,000 nested objects | 20,075 | 25,075 | 20,075 |
| MessagePack decode, 5,000 nested objects | 15,003 | 20,003 | 15,003 |

The control removed native map history in an isolated Go test overlay and
returned a summary for every native object. This recovered the allocation
baselines, including 56 B per empty object versus 66 B currently.
`TestObjectClearDetachesSurvivingCopy` passed. `TestDebugNamedPage/map`
failed at the first property with `errors.ErrUnsupported`. A debug session
given the fresh and pre-reserved singletons as host parameters described
both as `Object(1)`. The current implementation returned an expansion handle
for the fresh singleton only. The control also lost that handle. The codec and
constructor benchmark commands are in the [debugger guide](debugger.md#measuring-a-fixed-page).

## Work before the first page

Every fixture below contains exactly `{"keep": true}` when inspection starts.
The first-key benchmark stops the Go range as soon as it yields a name. This
omits sorting, formatting, and handle creation. The known-key benchmark
directly looks up the same name to isolate the cost after discovery.

| Backing history | First key, Go 1.26.9 | First key, Go 1.27.3 |
| --- | ---: | ---: |
| Fresh singleton | 29 ns | 27 ns |
| Reserved for 65,536 entries | 26.5 µs | 26.7 µs |
| Reserved for 1,000,000 entries | 1.02 ms | 1.04 ms |
| Grew to 1,000,000, then deleted 999,999 | 1.10 ms | 1.02 ms |

These are medians of five samples on an Apple M2 Max. Known-key lookup took
8–10 ns across the fixtures. Both operations allocated zero bytes per
iteration. Fixture construction and deletion were outside the timed loops.

```sh
go test ./pkg/runtime -run '^$' -bench '^BenchmarkObjectInspectionKeyDiscovery$' \
  -benchmem -cpu=1 -benchtime=150ms -count=5
```

The supported Go map API exposes the live count. Its iterator can cross
empty groups and tables before yielding control to the caller. Both tested
toolchains implement that traversal in `internal/runtime/maps.Iter.Next`.
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

The allocation remains open. A replacement must recover the ordinary
constructor and codec baselines, expand a fresh small native object, and
bound the first page on an indistinguishable sparse singleton. The existing
alias and `Clear` tests remain part of that acceptance check.
