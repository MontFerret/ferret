package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

// StandardDeviationPopulation returns the population standard deviation of the values in a given array.
// @params (Array) - Array of numbers.
// @returns (Float) - The population standard deviation.
func StandardDeviationPopulation(_ context.Context, args ...core.Value) (core.Value, error) {
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

	vp, err := variance(arr, values.NewInt(0))

	if err != nil {
		return values.NewFloat(math.NaN()), err
	}

	return values.NewFloat(math.Pow(float64(vp), 0.5)), nil
}
