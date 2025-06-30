package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// EXP2 returns 2 raised to the power of value.
// @param {Int | Float} number - Input number.
// @return {Float} - 2 raised to the power of value.
func Exp2(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.AssertNumber(arg); err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Exp2(toFloat(arg))), nil
}
