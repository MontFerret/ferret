package data

import (
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"
)

// KeyCollector is a structure used to collect and group unique runtime values by their string keys.
// It collects only unique keys without values and sorts them in ascending order when iterated.
type KeyCollector struct {
	*runtime.Box[runtime.List]
	alloc    runtime.Allocator
	grouping map[string]runtime.Value
	sorted   bool
}

func NewKeyCollector(alloc runtime.Allocator) Transformer {
	return &KeyCollector{
		Box: &runtime.Box[runtime.List]{
			Value: alloc.Array(8),
		},
		alloc:    alloc,
		grouping: make(map[string]runtime.Value),
	}
}

func (c *KeyCollector) Iterate(ctx runtime.Context) (runtime.Iterator, error) {
	if !c.sorted {
		if err := runtime.SortAsc(ctx, c.Value); err != nil {
			return nil, err
		}

		c.sorted = true
	}

	return c.Value.Iterate(ctx)
}

func (c *KeyCollector) Set(ctx runtime.Context, key, _ runtime.Value) error {
	k, err := Stringify(ctx, key)

	if err != nil {
		return err
	}

	_, exists := c.grouping[k]

	if !exists {
		c.grouping[k] = runtime.None

		return c.Value.Append(ctx, key)
	}

	return nil
}

func (c *KeyCollector) Get(ctx runtime.Context, key runtime.Value) (runtime.Value, error) {
	k, err := Stringify(ctx, key)

	if err != nil {
		return nil, err
	}

	v, ok := c.grouping[k]

	if !ok {
		return runtime.None, runtime.Errorf(runtime.ErrNotFound, "collector key: %s", k)
	}

	return v, nil
}

func (c *KeyCollector) Length(ctx runtime.Context) (runtime.Int, error) {
	return c.Value.Length(ctx)
}

func (c *KeyCollector) Close() error {
	val := c.Value
	c.Value = nil
	c.grouping = nil

	if closer := val.(io.Closer); closer != nil {
		return closer.Close()
	}

	return nil
}
