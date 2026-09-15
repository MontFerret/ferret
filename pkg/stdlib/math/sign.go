package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// sign returns the sign of a number, including positive and negative infinity.
// @param value {Int | Float} Input number; must not be NaN.
// @return {Int} Minus one for negative values, zero for either signed zero, or one for positive values.
// @throws {TypeError} The argument is not numeric.
// @throws {InvalidArgument} The argument is NaN and has no ordered sign.
func Sign(_ context.Context, value runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(value, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	number := toFloat(value)
	if math.IsNaN(number) {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "sign is undefined for NaN"), 0)
	}

	if number < 0 {
		return runtime.Int(-1), nil
	}

	if number > 0 {
		return runtime.Int(1), nil
	}

	return runtime.Int(0), nil
}
