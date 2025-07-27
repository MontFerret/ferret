package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// AVERAGE Returns the average (arithmetic mean) of the values in array.
// @param {Int[] | Float[]} array - arrayList of numbers.
// @return {Float} - The average of the values in array.
func Average(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	arr, err := runtime.CastList(arg)

	if err != nil {
		return runtime.None, err
	}

	var (
		sum   float64
		count int
	)

	err = arr.ForEach(ctx, func(c context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		if !runtime.IsNumber(value) {
			return true, nil // skip non-numbers/nulls
		}

		sum += toFloat(value)
		count++

		return true, nil
	})

	if err != nil {
		return runtime.None, err
	}

	if count == 0 {
		return runtime.ZeroFloat, nil
	}

	return runtime.Float(sum / float64(count)), nil
}
