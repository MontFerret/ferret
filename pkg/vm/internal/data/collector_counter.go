package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// CounterCollector is a Transformer implementation that tracks and increments a counter of processed values.
type CounterCollector struct {
	*runtime.Box[runtime.Int]
}

func NewCounterCollector() Transformer {
	return &CounterCollector{
		Box: &runtime.Box[runtime.Int]{
			Value: 0,
		},
	}
}

func (c *CounterCollector) Iterate(_ context.Context) (runtime.Iterator, error) {
	return &counterIterator{value: c.Value}, nil
}

func (c *CounterCollector) Set(_ context.Context, _, _ runtime.Value) error {
	c.Increment()

	return nil
}

func (c *CounterCollector) Increment() {
	c.Value++
}

func (c *CounterCollector) Get(_ context.Context, _ runtime.Value) (runtime.Value, error) {
	return c.Value, nil
}

func (c *CounterCollector) Length(_ context.Context) (runtime.Int, error) {
	return 1, nil
}

func (c *CounterCollector) Close() error {
	return nil
}

type counterIterator struct {
	value runtime.Int
	done  bool
}

func (it *counterIterator) Next(_ context.Context) (runtime.Value, runtime.Value, error) {
	if it.done {
		return runtime.None, runtime.None, io.EOF
	}

	it.done = true

	return it.value, runtime.ZeroInt, nil
}
