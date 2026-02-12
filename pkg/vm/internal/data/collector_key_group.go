package data

import (
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type KeyGroupCollector struct {
	*runtime.Box[runtime.List]
	grouping map[string]runtime.List
	sorted   bool
	alloc    runtime.Allocator
}

func NewKeyGroupCollector(alloc runtime.Allocator) Transformer {
	return &KeyGroupCollector{
		Box: &runtime.Box[runtime.List]{
			Value: alloc.Array(8),
		},
		grouping: make(map[string]runtime.List),
	}
}

func (c *KeyGroupCollector) Iterate(ctx runtime.Context) (runtime.Iterator, error) {
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

func (c *KeyGroupCollector) Set(ctx runtime.Context, key, value runtime.Value) error {
	k, err := Stringify(ctx, key)

	if err != nil {
		return err
	}

	group, exists := c.grouping[k]

	if !exists {
		group = c.alloc.Array(4)

		c.grouping[k] = group

		err = c.Value.Append(ctx, NewKV(key, group))

		if err != nil {
			return err
		}
	}

	return group.Append(ctx, value)
}

func (c *KeyGroupCollector) sort(ctx runtime.Context) error {
	return runtime.SortListWith(ctx, c.Value, func(c runtime.Context, first, second runtime.Value) int64 {
		firstKV, firstOk := first.(*KV)
		secondKV, secondOk := second.(*KV)

		var comp int64

		if firstOk && secondOk {
			comp = runtime.CompareValues(c, firstKV.Key, secondKV.Key)
		} else {
			comp = runtime.CompareValues(c, first, second)
		}

		return comp
	})
}

func (c *KeyGroupCollector) Get(ctx runtime.Context, key runtime.Value) (runtime.Value, error) {
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

func (c *KeyGroupCollector) Length(ctx runtime.Context) (runtime.Int, error) {
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
