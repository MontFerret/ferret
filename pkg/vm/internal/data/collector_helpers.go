package data

import (
	"context"
	"sort"

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

func collectorKeyNotFoundValue(ctx context.Context, key runtime.Value) error {
	keyStr, err := collectorKeyString(ctx, key)
	if err != nil {
		return err
	}

	return collectorKeyNotFound(keyStr)
}

func sortCollectorList(ctx context.Context, list runtime.List) error {
	return runtime.SortListWith(ctx, list, func(first, second runtime.Value) int {
		firstKV, firstOK := first.(*KV)
		secondKV, secondOK := second.(*KV)

		if firstOK && secondOK {
			return runtime.CompareValues(firstKV.Key, secondKV.Key)
		}

		return runtime.CompareValues(first, second)
	})
}

func promoteSingleGroup[T any](groups map[string]T, singleKey string, singleValue T) map[string]T {
	if groups == nil {
		groups = map[string]T{}
	}

	groups[singleKey] = singleValue

	return groups
}

func sortKVEntries(entries []*KV) {
	sort.Slice(entries, func(i, j int) bool {
		return runtime.CompareValues(entries[i].Key, entries[j].Key) < 0
	})
}

func collectorKeyNotFound(key string) error {
	return runtime.Errorf(runtime.ErrNotFound, "collector key: %s", key)
}
