package values

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type ObjectIterator struct {
	keys []string
	data map[string]core.Value
	pos  int
}

func NewObjectIterator(obj *Object) core.Iterator {
	iter := &ObjectIterator{data: obj.data, keys: make([]string, len(obj.data))}

	i := 0

	for key := range iter.data {
		iter.keys[i] = key
		i++
	}

	return iter
}

func (iter *ObjectIterator) HasNext(_ context.Context) (bool, error) {
	return len(iter.keys) > iter.pos, nil
}

func (iter *ObjectIterator) Next(_ context.Context) (core.Value, core.Value, error) {
	iter.pos++

	value := iter.data[iter.keys[iter.pos-1]]
	key := String(iter.keys[iter.pos-1])

	return value, key, nil
}
