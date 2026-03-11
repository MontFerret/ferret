window.BENCHMARK_DATA = {
  "lastUpdate": 1773244722099,
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
      }
    ]
  }
}