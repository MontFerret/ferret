package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
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
		return values.None, err
	}

	err = values.AssertArray(args[0])

	if err != nil {
		return values.None, err
	}

	err = values.AssertInt(args[1])

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)
	idx := args[1].(values.Int)

	return arr.Get(int(idx)), nil
}
