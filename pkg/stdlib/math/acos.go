package math

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"math"
)

/*
 * Returns the arccosine, in radians, of a given number.
 * @param number (Int|Float) - Input number.
 * @returns (Float) - The arccosine, in radians, of a given number.
 */
func Acos(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.IntType, core.FloatType)

	if err != nil {
		return values.None, err
	}

	arg := args[0]

	if arg.Type() == core.IntType {
		return values.NewFloat(math.Acos(float64(arg.(values.Int)))), nil
	}

	return values.NewFloat(math.Acos(float64(arg.(values.Float)))), nil
}
