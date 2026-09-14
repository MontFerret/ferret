package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// trunc removes the fractional part of a number by rounding toward zero.
// @param value {Int | Float} Input number.
// @return {Int | Float} The input Int unchanged, or a truncated Float preserving NaN, infinities, and signed zero.
// @throws {TypeError} The argument is not numeric.
func Trunc(_ context.Context, value runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(value, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	if _, ok := value.(runtime.Int); ok {
		return value, nil
	}

	return runtime.NewFloat(math.Trunc(toFloat(value))), nil
}
