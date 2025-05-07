package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DEGREES returns the angle converted from radians to degrees.
// @param {Int | Float} number - The input number.
// @return {Float} - The angle in degrees
func Degrees(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return runtime.None, err
	}

	r := toFloat(args[0])

	return runtime.NewFloat(r * RadToDeg), nil
}
