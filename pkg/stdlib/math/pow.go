package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// pow returns the base to the exponent value.
// @param base {Int | Float} The base value.
// @param exp {Int | Float} The exponent value.
// @return {Float} The exponentiated value.
func Pow(_ context.Context, arg1, arg2 runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg1, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgValue(arg2, 1, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Pow(toFloat(arg1), toFloat(arg2))), nil
}
