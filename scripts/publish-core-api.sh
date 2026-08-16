#!/bin/sh

set -eu

if [ "$#" -lt 3 ] || [ "$#" -gt 5 ]; then
	printf 'Usage: %s <api-reference> <api-catalog> <gh-pages-checkout> [remote] [branch]\n' "$0" >&2
	exit 2
fi

script_dir=$(CDPATH='' cd -- "$(dirname "$0")" && pwd)
source_root=$(CDPATH='' cd -- "$script_dir/.." && pwd)
reference_dir=$(CDPATH='' cd -- "$(dirname "$1")" && pwd)
reference="$reference_dir/$(basename "$1")"
catalog_dir=$(CDPATH='' cd -- "$(dirname "$2")" && pwd)
catalog="$catalog_dir/$(basename "$2")"
pages_root=$(CDPATH='' cd -- "$3" && pwd)
remote=${4:-origin}
branch=${5:-gh-pages}

if [ ! -f "$reference" ]; then
	printf 'API Reference not found: %s\n' "$reference" >&2
	exit 1
fi

if [ ! -f "$catalog" ]; then
	printf 'API Catalog not found: %s\n' "$catalog" >&2
	exit 1
fi

if ! git -C "$pages_root" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
	printf 'Pages checkout is not a Git worktree: %s\n' "$pages_root" >&2
	exit 1
fi

if [ -n "$(git -C "$pages_root" status --porcelain)" ]; then
	printf 'Pages checkout must be clean before publication: %s\n' "$pages_root" >&2
	exit 1
fi

version=$(jq -er '.version | strings | select(length > 0)' "$reference")

go -C "$source_root/tools/apipublish" run . \
	-reference "$reference" \
	-catalog "$catalog" \
	-pages "$pages_root"

reference_artifact="versions/$version/api.json"
catalog_artifact="versions/$version/catalog.json"
git -C "$pages_root" add -- index.json "$reference_artifact" "$catalog_artifact"

if git -C "$pages_root" diff --cached --quiet; then
	printf 'Publication produced no staged changes for %s\n' "$version" >&2
	exit 1
fi

git -C "$pages_root" commit -m "Publish montferret/core API artifacts $version" -- index.json "$reference_artifact" "$catalog_artifact"
git -C "$pages_root" push "$remote" "HEAD:refs/heads/$branch"
