package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// LOG2 returns the binary logarithm of a given value.
// @param {Int | Float} number - Input number.
// @return {Float} - The binary logarithm of a given value.
func Log2(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertNumber(arg); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Log2(toFloat(arg))), nil
}
