package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// MIN returns the smallest (arithmetic mean) of the values in array.
// @param {Int[] | Float[]} array - arrayList of numbers.
// @return {Float} - The smallest of the values in array.
func Min(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
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
	hasNumericValues := false

	err = arr.ForEach(ctx, func(c context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		if !runtime.IsNumber(value) {
			return true, nil // Skip non-numeric values
		}

		fv := toFloat(value)

		if !hasNumericValues || res > fv {
			res = fv
			hasNumericValues = true
		}

		return true, nil
	})

	if err != nil {
		return runtime.None, err
	}

	if !hasNumericValues {
		return runtime.None, nil
	}

	return runtime.NewFloat(res), nil
}
