package values

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type RangeIterator struct {
	values *Range
	dir    int64
	pos    int64
}

func NewRangeIterator(values *Range) core.Iterator {
	if values.start > values.end {
		return &RangeIterator{values: values, dir: -1, pos: values.start}
	}

	return &RangeIterator{values: values, dir: 1, pos: values.start}
}

func (iterator *RangeIterator) HasNext(_ context.Context) (bool, error) {
	if iterator.dir == 1 {
		return iterator.values.end > (iterator.pos - 1), nil
	}

	return iterator.values.start > iterator.pos, nil
}

func (iterator *RangeIterator) Next(_ context.Context) (value core.Value, key core.Value, err error) {
	val := NewInt64(iterator.pos)

	iterator.pos += iterator.dir

	return val, val, nil
}
