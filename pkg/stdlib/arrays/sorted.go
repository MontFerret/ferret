package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// SORTED sorts all elements in anyArray.
// The function will use the default comparison order for FQL value types.
// @param {Any[]} array - Target array.
// @return {Any[]} - Sorted array.
func Sorted(_ context.Context, args ...core.Value) (core.Value, error) {
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

	return arr.Sort(), nil
}
