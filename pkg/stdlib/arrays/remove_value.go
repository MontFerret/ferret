package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// REMOVE_VALUE returns a new array with removed all occurrences of value in a given array.
// Optionally with a limit to the number of removals.
// @param {Any[]} array - Source array.
// @param {Any} value - Target value.
// @param {Int} [limit] - A limit to the number of removals.
// @return {Any[]} - A new array with removed all occurrences of value in a given array.
func RemoveValue(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 3)

	if err != nil {
		return core.None, err
	}

	err = core.AssertList(args[0])

	if err != nil {
		return core.None, err
	}

	arr := args[0].(*internal.Array)
	value := args[1]
	limit := -1

	if len(args) > 2 {
		err = core.AssertInt(args[2])

		if err != nil {
			return core.None, err
		}

		limit = int(args[2].(core.Int))
	}

	result := internal.NewArray(int(arr.Length()))

	counter := 0
	arr.ForEach(func(item core.Value, idx int) bool {
		remove := core.CompareValues(item, value) == 0

		if remove {
			if counter == limit {
				result.Push(item)
			}

			counter++
		} else {
			result.Push(item)
		}

		return true
	})

	return result, nil
}
