package vm_test

import (
	"context"

	"github.com/MontFerret/ferret/test/integration/base"

	"github.com/MontFerret/ferret/pkg/runtime"
)

func ForWhileHelpers() runtime.Functions {
	return runtime.NewFunctionsFromMap(map[string]runtime.Function{
		"UNTIL": base.StateFn[int](func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			state := base.GetFnState[int](ctx)
			untilCounter := state.Get()

			if untilCounter < int(runtime.ToIntSafe(ctx, args[0])) {
				untilCounter++
				state.Set(untilCounter)

				return runtime.True, nil
			}

			return runtime.False, nil
		}, func(ctx context.Context) int {
			return 0
		}),
		"W_POP": func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
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
		},
		"COUNTER": base.StateFn[int](func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			state := base.GetFnState[int](ctx)
			counter := state.Get()
			counter++

			state.Set(counter)

			return runtime.Int(counter), nil
		}, func(ctx context.Context) int {
			return -1
		}),
		"COUNTER2": base.StateFn[int](func(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
			state := base.GetFnState[int](ctx)
			counter := state.Get()
			counter++

			state.Set(counter)

			return runtime.Int(counter), nil
		}, func(ctx context.Context) int {
			return -1
		}),
	})
}
