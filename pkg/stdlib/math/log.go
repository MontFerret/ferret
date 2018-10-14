package math

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"math"
)

/*
 * Returns the natural logarithm of a given value.
 * @param number (Int|Float) - Input number.
 * @returns (Float) - The natural logarithm of a given value.
 */
func Log(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.IntType, core.FloatType)

	if err != nil {
		return values.None, err
	}

	return values.NewFloat(math.Log(toFloat(args[0]))), nil
}
