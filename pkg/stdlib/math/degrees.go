package math

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"math"
)

const (
	RadToDeg  = 180 / math.Pi
	DegToRad  = math.Pi / 180
	RadToGrad = 200 / math.Pi
	GradToDeg = math.Pi / 200
)

/*
 * Returns the angle converted from radians to degrees.
 * @param number (Float|Int) - The input number.
 * @returns (Float) - The angle in degrees.
 */
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
