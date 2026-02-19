package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type RangeIterator struct {
	values     *Range
	descending bool
	pos        int64
	counter    int64
}

func NewRangeIterator(values *Range) runtime.Iterator {
	if values.start <= values.end {
		return &RangeIterator{values: values, pos: values.start, counter: -1}
	}

	return &RangeIterator{values: values, pos: values.start, counter: -1, descending: true}
}

func (iter *RangeIterator) Next(_ context.Context) (value runtime.Value, key runtime.Value, err error) {
	if !iter.descending && iter.pos > iter.values.end {
		return runtime.None, runtime.None, io.EOF
	}

	if iter.descending && iter.pos < iter.values.end {
		return runtime.None, runtime.None, io.EOF
	}

	iter.counter++

	if !iter.descending {
		iter.pos++
	} else {
		iter.pos--
	}

	if !iter.descending {
		return runtime.Int(iter.pos - 1), runtime.Int(iter.counter), nil
	}

	return runtime.Int(iter.pos + 1), runtime.Int(iter.counter), nil
}
