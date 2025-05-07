package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// LOG10 returns the decimal logarithm of a given value.
// @param {Int | Float} number - Input number.
// @return {Float} - The decimal logarithm of a given value.
func Log10(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return core.None, err
	}

	return core.NewFloat(math.Log10(toFloat(args[0]))), nil
}
