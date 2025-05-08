package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// UNSHIFT prepends value to a given array.
// @param {Any[]} array - Target array.
// @param {Any} value - Target value to prepend.
// @param {Boolean} [unique=False] - Optional value indicating whether a value must be unique to be prepended. Default is false.
// @return {Any[]} - New array with prepended value.
func Unshift(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.None, err
	}

	list, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	value := args[1]
	uniq := runtime.False

	if len(args) > 2 {
		uniq, err = runtime.CastBoolean(args[2])

		if err != nil {
			return runtime.None, err
		}
	}

	size, err := list.Length(ctx)

	if err != nil {
		return runtime.None, err
	}

	result := runtime.NewArray64(size + 1)

	if !uniq {
		_ = result.Add(ctx, value)

		err = list.ForEach(ctx, func(ctx context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
			_ = result.Add(ctx, value)

			return runtime.True, nil
		})

		if err != nil {
			return runtime.None, err
		}

		return result, nil
	}

	_ = result.Add(ctx, value)

	err = list.ForEach(ctx, func(ctx context.Context, el runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		if runtime.CompareValues(el, value) != 0 {
			_ = result.Add(ctx, el)
		}

		return true, nil
	})

	if err != nil {
		return runtime.None, err
	}

	return result, nil
}
