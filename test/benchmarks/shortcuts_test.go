package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/vm"
	"github.com/MontFerret/ferret/v2/test/spec"
)

var (
	WithParam = spec.WithParam
)

func RunBenchmarkO0(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	spec.RunBenchmarkWithOptimization(b, expression, compiler.O0, opts...)
}

func RunBenchmarkO1(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	spec.RunBenchmarkWithOptimization(b, expression, compiler.O1, opts...)
}
