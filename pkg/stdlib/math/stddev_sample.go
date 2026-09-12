package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// stddev_sample returns the square root of the sample variance.
// @param numbers {Int[] | Float[]} A list containing only Int and Float values; it is not mutated.
// @return {Float} The sample standard deviation, or NaN for fewer than two numbers.
// @throws {TypeError} An argument or list element has an invalid type.
func StandardDeviationSample(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	value, err := variance(ctx, arg.(runtime.List), true)
	if err != nil {
		return runtime.None, err
	}

	return runtime.Float(math.Sqrt(float64(value))), nil
}
