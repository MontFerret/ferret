package benchmarks_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type benchmarkQueryable struct {
	result runtime.List
}

func newBenchmarkQueryable() *benchmarkQueryable {
	return &benchmarkQueryable{
		result: runtime.NewArrayWith(runtime.NewString("a")),
	}
}

func (q *benchmarkQueryable) Query(_ context.Context, _ runtime.Query) (runtime.List, error) {
	return q.result, nil
}

func (q *benchmarkQueryable) QueryOne(_ context.Context, _ runtime.Query) (runtime.Value, error) {
	return runtime.NewString("a"), nil
}

func (q *benchmarkQueryable) QueryCount(_ context.Context, _ runtime.Query) (runtime.Int, error) {
	return runtime.NewInt(1), nil
}

func (q *benchmarkQueryable) QueryExists(_ context.Context, _ runtime.Query) (runtime.Boolean, error) {
	return runtime.True, nil
}

func (q *benchmarkQueryable) MarshalJSON() ([]byte, error) {
	return []byte(`"queryable"`), nil
}

func (q *benchmarkQueryable) String() string {
	return "queryable"
}

func (q *benchmarkQueryable) Unwrap() any {
	return "queryable"
}

func (q *benchmarkQueryable) Hash() uint64 {
	return 0
}

func (q *benchmarkQueryable) Copy() runtime.Value {
	return q
}
