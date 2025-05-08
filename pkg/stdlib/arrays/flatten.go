package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// FLATTEN turns an array of arrays into a flat array.
// All array elements in array will be expanded in the result array.
// Non-array elements are added as they are.
// The function will recurse into sub-arrays up to the specified depth.
// Duplicates will not be removed.
// @param {Any[]} arr - Target array.
// @param {Int} [depth] - Depth level.
// @return {Any[]} - Flat array.
func Flatten(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 2); err != nil {
		return runtime.None, err
	}

	list, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	size, err := list.Length(ctx)

	if err != nil {
		return runtime.None, err
	}

	var level runtime.Int
	level = 1

	if len(args) > 1 {
		arg1, err := runtime.CastInt(args[1])

		if err != nil {
			return runtime.None, err
		}

		level = arg1
	}

	var currentLevel runtime.Int
	result := runtime.NewArray64(size * 2)
	var unwrap func(input runtime.List) error

	unwrap = func(input runtime.List) error {
		currentLevel++

		return input.ForEach(ctx, func(c context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
			valueArr, ok := value.(runtime.List)

			if !ok || currentLevel > level {
				_ = result.Add(c, value)
			} else {
				if err := unwrap(valueArr); err != nil {
					return false, err
				}

				currentLevel--
			}

			return true, nil
		})
	}

	if err := unwrap(list); err != nil {
		return runtime.None, err
	}

	return result, nil
}
