package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// EXP returns Euler's constant (2.71828...) raised to the power of value.
// @param {Int | Float} number - Input number.
// @return {Float} - Euler's constant raised to the power of value.
func Exp(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return core.None, err
	}

	return core.NewFloat(math.Exp(toFloat(args[0]))), nil
}
