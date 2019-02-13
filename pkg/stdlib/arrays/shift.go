package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Shift returns a new array without the first element.
// @param array (Array) - Target array.
// @returns (Array) - Copy of an array without the first element.
func Shift(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Array)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)

	length := int(arr.Length())
	result := values.NewArray(length)

	arr.ForEach(func(value core.Value, idx int) bool {
		if idx != 0 {
			result.Push(value)
		}

		return true
	})

	return result, nil
}
