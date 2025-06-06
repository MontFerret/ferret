package internal

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type CounterCollector struct {
	*BaseCollector

	counter runtime.Int
}

func NewCounterCollector() Collector {
	return &CounterCollector{
		BaseCollector: &BaseCollector{},
		counter:       0,
	}
}

func (c *CounterCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
	return runtime.NewArrayWith(c.counter).Iterate(ctx)
}

func (c *CounterCollector) Collect(ctx context.Context, key, value runtime.Value) error {
	c.counter++

	return nil
}
