window.BENCHMARK_DATA = {
  "lastUpdate": 1773253666663,
  "repoUrl": "https://github.com/MontFerret/ferret",
  "entries": {
    "Ferret Go Benchmarks": [
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
          "id": "8549314f329132441685a746aa2e713e1673fcdf",
          "message": "chore: add benchmark workflow and enhance Makefile with benchmarking options (#891)",
          "timestamp": "2026-03-11T11:43:00-04:00",
          "tree_id": "ce5cdf1b28c29871659f8ae65fcc837a377c7b94",
          "url": "https://github.com/MontFerret/ferret/commit/8549314f329132441685a746aa2e713e1673fcdf"
        },
        "date": 1773244721552,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3551,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "333066 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3551,
            "unit": "ns/op",
            "extra": "333066 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "333066 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "333066 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 358.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3343921 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 358.4,
            "unit": "ns/op",
            "extra": "3343921 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3343921 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3343921 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2545,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "462319 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2545,
            "unit": "ns/op",
            "extra": "462319 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "462319 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "462319 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/vm/internal/data)",
            "value": 966.1,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1255404 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/vm/internal/data) - ns/op",
            "value": 966.1,
            "unit": "ns/op",
            "extra": "1255404 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/vm/internal/data) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1255404 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/vm/internal/data) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1255404 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 105852,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 105852,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 101853,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 101853,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 112186,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "9747 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 112186,
            "unit": "ns/op",
            "extra": "9747 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "9747 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "9747 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 107413,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 107413,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 146012,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7969 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 146012,
            "unit": "ns/op",
            "extra": "7969 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7969 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7969 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 140508,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "8347 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 140508,
            "unit": "ns/op",
            "extra": "8347 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "8347 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "8347 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 149074,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7904 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 149074,
            "unit": "ns/op",
            "extra": "7904 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7904 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7904 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 145462,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "8120 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 145462,
            "unit": "ns/op",
            "extra": "8120 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "8120 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "8120 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 193489,
            "unit": "ns/op\t   79002 B/op\t    3501 allocs/op",
            "extra": "6697 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 193489,
            "unit": "ns/op",
            "extra": "6697 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 79002,
            "unit": "B/op",
            "extra": "6697 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "6697 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 174594,
            "unit": "ns/op\t   79002 B/op\t    3501 allocs/op",
            "extra": "6658 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 174594,
            "unit": "ns/op",
            "extra": "6658 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 79002,
            "unit": "B/op",
            "extra": "6658 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "6658 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 205373,
            "unit": "ns/op\t   70322 B/op\t    3501 allocs/op",
            "extra": "5715 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 205373,
            "unit": "ns/op",
            "extra": "5715 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 70322,
            "unit": "B/op",
            "extra": "5715 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5715 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 206156,
            "unit": "ns/op\t   70322 B/op\t    3501 allocs/op",
            "extra": "5704 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 206156,
            "unit": "ns/op",
            "extra": "5704 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 70322,
            "unit": "B/op",
            "extra": "5704 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5704 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 4532,
            "unit": "ns/op\t    2040 B/op\t      55 allocs/op",
            "extra": "263728 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 4532,
            "unit": "ns/op",
            "extra": "263728 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 2040,
            "unit": "B/op",
            "extra": "263728 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 55,
            "unit": "allocs/op",
            "extra": "263728 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 4349,
            "unit": "ns/op\t    2040 B/op\t      55 allocs/op",
            "extra": "245346 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 4349,
            "unit": "ns/op",
            "extra": "245346 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 2040,
            "unit": "B/op",
            "extra": "245346 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 55,
            "unit": "allocs/op",
            "extra": "245346 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2428686,
            "unit": "ns/op\t  157196 B/op\t   19512 allocs/op",
            "extra": "488 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2428686,
            "unit": "ns/op",
            "extra": "488 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 157196,
            "unit": "B/op",
            "extra": "488 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 19512,
            "unit": "allocs/op",
            "extra": "488 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2425904,
            "unit": "ns/op\t  157192 B/op\t   19512 allocs/op",
            "extra": "501 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2425904,
            "unit": "ns/op",
            "extra": "501 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 157192,
            "unit": "B/op",
            "extra": "501 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 19512,
            "unit": "allocs/op",
            "extra": "501 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 3965991,
            "unit": "ns/op\t 1782991 B/op\t   49527 allocs/op",
            "extra": "304 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 3965991,
            "unit": "ns/op",
            "extra": "304 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1782991,
            "unit": "B/op",
            "extra": "304 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 49527,
            "unit": "allocs/op",
            "extra": "304 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 3962120,
            "unit": "ns/op\t 1782983 B/op\t   49527 allocs/op",
            "extra": "302 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 3962120,
            "unit": "ns/op",
            "extra": "302 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1782983,
            "unit": "B/op",
            "extra": "302 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 49527,
            "unit": "allocs/op",
            "extra": "302 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 9943326,
            "unit": "ns/op\t 4144435 B/op\t  127649 allocs/op",
            "extra": "123 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 9943326,
            "unit": "ns/op",
            "extra": "123 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 4144435,
            "unit": "B/op",
            "extra": "123 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 127649,
            "unit": "allocs/op",
            "extra": "123 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 9421681,
            "unit": "ns/op\t 4144402 B/op\t  127649 allocs/op",
            "extra": "126 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 9421681,
            "unit": "ns/op",
            "extra": "126 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 4144402,
            "unit": "B/op",
            "extra": "126 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 127649,
            "unit": "allocs/op",
            "extra": "126 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 98856,
            "unit": "ns/op\t   53497 B/op\t    1050 allocs/op",
            "extra": "12116 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 98856,
            "unit": "ns/op",
            "extra": "12116 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 53497,
            "unit": "B/op",
            "extra": "12116 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1050,
            "unit": "allocs/op",
            "extra": "12116 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 96978,
            "unit": "ns/op\t   53497 B/op\t    1050 allocs/op",
            "extra": "12238 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 96978,
            "unit": "ns/op",
            "extra": "12238 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 53497,
            "unit": "B/op",
            "extra": "12238 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1050,
            "unit": "allocs/op",
            "extra": "12238 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2679029,
            "unit": "ns/op\t 1782643 B/op\t   49521 allocs/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2679029,
            "unit": "ns/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1782643,
            "unit": "B/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 49521,
            "unit": "allocs/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2740379,
            "unit": "ns/op\t 1782647 B/op\t   49521 allocs/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2740379,
            "unit": "ns/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1782647,
            "unit": "B/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 49521,
            "unit": "allocs/op",
            "extra": "440 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 78759,
            "unit": "ns/op\t   40697 B/op\t     850 allocs/op",
            "extra": "15111 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 78759,
            "unit": "ns/op",
            "extra": "15111 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 40697,
            "unit": "B/op",
            "extra": "15111 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 850,
            "unit": "allocs/op",
            "extra": "15111 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 78576,
            "unit": "ns/op\t   40697 B/op\t     850 allocs/op",
            "extra": "15321 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 78576,
            "unit": "ns/op",
            "extra": "15321 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 40697,
            "unit": "B/op",
            "extra": "15321 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 850,
            "unit": "allocs/op",
            "extra": "15321 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 80852,
            "unit": "ns/op\t   40697 B/op\t     850 allocs/op",
            "extra": "14838 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 80852,
            "unit": "ns/op",
            "extra": "14838 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 40697,
            "unit": "B/op",
            "extra": "14838 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 850,
            "unit": "allocs/op",
            "extra": "14838 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 78115,
            "unit": "ns/op\t   40697 B/op\t     850 allocs/op",
            "extra": "15284 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 78115,
            "unit": "ns/op",
            "extra": "15284 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 40697,
            "unit": "B/op",
            "extra": "15284 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 850,
            "unit": "allocs/op",
            "extra": "15284 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 8626,
            "unit": "ns/op\t     528 B/op\t      12 allocs/op",
            "extra": "137941 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 8626,
            "unit": "ns/op",
            "extra": "137941 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 528,
            "unit": "B/op",
            "extra": "137941 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "137941 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 8684,
            "unit": "ns/op\t     528 B/op\t      12 allocs/op",
            "extra": "137046 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 8684,
            "unit": "ns/op",
            "extra": "137046 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 528,
            "unit": "B/op",
            "extra": "137046 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "137046 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2885,
            "unit": "ns/op\t     640 B/op\t      17 allocs/op",
            "extra": "400936 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2885,
            "unit": "ns/op",
            "extra": "400936 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 640,
            "unit": "B/op",
            "extra": "400936 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "400936 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 1267,
            "unit": "ns/op\t     560 B/op\t       7 allocs/op",
            "extra": "864381 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 1267,
            "unit": "ns/op",
            "extra": "864381 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "864381 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "864381 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2476,
            "unit": "ns/op\t    1112 B/op\t      34 allocs/op",
            "extra": "469134 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2476,
            "unit": "ns/op",
            "extra": "469134 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "469134 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "469134 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2422,
            "unit": "ns/op\t    1112 B/op\t      34 allocs/op",
            "extra": "475922 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2422,
            "unit": "ns/op",
            "extra": "475922 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "475922 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "475922 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 179.8,
            "unit": "ns/op\t      96 B/op\t       1 allocs/op",
            "extra": "6644376 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 179.8,
            "unit": "ns/op",
            "extra": "6644376 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "6644376 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "6644376 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 180.9,
            "unit": "ns/op\t      96 B/op\t       1 allocs/op",
            "extra": "6378825 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 180.9,
            "unit": "ns/op",
            "extra": "6378825 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "6378825 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "6378825 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 65.18,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18257530 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 65.18,
            "unit": "ns/op",
            "extra": "18257530 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18257530 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18257530 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 65.42,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18277534 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 65.42,
            "unit": "ns/op",
            "extra": "18277534 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18277534 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18277534 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 65.71,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18300878 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 65.71,
            "unit": "ns/op",
            "extra": "18300878 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18300878 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18300878 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 65.53,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18277437 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 65.53,
            "unit": "ns/op",
            "extra": "18277437 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18277437 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18277437 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 78.32,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15804548 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 78.32,
            "unit": "ns/op",
            "extra": "15804548 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15804548 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15804548 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 78.35,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15586483 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 78.35,
            "unit": "ns/op",
            "extra": "15586483 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15586483 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15586483 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 107.3,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "11066721 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 107.3,
            "unit": "ns/op",
            "extra": "11066721 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "11066721 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "11066721 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 106.2,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "11204660 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 106.2,
            "unit": "ns/op",
            "extra": "11204660 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "11204660 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "11204660 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 84.98,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13783755 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 84.98,
            "unit": "ns/op",
            "extra": "13783755 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13783755 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13783755 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 85.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13863363 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 85.1,
            "unit": "ns/op",
            "extra": "13863363 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13863363 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13863363 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 124.1,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "9952183 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 124.1,
            "unit": "ns/op",
            "extra": "9952183 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "9952183 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9952183 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 120.6,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "9936592 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 120.6,
            "unit": "ns/op",
            "extra": "9936592 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "9936592 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9936592 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 90.91,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13411354 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 90.91,
            "unit": "ns/op",
            "extra": "13411354 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13411354 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13411354 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 89.96,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13433280 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 89.96,
            "unit": "ns/op",
            "extra": "13433280 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13433280 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13433280 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 135.2,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "8769676 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 135.2,
            "unit": "ns/op",
            "extra": "8769676 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "8769676 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8769676 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 135.3,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "7970413 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 135.3,
            "unit": "ns/op",
            "extra": "7970413 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "7970413 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7970413 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 94.98,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12618732 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 94.98,
            "unit": "ns/op",
            "extra": "12618732 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12618732 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12618732 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 95.61,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12642649 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 95.61,
            "unit": "ns/op",
            "extra": "12642649 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12642649 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12642649 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 148.8,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "8012346 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 148.8,
            "unit": "ns/op",
            "extra": "8012346 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "8012346 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "8012346 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 149.6,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "7973028 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 149.6,
            "unit": "ns/op",
            "extra": "7973028 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "7973028 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7973028 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 196.2,
            "unit": "ns/op\t      88 B/op\t       2 allocs/op",
            "extra": "6109878 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 196.2,
            "unit": "ns/op",
            "extra": "6109878 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 88,
            "unit": "B/op",
            "extra": "6109878 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6109878 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 195.7,
            "unit": "ns/op\t      88 B/op\t       2 allocs/op",
            "extra": "6059757 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 195.7,
            "unit": "ns/op",
            "extra": "6059757 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 88,
            "unit": "B/op",
            "extra": "6059757 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6059757 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 440.7,
            "unit": "ns/op\t     288 B/op\t       6 allocs/op",
            "extra": "2702908 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 440.7,
            "unit": "ns/op",
            "extra": "2702908 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 288,
            "unit": "B/op",
            "extra": "2702908 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "2702908 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 439.8,
            "unit": "ns/op\t     288 B/op\t       6 allocs/op",
            "extra": "2723917 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 439.8,
            "unit": "ns/op",
            "extra": "2723917 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 288,
            "unit": "B/op",
            "extra": "2723917 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "2723917 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 346.8,
            "unit": "ns/op\t     208 B/op\t       5 allocs/op",
            "extra": "3449740 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 346.8,
            "unit": "ns/op",
            "extra": "3449740 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "3449740 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3449740 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 343.6,
            "unit": "ns/op\t     208 B/op\t       5 allocs/op",
            "extra": "3500583 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 343.6,
            "unit": "ns/op",
            "extra": "3500583 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "3500583 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3500583 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 8955,
            "unit": "ns/op\t    4464 B/op\t      12 allocs/op",
            "extra": "132157 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 8955,
            "unit": "ns/op",
            "extra": "132157 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 4464,
            "unit": "B/op",
            "extra": "132157 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "132157 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 92.12,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13015437 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 92.12,
            "unit": "ns/op",
            "extra": "13015437 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13015437 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13015437 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 84.31,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14283470 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 84.31,
            "unit": "ns/op",
            "extra": "14283470 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14283470 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14283470 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 136.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8812136 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 136.8,
            "unit": "ns/op",
            "extra": "8812136 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8812136 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8812136 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 135.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8925081 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 135.5,
            "unit": "ns/op",
            "extra": "8925081 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8925081 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8925081 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 299.1,
            "unit": "ns/op\t      64 B/op\t       4 allocs/op",
            "extra": "4001656 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 299.1,
            "unit": "ns/op",
            "extra": "4001656 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "4001656 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "4001656 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 290,
            "unit": "ns/op\t      64 B/op\t       4 allocs/op",
            "extra": "4112940 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 290,
            "unit": "ns/op",
            "extra": "4112940 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "4112940 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "4112940 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 143299,
            "unit": "ns/op\t   65785 B/op\t     780 allocs/op",
            "extra": "7837 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 143299,
            "unit": "ns/op",
            "extra": "7837 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 65785,
            "unit": "B/op",
            "extra": "7837 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 780,
            "unit": "allocs/op",
            "extra": "7837 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 135752,
            "unit": "ns/op\t   65785 B/op\t     780 allocs/op",
            "extra": "8228 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 135752,
            "unit": "ns/op",
            "extra": "8228 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 65785,
            "unit": "B/op",
            "extra": "8228 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 780,
            "unit": "allocs/op",
            "extra": "8228 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 43.43,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "27551340 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 43.43,
            "unit": "ns/op",
            "extra": "27551340 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "27551340 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "27551340 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 38.83,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30591524 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 38.83,
            "unit": "ns/op",
            "extra": "30591524 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30591524 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30591524 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 92.11,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12933585 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 92.11,
            "unit": "ns/op",
            "extra": "12933585 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12933585 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12933585 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 83.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14385392 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 83.5,
            "unit": "ns/op",
            "extra": "14385392 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14385392 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14385392 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 60.06,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19890441 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 60.06,
            "unit": "ns/op",
            "extra": "19890441 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "19890441 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "19890441 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 60.16,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19853685 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 60.16,
            "unit": "ns/op",
            "extra": "19853685 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "19853685 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "19853685 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 120.7,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "10011410 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 120.7,
            "unit": "ns/op",
            "extra": "10011410 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "10011410 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10011410 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 118.9,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "9964245 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 118.9,
            "unit": "ns/op",
            "extra": "9964245 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "9964245 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9964245 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 60.66,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19928352 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 60.66,
            "unit": "ns/op",
            "extra": "19928352 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "19928352 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "19928352 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 60.13,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19922914 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 60.13,
            "unit": "ns/op",
            "extra": "19922914 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "19922914 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "19922914 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 72.95,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16405471 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 72.95,
            "unit": "ns/op",
            "extra": "16405471 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16405471 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16405471 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 72.66,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16486364 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 72.66,
            "unit": "ns/op",
            "extra": "16486364 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16486364 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16486364 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 72.84,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16352415 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 72.84,
            "unit": "ns/op",
            "extra": "16352415 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16352415 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16352415 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 72.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16420902 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 72.78,
            "unit": "ns/op",
            "extra": "16420902 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16420902 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16420902 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 194,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "6184566 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 194,
            "unit": "ns/op",
            "extra": "6184566 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "6184566 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "6184566 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 195.3,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "6137688 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 195.3,
            "unit": "ns/op",
            "extra": "6137688 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "6137688 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "6137688 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 792.1,
            "unit": "ns/op\t     432 B/op\t      15 allocs/op",
            "extra": "1492411 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 792.1,
            "unit": "ns/op",
            "extra": "1492411 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 432,
            "unit": "B/op",
            "extra": "1492411 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "1492411 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 787.5,
            "unit": "ns/op\t     432 B/op\t      15 allocs/op",
            "extra": "1531390 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 787.5,
            "unit": "ns/op",
            "extra": "1531390 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 432,
            "unit": "B/op",
            "extra": "1531390 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "1531390 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 114,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "10373551 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 114,
            "unit": "ns/op",
            "extra": "10373551 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "10373551 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10373551 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 113.9,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "10554352 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 113.9,
            "unit": "ns/op",
            "extra": "10554352 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "10554352 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10554352 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 228.3,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "5216178 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 228.3,
            "unit": "ns/op",
            "extra": "5216178 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "5216178 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5216178 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 228.6,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "5246725 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 228.6,
            "unit": "ns/op",
            "extra": "5246725 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "5246725 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5246725 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 85896,
            "unit": "ns/op\t   47000 B/op\t    1501 allocs/op",
            "extra": "14017 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 85896,
            "unit": "ns/op",
            "extra": "14017 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 47000,
            "unit": "B/op",
            "extra": "14017 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "14017 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 85418,
            "unit": "ns/op\t   47000 B/op\t    1501 allocs/op",
            "extra": "13996 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 85418,
            "unit": "ns/op",
            "extra": "13996 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 47000,
            "unit": "B/op",
            "extra": "13996 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "13996 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2708,
            "unit": "ns/op\t    1114 B/op\t      31 allocs/op",
            "extra": "418646 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2708,
            "unit": "ns/op",
            "extra": "418646 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1114,
            "unit": "B/op",
            "extra": "418646 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "418646 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2712,
            "unit": "ns/op\t    1115 B/op\t      31 allocs/op",
            "extra": "420524 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2712,
            "unit": "ns/op",
            "extra": "420524 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1115,
            "unit": "B/op",
            "extra": "420524 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "420524 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 83.75,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14416834 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 83.75,
            "unit": "ns/op",
            "extra": "14416834 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14416834 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14416834 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 84.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14374384 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 84.1,
            "unit": "ns/op",
            "extra": "14374384 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14374384 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14374384 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 131.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9143487 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 131.5,
            "unit": "ns/op",
            "extra": "9143487 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9143487 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9143487 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 131.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9158540 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 131.6,
            "unit": "ns/op",
            "extra": "9158540 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9158540 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9158540 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 88.05,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13365464 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 88.05,
            "unit": "ns/op",
            "extra": "13365464 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13365464 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13365464 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 88.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13696566 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 88.9,
            "unit": "ns/op",
            "extra": "13696566 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13696566 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13696566 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 143.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8107522 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 143.2,
            "unit": "ns/op",
            "extra": "8107522 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8107522 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8107522 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 143.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8259417 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 143.7,
            "unit": "ns/op",
            "extra": "8259417 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8259417 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8259417 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 95.73,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12401880 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 95.73,
            "unit": "ns/op",
            "extra": "12401880 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12401880 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12401880 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 96.28,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12619624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 96.28,
            "unit": "ns/op",
            "extra": "12619624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12619624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12619624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 169.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7532695 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 169.8,
            "unit": "ns/op",
            "extra": "7532695 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7532695 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7532695 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 156.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7661439 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 156.8,
            "unit": "ns/op",
            "extra": "7661439 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7661439 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7661439 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 100.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11707837 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 100.6,
            "unit": "ns/op",
            "extra": "11707837 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11707837 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11707837 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 101.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11975121 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 101.3,
            "unit": "ns/op",
            "extra": "11975121 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11975121 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11975121 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 168.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7015914 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 168.6,
            "unit": "ns/op",
            "extra": "7015914 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7015914 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7015914 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 168.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7146866 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 168.2,
            "unit": "ns/op",
            "extra": "7146866 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7146866 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7146866 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 105.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11333174 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 105.3,
            "unit": "ns/op",
            "extra": "11333174 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11333174 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11333174 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 105.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11350214 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 105.3,
            "unit": "ns/op",
            "extra": "11350214 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11350214 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11350214 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 181.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6561747 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 181.3,
            "unit": "ns/op",
            "extra": "6561747 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6561747 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6561747 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 180.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6638970 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 180.5,
            "unit": "ns/op",
            "extra": "6638970 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6638970 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6638970 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 116.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10389660 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 116.7,
            "unit": "ns/op",
            "extra": "10389660 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10389660 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10389660 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 115.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10305476 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 115.6,
            "unit": "ns/op",
            "extra": "10305476 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10305476 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10305476 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 203.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5931876 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 203.1,
            "unit": "ns/op",
            "extra": "5931876 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5931876 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5931876 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 201.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5818803 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 201.9,
            "unit": "ns/op",
            "extra": "5818803 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5818803 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5818803 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 67.46,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18007345 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 67.46,
            "unit": "ns/op",
            "extra": "18007345 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18007345 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18007345 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 66.05,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "18222381 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 66.05,
            "unit": "ns/op",
            "extra": "18222381 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "18222381 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "18222381 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 76.51,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15659737 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 76.51,
            "unit": "ns/op",
            "extra": "15659737 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15659737 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15659737 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 77.87,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15539929 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 77.87,
            "unit": "ns/op",
            "extra": "15539929 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15539929 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15539929 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 84.89,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14253006 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 84.89,
            "unit": "ns/op",
            "extra": "14253006 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14253006 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14253006 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 85.76,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14176858 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 85.76,
            "unit": "ns/op",
            "extra": "14176858 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14176858 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14176858 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 89.73,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13299771 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 89.73,
            "unit": "ns/op",
            "extra": "13299771 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13299771 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13299771 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 89.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13360456 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 89.78,
            "unit": "ns/op",
            "extra": "13360456 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13360456 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13360456 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 95.75,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12582211 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 95.75,
            "unit": "ns/op",
            "extra": "12582211 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12582211 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12582211 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 95.26,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12643404 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 95.26,
            "unit": "ns/op",
            "extra": "12643404 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12643404 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12643404 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 179.2,
            "unit": "ns/op\t      96 B/op\t       1 allocs/op",
            "extra": "6663208 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 179.2,
            "unit": "ns/op",
            "extra": "6663208 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "6663208 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "6663208 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 179.8,
            "unit": "ns/op\t      96 B/op\t       1 allocs/op",
            "extra": "6636816 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 179.8,
            "unit": "ns/op",
            "extra": "6636816 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "6636816 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "6636816 times\n4 procs"
          }
        ]
      },
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
          "id": "5f387ebcfc956840d42a146f7a55f49766e809ff",
          "message": "Fix/formatter (#887)\n\n* test: add integration test cases for formatter handling of literals (array, object, and general cases)\n\n* refactor: introduce column tracking in printer and enhance formatter logic for multiline and structured elements handling\n\n* Update test/integration/formatter/case_test.go\n\nCo-authored-by: Copilot <175728472+Copilot@users.noreply.github.com>\n\n* test: add test for nested object formatting respecting print width at line start\n\n* refactor: rewrite `writeRaw` for improved newline handling and byte width tracking, add unit tests\n\n---------\n\nCo-authored-by: Copilot <175728472+Copilot@users.noreply.github.com>",
          "timestamp": "2026-03-11T14:23:58-04:00",
          "tree_id": "b19fe79c9c7c604deabbce01cd1becc9baa3f372",
          "url": "https://github.com/MontFerret/ferret/commit/5f387ebcfc956840d42a146f7a55f49766e809ff"
        },
        "date": 1773253666191,
        "tool": "go",
        "benches": [
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 3629,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "324260 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 3629,
            "unit": "ns/op",
            "extra": "324260 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "324260 times\n4 procs"
          },
          {
            "name": "BenchmarkIteratorEOF (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "324260 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 364.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "3247560 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 364.5,
            "unit": "ns/op",
            "extra": "3247560 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "3247560 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "3247560 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime)",
            "value": 2369,
            "unit": "ns/op\t    1792 B/op\t       1 allocs/op",
            "extra": "508130 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - ns/op",
            "value": 2369,
            "unit": "ns/op",
            "extra": "508130 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - B/op",
            "value": 1792,
            "unit": "B/op",
            "extra": "508130 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectIterator (github.com/MontFerret/ferret/v2/pkg/runtime) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "508130 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/vm/internal/data)",
            "value": 986.2,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "1214046 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/vm/internal/data) - ns/op",
            "value": 986.2,
            "unit": "ns/op",
            "extra": "1214046 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/vm/internal/data) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "1214046 times\n4 procs"
          },
          {
            "name": "BenchmarkRangeIterator (github.com/MontFerret/ferret/v2/pkg/vm/internal/data) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "1214046 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 104260,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 104260,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 99468,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "12068 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 99468,
            "unit": "ns/op",
            "extra": "12068 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "12068 times\n4 procs"
          },
          {
            "name": "BenchmarkAddNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "12068 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 116263,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 116263,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 105030,
            "unit": "ns/op\t   47001 B/op\t    1501 allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 105030,
            "unit": "ns/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 47001,
            "unit": "B/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstNumericWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "10000 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 142394,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7957 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 142394,
            "unit": "ns/op",
            "extra": "7957 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7957 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7957 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 138441,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "8376 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 138441,
            "unit": "ns/op",
            "extra": "8376 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "8376 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstString_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "8376 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 150877,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7348 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 150877,
            "unit": "ns/op",
            "extra": "7348 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7348 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7348 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 144388,
            "unit": "ns/op\t   71002 B/op\t    3501 allocs/op",
            "extra": "7857 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 144388,
            "unit": "ns/op",
            "extra": "7857 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 71002,
            "unit": "B/op",
            "extra": "7857 times\n4 procs"
          },
          {
            "name": "BenchmarkAddConstStringWithParam_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "7857 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 173204,
            "unit": "ns/op\t   79002 B/op\t    3501 allocs/op",
            "extra": "6454 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 173204,
            "unit": "ns/op",
            "extra": "6454 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 79002,
            "unit": "B/op",
            "extra": "6454 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "6454 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 182222,
            "unit": "ns/op\t   79002 B/op\t    3501 allocs/op",
            "extra": "6466 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 182222,
            "unit": "ns/op",
            "extra": "6466 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 79002,
            "unit": "B/op",
            "extra": "6466 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralSimple_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "6466 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 201236,
            "unit": "ns/op\t   70322 B/op\t    3501 allocs/op",
            "extra": "5953 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 201236,
            "unit": "ns/op",
            "extra": "5953 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 70322,
            "unit": "B/op",
            "extra": "5953 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5953 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 201699,
            "unit": "ns/op\t   70322 B/op\t    3501 allocs/op",
            "extra": "5698 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 201699,
            "unit": "ns/op",
            "extra": "5698 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 70322,
            "unit": "B/op",
            "extra": "5698 times\n4 procs"
          },
          {
            "name": "BenchmarkTemplateLiteralNumeric_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3501,
            "unit": "allocs/op",
            "extra": "5698 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 4303,
            "unit": "ns/op\t    2040 B/op\t      55 allocs/op",
            "extra": "267664 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 4303,
            "unit": "ns/op",
            "extra": "267664 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 2040,
            "unit": "B/op",
            "extra": "267664 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 55,
            "unit": "allocs/op",
            "extra": "267664 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 4167,
            "unit": "ns/op\t    2040 B/op\t      55 allocs/op",
            "extra": "286878 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 4167,
            "unit": "ns/op",
            "extra": "286878 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 2040,
            "unit": "B/op",
            "extra": "286878 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregate_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 55,
            "unit": "allocs/op",
            "extra": "286878 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2472193,
            "unit": "ns/op\t  157196 B/op\t   19512 allocs/op",
            "extra": "511 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2472193,
            "unit": "ns/op",
            "extra": "511 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 157196,
            "unit": "B/op",
            "extra": "511 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 19512,
            "unit": "allocs/op",
            "extra": "511 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2313081,
            "unit": "ns/op\t  157193 B/op\t   19512 allocs/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2313081,
            "unit": "ns/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 157193,
            "unit": "B/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 19512,
            "unit": "allocs/op",
            "extra": "518 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 3917685,
            "unit": "ns/op\t 1782995 B/op\t   49527 allocs/op",
            "extra": "277 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 3917685,
            "unit": "ns/op",
            "extra": "277 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1782995,
            "unit": "B/op",
            "extra": "277 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 49527,
            "unit": "allocs/op",
            "extra": "277 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 3967657,
            "unit": "ns/op\t 1782991 B/op\t   49527 allocs/op",
            "extra": "304 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 3967657,
            "unit": "ns/op",
            "extra": "304 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1782991,
            "unit": "B/op",
            "extra": "304 times\n4 procs"
          },
          {
            "name": "BenchmarkGlobalCollectAggregateLargeInto_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 49527,
            "unit": "allocs/op",
            "extra": "304 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 10019754,
            "unit": "ns/op\t 4144437 B/op\t  127649 allocs/op",
            "extra": "120 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 10019754,
            "unit": "ns/op",
            "extra": "120 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 4144437,
            "unit": "B/op",
            "extra": "120 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 127649,
            "unit": "allocs/op",
            "extra": "120 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 9343662,
            "unit": "ns/op\t 4144405 B/op\t  127649 allocs/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 9343662,
            "unit": "ns/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 4144405,
            "unit": "B/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkGroupedCollectAggregateLarge_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 127649,
            "unit": "allocs/op",
            "extra": "128 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 95946,
            "unit": "ns/op\t   53497 B/op\t    1050 allocs/op",
            "extra": "12355 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 95946,
            "unit": "ns/op",
            "extra": "12355 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 53497,
            "unit": "B/op",
            "extra": "12355 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1050,
            "unit": "allocs/op",
            "extra": "12355 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 93580,
            "unit": "ns/op\t   53497 B/op\t    1050 allocs/op",
            "extra": "12734 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 93580,
            "unit": "ns/op",
            "extra": "12734 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 53497,
            "unit": "B/op",
            "extra": "12734 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_AllVars_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1050,
            "unit": "allocs/op",
            "extra": "12734 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2712464,
            "unit": "ns/op\t 1782654 B/op\t   49521 allocs/op",
            "extra": "436 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2712464,
            "unit": "ns/op",
            "extra": "436 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1782654,
            "unit": "B/op",
            "extra": "436 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 49521,
            "unit": "allocs/op",
            "extra": "436 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2701045,
            "unit": "ns/op\t 1782652 B/op\t   49521 allocs/op",
            "extra": "447 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2701045,
            "unit": "ns/op",
            "extra": "447 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1782652,
            "unit": "B/op",
            "extra": "447 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_SingleGroup_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 49521,
            "unit": "allocs/op",
            "extra": "447 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 77318,
            "unit": "ns/op\t   40697 B/op\t     850 allocs/op",
            "extra": "15486 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 77318,
            "unit": "ns/op",
            "extra": "15486 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 40697,
            "unit": "B/op",
            "extra": "15486 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 850,
            "unit": "allocs/op",
            "extra": "15486 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 77234,
            "unit": "ns/op\t   40696 B/op\t     850 allocs/op",
            "extra": "15726 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 77234,
            "unit": "ns/op",
            "extra": "15726 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 40696,
            "unit": "B/op",
            "extra": "15726 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Keep_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 850,
            "unit": "allocs/op",
            "extra": "15726 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 78272,
            "unit": "ns/op\t   40697 B/op\t     850 allocs/op",
            "extra": "15411 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 78272,
            "unit": "ns/op",
            "extra": "15411 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 40697,
            "unit": "B/op",
            "extra": "15411 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 850,
            "unit": "allocs/op",
            "extra": "15411 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 76505,
            "unit": "ns/op\t   40697 B/op\t     850 allocs/op",
            "extra": "15666 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 76505,
            "unit": "ns/op",
            "extra": "15666 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 40697,
            "unit": "B/op",
            "extra": "15666 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Custom_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 850,
            "unit": "allocs/op",
            "extra": "15666 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 8758,
            "unit": "ns/op\t     528 B/op\t      12 allocs/op",
            "extra": "136338 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 8758,
            "unit": "ns/op",
            "extra": "136338 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 528,
            "unit": "B/op",
            "extra": "136338 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "136338 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 8756,
            "unit": "ns/op\t     528 B/op\t      12 allocs/op",
            "extra": "134684 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 8756,
            "unit": "ns/op",
            "extra": "134684 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 528,
            "unit": "B/op",
            "extra": "134684 times\n4 procs"
          },
          {
            "name": "BenchmarkCollectProjection_Count_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "134684 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2962,
            "unit": "ns/op\t     640 B/op\t      17 allocs/op",
            "extra": "378489 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2962,
            "unit": "ns/op",
            "extra": "378489 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 640,
            "unit": "B/op",
            "extra": "378489 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 17,
            "unit": "allocs/op",
            "extra": "378489 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 1252,
            "unit": "ns/op\t     560 B/op\t       7 allocs/op",
            "extra": "954589 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 1252,
            "unit": "ns/op",
            "extra": "954589 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 560,
            "unit": "B/op",
            "extra": "954589 times\n4 procs"
          },
          {
            "name": "BenchmarkConstPropagation_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 7,
            "unit": "allocs/op",
            "extra": "954589 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2357,
            "unit": "ns/op\t    1112 B/op\t      34 allocs/op",
            "extra": "486460 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2357,
            "unit": "ns/op",
            "extra": "486460 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "486460 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "486460 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2311,
            "unit": "ns/op\t    1112 B/op\t      34 allocs/op",
            "extra": "496441 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2311,
            "unit": "ns/op",
            "extra": "496441 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "496441 times\n4 procs"
          },
          {
            "name": "BenchmarkForSort_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 34,
            "unit": "allocs/op",
            "extra": "496441 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 162.8,
            "unit": "ns/op\t      96 B/op\t       1 allocs/op",
            "extra": "7390270 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 162.8,
            "unit": "ns/op",
            "extra": "7390270 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "7390270 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7390270 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 162.9,
            "unit": "ns/op\t      96 B/op\t       1 allocs/op",
            "extra": "7305333 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 162.9,
            "unit": "ns/op",
            "extra": "7305333 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "7305333 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7305333 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 51.21,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "23768971 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 51.21,
            "unit": "ns/op",
            "extra": "23768971 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "23768971 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "23768971 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 52.27,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "23831108 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 52.27,
            "unit": "ns/op",
            "extra": "23831108 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "23831108 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "23831108 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 50.25,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "23991195 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 50.25,
            "unit": "ns/op",
            "extra": "23991195 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "23991195 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "23991195 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 50.49,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "23958111 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 50.49,
            "unit": "ns/op",
            "extra": "23958111 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "23958111 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall0Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "23958111 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 56.12,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21354576 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 56.12,
            "unit": "ns/op",
            "extra": "21354576 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21354576 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21354576 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 56.92,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21317586 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 56.92,
            "unit": "ns/op",
            "extra": "21317586 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21317586 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21317586 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 84.34,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "14038174 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 84.34,
            "unit": "ns/op",
            "extra": "14038174 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "14038174 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "14038174 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 84.21,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "14139204 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 84.21,
            "unit": "ns/op",
            "extra": "14139204 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "14139204 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall1Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "14139204 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 61.68,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19522198 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 61.68,
            "unit": "ns/op",
            "extra": "19522198 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "19522198 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "19522198 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 61.56,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19413520 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 61.56,
            "unit": "ns/op",
            "extra": "19413520 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "19413520 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "19413520 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 97.34,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "12156508 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 97.34,
            "unit": "ns/op",
            "extra": "12156508 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "12156508 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "12156508 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 97.72,
            "unit": "ns/op\t      32 B/op\t       1 allocs/op",
            "extra": "12187550 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 97.72,
            "unit": "ns/op",
            "extra": "12187550 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 32,
            "unit": "B/op",
            "extra": "12187550 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall2Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "12187550 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 70.06,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17680798 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 70.06,
            "unit": "ns/op",
            "extra": "17680798 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17680798 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17680798 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 67.78,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17653058 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 67.78,
            "unit": "ns/op",
            "extra": "17653058 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17653058 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17653058 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 113.5,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "10279296 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 113.5,
            "unit": "ns/op",
            "extra": "10279296 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "10279296 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10279296 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 114.2,
            "unit": "ns/op\t      48 B/op\t       1 allocs/op",
            "extra": "10371416 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 114.2,
            "unit": "ns/op",
            "extra": "10371416 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "10371416 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall3Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10371416 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 73.67,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16329934 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 73.67,
            "unit": "ns/op",
            "extra": "16329934 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16329934 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16329934 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 73.75,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16314582 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 73.75,
            "unit": "ns/op",
            "extra": "16314582 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16314582 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16314582 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 128.9,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "9388777 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 128.9,
            "unit": "ns/op",
            "extra": "9388777 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "9388777 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9388777 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 128.6,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "9384423 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 128.6,
            "unit": "ns/op",
            "extra": "9384423 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "9384423 times\n4 procs"
          },
          {
            "name": "BenchmarkFunctionCall4Fallback_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "9384423 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 192.3,
            "unit": "ns/op\t      88 B/op\t       2 allocs/op",
            "extra": "6276771 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 192.3,
            "unit": "ns/op",
            "extra": "6276771 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 88,
            "unit": "B/op",
            "extra": "6276771 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6276771 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 197.7,
            "unit": "ns/op\t      88 B/op\t       2 allocs/op",
            "extra": "6268132 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 197.7,
            "unit": "ns/op",
            "extra": "6268132 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 88,
            "unit": "B/op",
            "extra": "6268132 times\n4 procs"
          },
          {
            "name": "BenchmarkArrayLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 2,
            "unit": "allocs/op",
            "extra": "6268132 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 428.2,
            "unit": "ns/op\t     288 B/op\t       6 allocs/op",
            "extra": "2782654 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 428.2,
            "unit": "ns/op",
            "extra": "2782654 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 288,
            "unit": "B/op",
            "extra": "2782654 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "2782654 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 426.3,
            "unit": "ns/op\t     288 B/op\t       6 allocs/op",
            "extra": "2790900 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 426.3,
            "unit": "ns/op",
            "extra": "2790900 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 288,
            "unit": "B/op",
            "extra": "2790900 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 6,
            "unit": "allocs/op",
            "extra": "2790900 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 335.5,
            "unit": "ns/op\t     208 B/op\t       5 allocs/op",
            "extra": "3585825 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 335.5,
            "unit": "ns/op",
            "extra": "3585825 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "3585825 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3585825 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 334.1,
            "unit": "ns/op\t     208 B/op\t       5 allocs/op",
            "extra": "3566072 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 334.1,
            "unit": "ns/op",
            "extra": "3566072 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 208,
            "unit": "B/op",
            "extra": "3566072 times\n4 procs"
          },
          {
            "name": "BenchmarkObjectComputedLiterals_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 5,
            "unit": "allocs/op",
            "extra": "3566072 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 9140,
            "unit": "ns/op\t    4464 B/op\t      12 allocs/op",
            "extra": "127712 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 9140,
            "unit": "ns/op",
            "extra": "127712 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 4464,
            "unit": "B/op",
            "extra": "127712 times\n4 procs"
          },
          {
            "name": "BenchmarkLoop_Constants (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 12,
            "unit": "allocs/op",
            "extra": "127712 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 101.7,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11932207 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 101.7,
            "unit": "ns/op",
            "extra": "11932207 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11932207 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11932207 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 91.36,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13579639 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 91.36,
            "unit": "ns/op",
            "extra": "13579639 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13579639 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Scrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13579639 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 132.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9070084 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 132.2,
            "unit": "ns/op",
            "extra": "9070084 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9070084 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9070084 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 129.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9329118 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 129.2,
            "unit": "ns/op",
            "extra": "9329118 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9329118 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_Guard_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9329118 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 338.3,
            "unit": "ns/op\t      64 B/op\t       4 allocs/op",
            "extra": "3565386 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 338.3,
            "unit": "ns/op",
            "extra": "3565386 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "3565386 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "3565386 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 329.7,
            "unit": "ns/op\t      64 B/op\t       4 allocs/op",
            "extra": "3641410 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 329.7,
            "unit": "ns/op",
            "extra": "3641410 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "3641410 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ObjectPattern_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 4,
            "unit": "allocs/op",
            "extra": "3641410 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 163503,
            "unit": "ns/op\t   65785 B/op\t     780 allocs/op",
            "extra": "7740 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 163503,
            "unit": "ns/op",
            "extra": "7740 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 65785,
            "unit": "B/op",
            "extra": "7740 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 780,
            "unit": "allocs/op",
            "extra": "7740 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 149024,
            "unit": "ns/op\t   65785 B/op\t     780 allocs/op",
            "extra": "7532 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 149024,
            "unit": "ns/op",
            "extra": "7532 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 65785,
            "unit": "B/op",
            "extra": "7532 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_LoopMix_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 780,
            "unit": "allocs/op",
            "extra": "7532 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 44.14,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "26857288 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 44.14,
            "unit": "ns/op",
            "extra": "26857288 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "26857288 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "26857288 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 38.74,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "30391626 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 38.74,
            "unit": "ns/op",
            "extra": "30391626 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "30391626 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_ConstScrutinee_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "30391626 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 101.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12034665 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 101.1,
            "unit": "ns/op",
            "extra": "12034665 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12034665 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12034665 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 87.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13428236 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 87.5,
            "unit": "ns/op",
            "extra": "13428236 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13428236 times\n4 procs"
          },
          {
            "name": "BenchmarkMatch_MergePureLiteralResults_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13428236 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 60.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19777791 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 60.9,
            "unit": "ns/op",
            "extra": "19777791 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "19777791 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "19777791 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 60.38,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19952862 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 60.38,
            "unit": "ns/op",
            "extra": "19952862 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "19952862 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "19952862 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 116.5,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "10183423 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 116.5,
            "unit": "ns/op",
            "extra": "10183423 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "10183423 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10183423 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 116.2,
            "unit": "ns/op\t      64 B/op\t       1 allocs/op",
            "extra": "10375772 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 116.2,
            "unit": "ns/op",
            "extra": "10375772 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 64,
            "unit": "B/op",
            "extra": "10375772 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Short2_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10375772 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 60.61,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19744786 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 60.61,
            "unit": "ns/op",
            "extra": "19744786 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "19744786 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "19744786 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 61,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19880524 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 61,
            "unit": "ns/op",
            "extra": "19880524 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "19880524 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "19880524 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 80.52,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13989326 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 80.52,
            "unit": "ns/op",
            "extra": "13989326 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13989326 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13989326 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 80.03,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14962550 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 80.03,
            "unit": "ns/op",
            "extra": "14962550 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14962550 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14962550 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 84.46,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "14830362 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 84.46,
            "unit": "ns/op",
            "extra": "14830362 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "14830362 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "14830362 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 80.56,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "15215926 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 80.56,
            "unit": "ns/op",
            "extra": "15215926 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "15215926 times\n4 procs"
          },
          {
            "name": "BenchmarkOptionalUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "15215926 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 191.3,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "6264601 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 191.3,
            "unit": "ns/op",
            "extra": "6264601 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "6264601 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "6264601 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 189.7,
            "unit": "ns/op\t      96 B/op\t       3 allocs/op",
            "extra": "6277833 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 189.7,
            "unit": "ns/op",
            "extra": "6277833 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "6277833 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "6277833 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 768,
            "unit": "ns/op\t     432 B/op\t      15 allocs/op",
            "extra": "1567622 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 768,
            "unit": "ns/op",
            "extra": "1567622 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 432,
            "unit": "B/op",
            "extra": "1567622 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "1567622 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 762.9,
            "unit": "ns/op\t     432 B/op\t      15 allocs/op",
            "extra": "1555561 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 762.9,
            "unit": "ns/op",
            "extra": "1555561 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 432,
            "unit": "B/op",
            "extra": "1555561 times\n4 procs"
          },
          {
            "name": "BenchmarkMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 15,
            "unit": "allocs/op",
            "extra": "1555561 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 110.2,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "10605638 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 110.2,
            "unit": "ns/op",
            "extra": "10605638 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "10605638 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10605638 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 108.7,
            "unit": "ns/op\t      16 B/op\t       1 allocs/op",
            "extra": "10846504 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 108.7,
            "unit": "ns/op",
            "extra": "10846504 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 16,
            "unit": "B/op",
            "extra": "10846504 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "10846504 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 214.8,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "5525203 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 214.8,
            "unit": "ns/op",
            "extra": "5525203 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "5525203 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5525203 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 222.8,
            "unit": "ns/op\t      48 B/op\t       3 allocs/op",
            "extra": "5563302 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 222.8,
            "unit": "ns/op",
            "extra": "5563302 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 48,
            "unit": "B/op",
            "extra": "5563302 times\n4 procs"
          },
          {
            "name": "BenchmarkUnknownMemberAccess_Long_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 3,
            "unit": "allocs/op",
            "extra": "5563302 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 86848,
            "unit": "ns/op\t   47000 B/op\t    1501 allocs/op",
            "extra": "13708 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 86848,
            "unit": "ns/op",
            "extra": "13708 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 47000,
            "unit": "B/op",
            "extra": "13708 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "13708 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 86765,
            "unit": "ns/op\t   47000 B/op\t    1501 allocs/op",
            "extra": "13774 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 86765,
            "unit": "ns/op",
            "extra": "13774 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 47000,
            "unit": "B/op",
            "extra": "13774 times\n4 procs"
          },
          {
            "name": "BenchmarkParamLoop_Short_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1501,
            "unit": "allocs/op",
            "extra": "13774 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2566,
            "unit": "ns/op\t    1112 B/op\t      31 allocs/op",
            "extra": "431881 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2566,
            "unit": "ns/op",
            "extra": "431881 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "431881 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "431881 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 2579,
            "unit": "ns/op\t    1112 B/op\t      31 allocs/op",
            "extra": "429963 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 2579,
            "unit": "ns/op",
            "extra": "429963 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 1112,
            "unit": "B/op",
            "extra": "429963 times\n4 procs"
          },
          {
            "name": "BenchmarkRegexp_Loop_O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 31,
            "unit": "allocs/op",
            "extra": "429963 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 87.33,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13789087 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 87.33,
            "unit": "ns/op",
            "extra": "13789087 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13789087 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13789087 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 100.4,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "13761595 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 100.4,
            "unit": "ns/op",
            "extra": "13761595 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "13761595 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "13761595 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 136.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8562477 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 136.3,
            "unit": "ns/op",
            "extra": "8562477 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8562477 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8562477 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 135.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8843200 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 135.8,
            "unit": "ns/op",
            "extra": "8843200 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8843200 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8843200 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 91.71,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12851340 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 91.71,
            "unit": "ns/op",
            "extra": "12851340 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12851340 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12851340 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 91.79,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12980356 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 91.79,
            "unit": "ns/op",
            "extra": "12980356 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12980356 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12980356 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 152.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8017134 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 152.1,
            "unit": "ns/op",
            "extra": "8017134 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8017134 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8017134 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 149.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "8018516 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 149.5,
            "unit": "ns/op",
            "extra": "8018516 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "8018516 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "8018516 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 102.3,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12263293 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 102.3,
            "unit": "ns/op",
            "extra": "12263293 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12263293 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12263293 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 98.29,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "12261996 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 98.29,
            "unit": "ns/op",
            "extra": "12261996 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "12261996 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "12261996 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 161.6,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7411014 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 161.6,
            "unit": "ns/op",
            "extra": "7411014 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7411014 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7411014 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 170.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "7015372 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 170.9,
            "unit": "ns/op",
            "extra": "7015372 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "7015372 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "7015372 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 105.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11097468 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 105.1,
            "unit": "ns/op",
            "extra": "11097468 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11097468 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11097468 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 104,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "11563004 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 104,
            "unit": "ns/op",
            "extra": "11563004 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "11563004 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "11563004 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 175,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6869760 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 175,
            "unit": "ns/op",
            "extra": "6869760 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6869760 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6869760 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 174,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6860244 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 174,
            "unit": "ns/op",
            "extra": "6860244 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6860244 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6860244 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 111.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10776618 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 111.2,
            "unit": "ns/op",
            "extra": "10776618 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10776618 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10776618 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 110.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "10832786 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 110.1,
            "unit": "ns/op",
            "extra": "10832786 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "10832786 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "10832786 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 188.1,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6352108 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 188.1,
            "unit": "ns/op",
            "extra": "6352108 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6352108 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6352108 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 187.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "6363624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 187.9,
            "unit": "ns/op",
            "extra": "6363624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "6363624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "6363624 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 122.2,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9773935 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 122.2,
            "unit": "ns/op",
            "extra": "9773935 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9773935 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9773935 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 122.9,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "9830868 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 122.9,
            "unit": "ns/op",
            "extra": "9830868 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "9830868 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "9830868 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 211.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5653602 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 211.5,
            "unit": "ns/op",
            "extra": "5653602 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5653602 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5653602 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 211.8,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "5650834 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 211.8,
            "unit": "ns/op",
            "extra": "5650834 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "5650834 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls/UDF/Nested/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "5650834 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 50.84,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "22480276 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 50.84,
            "unit": "ns/op",
            "extra": "22480276 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "22480276 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "22480276 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 50.93,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "23420349 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 50.93,
            "unit": "ns/op",
            "extra": "23420349 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "23420349 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A0/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "23420349 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 56.47,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21335055 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 56.47,
            "unit": "ns/op",
            "extra": "21335055 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21335055 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21335055 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 62.54,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "21074583 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 62.54,
            "unit": "ns/op",
            "extra": "21074583 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "21074583 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A1/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "21074583 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 61.89,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19403533 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 61.89,
            "unit": "ns/op",
            "extra": "19403533 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "19403533 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "19403533 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 62.5,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "19392284 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 62.5,
            "unit": "ns/op",
            "extra": "19392284 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "19392284 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A2/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "19392284 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 68.59,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17531859 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 68.59,
            "unit": "ns/op",
            "extra": "17531859 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17531859 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17531859 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 68.22,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "17501130 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 68.22,
            "unit": "ns/op",
            "extra": "17501130 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "17501130 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A3/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "17501130 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 74.06,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16261922 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 74.06,
            "unit": "ns/op",
            "extra": "16261922 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16261922 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16261922 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 73.94,
            "unit": "ns/op\t       0 B/op\t       0 allocs/op",
            "extra": "16218085 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 73.94,
            "unit": "ns/op",
            "extra": "16218085 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 0,
            "unit": "B/op",
            "extra": "16218085 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A4/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 0,
            "unit": "allocs/op",
            "extra": "16218085 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 158.9,
            "unit": "ns/op\t      96 B/op\t       1 allocs/op",
            "extra": "7489053 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 158.9,
            "unit": "ns/op",
            "extra": "7489053 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "7489053 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O0 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7489053 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks)",
            "value": 159.7,
            "unit": "ns/op\t      96 B/op\t       1 allocs/op",
            "extra": "7554084 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - ns/op",
            "value": 159.7,
            "unit": "ns/op",
            "extra": "7554084 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - B/op",
            "value": 96,
            "unit": "B/op",
            "extra": "7554084 times\n4 procs"
          },
          {
            "name": "BenchmarkUdfCalls_HostBaseline/Host/TopLevel/A6/O1 (github.com/MontFerret/ferret/v2/test/integration/benchmarks) - allocs/op",
            "value": 1,
            "unit": "allocs/op",
            "extra": "7554084 times\n4 procs"
          }
        ]
      }
    ]
  }
}