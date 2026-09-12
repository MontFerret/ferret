package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// variance_sample returns the sample variance, dividing squared deviations from the mean by N - 1.
// @param numbers {Int[] | Float[]} A list containing only Int and Float values; it is not mutated.
// @return {Float} The sample variance, or NaN for fewer than two numbers.
// @throws {TypeError} An argument or list element has an invalid type.
func SampleVariance(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	return variance(ctx, arg.(runtime.List), true)
}
