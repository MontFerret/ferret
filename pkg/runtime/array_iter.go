package runtime

import (
	"context"
	"io"
)

type ArrayIterator struct {
	values *Array
	length int
	pos    int
}

func NewArrayIterator(values *Array) Iterator {
	return &ArrayIterator{values: values, length: len(values.data), pos: 0}
}

func (iter *ArrayIterator) Next(_ context.Context) (value Value, key Value, err error) {
	if iter.pos >= iter.length {
		return None, None, io.EOF
	}

	value = iter.values.data[iter.pos]
	key = NewInt(iter.pos)
	iter.pos++

	return
}
