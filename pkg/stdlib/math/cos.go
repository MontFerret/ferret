package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// COS returns the cosine of a given number.
// @param {Int | Float} number - Input number.
// @return {Float} - The cosine of a given number.
func Cos(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertNumber(arg); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Cos(toFloat(arg))), nil
}
