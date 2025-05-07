package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// EXP returns Euler's constant (2.71828...) raised to the power of value.
// @param {Int | Float} number - Input number.
// @return {Float} - Euler's constant raised to the power of value.
func Exp(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Exp(toFloat(args[0]))), nil
}
