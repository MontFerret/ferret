package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// RADIANS returns the angle converted from degrees to radians.
// @param {Int | Float} number - The input number.
// @return {Float} - The angle in radians.
func Radians(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return runtime.None, err
	}

	r := toFloat(args[0])

	return runtime.NewFloat(r * DegToRad), nil
}
