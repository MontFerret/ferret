package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/vm"
	"github.com/MontFerret/ferret/test/integration/base"
)

func RunBenchmarkO0(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	base.RunBenchmarkWithOptimization(b, expression, compiler.O0, opts...)
}

func RunBenchmarkO1(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	base.RunBenchmarkWithOptimization(b, expression, compiler.O1, opts...)
}
