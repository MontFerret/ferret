package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/integration/base"
)

var (
	WithParam = base.WithParam
)

func RunBenchmarkO0(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	base.RunBenchmarkWithOptimization(b, expression, compiler.O0, opts...)
}

func RunBenchmarkO1(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	base.RunBenchmarkWithOptimization(b, expression, compiler.O1, opts...)
}
