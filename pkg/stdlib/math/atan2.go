package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
)

// Atan2 returns the arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.
// @param number1 (Int|Float) - Input number.
// @param number2 (Int|Float) - Input number.
// @returns (Float) - The arc tangent of y/x, using the signs of the two to determine the quadrant of the return value.
func Atan2(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 2, 2)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], types.Int, types.Float)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[1], types.Int, types.Float)

	if err != nil {
		return values.None, err
	}

	arg1 := toFloat(args[0])
	arg2 := toFloat(args[1])

	return values.NewFloat(math.Atan2(arg1, arg2)), nil
}
