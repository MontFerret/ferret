package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// INTERSECTION return the intersection of all arrays specified.
// The result is an array of values that occur in all arguments.
// The element order is random. Duplicates are removed.
// @param {Any[], repeated} arrays - An arbitrary number of arrays as multiple arguments (at least 2).
// @return {Any[]} - A single array with only the elements, which exist in all provided arrays.
func Intersection(_ context.Context, args ...core.Value) (core.Value, error) {
	return sections(args, len(args))
}

func sections(args []core.Value, count int) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return values.None, err
	}

	intersections := make(map[uint64][]core.Value)
	capacity := len(args)

	for _, i := range args {
		err := core.ValidateType(i, types.Array)

		if err != nil {
			return values.None, err
		}

		arr := i.(*values.Array)

		arr.ForEach(func(value core.Value, idx int) bool {
			h := value.Hash()

			bucket, exists := intersections[h]

			if !exists {
				bucket = make([]core.Value, 0, 5)
			}

			bucket = append(bucket, value)
			intersections[h] = bucket
			bucketLen := len(bucket)

			if bucketLen > capacity {
				capacity = bucketLen
			}

			return true
		})
	}

	result := values.NewArray(capacity)
	required := count

	for _, bucket := range intersections {
		if len(bucket) == required {
			result.Push(bucket[0])
		}
	}

	return result, nil
}
