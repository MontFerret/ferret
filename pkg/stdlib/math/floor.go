package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// FLOOR returns the greatest integer value less than or equal to a given value.
// @param {Int | Float} number - Input number.
// @return {Int} - The greatest integer value less than or equal to a given value.
func Floor(_ context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return runtime.None, err
	}

	return runtime.NewInt(int(math.Floor(toFloat(args[0])))), nil
}
