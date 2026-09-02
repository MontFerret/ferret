package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const coalesceBenchmarkQuery = "RETURN @value ?? @fallback"

const coalesceCompilerBenchmarkQuery = `
LET primary = @primary
RETURN primary ?? @secondary ?? @tertiary ?? "fallback"
`

func BenchmarkCoalescePresent_None(b *testing.B) {
	RunBenchmarkNone(
		b,
		coalesceBenchmarkQuery,
		WithParam("value", runtime.ZeroInt),
		WithParam("fallback", runtime.NewInt(42)),
	)
}

func BenchmarkCoalescePresent_Basic(b *testing.B) {
	RunBenchmarkBasic(
		b,
		coalesceBenchmarkQuery,
		WithParam("value", runtime.ZeroInt),
		WithParam("fallback", runtime.NewInt(42)),
	)
}

func BenchmarkCoalescePresent_Full(b *testing.B) {
	RunBenchmarkFull(
		b,
		coalesceBenchmarkQuery,
		WithParam("value", runtime.ZeroInt),
		WithParam("fallback", runtime.NewInt(42)),
	)
}

func BenchmarkCoalesceFallback_None(b *testing.B) {
	RunBenchmarkNone(
		b,
		coalesceBenchmarkQuery,
		WithParam("value", runtime.None),
		WithParam("fallback", runtime.NewInt(42)),
	)
}

func BenchmarkCoalesceFallback_Basic(b *testing.B) {
	RunBenchmarkBasic(
		b,
		coalesceBenchmarkQuery,
		WithParam("value", runtime.None),
		WithParam("fallback", runtime.NewInt(42)),
	)
}

func BenchmarkCoalesceFallback_Full(b *testing.B) {
	RunBenchmarkFull(
		b,
		coalesceBenchmarkQuery,
		WithParam("value", runtime.None),
		WithParam("fallback", runtime.NewInt(42)),
	)
}

func BenchmarkCompilerCompileCoalesce_None(b *testing.B) {
	benchmarkCompileQuery(b, coalesceCompilerBenchmarkQuery, compiler.OptimizationNone)
}

func BenchmarkCompilerCompileCoalesce_Basic(b *testing.B) {
	benchmarkCompileQuery(b, coalesceCompilerBenchmarkQuery, compiler.OptimizationBasic)
}

func BenchmarkCompilerCompileCoalesce_Full(b *testing.B) {
	benchmarkCompileQuery(b, coalesceCompilerBenchmarkQuery, compiler.OptimizationFull)
}
