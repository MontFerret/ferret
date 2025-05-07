package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// STDDEV_SAMPLE returns the sample standard deviation of the values in a given array.
// @param {Int[] | Float[]} numbers - arrayList of numbers.
// @return {Float} - The sample standard deviation.
func StandardDeviationSample(ctx context.Context, args ...runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	arr, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	size, err := arr.Length(ctx)

	if err != nil {
		return runtime.NaN(), err
	}

	if size == 0 {
		return runtime.NaN(), nil
	}

	vp, err := variance(ctx, arr, runtime.NewInt(1))

	if err != nil {
		return runtime.NaN(), err
	}

	return runtime.NewFloat(math.Pow(float64(vp), 0.5)), nil
}
