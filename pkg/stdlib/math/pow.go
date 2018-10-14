package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Pow returns the base to the exponent value.
// @param base (Int|Float) - The base value.
// @param exp (Int|Float) - The exponent value.
// @returns (Float) - The exponentiated value.
func Pow(_ context.Context, args ...core.Value) (core.Value, error) {
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

	return values.NewFloat(math.Pow(toFloat(args[0]), toFloat(args[1]))), nil
}
