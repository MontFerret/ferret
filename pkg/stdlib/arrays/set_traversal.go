package arrays

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/internal/valueset"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// collectSet keeps the first representative of each value in encounter order.
func collectSet(ctx context.Context, list runtime.List) (*valueset.Set, *runtime.Array, error) {
	capacity := nativeArrayCapacity(ctx, list)
	seen := valueset.New(capacity)
	ordered := runtime.NewArray(capacity)
	err := list.ForEach(ctx, func(ctx context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		added, err := seen.Add(ctx, value)
		if err != nil {
			return false, err
		}

		if added {
			_ = ordered.Append(ctx, value)
		}

		return true, nil
	})
	if err != nil {
		return nil, nil, err
	}

	return seen, ordered, nil
}

func filterSet(ctx context.Context, ordered *runtime.Array, set *valueset.Set, keep bool) (runtime.Value, error) {
	return ordered.Filter(ctx, func(ctx context.Context, value runtime.Value, _ runtime.Int) (runtime.Boolean, error) {
		found, err := set.Contains(ctx, value)

		return runtime.Boolean(found == keep), err
	})
}

// Native lengths are cheap allocation hints. Do not introduce potentially
// fallible Length calls on host lists solely to size temporary storage.
func nativeArrayCapacity(ctx context.Context, list runtime.List) int {
	if array, ok := list.(*runtime.Array); ok {
		size, _ := array.Length(ctx)

		return int(size)
	}

	return 0
}
