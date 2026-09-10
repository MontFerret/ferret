package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Contains reports whether an equal value exists in the array.
// @param array {Any[]} Source array.
// @param value {Any} Value to find.
// @return {Boolean} Whether the value exists.
func Contains(ctx context.Context, array, value runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](array, 0)
	if err != nil {
		return runtime.None, err
	}

	index, err := list.IndexOf(ctx, value)
	if err != nil {
		return runtime.None, err
	}

	return runtime.Boolean(index >= 0), nil
}

// IndexOf finds the first equal value, returning -1 when it is absent.
// @param array {Any[]} Source array.
// @param value {Any} Value to find.
// @return {Int} First matching index, or -1.
func IndexOf(ctx context.Context, array, value runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](array, 0)
	if err != nil {
		return runtime.None, err
	}

	index, err := list.IndexOf(ctx, value)
	if err != nil {
		return runtime.None, err
	}

	return index, nil
}
