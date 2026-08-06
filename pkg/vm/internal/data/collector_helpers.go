package data

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func normalizeCollectorKey(ctx context.Context, key runtime.Value) (string, error) {
	return collectorKeyString(ctx, key)
}

func collectorKeyString(ctx context.Context, key runtime.Value) (string, error) {
	if str, ok := key.(runtime.String); ok {
		return str.String(), nil
	}

	return Stringify(ctx, key)
}

func collectorKeyNotFoundValue(_ context.Context, key runtime.Value) error {
	return runtime.Errorf(runtime.ErrNotFound, "collector key of type %T", key)
}

func sortCollectorList(ctx context.Context, list runtime.List) error {
	return runtime.SortListWith(ctx, list, func(ctx context.Context, first, second runtime.Value) (runtime.Ordering, error) {
		firstKV, firstOK := first.(*KV)
		secondKV, secondOK := second.(*KV)

		if firstOK && secondOK {
			return runtime.CompareValues(ctx, firstKV.Key, secondKV.Key)
		}

		return runtime.CompareValues(ctx, first, second)
	})
}

func promoteSingleGroup[T any](groups map[string]T, singleKey string, singleValue T) map[string]T {
	if groups == nil {
		groups = map[string]T{}
	}

	groups[singleKey] = singleValue

	return groups
}

func sortKVEntries(ctx context.Context, entries []*KV) error {
	values := make([]runtime.Value, len(entries))
	for idx, entry := range entries {
		values[idx] = entry
	}

	err := runtime.SortSliceWith(ctx, values, func(ctx context.Context, first, second runtime.Value) (runtime.Ordering, error) {
		return runtime.CompareValues(ctx, first.(*KV).Key, second.(*KV).Key)
	})
	if err != nil {
		return err
	}

	for idx, value := range values {
		entries[idx] = value.(*KV)
	}

	return nil
}

func collectorKeyNotFound(key string) error {
	return runtime.Errorf(runtime.ErrNotFound, "collector key: %s", key)
}
