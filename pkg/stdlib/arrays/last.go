package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// LAST returns the last element of an array.
// @param {Any[]} array - The target array.
// @return {Any} - Last element of an array.
func Last(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return core.None, err
	}

	err = core.AssertList(args[0])

	if err != nil {
		return core.None, nil
	}

	arr := args[0].(*internal.Array)

	return arr.Get(arr.Length() - 1), nil
}
