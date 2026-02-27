package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// VARIANCE_SAMPLE returns the sample variance of the values in a given array.
// @param {Int[] | Float[]} numbers - arrayList of numbers.
// @return {Float} - The sample variance.
func SampleVariance(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
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

	return variance(ctx, arr, runtime.NewInt(1))
}
