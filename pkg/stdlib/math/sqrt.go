package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// SQRT returns the square root of a given number.
// @param {Int | Float} value - A number.
// @return {Float} - The square root.
func Sqrt(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertNumber(arg); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Sqrt(toFloat(arg))), nil
}
