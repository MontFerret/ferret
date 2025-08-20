package data

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
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

func (c *CounterCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
	return runtime.NewArrayWith(c.Value).Iterate(ctx)
}

func (c *CounterCollector) Add(_ context.Context, _, _ runtime.Value) error {
	c.Value++

	return nil
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
