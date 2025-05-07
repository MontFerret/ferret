package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// ROUND returns the nearest integer, rounding half away from zero.
// @param {Int | Float} number - Input number.
// @return {Int} - The nearest integer, rounding half away from zero.
func Round(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return core.None, err
	}

	return core.NewInt(int(math.Round(toFloat(args[0])))), nil
}
