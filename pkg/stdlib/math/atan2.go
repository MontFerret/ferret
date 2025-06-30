package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// ATAN2 returns the arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.
// @param {Int | Float} number1 - Input number.
// @param {Int | Float} number2 - Input number.
// @return {Float} - The arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.
func Atan2(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertNumber(arg1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(arg2); err != nil {
		return runtime.None, err
	}

	argf1 := toFloat(arg1)
	argf2 := toFloat(arg2)

	return runtime.NewFloat(math.Atan2(argf1, argf2)), nil
}
