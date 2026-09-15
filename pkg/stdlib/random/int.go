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
	lower, err := runtime.CastArg[runtime.Int](min, 0)
	if err != nil {
		return runtime.None, err
	}

	upper, err := runtime.CastArg[runtime.Int](max, 1)
	if err != nil {
		return runtime.None, err
	}

	if lower > upper {
		return runtime.None, runtime.ArgError(runtime.Error(runtime.ErrInvalidArgument, "minimum must not exceed maximum"), 0)
	}

	src, err := source(ctx)
	if err != nil {
		return runtime.None, err
	}

	return runtime.Int(src.Int64(int64(lower), int64(upper))), nil
}
