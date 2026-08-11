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

func BenchmarkCoalescePresent_O0(b *testing.B) {
	RunBenchmarkO0(
		b,
		coalesceBenchmarkQuery,
		WithParam("value", runtime.ZeroInt),
		WithParam("fallback", runtime.NewInt(42)),
	)
}

func BenchmarkCoalescePresent_O1(b *testing.B) {
	RunBenchmarkO1(
		b,
		coalesceBenchmarkQuery,
		WithParam("value", runtime.ZeroInt),
		WithParam("fallback", runtime.NewInt(42)),
	)
}

func BenchmarkCoalesceFallback_O0(b *testing.B) {
	RunBenchmarkO0(
		b,
		coalesceBenchmarkQuery,
		WithParam("value", runtime.None),
		WithParam("fallback", runtime.NewInt(42)),
	)
}

func BenchmarkCoalesceFallback_O1(b *testing.B) {
	RunBenchmarkO1(
		b,
		coalesceBenchmarkQuery,
		WithParam("value", runtime.None),
		WithParam("fallback", runtime.NewInt(42)),
	)
}

func BenchmarkCompilerCompileCoalesce_O0(b *testing.B) {
	benchmarkCompileQuery(b, coalesceCompilerBenchmarkQuery, compiler.O0)
}

func BenchmarkCompilerCompileCoalesce_O1(b *testing.B) {
	benchmarkCompileQuery(b, coalesceCompilerBenchmarkQuery, compiler.O1)
}
