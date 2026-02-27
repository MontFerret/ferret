package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// UNIQUE returns all unique elements from a given array.
// @param {Any[]} array - Target array.
// @return {Any[]} - New array without duplicates.
func Unique(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	return ToUniqueList(ctx, list)
}
