package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// variance_population returns the population variance, dividing squared deviations from the mean by N.
// @param numbers {Int[] | Float[]} A list containing only Int and Float values; it is not mutated.
// @return {Float} The population variance, or NaN for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
func PopulationVariance(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	return variance(ctx, arg.(runtime.List), false)
}
