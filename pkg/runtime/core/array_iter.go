package core

import (
	"context"
)

type ArrayIterator struct {
	values *arrayList
	length int
	pos    int
}

func NewArrayIterator(values *arrayList) Iterator {
	return &ArrayIterator{values: values, length: len(values.data), pos: 0}
}

func (iter *ArrayIterator) HasNext(_ context.Context) (bool, error) {
	return iter.length > iter.pos, nil
}

func (iter *ArrayIterator) Next(_ context.Context) (value Value, key Value, err error) {
	iter.pos++

	return iter.values.data[iter.pos-1], NewInt(iter.pos - 1), nil
}
