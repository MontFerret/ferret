package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// POW returns the base to the exponent value.
// @param {Int | Float} base - The base value.
// @param {Int | Float} exp - The exponent value.
// @return {Float} - The exponentiated value.
func Pow(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 2, 2); err != nil {
		return core.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return core.None, err
	}

	if err := runtime.AssertNumber(args[1]); err != nil {
		return core.None, err
	}

	return core.NewFloat(math.Pow(toFloat(args[0]), toFloat(args[1]))), nil
}
