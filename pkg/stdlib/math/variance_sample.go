package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// variance_sample returns the sample variance, dividing squared deviations from the mean by the numeric count N minus one.
// @param numbers {Any[]} A list whose Int and Float elements are used; other elements are ignored. It is not mutated.
// @return {Float} The sample variance, or NaN for fewer than two numbers.
// @throws {TypeError} An argument has an invalid type.
func SampleVariance(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	return variance(ctx, arg.(runtime.List), true)
}
