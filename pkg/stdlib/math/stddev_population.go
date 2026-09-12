package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// stddev_population returns the square root of the population variance.
// @param numbers {Int[] | Float[]} A list containing only Int and Float values; it is not mutated.
// @return {Float} The population standard deviation, or NaN for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
func StandardDeviationPopulation(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	value, err := variance(ctx, arg.(runtime.List), false)
	if err != nil {
		return runtime.None, err
	}

	return runtime.Float(math.Sqrt(float64(value))), nil
}
