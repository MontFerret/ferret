package math

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// stddev_population returns the square root of the population variance.
// @param numbers {Any[]} A list whose Int and Float elements are used; other elements are ignored. It is not mutated.
// @return {Float} The population standard deviation, or NaN when no numbers remain.
// @throws {TypeError} An argument has an invalid type.
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
