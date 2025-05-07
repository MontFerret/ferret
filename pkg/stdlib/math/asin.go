package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// ASIN returns the arcsine, in radians, of a given number.
// @param {Int | Float} number - Input number.
// @return {Float} - The arcsine, in radians, of a given number.
func Asin(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Asin(toFloat(args[0]))), nil
}
