package objects

import (
	"context"
	"fmt"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func filterKeys(ctx context.Context, args []runtime.Value, apply func(context.Context, runtime.Map, ...runtime.String) error) (runtime.Value, error) {
	src, keys, err := normalizeKeyArgs(ctx, args)
	if err != nil {
		return runtime.None, err
	}

	copied, err := runtime.CloneOrCopy(ctx, src)
	if err != nil {
		return runtime.None, runtime.ArgError(err, 0)
	}

	dst, ok := copied.(runtime.Map)
	if !ok {
		return runtime.None, runtime.ArgError(runtime.TypeErrorOf(copied, runtime.TypeMap), 0)
	}

	if err := apply(ctx, dst, keys...); err != nil {
		return runtime.None, runtime.ArgError(err, 0)
	}

	return dst, nil
}

func normalizeKeyArgs(ctx context.Context, args []runtime.Value) (runtime.Map, []runtime.String, error) {
	if err := runtime.ValidateArgs(args, 2, runtime.MaxArgs); err != nil {
		return nil, nil, err
	}

	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}

	src, err := runtime.CastArg[runtime.Map](args[0], 0)
	if err != nil {
		return nil, nil, err
	}

	if len(args) == 2 {
		if list, ok := args[1].(runtime.List); ok {
			var keys []runtime.String
			err := list.ForEach(ctx, func(ctx context.Context, value runtime.Value, index runtime.Int) (runtime.Boolean, error) {
				if err := ctx.Err(); err != nil {
					return false, err
				}

				key, err := runtime.CastString(value)
				if err != nil {
					return false, fmt.Errorf("item %d: %w", index, err)
				}

				keys = append(keys, key)

				return true, nil
			})
			if err != nil {
				return nil, nil, runtime.ArgError(err, 1)
			}

			return src, keys, nil
		}
	}

	keys := make([]runtime.String, len(args)-1)
	for index, value := range args[1:] {
		key, err := runtime.CastArg[runtime.String](value, index+1)
		if err != nil {
			return nil, nil, err
		}

		keys[index] = key
	}

	return src, keys, nil
}
