package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// SORTED_UNIQUE sorts all elements in anyArray.
// The function will use the default comparison order for FQL value types.
// Additionally, the values in the result array will be made unique
// @param {Any[]} array - Target array.
// @return {Any[]} - Sorted array.
func SortedUnique(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Array)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)

	if arr.Length() == 0 {
		return values.NewArray(0), nil
	}

	return ToUniqueArray(arr.Sort()), nil
}
