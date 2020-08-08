package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
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
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Array)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)
	level := 1

	if len(args) > 1 {
		err = core.ValidateType(args[1], types.Int)

		if err != nil {
			return values.None, err
		}

		level = int(args[1].(values.Int))
	}

	currentLevel := 0
	result := values.NewArray(int(arr.Length()) * 2)
	var unwrap func(input *values.Array)

	unwrap = func(input *values.Array) {
		currentLevel++

		input.ForEach(func(value core.Value, idx int) bool {
			if value.Type() != types.Array || currentLevel > level {
				result.Push(value)
			} else {
				unwrap(value.(*values.Array))
				currentLevel--
			}

			return true
		})
	}

	unwrap(arr)

	return result, nil
}
