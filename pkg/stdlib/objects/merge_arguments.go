package objects

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func normalizeMergeArgs(ctx context.Context, args []runtime.Value, offset int) ([]runtime.Map, bool, error) {
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
				return nil, true, runtime.ArgError(err, offset)
			}

			return sources, true, nil
		}
	}

	sources := make([]runtime.Map, len(args))
	for index, value := range args {
		source, err := runtime.CastArg[runtime.Map](value, index+offset)
		if err != nil {
			return nil, false, err
		}

		sources[index] = source
	}

	return sources, false, nil
}

func mergeArgumentError(err error, index int, listForm bool, offset int) error {
	if listForm {
		return runtime.ArgError(fmt.Errorf("item %d: %w", index, err), offset)
	}

	return runtime.ArgError(err, index+offset)
}
