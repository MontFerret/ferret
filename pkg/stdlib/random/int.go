package random

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// int draws a non-cryptographic pseudo-random integer between inclusive bounds.
// @param min {Int} Inclusive lower bound; must not exceed max. Floats are not coerced.
// @param max {Int} Inclusive upper bound, including the maximum supported Int.
// @return {Int} A uniformly distributed integer in [min, max]. Equal bounds return min without drawing.
// @throws {TypeError} A bound is not an Int.
// @throws {InvalidArgument} The minimum exceeds the maximum.
func Int(ctx context.Context, min, max runtime.Value) (runtime.Value, error) {
	if err := runtime.ValidateArgValue(min, 0, runtime.AssertInt); err != nil {
		return runtime.None, err
	}

	if err := runtime.ValidateArgValue(max, 1, runtime.AssertInt); err != nil {
		return runtime.None, err
	}

	lower, upper := int64(min.(runtime.Int)), int64(max.(runtime.Int))
	if lower > upper {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "minimum must not exceed maximum"), 0)
	}

	src, err := source(ctx)
	if err != nil {
		return runtime.None, err
	}

	return runtime.Int(src.Int64(lower, upper)), nil
}
