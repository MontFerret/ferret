package arrays

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

/*
 * Return a new array with its elements reversed.
 * @param array (Array) - Target array.
 * @returns (Array) - A new array with its elements reversed.
 */
func Reverse(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.ArrayType)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)
	size := int(arr.Length())
	result := values.NewArray(size)

	for i := size - 1; i >= 0; i-- {
		result.Push(arr.Get(values.NewInt(i)))
	}

	return result, nil
}
