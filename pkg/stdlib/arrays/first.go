package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// FIRST returns a first element from a given array.
// @param {Any[]} arr - Target array.
// @return {Any} - First element in a given array.
func First(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	arr, err := runtime.CastArg[runtime.List](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	return arr.First(ctx)
}
