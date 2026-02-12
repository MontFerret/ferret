package runtime

import (
	"context"
)

type ObjectIterator struct {
	keys []string
	data map[string]Value
	pos  int
}

func NewObjectIterator(obj *Object) Iterator {
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
	if iter.pos >= len(iter.keys) {
		return None, None, Error(ErrInvalidOperation, "no more elements")
	}

	key := iter.keys[iter.pos]
	value := iter.data[key]
	iter.pos++

	return value, String(key), nil
}
