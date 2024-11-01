package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// LAST returns the last element of an array.
// @param {Any[]} array - The target array.
// @return {Any} - Last element of an array.
func Last(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = values.AssertArray(args[0])

	if err != nil {
		return values.None, nil
	}

	arr := args[0].(*values.Array)

	return arr.Get(arr.Length() - 1), nil
}
