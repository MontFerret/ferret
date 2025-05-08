package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// POP returns a new array without last element.
// @param {Any[]} array - Target array.
// @return {Any[]} - Copy of an array without last element.
func Pop(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	arr, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	size, err := arr.Length(ctx)

	if err != nil {
		return runtime.None, err
	}

	return arr.Slice(ctx, 0, size-1)
}
