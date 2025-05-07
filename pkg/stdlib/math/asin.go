package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// ASIN returns the arcsine, in radians, of a given number.
// @param {Int | Float} number - Input number.
// @return {Float} - The arcsine, in radians, of a given number.
func Asin(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return runtime.None, err
	}

	return core.NewFloat(math.Asin(toFloat(args[0]))), nil
}
