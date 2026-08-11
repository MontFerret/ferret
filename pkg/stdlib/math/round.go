package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ROUND returns the nearest integer, rounding half away from zero.
// @param number {Int | Float} Input number.
// @return {Int} The nearest integer, rounding half away from zero.
func Round(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	return runtime.NewInt(int(math.Round(toFloat(arg)))), nil
}
