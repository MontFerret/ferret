package math

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"math"
)

/*
 * Returns the integer closest but not less than value.
 * @param number (Int|Float) - Input number.
 * @returns (Int) - The integer closest but not less than value.
 */
func Ceil(_ context.Context, args ...core.Value) (core.Value, error) {
	err := core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.IntType, core.FloatType)

	if err != nil {
		return values.None, err
	}

	arg := args[0]
	var fv float64

	if arg.Type() == core.IntType {
		fv = float64(arg.(values.Int))
	} else {
		fv = float64(arg.(values.Float))
	}

	return values.NewInt(int(math.Ceil(fv))), nil
}
