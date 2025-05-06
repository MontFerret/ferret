package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// SHIFT returns a new array without the first element.
// @param {Any[]} array - Target array.
// @return {Any[]} - Copy of an array without the first element.
func Shift(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return core.None, err
	}

	err = core.AssertList(args[0])

	if err != nil {
		return core.None, err
	}

	arr := args[0].(*internal.Array)

	length := int(arr.Length())
	result := internal.NewArray(length)

	arr.ForEach(func(value core.Value, idx int) bool {
		if idx != 0 {
			result.Push(value)
		}

		return true
	})

	return result, nil
}
