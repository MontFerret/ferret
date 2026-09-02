package benchmarks_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

const dispatchLongStatementQuery = `
DISPATCH "click" IN @d
RETURN 1`

const dispatchShorthandStatementQuery = `
@d <- "click"
RETURN 1`

const dispatchGroupedTargetQuery = `
DISPATCH "click" IN (QUERY ONE "#submit" IN @page USING css)
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

func BenchmarkDispatchLongStatement_None(b *testing.B) {
	RunBenchmarkNone(b, dispatchLongStatementQuery, vm.WithParam("d", &benchmarkDispatcher{}))
}

func BenchmarkDispatchLongStatement_Full(b *testing.B) {
	RunBenchmarkFull(b, dispatchLongStatementQuery, vm.WithParam("d", &benchmarkDispatcher{}))
}

func BenchmarkDispatchShorthandStatement_None(b *testing.B) {
	RunBenchmarkNone(b, dispatchShorthandStatementQuery, vm.WithParam("d", &benchmarkDispatcher{}))
}

func BenchmarkDispatchShorthandStatement_Full(b *testing.B) {
	RunBenchmarkFull(b, dispatchShorthandStatementQuery, vm.WithParam("d", &benchmarkDispatcher{}))
}

func BenchmarkCompilerCompileDispatchGroupedTarget_None(b *testing.B) {
	benchmarkCompileQuery(b, dispatchGroupedTargetQuery, compiler.OptimizationNone)
}

func BenchmarkCompilerCompileDispatchGroupedTarget_Full(b *testing.B) {
	benchmarkCompileQuery(b, dispatchGroupedTargetQuery, compiler.OptimizationFull)
}
