package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// UNIQUE returns all unique elements from a given array.
// @param {Any[]} array - Target array.
// @return {Any[]} - New array without duplicates.
func Unique(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return core.None, err
	}

	err = core.AssertList(args[0])

	if err != nil {
		return core.None, err
	}

	arr := args[0].(*internal.Array)

	if arr.Length() == 0 {
		return internal.NewArray(0), nil
	}

	return ToUniqueArray(arr), nil
}
