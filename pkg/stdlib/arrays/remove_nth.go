package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// REMOVE_NTH returns a new array without an element by a given position.
// @param {Any[]} array - Source array.
// @param {Int} position - Target element position.
// @return {Any[]} - A new array without an element by a given position.
func RemoveNth(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 2); err != nil {
		return runtime.None, err
	}

	list, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	index, err := runtime.CastInt(args[1])

	if err != nil {
		return runtime.None, err
	}

	next := list.Copy().(runtime.List)

	//if err != nil {
	//	return runtime.None, err
	//}

	_, err = next.RemoveAt(ctx, index)

	return next, err
}
