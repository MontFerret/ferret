package vm

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
	vmtest "github.com/MontFerret/ferret/v2/pkg/vm/test"
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

	for _, path := range []struct {
		name string
		grow bool
	}{{"Reuse", false}, {"Growth", true}} {
		for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
			b.Run("OwnedAliases/"+path.name+"/"+level.String(), func(b *testing.B) {
				benchmarkTailCallOwnedAliases(b, level, path.grow)
			})
		}
	}
}

func benchmarkTailCallOwnedAliases(b *testing.B, level compiler.OptimizationLevel, grow bool) {
	query := tailCallLifecycleQuery(grow, "return forward(@borrowed)")
	program, err := mustNewCompiler(b, compiler.WithOptimizationLevel(level)).Compile(b.Context(), source.NewAnonymous(query))
	if err != nil {
		b.Fatal(err)
	}

	instance, err := NewWith(program, WithTesting(vmtest.WithBenchmarkMode()))
	if err != nil {
		b.Fatal(err)
	}
	defer func() { _ = instance.Close() }()

	candidate := assertTailCallCandidate(b, instance, level, "forward", "leaf", true)
	leaf := program.Functions.UserDefined[candidate.ID]
	borrowed := newTrackingCloser("borrowed")
	var resources [2]*trackingCloser
	var oldWindow *runtime.Value
	openCalls, depth := 0, 0
	validate := true
	noop := func(context.Context, ...runtime.Value) (runtime.Value, error) { return runtime.None, nil }
	env, err := NewEnvironment([]EnvironmentOption{
		WithParam("borrowed", borrowed),
		WithFunction("OPEN", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			resource := newTrackingCloser("owned")
			resources[openCalls%len(resources)] = resource
			openCalls++

			return resource, nil
		}),
		WithFunction("VALUE", func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
			return args[0], nil
		}),
		WithFunction("PAD", noop),
		WithFunction("SNAPSHOT", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			if validate {
				state := &instance.state
				if (cap(state.registers) < leaf.Registers) != grow {
					b.Fatalf("window capacity = %d, leaf registers = %d, growth = %v", cap(state.registers), leaf.Registers, grow)
				}

				oldWindow = &state.registers[0]
				depth = state.frames.Len()
			}

			return runtime.None, nil
		}),
		WithFunction("ENTRY", func(context.Context, ...runtime.Value) (runtime.Value, error) {
			state := &instance.state
			if validate {
				wantDepth := depth
				if level == compiler.None {
					wantDepth++
				}

				if state.frames.Len() != wantDepth || state.frames.CurrentFunctionID() != candidate.ID {
					b.Fatal("owned-alias benchmark did not enter the intended callee")
				}

				if level != compiler.None {
					key, _, _ := mem.ResourceKeyOf(resources[0])
					if !state.owned.Owns(resources[0]) || state.aliases.Count(key) != 2 || (&state.registers[0] != oldWindow) != grow {
						b.Fatal("owned-alias benchmark did not transfer both aliases along the intended window path")
					}
				}
			}

			state.clearRegister(bytecode.NewRegister(1))

			return runtime.None, nil
		}),
		WithFunction("CHECK_SURVIVOR", func(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
			if validate && (args[0] != resources[0] || args[1] != borrowed || resources[0].closed != 0) {
				b.Fatal("owned-alias benchmark lost its surviving resource")
			}

			return runtime.None, nil
		}),
		WithFunction("STEP", noop),
		WithFunction("AFTER_STEP", noop),
	})
	if err != nil {
		b.Fatal(err)
	}

	result, err := instance.Run(b.Context(), env)
	if err != nil {
		b.Fatal(err)
	}
	if result.Root() != resources[0] {
		b.Fatal("owned-alias benchmark returned the wrong resource")
	}

	if err := result.Close(); err != nil {
		b.Fatal(err)
	}

	validate = false
	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		result, err := instance.Run(b.Context(), env)
		if err != nil {
			b.Fatal(err)
		}

		if err := result.Close(); err != nil {
			b.Fatal(err)
		}
	}

	b.StopTimer()
	if resources[0].closed != 1 || resources[1].closed != 1 || borrowed.closed != 0 {
		b.Fatal("owned-alias benchmark resource cleanup changed")
	}
}
