package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// MEDIAN returns the median of the values in array.
// @param {Int[] | Float[]} array - arrayList of numbers.
// @return {Float} - The median of the values in array.
func Median(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
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

	// Filter numeric values into a new array
	numericValues := runtime.NewArray(0)
	err = arr.ForEach(ctx, func(c context.Context, value runtime.Value, idx runtime.Int) (runtime.Boolean, error) {
		if runtime.IsNumber(value) {
			err := numericValues.Add(ctx, value)
			if err != nil {
				return false, err
			}
		}
		return true, nil
	})

	if err != nil {
		return runtime.None, err
	}

	numericSize, err := numericValues.Length(ctx)
	if err != nil {
		return runtime.None, err
	}

	if numericSize == 0 {
		return runtime.None, nil
	}

	sorted := numericValues.Copy().(runtime.List)

	if err := runtime.SortDesc(ctx, sorted); err != nil {
		return runtime.None, err
	}

	size, err = sorted.Length(ctx)

	if err != nil {
		return runtime.None, err
	}

	switch {
	case size == 0:
		return runtime.NewFloat(math.NaN()), nil
	case size%2 == 0:
		sliced, err := sorted.Slice(ctx, 0, size)

		if err != nil {
			return runtime.None, err
		}

		return mean(ctx, sliced)
	default:
		return sorted.Get(ctx, size/2)
	}
}
