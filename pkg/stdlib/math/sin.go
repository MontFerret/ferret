package math

import (
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// SIN returns the sine of the radian argument.
// @param {Int | Float} number - Input number.
// @return {Float} - The sin, in radians, of a given number.
func Sin(_ runtime.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertNumber(arg); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Sin(toFloat(arg))), nil
}
