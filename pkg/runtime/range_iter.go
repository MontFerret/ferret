package runtime

import (
	"context"
	"io"
)

type RangeIterator struct {
	values     *Range
	descending bool
	pos        Int
	counter    Int
}

func NewRangeIterator(values *Range) Iterator {
	if values.start <= values.end {
		return &RangeIterator{values: values, pos: values.start, counter: -1}
	}

	return &RangeIterator{values: values, pos: values.start, counter: -1, descending: true}
}

func (iter *RangeIterator) Next(_ context.Context) (value Value, key Value, err error) {
	if !iter.descending && iter.pos > iter.values.end {
		return None, None, io.EOF
	}

	if iter.descending && iter.pos < iter.values.end {
		return None, None, io.EOF
	}

	iter.counter++

	if !iter.descending {
		iter.pos++
	} else {
		iter.pos--
	}

	if !iter.descending {
		return iter.pos - 1, iter.counter, nil
	}

	return iter.pos + 1, iter.counter, nil
}
