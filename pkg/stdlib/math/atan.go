package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// ATAN returns the arctangent, in radians, of a given number.
// @param {Int | Float} number - Input number.
// @return {Float} - The arctangent, in radians, of a given number.
func Atan(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertNumber(arg); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Atan(toFloat(arg))), nil
}
