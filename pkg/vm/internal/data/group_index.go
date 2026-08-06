package data

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	groupIndex[T any] struct {
		buckets map[uint64][]groupIndexEntry[T]
		count   int
	}

	groupIndexEntry[T any] struct {
		key   runtime.Value
		value T
	}
)

func (idx *groupIndex[T]) get(ctx context.Context, key runtime.Value) (T, bool, error) {
	var zero T

	if idx == nil || len(idx.buckets) == 0 {
		return zero, false, nil
	}

	bucket := idx.buckets[key.Hash()]
	for _, entry := range bucket {
		equal, err := runtime.EqualValues(ctx, entry.key, key)
		if err != nil {
			return zero, false, err
		}

		if equal {
			return entry.value, true, nil
		}
	}

	return zero, false, nil
}

// insertUnique adds a key the caller has already proved is not equal to any
// key in its hash bucket.
func (idx *groupIndex[T]) insertUnique(_ context.Context, key runtime.Value, value T) error {
	if idx.buckets == nil {
		idx.buckets = make(map[uint64][]groupIndexEntry[T], 8)
	}

	hash := key.Hash()
	idx.buckets[hash] = append(idx.buckets[hash], groupIndexEntry[T]{
		key:   key,
		value: value,
	})
	idx.count++

	return nil
}

func (idx *groupIndex[T]) loadOrStore(ctx context.Context, key runtime.Value, value T) (T, bool, error) {
	return idx.loadOrCreate(ctx, key, func() T { return value })
}

func (idx *groupIndex[T]) loadOrCreate(ctx context.Context, key runtime.Value, create func() T) (T, bool, error) {
	if idx.buckets == nil {
		idx.buckets = make(map[uint64][]groupIndexEntry[T], 8)
	}

	hash := key.Hash()
	bucket := idx.buckets[hash]

	for i := range bucket {
		equal, err := runtime.EqualValues(ctx, bucket[i].key, key)
		if err != nil {
			var zero T
			return zero, false, err
		}

		if equal {
			return bucket[i].value, true, nil
		}
	}

	value := create()
	idx.buckets[hash] = append(bucket, groupIndexEntry[T]{
		key:   key,
		value: value,
	})
	idx.count++

	return value, false, nil
}

func (idx *groupIndex[T]) len() int {
	if idx == nil {
		return 0
	}

	return idx.count
}
