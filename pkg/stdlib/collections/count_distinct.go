package collections

import (
	"context"
	"errors"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// CountDistinct counts distinct yielded values using canonical runtime equality.
// Maps contribute values, not keys. One traversal may consume a one-shot source,
// perform I/O, or never finish for an unbounded source. The acquired iterator is
// closed; the source and yielded values are neither cloned nor closed. Iteration,
// equality, cancellation, and close failures are returned without a partial count.
// @param collection {Iterable} Source whose distinct yielded values are counted.
// @return {Int} Number of distinct yielded values.
func CountDistinct(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := ctx.Err(); err != nil {
		return runtime.ZeroInt, err
	}

	collection, ok := arg.(runtime.Iterable)
	if !ok {
		return runtime.ZeroInt, runtime.ArgError(runtime.TypeErrorOf(arg, runtime.TypeIterable), 0)
	}

	seen := valueset.New(0)

	err := runtime.ForEach(ctx, collection, func(c context.Context, value, idx runtime.Value) (runtime.Boolean, error) {
		if err := c.Err(); err != nil {
			return false, err
		}

		if _, err := seen.Add(c, value); err != nil {
			return false, err
		}

		return true, c.Err()
	})

	if canceled := ctx.Err(); canceled != nil && !errors.Is(err, canceled) {
		err = errors.Join(err, canceled)
	}

	if err != nil {
		return runtime.ZeroInt, err
	}

	return runtime.Int(seen.Len()), nil
}
