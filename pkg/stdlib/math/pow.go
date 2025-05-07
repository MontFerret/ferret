package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// POW returns the base to the exponent value.
// @param {Int | Float} base - The base value.
// @param {Int | Float} exp - The exponent value.
// @return {Float} - The exponentiated value.
func Pow(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 2, 2); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[1]); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Pow(toFloat(args[0]), toFloat(args[1]))), nil
}
