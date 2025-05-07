package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// MIN returns the smallest (arithmetic mean) of the values in array.
// @param {Int[] | Float[]} array - arrayList of numbers.
// @return {Float} - The smallest of the values in array.
func Min(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	arr, err := runtime.CastList(args[0])

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

	var min float64

	err = arr.ForEach(ctx, func(c context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		err = runtime.AssertNumber(value)

		if err != nil {
			return false, nil
		}

		fv := toFloat(value)

		if min > fv || idx == 0 {
			min = fv
		}

		return true, nil
	})

	if err != nil {
		return runtime.None, nil
	}

	return runtime.NewFloat(min), nil
}
