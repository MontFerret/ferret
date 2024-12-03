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
	return &ArrayIterator{values: values, length: values.Length(), pos: 0}
}

func (iter *ArrayIterator) HasNext(_ context.Context) (bool, error) {
	return iter.length > iter.pos, nil
}

func (iter *ArrayIterator) Next(_ context.Context) (value core.Value, key core.Value, err error) {
	iter.pos++

	return iter.values.data[iter.pos-1], NewInt(iter.pos - 1), nil
}
