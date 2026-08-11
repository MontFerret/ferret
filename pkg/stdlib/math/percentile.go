package math

import (
	"context"
	"errors"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// PERCENTILE returns the nth percentile of the values in a given array.
// @param array {Int[] | Float[]} arrayList of numbers.
// @param number {Int} A number which must be between 0 (excluded) and 100 (included).
// @param method {String} "rank" or "interpolation".
// @return {Float} The nth percentile, or null if the array is empty or only null values are contained in it or the percentile cannot be calculated.
func Percentile(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.None, err
	}

	if len(args) == 2 {
		return percentile2(ctx, args[0], args[1])
	}

	return percentile3(ctx, args[0], args[1], args[2])
}

// PERCENTILE returns the nth percentile of the values in a given array.
// @param array {Int[] | Float[]} arrayList of numbers.
// @param number {Int} A number which must be between 0 (excluded) and 100 (included).
// @return {Float} The nth percentile, or null if the array is empty or only null values are contained in it or the percentile cannot be calculated.
func percentile2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return percentile(ctx, arg1, arg2, "rank")
}

// PERCENTILE returns the nth percentile of the values in a given array.
// @param array {Int[] | Float[]} arrayList of numbers.
// @param number {Int} A number which must be between 0 (excluded) and 100 (included).
// @param method {String} "rank" or "interpolation".
// @return {Float} The nth percentile, or null if the array is empty or only null values are contained in it or the percentile cannot be calculated.
func percentile3(ctx context.Context, arg1, arg2, arg3 runtime.Value) (runtime.Value, error) {
	method, err := runtime.CastArg[runtime.String](arg3, 2)
	if err != nil {
		return runtime.None, err
	}

	return percentile(ctx, arg1, arg2, method.String())
}

func percentile(ctx context.Context, arg1, arg2 runtime.Value, method string) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg1, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	arr := arg1.(runtime.List)
	size, err := arr.Length(ctx)

	if err != nil {
		return runtime.None, err
	}

	if size == 0 {
		return runtime.NewFloat(math.NaN()), nil
	}

	num, err := runtime.CastArg[runtime.Int](arg2, 1)

	if err != nil {
		return runtime.None, err
	}

	percent := runtime.Float(num)

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

	switch method {
	case "interpolation":
		if size == 1 {
			return sorted.At(ctx, 0)
		}

		pos := (float64(percent) / 100.0) * float64(size-1)
		lower := int(math.Floor(pos))
		upper := int(math.Ceil(pos))

		if lower == upper {
			return sorted.At(ctx, runtime.Int(lower))
		}

		lowerVal, err := sorted.At(ctx, runtime.Int(lower))
		if err != nil {
			return runtime.None, err
		}

		upperVal, err := sorted.At(ctx, runtime.Int(upper))
		if err != nil {
			return runtime.None, err
		}

		if err := runtime.AssertNumber(lowerVal); err != nil {
			return runtime.None, err
		}

		if err := runtime.AssertNumber(upperVal); err != nil {
			return runtime.None, err
		}

		frac := pos - float64(lower)
		result := toFloat(lowerVal) + (toFloat(upperVal)-toFloat(lowerVal))*frac

		return runtime.NewFloat(result), nil
	default:
		pos := math.Ceil((float64(percent) / 100.0) * float64(size))
		if pos < 1 || pos > float64(size) {
			return runtime.NaN(), errors.New("input is outside of range")
		}

		return sorted.At(ctx, runtime.Int(pos-1))
	}
}
