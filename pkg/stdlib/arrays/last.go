package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// LAST returns the last element of an array.
// @param {Any[]} array - The target array.
// @return {Any} - Last element of an array.
func Last(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	list, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, nil
	}

	return list.Last(ctx)
}
