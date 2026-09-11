package collections

import (
	"context"
	"errors"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Count returns the number of values in an iterable, including map values.
// A measurable source supplies its length without traversal; computing that length
// may still perform I/O. Otherwise one traversal may consume a one-shot source,
// perform I/O, or never finish for an unbounded source. The acquired iterator is
// closed; the source and its values remain borrowed. Traversal, cancellation,
// overflow, and iterator close failures are returned without a partial count.
// @param collection {Iterable} Source whose yielded values are counted.
// @return {Int} Number of yielded values, or the source's measured length.
func Count(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	if err := ctx.Err(); err != nil {
		return runtime.ZeroInt, err
	}

	collection, ok := arg.(runtime.Iterable)
	if !ok {
		return runtime.ZeroInt, runtime.ArgError(runtime.TypeErrorOf(arg, runtime.TypeIterable), 0)
	}

	var (
		count runtime.Int
		err   error
	)

	if measurable, ok := collection.(runtime.Measurable); ok {
		count, err = measurable.Length(ctx)
	} else {
		err = runtime.ForEach(ctx, collection, func(c context.Context, _, _ runtime.Value) (runtime.Boolean, error) {
			if err := c.Err(); err != nil {
				return false, err
			}

			var err error
			count, err = incrementCount(count)

			return err == nil, err
		})
	}

	if canceled := ctx.Err(); canceled != nil && !errors.Is(err, canceled) {
		err = errors.Join(err, canceled)
	}

	if err != nil {
		return runtime.ZeroInt, err
	}

	return count, nil
}

func incrementCount(count runtime.Int) (runtime.Int, error) {
	if count == math.MaxInt64 {
		return runtime.ZeroInt, runtime.Error(runtime.ErrRange, "count exceeds runtime.Int capacity")
	}

	return count + 1, nil
}
