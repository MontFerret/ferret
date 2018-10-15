package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/collections"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Unique returns all unique elements from a given array.
// @param array (Array) - Target array.
// @returns (Array) - New array without duplicates.
func Unique(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.ArrayType)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)

	if arr.Length() == 0 {
		return values.NewArray(0), nil
	}

	iterator, err := collections.NewUniqueIterator(
		collections.NewArrayIterator(arr),
	)

	if err != nil {
		return values.None, err
	}

	return collections.ToArray(iterator)
}
