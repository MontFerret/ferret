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

func (iter *ArrayIterator) Next(_ context.Context) error {
	iter.pos++

	return nil
}

func (iter *ArrayIterator) Value() core.Value {
	return iter.values.data[iter.pos-1]
}

func (iter *ArrayIterator) Key() core.Value {
	return NewInt(iter.pos - 1)
}
