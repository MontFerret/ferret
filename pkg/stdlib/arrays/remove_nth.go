package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// RemoveNth returns a new array without an element by a given position.
// @param array (Array) - Source array.
// @param position (Int) - Target element position.
// @return (Array) - A new array without an element by a given position.
func RemoveNth(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Array)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], types.Int)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)
	index := int(args[1].(values.Int))
	result := values.NewArray(int(arr.Length() - 1))

	arr.ForEach(func(value core.Value, idx int) bool {
		if idx != index {
			result.Push(value)
		}

		return true
	})

	return result, nil
}
