package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// ROUND returns the nearest integer, rounding half away from zero.
// @param {Int | Float} number - Input number.
// @return {Int} - The nearest integer, rounding half away from zero.
func Round(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return runtime.None, err
	}

	return runtime.NewInt(int(math.Round(toFloat(args[0])))), nil
}
