package values

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type ArrayIterator struct {
	values *Array
	length int
	pos    int
}

func NewArrayIterator(values *Array) core.Iterator {
	return &ArrayIterator{values: values, length: int(values.Length()), pos: 0}
}

func (iterator *ArrayIterator) HasNext(_ context.Context) (bool, error) {
	return iterator.length > iterator.pos, nil
}

func (iterator *ArrayIterator) Next(_ context.Context) (value core.Value, key core.Value, err error) {
	idx := NewInt(iterator.pos)
	val := iterator.values.Get(idx)

	iterator.pos++

	return val, idx, nil
}
