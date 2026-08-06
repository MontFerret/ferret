package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// KeyCollector collects semantically unique keys and sorts them in ascending order when iterated.
type KeyCollector struct {
	*runtime.Box[runtime.List]
	grouping groupIndex[runtime.Value]
	sorted   bool
}

func NewKeyCollector() Transformer {
	return &KeyCollector{
		Box: &runtime.Box[runtime.List]{
			Value: runtime.NewArray(8),
		},
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

func (c *KeyCollector) Set(ctx context.Context, key, _ runtime.Value) error {
	_, exists, err := c.grouping.loadOrStore(ctx, key, runtime.None)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	c.sorted = false

	return c.Value.Append(ctx, key)
}

func (c *KeyCollector) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
	v, ok, err := c.grouping.get(ctx, key)
	if err != nil {
		return nil, err
	}

	if !ok {
		return runtime.None, collectorKeyNotFoundValue(ctx, key)
	}

	return v, nil
}

func (c *KeyCollector) Length(ctx context.Context) (runtime.Int, error) {
	return c.Value.Length(ctx)
}

func (c *KeyCollector) Close() error {
	val := c.Value
	c.Value = nil
	c.grouping = groupIndex[runtime.Value]{}

	if closer, ok := val.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
