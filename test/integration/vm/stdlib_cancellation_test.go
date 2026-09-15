package vm_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/collections"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestStdlibCompletesBeforeVMCancellationBoundary(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.None, compiler.Basic, compiler.Full} {
		for name, query := range map[string]string{
			"return":    `RETURN COUNT(@source) ON ERROR RETURN AFTER()`,
			"next call": "LET n = COUNT(@source)\nRETURN AFTER(n)",
		} {
			t.Run(optimizationName(level)+"/"+name, func(t *testing.T) {
				instance := compileCancellationVM(t, level, query)
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				source := &stdlibCancellationSource{Value: runtime.None, expected: ctx, cancel: cancel}
				var callResult runtime.Value
				var callErr error
				afterCalls := 0
				env := cancellationEnvironment(t,
					vm.WithParam("source", source),
					vm.WithFunction("COUNT", func(c context.Context, args ...runtime.Value) (runtime.Value, error) {
						callResult, callErr = collections.Count(c, args[0])

						return callResult, callErr
					}),
					vm.WithFunction("AFTER", func(context.Context, ...runtime.Value) (runtime.Value, error) {
						afterCalls++

						return runtime.True, nil
					}),
				)

				result, err := instance.Run(ctx, env)
				if result != nil {
					_ = result.Close()
				}

				if callErr != nil || callResult != runtime.Int(3) || source.visits != 3 || source.closes != 1 {
					t.Fatalf("stdlib did not complete: %v, %v, visits %d, closes %d", callResult, callErr, source.visits, source.closes)
				}

				if !errors.Is(err, context.Canceled) || afterCalls != 0 {
					t.Fatalf("VM error = %v, subsequent calls = %d", err, afterCalls)
				}
			})
		}
	}
}
