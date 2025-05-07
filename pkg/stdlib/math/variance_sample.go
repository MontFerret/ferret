package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

// VARIANCE_SAMPLE returns the sample variance of the values in a given array.
// @param {Int[] | Float[]} numbers - arrayList of numbers.
// @return {Float} - The sample variance.
func SampleVariance(ctx context.Context, args ...core.Value) (core.Value, error) {
	if err := runtime.ValidateArgs(args, 1, 1); err != nil {
		return runtime.None, err
	}

	arr, err := runtime.CastList(args[0])

	if err != nil {
		return runtime.None, err
	}

	size, err := arr.Length(ctx)

	if err != nil {
		return runtime.None, err
	}

	if size == 0 {
		return runtime.NewFloat(math.NaN()), nil
	}

	return variance(ctx, arr, core.NewInt(1))
}
