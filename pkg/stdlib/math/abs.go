package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// ABS returns the absolute value of a given number.
// @param {Int | Float} number - Input number.
// @return {Float} - The absolute value of a given number.
func Abs(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return runtime.None, err
	}

	return core.NewFloat(math.Abs(toFloat(args[0]))), nil
}
