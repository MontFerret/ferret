package arrays

import (
	"github.com/MontFerret/ferret/pkg/runtime"
)

// SHIFT returns a new array without the first element.
// @param {Any[]} array - Target array.
// @return {Any[]} - Copy of an array without the first element.
func Shift(ctx runtime.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	list, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	return list.Find(ctx, func(ctx runtime.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		return idx != 0, nil
	})
}
