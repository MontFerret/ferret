package arrays

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// PUSH create a new array with appended value.
// @param {Any[]} array - Source array.
// @param {Any} value - Target value.
// @param {Boolean} [unique=False] - Read indicating whether to do uniqueness check.
// @return {Any[]} - A new array with appended value.
func Push(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 2, 3); err != nil {
		return values.None, err
	}

	arr, err := values.CastArray(args[0])

	if err != nil {
		return values.None, err
	}
	
	value := args[1]
	uniq := false

	if len(args) > 2 {
		err = values.AssertBoolean(args[2])

		if err != nil {
			return values.None, err
		}

		uniq = values.Compare(args[2], values.True) == 0
	}

	result := values.NewArray(int(arr.Length() + 1))
	push := true

	arr.ForEach(func(item core.Value, idx int) bool {
		if uniq && push {
			push = values.Compare(item, value) != 0
		}

		result.Push(item)

		return true
	})

	if push {
		result.Push(value)
	}

	return result, nil
}
