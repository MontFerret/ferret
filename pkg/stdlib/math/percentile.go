package math

import (
	"context"
	"errors"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// percentile returns the nth percentile of the values in a given array.
// @param array {Any[]} A list whose Int and Float elements are used; other elements are ignored. It is not mutated.
// @param number {Int} A number which must be between 0 (excluded) and 100 (included).
// @param method {String} "interpolation" uses linear interpolation; all other strings select nearest rank.
// @return {Int | Float} The selected number or interpolated Float, or NaN when no numbers remain.
// @throws {TypeError} An argument has an invalid type.
// @deprecated Use math::percentile instead.
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
// @param array {Any[]} A list whose Int and Float elements are used; other elements are ignored. It is not mutated.
// @param number {Int} A number which must be between 0 (excluded) and 100 (included).
// @return {Int | Float} The selected number or interpolated Float, or NaN when no numbers remain.
// @throws {TypeError} An argument has an invalid type.
// @deprecated Use math::percentile instead.
func percentile2(ctx context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	return percentile(ctx, arg1, arg2, "rank")
}

// percentile returns the nth percentile of the values in a given array.
// @param array {Any[]} A list whose Int and Float elements are used; other elements are ignored. It is not mutated.
// @param number {Int} A number which must be between 0 (excluded) and 100 (included).
// @param method {String} "interpolation" uses linear interpolation; all other strings select nearest rank.
// @return {Int | Float} The selected number or interpolated Float, or NaN when no numbers remain.
// @throws {TypeError} An argument has an invalid type.
// @deprecated Use math::percentile instead.
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

	values, err := snapshotNumbers(ctx, arg1.(runtime.List), filterNonNumbers)
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
		return interpolatePercentile(values, float64(percent)), nil
	default:
		pos := math.Ceil((float64(percent) / 100.0) * float64(size))
		if pos < 1 || pos > float64(size) {
			return runtime.NaN(), errors.New("input is outside of range")
		}

		return values[int(pos)-1], nil
	}
}

// percentile returns the percentile using linear interpolation at (p / 100) * (N - 1).
// @param values {Int[] | Float[]} Numeric list; non-numeric elements are rejected. The list is not mutated.
// @param p {Int | Float} Finite percentile from zero through 100, inclusive; validated even for an empty list.
// @return {Int | Float} The selected number or interpolated Float, or NaN for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
// @throws {InvalidArgument} The percentile is non-finite or outside zero through 100.
func canonicalPercentile(ctx context.Context, arg, percent runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgValue(percent, 1, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	p := toFloat(percent)
	if math.IsNaN(p) || math.IsInf(p, 0) || p < 0 || p > 100 {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "percentile must be finite and between 0 and 100"), 1)
	}

	values, err := snapshotNumbers(ctx, arg.(runtime.List), strictNumbers)
	if err != nil {
		return runtime.None, err
	}

	if len(values) == 0 {
		return runtime.NaN(), nil
	}

	if err := sortNumbers(ctx, values, true); err != nil {
		return runtime.None, err
	}

	return interpolatePercentile(values, p), nil
}

// Callers supply a nonempty, ascending numeric snapshot and a valid percentile.
func interpolatePercentile(values []runtime.Value, percent float64) runtime.Value {
	pos := (percent / 100) * float64(len(values)-1)
	lower := int(math.Floor(pos))
	upper := int(math.Ceil(pos))
	if lower == upper {
		return values[lower]
	}

	fraction := pos - float64(lower)
	lowerValue, upperValue := toFloat(values[lower]), toFloat(values[upper])

	return runtime.NewFloat(lowerValue + (upperValue-lowerValue)*fraction)
}
