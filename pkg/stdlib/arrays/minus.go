package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// MINUS return the difference of all arrays specified.
// The order of the result array is undefined and should not be relied on. Duplicates will be removed.
// @param {Any[], repeated} arrays - An arbitrary number of arrays as multiple arguments (at least 2).
// @return {Any[]} - An array of values that occur in the first array, but not in any of the subsequent arrays.
func Minus(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, core.MaxArgs)

	if err != nil {
		return values.None, err
	}

	intersections := make(map[uint64]core.Value)
	capacity := values.NewInt(0)

	for idx, i := range args {
		idx := idx
		err := core.ValidateType(i, types.Array)

		if err != nil {
			return values.None, err
		}

		arr := i.(*values.Array)

		arr.ForEach(func(value core.Value, _ int) bool {
			h := value.Hash()

			// first array, fill out the map
			if idx == 0 {
				capacity = arr.Length()
				intersections[h] = value

				return true
			}

			_, exists := intersections[h]

			// if it exists in the first array, remove it
			if exists {
				delete(intersections, h)
			}

			return true
		})
	}

	result := values.NewArray(int(capacity))

	for _, item := range intersections {
		result.Push(item)
	}

	return result, nil
}
