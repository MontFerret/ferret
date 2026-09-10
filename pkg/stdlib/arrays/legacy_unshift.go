package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// unshift prepends value to a given array.
// @param array {Any[]} Target array.
// @param value {Any} Target value to prepend.
// @return {Any[]} New array with prepended value.
// @deprecated Use arrays::concat([value], array) and arrays::remove as needed.
func legacyUnshift2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return legacyUnshift3(ctx, arg1, arg2, runtime.False)
}

// unshift prepends value to a given array.
// @param array {Any[]} Target array.
// @param value {Any} Target value to prepend.
// @param unique {Boolean} Optional value indicating whether a value must be unique to be prepended. Default is false.
// @return {Any[]} New array with prepended value.
// @deprecated Use arrays::concat([value], array) and arrays::remove as needed.
func legacyUnshift3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](arg1, 0)
	if err != nil {
		return runtime.None, err
	}

	uniq, err := runtime.CastArg[runtime.Boolean](arg3, 2)
	if err != nil {
		return runtime.None, err
	}

	size, err := list.Length(ctx)
	if err != nil {
		return runtime.None, err
	}

	result := runtime.NewArray64(size + 1)

	if !uniq {
		_ = result.Append(ctx, arg2)

		err = list.ForEach(ctx, func(ctx context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
			_ = result.Append(ctx, value)

			return runtime.True, nil
		})
		if err != nil {
			return runtime.None, err
		}

		return result, nil
	}

	_ = result.Append(ctx, arg2)

	err = list.ForEach(ctx, func(ctx context.Context, el runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		equal, err := runtime.EqualValues(ctx, el, arg2)
		if err != nil {
			return false, err
		}

		if !equal {
			_ = result.Append(ctx, el)
		}

		return true, nil
	})
	if err != nil {
		return runtime.None, err
	}

	return result, nil
}
