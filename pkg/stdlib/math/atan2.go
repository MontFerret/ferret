package math

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"math"
)

/*
 * Returns the arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.
 * @param number1 (Int|Float) - Input number.
 * @param number2 (Int|Float) - Input number.
 * @returns (Float) - The arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.
 */
func Atan2(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.IntType, core.FloatType)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], core.IntType, core.FloatType)

	if err != nil {
		return values.None, err
	}

	var arg1 float64
	var arg2 float64

	if args[0].Type() == core.IntType {
		arg1 = float64(args[0].(values.Int))
	} else {
		arg1 = float64(args[0].(values.Float))
	}

	if args[1].Type() == core.IntType {
		arg2 = float64(args[1].(values.Int))
	} else {
		arg2 = float64(args[1].(values.Float))
	}

	return values.NewFloat(math.Atan2(arg1, arg2)), nil
}
