package data

import (
	"github.com/MontFerret/ferret/pkg/runtime"
)

// CounterCollector is a Transformer implementation that tracks and increments a counter of processed values.
type CounterCollector struct {
	*runtime.Box[runtime.Int]
	alloc runtime.Allocator
}

func NewCounterCollector(alloc runtime.Allocator) Transformer {
	return &CounterCollector{
		Box: &runtime.Box[runtime.Int]{
			Value: 0,
		},
		alloc: alloc,
	}
}

func (c *CounterCollector) Iterate(ctx runtime.Context) (runtime.Iterator, error) {
	arr := c.alloc.Array(1)
	_ = arr.Append(ctx, c.Value)

	return arr.Iterate(ctx)
}

func (c *CounterCollector) Set(_ runtime.Context, key, value runtime.Value) error {
	c.Value++

	return nil
}

func (c *CounterCollector) Get(_ runtime.Context, _ runtime.Value) (runtime.Value, error) {
	return c.Value, nil
}

func (c *CounterCollector) Length(_ runtime.Context) (runtime.Int, error) {
	return 1, nil
}

func (c *CounterCollector) Close() error {
	return nil
}
