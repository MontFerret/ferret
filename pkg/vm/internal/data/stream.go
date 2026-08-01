package data

import "github.com/MontFerret/ferret/v2/pkg/runtime"

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

	return WrapIterator(newStreamIterator(v.Value, timeout))
}

func (v *StreamValue) Close() error {
	return v.Value.Close()
}
