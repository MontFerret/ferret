package internal

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime"
)

type ValueCollector struct {
	*runtime.Box[runtime.List]
	sorted bool
}

func NewValueCollector() Transformer {
	return &ValueCollector{
		Box: &runtime.Box[runtime.List]{
			Value: runtime.NewArray(16),
		}}
}

func (c *ValueCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
	if !c.sorted {
		if err := runtime.SortAsc(ctx, c.Value); err != nil {
			return nil, err
		}

		c.sorted = true
	}

	return c.Value.Iterate(ctx)
}

func (c *ValueCollector) Add(ctx context.Context, _, value runtime.Value) error {
	return c.Value.Add(ctx, value)
}
