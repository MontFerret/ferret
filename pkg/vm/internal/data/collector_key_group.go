package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type KeyGroupCollector struct {
	singleGroup runtime.List
	*runtime.Box[runtime.List]
	grouping       map[string]runtime.List
	singleKey      string
	hasSingleGroup bool
	sorted         bool
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
	keyStr, err := normalizeCollectorKey(ctx, key)
	if err != nil {
		return err
	}

	// Fast path: first key stays in singleKey/singleGroup to avoid map allocation.
	if c.grouping == nil && !c.hasSingleGroup {
		group := runtime.NewArray(4)
		c.singleKey = keyStr
		c.singleGroup = group
		c.hasSingleGroup = true

		if err := c.Value.Append(ctx, NewKV(key, group)); err != nil {
			return err
		}

		return group.Append(ctx, value)
	}

	// Promote to map when a second distinct key appears.
	if c.grouping == nil {
		if c.hasSingleGroup && keyStr == c.singleKey {
			return c.singleGroup.Append(ctx, value)
		}

		if c.hasSingleGroup {
			c.grouping = promoteSingleGroup(c.grouping, c.singleKey, c.singleGroup)
		} else {
			c.grouping = map[string]runtime.List{}
		}

		c.hasSingleGroup = false
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
	return sortCollectorList(ctx, c.Value)
}

func (c *KeyGroupCollector) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
	keyStr, err := normalizeCollectorKey(ctx, key)
	if err != nil {
		return nil, err
	}

	if c.grouping == nil {
		if c.hasSingleGroup && keyStr == c.singleKey {
			return c.singleGroup, nil
		}

		return runtime.None, collectorKeyNotFound(keyStr)
	}

	v, ok := c.grouping[keyStr]

	if !ok {
		return runtime.None, collectorKeyNotFound(keyStr)
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
	c.hasSingleGroup = false
	c.singleKey = ""
	c.singleGroup = nil

	if closer, ok := val.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
