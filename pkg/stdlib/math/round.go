package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// ROUND returns the nearest integer, rounding half away from zero.
// @param {Int | Float} number - Input number.
// @return {Int} - The nearest integer, rounding half away from zero.
func Round(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertNumber(arg); err != nil {
		return runtime.None, err
	}

	return runtime.NewInt(int(math.Round(toFloat(arg)))), nil
}
