package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// SHIFT returns a new array without the first element.
// @param {Any[]} array - Target array.
// @return {Any[]} - Copy of an array without the first element.
func Shift(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](arg, 0)

	if err != nil {
		return runtime.None, err
	}

	return list.Filter(ctx, func(ctx context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		return idx != 0, nil
	})
}
