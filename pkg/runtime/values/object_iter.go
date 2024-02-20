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

func (iterator *ObjectIterator) HasNext(_ context.Context) (bool, error) {
	return len(iterator.keys) > iterator.pos, nil
}

func (iterator *ObjectIterator) Next(_ context.Context) (core.Value, core.Value, error) {
	key := iterator.keys[iterator.pos]
	val := iterator.data[key]

	iterator.pos++

	return val, String(key), nil
}
