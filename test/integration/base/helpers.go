package base

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func ForWhileHelpers() runtime.FunctionDefs {
	builder := runtime.NewFunctionsBuilder()

	builder.Var().
		Add("UNTIL", StateFn[int](func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			state := GetFnState[int](ctx)
			untilCounter := state.Get()

			if untilCounter < int(runtime.ToIntSafe(ctx, args[0])) {
				untilCounter++
				state.Set(untilCounter)

				return runtime.True, nil
			}

			return runtime.False, nil
		}, func(ctx context.Context) int {
			return 0
		})).
		Add("W_POP", func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			list, err := runtime.CastList(args[0])

			if err != nil {
				return runtime.False, err
			}

			size, _ := list.Length(ctx)

			if size > 0 {
				_, err := list.RemoveAt(ctx, 0)

				if err != nil {
					return runtime.False, err
				}

				return runtime.True, nil
			}

			return runtime.False, nil
		}).
		Add("COUNTER", StateFn[int](func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			state := GetFnState[int](ctx)
			counter := state.Get()
			counter++

			state.Set(counter)

			return runtime.Int(counter), nil
		}, func(ctx context.Context) int {
			return -1
		}))

	return builder
}

func WithParam(name string, value any) vm.EnvironmentOption {
	parsed, err := runtime.ValueOf(value)

	if err != nil {
		panic(err)
	}

	return vm.WithParam(name, parsed)
}
