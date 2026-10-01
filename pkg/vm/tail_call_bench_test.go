package vm

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func BenchmarkTailCallEliminationExecution(b *testing.B) {
	queries := []struct {
		name  string
		query string
	}{
		{"Chain", `func first(value) => second(value)
func second(value) => third(value)
func third(value) => value * 3
return for item in 1..128 return first(item)`},
		{"Recursive", `func recur(value) => recur(NEXT(value))
return recur(256) on error return none`},
	}
	stop := errors.New("recursion complete")
	env, err := NewEnvironment([]EnvironmentOption{WithFunction("NEXT", func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
		value := args[0].(runtime.Int)
		if value == 0 {
			return runtime.None, stop
		}

		return value - 1, nil
	})})
	if err != nil {
		b.Fatal(err)
	}

	for _, query := range queries {
		for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
			b.Run(query.name+"/"+level.String(), func(b *testing.B) {
				program, err := mustNewCompiler(b, compiler.WithOptimizationLevel(level)).Compile(b.Context(), source.NewAnonymous(query.query))
				if err != nil {
					b.Fatal(err)
				}

				benchmarkCancellationSafepoint(b, program, env, context.Background())
			})
		}
	}
}
