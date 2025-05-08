package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// INTERSECTION return the intersection of all arrays specified.
// The result is an array of values that occur in all arguments.
// The element order is random. Duplicates are removed.
// @param {Any[], repeated} arrays - An arbitrary number of arrays as multiple arguments (at least 2).
// @return {Any[]} - A single array with only the elements, which exist in all provided arrays.
func Intersection(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	return sections(ctx, args, len(args))
}

func sections(ctx context.Context, args []runtime.Value, count int) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, runtime.MaxArgs); err != nil {
		return runtime.None, err
	}

	intersections := make(map[uint64][]runtime.Value)
	capacity := len(args)

	for _, i := range args {
		list, err := runtime.CastList(i)

		if err != nil {
			return runtime.None, err
		}

		err = list.ForEach(ctx, func(c context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
			h := value.Hash()

			bucket, exists := intersections[h]

			if !exists {
				bucket = make([]runtime.Value, 0, 5)
			}

			bucket = append(bucket, value)
			intersections[h] = bucket
			bucketLen := len(bucket)

			if bucketLen > capacity {
				capacity = bucketLen
			}

			return true, nil
		})
	}

	result := runtime.NewArray(capacity)
	required := count

	for _, bucket := range intersections {
		if len(bucket) == required {
			// It's safe to ignore the error here because we know that it's runtime.Array
			_ = result.Add(ctx, bucket[0])
		}
	}

	return result, nil
}
