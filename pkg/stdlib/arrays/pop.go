package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// POP returns a new array without last element.
// @param {Any[]} array - Target array.
// @return {Any[]} - Copy of an array without last element.
func Pop(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	arr, err := runtime.CastArg[runtime.List](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	size, err := arr.Length(ctx)

	if err != nil {
		return runtime.None, err
	}

	return arr.Slice(ctx, 0, size-1)
}
