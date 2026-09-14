package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// cbrt returns the real cube root of a number, including negative numbers.
// @param value {Int | Float} Input number.
// @return {Float} The cube root, preserving NaN, infinities, and signed zero.
// @throws {TypeError} The argument is not numeric.
func Cbrt(_ context.Context, value runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(value, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Cbrt(toFloat(value))), nil
}
