package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// Floor returns the greatest integer value less than or equal to a given value.
// @param number (Int|Float) - Input number.
// @returns (Int) - The greatest integer value less than or equal to a given value.
func Floor(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.IntType, core.FloatType)

	if err != nil {
		return values.None, err
	}

	return values.NewInt(int(math.Floor(toFloat(args[0])))), nil
}
