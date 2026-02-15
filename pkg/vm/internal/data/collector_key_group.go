package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type KeyGroupCollector struct {
	*runtime.Box[runtime.List]
	grouping map[string]runtime.List
	// Fast path for the common single-key case: keep first group without a map.
	singleKey   string
	singleGroup runtime.List
	sorted      bool
}

func NewKeyGroupCollector() Transformer {
	return &KeyGroupCollector{
		Box: &runtime.Box[runtime.List]{
			Value: runtime.NewArray(8),
		},
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

func (c *KeyGroupCollector) Set(ctx context.Context, key, value runtime.Value) error {
	var keyStr string

	switch k := key.(type) {
	case runtime.String:
		keyStr = k.String()
	default:
		var err error
		keyStr, err = Stringify(ctx, key)
		if err != nil {
			return err
		}
	}

	// Fast path: first key stays in singleKey/singleGroup to avoid map allocation.
	if c.grouping == nil && c.singleKey == "" {
		group := runtime.NewArray(4)
		c.singleKey = keyStr
		c.singleGroup = group

		if err := c.Value.Append(ctx, NewKV(key, group)); err != nil {
			return err
		}

		return group.Append(ctx, value)
	}

	// Promote to map when a second distinct key appears.
	if c.grouping == nil {
		if keyStr == c.singleKey {
			return c.singleGroup.Append(ctx, value)
		}

		c.grouping = map[string]runtime.List{
			c.singleKey: c.singleGroup,
		}
		c.singleKey = ""
		c.singleGroup = nil
	}

	group, exists := c.grouping[keyStr]

	if !exists {
		group = runtime.NewArray(4)
		c.grouping[keyStr] = group

		if err := c.Value.Append(ctx, NewKV(key, group)); err != nil {
			return err
		}
	}

	return group.Append(ctx, value)
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
	var keyStr string

	switch k := key.(type) {
	case runtime.String:
		keyStr = k.String()
	default:
		var err error
		keyStr, err = Stringify(ctx, key)
		if err != nil {
			return nil, err
		}
	}

	if c.grouping == nil {
		if keyStr == c.singleKey {
			return c.singleGroup, nil
		}

		return runtime.None, runtime.Errorf(runtime.ErrNotFound, "collector key: %s", keyStr)
	}

	v, ok := c.grouping[keyStr]

	if !ok {
		return runtime.None, runtime.Errorf(runtime.ErrNotFound, "collector key: %s", keyStr)
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
	c.singleKey = ""
	c.singleGroup = nil

	if closer := val.(io.Closer); closer != nil {
		return closer.Close()
	}

	return nil
}
