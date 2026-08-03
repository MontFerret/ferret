package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// REMOVE_VALUE returns a new array with removed all occurrences of value in a given array.
// Optionally with a limit to the number of removals.
// @param {Any[]} array - Source array.
// @param {Any} value - Target value.
// @param {Int} [limit] - A limit to the number of removals.
// @return {Any[]} - A new array with removed all occurrences of value in a given array.
func RemoveValue(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.None, err
	}

	if len(args) == 2 {
		return removeValue2(ctx, args[0], args[1])
	}

	return removeValue3(ctx, args[0], args[1], args[2])
}

func removeValue2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return removeValue3(ctx, arg1, arg2, runtime.Int(-1))
}

func removeValue3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	arr, err := runtime.CastArg[runtime.List](arg1, 0)

	if err != nil {
		return runtime.None, err
	}

	limit, err := runtime.CastArg[runtime.Int](arg3, 2)
	if err != nil {
		return runtime.None, err
	}

	var counter runtime.Int

	return arr.Filter(ctx, func(ctx context.Context, item runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		remove := runtime.CompareValues(item, arg2) == 0

		if remove {
			counter++

			// If limit is 0, don't remove anything
			if limit == 0 {
				return true, nil
			}

			// If limit is negative or we haven't reached the limit, remove the item
			if limit < 0 || counter <= limit {
				return false, nil
			}
		}

		return true, nil
	})
}
