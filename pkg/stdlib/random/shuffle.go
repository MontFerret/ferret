package random

import (
	"context"
	"errors"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// shuffle randomizes the input values using the Session's pseudo-random source.
// Materialization and Fisher-Yates use O(n) append/swap operations; their cost
// and the destination's storage depend on the source's collection backend.
// @param values {List} Input list, traversed once and never modified. Its factory creates the destination; values retain their identity.
// @return {List} A new list preserving the source's implementation family and backend configuration, containing the original values in randomized order. Empty and singleton lists consume no randomness.
// @throws {TypeError} The argument is not a List.
func Shuffle(ctx context.Context, values runtime.Value) (_ runtime.Value, err error) {
	list, err := runtime.CastArg[runtime.List](values, 0)
	if err != nil {
		return runtime.None, err
	}

	src, err := source(ctx)
	if err != nil {
		return runtime.None, err
	}

	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	destination, err := list.New(ctx)
	if err != nil {
		return runtime.None, err
	}

	// A successful factory transfers the independent destination to us until
	// success. The source and yielded values remain borrowed throughout.
	defer func() {
		if err != nil {
			if closer, ok := destination.(io.Closer); ok {
				if closeErr := closer.Close(); closeErr != nil {
					err = errors.Join(err, closeErr)
				}
			}
		}
	}()

	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	err = list.ForEach(ctx, func(c context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		if err := c.Err(); err != nil {
			return false, err
		}

		if err := destination.Append(c, value); err != nil {
			return false, err
		}

		if err := c.Err(); err != nil {
			return false, err
		}

		return true, nil
	})
	if err != nil {
		return runtime.None, err
	}

	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	size, err := destination.Length(ctx)
	if err != nil {
		return runtime.None, err
	}

	if size < 0 {
		return runtime.None, runtime.Error(runtime.ErrInvalidOperation, "negative list length")
	}

	// Fisher-Yates operates only on the destination, after traversal cleanup.
	for i := size - 1; i > 0; i-- {
		if err := ctx.Err(); err != nil {
			return runtime.None, err
		}

		j := runtime.Int(src.Int64(0, int64(i)))
		if err := destination.Swap(ctx, i, j); err != nil {
			return runtime.None, err
		}
	}

	if err := ctx.Err(); err != nil {
		return runtime.None, err
	}

	return destination, nil
}
