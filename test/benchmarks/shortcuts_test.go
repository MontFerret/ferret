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

func RunBenchmarkNone(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	spec.RunBenchmarkWithOptimization(b, expression, compiler.None, opts...)
}

func RunBenchmarkBasic(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	spec.RunBenchmarkWithOptimization(b, expression, compiler.Basic, opts...)
}

func RunBenchmarkFull(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	spec.RunBenchmarkWithOptimization(b, expression, compiler.Full, opts...)
}
