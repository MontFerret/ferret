window.BENCHMARK_DATA = {
  "lastUpdate": 1773946600593,
  "repoUrl": "https://github.com/MontFerret/ferret",
  "entries": {
    "Ferret Go Benchmarks - Integration": [
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
        "date": 1773946599402,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkAddNumeric_O0",
            "value": 96827,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "12405 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 - ns/op",
            "value": 96827,
            "unit": "ns/op",
            "extra": "12405 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "12405 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12405 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0",
            "value": 97727,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "12032 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 - ns/op",
            "value": 97727,
            "unit": "ns/op",
            "extra": "12032 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "12032 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12032 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0",
            "value": 97146,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "12290 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 - ns/op",
            "value": 97146,
            "unit": "ns/op",
            "extra": "12290 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "12290 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12290 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0",
            "value": 97933,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "12163 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 - ns/op",
            "value": 97933,
            "unit": "ns/op",
            "extra": "12163 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "12163 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12163 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0",
            "value": 105654,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "12264 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 - ns/op",
            "value": 105654,
            "unit": "ns/op",
            "extra": "12264 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "12264 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12264 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1",
            "value": 92303,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "12968 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 - ns/op",
            "value": 92303,
            "unit": "ns/op",
            "extra": "12968 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "12968 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12968 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1",
            "value": 94054,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "12537 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 - ns/op",
            "value": 94054,
            "unit": "ns/op",
            "extra": "12537 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "12537 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12537 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1",
            "value": 94253,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "12910 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 - ns/op",
            "value": 94253,
            "unit": "ns/op",
            "extra": "12910 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "12910 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12910 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1",
            "value": 94607,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "12752 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 - ns/op",
            "value": 94607,
            "unit": "ns/op",
            "extra": "12752 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "12752 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12752 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1",
            "value": 93909,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "12537 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 - ns/op",
            "value": 93909,
            "unit": "ns/op",
            "extra": "12537 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "12537 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12537 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0",
            "value": 108606,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 - ns/op",
            "value": 108606,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0",
            "value": 107773,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 - ns/op",
            "value": 107773,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0",
            "value": 109252,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 - ns/op",
            "value": 109252,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0",
            "value": 106562,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 - ns/op",
            "value": 106562,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0",
            "value": 108829,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 - ns/op",
            "value": 108829,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1",
            "value": 101796,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 - ns/op",
            "value": 101796,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1",
            "value": 101854,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 - ns/op",
            "value": 101854,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1",
            "value": 102321,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 - ns/op",
            "value": 102321,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1",
            "value": 102234,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 - ns/op",
            "value": 102234,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1",
            "value": 101991,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 - ns/op",
            "value": 101991,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0",
            "value": 148664,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7396 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 - ns/op",
            "value": 148664,
            "unit": "ns/op",
            "extra": "7396 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7396 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7396 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0",
            "value": 149715,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7758 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 - ns/op",
            "value": 149715,
            "unit": "ns/op",
            "extra": "7758 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7758 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7758 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0",
            "value": 149263,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7581 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 - ns/op",
            "value": 149263,
            "unit": "ns/op",
            "extra": "7581 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7581 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7581 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0",
            "value": 150122,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7339 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 - ns/op",
            "value": 150122,
            "unit": "ns/op",
            "extra": "7339 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7339 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7339 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0",
            "value": 150198,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7482 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 - ns/op",
            "value": 150198,
            "unit": "ns/op",
            "extra": "7482 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7482 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7482 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1",
            "value": 147457,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7789 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 - ns/op",
            "value": 147457,
            "unit": "ns/op",
            "extra": "7789 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7789 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7789 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1",
            "value": 148993,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "8014 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 - ns/op",
            "value": 148993,
            "unit": "ns/op",
            "extra": "8014 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "8014 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "8014 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1",
            "value": 147626,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7657 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 - ns/op",
            "value": 147626,
            "unit": "ns/op",
            "extra": "7657 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7657 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7657 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1",
            "value": 149077,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7776 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 - ns/op",
            "value": 149077,
            "unit": "ns/op",
            "extra": "7776 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7776 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7776 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1",
            "value": 147974,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7844 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 - ns/op",
            "value": 147974,
            "unit": "ns/op",
            "extra": "7844 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7844 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7844 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0",
            "value": 156635,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7380 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 - ns/op",
            "value": 156635,
            "unit": "ns/op",
            "extra": "7380 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7380 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7380 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0",
            "value": 157061,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7306 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 - ns/op",
            "value": 157061,
            "unit": "ns/op",
            "extra": "7306 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7306 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7306 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0",
            "value": 157464,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7698 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 - ns/op",
            "value": 157464,
            "unit": "ns/op",
            "extra": "7698 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7698 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7698 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0",
            "value": 156641,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7410 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 - ns/op",
            "value": 156641,
            "unit": "ns/op",
            "extra": "7410 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7410 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7410 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0",
            "value": 156760,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7020 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 - ns/op",
            "value": 156760,
            "unit": "ns/op",
            "extra": "7020 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7020 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7020 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1",
            "value": 156091,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7244 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 - ns/op",
            "value": 156091,
            "unit": "ns/op",
            "extra": "7244 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7244 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7244 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1",
            "value": 156740,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7257 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 - ns/op",
            "value": 156740,
            "unit": "ns/op",
            "extra": "7257 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7257 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7257 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1",
            "value": 156969,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7242 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 - ns/op",
            "value": 156969,
            "unit": "ns/op",
            "extra": "7242 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7242 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7242 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1",
            "value": 155908,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7185 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 - ns/op",
            "value": 155908,
            "unit": "ns/op",
            "extra": "7185 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7185 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7185 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1",
            "value": 155745,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7180 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 - ns/op",
            "value": 155745,
            "unit": "ns/op",
            "extra": "7180 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7180 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7180 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0",
            "value": 202514,
            "unit": "ns/op\t   79002 B/op\t    3501 allocs/op",
            "extra": "6070 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 - ns/op",
            "value": 202514,
            "unit": "ns/op",
            "extra": "6070 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 - B/op",
            "value": 79002,
            "unit": "B/op",
            "extra": "6070 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "6070 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0",
            "value": 200599,
            "unit": "ns/op\t   79002 B/op\t    3501 allocs/op",
            "extra": "5839 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 - ns/op",
            "value": 200599,
            "unit": "ns/op",
            "extra": "5839 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 - B/op",
            "value": 79002,
            "unit": "B/op",
            "extra": "5839 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5839 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0",
            "value": 201206,
            "unit": "ns/op\t   79002 B/op\t    3501 allocs/op",
            "extra": "5959 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 - ns/op",
            "value": 201206,
            "unit": "ns/op",
            "extra": "5959 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 - B/op",
            "value": 79002,
            "unit": "B/op",
            "extra": "5959 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5959 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0",
            "value": 200215,
            "unit": "ns/op\t   79002 B/op\t    3501 allocs/op",
            "extra": "5821 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 - ns/op",
            "value": 200215,
            "unit": "ns/op",
            "extra": "5821 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 - B/op",
            "value": 79002,
            "unit": "B/op",
            "extra": "5821 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5821 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0",
            "value": 201810,
            "unit": "ns/op\t   79002 B/op\t    3501 allocs/op",
            "extra": "5908 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 - ns/op",
            "value": 201810,
            "unit": "ns/op",
            "extra": "5908 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 - B/op",
            "value": 79002,
            "unit": "B/op",
            "extra": "5908 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5908 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1",
            "value": 203232,
            "unit": "ns/op\t   79002 B/op\t    3501 allocs/op",
            "extra": "6020 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 - ns/op",
            "value": 203232,
            "unit": "ns/op",
            "extra": "6020 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 - B/op",
            "value": 79002,
            "unit": "B/op",
            "extra": "6020 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "6020 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1",
            "value": 202598,
            "unit": "ns/op\t   79002 B/op\t    3501 allocs/op",
            "extra": "5780 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 - ns/op",
            "value": 202598,
            "unit": "ns/op",
            "extra": "5780 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 - B/op",
            "value": 79002,
            "unit": "B/op",
            "extra": "5780 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5780 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1",
            "value": 204130,
            "unit": "ns/op\t   79002 B/op\t    3501 allocs/op",
            "extra": "5949 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 - ns/op",
            "value": 204130,
            "unit": "ns/op",
            "extra": "5949 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 - B/op",
            "value": 79002,
            "unit": "B/op",
            "extra": "5949 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5949 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1",
            "value": 203427,
            "unit": "ns/op\t   79002 B/op\t    3501 allocs/op",
            "extra": "5772 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 - ns/op",
            "value": 203427,
            "unit": "ns/op",
            "extra": "5772 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 - B/op",
            "value": 79002,
            "unit": "B/op",
            "extra": "5772 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5772 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1",
            "value": 201764,
            "unit": "ns/op\t   79002 B/op\t    3501 allocs/op",
            "extra": "5804 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 - ns/op",
            "value": 201764,
            "unit": "ns/op",
            "extra": "5804 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 - B/op",
            "value": 79002,
            "unit": "B/op",
            "extra": "5804 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5804 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0",
            "value": 217994,
            "unit": "ns/op\t   70322 B/op\t    3501 allocs/op",
            "extra": "5360 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 - ns/op",
            "value": 217994,
            "unit": "ns/op",
            "extra": "5360 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 - B/op",
            "value": 70322,
            "unit": "B/op",
            "extra": "5360 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5360 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0",
            "value": 217418,
            "unit": "ns/op\t   70322 B/op\t    3501 allocs/op",
            "extra": "5247 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 - ns/op",
            "value": 217418,
            "unit": "ns/op",
            "extra": "5247 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 - B/op",
            "value": 70322,
            "unit": "B/op",
            "extra": "5247 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5247 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0",
            "value": 215360,
            "unit": "ns/op\t   70322 B/op\t    3501 allocs/op",
            "extra": "5504 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 - ns/op",
            "value": 215360,
            "unit": "ns/op",
            "extra": "5504 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 - B/op",
            "value": 70322,
            "unit": "B/op",
            "extra": "5504 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5504 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0",
            "value": 217193,
            "unit": "ns/op\t   70322 B/op\t    3501 allocs/op",
            "extra": "5460 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 - ns/op",
            "value": 217193,
            "unit": "ns/op",
            "extra": "5460 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 - B/op",
            "value": 70322,
            "unit": "B/op",
            "extra": "5460 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5460 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0",
            "value": 217009,
            "unit": "ns/op\t   70322 B/op\t    3501 allocs/op",
            "extra": "5515 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 - ns/op",
            "value": 217009,
            "unit": "ns/op",
            "extra": "5515 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 - B/op",
            "value": 70322,
            "unit": "B/op",
            "extra": "5515 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5515 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1",
            "value": 216696,
            "unit": "ns/op\t   70322 B/op\t    3501 allocs/op",
            "extra": "5367 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 - ns/op",
            "value": 216696,
            "unit": "ns/op",
            "extra": "5367 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 - B/op",
            "value": 70322,
            "unit": "B/op",
            "extra": "5367 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5367 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1",
            "value": 217762,
            "unit": "ns/op\t   70322 B/op\t    3501 allocs/op",
            "extra": "5481 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 - ns/op",
            "value": 217762,
            "unit": "ns/op",
            "extra": "5481 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 - B/op",
            "value": 70322,
            "unit": "B/op",
            "extra": "5481 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5481 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1",
            "value": 217698,
            "unit": "ns/op\t   70322 B/op\t    3501 allocs/op",
            "extra": "5281 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 - ns/op",
            "value": 217698,
            "unit": "ns/op",
            "extra": "5281 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 - B/op",
            "value": 70322,
            "unit": "B/op",
            "extra": "5281 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5281 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1",
            "value": 216271,
            "unit": "ns/op\t   70322 B/op\t    3501 allocs/op",
            "extra": "5337 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 - ns/op",
            "value": 216271,
            "unit": "ns/op",
            "extra": "5337 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 - B/op",
            "value": 70322,
            "unit": "B/op",
            "extra": "5337 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5337 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1",
            "value": 219511,
            "unit": "ns/op\t   70322 B/op\t    3501 allocs/op",
            "extra": "5092 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 - ns/op",
            "value": 219511,
            "unit": "ns/op",
            "extra": "5092 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 - B/op",
            "value": 70322,
            "unit": "B/op",
            "extra": "5092 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5092 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0",
            "value": 5273,
            "unit": "ns/op\t    2040 B/op\t      55 allocs/op",
            "extra": "221684 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 - ns/op",
            "value": 5273,
            "unit": "ns/op",
            "extra": "221684 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 - B/op",
            "value": 2040,
            "unit": "B/op",
            "extra": "221684 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 - allocs/op",
            "value": 55,
            "unit": "allocs/op",
            "extra": "221684 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0",
            "value": 5295,
            "unit": "ns/op\t    2040 B/op\t      55 allocs/op",
            "extra": "227132 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 - ns/op",
            "value": 5295,
            "unit": "ns/op",
            "extra": "227132 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 - B/op",
            "value": 2040,
            "unit": "B/op",
            "extra": "227132 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 - allocs/op",
            "value": 55,
            "unit": "allocs/op",
            "extra": "227132 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0",
            "value": 5275,
            "unit": "ns/op\t    2040 B/op\t      55 allocs/op",
            "extra": "225901 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 - ns/op",
            "value": 5275,
            "unit": "ns/op",
            "extra": "225901 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 - B/op",
            "value": 2040,
            "unit": "B/op",
            "extra": "225901 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 - allocs/op",
            "value": 55,
            "unit": "allocs/op",
            "extra": "225901 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0",
            "value": 5334,
            "unit": "ns/op\t    2040 B/op\t      55 allocs/op",
            "extra": "227974 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 - ns/op",
            "value": 5334,
            "unit": "ns/op",
            "extra": "227974 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 - B/op",
            "value": 2040,
            "unit": "B/op",
            "extra": "227974 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 - allocs/op",
            "value": 55,
            "unit": "allocs/op",
            "extra": "227974 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0",
            "value": 5295,
            "unit": "ns/op\t    2040 B/op\t      55 allocs/op",
            "extra": "225130 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 - ns/op",
            "value": 5295,
            "unit": "ns/op",
            "extra": "225130 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 - B/op",
            "value": 2040,
            "unit": "B/op",
            "extra": "225130 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 - allocs/op",
            "value": 55,
            "unit": "allocs/op",
            "extra": "225130 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1",
            "value": 4956,
            "unit": "ns/op\t    2040 B/op\t      55 allocs/op",
            "extra": "236066 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 - ns/op",
            "value": 4956,
            "unit": "ns/op",
            "extra": "236066 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 - B/op",
            "value": 2040,
            "unit": "B/op",
            "extra": "236066 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 - allocs/op",
            "value": 55,
            "unit": "allocs/op",
            "extra": "236066 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1",
            "value": 4979,
            "unit": "ns/op\t    2040 B/op\t      55 allocs/op",
            "extra": "241288 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 - ns/op",
            "value": 4979,
            "unit": "ns/op",
            "extra": "241288 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 - B/op",
            "value": 2040,
            "unit": "B/op",
            "extra": "241288 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 - allocs/op",
            "value": 55,
            "unit": "allocs/op",
            "extra": "241288 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1",
            "value": 5007,
            "unit": "ns/op\t    2040 B/op\t      55 allocs/op",
            "extra": "236376 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 - ns/op",
            "value": 5007,
            "unit": "ns/op",
            "extra": "236376 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 - B/op",
            "value": 2040,
            "unit": "B/op",
            "extra": "236376 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 - allocs/op",
            "value": 55,
            "unit": "allocs/op",
            "extra": "236376 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1",
            "value": 5073,
            "unit": "ns/op\t    2040 B/op\t      55 allocs/op",
            "extra": "240090 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 - ns/op",
            "value": 5073,
            "unit": "ns/op",
            "extra": "240090 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 - B/op",
            "value": 2040,
            "unit": "B/op",
            "extra": "240090 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 - allocs/op",
            "value": 55,
            "unit": "allocs/op",
            "extra": "240090 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1",
            "value": 4987,
            "unit": "ns/op\t    2040 B/op\t      55 allocs/op",
            "extra": "245318 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 - ns/op",
            "value": 4987,
            "unit": "ns/op",
            "extra": "245318 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 - B/op",
            "value": 2040,
            "unit": "B/op",
            "extra": "245318 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 - allocs/op",
            "value": 55,
            "unit": "allocs/op",
            "extra": "245318 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0",
            "value": 2082949,
            "unit": "ns/op\t  157193 B/op\t   19512 allocs/op",
            "extra": "564 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 - ns/op",
            "value": 2082949,
            "unit": "ns/op",
            "extra": "564 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 - B/op",
            "value": 157193,
            "unit": "B/op",
            "extra": "564 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 - allocs/op",
            "value": 19512,
            "unit": "allocs/op",
            "extra": "564 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0",
            "value": 2096246,
            "unit": "ns/op\t  157193 B/op\t   19512 allocs/op",
            "extra": "565 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 - ns/op",
            "value": 2096246,
            "unit": "ns/op",
            "extra": "565 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 - B/op",
            "value": 157193,
            "unit": "B/op",
            "extra": "565 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 - allocs/op",
            "value": 19512,
            "unit": "allocs/op",
            "extra": "565 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0",
            "value": 2083970,
            "unit": "ns/op\t  157194 B/op\t   19512 allocs/op",
            "extra": "570 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 - ns/op",
            "value": 2083970,
            "unit": "ns/op",
            "extra": "570 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 - B/op",
            "value": 157194,
            "unit": "B/op",
            "extra": "570 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 - allocs/op",
            "value": 19512,
            "unit": "allocs/op",
            "extra": "570 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0",
            "value": 2093896,
            "unit": "ns/op\t  157195 B/op\t   19512 allocs/op",
            "extra": "571 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 - ns/op",
            "value": 2093896,
            "unit": "ns/op",
            "extra": "571 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 - B/op",
            "value": 157195,
            "unit": "B/op",
            "extra": "571 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 - allocs/op",
            "value": 19512,
            "unit": "allocs/op",
            "extra": "571 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0",
            "value": 2092614,
            "unit": "ns/op\t  157196 B/op\t   19512 allocs/op",
            "extra": "572 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 - ns/op",
            "value": 2092614,
            "unit": "ns/op",
            "extra": "572 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 - B/op",
            "value": 157196,
            "unit": "B/op",
            "extra": "572 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 - allocs/op",
            "value": 19512,
            "unit": "allocs/op",
            "extra": "572 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1",
            "value": 2011099,
            "unit": "ns/op\t  157195 B/op\t   19512 allocs/op",
            "extra": "612 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 - ns/op",
            "value": 2011099,
            "unit": "ns/op",
            "extra": "612 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 - B/op",
            "value": 157195,
            "unit": "B/op",
            "extra": "612 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 - allocs/op",
            "value": 19512,
            "unit": "allocs/op",
            "extra": "612 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1",
            "value": 1978210,
            "unit": "ns/op\t  157196 B/op\t   19512 allocs/op",
            "extra": "595 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 - ns/op",
            "value": 1978210,
            "unit": "ns/op",
            "extra": "595 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 - B/op",
            "value": 157196,
            "unit": "B/op",
            "extra": "595 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 - allocs/op",
            "value": 19512,
            "unit": "allocs/op",
            "extra": "595 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1",
            "value": 2016127,
            "unit": "ns/op\t  157195 B/op\t   19512 allocs/op",
            "extra": "574 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 - ns/op",
            "value": 2016127,
            "unit": "ns/op",
            "extra": "574 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 - B/op",
            "value": 157195,
            "unit": "B/op",
            "extra": "574 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 - allocs/op",
            "value": 19512,
            "unit": "allocs/op",
            "extra": "574 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1",
            "value": 1990806,
            "unit": "ns/op\t  157194 B/op\t   19512 allocs/op",
            "extra": "598 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 - ns/op",
            "value": 1990806,
            "unit": "ns/op",
            "extra": "598 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 - B/op",
            "value": 157194,
            "unit": "B/op",
            "extra": "598 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 - allocs/op",
            "value": 19512,
            "unit": "allocs/op",
            "extra": "598 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1",
            "value": 2011050,
            "unit": "ns/op\t  157194 B/op\t   19512 allocs/op",
            "extra": "591 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 - ns/op",
            "value": 2011050,
            "unit": "ns/op",
            "extra": "591 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 - B/op",
            "value": 157194,
            "unit": "B/op",
            "extra": "591 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 - allocs/op",
            "value": 19512,
            "unit": "allocs/op",
            "extra": "591 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0",
            "value": 3736895,
            "unit": "ns/op\t 1782931 B/op\t   49526 allocs/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 - ns/op",
            "value": 3736895,
            "unit": "ns/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 - B/op",
            "value": 1782931,
            "unit": "B/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 - allocs/op",
            "value": 49526,
            "unit": "allocs/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0",
            "value": 3739579,
            "unit": "ns/op\t 1782928 B/op\t   49526 allocs/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 - ns/op",
            "value": 3739579,
            "unit": "ns/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 - B/op",
            "value": 1782928,
            "unit": "B/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 - allocs/op",
            "value": 49526,
            "unit": "allocs/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0",
            "value": 3731146,
            "unit": "ns/op\t 1782924 B/op\t   49526 allocs/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 - ns/op",
            "value": 3731146,
            "unit": "ns/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 - B/op",
            "value": 1782924,
            "unit": "B/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 - allocs/op",
            "value": 49526,
            "unit": "allocs/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0",
            "value": 3733782,
            "unit": "ns/op\t 1782934 B/op\t   49526 allocs/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 - ns/op",
            "value": 3733782,
            "unit": "ns/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 - B/op",
            "value": 1782934,
            "unit": "B/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 - allocs/op",
            "value": 49526,
            "unit": "allocs/op",
            "extra": "321 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0",
            "value": 3716341,
            "unit": "ns/op\t 1782932 B/op\t   49526 allocs/op",
            "extra": "319 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 - ns/op",
            "value": 3716341,
            "unit": "ns/op",
            "extra": "319 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 - B/op",
            "value": 1782932,
            "unit": "B/op",
            "extra": "319 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 - allocs/op",
            "value": 49526,
            "unit": "allocs/op",
            "extra": "319 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1",
            "value": 3736293,
            "unit": "ns/op\t 1782933 B/op\t   49526 allocs/op",
            "extra": "318 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 - ns/op",
            "value": 3736293,
            "unit": "ns/op",
            "extra": "318 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 - B/op",
            "value": 1782933,
            "unit": "B/op",
            "extra": "318 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 - allocs/op",
            "value": 49526,
            "unit": "allocs/op",
            "extra": "318 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1",
            "value": 3757928,
            "unit": "ns/op\t 1782932 B/op\t   49526 allocs/op",
            "extra": "308 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 - ns/op",
            "value": 3757928,
            "unit": "ns/op",
            "extra": "308 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 - B/op",
            "value": 1782932,
            "unit": "B/op",
            "extra": "308 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 - allocs/op",
            "value": 49526,
            "unit": "allocs/op",
            "extra": "308 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1",
            "value": 3725442,
            "unit": "ns/op\t 1782932 B/op\t   49526 allocs/op",
            "extra": "320 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 - ns/op",
            "value": 3725442,
            "unit": "ns/op",
            "extra": "320 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 - B/op",
            "value": 1782932,
            "unit": "B/op",
            "extra": "320 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 - allocs/op",
            "value": 49526,
            "unit": "allocs/op",
            "extra": "320 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1",
            "value": 3757859,
            "unit": "ns/op\t 1782931 B/op\t   49526 allocs/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 - ns/op",
            "value": 3757859,
            "unit": "ns/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 - B/op",
            "value": 1782931,
            "unit": "B/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 - allocs/op",
            "value": 49526,
            "unit": "allocs/op",
            "extra": "322 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1",
            "value": 3714049,
            "unit": "ns/op\t 1782926 B/op\t   49526 allocs/op",
            "extra": "316 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 - ns/op",
            "value": 3714049,
            "unit": "ns/op",
            "extra": "316 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 - B/op",
            "value": 1782926,
            "unit": "B/op",
            "extra": "316 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 - allocs/op",
            "value": 49526,
            "unit": "allocs/op",
            "extra": "316 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0",
            "value": 5895068,
            "unit": "ns/op\t  266483 B/op\t   21519 allocs/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 - ns/op",
            "value": 5895068,
            "unit": "ns/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 - B/op",
            "value": 266483,
            "unit": "B/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 - allocs/op",
            "value": 21519,
            "unit": "allocs/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0",
            "value": 5924837,
            "unit": "ns/op\t  266481 B/op\t   21519 allocs/op",
            "extra": "201 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 - ns/op",
            "value": 5924837,
            "unit": "ns/op",
            "extra": "201 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 - B/op",
            "value": 266481,
            "unit": "B/op",
            "extra": "201 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 - allocs/op",
            "value": 21519,
            "unit": "allocs/op",
            "extra": "201 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0",
            "value": 5943467,
            "unit": "ns/op\t  266483 B/op\t   21519 allocs/op",
            "extra": "199 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 - ns/op",
            "value": 5943467,
            "unit": "ns/op",
            "extra": "199 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 - B/op",
            "value": 266483,
            "unit": "B/op",
            "extra": "199 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 - allocs/op",
            "value": 21519,
            "unit": "allocs/op",
            "extra": "199 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0",
            "value": 5884254,
            "unit": "ns/op\t  266481 B/op\t   21519 allocs/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 - ns/op",
            "value": 5884254,
            "unit": "ns/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 - B/op",
            "value": 266481,
            "unit": "B/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 - allocs/op",
            "value": 21519,
            "unit": "allocs/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0",
            "value": 5878683,
            "unit": "ns/op\t  266481 B/op\t   21519 allocs/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 - ns/op",
            "value": 5878683,
            "unit": "ns/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 - B/op",
            "value": 266481,
            "unit": "B/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 - allocs/op",
            "value": 21519,
            "unit": "allocs/op",
            "extra": "204 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1",
            "value": 5768189,
            "unit": "ns/op\t  266480 B/op\t   21519 allocs/op",
            "extra": "211 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 - ns/op",
            "value": 5768189,
            "unit": "ns/op",
            "extra": "211 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 - B/op",
            "value": 266480,
            "unit": "B/op",
            "extra": "211 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 - allocs/op",
            "value": 21519,
            "unit": "allocs/op",
            "extra": "211 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1",
            "value": 5698142,
            "unit": "ns/op\t  266482 B/op\t   21519 allocs/op",
            "extra": "210 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 - ns/op",
            "value": 5698142,
            "unit": "ns/op",
            "extra": "210 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 - B/op",
            "value": 266482,
            "unit": "B/op",
            "extra": "210 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 - allocs/op",
            "value": 21519,
            "unit": "allocs/op",
            "extra": "210 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1",
            "value": 5736151,
            "unit": "ns/op\t  266483 B/op\t   21519 allocs/op",
            "extra": "208 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 - ns/op",
            "value": 5736151,
            "unit": "ns/op",
            "extra": "208 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 - B/op",
            "value": 266483,
            "unit": "B/op",
            "extra": "208 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 - allocs/op",
            "value": 21519,
            "unit": "allocs/op",
            "extra": "208 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1",
            "value": 5731486,
            "unit": "ns/op\t  266482 B/op\t   21519 allocs/op",
            "extra": "208 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 - ns/op",
            "value": 5731486,
            "unit": "ns/op",
            "extra": "208 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 - B/op",
            "value": 266482,
            "unit": "B/op",
            "extra": "208 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 - allocs/op",
            "value": 21519,
            "unit": "allocs/op",
            "extra": "208 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1",
            "value": 5736580,
            "unit": "ns/op\t  266480 B/op\t   21519 allocs/op",
            "extra": "208 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 - ns/op",
            "value": 5736580,
            "unit": "ns/op",
            "extra": "208 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 - B/op",
            "value": 266480,
            "unit": "B/op",
            "extra": "208 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 - allocs/op",
            "value": 21519,
            "unit": "allocs/op",
            "extra": "208 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0",
            "value": 106166,
            "unit": "ns/op\t   53697 B/op\t    1055 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 - ns/op",
            "value": 106166,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 - B/op",
            "value": 53697,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 - allocs/op",
            "value": 1055,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0",
            "value": 106561,
            "unit": "ns/op\t   53697 B/op\t    1055 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 - ns/op",
            "value": 106561,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 - B/op",
            "value": 53697,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 - allocs/op",
            "value": 1055,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0",
            "value": 106851,
            "unit": "ns/op\t   53697 B/op\t    1055 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 - ns/op",
            "value": 106851,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 - B/op",
            "value": 53697,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 - allocs/op",
            "value": 1055,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0",
            "value": 107333,
            "unit": "ns/op\t   53698 B/op\t    1055 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 - ns/op",
            "value": 107333,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 - B/op",
            "value": 53698,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 - allocs/op",
            "value": 1055,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0",
            "value": 108730,
            "unit": "ns/op\t   53697 B/op\t    1055 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 - ns/op",
            "value": 108730,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 - B/op",
            "value": 53697,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 - allocs/op",
            "value": 1055,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1",
            "value": 108223,
            "unit": "ns/op\t   53697 B/op\t    1055 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 - ns/op",
            "value": 108223,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 - B/op",
            "value": 53697,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 - allocs/op",
            "value": 1055,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1",
            "value": 106360,
            "unit": "ns/op\t   53697 B/op\t    1055 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 - ns/op",
            "value": 106360,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 - B/op",
            "value": 53697,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 - allocs/op",
            "value": 1055,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1",
            "value": 105167,
            "unit": "ns/op\t   53697 B/op\t    1055 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 - ns/op",
            "value": 105167,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 - B/op",
            "value": 53697,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 - allocs/op",
            "value": 1055,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1",
            "value": 104102,
            "unit": "ns/op\t   53698 B/op\t    1055 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 - ns/op",
            "value": 104102,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 - B/op",
            "value": 53698,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 - allocs/op",
            "value": 1055,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1",
            "value": 104587,
            "unit": "ns/op\t   53698 B/op\t    1055 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 - ns/op",
            "value": 104587,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 - B/op",
            "value": 53698,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 - allocs/op",
            "value": 1055,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0",
            "value": 2996087,
            "unit": "ns/op\t 1782600 B/op\t   49520 allocs/op",
            "extra": "394 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 - ns/op",
            "value": 2996087,
            "unit": "ns/op",
            "extra": "394 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 - B/op",
            "value": 1782600,
            "unit": "B/op",
            "extra": "394 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 - allocs/op",
            "value": 49520,
            "unit": "allocs/op",
            "extra": "394 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0",
            "value": 2983786,
            "unit": "ns/op\t 1782601 B/op\t   49520 allocs/op",
            "extra": "398 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 - ns/op",
            "value": 2983786,
            "unit": "ns/op",
            "extra": "398 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 - B/op",
            "value": 1782601,
            "unit": "B/op",
            "extra": "398 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 - allocs/op",
            "value": 49520,
            "unit": "allocs/op",
            "extra": "398 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0",
            "value": 2997636,
            "unit": "ns/op\t 1782594 B/op\t   49520 allocs/op",
            "extra": "398 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 - ns/op",
            "value": 2997636,
            "unit": "ns/op",
            "extra": "398 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 - B/op",
            "value": 1782594,
            "unit": "B/op",
            "extra": "398 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 - allocs/op",
            "value": 49520,
            "unit": "allocs/op",
            "extra": "398 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0",
            "value": 3005034,
            "unit": "ns/op\t 1782592 B/op\t   49520 allocs/op",
            "extra": "396 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 - ns/op",
            "value": 3005034,
            "unit": "ns/op",
            "extra": "396 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 - B/op",
            "value": 1782592,
            "unit": "B/op",
            "extra": "396 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 - allocs/op",
            "value": 49520,
            "unit": "allocs/op",
            "extra": "396 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0",
            "value": 2999066,
            "unit": "ns/op\t 1782599 B/op\t   49520 allocs/op",
            "extra": "400 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 - ns/op",
            "value": 2999066,
            "unit": "ns/op",
            "extra": "400 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 - B/op",
            "value": 1782599,
            "unit": "B/op",
            "extra": "400 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 - allocs/op",
            "value": 49520,
            "unit": "allocs/op",
            "extra": "400 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1",
            "value": 2995632,
            "unit": "ns/op\t 1782601 B/op\t   49520 allocs/op",
            "extra": "402 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 - ns/op",
            "value": 2995632,
            "unit": "ns/op",
            "extra": "402 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 - B/op",
            "value": 1782601,
            "unit": "B/op",
            "extra": "402 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 - allocs/op",
            "value": 49520,
            "unit": "allocs/op",
            "extra": "402 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1",
            "value": 2994083,
            "unit": "ns/op\t 1782604 B/op\t   49520 allocs/op",
            "extra": "403 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 - ns/op",
            "value": 2994083,
            "unit": "ns/op",
            "extra": "403 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 - B/op",
            "value": 1782604,
            "unit": "B/op",
            "extra": "403 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 - allocs/op",
            "value": 49520,
            "unit": "allocs/op",
            "extra": "403 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1",
            "value": 3063914,
            "unit": "ns/op\t 1782602 B/op\t   49520 allocs/op",
            "extra": "402 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 - ns/op",
            "value": 3063914,
            "unit": "ns/op",
            "extra": "402 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 - B/op",
            "value": 1782602,
            "unit": "B/op",
            "extra": "402 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 - allocs/op",
            "value": 49520,
            "unit": "allocs/op",
            "extra": "402 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1",
            "value": 3026156,
            "unit": "ns/op\t 1782600 B/op\t   49520 allocs/op",
            "extra": "396 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 - ns/op",
            "value": 3026156,
            "unit": "ns/op",
            "extra": "396 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 - B/op",
            "value": 1782600,
            "unit": "B/op",
            "extra": "396 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 - allocs/op",
            "value": 49520,
            "unit": "allocs/op",
            "extra": "396 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1",
            "value": 3015568,
            "unit": "ns/op\t 1782604 B/op\t   49520 allocs/op",
            "extra": "399 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 - ns/op",
            "value": 3015568,
            "unit": "ns/op",
            "extra": "399 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 - B/op",
            "value": 1782604,
            "unit": "B/op",
            "extra": "399 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 - allocs/op",
            "value": 49520,
            "unit": "allocs/op",
            "extra": "399 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0",
            "value": 87687,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13618 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 - ns/op",
            "value": 87687,
            "unit": "ns/op",
            "extra": "13618 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13618 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13618 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0",
            "value": 88091,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13687 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 - ns/op",
            "value": 88091,
            "unit": "ns/op",
            "extra": "13687 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13687 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13687 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0",
            "value": 87516,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13626 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 - ns/op",
            "value": 87516,
            "unit": "ns/op",
            "extra": "13626 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13626 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13626 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0",
            "value": 87796,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13678 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 - ns/op",
            "value": 87796,
            "unit": "ns/op",
            "extra": "13678 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13678 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13678 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0",
            "value": 88220,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13699 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 - ns/op",
            "value": 88220,
            "unit": "ns/op",
            "extra": "13699 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13699 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13699 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1",
            "value": 86556,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13833 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 - ns/op",
            "value": 86556,
            "unit": "ns/op",
            "extra": "13833 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13833 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13833 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1",
            "value": 86472,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13900 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 - ns/op",
            "value": 86472,
            "unit": "ns/op",
            "extra": "13900 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13900 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13900 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1",
            "value": 86586,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13654 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 - ns/op",
            "value": 86586,
            "unit": "ns/op",
            "extra": "13654 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13654 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13654 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1",
            "value": 86159,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13826 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 - ns/op",
            "value": 86159,
            "unit": "ns/op",
            "extra": "13826 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13826 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13826 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1",
            "value": 86443,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13891 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 - ns/op",
            "value": 86443,
            "unit": "ns/op",
            "extra": "13891 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13891 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13891 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0",
            "value": 88172,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13513 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 - ns/op",
            "value": 88172,
            "unit": "ns/op",
            "extra": "13513 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13513 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13513 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0",
            "value": 88894,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13519 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 - ns/op",
            "value": 88894,
            "unit": "ns/op",
            "extra": "13519 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13519 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13519 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0",
            "value": 88662,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13503 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 - ns/op",
            "value": 88662,
            "unit": "ns/op",
            "extra": "13503 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13503 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13503 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0",
            "value": 88994,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13531 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 - ns/op",
            "value": 88994,
            "unit": "ns/op",
            "extra": "13531 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13531 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13531 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0",
            "value": 88316,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13497 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 - ns/op",
            "value": 88316,
            "unit": "ns/op",
            "extra": "13497 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13497 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13497 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1",
            "value": 89227,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13891 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 - ns/op",
            "value": 89227,
            "unit": "ns/op",
            "extra": "13891 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13891 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13891 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1",
            "value": 86605,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13870 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 - ns/op",
            "value": 86605,
            "unit": "ns/op",
            "extra": "13870 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13870 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13870 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1",
            "value": 86284,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13813 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 - ns/op",
            "value": 86284,
            "unit": "ns/op",
            "extra": "13813 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13813 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13813 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1",
            "value": 86401,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13946 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 - ns/op",
            "value": 86401,
            "unit": "ns/op",
            "extra": "13946 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13946 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13946 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1",
            "value": 86251,
            "unit": "ns/op\t   40897 B/op\t     855 allocs/op",
            "extra": "13929 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 - ns/op",
            "value": 86251,
            "unit": "ns/op",
            "extra": "13929 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 - B/op",
            "value": 40897,
            "unit": "B/op",
            "extra": "13929 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 - allocs/op",
            "value": 855,
            "unit": "allocs/op",
            "extra": "13929 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0",
            "value": 8962,
            "unit": "ns/op\t     480 B/op\t      10 allocs/op",
            "extra": "131683 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 - ns/op",
            "value": 8962,
            "unit": "ns/op",
            "extra": "131683 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 - B/op",
            "value": 480,
            "unit": "B/op",
            "extra": "131683 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "131683 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0",
            "value": 8918,
            "unit": "ns/op\t     480 B/op\t      10 allocs/op",
            "extra": "132900 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 - ns/op",
            "value": 8918,
            "unit": "ns/op",
            "extra": "132900 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 - B/op",
            "value": 480,
            "unit": "B/op",
            "extra": "132900 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "132900 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0",
            "value": 8942,
            "unit": "ns/op\t     480 B/op\t      10 allocs/op",
            "extra": "133658 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 - ns/op",
            "value": 8942,
            "unit": "ns/op",
            "extra": "133658 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 - B/op",
            "value": 480,
            "unit": "B/op",
            "extra": "133658 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "133658 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0",
            "value": 8915,
            "unit": "ns/op\t     480 B/op\t      10 allocs/op",
            "extra": "133594 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 - ns/op",
            "value": 8915,
            "unit": "ns/op",
            "extra": "133594 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 - B/op",
            "value": 480,
            "unit": "B/op",
            "extra": "133594 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "133594 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0",
            "value": 8913,
            "unit": "ns/op\t     480 B/op\t      10 allocs/op",
            "extra": "134695 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 - ns/op",
            "value": 8913,
            "unit": "ns/op",
            "extra": "134695 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 - B/op",
            "value": 480,
            "unit": "B/op",
            "extra": "134695 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "134695 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1",
            "value": 9001,
            "unit": "ns/op\t     480 B/op\t      10 allocs/op",
            "extra": "131997 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 - ns/op",
            "value": 9001,
            "unit": "ns/op",
            "extra": "131997 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 - B/op",
            "value": 480,
            "unit": "B/op",
            "extra": "131997 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "131997 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1",
            "value": 9050,
            "unit": "ns/op\t     480 B/op\t      10 allocs/op",
            "extra": "133298 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 - ns/op",
            "value": 9050,
            "unit": "ns/op",
            "extra": "133298 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 - B/op",
            "value": 480,
            "unit": "B/op",
            "extra": "133298 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "133298 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1",
            "value": 9001,
            "unit": "ns/op\t     480 B/op\t      10 allocs/op",
            "extra": "133410 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 - ns/op",
            "value": 9001,
            "unit": "ns/op",
            "extra": "133410 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 - B/op",
            "value": 480,
            "unit": "B/op",
            "extra": "133410 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "133410 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1",
            "value": 9055,
            "unit": "ns/op\t     480 B/op\t      10 allocs/op",
            "extra": "133959 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 - ns/op",
            "value": 9055,
            "unit": "ns/op",
            "extra": "133959 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 - B/op",
            "value": 480,
            "unit": "B/op",
            "extra": "133959 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "133959 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1",
            "value": 9048,
            "unit": "ns/op\t     480 B/op\t      10 allocs/op",
            "extra": "133405 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 - ns/op",
            "value": 9048,
            "unit": "ns/op",
            "extra": "133405 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 - B/op",
            "value": 480,
            "unit": "B/op",
            "extra": "133405 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 - allocs/op",
            "value": 10,
            "unit": "allocs/op",
            "extra": "133405 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0",
            "value": 2686,
            "unit": "ns/op\t     640 B/op\t      17 allocs/op",
            "extra": "385279 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 - ns/op",
            "value": 2686,
            "unit": "ns/op",
            "extra": "385279 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 - B/op",
            "value": 640,
            "unit": "B/op",
            "extra": "385279 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "385279 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0",
            "value": 2665,
            "unit": "ns/op\t     640 B/op\t      17 allocs/op",
            "extra": "394228 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 - ns/op",
            "value": 2665,
            "unit": "ns/op",
            "extra": "394228 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 - B/op",
            "value": 640,
            "unit": "B/op",
            "extra": "394228 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "394228 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0",
            "value": 2650,
            "unit": "ns/op\t     640 B/op\t      17 allocs/op",
            "extra": "447711 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 - ns/op",
            "value": 2650,
            "unit": "ns/op",
            "extra": "447711 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 - B/op",
            "value": 640,
            "unit": "B/op",
            "extra": "447711 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "447711 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0",
            "value": 2653,
            "unit": "ns/op\t     640 B/op\t      17 allocs/op",
            "extra": "444650 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 - ns/op",
            "value": 2653,
            "unit": "ns/op",
            "extra": "444650 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 - B/op",
            "value": 640,
            "unit": "B/op",
            "extra": "444650 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "444650 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0",
            "value": 2660,
            "unit": "ns/op\t     640 B/op\t      17 allocs/op",
            "extra": "453333 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 - ns/op",
            "value": 2660,
            "unit": "ns/op",
            "extra": "453333 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 - B/op",
            "value": 640,
            "unit": "B/op",
            "extra": "453333 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "453333 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1",
            "value": 1394,
            "unit": "ns/op\t     560 B/op\t       7 allocs/op",
            "extra": "805978 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 - ns/op",
            "value": 1394,
            "unit": "ns/op",
            "extra": "805978 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "805978 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "805978 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1",
            "value": 1398,
            "unit": "ns/op\t     560 B/op\t       7 allocs/op",
            "extra": "795014 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 - ns/op",
            "value": 1398,
            "unit": "ns/op",
            "extra": "795014 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "795014 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "795014 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1",
            "value": 1392,
            "unit": "ns/op\t     560 B/op\t       7 allocs/op",
            "extra": "817506 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 - ns/op",
            "value": 1392,
            "unit": "ns/op",
            "extra": "817506 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "817506 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "817506 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1",
            "value": 1399,
            "unit": "ns/op\t     560 B/op\t       7 allocs/op",
            "extra": "812610 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 - ns/op",
            "value": 1399,
            "unit": "ns/op",
            "extra": "812610 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "812610 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "812610 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1",
            "value": 1392,
            "unit": "ns/op\t     560 B/op\t       7 allocs/op",
            "extra": "803698 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 - ns/op",
            "value": 1392,
            "unit": "ns/op",
            "extra": "803698 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "803698 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "803698 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0",
            "value": 3618,
            "unit": "ns/op\t    1112 B/op\t      34 allocs/op",
            "extra": "317169 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 - ns/op",
            "value": 3618,
            "unit": "ns/op",
            "extra": "317169 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "317169 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "317169 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0",
            "value": 3623,
            "unit": "ns/op\t    1112 B/op\t      34 allocs/op",
            "extra": "322653 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 - ns/op",
            "value": 3623,
            "unit": "ns/op",
            "extra": "322653 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "322653 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "322653 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0",
            "value": 3744,
            "unit": "ns/op\t    1112 B/op\t      34 allocs/op",
            "extra": "322934 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 - ns/op",
            "value": 3744,
            "unit": "ns/op",
            "extra": "322934 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "322934 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "322934 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0",
            "value": 3617,
            "unit": "ns/op\t    1112 B/op\t      34 allocs/op",
            "extra": "326234 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 - ns/op",
            "value": 3617,
            "unit": "ns/op",
            "extra": "326234 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "326234 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "326234 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0",
            "value": 3673,
            "unit": "ns/op\t    1112 B/op\t      34 allocs/op",
            "extra": "322821 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 - ns/op",
            "value": 3673,
            "unit": "ns/op",
            "extra": "322821 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "322821 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "322821 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1",
            "value": 3844,
            "unit": "ns/op\t    1112 B/op\t      34 allocs/op",
            "extra": "308421 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 - ns/op",
            "value": 3844,
            "unit": "ns/op",
            "extra": "308421 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "308421 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "308421 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1",
            "value": 3809,
            "unit": "ns/op\t    1112 B/op\t      34 allocs/op",
            "extra": "313830 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 - ns/op",
            "value": 3809,
            "unit": "ns/op",
            "extra": "313830 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "313830 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "313830 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1",
            "value": 3736,
            "unit": "ns/op\t    1112 B/op\t      34 allocs/op",
            "extra": "316017 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 - ns/op",
            "value": 3736,
            "unit": "ns/op",
            "extra": "316017 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "316017 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "316017 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1",
            "value": 3748,
            "unit": "ns/op\t    1112 B/op\t      34 allocs/op",
            "extra": "316442 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 - ns/op",
            "value": 3748,
            "unit": "ns/op",
            "extra": "316442 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "316442 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "316442 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1",
            "value": 3743,
            "unit": "ns/op\t    1112 B/op\t      34 allocs/op",
            "extra": "302745 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 - ns/op",
            "value": 3743,
            "unit": "ns/op",
            "extra": "302745 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "302745 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "302745 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0",
            "value": 155.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7501658 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 - ns/op",
            "value": 155.6,
            "unit": "ns/op",
            "extra": "7501658 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7501658 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7501658 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0",
            "value": 157.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7875397 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 - ns/op",
            "value": 157.1,
            "unit": "ns/op",
            "extra": "7875397 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7875397 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7875397 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0",
            "value": 156.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7815556 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 - ns/op",
            "value": 156.2,
            "unit": "ns/op",
            "extra": "7815556 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7815556 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7815556 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0",
            "value": 152.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7812962 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 - ns/op",
            "value": 152.9,
            "unit": "ns/op",
            "extra": "7812962 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7812962 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7812962 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0",
            "value": 153.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7754881 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 - ns/op",
            "value": 153.5,
            "unit": "ns/op",
            "extra": "7754881 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7754881 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7754881 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1",
            "value": 152.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7813900 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 - ns/op",
            "value": 152.9,
            "unit": "ns/op",
            "extra": "7813900 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7813900 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7813900 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1",
            "value": 153.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7868190 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 - ns/op",
            "value": 153.3,
            "unit": "ns/op",
            "extra": "7868190 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7868190 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7868190 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1",
            "value": 153.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7811566 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 - ns/op",
            "value": 153.7,
            "unit": "ns/op",
            "extra": "7811566 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7811566 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7811566 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1",
            "value": 153.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7825791 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 - ns/op",
            "value": 153.6,
            "unit": "ns/op",
            "extra": "7825791 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7825791 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7825791 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1",
            "value": 153.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7816776 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 - ns/op",
            "value": 153.5,
            "unit": "ns/op",
            "extra": "7816776 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7816776 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7816776 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0",
            "value": 83.49,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14425888 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 - ns/op",
            "value": 83.49,
            "unit": "ns/op",
            "extra": "14425888 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14425888 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14425888 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0",
            "value": 83.11,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14395502 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 - ns/op",
            "value": 83.11,
            "unit": "ns/op",
            "extra": "14395502 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14395502 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14395502 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0",
            "value": 83.13,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14483920 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 - ns/op",
            "value": 83.13,
            "unit": "ns/op",
            "extra": "14483920 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14483920 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14483920 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0",
            "value": 82.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14462254 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 - ns/op",
            "value": 82.62,
            "unit": "ns/op",
            "extra": "14462254 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14462254 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14462254 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0",
            "value": 82.66,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14454368 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 - ns/op",
            "value": 82.66,
            "unit": "ns/op",
            "extra": "14454368 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14454368 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14454368 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1",
            "value": 83.04,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14557254 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 - ns/op",
            "value": 83.04,
            "unit": "ns/op",
            "extra": "14557254 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14557254 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14557254 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1",
            "value": 83.73,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14508884 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 - ns/op",
            "value": 83.73,
            "unit": "ns/op",
            "extra": "14508884 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14508884 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14508884 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1",
            "value": 82.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14508552 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 - ns/op",
            "value": 82.67,
            "unit": "ns/op",
            "extra": "14508552 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14508552 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14508552 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1",
            "value": 82.56,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13617740 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 - ns/op",
            "value": 82.56,
            "unit": "ns/op",
            "extra": "13617740 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13617740 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13617740 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1",
            "value": 83.13,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14445940 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 - ns/op",
            "value": 83.13,
            "unit": "ns/op",
            "extra": "14445940 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14445940 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14445940 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0",
            "value": 82.68,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14608970 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 - ns/op",
            "value": 82.68,
            "unit": "ns/op",
            "extra": "14608970 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14608970 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14608970 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0",
            "value": 82.69,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14621994 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 - ns/op",
            "value": 82.69,
            "unit": "ns/op",
            "extra": "14621994 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14621994 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14621994 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0",
            "value": 82.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14435115 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 - ns/op",
            "value": 82.3,
            "unit": "ns/op",
            "extra": "14435115 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14435115 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14435115 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0",
            "value": 82.31,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14525028 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 - ns/op",
            "value": 82.31,
            "unit": "ns/op",
            "extra": "14525028 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14525028 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14525028 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0",
            "value": 82.05,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14649082 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 - ns/op",
            "value": 82.05,
            "unit": "ns/op",
            "extra": "14649082 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14649082 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14649082 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1",
            "value": 82.22,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14682303 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 - ns/op",
            "value": 82.22,
            "unit": "ns/op",
            "extra": "14682303 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14682303 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14682303 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1",
            "value": 82.31,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14584908 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 - ns/op",
            "value": 82.31,
            "unit": "ns/op",
            "extra": "14584908 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14584908 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14584908 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1",
            "value": 82.45,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14640684 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 - ns/op",
            "value": 82.45,
            "unit": "ns/op",
            "extra": "14640684 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14640684 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14640684 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1",
            "value": 82.86,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13312414 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 - ns/op",
            "value": 82.86,
            "unit": "ns/op",
            "extra": "13312414 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13312414 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13312414 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1",
            "value": 82.14,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14634301 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 - ns/op",
            "value": 82.14,
            "unit": "ns/op",
            "extra": "14634301 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14634301 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14634301 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0",
            "value": 88.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13370701 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 - ns/op",
            "value": 88.78,
            "unit": "ns/op",
            "extra": "13370701 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13370701 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13370701 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0",
            "value": 88.81,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13513316 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 - ns/op",
            "value": 88.81,
            "unit": "ns/op",
            "extra": "13513316 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13513316 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13513316 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0",
            "value": 88.82,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13539759 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 - ns/op",
            "value": 88.82,
            "unit": "ns/op",
            "extra": "13539759 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13539759 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13539759 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0",
            "value": 89.97,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13531107 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 - ns/op",
            "value": 89.97,
            "unit": "ns/op",
            "extra": "13531107 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13531107 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13531107 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0",
            "value": 88.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13581109 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 - ns/op",
            "value": 88.62,
            "unit": "ns/op",
            "extra": "13581109 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13581109 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13581109 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1",
            "value": 88.75,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13500841 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 - ns/op",
            "value": 88.75,
            "unit": "ns/op",
            "extra": "13500841 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13500841 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13500841 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1",
            "value": 88.73,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13568169 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 - ns/op",
            "value": 88.73,
            "unit": "ns/op",
            "extra": "13568169 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13568169 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13568169 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1",
            "value": 88.68,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13564748 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 - ns/op",
            "value": 88.68,
            "unit": "ns/op",
            "extra": "13564748 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13564748 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13564748 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1",
            "value": 88.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13571179 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 - ns/op",
            "value": 88.7,
            "unit": "ns/op",
            "extra": "13571179 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13571179 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13571179 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1",
            "value": 88.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13549791 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 - ns/op",
            "value": 88.83,
            "unit": "ns/op",
            "extra": "13549791 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13549791 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13549791 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0",
            "value": 124.6,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "9562075 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 - ns/op",
            "value": 124.6,
            "unit": "ns/op",
            "extra": "9562075 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "9562075 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9562075 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0",
            "value": 124.4,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "9595293 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 - ns/op",
            "value": 124.4,
            "unit": "ns/op",
            "extra": "9595293 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "9595293 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9595293 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0",
            "value": 123.8,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "9495420 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 - ns/op",
            "value": 123.8,
            "unit": "ns/op",
            "extra": "9495420 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "9495420 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9495420 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0",
            "value": 124.7,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "9574527 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 - ns/op",
            "value": 124.7,
            "unit": "ns/op",
            "extra": "9574527 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "9574527 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9574527 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0",
            "value": 124.2,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "9538908 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 - ns/op",
            "value": 124.2,
            "unit": "ns/op",
            "extra": "9538908 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "9538908 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9538908 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1",
            "value": 124.2,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "9590568 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 - ns/op",
            "value": 124.2,
            "unit": "ns/op",
            "extra": "9590568 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "9590568 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9590568 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1",
            "value": 124.3,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "9572150 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 - ns/op",
            "value": 124.3,
            "unit": "ns/op",
            "extra": "9572150 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "9572150 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9572150 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1",
            "value": 124.1,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "9597690 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 - ns/op",
            "value": 124.1,
            "unit": "ns/op",
            "extra": "9597690 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "9597690 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9597690 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1",
            "value": 124,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "9571884 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 - ns/op",
            "value": 124,
            "unit": "ns/op",
            "extra": "9571884 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "9571884 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9571884 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1",
            "value": 124.3,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "9524446 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 - ns/op",
            "value": 124.3,
            "unit": "ns/op",
            "extra": "9524446 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "9524446 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9524446 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0",
            "value": 94.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12795398 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 - ns/op",
            "value": 94.2,
            "unit": "ns/op",
            "extra": "12795398 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12795398 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12795398 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0",
            "value": 94.01,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12852854 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 - ns/op",
            "value": 94.01,
            "unit": "ns/op",
            "extra": "12852854 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12852854 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12852854 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0",
            "value": 96.41,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12733458 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 - ns/op",
            "value": 96.41,
            "unit": "ns/op",
            "extra": "12733458 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12733458 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12733458 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0",
            "value": 93.81,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12749264 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 - ns/op",
            "value": 93.81,
            "unit": "ns/op",
            "extra": "12749264 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12749264 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12749264 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0",
            "value": 93.47,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12841438 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 - ns/op",
            "value": 93.47,
            "unit": "ns/op",
            "extra": "12841438 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12841438 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12841438 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1",
            "value": 93.44,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12749841 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 - ns/op",
            "value": 93.44,
            "unit": "ns/op",
            "extra": "12749841 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12749841 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12749841 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1",
            "value": 93.87,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12823022 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 - ns/op",
            "value": 93.87,
            "unit": "ns/op",
            "extra": "12823022 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12823022 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12823022 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1",
            "value": 93.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12814204 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 - ns/op",
            "value": 93.62,
            "unit": "ns/op",
            "extra": "12814204 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12814204 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12814204 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1",
            "value": 94.19,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12847838 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 - ns/op",
            "value": 94.19,
            "unit": "ns/op",
            "extra": "12847838 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12847838 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12847838 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1",
            "value": 93.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12818265 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 - ns/op",
            "value": 93.62,
            "unit": "ns/op",
            "extra": "12818265 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12818265 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12818265 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0",
            "value": 139.1,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "8689231 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 - ns/op",
            "value": 139.1,
            "unit": "ns/op",
            "extra": "8689231 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "8689231 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8689231 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0",
            "value": 137.2,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "8731126 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 - ns/op",
            "value": 137.2,
            "unit": "ns/op",
            "extra": "8731126 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "8731126 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8731126 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0",
            "value": 136.7,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "8710275 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 - ns/op",
            "value": 136.7,
            "unit": "ns/op",
            "extra": "8710275 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "8710275 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8710275 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0",
            "value": 137.4,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "8703916 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 - ns/op",
            "value": 137.4,
            "unit": "ns/op",
            "extra": "8703916 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "8703916 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8703916 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0",
            "value": 136.8,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "8734682 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 - ns/op",
            "value": 136.8,
            "unit": "ns/op",
            "extra": "8734682 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "8734682 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8734682 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1",
            "value": 137.6,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "8687752 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 - ns/op",
            "value": 137.6,
            "unit": "ns/op",
            "extra": "8687752 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "8687752 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8687752 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1",
            "value": 137,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "8710214 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 - ns/op",
            "value": 137,
            "unit": "ns/op",
            "extra": "8710214 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "8710214 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8710214 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1",
            "value": 137,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "8726109 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 - ns/op",
            "value": 137,
            "unit": "ns/op",
            "extra": "8726109 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "8726109 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8726109 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1",
            "value": 137.6,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "8698602 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 - ns/op",
            "value": 137.6,
            "unit": "ns/op",
            "extra": "8698602 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "8698602 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8698602 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1",
            "value": 137.2,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "7912719 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 - ns/op",
            "value": 137.2,
            "unit": "ns/op",
            "extra": "7912719 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "7912719 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7912719 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0",
            "value": 99.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12062649 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 - ns/op",
            "value": 99.4,
            "unit": "ns/op",
            "extra": "12062649 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12062649 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12062649 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0",
            "value": 100.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12064495 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 - ns/op",
            "value": 100.6,
            "unit": "ns/op",
            "extra": "12064495 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12064495 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12064495 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0",
            "value": 99.88,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10208731 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 - ns/op",
            "value": 99.88,
            "unit": "ns/op",
            "extra": "10208731 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10208731 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10208731 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0",
            "value": 99.54,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12032980 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 - ns/op",
            "value": 99.54,
            "unit": "ns/op",
            "extra": "12032980 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12032980 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12032980 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0",
            "value": 100.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11943907 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 - ns/op",
            "value": 100.8,
            "unit": "ns/op",
            "extra": "11943907 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11943907 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11943907 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1",
            "value": 100.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11845520 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 - ns/op",
            "value": 100.8,
            "unit": "ns/op",
            "extra": "11845520 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11845520 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11845520 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1",
            "value": 100.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11995868 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 - ns/op",
            "value": 100.9,
            "unit": "ns/op",
            "extra": "11995868 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11995868 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11995868 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1",
            "value": 100.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11950438 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 - ns/op",
            "value": 100.7,
            "unit": "ns/op",
            "extra": "11950438 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11950438 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11950438 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1",
            "value": 100.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11849074 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 - ns/op",
            "value": 100.8,
            "unit": "ns/op",
            "extra": "11849074 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11849074 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11849074 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1",
            "value": 103.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11835244 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 - ns/op",
            "value": 103.4,
            "unit": "ns/op",
            "extra": "11835244 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11835244 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11835244 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0",
            "value": 154.7,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "7792447 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 - ns/op",
            "value": 154.7,
            "unit": "ns/op",
            "extra": "7792447 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "7792447 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7792447 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0",
            "value": 154.9,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "7701483 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 - ns/op",
            "value": 154.9,
            "unit": "ns/op",
            "extra": "7701483 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "7701483 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7701483 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0",
            "value": 155.4,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "7683662 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 - ns/op",
            "value": 155.4,
            "unit": "ns/op",
            "extra": "7683662 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "7683662 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7683662 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0",
            "value": 154.7,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "7692223 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 - ns/op",
            "value": 154.7,
            "unit": "ns/op",
            "extra": "7692223 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "7692223 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7692223 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0",
            "value": 154.9,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "7658227 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 - ns/op",
            "value": 154.9,
            "unit": "ns/op",
            "extra": "7658227 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "7658227 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7658227 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1",
            "value": 154.6,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "7712180 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 - ns/op",
            "value": 154.6,
            "unit": "ns/op",
            "extra": "7712180 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "7712180 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7712180 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1",
            "value": 154.6,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "7650351 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 - ns/op",
            "value": 154.6,
            "unit": "ns/op",
            "extra": "7650351 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "7650351 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7650351 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1",
            "value": 154.4,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "7701404 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 - ns/op",
            "value": 154.4,
            "unit": "ns/op",
            "extra": "7701404 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "7701404 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7701404 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1",
            "value": 154.6,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "7759485 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 - ns/op",
            "value": 154.6,
            "unit": "ns/op",
            "extra": "7759485 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "7759485 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7759485 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1",
            "value": 160.6,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "7767367 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 - ns/op",
            "value": 160.6,
            "unit": "ns/op",
            "extra": "7767367 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "7767367 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7767367 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0",
            "value": 116,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10361072 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 - ns/op",
            "value": 116,
            "unit": "ns/op",
            "extra": "10361072 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10361072 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10361072 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0",
            "value": 115.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10376265 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 - ns/op",
            "value": 115.6,
            "unit": "ns/op",
            "extra": "10376265 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10376265 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10376265 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0",
            "value": 115.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10359692 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 - ns/op",
            "value": 115.6,
            "unit": "ns/op",
            "extra": "10359692 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10359692 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10359692 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0",
            "value": 115.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10337361 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 - ns/op",
            "value": 115.9,
            "unit": "ns/op",
            "extra": "10337361 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10337361 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10337361 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0",
            "value": 115.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10263349 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 - ns/op",
            "value": 115.7,
            "unit": "ns/op",
            "extra": "10263349 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10263349 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10263349 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1",
            "value": 116,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10422579 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 - ns/op",
            "value": 116,
            "unit": "ns/op",
            "extra": "10422579 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10422579 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10422579 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1",
            "value": 115.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10378720 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 - ns/op",
            "value": 115.5,
            "unit": "ns/op",
            "extra": "10378720 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10378720 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10378720 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1",
            "value": 115.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10367838 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 - ns/op",
            "value": 115.7,
            "unit": "ns/op",
            "extra": "10367838 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10367838 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10367838 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1",
            "value": 115.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10404327 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 - ns/op",
            "value": 115.5,
            "unit": "ns/op",
            "extra": "10404327 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10404327 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10404327 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1",
            "value": 118.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10353735 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 - ns/op",
            "value": 118.9,
            "unit": "ns/op",
            "extra": "10353735 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10353735 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10353735 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0",
            "value": 172.2,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "6904114 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 - ns/op",
            "value": 172.2,
            "unit": "ns/op",
            "extra": "6904114 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "6904114 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "6904114 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0",
            "value": 173.3,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "6931413 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 - ns/op",
            "value": 173.3,
            "unit": "ns/op",
            "extra": "6931413 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "6931413 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "6931413 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0",
            "value": 172.1,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "6910794 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 - ns/op",
            "value": 172.1,
            "unit": "ns/op",
            "extra": "6910794 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "6910794 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "6910794 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0",
            "value": 172.3,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "7024356 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 - ns/op",
            "value": 172.3,
            "unit": "ns/op",
            "extra": "7024356 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "7024356 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7024356 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0",
            "value": 173.1,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "6837564 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 - ns/op",
            "value": 173.1,
            "unit": "ns/op",
            "extra": "6837564 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "6837564 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "6837564 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1",
            "value": 172.1,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "6927885 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 - ns/op",
            "value": 172.1,
            "unit": "ns/op",
            "extra": "6927885 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "6927885 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "6927885 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1",
            "value": 171.6,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "6919380 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 - ns/op",
            "value": 171.6,
            "unit": "ns/op",
            "extra": "6919380 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "6919380 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "6919380 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1",
            "value": 172.4,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "6965632 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 - ns/op",
            "value": 172.4,
            "unit": "ns/op",
            "extra": "6965632 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "6965632 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "6965632 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1",
            "value": 172.1,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "6934784 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 - ns/op",
            "value": 172.1,
            "unit": "ns/op",
            "extra": "6934784 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "6934784 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "6934784 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1",
            "value": 172.1,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "6946594 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 - ns/op",
            "value": 172.1,
            "unit": "ns/op",
            "extra": "6946594 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "6946594 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "6946594 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0",
            "value": 221,
            "unit": "ns/op\t      88 B/op\t       2 allocs/op",
            "extra": "5409986 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 - ns/op",
            "value": 221,
            "unit": "ns/op",
            "extra": "5409986 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 - B/op",
            "value": 88,
            "unit": "B/op",
            "extra": "5409986 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5409986 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0",
            "value": 227.6,
            "unit": "ns/op\t      88 B/op\t       2 allocs/op",
            "extra": "5397913 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 - ns/op",
            "value": 227.6,
            "unit": "ns/op",
            "extra": "5397913 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 - B/op",
            "value": 88,
            "unit": "B/op",
            "extra": "5397913 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5397913 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0",
            "value": 220.9,
            "unit": "ns/op\t      88 B/op\t       2 allocs/op",
            "extra": "5414281 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 - ns/op",
            "value": 220.9,
            "unit": "ns/op",
            "extra": "5414281 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 - B/op",
            "value": 88,
            "unit": "B/op",
            "extra": "5414281 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5414281 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0",
            "value": 224.1,
            "unit": "ns/op\t      88 B/op\t       2 allocs/op",
            "extra": "5417862 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 - ns/op",
            "value": 224.1,
            "unit": "ns/op",
            "extra": "5417862 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 - B/op",
            "value": 88,
            "unit": "B/op",
            "extra": "5417862 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5417862 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0",
            "value": 223.3,
            "unit": "ns/op\t      88 B/op\t       2 allocs/op",
            "extra": "5416743 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 - ns/op",
            "value": 223.3,
            "unit": "ns/op",
            "extra": "5416743 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 - B/op",
            "value": 88,
            "unit": "B/op",
            "extra": "5416743 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5416743 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1",
            "value": 222.7,
            "unit": "ns/op\t      88 B/op\t       2 allocs/op",
            "extra": "5499036 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 - ns/op",
            "value": 222.7,
            "unit": "ns/op",
            "extra": "5499036 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 - B/op",
            "value": 88,
            "unit": "B/op",
            "extra": "5499036 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5499036 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1",
            "value": 218.8,
            "unit": "ns/op\t      88 B/op\t       2 allocs/op",
            "extra": "5498374 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 - ns/op",
            "value": 218.8,
            "unit": "ns/op",
            "extra": "5498374 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 - B/op",
            "value": 88,
            "unit": "B/op",
            "extra": "5498374 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5498374 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1",
            "value": 222.4,
            "unit": "ns/op\t      88 B/op\t       2 allocs/op",
            "extra": "5497279 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 - ns/op",
            "value": 222.4,
            "unit": "ns/op",
            "extra": "5497279 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 - B/op",
            "value": 88,
            "unit": "B/op",
            "extra": "5497279 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5497279 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1",
            "value": 218.1,
            "unit": "ns/op\t      88 B/op\t       2 allocs/op",
            "extra": "5486943 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 - ns/op",
            "value": 218.1,
            "unit": "ns/op",
            "extra": "5486943 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 - B/op",
            "value": 88,
            "unit": "B/op",
            "extra": "5486943 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5486943 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1",
            "value": 219.7,
            "unit": "ns/op\t      88 B/op\t       2 allocs/op",
            "extra": "5442225 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 - ns/op",
            "value": 219.7,
            "unit": "ns/op",
            "extra": "5442225 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 - B/op",
            "value": 88,
            "unit": "B/op",
            "extra": "5442225 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "5442225 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0",
            "value": 487.8,
            "unit": "ns/op\t     288 B/op\t       6 allocs/op",
            "extra": "2472142 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 - ns/op",
            "value": 487.8,
            "unit": "ns/op",
            "extra": "2472142 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 - B/op",
            "value": 288,
            "unit": "B/op",
            "extra": "2472142 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "2472142 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0",
            "value": 498.5,
            "unit": "ns/op\t     288 B/op\t       6 allocs/op",
            "extra": "2463105 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 - ns/op",
            "value": 498.5,
            "unit": "ns/op",
            "extra": "2463105 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 - B/op",
            "value": 288,
            "unit": "B/op",
            "extra": "2463105 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "2463105 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0",
            "value": 487.8,
            "unit": "ns/op\t     288 B/op\t       6 allocs/op",
            "extra": "2469661 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 - ns/op",
            "value": 487.8,
            "unit": "ns/op",
            "extra": "2469661 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 - B/op",
            "value": 288,
            "unit": "B/op",
            "extra": "2469661 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "2469661 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0",
            "value": 490.5,
            "unit": "ns/op\t     288 B/op\t       6 allocs/op",
            "extra": "2460769 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 - ns/op",
            "value": 490.5,
            "unit": "ns/op",
            "extra": "2460769 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 - B/op",
            "value": 288,
            "unit": "B/op",
            "extra": "2460769 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "2460769 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0",
            "value": 487.5,
            "unit": "ns/op\t     288 B/op\t       6 allocs/op",
            "extra": "2452971 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 - ns/op",
            "value": 487.5,
            "unit": "ns/op",
            "extra": "2452971 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 - B/op",
            "value": 288,
            "unit": "B/op",
            "extra": "2452971 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "2452971 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1",
            "value": 483.6,
            "unit": "ns/op\t     288 B/op\t       6 allocs/op",
            "extra": "2478141 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 - ns/op",
            "value": 483.6,
            "unit": "ns/op",
            "extra": "2478141 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 - B/op",
            "value": 288,
            "unit": "B/op",
            "extra": "2478141 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "2478141 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1",
            "value": 481.1,
            "unit": "ns/op\t     288 B/op\t       6 allocs/op",
            "extra": "2482671 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 - ns/op",
            "value": 481.1,
            "unit": "ns/op",
            "extra": "2482671 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 - B/op",
            "value": 288,
            "unit": "B/op",
            "extra": "2482671 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "2482671 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1",
            "value": 489.8,
            "unit": "ns/op\t     288 B/op\t       6 allocs/op",
            "extra": "2474060 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 - ns/op",
            "value": 489.8,
            "unit": "ns/op",
            "extra": "2474060 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 - B/op",
            "value": 288,
            "unit": "B/op",
            "extra": "2474060 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "2474060 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1",
            "value": 481.6,
            "unit": "ns/op\t     288 B/op\t       6 allocs/op",
            "extra": "2486028 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 - ns/op",
            "value": 481.6,
            "unit": "ns/op",
            "extra": "2486028 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 - B/op",
            "value": 288,
            "unit": "B/op",
            "extra": "2486028 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "2486028 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1",
            "value": 482.7,
            "unit": "ns/op\t     288 B/op\t       6 allocs/op",
            "extra": "2472897 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 - ns/op",
            "value": 482.7,
            "unit": "ns/op",
            "extra": "2472897 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 - B/op",
            "value": 288,
            "unit": "B/op",
            "extra": "2472897 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "2472897 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0",
            "value": 386.4,
            "unit": "ns/op\t     208 B/op\t       5 allocs/op",
            "extra": "3073358 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 - ns/op",
            "value": 386.4,
            "unit": "ns/op",
            "extra": "3073358 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "3073358 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3073358 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0",
            "value": 386.7,
            "unit": "ns/op\t     208 B/op\t       5 allocs/op",
            "extra": "3098608 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 - ns/op",
            "value": 386.7,
            "unit": "ns/op",
            "extra": "3098608 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "3098608 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3098608 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0",
            "value": 386.2,
            "unit": "ns/op\t     208 B/op\t       5 allocs/op",
            "extra": "3091423 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 - ns/op",
            "value": 386.2,
            "unit": "ns/op",
            "extra": "3091423 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "3091423 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3091423 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0",
            "value": 388.2,
            "unit": "ns/op\t     208 B/op\t       5 allocs/op",
            "extra": "3089595 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 - ns/op",
            "value": 388.2,
            "unit": "ns/op",
            "extra": "3089595 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "3089595 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3089595 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0",
            "value": 388.3,
            "unit": "ns/op\t     208 B/op\t       5 allocs/op",
            "extra": "3099790 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 - ns/op",
            "value": 388.3,
            "unit": "ns/op",
            "extra": "3099790 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "3099790 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3099790 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1",
            "value": 387.5,
            "unit": "ns/op\t     208 B/op\t       5 allocs/op",
            "extra": "3120942 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 - ns/op",
            "value": 387.5,
            "unit": "ns/op",
            "extra": "3120942 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "3120942 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3120942 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1",
            "value": 385.3,
            "unit": "ns/op\t     208 B/op\t       5 allocs/op",
            "extra": "3128103 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 - ns/op",
            "value": 385.3,
            "unit": "ns/op",
            "extra": "3128103 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "3128103 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3128103 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1",
            "value": 385.1,
            "unit": "ns/op\t     208 B/op\t       5 allocs/op",
            "extra": "3119656 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 - ns/op",
            "value": 385.1,
            "unit": "ns/op",
            "extra": "3119656 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "3119656 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3119656 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1",
            "value": 383.4,
            "unit": "ns/op\t     208 B/op\t       5 allocs/op",
            "extra": "3125482 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 - ns/op",
            "value": 383.4,
            "unit": "ns/op",
            "extra": "3125482 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "3125482 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3125482 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1",
            "value": 384.9,
            "unit": "ns/op\t     208 B/op\t       5 allocs/op",
            "extra": "3110023 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 - ns/op",
            "value": 384.9,
            "unit": "ns/op",
            "extra": "3110023 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "3110023 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3110023 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants",
            "value": 9757,
            "unit": "ns/op\t    4464 B/op\t      12 allocs/op",
            "extra": "124474 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants - ns/op",
            "value": 9757,
            "unit": "ns/op",
            "extra": "124474 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants - B/op",
            "value": 4464,
            "unit": "B/op",
            "extra": "124474 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "124474 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants",
            "value": 9740,
            "unit": "ns/op\t    4464 B/op\t      12 allocs/op",
            "extra": "121590 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants - ns/op",
            "value": 9740,
            "unit": "ns/op",
            "extra": "121590 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants - B/op",
            "value": 4464,
            "unit": "B/op",
            "extra": "121590 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "121590 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants",
            "value": 9768,
            "unit": "ns/op\t    4464 B/op\t      12 allocs/op",
            "extra": "121387 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants - ns/op",
            "value": 9768,
            "unit": "ns/op",
            "extra": "121387 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants - B/op",
            "value": 4464,
            "unit": "B/op",
            "extra": "121387 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "121387 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants",
            "value": 9718,
            "unit": "ns/op\t    4464 B/op\t      12 allocs/op",
            "extra": "124384 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants - ns/op",
            "value": 9718,
            "unit": "ns/op",
            "extra": "124384 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants - B/op",
            "value": 4464,
            "unit": "B/op",
            "extra": "124384 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "124384 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants",
            "value": 9776,
            "unit": "ns/op\t    4464 B/op\t      12 allocs/op",
            "extra": "124219 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants - ns/op",
            "value": 9776,
            "unit": "ns/op",
            "extra": "124219 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants - B/op",
            "value": 4464,
            "unit": "B/op",
            "extra": "124219 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "124219 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1",
            "value": 9793,
            "unit": "ns/op\t    4464 B/op\t      12 allocs/op",
            "extra": "124076 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1 - ns/op",
            "value": 9793,
            "unit": "ns/op",
            "extra": "124076 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1 - B/op",
            "value": 4464,
            "unit": "B/op",
            "extra": "124076 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1 - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "124076 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1",
            "value": 9867,
            "unit": "ns/op\t    4464 B/op\t      12 allocs/op",
            "extra": "122827 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1 - ns/op",
            "value": 9867,
            "unit": "ns/op",
            "extra": "122827 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1 - B/op",
            "value": 4464,
            "unit": "B/op",
            "extra": "122827 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1 - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "122827 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1",
            "value": 9752,
            "unit": "ns/op\t    4464 B/op\t      12 allocs/op",
            "extra": "124056 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1 - ns/op",
            "value": 9752,
            "unit": "ns/op",
            "extra": "124056 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1 - B/op",
            "value": 4464,
            "unit": "B/op",
            "extra": "124056 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1 - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "124056 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1",
            "value": 9745,
            "unit": "ns/op\t    4464 B/op\t      12 allocs/op",
            "extra": "123297 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1 - ns/op",
            "value": 9745,
            "unit": "ns/op",
            "extra": "123297 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1 - B/op",
            "value": 4464,
            "unit": "B/op",
            "extra": "123297 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1 - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "123297 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1",
            "value": 9776,
            "unit": "ns/op\t    4464 B/op\t      12 allocs/op",
            "extra": "123799 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1 - ns/op",
            "value": 9776,
            "unit": "ns/op",
            "extra": "123799 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1 - B/op",
            "value": 4464,
            "unit": "B/op",
            "extra": "123799 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants_O1 - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "123799 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0",
            "value": 152.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7858474 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 - ns/op",
            "value": 152.2,
            "unit": "ns/op",
            "extra": "7858474 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7858474 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7858474 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0",
            "value": 152.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7908614 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 - ns/op",
            "value": 152.2,
            "unit": "ns/op",
            "extra": "7908614 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7908614 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7908614 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0",
            "value": 152.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7888873 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 - ns/op",
            "value": 152.1,
            "unit": "ns/op",
            "extra": "7888873 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7888873 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7888873 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0",
            "value": 152.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7890960 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 - ns/op",
            "value": 152.5,
            "unit": "ns/op",
            "extra": "7890960 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7890960 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7890960 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0",
            "value": 152,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7882425 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 - ns/op",
            "value": 152,
            "unit": "ns/op",
            "extra": "7882425 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7882425 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7882425 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1",
            "value": 128.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9393438 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 - ns/op",
            "value": 128.2,
            "unit": "ns/op",
            "extra": "9393438 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9393438 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9393438 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1",
            "value": 127.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9382437 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 - ns/op",
            "value": 127.8,
            "unit": "ns/op",
            "extra": "9382437 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9382437 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9382437 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1",
            "value": 128.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9385030 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 - ns/op",
            "value": 128.6,
            "unit": "ns/op",
            "extra": "9385030 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9385030 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9385030 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1",
            "value": 127.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9380270 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 - ns/op",
            "value": 127.9,
            "unit": "ns/op",
            "extra": "9380270 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9380270 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9380270 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1",
            "value": 128,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9387154 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 - ns/op",
            "value": 128,
            "unit": "ns/op",
            "extra": "9387154 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9387154 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9387154 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0",
            "value": 202,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5970842 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 - ns/op",
            "value": 202,
            "unit": "ns/op",
            "extra": "5970842 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5970842 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5970842 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0",
            "value": 204.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5895710 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 - ns/op",
            "value": 204.5,
            "unit": "ns/op",
            "extra": "5895710 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5895710 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5895710 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0",
            "value": 207.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5944450 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 - ns/op",
            "value": 207.1,
            "unit": "ns/op",
            "extra": "5944450 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5944450 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5944450 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0",
            "value": 203.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5943586 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 - ns/op",
            "value": 203.5,
            "unit": "ns/op",
            "extra": "5943586 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5943586 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5943586 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0",
            "value": 205.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5968042 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 - ns/op",
            "value": 205.8,
            "unit": "ns/op",
            "extra": "5968042 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5968042 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5968042 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1",
            "value": 198.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6141798 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 - ns/op",
            "value": 198.1,
            "unit": "ns/op",
            "extra": "6141798 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6141798 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6141798 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1",
            "value": 198.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6119146 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 - ns/op",
            "value": 198.6,
            "unit": "ns/op",
            "extra": "6119146 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6119146 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6119146 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1",
            "value": 197.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6135226 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 - ns/op",
            "value": 197.7,
            "unit": "ns/op",
            "extra": "6135226 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6135226 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6135226 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1",
            "value": 198.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6142371 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 - ns/op",
            "value": 198.8,
            "unit": "ns/op",
            "extra": "6142371 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6142371 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6142371 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1",
            "value": 197,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6108817 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 - ns/op",
            "value": 197,
            "unit": "ns/op",
            "extra": "6108817 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6108817 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6108817 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0",
            "value": 275.6,
            "unit": "ns/op\t      32 B/op\t       2 allocs/op",
            "extra": "4331944 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 - ns/op",
            "value": 275.6,
            "unit": "ns/op",
            "extra": "4331944 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "4331944 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4331944 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0",
            "value": 276.6,
            "unit": "ns/op\t      32 B/op\t       2 allocs/op",
            "extra": "4328479 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 - ns/op",
            "value": 276.6,
            "unit": "ns/op",
            "extra": "4328479 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "4328479 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4328479 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0",
            "value": 276,
            "unit": "ns/op\t      32 B/op\t       2 allocs/op",
            "extra": "4305837 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 - ns/op",
            "value": 276,
            "unit": "ns/op",
            "extra": "4305837 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "4305837 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4305837 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0",
            "value": 275.6,
            "unit": "ns/op\t      32 B/op\t       2 allocs/op",
            "extra": "4331562 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 - ns/op",
            "value": 275.6,
            "unit": "ns/op",
            "extra": "4331562 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "4331562 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4331562 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0",
            "value": 276.2,
            "unit": "ns/op\t      32 B/op\t       2 allocs/op",
            "extra": "4336833 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 - ns/op",
            "value": 276.2,
            "unit": "ns/op",
            "extra": "4336833 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "4336833 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4336833 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1",
            "value": 257.1,
            "unit": "ns/op\t      32 B/op\t       2 allocs/op",
            "extra": "4651413 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 - ns/op",
            "value": 257.1,
            "unit": "ns/op",
            "extra": "4651413 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "4651413 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4651413 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1",
            "value": 257,
            "unit": "ns/op\t      32 B/op\t       2 allocs/op",
            "extra": "4648777 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 - ns/op",
            "value": 257,
            "unit": "ns/op",
            "extra": "4648777 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "4648777 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4648777 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1",
            "value": 256.5,
            "unit": "ns/op\t      32 B/op\t       2 allocs/op",
            "extra": "4653130 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 - ns/op",
            "value": 256.5,
            "unit": "ns/op",
            "extra": "4653130 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "4653130 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4653130 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1",
            "value": 260,
            "unit": "ns/op\t      32 B/op\t       2 allocs/op",
            "extra": "4654501 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 - ns/op",
            "value": 260,
            "unit": "ns/op",
            "extra": "4654501 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "4654501 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4654501 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1",
            "value": 257.5,
            "unit": "ns/op\t      32 B/op\t       2 allocs/op",
            "extra": "4610667 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 - ns/op",
            "value": 257.5,
            "unit": "ns/op",
            "extra": "4610667 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "4610667 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "4610667 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0",
            "value": 140152,
            "unit": "ns/op\t   65785 B/op\t     780 allocs/op",
            "extra": "8608 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 - ns/op",
            "value": 140152,
            "unit": "ns/op",
            "extra": "8608 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 - B/op",
            "value": 65785,
            "unit": "B/op",
            "extra": "8608 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 - allocs/op",
            "value": 780,
            "unit": "allocs/op",
            "extra": "8608 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0",
            "value": 138669,
            "unit": "ns/op\t   65785 B/op\t     780 allocs/op",
            "extra": "8330 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 - ns/op",
            "value": 138669,
            "unit": "ns/op",
            "extra": "8330 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 - B/op",
            "value": 65785,
            "unit": "B/op",
            "extra": "8330 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 - allocs/op",
            "value": 780,
            "unit": "allocs/op",
            "extra": "8330 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0",
            "value": 138575,
            "unit": "ns/op\t   65785 B/op\t     780 allocs/op",
            "extra": "8443 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 - ns/op",
            "value": 138575,
            "unit": "ns/op",
            "extra": "8443 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 - B/op",
            "value": 65785,
            "unit": "B/op",
            "extra": "8443 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 - allocs/op",
            "value": 780,
            "unit": "allocs/op",
            "extra": "8443 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0",
            "value": 138397,
            "unit": "ns/op\t   65785 B/op\t     780 allocs/op",
            "extra": "8202 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 - ns/op",
            "value": 138397,
            "unit": "ns/op",
            "extra": "8202 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 - B/op",
            "value": 65785,
            "unit": "B/op",
            "extra": "8202 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 - allocs/op",
            "value": 780,
            "unit": "allocs/op",
            "extra": "8202 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0",
            "value": 138039,
            "unit": "ns/op\t   65785 B/op\t     780 allocs/op",
            "extra": "8140 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 - ns/op",
            "value": 138039,
            "unit": "ns/op",
            "extra": "8140 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 - B/op",
            "value": 65785,
            "unit": "B/op",
            "extra": "8140 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 - allocs/op",
            "value": 780,
            "unit": "allocs/op",
            "extra": "8140 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1",
            "value": 128105,
            "unit": "ns/op\t   65785 B/op\t     780 allocs/op",
            "extra": "9081 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 - ns/op",
            "value": 128105,
            "unit": "ns/op",
            "extra": "9081 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 - B/op",
            "value": 65785,
            "unit": "B/op",
            "extra": "9081 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 - allocs/op",
            "value": 780,
            "unit": "allocs/op",
            "extra": "9081 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1",
            "value": 126678,
            "unit": "ns/op\t   65785 B/op\t     780 allocs/op",
            "extra": "9556 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 - ns/op",
            "value": 126678,
            "unit": "ns/op",
            "extra": "9556 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 - B/op",
            "value": 65785,
            "unit": "B/op",
            "extra": "9556 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 - allocs/op",
            "value": 780,
            "unit": "allocs/op",
            "extra": "9556 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1",
            "value": 126888,
            "unit": "ns/op\t   65785 B/op\t     780 allocs/op",
            "extra": "9282 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 - ns/op",
            "value": 126888,
            "unit": "ns/op",
            "extra": "9282 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 - B/op",
            "value": 65785,
            "unit": "B/op",
            "extra": "9282 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 - allocs/op",
            "value": 780,
            "unit": "allocs/op",
            "extra": "9282 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1",
            "value": 126085,
            "unit": "ns/op\t   65785 B/op\t     780 allocs/op",
            "extra": "9114 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 - ns/op",
            "value": 126085,
            "unit": "ns/op",
            "extra": "9114 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 - B/op",
            "value": 65785,
            "unit": "B/op",
            "extra": "9114 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 - allocs/op",
            "value": 780,
            "unit": "allocs/op",
            "extra": "9114 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1",
            "value": 126809,
            "unit": "ns/op\t   65785 B/op\t     780 allocs/op",
            "extra": "8678 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 - ns/op",
            "value": 126809,
            "unit": "ns/op",
            "extra": "8678 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 - B/op",
            "value": 65785,
            "unit": "B/op",
            "extra": "8678 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 - allocs/op",
            "value": 780,
            "unit": "allocs/op",
            "extra": "8678 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0",
            "value": 84.99,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14129383 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 - ns/op",
            "value": 84.99,
            "unit": "ns/op",
            "extra": "14129383 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14129383 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14129383 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0",
            "value": 85.07,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14171721 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 - ns/op",
            "value": 85.07,
            "unit": "ns/op",
            "extra": "14171721 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14171721 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14171721 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0",
            "value": 84.66,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14121855 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 - ns/op",
            "value": 84.66,
            "unit": "ns/op",
            "extra": "14121855 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14121855 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14121855 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0",
            "value": 85.05,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14101370 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 - ns/op",
            "value": 85.05,
            "unit": "ns/op",
            "extra": "14101370 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14101370 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14101370 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0",
            "value": 84.93,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14126239 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 - ns/op",
            "value": 84.93,
            "unit": "ns/op",
            "extra": "14126239 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14126239 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14126239 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1",
            "value": 63.44,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18862424 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 - ns/op",
            "value": 63.44,
            "unit": "ns/op",
            "extra": "18862424 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18862424 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18862424 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1",
            "value": 63.55,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18925491 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 - ns/op",
            "value": 63.55,
            "unit": "ns/op",
            "extra": "18925491 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18925491 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18925491 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1",
            "value": 63.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18941752 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 - ns/op",
            "value": 63.6,
            "unit": "ns/op",
            "extra": "18941752 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18941752 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18941752 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1",
            "value": 69.22,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18951601 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 - ns/op",
            "value": 69.22,
            "unit": "ns/op",
            "extra": "18951601 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18951601 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18951601 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1",
            "value": 63.42,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18913869 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 - ns/op",
            "value": 63.42,
            "unit": "ns/op",
            "extra": "18913869 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18913869 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18913869 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0",
            "value": 151.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7900825 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 - ns/op",
            "value": 151.7,
            "unit": "ns/op",
            "extra": "7900825 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7900825 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7900825 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0",
            "value": 151.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7900696 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 - ns/op",
            "value": 151.9,
            "unit": "ns/op",
            "extra": "7900696 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7900696 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7900696 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0",
            "value": 153.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7906113 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 - ns/op",
            "value": 153.1,
            "unit": "ns/op",
            "extra": "7906113 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7906113 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7906113 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0",
            "value": 152.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7880155 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 - ns/op",
            "value": 152.2,
            "unit": "ns/op",
            "extra": "7880155 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7880155 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7880155 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0",
            "value": 152,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7913252 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 - ns/op",
            "value": 152,
            "unit": "ns/op",
            "extra": "7913252 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7913252 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7913252 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1",
            "value": 127.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9389181 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 - ns/op",
            "value": 127.7,
            "unit": "ns/op",
            "extra": "9389181 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9389181 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9389181 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1",
            "value": 130.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9407163 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 - ns/op",
            "value": 130.3,
            "unit": "ns/op",
            "extra": "9407163 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9407163 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9407163 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1",
            "value": 128.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9380218 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 - ns/op",
            "value": 128.3,
            "unit": "ns/op",
            "extra": "9380218 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9380218 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9380218 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1",
            "value": 128.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9407148 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 - ns/op",
            "value": 128.2,
            "unit": "ns/op",
            "extra": "9407148 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9407148 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9407148 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1",
            "value": 128.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9429956 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 - ns/op",
            "value": 128.7,
            "unit": "ns/op",
            "extra": "9429956 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9429956 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9429956 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0",
            "value": 86.18,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13963248 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 - ns/op",
            "value": 86.18,
            "unit": "ns/op",
            "extra": "13963248 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13963248 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13963248 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0",
            "value": 86.33,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13963861 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 - ns/op",
            "value": 86.33,
            "unit": "ns/op",
            "extra": "13963861 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13963861 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13963861 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0",
            "value": 86.19,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13942538 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 - ns/op",
            "value": 86.19,
            "unit": "ns/op",
            "extra": "13942538 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13942538 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13942538 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0",
            "value": 86.13,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13956385 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 - ns/op",
            "value": 86.13,
            "unit": "ns/op",
            "extra": "13956385 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13956385 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13956385 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0",
            "value": 86.62,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13873333 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 - ns/op",
            "value": 86.62,
            "unit": "ns/op",
            "extra": "13873333 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13873333 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13873333 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1",
            "value": 85.03,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14099212 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 - ns/op",
            "value": 85.03,
            "unit": "ns/op",
            "extra": "14099212 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14099212 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14099212 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1",
            "value": 85.53,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14127861 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 - ns/op",
            "value": 85.53,
            "unit": "ns/op",
            "extra": "14127861 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14127861 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14127861 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1",
            "value": 86.92,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14101554 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 - ns/op",
            "value": 86.92,
            "unit": "ns/op",
            "extra": "14101554 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14101554 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14101554 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1",
            "value": 85.06,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14088171 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 - ns/op",
            "value": 85.06,
            "unit": "ns/op",
            "extra": "14088171 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14088171 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14088171 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1",
            "value": 85.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14094051 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 - ns/op",
            "value": 85.67,
            "unit": "ns/op",
            "extra": "14094051 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14094051 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14094051 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0",
            "value": 146.9,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "8246868 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 - ns/op",
            "value": 146.9,
            "unit": "ns/op",
            "extra": "8246868 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "8246868 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8246868 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0",
            "value": 144.8,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "8150394 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 - ns/op",
            "value": 144.8,
            "unit": "ns/op",
            "extra": "8150394 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "8150394 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8150394 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0",
            "value": 143.8,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "8287724 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 - ns/op",
            "value": 143.8,
            "unit": "ns/op",
            "extra": "8287724 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "8287724 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8287724 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0",
            "value": 145,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "8319501 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 - ns/op",
            "value": 145,
            "unit": "ns/op",
            "extra": "8319501 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "8319501 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8319501 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0",
            "value": 145.6,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "8255936 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 - ns/op",
            "value": 145.6,
            "unit": "ns/op",
            "extra": "8255936 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "8255936 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8255936 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1",
            "value": 146.2,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "8194750 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 - ns/op",
            "value": 146.2,
            "unit": "ns/op",
            "extra": "8194750 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "8194750 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8194750 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1",
            "value": 147.1,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "8151402 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 - ns/op",
            "value": 147.1,
            "unit": "ns/op",
            "extra": "8151402 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "8151402 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8151402 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1",
            "value": 146.3,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "8133973 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 - ns/op",
            "value": 146.3,
            "unit": "ns/op",
            "extra": "8133973 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "8133973 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8133973 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1",
            "value": 146.6,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "8111515 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 - ns/op",
            "value": 146.6,
            "unit": "ns/op",
            "extra": "8111515 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "8111515 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8111515 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1",
            "value": 146.9,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "8044668 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 - ns/op",
            "value": 146.9,
            "unit": "ns/op",
            "extra": "8044668 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "8044668 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8044668 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0",
            "value": 86.12,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14030958 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 - ns/op",
            "value": 86.12,
            "unit": "ns/op",
            "extra": "14030958 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14030958 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14030958 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0",
            "value": 85.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14005137 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 - ns/op",
            "value": 85.83,
            "unit": "ns/op",
            "extra": "14005137 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14005137 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14005137 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0",
            "value": 86.03,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14008146 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 - ns/op",
            "value": 86.03,
            "unit": "ns/op",
            "extra": "14008146 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14008146 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14008146 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0",
            "value": 85.86,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14012247 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 - ns/op",
            "value": 85.86,
            "unit": "ns/op",
            "extra": "14012247 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14012247 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14012247 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0",
            "value": 86.05,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14028603 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 - ns/op",
            "value": 86.05,
            "unit": "ns/op",
            "extra": "14028603 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14028603 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14028603 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1",
            "value": 88.22,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13718150 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 - ns/op",
            "value": 88.22,
            "unit": "ns/op",
            "extra": "13718150 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13718150 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13718150 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1",
            "value": 87.97,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13649266 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 - ns/op",
            "value": 87.97,
            "unit": "ns/op",
            "extra": "13649266 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13649266 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13649266 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1",
            "value": 86.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13684801 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 - ns/op",
            "value": 86.3,
            "unit": "ns/op",
            "extra": "13684801 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13684801 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13684801 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1",
            "value": 88.11,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13627645 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 - ns/op",
            "value": 88.11,
            "unit": "ns/op",
            "extra": "13627645 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13627645 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13627645 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1",
            "value": 88.07,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13692712 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 - ns/op",
            "value": 88.07,
            "unit": "ns/op",
            "extra": "13692712 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13692712 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13692712 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0",
            "value": 108.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11119737 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 - ns/op",
            "value": 108.2,
            "unit": "ns/op",
            "extra": "11119737 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11119737 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11119737 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0",
            "value": 108.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11091417 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 - ns/op",
            "value": 108.7,
            "unit": "ns/op",
            "extra": "11091417 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11091417 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11091417 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0",
            "value": 108.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11101640 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 - ns/op",
            "value": 108.4,
            "unit": "ns/op",
            "extra": "11101640 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11101640 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11101640 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0",
            "value": 108.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11067547 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 - ns/op",
            "value": 108.6,
            "unit": "ns/op",
            "extra": "11067547 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11067547 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11067547 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0",
            "value": 108.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11082026 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 - ns/op",
            "value": 108.3,
            "unit": "ns/op",
            "extra": "11082026 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11082026 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11082026 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1",
            "value": 107.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11181669 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 - ns/op",
            "value": 107.5,
            "unit": "ns/op",
            "extra": "11181669 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11181669 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11181669 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1",
            "value": 107.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11196030 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 - ns/op",
            "value": 107.4,
            "unit": "ns/op",
            "extra": "11196030 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11196030 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11196030 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1",
            "value": 107.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11186373 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 - ns/op",
            "value": 107.8,
            "unit": "ns/op",
            "extra": "11186373 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11186373 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11186373 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1",
            "value": 107.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10402610 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 - ns/op",
            "value": 107.5,
            "unit": "ns/op",
            "extra": "10402610 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10402610 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10402610 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1",
            "value": 107.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11016138 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 - ns/op",
            "value": 107.7,
            "unit": "ns/op",
            "extra": "11016138 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11016138 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11016138 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0",
            "value": 108.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11107401 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 - ns/op",
            "value": 108.2,
            "unit": "ns/op",
            "extra": "11107401 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11107401 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11107401 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0",
            "value": 109,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11047687 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 - ns/op",
            "value": 109,
            "unit": "ns/op",
            "extra": "11047687 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11047687 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11047687 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0",
            "value": 108.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11119910 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 - ns/op",
            "value": 108.5,
            "unit": "ns/op",
            "extra": "11119910 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11119910 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11119910 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0",
            "value": 108.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11095549 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 - ns/op",
            "value": 108.3,
            "unit": "ns/op",
            "extra": "11095549 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11095549 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11095549 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0",
            "value": 108.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11057032 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 - ns/op",
            "value": 108.6,
            "unit": "ns/op",
            "extra": "11057032 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11057032 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11057032 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1",
            "value": 110.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11072781 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 - ns/op",
            "value": 110.6,
            "unit": "ns/op",
            "extra": "11072781 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11072781 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11072781 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1",
            "value": 108.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11114792 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 - ns/op",
            "value": 108.3,
            "unit": "ns/op",
            "extra": "11114792 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11114792 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11114792 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1",
            "value": 108.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11053934 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 - ns/op",
            "value": 108.5,
            "unit": "ns/op",
            "extra": "11053934 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11053934 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11053934 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1",
            "value": 108.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11107714 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 - ns/op",
            "value": 108.5,
            "unit": "ns/op",
            "extra": "11107714 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11107714 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11107714 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1",
            "value": 108.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11073693 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 - ns/op",
            "value": 108.2,
            "unit": "ns/op",
            "extra": "11073693 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11073693 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11073693 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0",
            "value": 231.8,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "5100190 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 - ns/op",
            "value": 231.8,
            "unit": "ns/op",
            "extra": "5100190 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "5100190 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5100190 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0",
            "value": 232.8,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "5173258 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 - ns/op",
            "value": 232.8,
            "unit": "ns/op",
            "extra": "5173258 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "5173258 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5173258 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0",
            "value": 231.7,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "5101927 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 - ns/op",
            "value": 231.7,
            "unit": "ns/op",
            "extra": "5101927 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "5101927 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5101927 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0",
            "value": 231.2,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "5164168 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 - ns/op",
            "value": 231.2,
            "unit": "ns/op",
            "extra": "5164168 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "5164168 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5164168 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0",
            "value": 232.6,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "5123216 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 - ns/op",
            "value": 232.6,
            "unit": "ns/op",
            "extra": "5123216 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "5123216 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5123216 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1",
            "value": 238.7,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "4993858 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 - ns/op",
            "value": 238.7,
            "unit": "ns/op",
            "extra": "4993858 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "4993858 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4993858 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1",
            "value": 237.3,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "5014423 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 - ns/op",
            "value": 237.3,
            "unit": "ns/op",
            "extra": "5014423 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "5014423 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5014423 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1",
            "value": 237.7,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "5015492 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 - ns/op",
            "value": 237.7,
            "unit": "ns/op",
            "extra": "5015492 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "5015492 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5015492 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1",
            "value": 238.3,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "5010949 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 - ns/op",
            "value": 238.3,
            "unit": "ns/op",
            "extra": "5010949 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "5010949 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5010949 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1",
            "value": 239.2,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "5019198 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 - ns/op",
            "value": 239.2,
            "unit": "ns/op",
            "extra": "5019198 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "5019198 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5019198 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0",
            "value": 862.5,
            "unit": "ns/op\t     432 B/op\t      15 allocs/op",
            "extra": "1395860 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 - ns/op",
            "value": 862.5,
            "unit": "ns/op",
            "extra": "1395860 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 - B/op",
            "value": 432,
            "unit": "B/op",
            "extra": "1395860 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "1395860 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0",
            "value": 856.9,
            "unit": "ns/op\t     432 B/op\t      15 allocs/op",
            "extra": "1399366 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 - ns/op",
            "value": 856.9,
            "unit": "ns/op",
            "extra": "1399366 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 - B/op",
            "value": 432,
            "unit": "B/op",
            "extra": "1399366 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "1399366 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0",
            "value": 853.9,
            "unit": "ns/op\t     432 B/op\t      15 allocs/op",
            "extra": "1389764 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 - ns/op",
            "value": 853.9,
            "unit": "ns/op",
            "extra": "1389764 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 - B/op",
            "value": 432,
            "unit": "B/op",
            "extra": "1389764 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "1389764 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0",
            "value": 859.4,
            "unit": "ns/op\t     432 B/op\t      15 allocs/op",
            "extra": "1398045 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 - ns/op",
            "value": 859.4,
            "unit": "ns/op",
            "extra": "1398045 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 - B/op",
            "value": 432,
            "unit": "B/op",
            "extra": "1398045 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "1398045 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0",
            "value": 861.7,
            "unit": "ns/op\t     432 B/op\t      15 allocs/op",
            "extra": "1399632 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 - ns/op",
            "value": 861.7,
            "unit": "ns/op",
            "extra": "1399632 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 - B/op",
            "value": 432,
            "unit": "B/op",
            "extra": "1399632 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "1399632 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1",
            "value": 855.4,
            "unit": "ns/op\t     432 B/op\t      15 allocs/op",
            "extra": "1403966 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 - ns/op",
            "value": 855.4,
            "unit": "ns/op",
            "extra": "1403966 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 - B/op",
            "value": 432,
            "unit": "B/op",
            "extra": "1403966 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "1403966 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1",
            "value": 853.2,
            "unit": "ns/op\t     432 B/op\t      15 allocs/op",
            "extra": "1399197 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 - ns/op",
            "value": 853.2,
            "unit": "ns/op",
            "extra": "1399197 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 - B/op",
            "value": 432,
            "unit": "B/op",
            "extra": "1399197 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "1399197 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1",
            "value": 851.2,
            "unit": "ns/op\t     432 B/op\t      15 allocs/op",
            "extra": "1407163 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 - ns/op",
            "value": 851.2,
            "unit": "ns/op",
            "extra": "1407163 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 - B/op",
            "value": 432,
            "unit": "B/op",
            "extra": "1407163 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "1407163 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1",
            "value": 851.1,
            "unit": "ns/op\t     432 B/op\t      15 allocs/op",
            "extra": "1410936 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 - ns/op",
            "value": 851.1,
            "unit": "ns/op",
            "extra": "1410936 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 - B/op",
            "value": 432,
            "unit": "B/op",
            "extra": "1410936 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "1410936 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1",
            "value": 855.2,
            "unit": "ns/op\t     432 B/op\t      15 allocs/op",
            "extra": "1404094 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 - ns/op",
            "value": 855.2,
            "unit": "ns/op",
            "extra": "1404094 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 - B/op",
            "value": 432,
            "unit": "B/op",
            "extra": "1404094 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "1404094 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0",
            "value": 157.1,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "7532696 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 - ns/op",
            "value": 157.1,
            "unit": "ns/op",
            "extra": "7532696 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "7532696 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7532696 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0",
            "value": 156.6,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "7632870 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 - ns/op",
            "value": 156.6,
            "unit": "ns/op",
            "extra": "7632870 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "7632870 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7632870 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0",
            "value": 156.7,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "7653835 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 - ns/op",
            "value": 156.7,
            "unit": "ns/op",
            "extra": "7653835 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "7653835 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7653835 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0",
            "value": 156.7,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "7618375 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 - ns/op",
            "value": 156.7,
            "unit": "ns/op",
            "extra": "7618375 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "7618375 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7618375 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0",
            "value": 157,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "7613172 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 - ns/op",
            "value": 157,
            "unit": "ns/op",
            "extra": "7613172 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "7613172 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7613172 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1",
            "value": 159.9,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "7561069 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 - ns/op",
            "value": 159.9,
            "unit": "ns/op",
            "extra": "7561069 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "7561069 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7561069 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1",
            "value": 158.9,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "7573760 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 - ns/op",
            "value": 158.9,
            "unit": "ns/op",
            "extra": "7573760 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "7573760 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7573760 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1",
            "value": 158.4,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "7542921 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 - ns/op",
            "value": 158.4,
            "unit": "ns/op",
            "extra": "7542921 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "7542921 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7542921 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1",
            "value": 159.2,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "7482316 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 - ns/op",
            "value": 159.2,
            "unit": "ns/op",
            "extra": "7482316 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "7482316 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7482316 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1",
            "value": 160.5,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "7541359 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 - ns/op",
            "value": 160.5,
            "unit": "ns/op",
            "extra": "7541359 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "7541359 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7541359 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0",
            "value": 280.3,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "4271670 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 - ns/op",
            "value": 280.3,
            "unit": "ns/op",
            "extra": "4271670 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "4271670 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4271670 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0",
            "value": 280,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "4284414 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 - ns/op",
            "value": 280,
            "unit": "ns/op",
            "extra": "4284414 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "4284414 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4284414 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0",
            "value": 281.7,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "4287580 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 - ns/op",
            "value": 281.7,
            "unit": "ns/op",
            "extra": "4287580 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "4287580 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4287580 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0",
            "value": 280.1,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "4263104 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 - ns/op",
            "value": 280.1,
            "unit": "ns/op",
            "extra": "4263104 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "4263104 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4263104 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0",
            "value": 280.1,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "4273574 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 - ns/op",
            "value": 280.1,
            "unit": "ns/op",
            "extra": "4273574 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "4273574 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4273574 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1",
            "value": 284.3,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "4192734 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 - ns/op",
            "value": 284.3,
            "unit": "ns/op",
            "extra": "4192734 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "4192734 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4192734 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1",
            "value": 285.5,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "4225764 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 - ns/op",
            "value": 285.5,
            "unit": "ns/op",
            "extra": "4225764 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "4225764 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4225764 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1",
            "value": 284.4,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "4165302 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 - ns/op",
            "value": 284.4,
            "unit": "ns/op",
            "extra": "4165302 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "4165302 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4165302 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1",
            "value": 284.5,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "4202060 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 - ns/op",
            "value": 284.5,
            "unit": "ns/op",
            "extra": "4202060 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "4202060 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4202060 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1",
            "value": 285.2,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "4204502 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 - ns/op",
            "value": 285.2,
            "unit": "ns/op",
            "extra": "4204502 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "4204502 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "4204502 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0",
            "value": 92786,
            "unit": "ns/op\t   47000 B/op\t    1501 allocs/op",
            "extra": "13040 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 - ns/op",
            "value": 92786,
            "unit": "ns/op",
            "extra": "13040 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 - B/op",
            "value": 47000,
            "unit": "B/op",
            "extra": "13040 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "13040 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0",
            "value": 93061,
            "unit": "ns/op\t   47000 B/op\t    1501 allocs/op",
            "extra": "12969 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 - ns/op",
            "value": 93061,
            "unit": "ns/op",
            "extra": "12969 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 - B/op",
            "value": 47000,
            "unit": "B/op",
            "extra": "12969 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12969 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0",
            "value": 93980,
            "unit": "ns/op\t   47000 B/op\t    1501 allocs/op",
            "extra": "12895 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 - ns/op",
            "value": 93980,
            "unit": "ns/op",
            "extra": "12895 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 - B/op",
            "value": 47000,
            "unit": "B/op",
            "extra": "12895 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12895 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0",
            "value": 92323,
            "unit": "ns/op\t   47000 B/op\t    1501 allocs/op",
            "extra": "12916 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 - ns/op",
            "value": 92323,
            "unit": "ns/op",
            "extra": "12916 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 - B/op",
            "value": 47000,
            "unit": "B/op",
            "extra": "12916 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12916 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0",
            "value": 93034,
            "unit": "ns/op\t   47000 B/op\t    1501 allocs/op",
            "extra": "13038 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 - ns/op",
            "value": 93034,
            "unit": "ns/op",
            "extra": "13038 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 - B/op",
            "value": 47000,
            "unit": "B/op",
            "extra": "13038 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "13038 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1",
            "value": 92585,
            "unit": "ns/op\t   47000 B/op\t    1501 allocs/op",
            "extra": "12878 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 - ns/op",
            "value": 92585,
            "unit": "ns/op",
            "extra": "12878 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 - B/op",
            "value": 47000,
            "unit": "B/op",
            "extra": "12878 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12878 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1",
            "value": 93348,
            "unit": "ns/op\t   47000 B/op\t    1501 allocs/op",
            "extra": "12879 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 - ns/op",
            "value": 93348,
            "unit": "ns/op",
            "extra": "12879 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 - B/op",
            "value": 47000,
            "unit": "B/op",
            "extra": "12879 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12879 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1",
            "value": 92946,
            "unit": "ns/op\t   47000 B/op\t    1501 allocs/op",
            "extra": "12990 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 - ns/op",
            "value": 92946,
            "unit": "ns/op",
            "extra": "12990 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 - B/op",
            "value": 47000,
            "unit": "B/op",
            "extra": "12990 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12990 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1",
            "value": 93438,
            "unit": "ns/op\t   47000 B/op\t    1501 allocs/op",
            "extra": "12920 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 - ns/op",
            "value": 93438,
            "unit": "ns/op",
            "extra": "12920 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 - B/op",
            "value": 47000,
            "unit": "B/op",
            "extra": "12920 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12920 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1",
            "value": 96209,
            "unit": "ns/op\t   47000 B/op\t    1501 allocs/op",
            "extra": "12742 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 - ns/op",
            "value": 96209,
            "unit": "ns/op",
            "extra": "12742 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 - B/op",
            "value": 47000,
            "unit": "B/op",
            "extra": "12742 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12742 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0",
            "value": 2922,
            "unit": "ns/op\t    1112 B/op\t      31 allocs/op",
            "extra": "407379 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 - ns/op",
            "value": 2922,
            "unit": "ns/op",
            "extra": "407379 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "407379 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "407379 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0",
            "value": 2859,
            "unit": "ns/op\t    1113 B/op\t      31 allocs/op",
            "extra": "419437 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 - ns/op",
            "value": 2859,
            "unit": "ns/op",
            "extra": "419437 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 - B/op",
            "value": 1113,
            "unit": "B/op",
            "extra": "419437 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "419437 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0",
            "value": 2879,
            "unit": "ns/op\t    1112 B/op\t      31 allocs/op",
            "extra": "402775 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 - ns/op",
            "value": 2879,
            "unit": "ns/op",
            "extra": "402775 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "402775 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "402775 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0",
            "value": 2867,
            "unit": "ns/op\t    1111 B/op\t      31 allocs/op",
            "extra": "399067 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 - ns/op",
            "value": 2867,
            "unit": "ns/op",
            "extra": "399067 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 - B/op",
            "value": 1111,
            "unit": "B/op",
            "extra": "399067 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "399067 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0",
            "value": 2895,
            "unit": "ns/op\t    1111 B/op\t      31 allocs/op",
            "extra": "408904 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 - ns/op",
            "value": 2895,
            "unit": "ns/op",
            "extra": "408904 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 - B/op",
            "value": 1111,
            "unit": "B/op",
            "extra": "408904 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "408904 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1",
            "value": 2846,
            "unit": "ns/op\t    1111 B/op\t      31 allocs/op",
            "extra": "413346 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 - ns/op",
            "value": 2846,
            "unit": "ns/op",
            "extra": "413346 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 - B/op",
            "value": 1111,
            "unit": "B/op",
            "extra": "413346 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "413346 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1",
            "value": 2840,
            "unit": "ns/op\t    1112 B/op\t      31 allocs/op",
            "extra": "397503 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 - ns/op",
            "value": 2840,
            "unit": "ns/op",
            "extra": "397503 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "397503 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "397503 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1",
            "value": 2857,
            "unit": "ns/op\t    1112 B/op\t      31 allocs/op",
            "extra": "389964 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 - ns/op",
            "value": 2857,
            "unit": "ns/op",
            "extra": "389964 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "389964 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "389964 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1",
            "value": 2854,
            "unit": "ns/op\t    1112 B/op\t      31 allocs/op",
            "extra": "397635 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 - ns/op",
            "value": 2854,
            "unit": "ns/op",
            "extra": "397635 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "397635 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "397635 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1",
            "value": 2874,
            "unit": "ns/op\t    1111 B/op\t      31 allocs/op",
            "extra": "409723 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 - ns/op",
            "value": 2874,
            "unit": "ns/op",
            "extra": "409723 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 - B/op",
            "value": 1111,
            "unit": "B/op",
            "extra": "409723 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "409723 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0",
            "value": 129.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9293335 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 - ns/op",
            "value": 129.5,
            "unit": "ns/op",
            "extra": "9293335 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9293335 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9293335 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0",
            "value": 129.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9325393 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 - ns/op",
            "value": 129.6,
            "unit": "ns/op",
            "extra": "9325393 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9325393 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9325393 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0",
            "value": 129.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9275126 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 - ns/op",
            "value": 129.3,
            "unit": "ns/op",
            "extra": "9275126 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9275126 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9275126 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0",
            "value": 129.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9279312 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 - ns/op",
            "value": 129.9,
            "unit": "ns/op",
            "extra": "9279312 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9279312 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9279312 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0",
            "value": 130.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9289604 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 - ns/op",
            "value": 130.9,
            "unit": "ns/op",
            "extra": "9289604 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9289604 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9289604 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1",
            "value": 129.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9219590 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 - ns/op",
            "value": 129.6,
            "unit": "ns/op",
            "extra": "9219590 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9219590 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9219590 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1",
            "value": 131.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9235052 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 - ns/op",
            "value": 131.1,
            "unit": "ns/op",
            "extra": "9235052 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9235052 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9235052 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1",
            "value": 131.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9278606 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 - ns/op",
            "value": 131.2,
            "unit": "ns/op",
            "extra": "9278606 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9278606 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9278606 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1",
            "value": 130.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9258922 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 - ns/op",
            "value": 130.3,
            "unit": "ns/op",
            "extra": "9258922 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9258922 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9258922 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1",
            "value": 131.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9229473 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 - ns/op",
            "value": 131.3,
            "unit": "ns/op",
            "extra": "9229473 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9229473 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9229473 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0",
            "value": 192,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6376046 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 - ns/op",
            "value": 192,
            "unit": "ns/op",
            "extra": "6376046 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6376046 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6376046 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0",
            "value": 188.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6400975 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 - ns/op",
            "value": 188.4,
            "unit": "ns/op",
            "extra": "6400975 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6400975 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6400975 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0",
            "value": 187.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6361686 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 - ns/op",
            "value": 187.7,
            "unit": "ns/op",
            "extra": "6361686 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6361686 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6361686 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0",
            "value": 188.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6386355 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 - ns/op",
            "value": 188.8,
            "unit": "ns/op",
            "extra": "6386355 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6386355 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6386355 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0",
            "value": 188.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6386836 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 - ns/op",
            "value": 188.5,
            "unit": "ns/op",
            "extra": "6386836 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6386836 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6386836 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1",
            "value": 188.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6401709 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 - ns/op",
            "value": 188.8,
            "unit": "ns/op",
            "extra": "6401709 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6401709 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6401709 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1",
            "value": 188.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6395656 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 - ns/op",
            "value": 188.2,
            "unit": "ns/op",
            "extra": "6395656 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6395656 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6395656 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1",
            "value": 188,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6388443 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 - ns/op",
            "value": 188,
            "unit": "ns/op",
            "extra": "6388443 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6388443 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6388443 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1",
            "value": 188.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6391448 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 - ns/op",
            "value": 188.2,
            "unit": "ns/op",
            "extra": "6391448 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6391448 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6391448 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1",
            "value": 188.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6376478 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 - ns/op",
            "value": 188.4,
            "unit": "ns/op",
            "extra": "6376478 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6376478 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6376478 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0",
            "value": 140.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8795570 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 - ns/op",
            "value": 140.3,
            "unit": "ns/op",
            "extra": "8795570 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8795570 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8795570 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0",
            "value": 140,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8812141 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 - ns/op",
            "value": 140,
            "unit": "ns/op",
            "extra": "8812141 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8812141 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8812141 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0",
            "value": 136.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8376914 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 - ns/op",
            "value": 136.3,
            "unit": "ns/op",
            "extra": "8376914 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8376914 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8376914 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0",
            "value": 138,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8809975 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 - ns/op",
            "value": 138,
            "unit": "ns/op",
            "extra": "8809975 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8809975 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8809975 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0",
            "value": 136.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8804671 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 - ns/op",
            "value": 136.8,
            "unit": "ns/op",
            "extra": "8804671 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8804671 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8804671 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1",
            "value": 138.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8723859 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 - ns/op",
            "value": 138.3,
            "unit": "ns/op",
            "extra": "8723859 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8723859 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8723859 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1",
            "value": 138.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8783764 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 - ns/op",
            "value": 138.3,
            "unit": "ns/op",
            "extra": "8783764 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8783764 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8783764 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1",
            "value": 136.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8802550 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 - ns/op",
            "value": 136.5,
            "unit": "ns/op",
            "extra": "8802550 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8802550 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8802550 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1",
            "value": 137.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8747884 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 - ns/op",
            "value": 137.5,
            "unit": "ns/op",
            "extra": "8747884 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8747884 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8747884 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1",
            "value": 136.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8746959 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 - ns/op",
            "value": 136.2,
            "unit": "ns/op",
            "extra": "8746959 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8746959 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8746959 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0",
            "value": 207.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5786302 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 - ns/op",
            "value": 207.9,
            "unit": "ns/op",
            "extra": "5786302 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5786302 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5786302 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0",
            "value": 208,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5780277 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 - ns/op",
            "value": 208,
            "unit": "ns/op",
            "extra": "5780277 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5780277 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5780277 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0",
            "value": 207.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5773162 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 - ns/op",
            "value": 207.7,
            "unit": "ns/op",
            "extra": "5773162 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5773162 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5773162 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0",
            "value": 208.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5790459 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 - ns/op",
            "value": 208.6,
            "unit": "ns/op",
            "extra": "5790459 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5790459 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5790459 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0",
            "value": 210.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5751585 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 - ns/op",
            "value": 210.5,
            "unit": "ns/op",
            "extra": "5751585 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5751585 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5751585 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1",
            "value": 208.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5749809 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 - ns/op",
            "value": 208.9,
            "unit": "ns/op",
            "extra": "5749809 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5749809 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5749809 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1",
            "value": 208.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5750853 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 - ns/op",
            "value": 208.8,
            "unit": "ns/op",
            "extra": "5750853 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5750853 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5750853 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1",
            "value": 209.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5778616 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 - ns/op",
            "value": 209.2,
            "unit": "ns/op",
            "extra": "5778616 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5778616 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5778616 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1",
            "value": 208,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5796555 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 - ns/op",
            "value": 208,
            "unit": "ns/op",
            "extra": "5796555 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5796555 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5796555 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1",
            "value": 207.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5781722 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 - ns/op",
            "value": 207.9,
            "unit": "ns/op",
            "extra": "5781722 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5781722 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5781722 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0",
            "value": 144.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8399944 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 - ns/op",
            "value": 144.1,
            "unit": "ns/op",
            "extra": "8399944 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8399944 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8399944 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0",
            "value": 143.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8323638 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 - ns/op",
            "value": 143.4,
            "unit": "ns/op",
            "extra": "8323638 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8323638 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8323638 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0",
            "value": 145.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8427981 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 - ns/op",
            "value": 145.2,
            "unit": "ns/op",
            "extra": "8427981 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8427981 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8427981 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0",
            "value": 143.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8361003 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 - ns/op",
            "value": 143.6,
            "unit": "ns/op",
            "extra": "8361003 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8361003 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8361003 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0",
            "value": 143.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8365670 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 - ns/op",
            "value": 143.3,
            "unit": "ns/op",
            "extra": "8365670 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8365670 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8365670 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1",
            "value": 144.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8352651 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 - ns/op",
            "value": 144.9,
            "unit": "ns/op",
            "extra": "8352651 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8352651 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8352651 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1",
            "value": 144.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8267613 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 - ns/op",
            "value": 144.2,
            "unit": "ns/op",
            "extra": "8267613 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8267613 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8267613 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1",
            "value": 145.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8347878 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 - ns/op",
            "value": 145.4,
            "unit": "ns/op",
            "extra": "8347878 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8347878 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8347878 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1",
            "value": 145.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8341188 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 - ns/op",
            "value": 145.5,
            "unit": "ns/op",
            "extra": "8341188 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8341188 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8341188 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1",
            "value": 142.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8363323 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 - ns/op",
            "value": 142.9,
            "unit": "ns/op",
            "extra": "8363323 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8363323 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8363323 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0",
            "value": 220.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5418678 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 - ns/op",
            "value": 220.4,
            "unit": "ns/op",
            "extra": "5418678 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5418678 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5418678 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0",
            "value": 228,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5432013 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 - ns/op",
            "value": 228,
            "unit": "ns/op",
            "extra": "5432013 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5432013 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5432013 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0",
            "value": 221.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5443982 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 - ns/op",
            "value": 221.4,
            "unit": "ns/op",
            "extra": "5443982 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5443982 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5443982 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0",
            "value": 221.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5428704 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 - ns/op",
            "value": 221.2,
            "unit": "ns/op",
            "extra": "5428704 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5428704 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5428704 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0",
            "value": 222.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4930062 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 - ns/op",
            "value": 222.2,
            "unit": "ns/op",
            "extra": "4930062 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4930062 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4930062 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1",
            "value": 221.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5392315 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 - ns/op",
            "value": 221.1,
            "unit": "ns/op",
            "extra": "5392315 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5392315 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5392315 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1",
            "value": 221.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5451300 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 - ns/op",
            "value": 221.3,
            "unit": "ns/op",
            "extra": "5451300 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5451300 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5451300 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1",
            "value": 225.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5109862 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 - ns/op",
            "value": 225.4,
            "unit": "ns/op",
            "extra": "5109862 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5109862 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5109862 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1",
            "value": 224.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5437371 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 - ns/op",
            "value": 224.2,
            "unit": "ns/op",
            "extra": "5437371 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5437371 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5437371 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1",
            "value": 221,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5430980 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 - ns/op",
            "value": 221,
            "unit": "ns/op",
            "extra": "5430980 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5430980 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5430980 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0",
            "value": 154.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7969003 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 - ns/op",
            "value": 154.1,
            "unit": "ns/op",
            "extra": "7969003 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7969003 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7969003 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0",
            "value": 151.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7881084 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 - ns/op",
            "value": 151.9,
            "unit": "ns/op",
            "extra": "7881084 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7881084 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7881084 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0",
            "value": 153.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7955359 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 - ns/op",
            "value": 153.2,
            "unit": "ns/op",
            "extra": "7955359 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7955359 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7955359 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0",
            "value": 154.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7910066 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 - ns/op",
            "value": 154.1,
            "unit": "ns/op",
            "extra": "7910066 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7910066 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7910066 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0",
            "value": 154,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7905908 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 - ns/op",
            "value": 154,
            "unit": "ns/op",
            "extra": "7905908 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7905908 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7905908 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1",
            "value": 152.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7941379 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 - ns/op",
            "value": 152.3,
            "unit": "ns/op",
            "extra": "7941379 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7941379 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7941379 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1",
            "value": 152.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7960562 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 - ns/op",
            "value": 152.5,
            "unit": "ns/op",
            "extra": "7960562 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7960562 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7960562 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1",
            "value": 153.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7882656 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 - ns/op",
            "value": 153.2,
            "unit": "ns/op",
            "extra": "7882656 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7882656 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7882656 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1",
            "value": 153.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7872010 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 - ns/op",
            "value": 153.5,
            "unit": "ns/op",
            "extra": "7872010 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7872010 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7872010 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1",
            "value": 152.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7927310 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 - ns/op",
            "value": 152.4,
            "unit": "ns/op",
            "extra": "7927310 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7927310 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7927310 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0",
            "value": 235,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5002886 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 - ns/op",
            "value": 235,
            "unit": "ns/op",
            "extra": "5002886 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5002886 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5002886 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0",
            "value": 235.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5095000 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 - ns/op",
            "value": 235.8,
            "unit": "ns/op",
            "extra": "5095000 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5095000 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5095000 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0",
            "value": 236.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5072892 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 - ns/op",
            "value": 236.2,
            "unit": "ns/op",
            "extra": "5072892 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5072892 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5072892 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0",
            "value": 236.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5034711 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 - ns/op",
            "value": 236.6,
            "unit": "ns/op",
            "extra": "5034711 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5034711 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5034711 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0",
            "value": 236.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5074281 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 - ns/op",
            "value": 236.9,
            "unit": "ns/op",
            "extra": "5074281 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5074281 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5074281 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1",
            "value": 236.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5095260 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 - ns/op",
            "value": 236.2,
            "unit": "ns/op",
            "extra": "5095260 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5095260 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5095260 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1",
            "value": 237.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5094903 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 - ns/op",
            "value": 237.5,
            "unit": "ns/op",
            "extra": "5094903 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5094903 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5094903 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1",
            "value": 235.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5109420 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 - ns/op",
            "value": 235.8,
            "unit": "ns/op",
            "extra": "5109420 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5109420 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5109420 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1",
            "value": 236.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5078353 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 - ns/op",
            "value": 236.5,
            "unit": "ns/op",
            "extra": "5078353 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5078353 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5078353 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1",
            "value": 237.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5098536 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 - ns/op",
            "value": 237.1,
            "unit": "ns/op",
            "extra": "5098536 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5098536 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5098536 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0",
            "value": 159.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7573322 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 - ns/op",
            "value": 159.1,
            "unit": "ns/op",
            "extra": "7573322 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7573322 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7573322 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0",
            "value": 158.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7534489 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 - ns/op",
            "value": 158.8,
            "unit": "ns/op",
            "extra": "7534489 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7534489 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7534489 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0",
            "value": 160.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7588291 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 - ns/op",
            "value": 160.4,
            "unit": "ns/op",
            "extra": "7588291 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7588291 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7588291 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0",
            "value": 159.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7530684 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 - ns/op",
            "value": 159.1,
            "unit": "ns/op",
            "extra": "7530684 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7530684 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7530684 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0",
            "value": 158.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7535931 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 - ns/op",
            "value": 158.4,
            "unit": "ns/op",
            "extra": "7535931 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7535931 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7535931 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1",
            "value": 160.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7570906 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 - ns/op",
            "value": 160.3,
            "unit": "ns/op",
            "extra": "7570906 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7570906 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7570906 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1",
            "value": 159.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7478682 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 - ns/op",
            "value": 159.7,
            "unit": "ns/op",
            "extra": "7478682 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7478682 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7478682 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1",
            "value": 157.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7605842 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 - ns/op",
            "value": 157.8,
            "unit": "ns/op",
            "extra": "7605842 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7605842 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7605842 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1",
            "value": 159.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7623416 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 - ns/op",
            "value": 159.2,
            "unit": "ns/op",
            "extra": "7623416 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7623416 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7623416 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1",
            "value": 158.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7524444 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 - ns/op",
            "value": 158.4,
            "unit": "ns/op",
            "extra": "7524444 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7524444 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7524444 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0",
            "value": 257.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4779862 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 - ns/op",
            "value": 257.7,
            "unit": "ns/op",
            "extra": "4779862 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4779862 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4779862 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0",
            "value": 250.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4756012 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 - ns/op",
            "value": 250.3,
            "unit": "ns/op",
            "extra": "4756012 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4756012 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4756012 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0",
            "value": 248.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4817145 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 - ns/op",
            "value": 248.7,
            "unit": "ns/op",
            "extra": "4817145 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4817145 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4817145 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0",
            "value": 254.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4823110 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 - ns/op",
            "value": 254.9,
            "unit": "ns/op",
            "extra": "4823110 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4823110 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4823110 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0",
            "value": 248.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4834278 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 - ns/op",
            "value": 248.7,
            "unit": "ns/op",
            "extra": "4834278 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4834278 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4834278 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1",
            "value": 248.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4836904 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 - ns/op",
            "value": 248.6,
            "unit": "ns/op",
            "extra": "4836904 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4836904 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4836904 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1",
            "value": 250.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4826624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 - ns/op",
            "value": 250.1,
            "unit": "ns/op",
            "extra": "4826624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4826624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4826624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1",
            "value": 251.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4827946 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 - ns/op",
            "value": 251.5,
            "unit": "ns/op",
            "extra": "4827946 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4827946 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4827946 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1",
            "value": 249.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4831152 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 - ns/op",
            "value": 249.8,
            "unit": "ns/op",
            "extra": "4831152 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4831152 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4831152 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1",
            "value": 250.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4827130 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 - ns/op",
            "value": 250.4,
            "unit": "ns/op",
            "extra": "4827130 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4827130 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4827130 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0",
            "value": 173.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6858195 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 - ns/op",
            "value": 173.1,
            "unit": "ns/op",
            "extra": "6858195 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6858195 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6858195 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0",
            "value": 174.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6919281 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 - ns/op",
            "value": 174.4,
            "unit": "ns/op",
            "extra": "6919281 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6919281 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6919281 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0",
            "value": 176,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6740846 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 - ns/op",
            "value": 176,
            "unit": "ns/op",
            "extra": "6740846 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6740846 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6740846 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0",
            "value": 173.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6930778 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 - ns/op",
            "value": 173.7,
            "unit": "ns/op",
            "extra": "6930778 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6930778 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6930778 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0",
            "value": 179.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6933590 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 - ns/op",
            "value": 179.2,
            "unit": "ns/op",
            "extra": "6933590 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6933590 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6933590 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1",
            "value": 174.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6833484 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 - ns/op",
            "value": 174.1,
            "unit": "ns/op",
            "extra": "6833484 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6833484 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6833484 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1",
            "value": 176,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6839727 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 - ns/op",
            "value": 176,
            "unit": "ns/op",
            "extra": "6839727 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6839727 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6839727 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1",
            "value": 176.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6840873 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 - ns/op",
            "value": 176.1,
            "unit": "ns/op",
            "extra": "6840873 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6840873 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6840873 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1",
            "value": 179.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6932846 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 - ns/op",
            "value": 179.3,
            "unit": "ns/op",
            "extra": "6932846 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6932846 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6932846 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1",
            "value": 180.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6916624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 - ns/op",
            "value": 180.9,
            "unit": "ns/op",
            "extra": "6916624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6916624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6916624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0",
            "value": 281.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4261117 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 - ns/op",
            "value": 281.8,
            "unit": "ns/op",
            "extra": "4261117 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4261117 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4261117 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0",
            "value": 282.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4261206 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 - ns/op",
            "value": 282.6,
            "unit": "ns/op",
            "extra": "4261206 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4261206 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4261206 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0",
            "value": 280.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4283091 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 - ns/op",
            "value": 280.9,
            "unit": "ns/op",
            "extra": "4283091 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4283091 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4283091 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0",
            "value": 280.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4294932 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 - ns/op",
            "value": 280.5,
            "unit": "ns/op",
            "extra": "4294932 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4294932 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4294932 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0",
            "value": 280.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4279095 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 - ns/op",
            "value": 280.3,
            "unit": "ns/op",
            "extra": "4279095 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4279095 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4279095 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1",
            "value": 280.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4275891 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 - ns/op",
            "value": 280.4,
            "unit": "ns/op",
            "extra": "4275891 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4275891 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4275891 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1",
            "value": 282.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4310948 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 - ns/op",
            "value": 282.8,
            "unit": "ns/op",
            "extra": "4310948 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4310948 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4310948 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1",
            "value": 280,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4300677 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 - ns/op",
            "value": 280,
            "unit": "ns/op",
            "extra": "4300677 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4300677 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4300677 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1",
            "value": 280.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4143422 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 - ns/op",
            "value": 280.9,
            "unit": "ns/op",
            "extra": "4143422 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4143422 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4143422 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1",
            "value": 280.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "4275466 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 - ns/op",
            "value": 280.9,
            "unit": "ns/op",
            "extra": "4275466 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "4275466 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "4275466 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0",
            "value": 83.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14433673 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 - ns/op",
            "value": 83.1,
            "unit": "ns/op",
            "extra": "14433673 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14433673 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14433673 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0",
            "value": 82.87,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14392021 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 - ns/op",
            "value": 82.87,
            "unit": "ns/op",
            "extra": "14392021 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14392021 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14392021 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0",
            "value": 82.93,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14368720 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 - ns/op",
            "value": 82.93,
            "unit": "ns/op",
            "extra": "14368720 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14368720 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14368720 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0",
            "value": 83.09,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14315806 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 - ns/op",
            "value": 83.09,
            "unit": "ns/op",
            "extra": "14315806 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14315806 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14315806 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0",
            "value": 83.45,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14446684 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 - ns/op",
            "value": 83.45,
            "unit": "ns/op",
            "extra": "14446684 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14446684 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14446684 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1",
            "value": 83.23,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12487554 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 - ns/op",
            "value": 83.23,
            "unit": "ns/op",
            "extra": "12487554 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12487554 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12487554 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1",
            "value": 82.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14491395 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 - ns/op",
            "value": 82.9,
            "unit": "ns/op",
            "extra": "14491395 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14491395 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14491395 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1",
            "value": 82.94,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14383528 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 - ns/op",
            "value": 82.94,
            "unit": "ns/op",
            "extra": "14383528 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14383528 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14383528 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1",
            "value": 82.75,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14268015 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 - ns/op",
            "value": 82.75,
            "unit": "ns/op",
            "extra": "14268015 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14268015 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14268015 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1",
            "value": 82.97,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14409668 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 - ns/op",
            "value": 82.97,
            "unit": "ns/op",
            "extra": "14409668 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14409668 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14409668 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0",
            "value": 88.58,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13532083 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 - ns/op",
            "value": 88.58,
            "unit": "ns/op",
            "extra": "13532083 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13532083 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13532083 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0",
            "value": 88.47,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13498832 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 - ns/op",
            "value": 88.47,
            "unit": "ns/op",
            "extra": "13498832 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13498832 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13498832 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0",
            "value": 88.74,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13505058 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 - ns/op",
            "value": 88.74,
            "unit": "ns/op",
            "extra": "13505058 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13505058 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13505058 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0",
            "value": 89.55,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13450570 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 - ns/op",
            "value": 89.55,
            "unit": "ns/op",
            "extra": "13450570 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13450570 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13450570 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0",
            "value": 88.77,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13269656 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 - ns/op",
            "value": 88.77,
            "unit": "ns/op",
            "extra": "13269656 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13269656 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13269656 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1",
            "value": 88.69,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13503234 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 - ns/op",
            "value": 88.69,
            "unit": "ns/op",
            "extra": "13503234 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13503234 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13503234 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1",
            "value": 88.69,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13590438 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 - ns/op",
            "value": 88.69,
            "unit": "ns/op",
            "extra": "13590438 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13590438 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13590438 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1",
            "value": 88.74,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13485568 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 - ns/op",
            "value": 88.74,
            "unit": "ns/op",
            "extra": "13485568 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13485568 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13485568 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1",
            "value": 88.59,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13560970 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 - ns/op",
            "value": 88.59,
            "unit": "ns/op",
            "extra": "13560970 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13560970 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13560970 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1",
            "value": 90.57,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13517916 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 - ns/op",
            "value": 90.57,
            "unit": "ns/op",
            "extra": "13517916 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13517916 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13517916 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0",
            "value": 93.82,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12792255 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 - ns/op",
            "value": 93.82,
            "unit": "ns/op",
            "extra": "12792255 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12792255 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12792255 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0",
            "value": 94.04,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12167496 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 - ns/op",
            "value": 94.04,
            "unit": "ns/op",
            "extra": "12167496 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12167496 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12167496 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0",
            "value": 93.77,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12746578 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 - ns/op",
            "value": 93.77,
            "unit": "ns/op",
            "extra": "12746578 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12746578 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12746578 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0",
            "value": 95.47,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12712514 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 - ns/op",
            "value": 95.47,
            "unit": "ns/op",
            "extra": "12712514 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12712514 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12712514 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0",
            "value": 93.85,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12690840 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 - ns/op",
            "value": 93.85,
            "unit": "ns/op",
            "extra": "12690840 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12690840 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12690840 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1",
            "value": 94.03,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12822921 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 - ns/op",
            "value": 94.03,
            "unit": "ns/op",
            "extra": "12822921 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12822921 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12822921 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1",
            "value": 94.06,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12834709 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 - ns/op",
            "value": 94.06,
            "unit": "ns/op",
            "extra": "12834709 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12834709 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12834709 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1",
            "value": 93.74,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12385254 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 - ns/op",
            "value": 93.74,
            "unit": "ns/op",
            "extra": "12385254 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12385254 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12385254 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1",
            "value": 93.88,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12835242 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 - ns/op",
            "value": 93.88,
            "unit": "ns/op",
            "extra": "12835242 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12835242 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12835242 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1",
            "value": 93.91,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12780032 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 - ns/op",
            "value": 93.91,
            "unit": "ns/op",
            "extra": "12780032 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12780032 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12780032 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0",
            "value": 101.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11989232 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 - ns/op",
            "value": 101.3,
            "unit": "ns/op",
            "extra": "11989232 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11989232 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11989232 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0",
            "value": 99.99,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11943493 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 - ns/op",
            "value": 99.99,
            "unit": "ns/op",
            "extra": "11943493 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11943493 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11943493 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0",
            "value": 99.84,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11955578 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 - ns/op",
            "value": 99.84,
            "unit": "ns/op",
            "extra": "11955578 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11955578 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11955578 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0",
            "value": 100.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12025273 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 - ns/op",
            "value": 100.3,
            "unit": "ns/op",
            "extra": "12025273 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12025273 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12025273 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0",
            "value": 100.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11744622 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 - ns/op",
            "value": 100.5,
            "unit": "ns/op",
            "extra": "11744622 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11744622 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11744622 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1",
            "value": 99.79,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12023128 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 - ns/op",
            "value": 99.79,
            "unit": "ns/op",
            "extra": "12023128 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12023128 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12023128 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1",
            "value": 99.84,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11977346 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 - ns/op",
            "value": 99.84,
            "unit": "ns/op",
            "extra": "11977346 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11977346 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11977346 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1",
            "value": 99.92,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12017083 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 - ns/op",
            "value": 99.92,
            "unit": "ns/op",
            "extra": "12017083 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12017083 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12017083 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1",
            "value": 99.86,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12071625 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 - ns/op",
            "value": 99.86,
            "unit": "ns/op",
            "extra": "12071625 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12071625 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12071625 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1",
            "value": 100.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12012333 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 - ns/op",
            "value": 100.1,
            "unit": "ns/op",
            "extra": "12012333 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12012333 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12012333 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0",
            "value": 116.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10282582 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 - ns/op",
            "value": 116.1,
            "unit": "ns/op",
            "extra": "10282582 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10282582 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10282582 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0",
            "value": 116.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10376334 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 - ns/op",
            "value": 116.3,
            "unit": "ns/op",
            "extra": "10376334 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10376334 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10376334 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0",
            "value": 116.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10208674 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 - ns/op",
            "value": 116.2,
            "unit": "ns/op",
            "extra": "10208674 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10208674 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10208674 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0",
            "value": 116.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10298828 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 - ns/op",
            "value": 116.7,
            "unit": "ns/op",
            "extra": "10298828 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10298828 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10298828 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0",
            "value": 116.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10282780 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 - ns/op",
            "value": 116.2,
            "unit": "ns/op",
            "extra": "10282780 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10282780 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10282780 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1",
            "value": 116.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10405814 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 - ns/op",
            "value": 116.2,
            "unit": "ns/op",
            "extra": "10405814 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10405814 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10405814 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1",
            "value": 116.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10346498 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 - ns/op",
            "value": 116.3,
            "unit": "ns/op",
            "extra": "10346498 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10346498 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10346498 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1",
            "value": 116,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10289692 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 - ns/op",
            "value": 116,
            "unit": "ns/op",
            "extra": "10289692 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10289692 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10289692 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1",
            "value": 116.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10310028 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 - ns/op",
            "value": 116.1,
            "unit": "ns/op",
            "extra": "10310028 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10310028 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10310028 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1",
            "value": 116.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10364242 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 - ns/op",
            "value": 116.2,
            "unit": "ns/op",
            "extra": "10364242 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10364242 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10364242 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0",
            "value": 153.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7820160 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 - ns/op",
            "value": 153.3,
            "unit": "ns/op",
            "extra": "7820160 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7820160 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7820160 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0",
            "value": 152.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7795558 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 - ns/op",
            "value": 152.7,
            "unit": "ns/op",
            "extra": "7795558 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7795558 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7795558 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0",
            "value": 153.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7831783 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 - ns/op",
            "value": 153.4,
            "unit": "ns/op",
            "extra": "7831783 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7831783 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7831783 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0",
            "value": 153.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7819297 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 - ns/op",
            "value": 153.5,
            "unit": "ns/op",
            "extra": "7819297 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7819297 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7819297 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0",
            "value": 155.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7801302 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 - ns/op",
            "value": 155.8,
            "unit": "ns/op",
            "extra": "7801302 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7801302 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7801302 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1",
            "value": 153,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7850829 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 - ns/op",
            "value": 153,
            "unit": "ns/op",
            "extra": "7850829 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7850829 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7850829 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1",
            "value": 153.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7813304 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 - ns/op",
            "value": 153.7,
            "unit": "ns/op",
            "extra": "7813304 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7813304 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7813304 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1",
            "value": 154.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7814287 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 - ns/op",
            "value": 154.1,
            "unit": "ns/op",
            "extra": "7814287 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7814287 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7814287 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1",
            "value": 154.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7827500 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 - ns/op",
            "value": 154.4,
            "unit": "ns/op",
            "extra": "7827500 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7827500 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7827500 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1",
            "value": 155.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7810573 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 - ns/op",
            "value": 155.6,
            "unit": "ns/op",
            "extra": "7810573 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7810573 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7810573 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0",
            "value": 260142,
            "unit": "ns/op\t   38488 B/op\t    1480 allocs/op",
            "extra": "4376 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0 - ns/op",
            "value": 260142,
            "unit": "ns/op",
            "extra": "4376 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0 - B/op",
            "value": 38488,
            "unit": "B/op",
            "extra": "4376 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0 - allocs/op",
            "value": 1480,
            "unit": "allocs/op",
            "extra": "4376 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0",
            "value": 253137,
            "unit": "ns/op\t   38488 B/op\t    1480 allocs/op",
            "extra": "4448 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0 - ns/op",
            "value": 253137,
            "unit": "ns/op",
            "extra": "4448 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0 - B/op",
            "value": 38488,
            "unit": "B/op",
            "extra": "4448 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0 - allocs/op",
            "value": 1480,
            "unit": "allocs/op",
            "extra": "4448 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0",
            "value": 252757,
            "unit": "ns/op\t   38488 B/op\t    1480 allocs/op",
            "extra": "4474 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0 - ns/op",
            "value": 252757,
            "unit": "ns/op",
            "extra": "4474 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0 - B/op",
            "value": 38488,
            "unit": "B/op",
            "extra": "4474 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0 - allocs/op",
            "value": 1480,
            "unit": "allocs/op",
            "extra": "4474 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0",
            "value": 252323,
            "unit": "ns/op\t   38488 B/op\t    1480 allocs/op",
            "extra": "4442 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0 - ns/op",
            "value": 252323,
            "unit": "ns/op",
            "extra": "4442 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0 - B/op",
            "value": 38488,
            "unit": "B/op",
            "extra": "4442 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0 - allocs/op",
            "value": 1480,
            "unit": "allocs/op",
            "extra": "4442 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0",
            "value": 252904,
            "unit": "ns/op\t   38488 B/op\t    1480 allocs/op",
            "extra": "4431 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0 - ns/op",
            "value": 252904,
            "unit": "ns/op",
            "extra": "4431 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0 - B/op",
            "value": 38488,
            "unit": "B/op",
            "extra": "4431 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O0 - allocs/op",
            "value": 1480,
            "unit": "allocs/op",
            "extra": "4431 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1",
            "value": 233728,
            "unit": "ns/op\t   38488 B/op\t    1480 allocs/op",
            "extra": "4716 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1 - ns/op",
            "value": 233728,
            "unit": "ns/op",
            "extra": "4716 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1 - B/op",
            "value": 38488,
            "unit": "B/op",
            "extra": "4716 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1 - allocs/op",
            "value": 1480,
            "unit": "allocs/op",
            "extra": "4716 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1",
            "value": 233960,
            "unit": "ns/op\t   38488 B/op\t    1480 allocs/op",
            "extra": "4796 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1 - ns/op",
            "value": 233960,
            "unit": "ns/op",
            "extra": "4796 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1 - B/op",
            "value": 38488,
            "unit": "B/op",
            "extra": "4796 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1 - allocs/op",
            "value": 1480,
            "unit": "allocs/op",
            "extra": "4796 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1",
            "value": 233803,
            "unit": "ns/op\t   38488 B/op\t    1480 allocs/op",
            "extra": "4749 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1 - ns/op",
            "value": 233803,
            "unit": "ns/op",
            "extra": "4749 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1 - B/op",
            "value": 38488,
            "unit": "B/op",
            "extra": "4749 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1 - allocs/op",
            "value": 1480,
            "unit": "allocs/op",
            "extra": "4749 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1",
            "value": 234445,
            "unit": "ns/op\t   38488 B/op\t    1480 allocs/op",
            "extra": "4322 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1 - ns/op",
            "value": 234445,
            "unit": "ns/op",
            "extra": "4322 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1 - B/op",
            "value": 38488,
            "unit": "B/op",
            "extra": "4322 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1 - allocs/op",
            "value": 1480,
            "unit": "allocs/op",
            "extra": "4322 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1",
            "value": 234337,
            "unit": "ns/op\t   38488 B/op\t    1480 allocs/op",
            "extra": "4657 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1 - ns/op",
            "value": 234337,
            "unit": "ns/op",
            "extra": "4657 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1 - B/op",
            "value": 38488,
            "unit": "B/op",
            "extra": "4657 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_SingleCell_O1 - allocs/op",
            "value": 1480,
            "unit": "allocs/op",
            "extra": "4657 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0",
            "value": 627651,
            "unit": "ns/op\t   97890 B/op\t    4383 allocs/op",
            "extra": "1902 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0 - ns/op",
            "value": 627651,
            "unit": "ns/op",
            "extra": "1902 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0 - B/op",
            "value": 97890,
            "unit": "B/op",
            "extra": "1902 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0 - allocs/op",
            "value": 4383,
            "unit": "allocs/op",
            "extra": "1902 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0",
            "value": 627300,
            "unit": "ns/op\t   97890 B/op\t    4383 allocs/op",
            "extra": "1831 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0 - ns/op",
            "value": 627300,
            "unit": "ns/op",
            "extra": "1831 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0 - B/op",
            "value": 97890,
            "unit": "B/op",
            "extra": "1831 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0 - allocs/op",
            "value": 4383,
            "unit": "allocs/op",
            "extra": "1831 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0",
            "value": 632921,
            "unit": "ns/op\t   97890 B/op\t    4383 allocs/op",
            "extra": "1836 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0 - ns/op",
            "value": 632921,
            "unit": "ns/op",
            "extra": "1836 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0 - B/op",
            "value": 97890,
            "unit": "B/op",
            "extra": "1836 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0 - allocs/op",
            "value": 4383,
            "unit": "allocs/op",
            "extra": "1836 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0",
            "value": 628455,
            "unit": "ns/op\t   97890 B/op\t    4383 allocs/op",
            "extra": "1777 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0 - ns/op",
            "value": 628455,
            "unit": "ns/op",
            "extra": "1777 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0 - B/op",
            "value": 97890,
            "unit": "B/op",
            "extra": "1777 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0 - allocs/op",
            "value": 4383,
            "unit": "allocs/op",
            "extra": "1777 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0",
            "value": 624320,
            "unit": "ns/op\t   97890 B/op\t    4383 allocs/op",
            "extra": "1916 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0 - ns/op",
            "value": 624320,
            "unit": "ns/op",
            "extra": "1916 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0 - B/op",
            "value": 97890,
            "unit": "B/op",
            "extra": "1916 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O0 - allocs/op",
            "value": 4383,
            "unit": "allocs/op",
            "extra": "1916 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1",
            "value": 594580,
            "unit": "ns/op\t   97890 B/op\t    4383 allocs/op",
            "extra": "1957 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1 - ns/op",
            "value": 594580,
            "unit": "ns/op",
            "extra": "1957 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1 - B/op",
            "value": 97890,
            "unit": "B/op",
            "extra": "1957 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1 - allocs/op",
            "value": 4383,
            "unit": "allocs/op",
            "extra": "1957 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1",
            "value": 590722,
            "unit": "ns/op\t   97890 B/op\t    4383 allocs/op",
            "extra": "1940 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1 - ns/op",
            "value": 590722,
            "unit": "ns/op",
            "extra": "1940 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1 - B/op",
            "value": 97890,
            "unit": "B/op",
            "extra": "1940 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1 - allocs/op",
            "value": 4383,
            "unit": "allocs/op",
            "extra": "1940 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1",
            "value": 592436,
            "unit": "ns/op\t   97890 B/op\t    4383 allocs/op",
            "extra": "1924 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1 - ns/op",
            "value": 592436,
            "unit": "ns/op",
            "extra": "1924 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1 - B/op",
            "value": 97890,
            "unit": "B/op",
            "extra": "1924 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1 - allocs/op",
            "value": 4383,
            "unit": "allocs/op",
            "extra": "1924 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1",
            "value": 591546,
            "unit": "ns/op\t   97890 B/op\t    4383 allocs/op",
            "extra": "1936 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1 - ns/op",
            "value": 591546,
            "unit": "ns/op",
            "extra": "1936 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1 - B/op",
            "value": 97890,
            "unit": "B/op",
            "extra": "1936 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1 - allocs/op",
            "value": 4383,
            "unit": "allocs/op",
            "extra": "1936 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1",
            "value": 590062,
            "unit": "ns/op\t   97890 B/op\t    4383 allocs/op",
            "extra": "1980 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1 - ns/op",
            "value": 590062,
            "unit": "ns/op",
            "extra": "1980 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1 - B/op",
            "value": 97890,
            "unit": "B/op",
            "extra": "1980 times\n4 procs"
          },
          {
            "name": "BenchmarkVarCapture_MultiCell_O1 - allocs/op",
            "value": 4383,
            "unit": "allocs/op",
            "extra": "1980 times\n4 procs"
          }
        ]
      }
    ]
  }
}