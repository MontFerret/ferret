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

// Benchmark shortcuts and gate usage.
//
// This package exposes helpers for running the same benchmark expression
// under different compiler optimization levels (O0, O1, etc.).
//
// Benchmark parameters can be passed via WithParam and wired into the
// environment options used by the benchmarks. A typical pattern looks like:
//
//   env, _ := vm.New(WithParam("dataset", "small"))
//   // pass env as one of the vm.EnvironmentOption values when calling
//   // RunBenchmarkO0 / RunBenchmarkO1 helpers from individual benchmarks.
//
// From the command line, benchmarks are usually run with:
//
//   go test ./test/integration/benchmarks -bench=. -run=^$ -args dataset=small
//
// Individual benchmark files can interpret these parameters via WithParam
// to select specific datasets, query variants, or other benchmark gates.
func RunBenchmarkO0(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	base.RunBenchmarkWithOptimization(b, expression, compiler.O0, opts...)
}

func RunBenchmarkO1(b *testing.B, expression string, opts ...vm.EnvironmentOption) {
	base.RunBenchmarkWithOptimization(b, expression, compiler.O1, opts...)
}
