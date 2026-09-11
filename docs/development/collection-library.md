# Collection library contracts

The Collections group retains four global functions: unary `count`,
`count_distinct`, and `reverse`, and binary `includes`. Their Go implementations
live in `pkg/stdlib/collections`. FQL argument validation belongs here; collection
construction and canonical equality remain runtime contracts.

## Counting read-only sources

`count` and `count_distinct` require only `runtime.Iterable` in addition to the
`runtime.Value` input contract. Read-only host values do not need collection
mutation, cloning, membership, equality, or ordering methods. Native ascending
and descending ranges are supported. Strings and non-iterable scalars remain
invalid counting inputs.

After iterable validation, `count` calls `Measurable.Length` when available.
It does not create an iterator or retry a failed length operation by scanning.
Native ranges use this path, including propagation of range-length overflow.
A host length can perform I/O or take nonconstant time.

Without measurement, `count` traverses once through `runtime.ForEach`. Checked
increments return `ErrRange` before exceeding `runtime.Int`. No partial count is
returned after a failure. Scanning may consume a one-shot source, perform I/O,
or fail to terminate on an unbounded source.

`count_distinct` always traverses yielded values. Objects contribute values,
not keys or key/value pairs. It uses `pkg/internal/valueset`: hashes select
candidates, and canonical equality resolves collisions. Equivalent Int/Float
values share a distinct entry, Duration equality remains strict, and host
equality failures propagate. Stored references are borrowed, never cloned or closed.

## Membership and traversal ownership

`includes` dispatches to string handling first, then `runtime.Containable`, then
an iterable scan. String needles retain textual conversion, so
`includes("123", 123)` is true; the separate string `contains` remains strict.
Object membership searches values. A matching key alone is insufficient.

Fallback scans use `runtime.ForEach`, which closes only its acquired closable
iterator exactly once on exhaustion, early match, or failure. Iterator creation
failure acquires no cleanup obligation. Traversal/equality and close errors are
joined, including a close failure after a match. Sources and yielded values
remain borrowed. Delegated `Contains` owns its own internal resources.

The counting and membership wrappers check cancellation before host dispatch,
within traversal, and before successful completion. Equality/set operations are
followed by cancellation checks, and cancellation during EOF or iterator cleanup
cannot produce success. Hosts receive the caller's context and remain responsible
while they retain control. Entry checks cannot interrupt arbitrary blocking host
operations. Global iterator and equality semantics are unchanged.

## Reversal and destination ownership

String reversal reverses Unicode code points, not grapheme clusters. List
reversal obtains length, calls `source.New(ctx)`, reads descending indexes, and
checks every append. Negative host lengths fail before construction. It returns
the same implementation family with the source's
relevant configuration. Only the outer list is new; element references remain
shallow and the source is unchanged. Non-list iterables are not materialized.

Factories own construction failure, as specified by
[Factory](runtime.md#collection-construction-and-migration). After successful
creation, reversal owns the destination until success. Access, append, and
cancellation failures close that incomplete destination when closable and join
cleanup errors with the primary failure. Successful results transfer ownership
through normal VM/result lifecycle handling. The borrowed source and element
values are never explicitly closed.

## Verification and compatibility

Package tests use minimal read-only sources, configured lists/maps, deterministic
cancellation hooks, and closable iterators. Tests assert exact cleanup counts,
joined error identity, shallow references, and independent destinations. VM
contracts run at None, Basic, and Full optimization; API analyzer tests verify
capability types and arities. Benchmarks cover measured counting, distinct scans,
membership dispatch, and native/custom reversal, with new iterable-only counting
reported separately from previously supported paths.

The Factory rename intentionally breaks Go implementers of List and Map; existing
FQL names and arities are retained. Counting accepts additional read-only inputs,
and global reverse now preserves custom list implementations. `arrays::reverse`
and `is_empty` are outside this migration.
