package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

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

	if err := runtime.ValidateArgValueAt(args, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	arr := args[0].(runtime.List)
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

	method := "rank"

	if len(args) > 2 {
		if err := runtime.ValidateType(args[2], runtime.TypeString); err != nil {
			return runtime.None, err
		}

		if args[2].String() == "interpolation" {
			method = "interpolation"
		}
	}

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
