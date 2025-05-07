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
func Atan2(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 2); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[1]); err != nil {
		return runtime.None, err
	}

	arg1 := toFloat(args[0])
	arg2 := toFloat(args[1])

	return runtime.NewFloat(math.Atan2(arg1, arg2)), nil
}
