package sdk

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// MapIterator is an implementation of runtime.Iterator for maps.
// It iterates over the keys of the map and returns the corresponding values.
type MapIterator[TKey comparable, TValue any] struct {
	data map[TKey]TValue
	keys []TKey
	pos  int
}

// NewMapIterator creates a new MapIterator for the given map.
func NewMapIterator[TKey comparable, TValue any](data map[TKey]TValue) runtime.Iterator {
	iter := &MapIterator[TKey, TValue]{data: data, keys: make([]TKey, len(data))}

	i := 0

	for key := range iter.data {
		iter.keys[i] = key
		i++
	}

	return iter
}

func (iter *MapIterator[TKey, TValue]) Next(_ context.Context) (runtime.Value, runtime.Value, error) {
	if iter.pos >= len(iter.keys) {
		return runtime.None, runtime.None, io.EOF
	}

	key := iter.keys[iter.pos]
	value := iter.data[key]
	iter.pos++

	var runtimeKey runtime.Value

	if k, ok := any(key).(runtime.Value); ok {
		runtimeKey = k
	} else {
		runtimeKey = NewProxy[TKey](key)
	}

	var runtimeValue runtime.Value

	if v, ok := any(value).(runtime.Value); ok {
		runtimeValue = v
	} else {
		runtimeValue = NewProxy[TValue](value)
	}

	return runtimeKey, runtimeValue, nil
}
