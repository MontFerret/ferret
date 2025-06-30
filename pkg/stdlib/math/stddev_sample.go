package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// STDDEV_SAMPLE returns the sample standard deviation of the values in a given array.
// @param {Int[] | Float[]} numbers - arrayList of numbers.
// @return {Float} - The sample standard deviation.
func StandardDeviationSample(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	arr, err := runtime.CastList(arg)

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
