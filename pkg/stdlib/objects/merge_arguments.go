package objects

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func normalizeMergeArgs(ctx context.Context, args []runtime.Value) ([]runtime.Map, bool, error) {
	if err := runtime.ValidateArgs(args, 1, runtime.MaxArgs); err != nil {
		return nil, false, err
	}

	if err := ctx.Err(); err != nil {
		return nil, false, err
	}

	if len(args) == 1 {
		if list, ok := args[0].(runtime.List); ok {
			var sources []runtime.Map
			err := list.ForEach(ctx, func(ctx context.Context, value runtime.Value, index runtime.Int) (runtime.Boolean, error) {
				if err := ctx.Err(); err != nil {
					return false, err
				}

				source, ok := value.(runtime.Map)
				if !ok {
					return false, fmt.Errorf("item %d: %w", index, runtime.TypeErrorOf(value, runtime.TypeMap))
				}

				sources = append(sources, source)

				return true, nil
			})
			if err != nil {
				return nil, true, runtime.ArgError(err, 0)
			}

			return sources, true, nil
		}
	}

	sources, err := runtime.CastArgs[runtime.Map](args)

	return sources, false, err
}

func mergeArgumentError(err error, index int, listForm bool) error {
	if listForm {
		return runtime.ArgError(fmt.Errorf("item %d: %w", index, err), 0)
	}

	return runtime.ArgError(err, index)
}
