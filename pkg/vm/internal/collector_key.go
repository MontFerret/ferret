package internal

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type KeyCollector struct {
	*runtime.Box[runtime.List]
	grouping map[string]runtime.Value
	sorted   bool
}

func NewKeyCollector() Transformer {
	return &KeyCollector{
		Box: &runtime.Box[runtime.List]{
			Value: runtime.NewArray(16),
		},
		grouping: make(map[string]runtime.Value),
	}
}

func (c *KeyCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
	if !c.sorted {
		if err := runtime.SortAsc(ctx, c.Value); err != nil {
			return nil, err
		}

		c.sorted = true
	}

	return c.Value.Iterate(ctx)
}

func (c *KeyCollector) Add(ctx context.Context, key, _ runtime.Value) error {
	k, err := Stringify(ctx, key)

	if err != nil {
		return err
	}

	_, exists := c.grouping[k]

	if !exists {
		c.grouping[k] = runtime.None

		return c.Value.Add(ctx, key)
	}

	return nil
}
