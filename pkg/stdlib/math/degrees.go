package math

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Degrees returns the angle converted from radians to degrees.
// @param number (Float|Int) - The input number.
// @returns (Float) - The angle in degrees.l
func Degrees(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.IntType, core.FloatType)

	if err != nil {
		return values.None, err
	}

	r := toFloat(args[0])

	return values.NewFloat(r * RadToDeg), nil
}
