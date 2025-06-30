package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// ACOS returns the arccosine, in radians, of a given number.
// @param {Int | Float} number - Input number.
// @return {Float} - The arccosine, in radians, of a given number.
func Acos(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertNumber(arg); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Acos(toFloat(arg))), nil
}
