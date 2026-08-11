#!/bin/sh

set -eu

if [ "$#" -lt 2 ] || [ "$#" -gt 4 ]; then
	printf 'Usage: %s <api-reference> <gh-pages-checkout> [remote] [branch]\n' "$0" >&2
	exit 2
fi

script_dir=$(CDPATH='' cd -- "$(dirname "$0")" && pwd)
source_root=$(CDPATH='' cd -- "$script_dir/.." && pwd)
reference_dir=$(CDPATH='' cd -- "$(dirname "$1")" && pwd)
reference="$reference_dir/$(basename "$1")"
pages_root=$(CDPATH='' cd -- "$2" && pwd)
remote=${3:-origin}
branch=${4:-gh-pages}

if [ ! -f "$reference" ]; then
	printf 'API Reference not found: %s\n' "$reference" >&2
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
	-pages "$pages_root"

artifact="versions/$version/api.json"
git -C "$pages_root" add -- index.json "$artifact"

if git -C "$pages_root" diff --cached --quiet; then
	printf 'Publication produced no staged changes for %s\n' "$version" >&2
	exit 1
fi

git -C "$pages_root" commit -m "Publish montferret/core API $version" -- index.json "$artifact"
git -C "$pages_root" push "$remote" "HEAD:refs/heads/$branch"
