# Object library contracts

The Objects capability group registers ten immutable functions under
`object::`: `keys`, `values`, `entries`, `has_key`, `keep_keys`,
`omit_keys`, `merge`, `merge_deep`, `zip`, and `from_entries`.
Old global names are removed. Host-function lookup remains case-insensitive;
registration and generated API metadata use lowercase names. It also registers
`object::mut::{merge,merge_deep,keep_keys,omit_keys}`. The stdlib convention is
immutable by default, with explicit mutation through a `::mut::` subnamespace.

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

Mutable merges require a target map followed by variadic sources or one list of
sources. Target-only calls and explicit empty source lists are valid no-ops.
Mutable key filters retain the same required key argument and explicit-empty
semantics as immutable filters. Every mutable call returns the exact target
value; top-level aliases observe its changes. All arguments are normalized and
validated before the first mutation.

`zip` accepts two runtime lists of equal length and requires String keys.
`from_entries` consumes any runtime iterable once. Each entry must implement
`runtime.Measurable` and `runtime.IndexReadable`, have length two, and contain
a String key at index zero. Both constructors use last-key-wins; this deliberately
changes the former global ZIP's first-key-wins behavior.

## Ownership and shared operations

Runtime owns `MergeMapsInto`, `MergeMapsDeepInto`,
`MergeMapsDeepIsolatedInto`, `KeepMapKeys`, and `OmitMapKeys`.
FQL arity, list-versus-variadic normalization, and argument
attribution remain in stdlib. These operations use runtime interfaces without
depending on concrete Objects. Existing `Map.Merge` behavior is unchanged.

Immutable merges create an independent destination and clone or copy incoming
values once before insertion, without a final redundant clone. Deep merges
operate only on destination-owned nested maps. Key filters clone the complete
source once before applying the shared removal operation; a clone failure is
reported even when the failing value belongs to a key that would be removed.
Key snapshots are fully traversed before removal, supporting live host key views.

Immutable results have independent containers. Values use `runtime.CloneOrCopy`:
Cloneable values must honor their deep-clone contract; other values follow
their shallow `Value.Copy` contract. The library cannot strengthen a host
value's copy guarantees. Host implementations of `Empty` and `Clone` must
produce independent destinations.

Mutable shallow merge and key filters use the same runtime operations directly
on their target, without cloning it. Retained existing values are untouched.
`runtime.Map` already includes `Set` and `RemoveKey`; readable values lacking
that contract fail with a type error. Host mutation refusals propagate, with no
copy fallback. No-op calls do not test mutability by issuing speculative writes.

Mutable deep merge uses `MergeMapsDeepIsolatedInto`. It shares traversal,
conflict rules, and recursive merging with the exclusively owned helper.
Conflicting nested branches are cloned before modification, merged through
`MergeMapsDeepInto`, then replaced on the target after success. The root keeps
its identity; original nested branches, including branches shared by sources or
other target keys, remain unchanged. Untouched branches retain their identity.
This relies on the same host clone/copy contracts as immutable transformations.

Destination helpers borrow sources and may partially update the destination
on failure. A failed nested merge does not replace its original branch, but
earlier successful updates can remain. No transaction or rollback is promised,
including when a host mutation itself fails after changing state. Targets,
sources, and original nested branches remain borrowed by the operation.

## Errors, iteration, and validation

Failed calls return `runtime.None`; immutable calls preserve their sources while
mutable calls can leave their target partially updated.
Runtime errors retain their causes with argument and key/entry context.
Entry iteration closes the acquired iterator exactly once when it implements
`io.Closer`; the source iterable remains borrowed. Runtime traversal preserves
both iteration/predicate and close errors through `errors.Join`.

Semantic fixtures in `test/spec/objects` are shared by runtime destination
tests and both immutable and mutable wrapper tests. Host-value tests cover
copy-only and Cloneable values, fallible access and mutation, one-shot entry iteration,
cancellation, and cleanup. FQL integration covers None, Basic, and Full
optimization. Object benchmarks cover transformation time and allocations.

CLI source migration is owned by the sibling CLI repository. It rewrites
unqualified legacy calls and composes literal sorted KEYS calls with `sorted`.
Unsafe sorting expressions produce the established manual action, leaving the
file untouched. ZIP's duplicate-key change is documented rather than hidden
behind an alias.
