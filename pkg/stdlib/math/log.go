package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// LOG returns the natural logarithm of a given value.
// @param {Int | Float} number - Input number.
// @return {Float} - The natural logarithm of a given value.
func Log(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return core.None, err
	}

	return core.NewFloat(math.Log(toFloat(args[0]))), nil
}
