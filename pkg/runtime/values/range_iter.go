package values

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type RangeIterator struct {
	values     *Range
	descending bool
	pos        int64
	counter    int64
}

func NewRangeIterator(values *Range) core.Iterator {
	if values.start <= values.end {
		return &RangeIterator{values: values, pos: values.start, counter: -1}
	}

	return &RangeIterator{values: values, pos: values.start, counter: -1, descending: true}
}

func (iter *RangeIterator) HasNext(_ context.Context) (bool, error) {
	if !iter.descending {
		return iter.values.end >= iter.pos, nil
	}

	return iter.values.end <= iter.pos, nil
}

func (iter *RangeIterator) Next(_ context.Context) error {
	iter.counter++

	if !iter.descending {
		iter.pos++
	} else {
		iter.pos--
	}

	return nil
}

func (iter *RangeIterator) Value() core.Value {
	if !iter.descending {
		return Int(iter.pos - 1)
	}

	return Int(iter.pos + 1)
}

func (iter *RangeIterator) Key() core.Value {
	return Int(iter.counter)
}
