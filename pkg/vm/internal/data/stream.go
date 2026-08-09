package data

import (
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	// StreamIterable is the VM-internal iterator factory shared by singular and grouped streams.
	StreamIterable interface {
		runtime.Value
		Iterate(timeout runtime.Duration) IteratorState
	}

	// StreamGroupController exposes completion only to stream-group bytecode.
	StreamGroupController interface {
		runtime.Value
		ArmDone(index int) error
	}
)

type StreamValue struct {
	*runtime.Box[runtime.Stream]
}

func NewStreamValue(stream runtime.Stream) runtime.Value {
	return &StreamValue{
		Box: &runtime.Box[runtime.Stream]{
			Value: stream,
		},
	}
}

func (v *StreamValue) Iterate(timeout runtime.Duration) IteratorState {
	if timeout == 0 {
		return WrapIterator(runtime.NewStreamIterator(v.Value))
	}

	return WrapIterator(runtime.NewStreamIteratorWithTimeout(v.Value, time.Duration(timeout)))
}

func (v *StreamValue) Close() error {
	return v.Value.Close()
}
