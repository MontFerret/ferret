package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// MEDIAN returns the median of the values in array.
// @param {Int[] | Float[]} array - arrayList of numbers.
// @return {Float} - The median of the values in array.
func Median(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	arr, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	sorted, err := arr.SortDesc(ctx)

	if err != nil {
		return runtime.None, err
	}

	size, err := sorted.Length(ctx)

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
