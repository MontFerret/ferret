package internal

import (
	"github.com/MontFerret/ferret/pkg/runtime"
	"time"
)

type StreamValue = Box[runtime.Stream]

func NewStreamValue(stream runtime.Stream) runtime.Value {
	return &StreamValue{
		Value: stream,
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
