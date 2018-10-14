package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Acos returns the arccosine, in radians, of a given number.
// @param number (Int|Float) - Input number.
// @returns (Float) - The arccosine, in radians, of a given number.
func Acos(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.IntType, core.FloatType)

	if err != nil {
		return values.None, err
	}

	return values.NewFloat(math.Acos(toFloat(args[0]))), nil
}
