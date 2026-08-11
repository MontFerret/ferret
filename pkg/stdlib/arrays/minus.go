package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// MINUS return the difference of all arrays specified.
// The order of the result array is undefined and should not be relied on. Duplicates will be removed.
// @param arrays {Any[], repeated} An arbitrary number of arrays as multiple arguments (at least 2).
// @return {Any[]} An array of values that occur in the first array, but not in any of the subsequent arrays.
func Minus(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	intersections := make(map[uint64][]runtime.Value)
	var capacity runtime.Int

	for idx := range args {
		idx := idx
		list, err := runtime.CastArgAt[runtime.List](args, idx)

		if err != nil {
			return runtime.None, err
		}

		err = list.ForEach(ctx, func(c context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
			h := value.Hash()
			bucket := intersections[h]

			// first array, fill out the map
			if idx == 0 {
				size, err := list.Length(c)

				if err != nil {
					return false, err
				}

				capacity = size
				equalIdx, err := findEqualValue(c, bucket, value)
				if err != nil {
					return false, err
				}
				if equalIdx < 0 {
					intersections[h] = append(bucket, value)
				}

				return true, nil
			}

			match, err := findEqualValue(c, bucket, value)
			if err != nil {
				return false, err
			}
			if match >= 0 {
				bucket = append(bucket[:match], bucket[match+1:]...)
				if len(bucket) == 0 {
					delete(intersections, h)
				} else {
					intersections[h] = bucket
				}
			}

			return true, nil
		})

		if err != nil {
			return runtime.None, err
		}
	}

	result := runtime.NewArray64(capacity)

	for _, bucket := range intersections {
		for _, item := range bucket {
			_ = result.Append(ctx, item)
		}
	}

	return result, nil
}
