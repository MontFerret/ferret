window.BENCHMARK_DATA = {
  "lastUpdate": 1773948092229,
  "repoUrl": "https://github.com/MontFerret/ferret",
  "entries": {
    "Ferret Go Benchmarks - Unit": [
      {
        "commit": {
          "author": {
            "email": "ziflex@users.noreply.github.com",
            "name": "Tim Voronov",
            "username": "ziflex"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4dfafb529e96e2598a0f600da612ff2feb530f51",
          "message": "Chore/struct split (#899)\n\n* refactor: restructure `mem` package by splitting `CellStore` and `CellHandle` into separate files\n\n* refactor: split `collector` code into smaller modules for improved maintainability\n\n* docs: clarify commenting guidelines and update formatting in AGENTS.md\n\n* refactor: consolidate type and struct definitions across core packages\n\n* docs: add benchmarking guidelines for significant changes in AGENTS.md\n\n* docs: revise type declaration guidelines in AGENTS.md for clarity and alignment\n\n* refactor: enhance optimization by introducing use-def collector and basic block structures\n\n* chore: expand and reorganize benchmarks with unit and integration workflows\n\n* Potential fix for pull request finding\n\nCo-authored-by: Copilot Autofix powered by AI <175728472+Copilot@users.noreply.github.com>\n\n* refactor: enhance dependency management in optimization pipeline with fresh metadata handling\n\n* chore: reorganize benchmark workflows and refactor mapsEqual function\n\n* chore: update integration benchmark label in workflow configuration\n\n---------\n\nCo-authored-by: Copilot Autofix powered by AI <175728472+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-19T14:39:43-04:00",
          "tree_id": "ec1c2cfce7100e21957ad36f7c2c1625d72ddcb5",
          "url": "https://github.com/MontFerret/ferret/commit/4dfafb529e96e2598a0f600da612ff2feb530f51"
        },
        "date": 1773945818671,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 47470,
            "unit": "ns/op\t  84.50 MB/s\t   11152 B/op\t     932 allocs/op",
            "extra": "25255 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 47470,
            "unit": "ns/op",
            "extra": "25255 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 84.5,
            "unit": "MB/s",
            "extra": "25255 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 11152,
            "unit": "B/op",
            "extra": "25255 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 932,
            "unit": "allocs/op",
            "extra": "25255 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 47449,
            "unit": "ns/op\t  84.53 MB/s\t   11152 B/op\t     932 allocs/op",
            "extra": "25052 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 47449,
            "unit": "ns/op",
            "extra": "25052 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 84.53,
            "unit": "MB/s",
            "extra": "25052 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 11152,
            "unit": "B/op",
            "extra": "25052 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 932,
            "unit": "allocs/op",
            "extra": "25052 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 48568,
            "unit": "ns/op\t  82.58 MB/s\t   11152 B/op\t     932 allocs/op",
            "extra": "24796 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 48568,
            "unit": "ns/op",
            "extra": "24796 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 82.58,
            "unit": "MB/s",
            "extra": "24796 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 11152,
            "unit": "B/op",
            "extra": "24796 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 932,
            "unit": "allocs/op",
            "extra": "24796 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 48067,
            "unit": "ns/op\t  83.45 MB/s\t   11152 B/op\t     932 allocs/op",
            "extra": "25116 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 48067,
            "unit": "ns/op",
            "extra": "25116 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 83.45,
            "unit": "MB/s",
            "extra": "25116 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 11152,
            "unit": "B/op",
            "extra": "25116 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 932,
            "unit": "allocs/op",
            "extra": "25116 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 47650,
            "unit": "ns/op\t  84.18 MB/s\t   11152 B/op\t     932 allocs/op",
            "extra": "24042 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 47650,
            "unit": "ns/op",
            "extra": "24042 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 84.18,
            "unit": "MB/s",
            "extra": "24042 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 11152,
            "unit": "B/op",
            "extra": "24042 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 932,
            "unit": "allocs/op",
            "extra": "24042 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 35326,
            "unit": "ns/op\t  66.27 MB/s\t   12867 B/op\t     422 allocs/op",
            "extra": "33895 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 35326,
            "unit": "ns/op",
            "extra": "33895 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 66.27,
            "unit": "MB/s",
            "extra": "33895 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 12867,
            "unit": "B/op",
            "extra": "33895 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 422,
            "unit": "allocs/op",
            "extra": "33895 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 35800,
            "unit": "ns/op\t  65.39 MB/s\t   12867 B/op\t     422 allocs/op",
            "extra": "34452 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 35800,
            "unit": "ns/op",
            "extra": "34452 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 65.39,
            "unit": "MB/s",
            "extra": "34452 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 12867,
            "unit": "B/op",
            "extra": "34452 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 422,
            "unit": "allocs/op",
            "extra": "34452 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 35601,
            "unit": "ns/op\t  65.76 MB/s\t   12867 B/op\t     422 allocs/op",
            "extra": "34022 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 35601,
            "unit": "ns/op",
            "extra": "34022 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 65.76,
            "unit": "MB/s",
            "extra": "34022 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 12867,
            "unit": "B/op",
            "extra": "34022 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 422,
            "unit": "allocs/op",
            "extra": "34022 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 35240,
            "unit": "ns/op\t  66.43 MB/s\t   12867 B/op\t     422 allocs/op",
            "extra": "32828 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 35240,
            "unit": "ns/op",
            "extra": "32828 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 66.43,
            "unit": "MB/s",
            "extra": "32828 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 12867,
            "unit": "B/op",
            "extra": "32828 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 422,
            "unit": "allocs/op",
            "extra": "32828 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 36014,
            "unit": "ns/op\t  65.00 MB/s\t   12867 B/op\t     422 allocs/op",
            "extra": "34389 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 36014,
            "unit": "ns/op",
            "extra": "34389 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 65,
            "unit": "MB/s",
            "extra": "34389 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 12867,
            "unit": "B/op",
            "extra": "34389 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 422,
            "unit": "allocs/op",
            "extra": "34389 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 515078,
            "unit": "ns/op\t  38.83 MB/s\t   65520 B/op\t      11 allocs/op",
            "extra": "2340 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 515078,
            "unit": "ns/op",
            "extra": "2340 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 38.83,
            "unit": "MB/s",
            "extra": "2340 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 65520,
            "unit": "B/op",
            "extra": "2340 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "2340 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 494692,
            "unit": "ns/op\t  40.43 MB/s\t   65520 B/op\t      11 allocs/op",
            "extra": "2186 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 494692,
            "unit": "ns/op",
            "extra": "2186 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 40.43,
            "unit": "MB/s",
            "extra": "2186 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 65520,
            "unit": "B/op",
            "extra": "2186 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "2186 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 502035,
            "unit": "ns/op\t  39.84 MB/s\t   65520 B/op\t      11 allocs/op",
            "extra": "2562 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 502035,
            "unit": "ns/op",
            "extra": "2562 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 39.84,
            "unit": "MB/s",
            "extra": "2562 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 65520,
            "unit": "B/op",
            "extra": "2562 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "2562 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 497905,
            "unit": "ns/op\t  40.17 MB/s\t   65520 B/op\t      11 allocs/op",
            "extra": "2420 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 497905,
            "unit": "ns/op",
            "extra": "2420 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 40.17,
            "unit": "MB/s",
            "extra": "2420 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 65520,
            "unit": "B/op",
            "extra": "2420 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "2420 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 513366,
            "unit": "ns/op\t  38.96 MB/s\t   65520 B/op\t      11 allocs/op",
            "extra": "2413 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 513366,
            "unit": "ns/op",
            "extra": "2413 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 38.96,
            "unit": "MB/s",
            "extra": "2413 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 65520,
            "unit": "B/op",
            "extra": "2413 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "2413 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1380016,
            "unit": "ns/op\t  21.74 MB/s\t  630530 B/op\t   15011 allocs/op",
            "extra": "856 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1380016,
            "unit": "ns/op",
            "extra": "856 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 21.74,
            "unit": "MB/s",
            "extra": "856 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 630530,
            "unit": "B/op",
            "extra": "856 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 15011,
            "unit": "allocs/op",
            "extra": "856 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1363445,
            "unit": "ns/op\t  22.00 MB/s\t  630530 B/op\t   15011 allocs/op",
            "extra": "879 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1363445,
            "unit": "ns/op",
            "extra": "879 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 22,
            "unit": "MB/s",
            "extra": "879 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 630530,
            "unit": "B/op",
            "extra": "879 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 15011,
            "unit": "allocs/op",
            "extra": "879 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1499165,
            "unit": "ns/op\t  20.01 MB/s\t  630530 B/op\t   15011 allocs/op",
            "extra": "837 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1499165,
            "unit": "ns/op",
            "extra": "837 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 20.01,
            "unit": "MB/s",
            "extra": "837 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 630530,
            "unit": "B/op",
            "extra": "837 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 15011,
            "unit": "allocs/op",
            "extra": "837 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1521880,
            "unit": "ns/op\t  19.71 MB/s\t  630530 B/op\t   15011 allocs/op",
            "extra": "704 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1521880,
            "unit": "ns/op",
            "extra": "704 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 19.71,
            "unit": "MB/s",
            "extra": "704 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 630530,
            "unit": "B/op",
            "extra": "704 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 15011,
            "unit": "allocs/op",
            "extra": "704 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1421853,
            "unit": "ns/op\t  21.10 MB/s\t  630530 B/op\t   15011 allocs/op",
            "extra": "836 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1421853,
            "unit": "ns/op",
            "extra": "836 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 21.1,
            "unit": "MB/s",
            "extra": "836 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 630530,
            "unit": "B/op",
            "extra": "836 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 15011,
            "unit": "allocs/op",
            "extra": "836 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 118510,
            "unit": "ns/op\t  33.85 MB/s\t   90121 B/op\t    1812 allocs/op",
            "extra": "8470 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 118510,
            "unit": "ns/op",
            "extra": "8470 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 33.85,
            "unit": "MB/s",
            "extra": "8470 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 90121,
            "unit": "B/op",
            "extra": "8470 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "8470 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 117090,
            "unit": "ns/op\t  34.26 MB/s\t   90121 B/op\t    1812 allocs/op",
            "extra": "10274 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 117090,
            "unit": "ns/op",
            "extra": "10274 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 34.26,
            "unit": "MB/s",
            "extra": "10274 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 90121,
            "unit": "B/op",
            "extra": "10274 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "10274 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 116930,
            "unit": "ns/op\t  34.30 MB/s\t   90121 B/op\t    1812 allocs/op",
            "extra": "9697 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 116930,
            "unit": "ns/op",
            "extra": "9697 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 34.3,
            "unit": "MB/s",
            "extra": "9697 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 90121,
            "unit": "B/op",
            "extra": "9697 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "9697 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 115600,
            "unit": "ns/op\t  34.70 MB/s\t   90121 B/op\t    1812 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 115600,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 34.7,
            "unit": "MB/s",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 90121,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 114589,
            "unit": "ns/op\t  35.00 MB/s\t   90121 B/op\t    1812 allocs/op",
            "extra": "9390 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 114589,
            "unit": "ns/op",
            "extra": "9390 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 35,
            "unit": "MB/s",
            "extra": "9390 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 90121,
            "unit": "B/op",
            "extra": "9390 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "9390 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 73269,
            "unit": "ns/op\t  31.95 MB/s\t   54065 B/op\t     781 allocs/op",
            "extra": "16548 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 73269,
            "unit": "ns/op",
            "extra": "16548 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 31.95,
            "unit": "MB/s",
            "extra": "16548 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 54065,
            "unit": "B/op",
            "extra": "16548 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 781,
            "unit": "allocs/op",
            "extra": "16548 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 72109,
            "unit": "ns/op\t  32.46 MB/s\t   54065 B/op\t     781 allocs/op",
            "extra": "16495 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 72109,
            "unit": "ns/op",
            "extra": "16495 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 32.46,
            "unit": "MB/s",
            "extra": "16495 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 54065,
            "unit": "B/op",
            "extra": "16495 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 781,
            "unit": "allocs/op",
            "extra": "16495 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 72074,
            "unit": "ns/op\t  32.48 MB/s\t   54065 B/op\t     781 allocs/op",
            "extra": "16297 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 72074,
            "unit": "ns/op",
            "extra": "16297 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 32.48,
            "unit": "MB/s",
            "extra": "16297 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 54065,
            "unit": "B/op",
            "extra": "16297 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 781,
            "unit": "allocs/op",
            "extra": "16297 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 71674,
            "unit": "ns/op\t  32.66 MB/s\t   54065 B/op\t     781 allocs/op",
            "extra": "16532 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 71674,
            "unit": "ns/op",
            "extra": "16532 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 32.66,
            "unit": "MB/s",
            "extra": "16532 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 54065,
            "unit": "B/op",
            "extra": "16532 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 781,
            "unit": "allocs/op",
            "extra": "16532 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 73519,
            "unit": "ns/op\t  31.84 MB/s\t   54065 B/op\t     781 allocs/op",
            "extra": "16447 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 73519,
            "unit": "ns/op",
            "extra": "16447 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 31.84,
            "unit": "MB/s",
            "extra": "16447 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 54065,
            "unit": "B/op",
            "extra": "16447 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 781,
            "unit": "allocs/op",
            "extra": "16447 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1876093,
            "unit": "ns/op\t  10.66 MB/s\t 2308002 B/op\t   20030 allocs/op",
            "extra": "646 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1876093,
            "unit": "ns/op",
            "extra": "646 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 10.66,
            "unit": "MB/s",
            "extra": "646 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2308002,
            "unit": "B/op",
            "extra": "646 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20030,
            "unit": "allocs/op",
            "extra": "646 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1832152,
            "unit": "ns/op\t  10.92 MB/s\t 2308000 B/op\t   20030 allocs/op",
            "extra": "633 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1832152,
            "unit": "ns/op",
            "extra": "633 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 10.92,
            "unit": "MB/s",
            "extra": "633 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2308000,
            "unit": "B/op",
            "extra": "633 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20030,
            "unit": "allocs/op",
            "extra": "633 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 2062527,
            "unit": "ns/op\t   9.70 MB/s\t 2308004 B/op\t   20030 allocs/op",
            "extra": "628 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 2062527,
            "unit": "ns/op",
            "extra": "628 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 9.7,
            "unit": "MB/s",
            "extra": "628 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2308004,
            "unit": "B/op",
            "extra": "628 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20030,
            "unit": "allocs/op",
            "extra": "628 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 2099974,
            "unit": "ns/op\t   9.52 MB/s\t 2308002 B/op\t   20030 allocs/op",
            "extra": "590 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 2099974,
            "unit": "ns/op",
            "extra": "590 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 9.52,
            "unit": "MB/s",
            "extra": "590 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2308002,
            "unit": "B/op",
            "extra": "590 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20030,
            "unit": "allocs/op",
            "extra": "590 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1837745,
            "unit": "ns/op\t  10.88 MB/s\t 2308000 B/op\t   20030 allocs/op",
            "extra": "642 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1837745,
            "unit": "ns/op",
            "extra": "642 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 10.88,
            "unit": "MB/s",
            "extra": "642 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2308000,
            "unit": "B/op",
            "extra": "642 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20030,
            "unit": "allocs/op",
            "extra": "642 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1693875,
            "unit": "ns/op\t  17.71 MB/s\t 2610264 B/op\t   20026 allocs/op",
            "extra": "667 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1693875,
            "unit": "ns/op",
            "extra": "667 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 17.71,
            "unit": "MB/s",
            "extra": "667 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2610264,
            "unit": "B/op",
            "extra": "667 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20026,
            "unit": "allocs/op",
            "extra": "667 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1588056,
            "unit": "ns/op\t  18.89 MB/s\t 2610263 B/op\t   20026 allocs/op",
            "extra": "751 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1588056,
            "unit": "ns/op",
            "extra": "751 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 18.89,
            "unit": "MB/s",
            "extra": "751 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2610263,
            "unit": "B/op",
            "extra": "751 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20026,
            "unit": "allocs/op",
            "extra": "751 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1586197,
            "unit": "ns/op\t  18.91 MB/s\t 2610265 B/op\t   20026 allocs/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1586197,
            "unit": "ns/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 18.91,
            "unit": "MB/s",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2610265,
            "unit": "B/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20026,
            "unit": "allocs/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1589511,
            "unit": "ns/op\t  18.87 MB/s\t 2610263 B/op\t   20026 allocs/op",
            "extra": "723 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1589511,
            "unit": "ns/op",
            "extra": "723 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 18.87,
            "unit": "MB/s",
            "extra": "723 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2610263,
            "unit": "B/op",
            "extra": "723 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20026,
            "unit": "allocs/op",
            "extra": "723 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1605267,
            "unit": "ns/op\t  18.69 MB/s\t 2610264 B/op\t   20026 allocs/op",
            "extra": "730 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1605267,
            "unit": "ns/op",
            "extra": "730 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 18.69,
            "unit": "MB/s",
            "extra": "730 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2610264,
            "unit": "B/op",
            "extra": "730 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20026,
            "unit": "allocs/op",
            "extra": "730 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 25620,
            "unit": "ns/op\t 105.04 MB/s\t    8177 B/op\t       8 allocs/op",
            "extra": "46813 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 25620,
            "unit": "ns/op",
            "extra": "46813 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 105.04,
            "unit": "MB/s",
            "extra": "46813 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8177,
            "unit": "B/op",
            "extra": "46813 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "46813 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 25554,
            "unit": "ns/op\t 105.31 MB/s\t    8177 B/op\t       8 allocs/op",
            "extra": "47396 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 25554,
            "unit": "ns/op",
            "extra": "47396 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 105.31,
            "unit": "MB/s",
            "extra": "47396 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8177,
            "unit": "B/op",
            "extra": "47396 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "47396 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 25810,
            "unit": "ns/op\t 104.26 MB/s\t    8177 B/op\t       8 allocs/op",
            "extra": "46182 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 25810,
            "unit": "ns/op",
            "extra": "46182 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 104.26,
            "unit": "MB/s",
            "extra": "46182 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8177,
            "unit": "B/op",
            "extra": "46182 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "46182 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 25559,
            "unit": "ns/op\t 105.29 MB/s\t    8177 B/op\t       8 allocs/op",
            "extra": "46110 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 25559,
            "unit": "ns/op",
            "extra": "46110 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 105.29,
            "unit": "MB/s",
            "extra": "46110 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8177,
            "unit": "B/op",
            "extra": "46110 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "46110 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 25867,
            "unit": "ns/op\t 104.03 MB/s\t    8177 B/op\t       8 allocs/op",
            "extra": "45375 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 25867,
            "unit": "ns/op",
            "extra": "45375 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 104.03,
            "unit": "MB/s",
            "extra": "45375 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8177,
            "unit": "B/op",
            "extra": "45375 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "45375 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 22143,
            "unit": "ns/op\t  58.75 MB/s\t    8257 B/op\t     264 allocs/op",
            "extra": "53476 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 22143,
            "unit": "ns/op",
            "extra": "53476 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 58.75,
            "unit": "MB/s",
            "extra": "53476 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8257,
            "unit": "B/op",
            "extra": "53476 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 264,
            "unit": "allocs/op",
            "extra": "53476 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 22819,
            "unit": "ns/op\t  57.01 MB/s\t    8257 B/op\t     264 allocs/op",
            "extra": "52346 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 22819,
            "unit": "ns/op",
            "extra": "52346 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 57.01,
            "unit": "MB/s",
            "extra": "52346 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8257,
            "unit": "B/op",
            "extra": "52346 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 264,
            "unit": "allocs/op",
            "extra": "52346 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 22876,
            "unit": "ns/op\t  56.87 MB/s\t    8257 B/op\t     264 allocs/op",
            "extra": "53650 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 22876,
            "unit": "ns/op",
            "extra": "53650 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 56.87,
            "unit": "MB/s",
            "extra": "53650 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8257,
            "unit": "B/op",
            "extra": "53650 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 264,
            "unit": "allocs/op",
            "extra": "53650 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 22216,
            "unit": "ns/op\t  58.56 MB/s\t    8257 B/op\t     264 allocs/op",
            "extra": "53290 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 22216,
            "unit": "ns/op",
            "extra": "53290 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 58.56,
            "unit": "MB/s",
            "extra": "53290 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8257,
            "unit": "B/op",
            "extra": "53290 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 264,
            "unit": "allocs/op",
            "extra": "53290 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 22860,
            "unit": "ns/op\t  56.91 MB/s\t    8257 B/op\t     264 allocs/op",
            "extra": "51684 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 22860,
            "unit": "ns/op",
            "extra": "51684 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 56.91,
            "unit": "MB/s",
            "extra": "51684 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8257,
            "unit": "B/op",
            "extra": "51684 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 264,
            "unit": "allocs/op",
            "extra": "51684 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 345605,
            "unit": "ns/op\t  28.94 MB/s\t   32755 B/op\t      10 allocs/op",
            "extra": "3315 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 345605,
            "unit": "ns/op",
            "extra": "3315 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 28.94,
            "unit": "MB/s",
            "extra": "3315 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 32755,
            "unit": "B/op",
            "extra": "3315 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "3315 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 375983,
            "unit": "ns/op\t  26.60 MB/s\t   32756 B/op\t      10 allocs/op",
            "extra": "3136 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 375983,
            "unit": "ns/op",
            "extra": "3136 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 26.6,
            "unit": "MB/s",
            "extra": "3136 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 32756,
            "unit": "B/op",
            "extra": "3136 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "3136 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 354184,
            "unit": "ns/op\t  28.24 MB/s\t   32755 B/op\t      10 allocs/op",
            "extra": "3097 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 354184,
            "unit": "ns/op",
            "extra": "3097 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 28.24,
            "unit": "MB/s",
            "extra": "3097 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 32755,
            "unit": "B/op",
            "extra": "3097 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "3097 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 364750,
            "unit": "ns/op\t  27.42 MB/s\t   32755 B/op\t      10 allocs/op",
            "extra": "3207 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 364750,
            "unit": "ns/op",
            "extra": "3207 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 27.42,
            "unit": "MB/s",
            "extra": "3207 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 32755,
            "unit": "B/op",
            "extra": "3207 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "3207 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 345372,
            "unit": "ns/op\t  28.96 MB/s\t   32756 B/op\t      10 allocs/op",
            "extra": "3098 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 345372,
            "unit": "ns/op",
            "extra": "3098 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 28.96,
            "unit": "MB/s",
            "extra": "3098 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 32756,
            "unit": "B/op",
            "extra": "3098 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "3098 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1267847,
            "unit": "ns/op\t  11.83 MB/s\t  512828 B/op\t   10010 allocs/op",
            "extra": "930 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1267847,
            "unit": "ns/op",
            "extra": "930 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 11.83,
            "unit": "MB/s",
            "extra": "930 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 512828,
            "unit": "B/op",
            "extra": "930 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10010,
            "unit": "allocs/op",
            "extra": "930 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1278415,
            "unit": "ns/op\t  11.73 MB/s\t  512828 B/op\t   10010 allocs/op",
            "extra": "1002 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1278415,
            "unit": "ns/op",
            "extra": "1002 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 11.73,
            "unit": "MB/s",
            "extra": "1002 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 512828,
            "unit": "B/op",
            "extra": "1002 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10010,
            "unit": "allocs/op",
            "extra": "1002 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1213367,
            "unit": "ns/op\t  12.36 MB/s\t  512828 B/op\t   10010 allocs/op",
            "extra": "974 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1213367,
            "unit": "ns/op",
            "extra": "974 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 12.36,
            "unit": "MB/s",
            "extra": "974 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 512828,
            "unit": "B/op",
            "extra": "974 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10010,
            "unit": "allocs/op",
            "extra": "974 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1270115,
            "unit": "ns/op\t  11.81 MB/s\t  512828 B/op\t   10010 allocs/op",
            "extra": "922 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1270115,
            "unit": "ns/op",
            "extra": "922 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 11.81,
            "unit": "MB/s",
            "extra": "922 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 512828,
            "unit": "B/op",
            "extra": "922 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10010,
            "unit": "allocs/op",
            "extra": "922 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1281504,
            "unit": "ns/op\t  11.71 MB/s\t  512825 B/op\t   10010 allocs/op",
            "extra": "940 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1281504,
            "unit": "ns/op",
            "extra": "940 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 11.71,
            "unit": "MB/s",
            "extra": "940 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 512825,
            "unit": "B/op",
            "extra": "940 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10010,
            "unit": "allocs/op",
            "extra": "940 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 55719,
            "unit": "ns/op\t  48.30 MB/s\t   24653 B/op\t     771 allocs/op",
            "extra": "21535 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 55719,
            "unit": "ns/op",
            "extra": "21535 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 48.3,
            "unit": "MB/s",
            "extra": "21535 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 24653,
            "unit": "B/op",
            "extra": "21535 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 771,
            "unit": "allocs/op",
            "extra": "21535 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 55792,
            "unit": "ns/op\t  48.23 MB/s\t   24653 B/op\t     771 allocs/op",
            "extra": "20946 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 55792,
            "unit": "ns/op",
            "extra": "20946 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 48.23,
            "unit": "MB/s",
            "extra": "20946 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 24653,
            "unit": "B/op",
            "extra": "20946 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 771,
            "unit": "allocs/op",
            "extra": "20946 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 54751,
            "unit": "ns/op\t  49.15 MB/s\t   24653 B/op\t     771 allocs/op",
            "extra": "21736 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 54751,
            "unit": "ns/op",
            "extra": "21736 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 49.15,
            "unit": "MB/s",
            "extra": "21736 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 24653,
            "unit": "B/op",
            "extra": "21736 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 771,
            "unit": "allocs/op",
            "extra": "21736 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 54322,
            "unit": "ns/op\t  49.54 MB/s\t   24653 B/op\t     771 allocs/op",
            "extra": "21780 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 54322,
            "unit": "ns/op",
            "extra": "21780 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 49.54,
            "unit": "MB/s",
            "extra": "21780 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 24653,
            "unit": "B/op",
            "extra": "21780 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 771,
            "unit": "allocs/op",
            "extra": "21780 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 54111,
            "unit": "ns/op\t  49.73 MB/s\t   24653 B/op\t     771 allocs/op",
            "extra": "21891 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 54111,
            "unit": "ns/op",
            "extra": "21891 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 49.73,
            "unit": "MB/s",
            "extra": "21891 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 24653,
            "unit": "B/op",
            "extra": "21891 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 771,
            "unit": "allocs/op",
            "extra": "21891 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 31997,
            "unit": "ns/op\t  40.66 MB/s\t   19300 B/op\t     252 allocs/op",
            "extra": "37215 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 31997,
            "unit": "ns/op",
            "extra": "37215 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 40.66,
            "unit": "MB/s",
            "extra": "37215 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 19300,
            "unit": "B/op",
            "extra": "37215 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 252,
            "unit": "allocs/op",
            "extra": "37215 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 32034,
            "unit": "ns/op\t  40.61 MB/s\t   19300 B/op\t     252 allocs/op",
            "extra": "37142 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 32034,
            "unit": "ns/op",
            "extra": "37142 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 40.61,
            "unit": "MB/s",
            "extra": "37142 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 19300,
            "unit": "B/op",
            "extra": "37142 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 252,
            "unit": "allocs/op",
            "extra": "37142 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 32119,
            "unit": "ns/op\t  40.51 MB/s\t   19300 B/op\t     252 allocs/op",
            "extra": "37260 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 32119,
            "unit": "ns/op",
            "extra": "37260 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 40.51,
            "unit": "MB/s",
            "extra": "37260 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 19300,
            "unit": "B/op",
            "extra": "37260 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 252,
            "unit": "allocs/op",
            "extra": "37260 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 32020,
            "unit": "ns/op\t  40.63 MB/s\t   19300 B/op\t     252 allocs/op",
            "extra": "37396 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 32020,
            "unit": "ns/op",
            "extra": "37396 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 40.63,
            "unit": "MB/s",
            "extra": "37396 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 19300,
            "unit": "B/op",
            "extra": "37396 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 252,
            "unit": "allocs/op",
            "extra": "37396 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 31873,
            "unit": "ns/op\t  40.82 MB/s\t   19300 B/op\t     252 allocs/op",
            "extra": "37942 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 31873,
            "unit": "ns/op",
            "extra": "37942 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 40.82,
            "unit": "MB/s",
            "extra": "37942 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 19300,
            "unit": "B/op",
            "extra": "37942 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 252,
            "unit": "allocs/op",
            "extra": "37942 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1085073,
            "unit": "ns/op\t   9.22 MB/s\t  400103 B/op\t   20001 allocs/op",
            "extra": "1083 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1085073,
            "unit": "ns/op",
            "extra": "1083 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.22,
            "unit": "MB/s",
            "extra": "1083 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 400103,
            "unit": "B/op",
            "extra": "1083 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 20001,
            "unit": "allocs/op",
            "extra": "1083 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1086264,
            "unit": "ns/op\t   9.21 MB/s\t  400103 B/op\t   20001 allocs/op",
            "extra": "1194 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1086264,
            "unit": "ns/op",
            "extra": "1194 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.21,
            "unit": "MB/s",
            "extra": "1194 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 400103,
            "unit": "B/op",
            "extra": "1194 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 20001,
            "unit": "allocs/op",
            "extra": "1194 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1094375,
            "unit": "ns/op\t   9.14 MB/s\t  400103 B/op\t   20001 allocs/op",
            "extra": "1090 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1094375,
            "unit": "ns/op",
            "extra": "1090 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.14,
            "unit": "MB/s",
            "extra": "1090 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 400103,
            "unit": "B/op",
            "extra": "1090 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 20001,
            "unit": "allocs/op",
            "extra": "1090 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1087557,
            "unit": "ns/op\t   9.20 MB/s\t  400102 B/op\t   20001 allocs/op",
            "extra": "1131 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1087557,
            "unit": "ns/op",
            "extra": "1131 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.2,
            "unit": "MB/s",
            "extra": "1131 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 400102,
            "unit": "B/op",
            "extra": "1131 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 20001,
            "unit": "allocs/op",
            "extra": "1131 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1081991,
            "unit": "ns/op\t   9.24 MB/s\t  400101 B/op\t   20001 allocs/op",
            "extra": "1089 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1081991,
            "unit": "ns/op",
            "extra": "1089 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.24,
            "unit": "MB/s",
            "extra": "1089 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 400101,
            "unit": "B/op",
            "extra": "1089 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 20001,
            "unit": "allocs/op",
            "extra": "1089 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1627876,
            "unit": "ns/op\t   9.22 MB/s\t 1720226 B/op\t   15002 allocs/op",
            "extra": "708 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1627876,
            "unit": "ns/op",
            "extra": "708 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.22,
            "unit": "MB/s",
            "extra": "708 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 1720226,
            "unit": "B/op",
            "extra": "708 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 15002,
            "unit": "allocs/op",
            "extra": "708 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1618105,
            "unit": "ns/op\t   9.27 MB/s\t 1720222 B/op\t   15002 allocs/op",
            "extra": "747 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1618105,
            "unit": "ns/op",
            "extra": "747 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.27,
            "unit": "MB/s",
            "extra": "747 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 1720222,
            "unit": "B/op",
            "extra": "747 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 15002,
            "unit": "allocs/op",
            "extra": "747 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1632049,
            "unit": "ns/op\t   9.19 MB/s\t 1720228 B/op\t   15002 allocs/op",
            "extra": "706 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1632049,
            "unit": "ns/op",
            "extra": "706 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.19,
            "unit": "MB/s",
            "extra": "706 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 1720228,
            "unit": "B/op",
            "extra": "706 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 15002,
            "unit": "allocs/op",
            "extra": "706 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1615555,
            "unit": "ns/op\t   9.29 MB/s\t 1720228 B/op\t   15002 allocs/op",
            "extra": "747 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1615555,
            "unit": "ns/op",
            "extra": "747 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.29,
            "unit": "MB/s",
            "extra": "747 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 1720228,
            "unit": "B/op",
            "extra": "747 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 15002,
            "unit": "allocs/op",
            "extra": "747 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1616033,
            "unit": "ns/op\t   9.28 MB/s\t 1720225 B/op\t   15002 allocs/op",
            "extra": "727 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1616033,
            "unit": "ns/op",
            "extra": "727 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.28,
            "unit": "MB/s",
            "extra": "727 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 1720225,
            "unit": "B/op",
            "extra": "727 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 15002,
            "unit": "allocs/op",
            "extra": "727 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3204,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "365018 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3204,
            "unit": "ns/op",
            "extra": "365018 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "365018 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "365018 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3210,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "374811 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3210,
            "unit": "ns/op",
            "extra": "374811 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "374811 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "374811 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3211,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "374377 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3211,
            "unit": "ns/op",
            "extra": "374377 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "374377 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "374377 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3207,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "373866 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3207,
            "unit": "ns/op",
            "extra": "373866 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "373866 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "373866 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3200,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "374274 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3200,
            "unit": "ns/op",
            "extra": "374274 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "374274 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "374274 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 356.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3672920 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 356.8,
            "unit": "ns/op",
            "extra": "3672920 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3672920 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3672920 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 348.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3552712 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 348.8,
            "unit": "ns/op",
            "extra": "3552712 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3552712 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3552712 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 349.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3673339 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 349.4,
            "unit": "ns/op",
            "extra": "3673339 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3673339 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3673339 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 357.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3670544 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 357.5,
            "unit": "ns/op",
            "extra": "3670544 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3670544 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3670544 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 355.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3407960 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 355.6,
            "unit": "ns/op",
            "extra": "3407960 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3407960 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3407960 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2620,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "480175 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2620,
            "unit": "ns/op",
            "extra": "480175 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "480175 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "480175 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2565,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "457401 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2565,
            "unit": "ns/op",
            "extra": "457401 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "457401 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "457401 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2595,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "450328 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2595,
            "unit": "ns/op",
            "extra": "450328 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "450328 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "450328 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2557,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "450826 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2557,
            "unit": "ns/op",
            "extra": "450826 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "450826 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "450826 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2642,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "451422 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2642,
            "unit": "ns/op",
            "extra": "451422 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "451422 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "451422 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 1026,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1215031 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 1026,
            "unit": "ns/op",
            "extra": "1215031 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1215031 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1215031 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 967.1,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1219318 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 967.1,
            "unit": "ns/op",
            "extra": "1219318 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1219318 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1219318 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 966.4,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1201112 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 966.4,
            "unit": "ns/op",
            "extra": "1201112 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1201112 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1201112 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 982.9,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1246306 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 982.9,
            "unit": "ns/op",
            "extra": "1246306 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1246306 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1246306 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 987.3,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1249111 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 987.3,
            "unit": "ns/op",
            "extra": "1249111 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1249111 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1249111 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 188.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8178872 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 188.5,
            "unit": "ns/op",
            "extra": "8178872 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8178872 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8178872 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 190,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8275530 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 190,
            "unit": "ns/op",
            "extra": "8275530 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8275530 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8275530 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 208.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9463428 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 208.5,
            "unit": "ns/op",
            "extra": "9463428 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9463428 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9463428 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 198.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6511638 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 198.3,
            "unit": "ns/op",
            "extra": "6511638 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6511638 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6511638 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 193.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9178334 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 193.7,
            "unit": "ns/op",
            "extra": "9178334 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9178334 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9178334 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 16.59,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "69929498 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 16.59,
            "unit": "ns/op",
            "extra": "69929498 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "69929498 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "69929498 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 16.59,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "71156962 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 16.59,
            "unit": "ns/op",
            "extra": "71156962 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "71156962 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "71156962 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 17.02,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "72616530 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 17.02,
            "unit": "ns/op",
            "extra": "72616530 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "72616530 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "72616530 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 16.56,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "72462856 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 16.56,
            "unit": "ns/op",
            "extra": "72462856 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "72462856 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "72462856 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 16.57,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "71113664 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 16.57,
            "unit": "ns/op",
            "extra": "71113664 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "71113664 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "71113664 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.39,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30490021 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.39,
            "unit": "ns/op",
            "extra": "30490021 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30490021 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30490021 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.32,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30499986 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.32,
            "unit": "ns/op",
            "extra": "30499986 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30499986 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30499986 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.63,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30598278 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.63,
            "unit": "ns/op",
            "extra": "30598278 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30598278 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30598278 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.52,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30638672 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.52,
            "unit": "ns/op",
            "extra": "30638672 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30638672 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30638672 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30436378 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.4,
            "unit": "ns/op",
            "extra": "30436378 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30436378 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30436378 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 75.24,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15970220 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 75.24,
            "unit": "ns/op",
            "extra": "15970220 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15970220 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15970220 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 75.11,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15894739 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 75.11,
            "unit": "ns/op",
            "extra": "15894739 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15894739 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15894739 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 75.32,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15936662 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 75.32,
            "unit": "ns/op",
            "extra": "15936662 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15936662 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15936662 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 75.16,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15801492 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 75.16,
            "unit": "ns/op",
            "extra": "15801492 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15801492 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15801492 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 75.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15952052 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 75.2,
            "unit": "ns/op",
            "extra": "15952052 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15952052 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15952052 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.74,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30153860 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.74,
            "unit": "ns/op",
            "extra": "30153860 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30153860 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30153860 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30099219 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.83,
            "unit": "ns/op",
            "extra": "30099219 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30099219 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30099219 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.81,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30122890 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.81,
            "unit": "ns/op",
            "extra": "30122890 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30122890 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30122890 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30087043 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.8,
            "unit": "ns/op",
            "extra": "30087043 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30087043 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30087043 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.82,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30255032 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.82,
            "unit": "ns/op",
            "extra": "30255032 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30255032 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30255032 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "ziflex@gmail.com",
            "name": "Tim Voronov",
            "username": "ziflex"
          },
          "committer": {
            "email": "ziflex@gmail.com",
            "name": "Tim Voronov",
            "username": "ziflex"
          },
          "distinct": true,
          "id": "d830ffcf66f1e736b6e9b0c72a7b22907bb2c27c",
          "message": "fix: update version constant to reflect the latest release",
          "timestamp": "2026-03-19T15:11:45-04:00",
          "tree_id": "c5aad4858435af518bbcec1ebb321ebde10aa29c",
          "url": "https://github.com/MontFerret/ferret/commit/d830ffcf66f1e736b6e9b0c72a7b22907bb2c27c"
        },
        "date": 1773947713222,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 50444,
            "unit": "ns/op\t  79.51 MB/s\t   11152 B/op\t     932 allocs/op",
            "extra": "23179 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 50444,
            "unit": "ns/op",
            "extra": "23179 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 79.51,
            "unit": "MB/s",
            "extra": "23179 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 11152,
            "unit": "B/op",
            "extra": "23179 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 932,
            "unit": "allocs/op",
            "extra": "23179 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 50861,
            "unit": "ns/op\t  78.86 MB/s\t   11152 B/op\t     932 allocs/op",
            "extra": "23583 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 50861,
            "unit": "ns/op",
            "extra": "23583 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 78.86,
            "unit": "MB/s",
            "extra": "23583 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 11152,
            "unit": "B/op",
            "extra": "23583 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 932,
            "unit": "allocs/op",
            "extra": "23583 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 50718,
            "unit": "ns/op\t  79.08 MB/s\t   11152 B/op\t     932 allocs/op",
            "extra": "23408 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 50718,
            "unit": "ns/op",
            "extra": "23408 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 79.08,
            "unit": "MB/s",
            "extra": "23408 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 11152,
            "unit": "B/op",
            "extra": "23408 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 932,
            "unit": "allocs/op",
            "extra": "23408 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 50923,
            "unit": "ns/op\t  78.77 MB/s\t   11152 B/op\t     932 allocs/op",
            "extra": "23872 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 50923,
            "unit": "ns/op",
            "extra": "23872 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 78.77,
            "unit": "MB/s",
            "extra": "23872 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 11152,
            "unit": "B/op",
            "extra": "23872 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 932,
            "unit": "allocs/op",
            "extra": "23872 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 50985,
            "unit": "ns/op\t  78.67 MB/s\t   11152 B/op\t     932 allocs/op",
            "extra": "23430 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 50985,
            "unit": "ns/op",
            "extra": "23430 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 78.67,
            "unit": "MB/s",
            "extra": "23430 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 11152,
            "unit": "B/op",
            "extra": "23430 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 932,
            "unit": "allocs/op",
            "extra": "23430 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 37548,
            "unit": "ns/op\t  62.35 MB/s\t   12867 B/op\t     422 allocs/op",
            "extra": "32472 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 37548,
            "unit": "ns/op",
            "extra": "32472 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 62.35,
            "unit": "MB/s",
            "extra": "32472 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 12867,
            "unit": "B/op",
            "extra": "32472 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 422,
            "unit": "allocs/op",
            "extra": "32472 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 37063,
            "unit": "ns/op\t  63.16 MB/s\t   12867 B/op\t     422 allocs/op",
            "extra": "32924 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 37063,
            "unit": "ns/op",
            "extra": "32924 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 63.16,
            "unit": "MB/s",
            "extra": "32924 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 12867,
            "unit": "B/op",
            "extra": "32924 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 422,
            "unit": "allocs/op",
            "extra": "32924 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 36962,
            "unit": "ns/op\t  63.33 MB/s\t   12867 B/op\t     422 allocs/op",
            "extra": "32924 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 36962,
            "unit": "ns/op",
            "extra": "32924 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 63.33,
            "unit": "MB/s",
            "extra": "32924 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 12867,
            "unit": "B/op",
            "extra": "32924 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 422,
            "unit": "allocs/op",
            "extra": "32924 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 37927,
            "unit": "ns/op\t  61.72 MB/s\t   12867 B/op\t     422 allocs/op",
            "extra": "31384 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 37927,
            "unit": "ns/op",
            "extra": "31384 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 61.72,
            "unit": "MB/s",
            "extra": "31384 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 12867,
            "unit": "B/op",
            "extra": "31384 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 422,
            "unit": "allocs/op",
            "extra": "31384 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 38144,
            "unit": "ns/op\t  61.37 MB/s\t   12867 B/op\t     422 allocs/op",
            "extra": "31070 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 38144,
            "unit": "ns/op",
            "extra": "31070 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 61.37,
            "unit": "MB/s",
            "extra": "31070 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 12867,
            "unit": "B/op",
            "extra": "31070 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 422,
            "unit": "allocs/op",
            "extra": "31070 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 531284,
            "unit": "ns/op\t  37.65 MB/s\t   65520 B/op\t      11 allocs/op",
            "extra": "2185 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 531284,
            "unit": "ns/op",
            "extra": "2185 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 37.65,
            "unit": "MB/s",
            "extra": "2185 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 65520,
            "unit": "B/op",
            "extra": "2185 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "2185 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 508818,
            "unit": "ns/op\t  39.31 MB/s\t   65520 B/op\t      11 allocs/op",
            "extra": "2362 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 508818,
            "unit": "ns/op",
            "extra": "2362 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 39.31,
            "unit": "MB/s",
            "extra": "2362 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 65520,
            "unit": "B/op",
            "extra": "2362 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "2362 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 526055,
            "unit": "ns/op\t  38.02 MB/s\t   65520 B/op\t      11 allocs/op",
            "extra": "2274 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 526055,
            "unit": "ns/op",
            "extra": "2274 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 38.02,
            "unit": "MB/s",
            "extra": "2274 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 65520,
            "unit": "B/op",
            "extra": "2274 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "2274 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 515783,
            "unit": "ns/op\t  38.78 MB/s\t   65520 B/op\t      11 allocs/op",
            "extra": "2382 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 515783,
            "unit": "ns/op",
            "extra": "2382 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 38.78,
            "unit": "MB/s",
            "extra": "2382 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 65520,
            "unit": "B/op",
            "extra": "2382 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "2382 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 504531,
            "unit": "ns/op\t  39.64 MB/s\t   65520 B/op\t      11 allocs/op",
            "extra": "2128 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 504531,
            "unit": "ns/op",
            "extra": "2128 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 39.64,
            "unit": "MB/s",
            "extra": "2128 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 65520,
            "unit": "B/op",
            "extra": "2128 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "2128 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1549640,
            "unit": "ns/op\t  19.36 MB/s\t  630530 B/op\t   15011 allocs/op",
            "extra": "790 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1549640,
            "unit": "ns/op",
            "extra": "790 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 19.36,
            "unit": "MB/s",
            "extra": "790 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 630530,
            "unit": "B/op",
            "extra": "790 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 15011,
            "unit": "allocs/op",
            "extra": "790 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1544545,
            "unit": "ns/op\t  19.42 MB/s\t  630530 B/op\t   15011 allocs/op",
            "extra": "754 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1544545,
            "unit": "ns/op",
            "extra": "754 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 19.42,
            "unit": "MB/s",
            "extra": "754 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 630530,
            "unit": "B/op",
            "extra": "754 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 15011,
            "unit": "allocs/op",
            "extra": "754 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1503156,
            "unit": "ns/op\t  19.96 MB/s\t  630530 B/op\t   15011 allocs/op",
            "extra": "775 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1503156,
            "unit": "ns/op",
            "extra": "775 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 19.96,
            "unit": "MB/s",
            "extra": "775 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 630530,
            "unit": "B/op",
            "extra": "775 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 15011,
            "unit": "allocs/op",
            "extra": "775 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1513637,
            "unit": "ns/op\t  19.82 MB/s\t  630530 B/op\t   15011 allocs/op",
            "extra": "788 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1513637,
            "unit": "ns/op",
            "extra": "788 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 19.82,
            "unit": "MB/s",
            "extra": "788 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 630530,
            "unit": "B/op",
            "extra": "788 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 15011,
            "unit": "allocs/op",
            "extra": "788 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1454371,
            "unit": "ns/op\t  20.63 MB/s\t  630530 B/op\t   15011 allocs/op",
            "extra": "865 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1454371,
            "unit": "ns/op",
            "extra": "865 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 20.63,
            "unit": "MB/s",
            "extra": "865 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 630530,
            "unit": "B/op",
            "extra": "865 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 15011,
            "unit": "allocs/op",
            "extra": "865 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 119936,
            "unit": "ns/op\t  33.44 MB/s\t   90121 B/op\t    1812 allocs/op",
            "extra": "10243 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 119936,
            "unit": "ns/op",
            "extra": "10243 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 33.44,
            "unit": "MB/s",
            "extra": "10243 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 90121,
            "unit": "B/op",
            "extra": "10243 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "10243 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 118593,
            "unit": "ns/op\t  33.82 MB/s\t   90121 B/op\t    1812 allocs/op",
            "extra": "10074 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 118593,
            "unit": "ns/op",
            "extra": "10074 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 33.82,
            "unit": "MB/s",
            "extra": "10074 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 90121,
            "unit": "B/op",
            "extra": "10074 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "10074 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 119821,
            "unit": "ns/op\t  33.47 MB/s\t   90121 B/op\t    1812 allocs/op",
            "extra": "9316 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 119821,
            "unit": "ns/op",
            "extra": "9316 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 33.47,
            "unit": "MB/s",
            "extra": "9316 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 90121,
            "unit": "B/op",
            "extra": "9316 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "9316 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 118781,
            "unit": "ns/op\t  33.77 MB/s\t   90121 B/op\t    1812 allocs/op",
            "extra": "8720 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 118781,
            "unit": "ns/op",
            "extra": "8720 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 33.77,
            "unit": "MB/s",
            "extra": "8720 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 90121,
            "unit": "B/op",
            "extra": "8720 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "8720 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 119558,
            "unit": "ns/op\t  33.55 MB/s\t   90121 B/op\t    1812 allocs/op",
            "extra": "9487 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 119558,
            "unit": "ns/op",
            "extra": "9487 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 33.55,
            "unit": "MB/s",
            "extra": "9487 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 90121,
            "unit": "B/op",
            "extra": "9487 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "9487 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 75765,
            "unit": "ns/op\t  30.90 MB/s\t   54081 B/op\t     781 allocs/op",
            "extra": "15939 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 75765,
            "unit": "ns/op",
            "extra": "15939 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 30.9,
            "unit": "MB/s",
            "extra": "15939 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 54081,
            "unit": "B/op",
            "extra": "15939 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 781,
            "unit": "allocs/op",
            "extra": "15939 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 75032,
            "unit": "ns/op\t  31.20 MB/s\t   54081 B/op\t     781 allocs/op",
            "extra": "15824 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 75032,
            "unit": "ns/op",
            "extra": "15824 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 31.2,
            "unit": "MB/s",
            "extra": "15824 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 54081,
            "unit": "B/op",
            "extra": "15824 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 781,
            "unit": "allocs/op",
            "extra": "15824 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 75402,
            "unit": "ns/op\t  31.05 MB/s\t   54081 B/op\t     781 allocs/op",
            "extra": "15519 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 75402,
            "unit": "ns/op",
            "extra": "15519 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 31.05,
            "unit": "MB/s",
            "extra": "15519 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 54081,
            "unit": "B/op",
            "extra": "15519 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 781,
            "unit": "allocs/op",
            "extra": "15519 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 73808,
            "unit": "ns/op\t  31.72 MB/s\t   54081 B/op\t     781 allocs/op",
            "extra": "15986 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 73808,
            "unit": "ns/op",
            "extra": "15986 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 31.72,
            "unit": "MB/s",
            "extra": "15986 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 54081,
            "unit": "B/op",
            "extra": "15986 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 781,
            "unit": "allocs/op",
            "extra": "15986 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 75353,
            "unit": "ns/op\t  31.07 MB/s\t   54081 B/op\t     781 allocs/op",
            "extra": "16171 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 75353,
            "unit": "ns/op",
            "extra": "16171 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 31.07,
            "unit": "MB/s",
            "extra": "16171 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 54081,
            "unit": "B/op",
            "extra": "16171 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 781,
            "unit": "allocs/op",
            "extra": "16171 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1918310,
            "unit": "ns/op\t  10.43 MB/s\t 2308002 B/op\t   20030 allocs/op",
            "extra": "622 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1918310,
            "unit": "ns/op",
            "extra": "622 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 10.43,
            "unit": "MB/s",
            "extra": "622 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2308002,
            "unit": "B/op",
            "extra": "622 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20030,
            "unit": "allocs/op",
            "extra": "622 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 2152230,
            "unit": "ns/op\t   9.29 MB/s\t 2308005 B/op\t   20030 allocs/op",
            "extra": "598 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 2152230,
            "unit": "ns/op",
            "extra": "598 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 9.29,
            "unit": "MB/s",
            "extra": "598 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2308005,
            "unit": "B/op",
            "extra": "598 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20030,
            "unit": "allocs/op",
            "extra": "598 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1999799,
            "unit": "ns/op\t  10.00 MB/s\t 2308000 B/op\t   20030 allocs/op",
            "extra": "576 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1999799,
            "unit": "ns/op",
            "extra": "576 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 10,
            "unit": "MB/s",
            "extra": "576 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2308000,
            "unit": "B/op",
            "extra": "576 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20030,
            "unit": "allocs/op",
            "extra": "576 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 2141626,
            "unit": "ns/op\t   9.34 MB/s\t 2308005 B/op\t   20030 allocs/op",
            "extra": "558 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 2141626,
            "unit": "ns/op",
            "extra": "558 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 9.34,
            "unit": "MB/s",
            "extra": "558 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2308005,
            "unit": "B/op",
            "extra": "558 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20030,
            "unit": "allocs/op",
            "extra": "558 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1949729,
            "unit": "ns/op\t  10.26 MB/s\t 2308000 B/op\t   20030 allocs/op",
            "extra": "607 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1949729,
            "unit": "ns/op",
            "extra": "607 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 10.26,
            "unit": "MB/s",
            "extra": "607 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2308000,
            "unit": "B/op",
            "extra": "607 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20030,
            "unit": "allocs/op",
            "extra": "607 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1650036,
            "unit": "ns/op\t  18.18 MB/s\t 2610263 B/op\t   20026 allocs/op",
            "extra": "702 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1650036,
            "unit": "ns/op",
            "extra": "702 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 18.18,
            "unit": "MB/s",
            "extra": "702 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2610263,
            "unit": "B/op",
            "extra": "702 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20026,
            "unit": "allocs/op",
            "extra": "702 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1664291,
            "unit": "ns/op\t  18.03 MB/s\t 2610263 B/op\t   20026 allocs/op",
            "extra": "679 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1664291,
            "unit": "ns/op",
            "extra": "679 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 18.03,
            "unit": "MB/s",
            "extra": "679 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2610263,
            "unit": "B/op",
            "extra": "679 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20026,
            "unit": "allocs/op",
            "extra": "679 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1662603,
            "unit": "ns/op\t  18.04 MB/s\t 2610263 B/op\t   20026 allocs/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1662603,
            "unit": "ns/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 18.04,
            "unit": "MB/s",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2610263,
            "unit": "B/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20026,
            "unit": "allocs/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1719924,
            "unit": "ns/op\t  17.44 MB/s\t 2610264 B/op\t   20026 allocs/op",
            "extra": "700 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1719924,
            "unit": "ns/op",
            "extra": "700 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 17.44,
            "unit": "MB/s",
            "extra": "700 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2610264,
            "unit": "B/op",
            "extra": "700 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20026,
            "unit": "allocs/op",
            "extra": "700 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1686082,
            "unit": "ns/op\t  17.79 MB/s\t 2610263 B/op\t   20026 allocs/op",
            "extra": "702 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1686082,
            "unit": "ns/op",
            "extra": "702 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 17.79,
            "unit": "MB/s",
            "extra": "702 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2610263,
            "unit": "B/op",
            "extra": "702 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20026,
            "unit": "allocs/op",
            "extra": "702 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 27596,
            "unit": "ns/op\t  97.51 MB/s\t    8177 B/op\t       8 allocs/op",
            "extra": "43926 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 27596,
            "unit": "ns/op",
            "extra": "43926 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 97.51,
            "unit": "MB/s",
            "extra": "43926 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8177,
            "unit": "B/op",
            "extra": "43926 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "43926 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 27686,
            "unit": "ns/op\t  97.20 MB/s\t    8177 B/op\t       8 allocs/op",
            "extra": "43232 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 27686,
            "unit": "ns/op",
            "extra": "43232 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 97.2,
            "unit": "MB/s",
            "extra": "43232 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8177,
            "unit": "B/op",
            "extra": "43232 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "43232 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 27405,
            "unit": "ns/op\t  98.20 MB/s\t    8177 B/op\t       8 allocs/op",
            "extra": "42800 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 27405,
            "unit": "ns/op",
            "extra": "42800 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 98.2,
            "unit": "MB/s",
            "extra": "42800 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8177,
            "unit": "B/op",
            "extra": "42800 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "42800 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 27554,
            "unit": "ns/op\t  97.66 MB/s\t    8177 B/op\t       8 allocs/op",
            "extra": "44649 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 27554,
            "unit": "ns/op",
            "extra": "44649 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 97.66,
            "unit": "MB/s",
            "extra": "44649 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8177,
            "unit": "B/op",
            "extra": "44649 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "44649 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 27723,
            "unit": "ns/op\t  97.07 MB/s\t    8177 B/op\t       8 allocs/op",
            "extra": "43843 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 27723,
            "unit": "ns/op",
            "extra": "43843 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 97.07,
            "unit": "MB/s",
            "extra": "43843 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8177,
            "unit": "B/op",
            "extra": "43843 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "43843 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 24049,
            "unit": "ns/op\t  54.10 MB/s\t    8257 B/op\t     264 allocs/op",
            "extra": "50176 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 24049,
            "unit": "ns/op",
            "extra": "50176 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 54.1,
            "unit": "MB/s",
            "extra": "50176 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8257,
            "unit": "B/op",
            "extra": "50176 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 264,
            "unit": "allocs/op",
            "extra": "50176 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 23691,
            "unit": "ns/op\t  54.91 MB/s\t    8257 B/op\t     264 allocs/op",
            "extra": "49851 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 23691,
            "unit": "ns/op",
            "extra": "49851 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 54.91,
            "unit": "MB/s",
            "extra": "49851 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8257,
            "unit": "B/op",
            "extra": "49851 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 264,
            "unit": "allocs/op",
            "extra": "49851 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 23221,
            "unit": "ns/op\t  56.03 MB/s\t    8257 B/op\t     264 allocs/op",
            "extra": "52720 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 23221,
            "unit": "ns/op",
            "extra": "52720 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 56.03,
            "unit": "MB/s",
            "extra": "52720 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8257,
            "unit": "B/op",
            "extra": "52720 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 264,
            "unit": "allocs/op",
            "extra": "52720 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 23504,
            "unit": "ns/op\t  55.35 MB/s\t    8257 B/op\t     264 allocs/op",
            "extra": "52215 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 23504,
            "unit": "ns/op",
            "extra": "52215 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 55.35,
            "unit": "MB/s",
            "extra": "52215 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8257,
            "unit": "B/op",
            "extra": "52215 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 264,
            "unit": "allocs/op",
            "extra": "52215 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 23646,
            "unit": "ns/op\t  55.02 MB/s\t    8257 B/op\t     264 allocs/op",
            "extra": "50740 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 23646,
            "unit": "ns/op",
            "extra": "50740 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 55.02,
            "unit": "MB/s",
            "extra": "50740 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8257,
            "unit": "B/op",
            "extra": "50740 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 264,
            "unit": "allocs/op",
            "extra": "50740 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 349376,
            "unit": "ns/op\t  28.63 MB/s\t   32756 B/op\t      10 allocs/op",
            "extra": "3250 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 349376,
            "unit": "ns/op",
            "extra": "3250 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 28.63,
            "unit": "MB/s",
            "extra": "3250 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 32756,
            "unit": "B/op",
            "extra": "3250 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "3250 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 341920,
            "unit": "ns/op\t  29.25 MB/s\t   32755 B/op\t      10 allocs/op",
            "extra": "3218 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 341920,
            "unit": "ns/op",
            "extra": "3218 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 29.25,
            "unit": "MB/s",
            "extra": "3218 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 32755,
            "unit": "B/op",
            "extra": "3218 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "3218 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 383238,
            "unit": "ns/op\t  26.10 MB/s\t   32755 B/op\t      10 allocs/op",
            "extra": "3073 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 383238,
            "unit": "ns/op",
            "extra": "3073 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 26.1,
            "unit": "MB/s",
            "extra": "3073 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 32755,
            "unit": "B/op",
            "extra": "3073 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "3073 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 363880,
            "unit": "ns/op\t  27.48 MB/s\t   32756 B/op\t      10 allocs/op",
            "extra": "2970 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 363880,
            "unit": "ns/op",
            "extra": "2970 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 27.48,
            "unit": "MB/s",
            "extra": "2970 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 32756,
            "unit": "B/op",
            "extra": "2970 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "2970 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 365361,
            "unit": "ns/op\t  27.37 MB/s\t   32756 B/op\t      10 allocs/op",
            "extra": "3066 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 365361,
            "unit": "ns/op",
            "extra": "3066 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 27.37,
            "unit": "MB/s",
            "extra": "3066 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 32756,
            "unit": "B/op",
            "extra": "3066 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "3066 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1307200,
            "unit": "ns/op\t  11.48 MB/s\t  512824 B/op\t   10010 allocs/op",
            "extra": "918 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1307200,
            "unit": "ns/op",
            "extra": "918 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 11.48,
            "unit": "MB/s",
            "extra": "918 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 512824,
            "unit": "B/op",
            "extra": "918 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10010,
            "unit": "allocs/op",
            "extra": "918 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1419003,
            "unit": "ns/op\t  10.57 MB/s\t  512829 B/op\t   10010 allocs/op",
            "extra": "962 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1419003,
            "unit": "ns/op",
            "extra": "962 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 10.57,
            "unit": "MB/s",
            "extra": "962 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 512829,
            "unit": "B/op",
            "extra": "962 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10010,
            "unit": "allocs/op",
            "extra": "962 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1384201,
            "unit": "ns/op\t  10.84 MB/s\t  512826 B/op\t   10010 allocs/op",
            "extra": "872 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1384201,
            "unit": "ns/op",
            "extra": "872 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 10.84,
            "unit": "MB/s",
            "extra": "872 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 512826,
            "unit": "B/op",
            "extra": "872 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10010,
            "unit": "allocs/op",
            "extra": "872 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1343107,
            "unit": "ns/op\t  11.17 MB/s\t  512825 B/op\t   10010 allocs/op",
            "extra": "862 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1343107,
            "unit": "ns/op",
            "extra": "862 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 11.17,
            "unit": "MB/s",
            "extra": "862 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 512825,
            "unit": "B/op",
            "extra": "862 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10010,
            "unit": "allocs/op",
            "extra": "862 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1317623,
            "unit": "ns/op\t  11.38 MB/s\t  512826 B/op\t   10010 allocs/op",
            "extra": "885 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1317623,
            "unit": "ns/op",
            "extra": "885 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 11.38,
            "unit": "MB/s",
            "extra": "885 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 512826,
            "unit": "B/op",
            "extra": "885 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10010,
            "unit": "allocs/op",
            "extra": "885 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 53116,
            "unit": "ns/op\t  50.66 MB/s\t   24653 B/op\t     771 allocs/op",
            "extra": "22231 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 53116,
            "unit": "ns/op",
            "extra": "22231 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 50.66,
            "unit": "MB/s",
            "extra": "22231 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 24653,
            "unit": "B/op",
            "extra": "22231 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 771,
            "unit": "allocs/op",
            "extra": "22231 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 53289,
            "unit": "ns/op\t  50.50 MB/s\t   24653 B/op\t     771 allocs/op",
            "extra": "22110 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 53289,
            "unit": "ns/op",
            "extra": "22110 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 50.5,
            "unit": "MB/s",
            "extra": "22110 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 24653,
            "unit": "B/op",
            "extra": "22110 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 771,
            "unit": "allocs/op",
            "extra": "22110 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 54409,
            "unit": "ns/op\t  49.46 MB/s\t   24653 B/op\t     771 allocs/op",
            "extra": "22335 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 54409,
            "unit": "ns/op",
            "extra": "22335 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 49.46,
            "unit": "MB/s",
            "extra": "22335 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 24653,
            "unit": "B/op",
            "extra": "22335 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 771,
            "unit": "allocs/op",
            "extra": "22335 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 53622,
            "unit": "ns/op\t  50.18 MB/s\t   24653 B/op\t     771 allocs/op",
            "extra": "22134 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 53622,
            "unit": "ns/op",
            "extra": "22134 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 50.18,
            "unit": "MB/s",
            "extra": "22134 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 24653,
            "unit": "B/op",
            "extra": "22134 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 771,
            "unit": "allocs/op",
            "extra": "22134 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 53943,
            "unit": "ns/op\t  49.89 MB/s\t   24653 B/op\t     771 allocs/op",
            "extra": "22477 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 53943,
            "unit": "ns/op",
            "extra": "22477 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 49.89,
            "unit": "MB/s",
            "extra": "22477 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 24653,
            "unit": "B/op",
            "extra": "22477 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 771,
            "unit": "allocs/op",
            "extra": "22477 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 31802,
            "unit": "ns/op\t  40.91 MB/s\t   19299 B/op\t     252 allocs/op",
            "extra": "37461 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 31802,
            "unit": "ns/op",
            "extra": "37461 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 40.91,
            "unit": "MB/s",
            "extra": "37461 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 19299,
            "unit": "B/op",
            "extra": "37461 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 252,
            "unit": "allocs/op",
            "extra": "37461 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 30974,
            "unit": "ns/op\t  42.00 MB/s\t   19299 B/op\t     252 allocs/op",
            "extra": "39274 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 30974,
            "unit": "ns/op",
            "extra": "39274 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 42,
            "unit": "MB/s",
            "extra": "39274 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 19299,
            "unit": "B/op",
            "extra": "39274 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 252,
            "unit": "allocs/op",
            "extra": "39274 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 30670,
            "unit": "ns/op\t  42.42 MB/s\t   19299 B/op\t     252 allocs/op",
            "extra": "39092 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 30670,
            "unit": "ns/op",
            "extra": "39092 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 42.42,
            "unit": "MB/s",
            "extra": "39092 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 19299,
            "unit": "B/op",
            "extra": "39092 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 252,
            "unit": "allocs/op",
            "extra": "39092 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 30500,
            "unit": "ns/op\t  42.66 MB/s\t   19299 B/op\t     252 allocs/op",
            "extra": "39345 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 30500,
            "unit": "ns/op",
            "extra": "39345 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 42.66,
            "unit": "MB/s",
            "extra": "39345 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 19299,
            "unit": "B/op",
            "extra": "39345 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 252,
            "unit": "allocs/op",
            "extra": "39345 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 30931,
            "unit": "ns/op\t  42.06 MB/s\t   19299 B/op\t     252 allocs/op",
            "extra": "38828 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 30931,
            "unit": "ns/op",
            "extra": "38828 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 42.06,
            "unit": "MB/s",
            "extra": "38828 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 19299,
            "unit": "B/op",
            "extra": "38828 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 252,
            "unit": "allocs/op",
            "extra": "38828 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1158703,
            "unit": "ns/op\t   8.63 MB/s\t  400103 B/op\t   20001 allocs/op",
            "extra": "1040 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1158703,
            "unit": "ns/op",
            "extra": "1040 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 8.63,
            "unit": "MB/s",
            "extra": "1040 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 400103,
            "unit": "B/op",
            "extra": "1040 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 20001,
            "unit": "allocs/op",
            "extra": "1040 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1155893,
            "unit": "ns/op\t   8.65 MB/s\t  400104 B/op\t   20001 allocs/op",
            "extra": "1029 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1155893,
            "unit": "ns/op",
            "extra": "1029 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 8.65,
            "unit": "MB/s",
            "extra": "1029 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 400104,
            "unit": "B/op",
            "extra": "1029 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 20001,
            "unit": "allocs/op",
            "extra": "1029 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1170601,
            "unit": "ns/op\t   8.54 MB/s\t  400104 B/op\t   20001 allocs/op",
            "extra": "1010 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1170601,
            "unit": "ns/op",
            "extra": "1010 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 8.54,
            "unit": "MB/s",
            "extra": "1010 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 400104,
            "unit": "B/op",
            "extra": "1010 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 20001,
            "unit": "allocs/op",
            "extra": "1010 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1134267,
            "unit": "ns/op\t   8.82 MB/s\t  400103 B/op\t   20001 allocs/op",
            "extra": "987 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1134267,
            "unit": "ns/op",
            "extra": "987 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 8.82,
            "unit": "MB/s",
            "extra": "987 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 400103,
            "unit": "B/op",
            "extra": "987 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 20001,
            "unit": "allocs/op",
            "extra": "987 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1130907,
            "unit": "ns/op\t   8.84 MB/s\t  400104 B/op\t   20001 allocs/op",
            "extra": "1096 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1130907,
            "unit": "ns/op",
            "extra": "1096 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 8.84,
            "unit": "MB/s",
            "extra": "1096 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 400104,
            "unit": "B/op",
            "extra": "1096 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 20001,
            "unit": "allocs/op",
            "extra": "1096 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1660301,
            "unit": "ns/op\t   9.04 MB/s\t 1720226 B/op\t   15002 allocs/op",
            "extra": "704 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1660301,
            "unit": "ns/op",
            "extra": "704 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.04,
            "unit": "MB/s",
            "extra": "704 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 1720226,
            "unit": "B/op",
            "extra": "704 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 15002,
            "unit": "allocs/op",
            "extra": "704 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1660202,
            "unit": "ns/op\t   9.04 MB/s\t 1720228 B/op\t   15002 allocs/op",
            "extra": "712 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1660202,
            "unit": "ns/op",
            "extra": "712 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.04,
            "unit": "MB/s",
            "extra": "712 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 1720228,
            "unit": "B/op",
            "extra": "712 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 15002,
            "unit": "allocs/op",
            "extra": "712 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1649718,
            "unit": "ns/op\t   9.09 MB/s\t 1720228 B/op\t   15002 allocs/op",
            "extra": "703 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1649718,
            "unit": "ns/op",
            "extra": "703 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.09,
            "unit": "MB/s",
            "extra": "703 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 1720228,
            "unit": "B/op",
            "extra": "703 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 15002,
            "unit": "allocs/op",
            "extra": "703 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1651649,
            "unit": "ns/op\t   9.08 MB/s\t 1720233 B/op\t   15002 allocs/op",
            "extra": "692 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1651649,
            "unit": "ns/op",
            "extra": "692 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.08,
            "unit": "MB/s",
            "extra": "692 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 1720233,
            "unit": "B/op",
            "extra": "692 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 15002,
            "unit": "allocs/op",
            "extra": "692 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1723006,
            "unit": "ns/op\t   8.71 MB/s\t 1720228 B/op\t   15002 allocs/op",
            "extra": "686 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1723006,
            "unit": "ns/op",
            "extra": "686 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 8.71,
            "unit": "MB/s",
            "extra": "686 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 1720228,
            "unit": "B/op",
            "extra": "686 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 15002,
            "unit": "allocs/op",
            "extra": "686 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3217,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "369506 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3217,
            "unit": "ns/op",
            "extra": "369506 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "369506 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "369506 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3231,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "374203 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3231,
            "unit": "ns/op",
            "extra": "374203 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "374203 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "374203 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3275,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "372232 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3275,
            "unit": "ns/op",
            "extra": "372232 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "372232 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "372232 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3223,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "373156 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3223,
            "unit": "ns/op",
            "extra": "373156 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "373156 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "373156 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3244,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "372344 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3244,
            "unit": "ns/op",
            "extra": "372344 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "372344 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "372344 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 359.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3432026 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 359.2,
            "unit": "ns/op",
            "extra": "3432026 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3432026 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3432026 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 349.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3488054 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 349.5,
            "unit": "ns/op",
            "extra": "3488054 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3488054 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3488054 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 353.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3347503 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 353.2,
            "unit": "ns/op",
            "extra": "3347503 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3347503 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3347503 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 358.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3334399 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 358.9,
            "unit": "ns/op",
            "extra": "3334399 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3334399 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3334399 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 351.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3361023 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 351.5,
            "unit": "ns/op",
            "extra": "3361023 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3361023 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3361023 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2633,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "435541 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2633,
            "unit": "ns/op",
            "extra": "435541 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "435541 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "435541 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2679,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "434616 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2679,
            "unit": "ns/op",
            "extra": "434616 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "434616 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "434616 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2729,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "429781 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2729,
            "unit": "ns/op",
            "extra": "429781 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "429781 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "429781 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2740,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "442738 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2740,
            "unit": "ns/op",
            "extra": "442738 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "442738 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "442738 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2728,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "443144 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2728,
            "unit": "ns/op",
            "extra": "443144 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "443144 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "443144 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 1010,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1227376 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 1010,
            "unit": "ns/op",
            "extra": "1227376 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1227376 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1227376 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 991.5,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1219794 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 991.5,
            "unit": "ns/op",
            "extra": "1219794 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1219794 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1219794 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 1003,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 1003,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 1032,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1226050 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 1032,
            "unit": "ns/op",
            "extra": "1226050 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1226050 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1226050 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 1001,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 1001,
            "unit": "ns/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1000000 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 247.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6023004 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 247.4,
            "unit": "ns/op",
            "extra": "6023004 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6023004 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6023004 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 218.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7794076 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 218.8,
            "unit": "ns/op",
            "extra": "7794076 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7794076 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7794076 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 219.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7690370 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 219.7,
            "unit": "ns/op",
            "extra": "7690370 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7690370 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7690370 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 220.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7670017 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 220.1,
            "unit": "ns/op",
            "extra": "7670017 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7670017 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7670017 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 219.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7697500 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 219.6,
            "unit": "ns/op",
            "extra": "7697500 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7697500 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7697500 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 16.57,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "69491592 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 16.57,
            "unit": "ns/op",
            "extra": "69491592 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "69491592 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "69491592 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 16.56,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "68488580 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 16.56,
            "unit": "ns/op",
            "extra": "68488580 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "68488580 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "68488580 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 16.57,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "69084684 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 16.57,
            "unit": "ns/op",
            "extra": "69084684 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "69084684 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "69084684 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 16.59,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "72274552 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 16.59,
            "unit": "ns/op",
            "extra": "72274552 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "72274552 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "72274552 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 16.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "70096413 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 16.7,
            "unit": "ns/op",
            "extra": "70096413 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "70096413 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "70096413 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.33,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30359779 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.33,
            "unit": "ns/op",
            "extra": "30359779 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30359779 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30359779 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.33,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30387423 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.33,
            "unit": "ns/op",
            "extra": "30387423 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30387423 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30387423 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.38,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30344460 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.38,
            "unit": "ns/op",
            "extra": "30344460 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30344460 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30344460 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.26,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30556845 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.26,
            "unit": "ns/op",
            "extra": "30556845 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30556845 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30556845 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.34,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30554943 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.34,
            "unit": "ns/op",
            "extra": "30554943 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30554943 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30554943 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 75.25,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15853800 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 75.25,
            "unit": "ns/op",
            "extra": "15853800 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15853800 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15853800 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 75.26,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15803944 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 75.26,
            "unit": "ns/op",
            "extra": "15803944 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15803944 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15803944 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 75.19,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15920377 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 75.19,
            "unit": "ns/op",
            "extra": "15920377 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15920377 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15920377 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 75.27,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15926722 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 75.27,
            "unit": "ns/op",
            "extra": "15926722 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15926722 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15926722 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 75.18,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15853014 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 75.18,
            "unit": "ns/op",
            "extra": "15853014 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15853014 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15853014 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.97,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "28832953 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.97,
            "unit": "ns/op",
            "extra": "28832953 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "28832953 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "28832953 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.85,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30115501 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.85,
            "unit": "ns/op",
            "extra": "30115501 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30115501 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30115501 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.77,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30147968 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.77,
            "unit": "ns/op",
            "extra": "30147968 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30147968 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30147968 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.81,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "29986632 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.81,
            "unit": "ns/op",
            "extra": "29986632 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "29986632 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "29986632 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.86,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "29980016 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.86,
            "unit": "ns/op",
            "extra": "29980016 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "29980016 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "29980016 times\n4 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "ziflex@gmail.com",
            "name": "Tim Voronov",
            "username": "ziflex"
          },
          "committer": {
            "email": "ziflex@gmail.com",
            "name": "Tim Voronov",
            "username": "ziflex"
          },
          "distinct": true,
          "id": "0d5e7e501772623feca6a12ca0bb6a92db0f68c1",
          "message": "docs: revise AGENTS.md for better clarity and expanded development practices",
          "timestamp": "2026-03-19T15:17:46-04:00",
          "tree_id": "6204e745991380a75ec507a280072e7c63e732ac",
          "url": "https://github.com/MontFerret/ferret/commit/0d5e7e501772623feca6a12ca0bb6a92db0f68c1"
        },
        "date": 1773948091826,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 48654,
            "unit": "ns/op\t  82.44 MB/s\t   11152 B/op\t     932 allocs/op",
            "extra": "24752 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 48654,
            "unit": "ns/op",
            "extra": "24752 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 82.44,
            "unit": "MB/s",
            "extra": "24752 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 11152,
            "unit": "B/op",
            "extra": "24752 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 932,
            "unit": "allocs/op",
            "extra": "24752 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 48460,
            "unit": "ns/op\t  82.77 MB/s\t   11152 B/op\t     932 allocs/op",
            "extra": "24914 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 48460,
            "unit": "ns/op",
            "extra": "24914 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 82.77,
            "unit": "MB/s",
            "extra": "24914 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 11152,
            "unit": "B/op",
            "extra": "24914 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 932,
            "unit": "allocs/op",
            "extra": "24914 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 48407,
            "unit": "ns/op\t  82.86 MB/s\t   11152 B/op\t     932 allocs/op",
            "extra": "24778 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 48407,
            "unit": "ns/op",
            "extra": "24778 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 82.86,
            "unit": "MB/s",
            "extra": "24778 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 11152,
            "unit": "B/op",
            "extra": "24778 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 932,
            "unit": "allocs/op",
            "extra": "24778 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 48878,
            "unit": "ns/op\t  82.06 MB/s\t   11152 B/op\t     932 allocs/op",
            "extra": "24922 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 48878,
            "unit": "ns/op",
            "extra": "24922 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 82.06,
            "unit": "MB/s",
            "extra": "24922 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 11152,
            "unit": "B/op",
            "extra": "24922 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 932,
            "unit": "allocs/op",
            "extra": "24922 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 48012,
            "unit": "ns/op\t  83.54 MB/s\t   11152 B/op\t     932 allocs/op",
            "extra": "24962 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 48012,
            "unit": "ns/op",
            "extra": "24962 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 83.54,
            "unit": "MB/s",
            "extra": "24962 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 11152,
            "unit": "B/op",
            "extra": "24962 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 932,
            "unit": "allocs/op",
            "extra": "24962 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 35538,
            "unit": "ns/op\t  65.87 MB/s\t   12867 B/op\t     422 allocs/op",
            "extra": "32568 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 35538,
            "unit": "ns/op",
            "extra": "32568 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 65.87,
            "unit": "MB/s",
            "extra": "32568 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 12867,
            "unit": "B/op",
            "extra": "32568 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 422,
            "unit": "allocs/op",
            "extra": "32568 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 37490,
            "unit": "ns/op\t  62.44 MB/s\t   12867 B/op\t     422 allocs/op",
            "extra": "32572 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 37490,
            "unit": "ns/op",
            "extra": "32572 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 62.44,
            "unit": "MB/s",
            "extra": "32572 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 12867,
            "unit": "B/op",
            "extra": "32572 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 422,
            "unit": "allocs/op",
            "extra": "32572 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 35622,
            "unit": "ns/op\t  65.72 MB/s\t   12867 B/op\t     422 allocs/op",
            "extra": "31945 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 35622,
            "unit": "ns/op",
            "extra": "31945 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 65.72,
            "unit": "MB/s",
            "extra": "31945 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 12867,
            "unit": "B/op",
            "extra": "31945 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 422,
            "unit": "allocs/op",
            "extra": "31945 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 36518,
            "unit": "ns/op\t  64.11 MB/s\t   12867 B/op\t     422 allocs/op",
            "extra": "33702 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 36518,
            "unit": "ns/op",
            "extra": "33702 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 64.11,
            "unit": "MB/s",
            "extra": "33702 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 12867,
            "unit": "B/op",
            "extra": "33702 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 422,
            "unit": "allocs/op",
            "extra": "33702 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 36401,
            "unit": "ns/op\t  64.31 MB/s\t   12867 B/op\t     422 allocs/op",
            "extra": "33247 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 36401,
            "unit": "ns/op",
            "extra": "33247 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 64.31,
            "unit": "MB/s",
            "extra": "33247 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 12867,
            "unit": "B/op",
            "extra": "33247 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 422,
            "unit": "allocs/op",
            "extra": "33247 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 503561,
            "unit": "ns/op\t  39.72 MB/s\t   65520 B/op\t      11 allocs/op",
            "extra": "2451 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 503561,
            "unit": "ns/op",
            "extra": "2451 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 39.72,
            "unit": "MB/s",
            "extra": "2451 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 65520,
            "unit": "B/op",
            "extra": "2451 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "2451 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 519413,
            "unit": "ns/op\t  38.51 MB/s\t   65520 B/op\t      11 allocs/op",
            "extra": "2419 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 519413,
            "unit": "ns/op",
            "extra": "2419 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 38.51,
            "unit": "MB/s",
            "extra": "2419 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 65520,
            "unit": "B/op",
            "extra": "2419 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "2419 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 506470,
            "unit": "ns/op\t  39.49 MB/s\t   65520 B/op\t      11 allocs/op",
            "extra": "2408 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 506470,
            "unit": "ns/op",
            "extra": "2408 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 39.49,
            "unit": "MB/s",
            "extra": "2408 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 65520,
            "unit": "B/op",
            "extra": "2408 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "2408 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 490161,
            "unit": "ns/op\t  40.80 MB/s\t   65520 B/op\t      11 allocs/op",
            "extra": "2518 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 490161,
            "unit": "ns/op",
            "extra": "2518 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 40.8,
            "unit": "MB/s",
            "extra": "2518 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 65520,
            "unit": "B/op",
            "extra": "2518 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "2518 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 515283,
            "unit": "ns/op\t  38.82 MB/s\t   65520 B/op\t      11 allocs/op",
            "extra": "2300 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 515283,
            "unit": "ns/op",
            "extra": "2300 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 38.82,
            "unit": "MB/s",
            "extra": "2300 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 65520,
            "unit": "B/op",
            "extra": "2300 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 11,
            "unit": "allocs/op",
            "extra": "2300 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1434271,
            "unit": "ns/op\t  20.92 MB/s\t  630530 B/op\t   15011 allocs/op",
            "extra": "825 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1434271,
            "unit": "ns/op",
            "extra": "825 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 20.92,
            "unit": "MB/s",
            "extra": "825 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 630530,
            "unit": "B/op",
            "extra": "825 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 15011,
            "unit": "allocs/op",
            "extra": "825 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1438332,
            "unit": "ns/op\t  20.86 MB/s\t  630530 B/op\t   15011 allocs/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1438332,
            "unit": "ns/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 20.86,
            "unit": "MB/s",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 630530,
            "unit": "B/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 15011,
            "unit": "allocs/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1396970,
            "unit": "ns/op\t  21.48 MB/s\t  630531 B/op\t   15011 allocs/op",
            "extra": "806 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1396970,
            "unit": "ns/op",
            "extra": "806 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 21.48,
            "unit": "MB/s",
            "extra": "806 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 630531,
            "unit": "B/op",
            "extra": "806 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 15011,
            "unit": "allocs/op",
            "extra": "806 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1448221,
            "unit": "ns/op\t  20.72 MB/s\t  630530 B/op\t   15011 allocs/op",
            "extra": "850 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1448221,
            "unit": "ns/op",
            "extra": "850 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 20.72,
            "unit": "MB/s",
            "extra": "850 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 630530,
            "unit": "B/op",
            "extra": "850 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 15011,
            "unit": "allocs/op",
            "extra": "850 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1413246,
            "unit": "ns/op\t  21.23 MB/s\t  630530 B/op\t   15011 allocs/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1413246,
            "unit": "ns/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 21.23,
            "unit": "MB/s",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 630530,
            "unit": "B/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 15011,
            "unit": "allocs/op",
            "extra": "841 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 116795,
            "unit": "ns/op\t  34.34 MB/s\t   90121 B/op\t    1812 allocs/op",
            "extra": "10347 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 116795,
            "unit": "ns/op",
            "extra": "10347 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 34.34,
            "unit": "MB/s",
            "extra": "10347 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 90121,
            "unit": "B/op",
            "extra": "10347 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "10347 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 115634,
            "unit": "ns/op\t  34.69 MB/s\t   90121 B/op\t    1812 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 115634,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 34.69,
            "unit": "MB/s",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 90121,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 114989,
            "unit": "ns/op\t  34.88 MB/s\t   90121 B/op\t    1812 allocs/op",
            "extra": "9154 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 114989,
            "unit": "ns/op",
            "extra": "9154 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 34.88,
            "unit": "MB/s",
            "extra": "9154 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 90121,
            "unit": "B/op",
            "extra": "9154 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "9154 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 115071,
            "unit": "ns/op\t  34.86 MB/s\t   90121 B/op\t    1812 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 115071,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 34.86,
            "unit": "MB/s",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 90121,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 116163,
            "unit": "ns/op\t  34.53 MB/s\t   90121 B/op\t    1812 allocs/op",
            "extra": "9939 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 116163,
            "unit": "ns/op",
            "extra": "9939 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 34.53,
            "unit": "MB/s",
            "extra": "9939 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 90121,
            "unit": "B/op",
            "extra": "9939 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 1812,
            "unit": "allocs/op",
            "extra": "9939 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 74101,
            "unit": "ns/op\t  31.59 MB/s\t   54065 B/op\t     781 allocs/op",
            "extra": "16477 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 74101,
            "unit": "ns/op",
            "extra": "16477 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 31.59,
            "unit": "MB/s",
            "extra": "16477 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 54065,
            "unit": "B/op",
            "extra": "16477 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 781,
            "unit": "allocs/op",
            "extra": "16477 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 73933,
            "unit": "ns/op\t  31.66 MB/s\t   54065 B/op\t     781 allocs/op",
            "extra": "16375 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 73933,
            "unit": "ns/op",
            "extra": "16375 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 31.66,
            "unit": "MB/s",
            "extra": "16375 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 54065,
            "unit": "B/op",
            "extra": "16375 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 781,
            "unit": "allocs/op",
            "extra": "16375 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 74389,
            "unit": "ns/op\t  31.47 MB/s\t   54065 B/op\t     781 allocs/op",
            "extra": "16125 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 74389,
            "unit": "ns/op",
            "extra": "16125 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 31.47,
            "unit": "MB/s",
            "extra": "16125 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 54065,
            "unit": "B/op",
            "extra": "16125 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 781,
            "unit": "allocs/op",
            "extra": "16125 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 75358,
            "unit": "ns/op\t  31.07 MB/s\t   54065 B/op\t     781 allocs/op",
            "extra": "16543 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 75358,
            "unit": "ns/op",
            "extra": "16543 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 31.07,
            "unit": "MB/s",
            "extra": "16543 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 54065,
            "unit": "B/op",
            "extra": "16543 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 781,
            "unit": "allocs/op",
            "extra": "16543 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 73107,
            "unit": "ns/op\t  32.02 MB/s\t   54065 B/op\t     781 allocs/op",
            "extra": "16347 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 73107,
            "unit": "ns/op",
            "extra": "16347 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 32.02,
            "unit": "MB/s",
            "extra": "16347 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 54065,
            "unit": "B/op",
            "extra": "16347 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 781,
            "unit": "allocs/op",
            "extra": "16347 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1902792,
            "unit": "ns/op\t  10.51 MB/s\t 2307998 B/op\t   20030 allocs/op",
            "extra": "638 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1902792,
            "unit": "ns/op",
            "extra": "638 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 10.51,
            "unit": "MB/s",
            "extra": "638 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2307998,
            "unit": "B/op",
            "extra": "638 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20030,
            "unit": "allocs/op",
            "extra": "638 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1922853,
            "unit": "ns/op\t  10.40 MB/s\t 2308002 B/op\t   20030 allocs/op",
            "extra": "620 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1922853,
            "unit": "ns/op",
            "extra": "620 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 10.4,
            "unit": "MB/s",
            "extra": "620 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2308002,
            "unit": "B/op",
            "extra": "620 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20030,
            "unit": "allocs/op",
            "extra": "620 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1958287,
            "unit": "ns/op\t  10.21 MB/s\t 2308002 B/op\t   20030 allocs/op",
            "extra": "594 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1958287,
            "unit": "ns/op",
            "extra": "594 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 10.21,
            "unit": "MB/s",
            "extra": "594 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2308002,
            "unit": "B/op",
            "extra": "594 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20030,
            "unit": "allocs/op",
            "extra": "594 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 2168369,
            "unit": "ns/op\t   9.22 MB/s\t 2308003 B/op\t   20030 allocs/op",
            "extra": "532 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 2168369,
            "unit": "ns/op",
            "extra": "532 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 9.22,
            "unit": "MB/s",
            "extra": "532 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2308003,
            "unit": "B/op",
            "extra": "532 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20030,
            "unit": "allocs/op",
            "extra": "532 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 2154189,
            "unit": "ns/op\t   9.28 MB/s\t 2308002 B/op\t   20030 allocs/op",
            "extra": "596 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 2154189,
            "unit": "ns/op",
            "extra": "596 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 9.28,
            "unit": "MB/s",
            "extra": "596 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2308002,
            "unit": "B/op",
            "extra": "596 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20030,
            "unit": "allocs/op",
            "extra": "596 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1685141,
            "unit": "ns/op\t  17.80 MB/s\t 2610265 B/op\t   20026 allocs/op",
            "extra": "714 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1685141,
            "unit": "ns/op",
            "extra": "714 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 17.8,
            "unit": "MB/s",
            "extra": "714 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2610265,
            "unit": "B/op",
            "extra": "714 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20026,
            "unit": "allocs/op",
            "extra": "714 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1656201,
            "unit": "ns/op\t  18.11 MB/s\t 2610262 B/op\t   20026 allocs/op",
            "extra": "716 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1656201,
            "unit": "ns/op",
            "extra": "716 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 18.11,
            "unit": "MB/s",
            "extra": "716 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2610262,
            "unit": "B/op",
            "extra": "716 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20026,
            "unit": "allocs/op",
            "extra": "716 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1669401,
            "unit": "ns/op\t  17.97 MB/s\t 2610263 B/op\t   20026 allocs/op",
            "extra": "733 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1669401,
            "unit": "ns/op",
            "extra": "733 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 17.97,
            "unit": "MB/s",
            "extra": "733 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2610263,
            "unit": "B/op",
            "extra": "733 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20026,
            "unit": "allocs/op",
            "extra": "733 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1664893,
            "unit": "ns/op\t  18.02 MB/s\t 2610263 B/op\t   20026 allocs/op",
            "extra": "699 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1664893,
            "unit": "ns/op",
            "extra": "699 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 18.02,
            "unit": "MB/s",
            "extra": "699 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2610263,
            "unit": "B/op",
            "extra": "699 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20026,
            "unit": "allocs/op",
            "extra": "699 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json)",
            "value": 1671182,
            "unit": "ns/op\t  17.95 MB/s\t 2610263 B/op\t   20026 allocs/op",
            "extra": "712 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - ns/op",
            "value": 1671182,
            "unit": "ns/op",
            "extra": "712 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - MB/s",
            "value": 17.95,
            "unit": "MB/s",
            "extra": "712 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - B/op",
            "value": 2610263,
            "unit": "B/op",
            "extra": "712 times\n4 procs"
          },
          {
            "name": "BenchmarkJSONCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/json) - allocs/op",
            "value": 20026,
            "unit": "allocs/op",
            "extra": "712 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 26502,
            "unit": "ns/op\t 101.54 MB/s\t    8177 B/op\t       8 allocs/op",
            "extra": "45178 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 26502,
            "unit": "ns/op",
            "extra": "45178 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 101.54,
            "unit": "MB/s",
            "extra": "45178 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8177,
            "unit": "B/op",
            "extra": "45178 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "45178 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 26718,
            "unit": "ns/op\t 100.72 MB/s\t    8177 B/op\t       8 allocs/op",
            "extra": "44583 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 26718,
            "unit": "ns/op",
            "extra": "44583 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 100.72,
            "unit": "MB/s",
            "extra": "44583 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8177,
            "unit": "B/op",
            "extra": "44583 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "44583 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 26597,
            "unit": "ns/op\t 101.18 MB/s\t    8177 B/op\t       8 allocs/op",
            "extra": "45854 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 26597,
            "unit": "ns/op",
            "extra": "45854 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 101.18,
            "unit": "MB/s",
            "extra": "45854 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8177,
            "unit": "B/op",
            "extra": "45854 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "45854 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 26028,
            "unit": "ns/op\t 103.39 MB/s\t    8177 B/op\t       8 allocs/op",
            "extra": "45967 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 26028,
            "unit": "ns/op",
            "extra": "45967 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 103.39,
            "unit": "MB/s",
            "extra": "45967 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8177,
            "unit": "B/op",
            "extra": "45967 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "45967 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 25811,
            "unit": "ns/op\t 104.26 MB/s\t    8177 B/op\t       8 allocs/op",
            "extra": "45547 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 25811,
            "unit": "ns/op",
            "extra": "45547 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 104.26,
            "unit": "MB/s",
            "extra": "45547 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8177,
            "unit": "B/op",
            "extra": "45547 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 8,
            "unit": "allocs/op",
            "extra": "45547 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 22910,
            "unit": "ns/op\t  56.79 MB/s\t    8257 B/op\t     264 allocs/op",
            "extra": "53095 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 22910,
            "unit": "ns/op",
            "extra": "53095 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 56.79,
            "unit": "MB/s",
            "extra": "53095 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8257,
            "unit": "B/op",
            "extra": "53095 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 264,
            "unit": "allocs/op",
            "extra": "53095 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 22753,
            "unit": "ns/op\t  57.18 MB/s\t    8257 B/op\t     264 allocs/op",
            "extra": "52514 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 22753,
            "unit": "ns/op",
            "extra": "52514 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 57.18,
            "unit": "MB/s",
            "extra": "52514 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8257,
            "unit": "B/op",
            "extra": "52514 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 264,
            "unit": "allocs/op",
            "extra": "52514 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 23108,
            "unit": "ns/op\t  56.30 MB/s\t    8257 B/op\t     264 allocs/op",
            "extra": "52963 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 23108,
            "unit": "ns/op",
            "extra": "52963 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 56.3,
            "unit": "MB/s",
            "extra": "52963 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8257,
            "unit": "B/op",
            "extra": "52963 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 264,
            "unit": "allocs/op",
            "extra": "52963 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 22574,
            "unit": "ns/op\t  57.63 MB/s\t    8257 B/op\t     264 allocs/op",
            "extra": "52518 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 22574,
            "unit": "ns/op",
            "extra": "52518 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 57.63,
            "unit": "MB/s",
            "extra": "52518 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8257,
            "unit": "B/op",
            "extra": "52518 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 264,
            "unit": "allocs/op",
            "extra": "52518 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 22605,
            "unit": "ns/op\t  57.55 MB/s\t    8257 B/op\t     264 allocs/op",
            "extra": "53682 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 22605,
            "unit": "ns/op",
            "extra": "53682 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 57.55,
            "unit": "MB/s",
            "extra": "53682 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 8257,
            "unit": "B/op",
            "extra": "53682 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 264,
            "unit": "allocs/op",
            "extra": "53682 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 330419,
            "unit": "ns/op\t  30.27 MB/s\t   32756 B/op\t      10 allocs/op",
            "extra": "3088 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 330419,
            "unit": "ns/op",
            "extra": "3088 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 30.27,
            "unit": "MB/s",
            "extra": "3088 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 32756,
            "unit": "B/op",
            "extra": "3088 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "3088 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 370610,
            "unit": "ns/op\t  26.99 MB/s\t   32756 B/op\t      10 allocs/op",
            "extra": "3146 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 370610,
            "unit": "ns/op",
            "extra": "3146 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 26.99,
            "unit": "MB/s",
            "extra": "3146 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 32756,
            "unit": "B/op",
            "extra": "3146 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "3146 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 375892,
            "unit": "ns/op\t  26.61 MB/s\t   32756 B/op\t      10 allocs/op",
            "extra": "2971 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 375892,
            "unit": "ns/op",
            "extra": "2971 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 26.61,
            "unit": "MB/s",
            "extra": "2971 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 32756,
            "unit": "B/op",
            "extra": "2971 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "2971 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 387674,
            "unit": "ns/op\t  25.80 MB/s\t   32756 B/op\t      10 allocs/op",
            "extra": "3067 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 387674,
            "unit": "ns/op",
            "extra": "3067 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 25.8,
            "unit": "MB/s",
            "extra": "3067 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 32756,
            "unit": "B/op",
            "extra": "3067 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "3067 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 351245,
            "unit": "ns/op\t  28.47 MB/s\t   32756 B/op\t      10 allocs/op",
            "extra": "3171 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 351245,
            "unit": "ns/op",
            "extra": "3171 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 28.47,
            "unit": "MB/s",
            "extra": "3171 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 32756,
            "unit": "B/op",
            "extra": "3171 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "3171 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1267297,
            "unit": "ns/op\t  11.84 MB/s\t  512826 B/op\t   10010 allocs/op",
            "extra": "933 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1267297,
            "unit": "ns/op",
            "extra": "933 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 11.84,
            "unit": "MB/s",
            "extra": "933 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 512826,
            "unit": "B/op",
            "extra": "933 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10010,
            "unit": "allocs/op",
            "extra": "933 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1307302,
            "unit": "ns/op\t  11.47 MB/s\t  512826 B/op\t   10010 allocs/op",
            "extra": "908 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1307302,
            "unit": "ns/op",
            "extra": "908 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 11.47,
            "unit": "MB/s",
            "extra": "908 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 512826,
            "unit": "B/op",
            "extra": "908 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10010,
            "unit": "allocs/op",
            "extra": "908 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1243067,
            "unit": "ns/op\t  12.07 MB/s\t  512831 B/op\t   10010 allocs/op",
            "extra": "912 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1243067,
            "unit": "ns/op",
            "extra": "912 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 12.07,
            "unit": "MB/s",
            "extra": "912 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 512831,
            "unit": "B/op",
            "extra": "912 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10010,
            "unit": "allocs/op",
            "extra": "912 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1282256,
            "unit": "ns/op\t  11.70 MB/s\t  512825 B/op\t   10010 allocs/op",
            "extra": "915 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1282256,
            "unit": "ns/op",
            "extra": "915 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 11.7,
            "unit": "MB/s",
            "extra": "915 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 512825,
            "unit": "B/op",
            "extra": "915 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10010,
            "unit": "allocs/op",
            "extra": "915 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1302181,
            "unit": "ns/op\t  11.52 MB/s\t  512827 B/op\t   10010 allocs/op",
            "extra": "926 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1302181,
            "unit": "ns/op",
            "extra": "926 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 11.52,
            "unit": "MB/s",
            "extra": "926 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 512827,
            "unit": "B/op",
            "extra": "926 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecEncode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 10010,
            "unit": "allocs/op",
            "extra": "926 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 52191,
            "unit": "ns/op\t  51.56 MB/s\t   24653 B/op\t     771 allocs/op",
            "extra": "22669 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 52191,
            "unit": "ns/op",
            "extra": "22669 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 51.56,
            "unit": "MB/s",
            "extra": "22669 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 24653,
            "unit": "B/op",
            "extra": "22669 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 771,
            "unit": "allocs/op",
            "extra": "22669 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 51474,
            "unit": "ns/op\t  52.28 MB/s\t   24653 B/op\t     771 allocs/op",
            "extra": "23252 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 51474,
            "unit": "ns/op",
            "extra": "23252 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 52.28,
            "unit": "MB/s",
            "extra": "23252 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 24653,
            "unit": "B/op",
            "extra": "23252 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 771,
            "unit": "allocs/op",
            "extra": "23252 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 52405,
            "unit": "ns/op\t  51.35 MB/s\t   24653 B/op\t     771 allocs/op",
            "extra": "22873 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 52405,
            "unit": "ns/op",
            "extra": "22873 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 51.35,
            "unit": "MB/s",
            "extra": "22873 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 24653,
            "unit": "B/op",
            "extra": "22873 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 771,
            "unit": "allocs/op",
            "extra": "22873 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 52507,
            "unit": "ns/op\t  51.25 MB/s\t   24653 B/op\t     771 allocs/op",
            "extra": "22908 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 52507,
            "unit": "ns/op",
            "extra": "22908 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 51.25,
            "unit": "MB/s",
            "extra": "22908 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 24653,
            "unit": "B/op",
            "extra": "22908 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 771,
            "unit": "allocs/op",
            "extra": "22908 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 52618,
            "unit": "ns/op\t  51.14 MB/s\t   24653 B/op\t     771 allocs/op",
            "extra": "23029 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 52618,
            "unit": "ns/op",
            "extra": "23029 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 51.14,
            "unit": "MB/s",
            "extra": "23029 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 24653,
            "unit": "B/op",
            "extra": "23029 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_array_1024 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 771,
            "unit": "allocs/op",
            "extra": "23029 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 30206,
            "unit": "ns/op\t  43.07 MB/s\t   19299 B/op\t     252 allocs/op",
            "extra": "40209 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 30206,
            "unit": "ns/op",
            "extra": "40209 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 43.07,
            "unit": "MB/s",
            "extra": "40209 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 19299,
            "unit": "B/op",
            "extra": "40209 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 252,
            "unit": "allocs/op",
            "extra": "40209 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 30142,
            "unit": "ns/op\t  43.16 MB/s\t   19299 B/op\t     252 allocs/op",
            "extra": "39253 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 30142,
            "unit": "ns/op",
            "extra": "39253 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 43.16,
            "unit": "MB/s",
            "extra": "39253 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 19299,
            "unit": "B/op",
            "extra": "39253 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 252,
            "unit": "allocs/op",
            "extra": "39253 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 30264,
            "unit": "ns/op\t  42.99 MB/s\t   19299 B/op\t     252 allocs/op",
            "extra": "40026 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 30264,
            "unit": "ns/op",
            "extra": "40026 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 42.99,
            "unit": "MB/s",
            "extra": "40026 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 19299,
            "unit": "B/op",
            "extra": "40026 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 252,
            "unit": "allocs/op",
            "extra": "40026 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 29660,
            "unit": "ns/op\t  43.86 MB/s\t   19299 B/op\t     252 allocs/op",
            "extra": "40396 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 29660,
            "unit": "ns/op",
            "extra": "40396 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 43.86,
            "unit": "MB/s",
            "extra": "40396 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 19299,
            "unit": "B/op",
            "extra": "40396 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 252,
            "unit": "allocs/op",
            "extra": "40396 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 30563,
            "unit": "ns/op\t  42.57 MB/s\t   19299 B/op\t     252 allocs/op",
            "extra": "39894 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 30563,
            "unit": "ns/op",
            "extra": "39894 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 42.57,
            "unit": "MB/s",
            "extra": "39894 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 19299,
            "unit": "B/op",
            "extra": "39894 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/flat_object_256 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 252,
            "unit": "allocs/op",
            "extra": "39894 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1130116,
            "unit": "ns/op\t   8.85 MB/s\t  400102 B/op\t   20001 allocs/op",
            "extra": "1092 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1130116,
            "unit": "ns/op",
            "extra": "1092 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 8.85,
            "unit": "MB/s",
            "extra": "1092 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 400102,
            "unit": "B/op",
            "extra": "1092 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 20001,
            "unit": "allocs/op",
            "extra": "1092 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1129925,
            "unit": "ns/op\t   8.85 MB/s\t  400103 B/op\t   20001 allocs/op",
            "extra": "1076 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1129925,
            "unit": "ns/op",
            "extra": "1076 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 8.85,
            "unit": "MB/s",
            "extra": "1076 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 400103,
            "unit": "B/op",
            "extra": "1076 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 20001,
            "unit": "allocs/op",
            "extra": "1076 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1104668,
            "unit": "ns/op\t   9.05 MB/s\t  400103 B/op\t   20001 allocs/op",
            "extra": "1002 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1104668,
            "unit": "ns/op",
            "extra": "1002 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.05,
            "unit": "MB/s",
            "extra": "1002 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 400103,
            "unit": "B/op",
            "extra": "1002 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 20001,
            "unit": "allocs/op",
            "extra": "1002 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1102602,
            "unit": "ns/op\t   9.07 MB/s\t  400105 B/op\t   20001 allocs/op",
            "extra": "1056 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1102602,
            "unit": "ns/op",
            "extra": "1056 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.07,
            "unit": "MB/s",
            "extra": "1056 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 400105,
            "unit": "B/op",
            "extra": "1056 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 20001,
            "unit": "allocs/op",
            "extra": "1056 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1100939,
            "unit": "ns/op\t   9.08 MB/s\t  400103 B/op\t   20001 allocs/op",
            "extra": "1084 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1100939,
            "unit": "ns/op",
            "extra": "1084 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.08,
            "unit": "MB/s",
            "extra": "1084 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 400103,
            "unit": "B/op",
            "extra": "1084 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_array_10000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 20001,
            "unit": "allocs/op",
            "extra": "1084 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1642535,
            "unit": "ns/op\t   9.13 MB/s\t 1720225 B/op\t   15002 allocs/op",
            "extra": "703 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1642535,
            "unit": "ns/op",
            "extra": "703 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.13,
            "unit": "MB/s",
            "extra": "703 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 1720225,
            "unit": "B/op",
            "extra": "703 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 15002,
            "unit": "allocs/op",
            "extra": "703 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1648161,
            "unit": "ns/op\t   9.10 MB/s\t 1720229 B/op\t   15002 allocs/op",
            "extra": "716 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1648161,
            "unit": "ns/op",
            "extra": "716 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.1,
            "unit": "MB/s",
            "extra": "716 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 1720229,
            "unit": "B/op",
            "extra": "716 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 15002,
            "unit": "allocs/op",
            "extra": "716 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1662615,
            "unit": "ns/op\t   9.02 MB/s\t 1720223 B/op\t   15002 allocs/op",
            "extra": "742 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1662615,
            "unit": "ns/op",
            "extra": "742 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.02,
            "unit": "MB/s",
            "extra": "742 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 1720223,
            "unit": "B/op",
            "extra": "742 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 15002,
            "unit": "allocs/op",
            "extra": "742 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1616613,
            "unit": "ns/op\t   9.28 MB/s\t 1720226 B/op\t   15002 allocs/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1616613,
            "unit": "ns/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.28,
            "unit": "MB/s",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 1720226,
            "unit": "B/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 15002,
            "unit": "allocs/op",
            "extra": "735 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack)",
            "value": 1653401,
            "unit": "ns/op\t   9.07 MB/s\t 1720228 B/op\t   15002 allocs/op",
            "extra": "759 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - ns/op",
            "value": 1653401,
            "unit": "ns/op",
            "extra": "759 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - MB/s",
            "value": 9.07,
            "unit": "MB/s",
            "extra": "759 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - B/op",
            "value": 1720228,
            "unit": "B/op",
            "extra": "759 times\n4 procs"
          },
          {
            "name": "BenchmarkMsgpackCodecDecode/nested_object_5000 (github.com/MontFerret/ferret/v2/pkg/encoding/msgpack) - allocs/op",
            "value": 15002,
            "unit": "allocs/op",
            "extra": "759 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3207,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "366848 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3207,
            "unit": "ns/op",
            "extra": "366848 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "366848 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "366848 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3207,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "374524 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3207,
            "unit": "ns/op",
            "extra": "374524 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "374524 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "374524 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3210,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "374510 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3210,
            "unit": "ns/op",
            "extra": "374510 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "374510 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "374510 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3215,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "374104 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3215,
            "unit": "ns/op",
            "extra": "374104 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "374104 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "374104 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3209,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "372862 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3209,
            "unit": "ns/op",
            "extra": "372862 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "372862 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "372862 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 356.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3421304 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 356.8,
            "unit": "ns/op",
            "extra": "3421304 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3421304 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3421304 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 357.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3657674 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 357.8,
            "unit": "ns/op",
            "extra": "3657674 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3657674 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3657674 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 355.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3374506 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 355.1,
            "unit": "ns/op",
            "extra": "3374506 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3374506 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3374506 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 355.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3670852 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 355.8,
            "unit": "ns/op",
            "extra": "3670852 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3670852 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3670852 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 356.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3472069 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 356.4,
            "unit": "ns/op",
            "extra": "3472069 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3472069 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3472069 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2497,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "453793 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2497,
            "unit": "ns/op",
            "extra": "453793 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "453793 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "453793 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2601,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "453992 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2601,
            "unit": "ns/op",
            "extra": "453992 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "453992 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "453992 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2577,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "471111 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2577,
            "unit": "ns/op",
            "extra": "471111 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "471111 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "471111 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2500,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "459560 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2500,
            "unit": "ns/op",
            "extra": "459560 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "459560 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "459560 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2554,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "468922 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2554,
            "unit": "ns/op",
            "extra": "468922 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "468922 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "468922 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 978.9,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1244442 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 978.9,
            "unit": "ns/op",
            "extra": "1244442 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1244442 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1244442 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 996.2,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1218358 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 996.2,
            "unit": "ns/op",
            "extra": "1218358 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1218358 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1218358 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 971.1,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1230114 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 971.1,
            "unit": "ns/op",
            "extra": "1230114 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1230114 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1230114 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 981,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1213753 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 981,
            "unit": "ns/op",
            "extra": "1213753 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1213753 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1213753 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 978.5,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1244167 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 978.5,
            "unit": "ns/op",
            "extra": "1244167 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1244167 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1244167 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 193.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7592046 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 193.2,
            "unit": "ns/op",
            "extra": "7592046 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7592046 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7592046 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 206.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7494202 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 206.2,
            "unit": "ns/op",
            "extra": "7494202 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7494202 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7494202 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 236.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6600673 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 236.9,
            "unit": "ns/op",
            "extra": "6600673 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6600673 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6600673 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 195.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7409101 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 195.8,
            "unit": "ns/op",
            "extra": "7409101 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7409101 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7409101 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 220.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7676739 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 220.5,
            "unit": "ns/op",
            "extra": "7676739 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7676739 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_New (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7676739 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 16.58,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "72268276 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 16.58,
            "unit": "ns/op",
            "extra": "72268276 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "72268276 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "72268276 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 16.52,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "70462699 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 16.52,
            "unit": "ns/op",
            "extra": "70462699 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "70462699 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "70462699 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 16.57,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "70881438 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 16.57,
            "unit": "ns/op",
            "extra": "70881438 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "70881438 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "70881438 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 16.54,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "70188490 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 16.54,
            "unit": "ns/op",
            "extra": "70188490 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "70188490 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "70188490 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 16.57,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "72394336 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 16.57,
            "unit": "ns/op",
            "extra": "72394336 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "72394336 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Get (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "72394336 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.37,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30445576 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.37,
            "unit": "ns/op",
            "extra": "30445576 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30445576 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30445576 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.31,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30624315 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.31,
            "unit": "ns/op",
            "extra": "30624315 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30624315 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30624315 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.32,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30586862 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.32,
            "unit": "ns/op",
            "extra": "30586862 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30586862 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30586862 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.31,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30293961 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.31,
            "unit": "ns/op",
            "extra": "30293961 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30293961 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30293961 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.33,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30486039 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.33,
            "unit": "ns/op",
            "extra": "30486039 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30486039 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_Set (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30486039 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 75.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15892442 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 75.29,
            "unit": "ns/op",
            "extra": "15892442 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15892442 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15892442 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 75.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15935083 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 75.29,
            "unit": "ns/op",
            "extra": "15935083 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15935083 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15935083 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 75.21,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15884740 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 75.21,
            "unit": "ns/op",
            "extra": "15884740 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15884740 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15884740 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 75.19,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15923805 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 75.19,
            "unit": "ns/op",
            "extra": "15923805 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15923805 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15923805 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 75.34,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15894831 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 75.34,
            "unit": "ns/op",
            "extra": "15894831 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15894831 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_DeleteThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15894831 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30150254 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.83,
            "unit": "ns/op",
            "extra": "30150254 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30150254 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30150254 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.77,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "29976309 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.77,
            "unit": "ns/op",
            "extra": "29976309 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "29976309 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "29976309 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "29771079 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.8,
            "unit": "ns/op",
            "extra": "29771079 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "29771079 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "29771079 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.85,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30155534 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.85,
            "unit": "ns/op",
            "extra": "30155534 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30155534 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30155534 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem)",
            "value": 39.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30110316 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - ns/op",
            "value": 39.83,
            "unit": "ns/op",
            "extra": "30110316 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30110316 times\n4 procs"
          },
          {
            "name": "BenchmarkCellStore_ResetThenNew (github.com/MontFerret/ferret/v2/pkg/vm/internal/mem) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30110316 times\n4 procs"
          }
        ]
      }
    ]
  }
}