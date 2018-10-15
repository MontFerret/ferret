package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Log2 returns the binary logarithm of a given value.
// @param number (Int|Float) - Input number.
// @returns (Float) - The binary logarithm of a given value.
func Log2(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.IntType, core.FloatType)

	if err != nil {
		return values.None, err
	}

	return values.NewFloat(math.Log2(toFloat(args[0]))), nil
}
