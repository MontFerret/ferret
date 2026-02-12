package math

import (
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// ABS returns the absolute value of a given number.
// @param {Int | Float} number - Input number.
// @return {Float} - The absolute value of a given number.
func Abs(_ runtime.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertNumber(arg); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Abs(toFloat(arg))), nil
}
