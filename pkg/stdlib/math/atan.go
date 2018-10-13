package math

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"math"
)

/*
 * Returns the arctangent, in radians, of a given number.
 * @param number (Int|Float) - Input number.
 * @returns (Float) - The arctangent, in radians, of a given number.
 */
func Atan(_ context.Context, args ...core.Value) (core.Value, error) {
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
		return values.NewFloat(math.Atan(float64(arg.(values.Int)))), nil
	}

	return values.NewFloat(math.Atan(float64(arg.(values.Float)))), nil
}
