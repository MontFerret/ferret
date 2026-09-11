package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RemoveAt returns a new array without an element by a given position.
// Negative and out-of-range positions return an unchanged shallow copy.
// @param array {Any[]} Source array.
// @param position {Int} Target element position.
// @return {Any[]} A new array without an element by a given position.
func RemoveAt(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](arg1, 0)
	if err != nil {
		return runtime.None, err
	}

	index, err := runtime.CastArg[runtime.Int](arg2, 1)
	if err != nil {
		return runtime.None, err
	}

	next, err := runtime.Copy(list)
	if err != nil {
		return runtime.None, err
	}

	if index < 0 {
		return next, nil
	}

	size, err := next.Length(ctx)
	if err != nil {
		return runtime.None, err
	}

	if index >= size {
		return next, nil
	}

	_, err = next.RemoveAt(ctx, index)

	return next, err
}
