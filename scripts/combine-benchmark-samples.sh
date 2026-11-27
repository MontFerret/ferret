#!/usr/bin/env bash
set -euo pipefail

directory=$1
suite=$2
output=$3
files=()

for sample in 1 2 3 4 5; do
  file="$directory/benchmark-output-main-$suite-$sample/${output##*/}"
  if [[ ! -s "$file" ]] || ! grep -Eq '^Benchmark[^[:space:]]*[[:space:]]+[[:digit:]]+[[:space:]]' "$file"; then
    echo "Missing benchmark results for $suite sample $sample in $file" >&2
    exit 1
  fi
  files+=("$file")
done

cat "${files[@]}" > "$output"
