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

func BenchmarkGlobalCollectAggregate_O0(b *testing.B) {
	RunBenchmarkO0(b, globalCollectAggregateQuery, vm.WithFunctions(base.Stdlib()))
}

func BenchmarkGlobalCollectAggregate_O1(b *testing.B) {
	RunBenchmarkO1(b, globalCollectAggregateQuery, vm.WithFunctions(base.Stdlib()))
}
