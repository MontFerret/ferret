package internal

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

func (iter *RangeIterator) HasNext(ctx context.Context) (bool, error) {
	if !iter.descending {
		return iter.values.end >= iter.pos, nil
	}

	return iter.values.end <= iter.pos, nil
}

func (iter *RangeIterator) Next(ctx context.Context) (value core.Value, key core.Value, err error) {
	iter.counter++

	if !iter.descending {
		iter.pos++
	} else {
		iter.pos--
	}

	if !iter.descending {
		return core.Int(iter.pos - 1), core.Int(iter.counter), nil
	}

	return core.Int(iter.pos + 1), core.Int(iter.counter), nil
}
