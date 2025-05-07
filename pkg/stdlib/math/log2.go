package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// LOG2 returns the binary logarithm of a given value.
// @param {Int | Float} number - Input number.
// @return {Float} - The binary logarithm of a given value.
func Log2(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return core.None, err
	}

	return core.NewFloat(math.Log2(toFloat(args[0]))), nil
}
