package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type KeyGroupCollector struct {
	*runtime.Box[runtime.List]
	grouping map[string]runtime.List
	sorted   bool
}

func NewKeyGroupCollector() Transformer {
	return &KeyGroupCollector{
		Box: &runtime.Box[runtime.List]{
			Value: runtime.NewArray(8),
		},
		grouping: make(map[string]runtime.List),
	}
}

func (c *KeyGroupCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
	if !c.sorted {
		if err := c.sort(ctx); err != nil {
			return nil, err
		}

		c.sorted = true
	}

	iter, err := c.Value.Iterate(ctx)

	if err != nil {
		return nil, err
	}

	return NewKVIterator(iter), nil
}

func (c *KeyGroupCollector) Add(ctx context.Context, key, value runtime.Value) error {
	k, err := Stringify(ctx, key)

	if err != nil {
		return err
	}

	group, exists := c.grouping[k]

	if !exists {
		group = runtime.NewArray(4)

		c.grouping[k] = group

		err = c.Value.Add(ctx, NewKV(key, group))

		if err != nil {
			return err
		}
	}

	return group.Add(ctx, value)
}

func (c *KeyGroupCollector) sort(ctx context.Context) error {
	return runtime.SortListWith(ctx, c.Value, func(first, second runtime.Value) int64 {
		firstKV, firstOk := first.(*KV)
		secondKV, secondOk := second.(*KV)

		var comp int64

		if firstOk && secondOk {
			comp = runtime.CompareValues(firstKV.Key, secondKV.Key)
		} else {
			comp = runtime.CompareValues(first, second)
		}

		return comp
	})
}

func (c *KeyGroupCollector) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
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

func (c *KeyGroupCollector) Length(ctx context.Context) (runtime.Int, error) {
	return c.Value.Length(ctx)
}

func (c *KeyGroupCollector) Close() error {
	val := c.Value
	c.Value = nil
	c.grouping = nil

	if closer := val.(io.Closer); closer != nil {
		return closer.Close()
	}

	return nil
}
