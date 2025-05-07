package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// SIN returns the sine of the radian argument.
// @param {Int | Float} number - Input number.
// @return {Float} - The sin, in radians, of a given number.
func Sin(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Sin(toFloat(args[0]))), nil
}
