package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// hypot computes sqrt(x*x + y*y), avoiding unnecessary overflow and underflow.
// @param x {Int | Float} First coordinate.
// @param y {Int | Float} Second coordinate.
// @return {Float} The nonnegative distance; infinity takes precedence over NaN.
// @throws {TypeError} An argument is not numeric.
func Hypot(_ context.Context, x, y runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(x, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgValue(y, 1, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Hypot(toFloat(x), toFloat(y))), nil
}
