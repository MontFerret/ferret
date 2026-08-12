package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// ceil returns the least integer value greater than or equal to a given value.
// @param number {Int | Float} Input number.
// @return {Int} The least integer value greater than or equal to a given value.
func Ceil(_ context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertNumber); err != nil {
		return runtime.None, err
	}

	return runtime.NewInt(int(math.Ceil(toFloat(arg)))), nil
}
