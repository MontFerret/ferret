package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
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

	type occurrence struct {
		value      runtime.Value
		lastSource int
		count      int
	}

	intersections := make(map[uint64][]*occurrence)
	capacity := len(args)

	for i, arg := range args {
		list, err := runtime.CastArg[runtime.List](arg, i)

		if err != nil {
			return runtime.None, err
		}

		err = list.ForEach(ctx, func(c context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
			h := value.Hash()
			bucket, exists := intersections[h]
			if exists {
				for _, entry := range bucket {
					equal, err := runtime.EqualValues(c, entry.value, value)
					if err != nil {
						return false, err
					}
					if !equal {
						continue
					}

					if entry.lastSource != i {
						entry.lastSource = i
						entry.count++
					}

					return true, nil
				}
			}

			bucket = append(bucket, &occurrence{value: value, lastSource: i, count: 1})
			intersections[h] = bucket

			return true, nil
		})

		if err != nil {
			return runtime.None, err
		}
	}

	result := runtime.NewArray(capacity)

	for _, bucket := range intersections {
		for _, entry := range bucket {
			if entry.count == count {
				// It's safe to ignore the error here because result is a runtime.Array.
				_ = result.Append(ctx, entry.value)
			}
		}
	}

	return result, nil
}
