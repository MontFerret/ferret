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
	count := 0

	err = arr.ForEach(ctx, func(c context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		if !runtime.IsNumber(value) {
			return true, nil // Skip non-numeric values
		}

		fv := toFloat(value)

		if count == 0 || fv > res {
			res = fv
		}
		count++

		return true, nil
	})

	if err != nil {
		return runtime.None, err
	}

	if count == 0 {
		return runtime.None, nil
	}

	return runtime.NewFloat(res), nil
}
