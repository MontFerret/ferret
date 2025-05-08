package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
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

	arr, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	value := args[1]
	var limit runtime.Int
	limit = -1

	if len(args) > 2 {
		arg3, err := runtime.CastInt(args[2])

		if err != nil {
			return runtime.None, err
		}

		limit = arg3
	}

	var counter runtime.Int

	return arr.Find(ctx, func(ctx context.Context, item runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		remove := runtime.CompareValues(item, value) == 0

		if remove {
			counter++

			if limit == -1 || counter <= limit {
				return false, nil
			}
		}

		return true, nil
	})
}
