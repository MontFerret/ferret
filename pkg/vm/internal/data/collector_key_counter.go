package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type KeyCounterCollector struct {
	*runtime.Box[runtime.List]
	grouping map[string]runtime.Int
	// Fast path for the common single-key case: keep first counter without a map.
	singleKey   string
	singleIndex runtime.Int
	singleKV    *KV
	sorted      bool
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

func (c *KeyCounterCollector) Set(ctx context.Context, key, _ runtime.Value) error {
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

	// Fast path: first key stays in singleKey/singleKV to avoid map allocation.
	if c.grouping == nil && c.singleKey == "" {
		kv := NewKV(key, runtime.ZeroInt)
		if err := c.Value.Append(ctx, kv); err != nil {
			return err
		}

		c.singleKey = keyStr
		c.singleIndex = 0
		c.singleKV = kv

		if count, ok := kv.Value.(runtime.Int); ok {
			kv.Value = count + 1
		} else {
			kv.Value = runtime.NewInt(1)
		}

		return nil
	}

	// Promote to map when a second distinct key appears.
	if c.grouping == nil {
		if keyStr == c.singleKey {
			kv := c.singleKV
			if count, ok := kv.Value.(runtime.Int); ok {
				kv.Value = count + 1
			} else {
				kv.Value = runtime.NewInt(1)
			}

			return nil
		}

		c.grouping = map[string]runtime.Int{
			c.singleKey: c.singleIndex,
		}
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
		value, err := c.Value.Get(ctx, idx)

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
			return c.singleIndex, nil
		}

		return runtime.None, runtime.Errorf(runtime.ErrNotFound, "collector key: %s", keyStr)
	}

	v, ok := c.grouping[keyStr]

	if !ok {
		return runtime.None, runtime.Errorf(runtime.ErrNotFound, "collector key: %s", keyStr)
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
	c.singleKey = ""
	c.singleIndex = 0
	c.singleKV = nil

	if closer := val.(io.Closer); closer != nil {
		return closer.Close()
	}

	return nil
}
