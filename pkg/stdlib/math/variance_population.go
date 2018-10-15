package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// PopulationVariance returns the population variance of the values in a given array.
// @params (Array) - Array of numbers.
// @returns (Float) - The population variance.
func PopulationVariance(_ context.Context, args ...core.Value) (core.Value, error) {
	var err error
	err = core.ValidateArgs(args, 1, 1)

	if err != nil {
		return values.None, err
	}

	err = core.ValidateType(args[0], core.ArrayType)

	if err != nil {
		return values.None, err
	}

	arr := args[0].(*values.Array)

	if arr.Length() == 0 {
		return values.NewFloat(math.NaN()), nil
	}

	return variance(arr, values.NewInt(0))
}
