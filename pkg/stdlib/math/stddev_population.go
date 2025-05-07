package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// STDDEV_POPULATION returns the population standard deviation of the values in a given array.
// @param {Int[] | Float[]} numbers - arrayList of numbers.
// @return {Float} - The population standard deviation.
func StandardDeviationPopulation(ctx context.Context, args ...core.Value) (core.Value, error) {
	if err := core.ValidateArgs(args, 1, 1); err != nil {
		return core.None, err
	}

	arr, err := runtime.CastList(args[0])

	if err != nil {
		return core.None, err
	}

	size, err := arr.Length(ctx)

	if err != nil {
		return core.None, err
	}

	if size == 0 {
		return core.NewFloat(math.NaN()), nil
	}

	vp, err := variance(ctx, arr, core.NewInt(0))

	if err != nil {
		return core.None, err
	}

	return core.NewFloat(math.Pow(float64(vp), 0.5)), nil
}
