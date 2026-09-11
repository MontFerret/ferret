# Array Library Contracts

The Arrays capability group registers the canonical immutable API under
`arrays::` and explicit mutation under `arrays::mut::`. Exported Go functions in
`pkg/stdlib/arrays` use the same vocabulary and semantics; mutable functions use
the `Mutable` suffix. The existing global FQL registrations are frozen migration
aliases, intended for removal in a later v2 minor release.

## Canonical vocabulary

| Arity | Operations |
| --- | --- |
| 1 | `first`, `last`, `unique`, `sorted` |
| 1 or 2 | `flatten(array[, depth])` |
| 2 | `at`, `append`, `contains`, `index_of`, `remove`, `remove_at`, `remove_any` |
| 2 or 3 | `slice(array, start[, length])` |
| At least 2 | `concat`, `union`, `intersection`, `difference`, `symmetric_difference` |

`append` appends one value, including an array as one nested element. `concat`
concatenates inputs without deduplication. `contains` returns Boolean;
`index_of` returns the first equal value's index or `-1`. `remove` removes every
match. `remove_any` removes every match to any value in its second list, retaining
other duplicates. `remove_at` removes one indexed element or returns an
unchanged shallow copy for a missing index. No canonical operation accepts a
Boolean mode switch or removal limit.

`flatten` retains its default depth of one; zero and negative depths do not
expand nested lists. `slice` retains its start/length convention, not an end-index
argument. Negative starts or lengths and starts beyond the list return an empty
array. Oversized lengths are capped using a subtraction comparison before
addition, preventing integer overflow. The exported Go `Slice` now rejects
undocumented extra arguments, matching its FQL overloads.

## Explicit mutation

Mutable operations change the original list. Aliases observe those changes;
there is no implicit copy and no global mutable registration or `mut::` alias.

| Call under `arrays::mut::` | Result |
| --- | --- |
| `push(array, value)` | Original array with one appended element |
| `unshift(array, value)` | Original array with one prepended element |
| `set(array, index, value)` | Original array with an existing element replaced |
| `insert(array, index, value)` | Original array with an element inserted |
| `remove(array, value)` | Original array with all equal values removed |
| `clear(array)` | Original array with all elements removed |
| `sort(array)` | Original array, stably sorted in ascending FQL order |
| `pop(array)` | Removed last value, or `none` when empty |
| `shift(array)` | Removed first value, or `none` when empty |
| `remove_at(array, index)` | Removed value |

Indexes must be Int. `set` and mutable `remove_at` require an existing index;
`insert` allows zero through length inclusive. Negative or missing indexes fail
before mutation with `ErrInvalidOperation` attributed to the index argument.
`set` does not grow the array. Immutable `remove_at` retains its unchanged-copy
behavior for missing indexes. There are no optional modes, removal limits, or
sort directions in this namespace.

Targets must implement `runtime.List`, whose existing contract includes every
required mutation method. Missing capabilities and nil targets fail with a type
error. Host mutation refusals propagate without copying the target; empty
`pop`/`shift` and removal with no matches do not issue speculative writes.
Negative host lengths are rejected wherever length is consumed.

The wrappers use existing runtime append, insert, indexed removal, clear, and
sorting primitives. Mutable and immutable sorting share `runtime.SortAsc`.
Immutable and legacy removal keep their host `Filter` dispatch. Unlimited
removal, including negative legacy limits, shares a stateless predicate with
mutable removal. Zero and positive legacy limits use a separate counted
predicate; equality is still evaluated after the limit is reached so host
comparison errors are preserved. Both predicates use `runtime.EqualValues`.
Mutable removal compacts survivors forward and trims the tail,
preserving survivor order without quadratic native middle removals.

The caller's context reaches host operations and comparisons. Mutable wrappers
reject pre-canceled calls and poll cancellation during removal. Errors preserve
their causes and return `runtime.None`; a host, comparison, or cancellation
failure can leave earlier writes applied. There is no rollback or new panic
recovery policy. Host implementations must honor their runtime contracts.
Targets and element values are borrowed: mutation neither clones nor closes
inserted, retained, or removed values.

## Equality and deterministic ordering

`unique` and `union` retain the first occurrence of each distinct value.
Intersection and difference follow the first input's order. Symmetric difference
retains values present in an odd number of input arrays; duplicates within each
input count once. Its result follows first encounter across inputs and retains
the first representative even when later equivalent representations toggle
membership. For example, `symmetric_difference([7,3], [7], [7])` is `[7,3]`.

`pkg/internal/valueset` owns the shared hash-bucket membership implementation.
Hash matches always require `runtime.EqualValues`; callers keep ordered output
separately from hash maps. Membership, removal, and insertion preserve equality
errors. Sets retain borrowed value references and do not clone or close them.
Arrays continue supporting `runtime.List` host implementations, passing the
caller's context to list and comparison operations. Runtime comparison dispatch
does not impose cancellation polling on host values.

## Copy and index ownership

Immutable array transformations do not modify their inputs. Native array copies
and slices have independent top-level backing storage, including when a slice
could otherwise reuse spare capacity. Copying remains shallow: nested arrays,
objects, and resources retain their identities. Slicing does not clone, acquire,
or close nested resources. Host lists remain responsible for their own `Copy`,
`Filter`, and `Slice` contracts.

Native `Array` enforces bounds itself:

- `At` returns `None` for negative or missing indexes; `LookupAt` also reports
  `found=false`.
- `RemoveAt` returns `None` without changing the array for missing indexes.
- `SetAt` and `Swap` require existing indexes. `Insert` accepts indexes from zero
  through length, inclusive. Invalid writes return `ErrInvalidOperation` before
  modifying storage.
- Runtime `Slice` uses `[start, end)`. Negative or reversed bounds and starts
  beyond the array return empty arrays; oversized ends are clamped.

These primitive guards do not alter the VM's existing indexed-read or
indexed-write policies. Legacy `pop` checks for an empty list before requesting
a slice, including for host lists.

## Frozen global compatibility

Matching global contracts bind the canonical implementation directly. Only
legacy-specific signatures and behavior use private adapters. Registry comments
document alias migration; adapters carry structured `@deprecated` documentation.
The registry has no per-alias deprecation metadata, so direct aliases share the
canonical function's documentation without deprecating that function.

- `union` still concatenates; `union_distinct` is distinct union.
- `nth`, `minus`, `remove_nth`, and `remove_values` map to `at`, `difference`,
  `remove_at`, and `remove_any`.
  Global `remove_nth` still delegates bounds policy directly to the copied host
  list; canonical `remove_at` checks its length before removing an element.
- `position` retains its optional Boolean/index switch.
- `append` and `push` retain the optional unique switch, which skips an existing
  target but does not remove duplicates already in the input.
- `remove_value` retains its optional limit: negative is unlimited, zero removes
  nothing, and positive limits removals in encounter order.
- `unshift` retains its optional switch: the true form prepends the value and
  removes existing matches to that value, preserving unrelated duplicates.
- `sorted_unique`, `pop`, and `shift` remain global-only.
- `outersection` retains values present in exactly one input array, counting
  duplicates within an input once. A value present in three arrays is excluded,
  unlike symmetric difference. Previously unspecified set output order is now
  deterministic.

Obsolete Go entry points are removed, including the former mode-bearing
functions and `ToUniqueList`. This Go API migration is intentional; compatibility
is supplied through global FQL registrations, not deprecated exported Go shims.

## Verification

Array package tests cover canonical registrations and direct Go calls; preserved
legacy tests invoke the actual runtime registry. VM integration covers both
`None` and `Full` optimization levels. Runtime tests exercise extreme bounds and
independent slice mutations. Valueset tests cover collisions, numeric/nested
equality, failed removal, cancellation propagation, and reuse after removal.
API generator tests check names, signatures, categorization, and adapter
deprecation metadata. Benchmark slice copying and set operations alongside
existing DISTINCT consumers whenever the common set implementation changes.

Mutable tests exercise exact result identity, strict bounds, host refusals,
partial failures, cancellation, shallow resource ownership, and mutation of
independent slices and copies. FQL mutation tests run at None, Basic, and Full
optimization. Copy, removal, sorting, and native/custom mutation benchmarks
track time and allocations; mutable benchmark setup is excluded from timing.
