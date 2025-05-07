package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// TAN returns the tangent of a given number.
// @param {Int | Float} number - A number.
// @return {Float} - The tangent.
func Tan(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Tan(toFloat(args[0]))), nil
}
