package internal

import (
	"time"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type StreamValue struct {
	*Box[runtime.Stream]
}

func NewStreamValue(stream runtime.Stream) runtime.Value {
	return &StreamValue{
		Box: &Box[runtime.Stream]{
			Value: stream,
		},
	}
}

func (v *StreamValue) Iterate(timeout runtime.Int) *Iterator {
	if timeout == 0 {
		return NewIterator(runtime.NewIterator(v.Value))
	}

	return NewIterator(runtime.NewIteratorWithTimeout(v.Value, time.Duration(timeout)))
}

func (v *StreamValue) Close() error {
	return v.Value.Close()
}
