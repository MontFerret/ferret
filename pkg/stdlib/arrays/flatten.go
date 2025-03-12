package arrays

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// FLATTEN turns an array of arrays into a flat array.
// All array elements in array will be expanded in the result array.
// Non-array elements are added as they are.
// The function will recurse into sub-arrays up to the specified depth.
// Duplicates will not be removed.
// @param {Any[]} arr - Target array.
// @param {Int} [depth] - Depth level.
// @return {Any[]} - Flat array.
func Flatten(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 2)

	if err != nil {
		return core.None, err
	}

	err = core.AssertList(args[0])

	if err != nil {
		return core.None, err
	}

	arr := args[0].(*internal.Array)
	level := 1

	if len(args) > 1 {
		err = core.AssertInt(args[1])

		if err != nil {
			return core.None, err
		}

		level = int(args[1].(core.Int))
	}

	currentLevel := 0
	result := internal.NewArray(int(arr.Length()) * 2)
	var unwrap func(input *internal.Array)

	unwrap = func(input *internal.Array) {
		currentLevel++

		input.ForEach(func(value core.Value, idx int) bool {
			valueArr, ok := value.(*internal.Array)

			if !ok || currentLevel > level {
				result.Push(value)
			} else {
				unwrap(valueArr)
				currentLevel--
			}

			return true
		})
	}

	unwrap(arr)

	return result, nil
}
