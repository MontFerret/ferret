package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// NTH returns the element of an array at a given position.
// It is the same as anyArray[position] for positive positions, but does not support negative positions.
// If position is negative or beyond the upper bound of the array, then NONE will be returned.
// @param {Any[]} array - An array with elements of arbitrary type.
// @param {Int} index - Position of desired element in array, positions start at 0.
// @return {Any} - The array element at the given position.
func Nth(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	list, idx, err := runtime.CastArgs2[runtime.List, runtime.Int](arg1, arg2)

	if err != nil {
		return runtime.None, err
	}

	// Handle negative index - return None as per documentation
	if idx < 0 {
		return runtime.None, nil
	}

	size, err := list.Length(ctx)

	if err != nil {
		return runtime.None, err
	}

	// Handle index beyond upper bound - return None as per documentation
	if idx >= size {
		return runtime.None, nil
	}

	return list.At(ctx, idx)
}
