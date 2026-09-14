package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// variance_population returns the population variance, dividing squared deviations from the mean by the numeric count N.
// @param numbers {Any[]} A list whose Int and Float elements are used; other elements are ignored. It is not mutated.
// @return {Float} The population variance, or NaN when no numbers remain.
// @throws {TypeError} An argument has an invalid type.
// @deprecated Use math::variance instead.
func PopulationVariance(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	return variance(ctx, arg.(runtime.List), false, filterNonNumbers)
}

// variance returns the population variance using N.
// @param values {Int[] | Float[]} Numeric list; non-numeric elements are rejected. The list is not mutated.
// @return {Float} The statistic, or NaN for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
func canonicalVariance(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	value, err := variance(ctx, arg.(runtime.List), false, strictNumbers)
	if err != nil {
		return runtime.None, err
	}

	return value, nil
}
