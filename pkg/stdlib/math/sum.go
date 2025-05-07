package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// SUM returns the sum of the values in a given array.
// @param {Int[] | Float[]} numbers - arrayList of numbers.
// @return {Float} - The sum of the values.
func Sum(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
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

	var sum float64

	err = arr.ForEach(ctx, func(c context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {

		if err = runtime.AssertNumber(value); err != nil {
			return false, err
		}

		sum += toFloat(value)

		return true, nil
	})

	if err != nil {
		return runtime.None, nil
	}

	return runtime.NewFloat(sum), nil
}
