package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// DEGREES returns the angle converted from radians to degrees.
// @param {Int | Float} number - The input number.
// @return {Float} - The angle in degrees
func Degrees(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertNumber(arg); err != nil {
		return runtime.None, err
	}

	r := toFloat(arg)

	return runtime.NewFloat(r * RadToDeg), nil
}
