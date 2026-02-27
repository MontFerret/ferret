package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// STDDEV_POPULATION returns the population standard deviation of the values in a given array.
// @param {Int[] | Float[]} numbers - arrayList of numbers.
// @return {Float} - The population standard deviation.
func StandardDeviationPopulation(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgType(arg, 0, runtime.TypeList); err != nil {
		return runtime.None, err
	}

	arr := arg.(runtime.List)
	size, err := arr.Length(ctx)

	if err != nil {
		return runtime.None, err
	}

	if size == 0 {
		return runtime.NewFloat(math.NaN()), nil
	}

	vp, err := variance(ctx, arr, runtime.NewInt(0))

	if err != nil {
		return runtime.None, err
	}

	return runtime.NewFloat(math.Pow(float64(vp), 0.5)), nil
}
