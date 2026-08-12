package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// sqrt returns the square root of a given number.
// @param value {Int | Float} A number.
// @return {Float} The square root.
func Sqrt(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Sqrt(toFloat(arg))), nil
}
