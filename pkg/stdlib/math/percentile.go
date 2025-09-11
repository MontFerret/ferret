package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"

	"errors"
)

// PERCENTILE returns the nth percentile of the values in a given array.
// @param {Int[] | Float[]} array - arrayList of numbers.
// @param {Int} number - A number which must be between 0 (excluded) and 100 (included).
// @param {String} [method="rank"] - "rank" or "interpolation".
// @return {Float} - The nth percentile, or null if the array is empty or only null values are contained in it or the percentile cannot be calculated.
func Percentile(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
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
		return runtime.NewFloat(math.NaN()), nil
	}

	num, err := runtime.CastInt(args[1])

	if err != nil {
		return runtime.None, err
	}

	percent := runtime.Float(num)

	// TODO: Implement different methods
	//method := "rank"
	//
	//if len(args) > 2 {
	//	err = runtime.ValidateType(args[2], runtime.StringType)
	//
	//	if err != nil {
	//		return values.None, err
	//	}
	//
	//	if args[2].String() == "interpolation" {
	//		method = "interpolation"
	//	}
	//}

	if percent <= 0 || percent > 100 {
		return runtime.NaN(), errors.New("input is outside of range")
	}

	sorted := arr.Copy().(runtime.List)

	//if err != nil {
	//	return runtime.NaN(), err
	//}

	if err := runtime.SortAsc(ctx, sorted); err != nil {
		return runtime.NaN(), err
	}

	// Multiply percent by length of input
	l := runtime.Float(size)
	index := (percent / 100) * l
	even := runtime.Float(runtime.Int(index))

	// Check if the index is a whole number
	switch {
	case index == even:
		i := runtime.Int(index)
		return sorted.Get(ctx, i-1)
	case index > 1:
		// Convert float to int via truncation
		i := runtime.Int(index)
		// Find the average of the index and following values
		aVal, err := sorted.Get(ctx, i-1)

		if err != nil {
			return runtime.None, err
		}

		bVal, err := sorted.Get(ctx, i)

		if err != nil {
			return runtime.None, err
		}

		return mean(ctx, runtime.NewArrayWith(aVal, bVal))
	default:
		return runtime.NaN(), errors.New("input is outside of range")
	}
}
