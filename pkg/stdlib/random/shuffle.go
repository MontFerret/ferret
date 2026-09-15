package random

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// shuffle randomizes the input values using the Session's pseudo-random source.
// Materialization and Fisher-Yates shuffling take O(n) time and O(n) storage.
// @param values {Any[]} Input list, traversed once and never modified. Values retain their identity.
// @return {Any[]} A new in-memory Array containing the original values in randomized order. The input is intentionally materialized; empty and singleton lists consume no randomness.
// @throws {TypeError} The argument is not a List.
func Shuffle(ctx context.Context, values runtime.Value) (runtime.Value, error) {
	list, err := runtime.CastArg[runtime.List](values, 0)
	if err != nil {
		return runtime.None, err
	}

	src, err := source(ctx)
	if err != nil {
		return runtime.None, err
	}

	// We do not know what the backend of the list is - hence it could be too large to shuffle in memory.
	// Therefore, we use Factory interface to create a new List with the same backend as the input list, and materialize the input into it.
	snapshot, err := list.New(ctx)
	if err != nil {
		return runtime.None, err
	}

	err = list.ForEach(ctx, func(c context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		if err = snapshot.Append(ctx, value); err != nil {
			return false, err
		}

		return true, nil
	})
	if err != nil {
		return runtime.None, err
	}

	size, err := snapshot.Length(ctx)
	if err != nil {
		return runtime.None, err
	}

	// Fisher-Yates operates only on private storage, after traversal cleanup.
	for i := size - 1; i > 0; i-- {
		if err := ctx.Err(); err != nil {
			return runtime.None, err
		}

		j := runtime.Int(src.Int64(0, int64(i)))

		iv, err := snapshot.At(ctx, i)
		if err != nil {
			return runtime.None, err
		}

		jv, err := snapshot.At(ctx, j)
		if err != nil {
			return runtime.None, err
		}

		if err := snapshot.SetAt(ctx, i, jv); err != nil {
			return runtime.None, err
		}

		if err := snapshot.SetAt(ctx, j, iv); err != nil {
			return runtime.None, err
		}
	}

	return snapshot, nil
}
