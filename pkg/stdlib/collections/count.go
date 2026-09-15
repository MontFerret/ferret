package collections

import (
	"context"
	"math"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Count returns the number of values in an iterable, including map values.
// A measurable source supplies its length without traversal; computing that length
// may still perform I/O. Negative measured lengths are rejected without traversal.
// Otherwise one traversal may consume a one-shot source,
// perform I/O, or never finish for an unbounded source. The acquired iterator is
// closed; the source and its values remain borrowed. Traversal, cancellation,
// overflow, and iterator close failures are returned without a partial count.
// @param collection {Iterable} Source whose yielded values are counted.
// @return {Int} Number of yielded values, or the source's measured length.
func Count(ctx context.Context, arg runtime.Value) (runtime.Value, error) {
	var (
		count runtime.Int
		err   error
	)

	switch collection := arg.(type) {
	case measuredIterable:
		count, err = collection.Length(ctx)
		if err == nil && count < 0 {
			err = runtime.Error(runtime.ErrInvalidOperation, "negative iterable length")
		}

		if err != nil {
			return runtime.ZeroInt, err
		}

		return count, nil
	case runtime.Iterable:
		err = runtime.ForEach(ctx, collection, func(_ context.Context, _, _ runtime.Value) (runtime.Boolean, error) {
			var err error
			count, err = incrementCount(count)

			return err == nil, err
		})
	default:
		return runtime.ZeroInt, runtime.ArgError(runtime.TypeErrorOf(arg, runtime.TypeIterable), 0)
	}

	if err != nil {
		return runtime.ZeroInt, err
	}

	return count, nil
}

type measuredIterable interface {
	runtime.Iterable
	runtime.Measurable
}

func incrementCount(count runtime.Int) (runtime.Int, error) {
	if count == math.MaxInt64 {
		return runtime.ZeroInt, runtime.Error(runtime.ErrRange, "count exceeds runtime.Int capacity")
	}

	return count + 1, nil
}
