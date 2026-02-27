package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// REMOVE_NTH returns a new array without an element by a given position.
// @param {Any[]} array - Source array.
// @param {Int} position - Target element position.
// @return {Any[]} - A new array without an element by a given position.
func RemoveNth(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](arg1, 0)

	if err != nil {
		return runtime.None, err
	}

	index, err := runtime.CastArg[runtime.Int](arg2, 1)

	if err != nil {
		return runtime.None, err
	}

	next := list.Copy().(runtime.List)

	_, err = next.RemoveAt(ctx, index)

	return next, err
}
