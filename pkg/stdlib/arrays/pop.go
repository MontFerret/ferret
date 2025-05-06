package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// POP returns a new array without last element.
// @param {Any[]} array - Target array.
// @return {Any[]} - Copy of an array without last element.
func Pop(_ context.Context, args ...core.Value) (core.Value, error) {
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
	lastIdx := length - 1

	arr.ForEach(func(value core.Value, idx int) bool {
		if idx == lastIdx {
			return false
		}

		result.Push(value)

		return true
	})

	return result, nil
}
