package runtime

import (
	"context"
	"io"
)

type ObjectIterator struct {
	data map[string]Value
	keys []string
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

func (iter *ObjectIterator) Next(_ context.Context) (Value, Value, error) {
	if iter.pos >= len(iter.keys) {
		return None, None, io.EOF
	}

	key := iter.keys[iter.pos]
	value := iter.data[key]
	iter.pos++

	return value, String(key), nil
}
