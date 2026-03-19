package data

import "github.com/MontFerret/ferret/v2/pkg/runtime"

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
