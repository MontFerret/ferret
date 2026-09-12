package math

import (
	"context"
	"errors"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// percentile returns the nth percentile of the values in a given array.
// @param array {Int[] | Float[]} A list containing only Int and Float values; it is not mutated.
// @param number {Int} A number which must be between 0 (excluded) and 100 (included).
// @param method {String} "interpolation" uses linear interpolation; all other strings select nearest rank.
// @return {Int | Float} The selected number or interpolated Float, or NaN for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
func Percentile(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 3); err != nil {
		return runtime.None, err
	}

	if len(args) == 2 {
		return percentile2(ctx, args[0], args[1])
	}

	return percentile3(ctx, args[0], args[1], args[2])
}

// percentile returns the nth percentile of the values in a given array.
// @param array {Int[] | Float[]} A list containing only Int and Float values; it is not mutated.
// @param number {Int} A number which must be between 0 (excluded) and 100 (included).
// @return {Int | Float} The selected number or interpolated Float, or NaN for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
func percentile2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return percentile(ctx, arg1, arg2, "rank")
}

// percentile returns the nth percentile of the values in a given array.
// @param array {Int[] | Float[]} A list containing only Int and Float values; it is not mutated.
// @param number {Int} A number which must be between 0 (excluded) and 100 (included).
// @param method {String} "interpolation" uses linear interpolation; all other strings select nearest rank.
// @return {Int | Float} The selected number or interpolated Float, or NaN for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
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

	values, err := snapshotNumbers(ctx, arg1.(runtime.List))
	if err != nil {
		return runtime.None, err
	}

	size := len(values)
	if size == 0 {
		return runtime.NaN(), nil
	}

	num, err := runtime.CastArg[runtime.Int](arg2, 1)
	if err != nil {
		return runtime.None, err
	}

	percent := runtime.Float(num)

	if percent <= 0 || percent > 100 {
		return runtime.NaN(), errors.New("input is outside of range")
	}

	if err := sortNumbers(ctx, values, true); err != nil {
		return runtime.NaN(), err
	}

	switch method {
	case "interpolation":
		if size == 1 {
			return values[0], nil
		}

		pos := (float64(percent) / 100.0) * float64(size-1)
		lower := int(math.Floor(pos))
		upper := int(math.Ceil(pos))

		if lower == upper {
			return values[lower], nil
		}

		lowerVal := values[lower]
		upperVal := values[upper]

		frac := pos - float64(lower)
		result := toFloat(lowerVal) + (toFloat(upperVal)-toFloat(lowerVal))*frac

		return runtime.NewFloat(result), nil
	default:
		pos := math.Ceil((float64(percent) / 100.0) * float64(size))
		if pos < 1 || pos > float64(size) {
			return runtime.NaN(), errors.New("input is outside of range")
		}

		return values[int(pos)-1], nil
	}
}
