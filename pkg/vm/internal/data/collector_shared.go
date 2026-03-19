package data

import (
	"context"
	"io"
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type groupIndexEntry[T any] struct {
	key   runtime.Value
	value T
}

type groupIndex[T any] struct {
	buckets map[uint64][]groupIndexEntry[T]
	count   int
}

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

func collectorKeyNotFound(key string) error {
	return runtime.Errorf(runtime.ErrNotFound, "collector key: %s", key)
}

func (idx *groupIndex[T]) get(key runtime.Value) (T, bool) {
	var zero T

	if idx == nil || len(idx.buckets) == 0 {
		return zero, false
	}

	bucket := idx.buckets[key.Hash()]
	for _, entry := range bucket {
		if runtime.CompareValues(entry.key, key) == 0 {
			return entry.value, true
		}
	}

	return zero, false
}

func (idx *groupIndex[T]) set(key runtime.Value, value T) {
	if idx.buckets == nil {
		idx.buckets = make(map[uint64][]groupIndexEntry[T], 8)
	}

	hash := key.Hash()
	bucket := idx.buckets[hash]

	for i := range bucket {
		if runtime.CompareValues(bucket[i].key, key) == 0 {
			bucket[i].value = value
			idx.buckets[hash] = bucket
			return
		}
	}

	idx.buckets[hash] = append(bucket, groupIndexEntry[T]{
		key:   key,
		value: value,
	})
	idx.count++
}

func (idx *groupIndex[T]) len() int {
	if idx == nil {
		return 0
	}

	return idx.count
}

func sortKVEntries(entries []*KV) {
	sort.Slice(entries, func(i, j int) bool {
		return runtime.CompareValues(entries[i].Key, entries[j].Key) < 0
	})
}

type kvEntriesIterator struct {
	entries []*KV
	idx     int
}

func newKVEntriesIterator(entries []*KV) runtime.Iterator {
	return &kvEntriesIterator{entries: entries}
}

func (it *kvEntriesIterator) Next(_ context.Context) (runtime.Value, runtime.Value, error) {
	if it == nil || it.idx >= len(it.entries) {
		return runtime.None, runtime.None, io.EOF
	}

	entry := it.entries[it.idx]
	it.idx++

	return entry.Value, entry.Key, nil
}
