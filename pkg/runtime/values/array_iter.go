package values

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type ArrayIterator struct {
	pos    int
	values *Array
}

func NewArrayIterator(values *Array) core.Iterator {
	return &ArrayIterator{values: values, pos: 0}
}

func (iterator *ArrayIterator) HasNext(_ context.Context) (bool, error) {
	return int(iterator.values.Length()) > iterator.pos, nil
}

func (iterator *ArrayIterator) Next(_ context.Context) (value core.Value, key core.Value, err error) {
	idx := NewInt(iterator.pos)
	val := iterator.values.Get(idx)

	iterator.pos++

	return val, idx, nil
}
