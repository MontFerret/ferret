package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// MINUS return the difference of all arrays specified.
// The order of the result array is undefined and should not be relied on. Duplicates will be removed.
// @param {Any[], repeated} arrays - An arbitrary number of arrays as multiple arguments (at least 2).
// @return {Any[]} - An array of values that occur in the first array, but not in any of the subsequent arrays.
func Minus(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	intersections := make(map[uint64]runtime.Value)
	var capacity runtime.Int

	for idx, i := range args {
		idx := idx
		list, err := runtime.CastList(i)

		if err != nil {
			return runtime.None, err
		}

		err = list.ForEach(ctx, func(c context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
			h := value.Hash()

			// first array, fill out the map
			if idx == 0 {
				size, err := list.Length(c)

				if err != nil {
					return false, err
				}

				capacity = size
				intersections[h] = value

				return true, nil
			}

			_, exists := intersections[h]

			// if it exists in the first array, remove it
			if exists {
				delete(intersections, h)
			}

			return true, nil
		})

		if err != nil {
			return runtime.None, err
		}
	}

	result := runtime.NewArray64(capacity)

	for _, item := range intersections {
		_ = result.Add(ctx, item)
	}

	return result, nil
}
