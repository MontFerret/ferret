package arrays

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// NTH returns the element of an array at a given position.
// It is the same as anyArray[position] for positive positions, but does not support negative positions.
// If position is negative or beyond the upper bound of the array, then NONE will be returned.
// @param {Any[]} array - An array with elements of arbitrary type.
// @param {Int} index - Position of desired element in array, positions start at 0.
// @return {Any} - The array element at the given position.
func Nth(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return core.None, err
	}

	err = core.AssertList(args[0])

	if err != nil {
		return core.None, err
	}

	err = core.AssertInt(args[1])

	if err != nil {
		return core.None, err
	}

	arr := args[0].(*internal.Array)
	idx := args[1].(core.Int)

	return arr.Get(int(idx)), nil
}
