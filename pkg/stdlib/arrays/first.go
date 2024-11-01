package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// FIRST returns a first element from a given array.
// @param {Any[]} arr - Target array.
// @return {Any} - First element in a given array.
func First(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = values.AssertArray(args[0])

	if err != nil {
		return values.None, nil
	}

	arr := args[0].(*values.Array)

	return arr.Get(0), nil
}
