package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// MAX returns the greatest (arithmetic mean) of the values in array.
// @param {Int[] | Float[]} array - arrayList of numbers.
// @return {Float} - The greatest of the values in array.
func Max(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	arr, err := runtime.CastList(arg)

	if err != nil {
		return runtime.None, err
	}

	size, err := arr.Length(ctx)

	if err != nil {
		return runtime.None, err
	}

	if size == 0 {
		return runtime.None, nil
	}

	var res float64

	err = arr.ForEach(ctx, func(c context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		fv := toFloat(value)

		if fv > res {
			res = fv
		}

		return true, nil
	})

	if err != nil {
		return runtime.None, nil
	}

	return runtime.NewFloat(res), nil
}
