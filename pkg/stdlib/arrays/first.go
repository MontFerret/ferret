package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// FIRST returns a first element from a given array.
// @param {Any[]} arr - Target array.
// @return {Any} - First element in a given array.
func First(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	arr, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	return arr.First(ctx)
}
