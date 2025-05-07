package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// ACOS returns the arccosine, in radians, of a given number.
// @param {Int | Float} number - Input number.
// @return {Float} - The arccosine, in radians, of a given number.
func Acos(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Acos(toFloat(args[0]))), nil
}
