package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// LOG2 returns the binary logarithm of a given value.
// @param {Int | Float} number - Input number.
// @return {Float} - The binary logarithm of a given value.
func Log2(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Log2(toFloat(args[0]))), nil
}
