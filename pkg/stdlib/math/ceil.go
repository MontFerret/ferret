package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// CEIL returns the least integer value greater than or equal to a given value.
// @param {Int | Float} number - Input number.
// @return {Int} - The least integer value greater than or equal to a given value.
func Ceil(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return runtime.None, err
	}

	return runtime.NewInt(int(math.Ceil(toFloat(args[0])))), nil
}
