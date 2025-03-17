package core

import (
	"context"
)

type ObjectIterator struct {
	keys []string
	data map[string]Value
	pos  int
}

func NewObjectIterator(obj *hashMap) Iterator {
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

func (iter *ObjectIterator) Next(_ context.Context) (Value, Value, error) {
	iter.pos++

	value := iter.data[iter.keys[iter.pos-1]]
	key := String(iter.keys[iter.pos-1])

	return value, key, nil
}
