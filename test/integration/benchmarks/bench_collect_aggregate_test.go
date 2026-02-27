package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/integration/base"
)

const globalCollectAggregateQuery = `
LET users = [
	{ age: 31, salary: 75000 },
	{ age: 25, salary: 60000 },
	{ age: 36, salary: 80000 },
	{ age: 69, salary: 95000 },
	{ age: 45, salary: 70000 }
]

FOR u IN users
	COLLECT AGGREGATE
		minAge = MIN(u.age),
		maxAge = MAX(u.age),
		totalSalary = SUM(u.salary)
	RETURN { minAge, maxAge, totalSalary }
`

const globalCollectAggregateLargeQuery = `
LET values = 1..10000

FOR v IN values
	COLLECT AGGREGATE
		cnt = COUNT(v),
		sum = SUM(v),
		min = MIN(v),
		max = MAX(v),
		avg = AVERAGE(v)
	RETURN { cnt, sum, min, max, avg }
`

const globalCollectAggregateLargeIntoQuery = `
LET values = 1..10000

FOR v IN values
	COLLECT AGGREGATE
		sum = SUM(v),
		min = MIN(v),
		max = MAX(v)
	INTO groups
	RETURN { sum, min, max, groups }
`

const groupedCollectAggregateLargeQuery = `
LET values = 1..10000

FOR v IN values
	COLLECT g = v % 100
	AGGREGATE
		cnt = COUNT(v),
		sum = SUM(v),
		min = MIN(v),
		max = MAX(v),
		avg = AVERAGE(v)
	RETURN { g, cnt, sum, min, max, avg }
`

func BenchmarkGlobalCollectAggregate_O0(b *testing.B) {
	RunBenchmarkO0(b, globalCollectAggregateQuery, vm.WithNamespace(base.Stdlib()))
}

func BenchmarkGlobalCollectAggregate_O1(b *testing.B) {
	RunBenchmarkO1(b, globalCollectAggregateQuery, vm.WithNamespace(base.Stdlib()))
}

func BenchmarkGlobalCollectAggregateLarge_O0(b *testing.B) {
	RunBenchmarkO0(b, globalCollectAggregateLargeQuery, vm.WithNamespace(base.Stdlib()))
}

func BenchmarkGlobalCollectAggregateLarge_O1(b *testing.B) {
	RunBenchmarkO1(b, globalCollectAggregateLargeQuery, vm.WithNamespace(base.Stdlib()))
}

func BenchmarkGlobalCollectAggregateLargeInto_O0(b *testing.B) {
	RunBenchmarkO0(b, globalCollectAggregateLargeIntoQuery, vm.WithNamespace(base.Stdlib()))
}

func BenchmarkGlobalCollectAggregateLargeInto_O1(b *testing.B) {
	RunBenchmarkO1(b, globalCollectAggregateLargeIntoQuery, vm.WithNamespace(base.Stdlib()))
}

func BenchmarkGroupedCollectAggregateLarge_O0(b *testing.B) {
	RunBenchmarkO0(b, groupedCollectAggregateLargeQuery, vm.WithNamespace(base.Stdlib()))
}

func BenchmarkGroupedCollectAggregateLarge_O1(b *testing.B) {
	RunBenchmarkO1(b, groupedCollectAggregateLargeQuery, vm.WithNamespace(base.Stdlib()))
}
