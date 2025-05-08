package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// NTH returns the element of an array at a given position.
// It is the same as anyArray[position] for positive positions, but does not support negative positions.
// If position is negative or beyond the upper bound of the array, then NONE will be returned.
// @param {Any[]} array - An array with elements of arbitrary type.
// @param {Int} index - Position of desired element in array, positions start at 0.
// @return {Any} - The array element at the given position.
func Nth(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 2); err != nil {
		return runtime.None, err
	}

	list, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	idx, err := runtime.CastInt(args[1])

	if err != nil {
		return runtime.None, err
	}

	return list.Get(ctx, idx)
}
