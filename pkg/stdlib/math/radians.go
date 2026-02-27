package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// RADIANS returns the angle converted from degrees to radians.
// @param {Int | Float} number - The input number.
// @return {Float} - The angle in radians.
func Radians(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	r := toFloat(arg)

	return runtime.NewFloat(r * DegToRad), nil
}
