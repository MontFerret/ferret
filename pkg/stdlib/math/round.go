package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// round returns the nearest integer, rounding half away from zero.
// @param number {Int | Float} Input number.
// @return {Int} The nearest integer, rounding half away from zero.
// @throws {RangeError} The rounded result is not a finite signed 64-bit integer.
func Round(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	if integer, ok := arg.(runtime.Int); ok {
		return integer, nil
	}

	return roundedInt(math.Round(float64(arg.(runtime.Float))))
}
