package benchmarks_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func BenchmarkCompilerTailCallElimination(b *testing.B) {
	queries := []struct {
		name  string
		query string
	}{
		{"NoUDF", `return @value + 1`},
		{"OrdinaryUDF", `func inc(value) => value + 1
return inc(1) + inc(2)`},
		{"TailChain", `func first(value) => second(value)
func second(value) => third(value)
func third(value) => value * 3
return first(2)`},
		{"Recursive", `func recur(value) => recur(NEXT(value))
return recur(256) on error return none`},
		{"Capture", compilerUdfLifecycleQuery},
		{"Recovery", `func target(value) => value / 0
func forward(value) => target(value) on error retry 2 or return value
return forward(3)`},
	}
	modes := []struct {
		name    string
		options []compiler.Option
	}{
		{"Default", nil},
		{"None", []compiler.Option{compiler.WithOptimizationLevel(compiler.None)}},
		{"Basic", []compiler.Option{compiler.WithOptimizationLevel(compiler.Basic)}},
		{"Full", []compiler.Option{compiler.WithOptimizationLevel(compiler.Full)}},
		{"Debug", []compiler.Option{compiler.WithOptimizationLevel(compiler.Full), compiler.WithDebugInfo()}},
	}

	for _, query := range queries {
		for _, mode := range modes {
			b.Run(query.name+"/"+mode.name, func(b *testing.B) {
				instance := mustNewCompiler(b, mode.options...)
				src := source.New("tail-call-benchmark.fql", query.query)

				b.ReportAllocs()
				b.ResetTimer()

				for range b.N {
					if _, err := instance.Compile(b.Context(), src); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
