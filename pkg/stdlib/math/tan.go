package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// TAN returns the tangent of a given number.
// @param {Int | Float} number - A number.
// @return {Float} - The tangent.
func Tan(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertNumber(arg); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Tan(toFloat(arg))), nil
}
