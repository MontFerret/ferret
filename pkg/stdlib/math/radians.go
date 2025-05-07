package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// RADIANS returns the angle converted from degrees to radians.
// @param {Int | Float} number - The input number.
// @return {Float} - The angle in radians.
func Radians(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return core.None, err
	}

	r := toFloat(args[0])

	return core.NewFloat(r * DegToRad), nil
}
