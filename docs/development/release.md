# Release Automation

This repository defines post-release automation for Ferret Core. It does not
define the release-creation or binary-packaging process for the separate
MontFerret CLI product.

Published GitHub releases with v2 tags trigger two independent workflows:

```text
published v2 release
    |-> notify dependent source repositories
    \-> publish versioned Ferret Core API artifacts
```

Both workflows also support explicitly validated manual dispatches. Their YAML
definitions are the authority for accepted inputs, permissions, target
repositories, and event payloads.

## Dependent repository notification

`.github/workflows/notify-dependents.yml` validates a v2 SemVer tag and either a
single allowed dependent repository or the complete configured set. Validation
finishes before authentication or dispatch.

For each selected repository, the workflow creates a repository-scoped GitHub
App token with content-write permission and sends a `repository_dispatch` event
named `ferret-core-released`. The payload contains the validated release version.
Each target repository owns its response to that event.

The configured repository names and event literal are security- and
compatibility-sensitive. Keep validation, the target allowlist, token scope, and
payload construction together. Do not accept arbitrary repository names, log
tokens or private keys, or include rejected raw input in downstream requests.

The matrix does not fail fast, so notification results remain visible for every
dependent. A failed dispatch must fail that matrix job rather than being reported
as a successful notification.

## Core API artifact publication

`.github/workflows/publish-core-api.yml` validates that its v2 SemVer tag names a
published, non-draft GitHub release. It then checks out the exact release tag and
generates two deterministic artifacts from that source:

* the callable API Reference;
* the sibling presentation Catalog.

The generator is the independent `tools/apiref` module. The registrations in
`stdlib.Full()` are the callable API authority; the generator's source analyzer
resolves and validates the structured documentation associated with those
registrations. Ferret pins the shared artifact contracts in the tool module,
not in the root runtime module.

Detailed authoring, local generation, validation, and wire-contract rules live
in the [Core API artifact maintainer guide](../maintainers/core-api-reference.md).

## Immutable `gh-pages` publication

After generation, the workflow checks out the existing `gh-pages` branch into a
separate worktree and invokes `scripts/publish-core-api.sh`. The script requires:

* existing reference and catalog files;
* a clean Git worktree for the Pages branch;
* a non-empty version in the reference;
* a publication that produces staged changes.

The script runs the independent `tools/apipublish` module, then commits only the
discovery index and the versioned `api.json` and `catalog.json` paths. Publication
preserves immutable version history and pushes the explicit branch ref. It does
not use Barn, modify the release tag checkout, or write partial artifacts when
validation fails.

`publish-core-api.yml` and the benchmark workflow both write to `gh-pages`. They
share the `gh-pages-writer` concurrency group with cancellation disabled so one
writer cannot discard another writer's completed history.

## Security and failure boundaries

Release workflows should keep event names as fixed literals and validate
user-controlled tags and target repositories before using credentials. Use the
minimum token permissions and repository scope required by each operation. Never
print GitHub App private keys, access tokens, authorization headers, or other
credential-bearing values.

The automatic notification path runs only for a published release event. Its
manual recovery path validates tag syntax and the target allowlist but does not
query GitHub to confirm that the supplied tag is published. In contrast, manual
API publication explicitly verifies the published, non-draft release before
checking out the tag. Keep that distinction visible when changing either flow.

Recovery paths must not silently overwrite divergent history or treat partial
publication as success.

Artifact generation and validation complete before the publisher stages or
pushes files. The Pages worktree must begin clean, and no-op publication is an
error rather than evidence that a new version was published.

## Documentation synchronization

Public language, embedding, SDK, and standard-library behavior is documented in
the website repository. Changes to public behavior should update the relevant
website documentation in the same task when that checkout is available. If it
is unavailable, record the specific required follow-up.

The API Reference and Catalog are generated from registered code and structured
source documentation. They complement, but do not replace, explanatory website
documentation.

## Related guides

* [Architecture](architecture.md)
* [Modules, SDK, and standard library](modules.md)
* [Development workflow](workflow.md)
* [Core API artifact maintainer guide](../maintainers/core-api-reference.md)
