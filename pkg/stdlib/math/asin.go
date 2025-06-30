package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// ASIN returns the arcsine, in radians, of a given number.
// @param {Int | Float} number - Input number.
// @return {Float} - The arcsine, in radians, of a given number.
func Asin(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertNumber(arg); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Asin(toFloat(arg))), nil
}
