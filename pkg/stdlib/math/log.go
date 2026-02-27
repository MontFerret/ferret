package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// LOG returns the natural logarithm of a given value.
// @param {Int | Float} number - Input number.
// @return {Float} - The natural logarithm of a given value.
func Log(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Log(toFloat(arg))), nil
}
