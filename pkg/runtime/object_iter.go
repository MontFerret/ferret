package runtime

import (
	"context"
	"io"
)

type ObjectIterator struct {
	backing Object
	keys    []string
	pos     int
}

func NewObjectIterator(obj *Object) Iterator {
	return &ObjectIterator{backing: *obj, keys: obj.propertyNames()}
}

func (iter *ObjectIterator) Next(_ context.Context) (Value, Value, error) {
	if iter.pos >= len(iter.keys) {
		return None, None, io.EOF
	}

	key := iter.keys[iter.pos]
	value, _ := iter.backing.lookup(key)
	iter.pos++

	return value, String(key), nil
}
