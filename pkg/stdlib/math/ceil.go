package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// CEIL returns the least integer value greater than or equal to a given value.
// @param {Int | Float} number - Input number.
// @return {Int} - The least integer value greater than or equal to a given value.
func Ceil(_ context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	if err := runtime.AssertNumber(args[0]); err != nil {
		return core.None, err
	}

	return core.NewInt(int(math.Ceil(toFloat(args[0])))), nil
}
