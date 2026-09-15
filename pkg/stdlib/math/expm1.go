package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// expm1 computes exp(value)-1 with improved numerical accuracy for values near zero.
// @param value {Int | Float} Input number.
// @return {Float} The exponential minus one; negative infinity gives -1, while NaN, positive infinity, and signed zero are preserved.
// @throws {TypeError} The argument is not numeric.
func Expm1(_ context.Context, value runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(value, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Expm1(toFloat(value))), nil
}
