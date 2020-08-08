package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// LOG10 returns the decimal logarithm of a given value.
// @param {Int | Float} number - Input number.
// @return {Float} - The decimal logarithm of a given value.
func Log10(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Int, types.Float)

	if err != nil {
		return values.None, err
	}

	return values.NewFloat(math.Log10(toFloat(args[0]))), nil
}
