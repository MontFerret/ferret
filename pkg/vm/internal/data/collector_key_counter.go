package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type KeyCounterCollector struct {
	*runtime.Box[runtime.List]
	grouping       map[string]runtime.Int
	singleKV       *KV
	singleKey      string
	singleIndex    runtime.Int
	hasSingleGroup bool
	sorted         bool
}

func NewKeyCounterCollector() Transformer {
	return &KeyCounterCollector{
		Box: &runtime.Box[runtime.List]{
			Value: runtime.NewArray(8),
		},
	}
}

func (c *KeyCounterCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
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

func (c *KeyCounterCollector) sort(ctx context.Context) error {
	return sortCollectorList(ctx, c.Value)
}

func (c *KeyCounterCollector) Set(ctx context.Context, key, _ runtime.Value) error {
	keyStr, err := normalizeCollectorKey(ctx, key)
	if err != nil {
		return err
	}

	// Fast path: first key stays in singleKey/singleKV to avoid map allocation.
	if c.grouping == nil && !c.hasSingleGroup {
		kv := NewKV(key, runtime.ZeroInt)

		if err := c.Value.Append(ctx, kv); err != nil {
			return err
		}

		c.singleKey = keyStr
		c.singleIndex = 0
		c.singleKV = kv
		c.hasSingleGroup = true

		if count, ok := kv.Value.(runtime.Int); ok {
			kv.Value = count + 1
		} else {
			kv.Value = runtime.NewInt(1)
		}

		return nil
	}

	// Promote to map when a second distinct key appears.
	if c.grouping == nil {
		if c.hasSingleGroup && keyStr == c.singleKey {
			kv := c.singleKV

			if count, ok := kv.Value.(runtime.Int); ok {
				kv.Value = count + 1
			} else {
				kv.Value = runtime.NewInt(1)
			}

			return nil
		}

		c.grouping = map[string]runtime.Int{}

		if c.hasSingleGroup {
			c.grouping[c.singleKey] = c.singleIndex
		}

		c.hasSingleGroup = false
		c.singleKey = ""
		c.singleIndex = 0
		c.singleKV = nil
	}

	idx, exists := c.grouping[keyStr]

	var kv *KV

	if !exists {
		size, err := c.Value.Length(ctx)

		if err != nil {
			return err
		}

		idx = size
		kv = NewKV(key, runtime.ZeroInt)

		if err := c.Value.Append(ctx, kv); err != nil {
			return err
		}

		c.grouping[keyStr] = idx
	} else {
		value, err := c.Value.At(ctx, idx)

		if err != nil {
			return err
		}

		kv = value.(*KV)
	}

	if count, ok := kv.Value.(runtime.Int); ok {
		kv.Value = count + 1
	} else {
		kv.Value = runtime.NewInt(1)
	}

	return nil
}

func (c *KeyCounterCollector) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
	keyStr, err := normalizeCollectorKey(ctx, key)
	if err != nil {
		return nil, err
	}

	if c.grouping == nil {
		if c.hasSingleGroup && keyStr == c.singleKey {
			return c.singleIndex, nil
		}

		return runtime.None, collectorKeyNotFound(keyStr)
	}

	v, ok := c.grouping[keyStr]

	if !ok {
		return runtime.None, collectorKeyNotFound(keyStr)
	}

	return v, nil
}

func (c *KeyCounterCollector) Length(ctx context.Context) (runtime.Int, error) {
	return c.Value.Length(ctx)
}

func (c *KeyCounterCollector) Close() error {
	val := c.Value
	c.Value = nil
	c.grouping = nil
	c.hasSingleGroup = false
	c.singleKey = ""
	c.singleIndex = 0
	c.singleKV = nil

	if closer := val.(io.Closer); closer != nil {
		return closer.Close()
	}

	return nil
}
