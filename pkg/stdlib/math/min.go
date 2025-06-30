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

	err = arr.ForEach(ctx, func(c context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		fv := toFloat(value)

		if res > fv || idx == 0 {
			res = fv
		}

		return true, nil
	})

	if err != nil {
		return runtime.None, nil
	}

	return runtime.NewFloat(res), nil
}
