package random

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// choice selects one uniformly distributed value using the Session's pseudo-random source.
// Reservoir sampling takes O(n) time and O(1) additional storage.
// @param values {Any[]} Input list, traversed once without copying or modifying it.
// @return {Any} One original value, or None for an empty list. Empty and singleton lists consume no randomness.
// @throws {TypeError} The argument is not a List.
func Choice(ctx context.Context, values runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](values, 0)
	if err != nil {
		return runtime.None, err
	}

	src, err := source(ctx)
	if err != nil {
		return runtime.None, err
	}

	return choose(ctx, list, src.Int64)
}

// Reservoir sampling counts visits rather than callback indices. The first
// value needs no draw; each subsequent value replaces it with probability 1/n.
func choose(ctx context.Context, values runtime.List, draw func(int64, int64) int64) (runtime.Value, error) {
	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	var selected runtime.Value = runtime.None
	var count int64

	err := values.ForEach(ctx, func(c context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		if err := c.Err(); err != nil {
			return false, err
		}

		if count == math.MaxInt64 {
			return false, runtime.Error(runtime.ErrInvalidOperation, "list exceeds supported element count")
		}

		count++

		if count == 1 || draw(0, count-1) == 0 {
			selected = value
		}

		return true, nil
	})
	if err != nil {
		return runtime.None, err
	}

	// Only successful traversal gets a final cancellation check; retain host
	// and iterator cleanup errors even when cancellation happens alongside them.
	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	return selected, nil
}
