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
	return &ObjectIterator{data: obj.data, keys: nil}
}

func (iterator *ObjectIterator) init() {
	iterator.keys = make([]string, len(iterator.data))

	i := 0

	for key := range iterator.data {
		iterator.keys[i] = key
		i++
	}
}

func (iterator *ObjectIterator) HasNext(_ context.Context) (bool, error) {
	return len(iterator.keys) > iterator.pos, nil
}

func (iterator *ObjectIterator) Next(ctx context.Context) (core.Value, core.Value, error) {
	if iterator.keys == nil {
		iterator.init()
	}

	key := iterator.keys[iterator.pos]
	val := iterator.data[key]

	iterator.pos++

	return val, String(key), nil
}
