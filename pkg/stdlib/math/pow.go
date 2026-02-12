package math

import (
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// POW returns the base to the exponent value.
// @param {Int | Float} base - The base value.
// @param {Int | Float} exp - The exponent value.
// @return {Float} - The exponentiated value.
func Pow(_ runtime.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertNumber(arg1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(arg2); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Pow(toFloat(arg1), toFloat(arg2))), nil
}
