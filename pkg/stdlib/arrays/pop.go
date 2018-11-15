package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Pop returns a new array without last element.
// @param array (Array) - Target array.
// @returns (Array) - Copy of an array without last element.
func Pop(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.ArrayType)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)

	length := int64(arr.Length())
	result := values.NewArray(length)
	lastIdx := length - 1

	arr.ForEach(func(value core.Value, idx int) bool {
		if int64(idx) == lastIdx {
			return false
		}

		result.Push(value)

		return true
	})

	return result, nil
}
