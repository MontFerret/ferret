package arrays

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// UNIQUE returns all unique elements from a given array.
// @param {Any[]} array - Target array.
// @return {Any[]} - New array without duplicates.
func Unique(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	list, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	return ToUniqueList(ctx, list)
}
