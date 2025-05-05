package runtime

import (
	"context"
)

type ArrayIterator struct {
	values *Array
	length int
	pos    int
}

func NewArrayIterator(values *Array) Iterator {
	return &ArrayIterator{values: values, length: len(values.data), pos: 0}
}

func (iter *ArrayIterator) HasNext(_ context.Context) (bool, error) {
	return iter.length > iter.pos, nil
}

func (iter *ArrayIterator) Next(_ context.Context) (value Value, key Value, err error) {
	iter.pos++

	return iter.values.data[iter.pos-1], NewInt(iter.pos - 1), nil
}
