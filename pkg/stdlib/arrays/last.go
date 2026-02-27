package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// LAST returns the last element of an array.
// @param {Any[]} array - The target array.
// @return {Any} - Last element of an array.
func Last(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	return list.Last(ctx)
}
