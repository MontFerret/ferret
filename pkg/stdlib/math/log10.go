package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// LOG10 returns the decimal logarithm of a given value.
// @param {Int | Float} number - Input number.
// @return {Float} - The decimal logarithm of a given value.
func Log10(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertNumber(arg); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Log10(toFloat(arg))), nil
}
