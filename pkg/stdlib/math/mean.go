package math

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// mean returns the arithmetic mean of a numeric list.
// @param values {Int[] | Float[]} Numeric list; non-numeric elements are rejected. The list is not mutated.
// @return {Float} The arithmetic mean, or NaN for an empty list.
// @throws {TypeError} An argument or list element has an invalid type.
func canonicalMean(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(arg, 0, runtime.AssertList); err != nil {
		return runtime.None, err
	}

	value, count, err := mean(ctx, arg.(runtime.List), strictNumbers)
	if err != nil {
		return runtime.None, err
	}

	if count == 0 {
		return runtime.NaN(), nil
	}

	return value, nil
}

func mean(ctx context.Context, input runtime.List, policy numberPolicy) (runtime.Float, runtime.Int, error) {
	sum, counts, err := sumNumbers(ctx, input, policy)
	if err != nil {
		return runtime.NaN(), 0, err
	}

	if counts.numeric == 0 {
		return runtime.ZeroFloat, 0, nil
	}

	return runtime.Float(sum / float64(counts.numeric)), counts.numeric, nil
}
