package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// COS returns the cosine of a given number.
// @param {Int | Float} number - Input number.
// @return {Float} - The cosine of a given number.
func Cos(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return values.None, err
	}

	if err := core.ValidateType(args[0], types.Int, types.Float); err != nil {
		return values.None, err
	}

	return values.NewFloat(math.Cos(toFloat(args[0]))), nil
}
