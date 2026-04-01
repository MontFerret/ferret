package benchmarks_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

const dispatchLongStatementQuery = `
DISPATCH "click" IN @d
RETURN 1`

const dispatchShorthandStatementQuery = `
"click" -> @d
RETURN 1`

type benchmarkDispatcher struct{}

func (d *benchmarkDispatcher) Dispatch(_ context.Context, event runtime.DispatchEvent) error {
	return nil
}

func (d *benchmarkDispatcher) MarshalJSON() ([]byte, error) {
	return []byte(`"dispatcher"`), nil
}

func (d *benchmarkDispatcher) String() string {
	return "dispatcher"
}

func (d *benchmarkDispatcher) Unwrap() any {
	return "dispatcher"
}

func (d *benchmarkDispatcher) Hash() uint64 {
	return 0
}

func (d *benchmarkDispatcher) Copy() runtime.Value {
	return d
}

func BenchmarkDispatchLongStatement_O0(b *testing.B) {
	RunBenchmarkO0(b, dispatchLongStatementQuery, vm.WithParam("d", &benchmarkDispatcher{}))
}

func BenchmarkDispatchLongStatement_O1(b *testing.B) {
	RunBenchmarkO1(b, dispatchLongStatementQuery, vm.WithParam("d", &benchmarkDispatcher{}))
}

func BenchmarkDispatchShorthandStatement_O0(b *testing.B) {
	RunBenchmarkO0(b, dispatchShorthandStatementQuery, vm.WithParam("d", &benchmarkDispatcher{}))
}

func BenchmarkDispatchShorthandStatement_O1(b *testing.B) {
	RunBenchmarkO1(b, dispatchShorthandStatementQuery, vm.WithParam("d", &benchmarkDispatcher{}))
}
