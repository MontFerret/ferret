# Ferret Core API Artifacts

Ferret publishes a versioned API Reference and sibling API Catalog with the canonical identity
`montferret/core`. This identity describes optional built-in functionality
provided by Ferret itself. It is not an installable Registry module and must not
be added to a module manifest as a dependency.

## Authority and scope

The function definitions produced by `stdlib.Full()` are the authority for the
published API. Generation enumerates that in-process registry, including every
fixed arity, variadic definition, overload, root function, and nested namespace.
Unregistered Go declarations are not published.

The generator then loads `pkg/stdlib` source with `go/packages` and resolves each
registered function to its declaration. The assertion descriptors used by `t`
and `t::not` are resolved statically, including their `Args.Min` and `Args.Max`
bounds. Any unresolved declaration, unsupported registration shape, documentation
error, arity contradiction, or runtime/source mismatch fails generation without
writing a partial artifact.

Host-function qualified names are case-insensitive in FQL and have one canonical
lowercase presentation in the registry and generated API Reference. This includes
every namespace segment; casing compatibility is represented by lookup rather
than duplicate aliases.

The API Reference, API Catalog, and discovery-index wire contracts belong to
[`github.com/MontFerret/specs`](https://github.com/MontFerret/specs). Ferret pins
the released Specs version in the independent generator and publisher modules,
which validate every completed reference and index. Ferret's root module does
not carry those tooling dependencies. Ferret does not depend on Barn for
generation or publication.

## Documentation authoring

Registered declarations use the strict structured format parsed by
`api.ParseDocumentation`:

```go
// split divides a string at each separator.
// @param value {String} Source string.
// @param separator {String} Separator string.
// @return {String[]} Split values.
// @throws {TypeError} The supplied values cannot be converted.
```

The rules are intentionally strict:

- Every generated signature has non-empty prose and exactly one `@return`.
- `@param` uses `@param name {Type} Description`; JSDoc optional-name brackets
  and dash separators are invalid.
- Fixed declarations document exactly their registered parameter count.
- Variadic declarations document at least one logical parameter and are emitted
  only as variadic signatures.
- Overload helpers carry their own complete documentation. Fixed overloads are
  ordered by arity before any variadic signature.
- Parameter and `@throws` order is retained as authored.
- Parameters are flat. Describe nested map fields in the parent parameter's
  description instead of using names such as `params.mode`.
- Assertion descriptor prose is namespace-neutral because the same descriptor
  documents both positive and `t::not` overloads. The descriptor documents its
  maximum argument list; each fixed overload receives the corresponding prefix.

Run the focused authoring and parity checks after changing stdlib registration
or documentation:

```sh
GOWORK=off go -C tools/apiref test ./internal/analyzer
```

## Local generation

Run the generator from the repository root with an unprefixed canonical SemVer:

```sh
GOWORK=off go -C tools/apiref run . \
  -version 2.0.0-alpha.45 \
  -o /tmp/montferret-core-api.json \
  -catalog /tmp/montferret-core-catalog.json
```

The command emits diagnostics only on stderr and writes deterministic,
two-space-indented JSON with one trailing newline. `api.json` remains the
canonical callable API. `catalog.json` contains presentation categories for
global and namespaced functions. Each catalog member identifies a function by
its namespace and name; the empty namespace identifies a global function.
Categories such as `math`, `io`, and `testing` are presentation concepts, not
callable Ferret namespaces. The generator contains no deployment domain.

## Release publication

The `Publish Ferret Core API Artifacts` workflow runs independently from
dependent-release notifications. It accepts `release.published` events and
manual dispatches, verifies that the selected canonical v2 tag belongs to a
published GitHub release, checks out that exact tag, strips the leading `v`, and
generates both artifacts.

The existing `gh-pages` branch is updated through the repository-local
`scripts/publish-core-api.sh` wrapper. Publication validates the existing index
and every referenced artifact before mutation, preserves unrelated Pages files,
atomically installs one version directory containing `api.json` and
`catalog.json`, recomputes the unchanged API index, commits once, and performs a
normal non-force push. Version directories, index entries, and hrefs are
immutable and cannot be republished even when bytes match. Existing API-only
version directories remain valid legacy state and are never backfilled. A stale
push fails; published history is never rewritten.

All `gh-pages` writers use the repository-wide `gh-pages-writer` concurrency key
with cancellation disabled. Unit and integration benchmarks still execute in
parallel, but their two results are persisted downstream by one serialized job.

Publication starts with the first release that contains this generator. Older
v2 alpha tags are not backfilled. A prerelease-only index omits `latest`; once a
stable release exists, `latest` identifies the greatest stable SemVer.

The public documents are:

- `https://ferretlang.org/ferret/index.json`
- `https://ferretlang.org/ferret/versions/<version>/api.json`
- `https://ferretlang.org/ferret/versions/<version>/catalog.json`

These URLs share the existing Ferret Pages root, so benchmark history and other
site files remain alongside the API Reference without being regenerated.
