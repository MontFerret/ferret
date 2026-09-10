package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Remove returns a shallow copy without any values equal to the target.
// @param array {Any[]} Source array.
// @param value {Any} Value whose occurrences are removed.
// @return {Any[]} Remaining elements in their original order.
func Remove(ctx context.Context, array, value runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](array, 0)
	if err != nil {
		return runtime.None, err
	}

	return removeLimited(ctx, list, value, -1)
}

func removeLimited(ctx context.Context, arr runtime.List, target runtime.Value, limit runtime.Int) (runtime.Value, error) {
	var counter runtime.Int

	return arr.Filter(ctx, func(ctx context.Context, item runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		remove, err := runtime.EqualValues(ctx, item, target)
		if err != nil {
			return false, err
		}

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
