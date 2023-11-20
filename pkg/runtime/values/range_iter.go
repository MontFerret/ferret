package values

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type RangeIterator struct {
	values *Range
	pos    uint64
}

func NewRangeIterator(values *Range) core.Iterator {
	return &RangeIterator{values: values, pos: values.start}
}

func (iterator *RangeIterator) HasNext(ctx context.Context) (bool, error) {
	return iterator.values.end > iterator.pos, nil
}

func (iterator *RangeIterator) Next(_ context.Context) (value core.Value, key core.Value, err error) {
	val := NewInt64(int64(iterator.pos))

	iterator.pos++

	return val, val, nil
}
