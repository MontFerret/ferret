package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// log1p computes ln(1+value) with improved numerical accuracy for values near zero.
// @param value {Int | Float} Input number.
// @return {Float} The logarithm; -1 gives negative infinity, values below -1 give NaN, and signed zero is preserved.
// @throws {TypeError} The argument is not numeric.
func Log1p(_ context.Context, value runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(value, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Log1p(toFloat(value))), nil
}
