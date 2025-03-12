package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// ASIN returns the arcsine, in radians, of a given number.
// @param {Int | Float} number - Input number.
// @return {Float} - The arcsine, in radians, of a given number.
func Asin(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return core.None, err
	}

	err = core.ValidateType(args[0], types.Int, types.Float)

	if err != nil {
		return core.None, err
	}

	return core.NewFloat(math.Asin(toFloat(args[0]))), nil
}
