# Object library contracts

The Objects capability group registers exactly ten immutable functions under
`object::`: `keys`, `values`, `entries`, `has_key`, `keep_keys`,
`omit_keys`, `merge`, `merge_deep`, `zip`, and `from_entries`.
Old global names are removed. Host-function lookup remains case-insensitive;
registration and generated API metadata use lowercase names.

## Arguments and results

Map operations accept `runtime.Map`, including host implementations. `keys`
has one argument; callers compose `sorted(object::keys(value))` for ordering.
No ordering is promised for keys, values, or entries. Entries traverses each
map once and constructs pairs directly, preserving key/value association.
Each entry contains a String key and a cloned or copied value. Keys are validated
with `runtime.CastString` and are never cloned or copied; a host map exposing a
non-string key fails the call with an argument-attributed type error.

Both merge functions accept variadic maps or one `runtime.List` of maps.
At least one argument is required; an empty list produces an empty object.
The first source's `Empty` method creates the destination. Later sources win.
Deep merge recurses only when both conflicting values implement `runtime.Map`;
arrays, scalars, and `none` are replaced.

Both key filters require a map and at least one key argument. Keys are variadic
Strings or one runtime list of Strings. Missing and repeated keys have no effect.
An explicit empty list keeps nothing or omits nothing.

`zip` accepts two runtime lists of equal length and requires String keys.
`from_entries` consumes any runtime iterable once. Each entry must implement
`runtime.Measurable` and `runtime.IndexReadable`, have length two, and contain
a String key at index zero. Both constructors use last-key-wins; this deliberately
changes the former global ZIP's first-key-wins behavior.

## Ownership and shared operations

Runtime owns `MergeMapsInto`, `MergeMapsDeepInto`, `KeepMapKeys`, and
`OmitMapKeys`. FQL arity, list-versus-variadic normalization, and argument
attribution remain in stdlib. These operations use runtime interfaces without
depending on concrete Objects. Existing `Map.Merge` behavior is unchanged.

Immutable merges create an independent destination and clone or copy incoming
values once before insertion, without a final redundant clone. Deep merges
operate only on destination-owned nested maps. Key filters clone the complete
source once before applying the shared removal operation; a clone failure is
reported even when the failing value belongs to a key that would be removed.
Key snapshots are fully traversed before removal, supporting live host key views.

Returned containers are independent. Values use `runtime.CloneOrCopy`:
Cloneable values must honor their deep-clone contract; other values follow
their shallow `Value.Copy` contract. The library cannot strengthen a host
value's copy guarantees. Host implementations of `Empty` and `Clone` must
produce independent destinations.

Destination helpers borrow sources and may partially update the destination
on failure. Their caller owns the destination and any nested maps mutated by
deep merge. Future `object::mut` wrappers can reuse these operations without
the immutable wrapper's destination creation or clone. No mutable FQL API
or transactional mutation guarantee is introduced here.

## Errors, iteration, and validation

Failed immutable calls return `runtime.None` and preserve their sources.
Runtime errors retain their causes with argument and key/entry context.
Entry iteration closes the acquired iterator exactly once when it implements
`io.Closer`; the source iterable remains borrowed. Runtime traversal preserves
both iteration/predicate and close errors through `errors.Join`.

Semantic fixtures in `test/spec/objects` are shared by runtime destination
tests and immutable wrapper tests. Host-value tests cover copy-only and
Cloneable values, fallible access and mutation, one-shot entry iteration,
cancellation, and cleanup. FQL integration covers None, Basic, and Full
optimization. Object benchmarks cover transformation time and allocations.

CLI source migration is owned by the sibling CLI repository. It rewrites
unqualified legacy calls and composes literal sorted KEYS calls with `sorted`.
Unsafe sorting expressions produce the established manual action, leaving the
file untouched. ZIP's duplicate-key change is documented rather than hidden
behind an alias.
