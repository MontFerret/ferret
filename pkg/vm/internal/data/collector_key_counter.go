package data

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type KeyCounterCollector struct {
	singleKey runtime.Value
	*runtime.Box[runtime.List]
	singleKV       *KV
	grouping       groupIndex[*KV]
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
	// Fast path: first key stays in singleKey/singleKV to avoid map allocation.
	if c.grouping.len() == 0 && !c.hasSingleGroup {
		kv := NewKV(key, runtime.ZeroInt)

		if err := c.Value.Append(ctx, kv); err != nil {
			return err
		}

		c.singleKey = key
		c.singleKV = kv
		c.hasSingleGroup = true
		c.sorted = false

		if count, ok := kv.Value.(runtime.Int); ok {
			kv.Value = count + 1
		} else {
			kv.Value = runtime.NewInt(1)
		}

		return nil
	}

	// Promote to map when a second distinct key appears.
	if c.grouping.len() == 0 {
		if c.hasSingleGroup {
			if key.Hash() == c.singleKey.Hash() {
				equal, err := runtime.EqualValues(ctx, key, c.singleKey)
				if err != nil {
					return err
				}

				if equal {
					kv := c.singleKV

					if count, ok := kv.Value.(runtime.Int); ok {
						kv.Value = count + 1
					} else {
						kv.Value = runtime.NewInt(1)
					}

					return nil
				}
			}

			if err := c.grouping.insertUnique(ctx, c.singleKey, c.singleKV); err != nil {
				return err
			}
		}

		c.hasSingleGroup = false
		c.singleKey = nil
		c.singleKV = nil

		kv := NewKV(key, runtime.NewInt(1))
		if err := c.grouping.insertUnique(ctx, key, kv); err != nil {
			return err
		}

		if err := c.Value.Append(ctx, kv); err != nil {
			return err
		}

		c.sorted = false

		return nil
	}

	kv, exists, err := c.grouping.loadOrCreate(ctx, key, func() *KV {
		return NewKV(key, runtime.ZeroInt)
	})
	if err != nil {
		return err
	}

	if !exists {
		if err := c.Value.Append(ctx, kv); err != nil {
			return err
		}

		c.sorted = false
	}

	if count, ok := kv.Value.(runtime.Int); ok {
		kv.Value = count + 1
	} else {
		kv.Value = runtime.NewInt(1)
	}

	return nil
}

func (c *KeyCounterCollector) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
	if c.grouping.len() == 0 {
		if c.hasSingleGroup && key.Hash() == c.singleKey.Hash() {
			equal, err := runtime.EqualValues(ctx, key, c.singleKey)
			if err != nil {
				return nil, err
			}

			if equal {
				return runtime.ZeroInt, nil
			}
		}

		return runtime.None, collectorKeyNotFoundValue(ctx, key)
	}

	kv, ok, err := c.grouping.get(ctx, key)
	if err != nil {
		return nil, err
	}

	if !ok {
		return runtime.None, collectorKeyNotFoundValue(ctx, key)
	}

	return c.indexOf(ctx, kv)
}

func (c *KeyCounterCollector) Length(ctx context.Context) (runtime.Int, error) {
	return c.Value.Length(ctx)
}

func (c *KeyCounterCollector) Close() error {
	val := c.Value
	c.Value = nil
	c.grouping = groupIndex[*KV]{}
	c.hasSingleGroup = false
	c.singleKey = nil
	c.singleKV = nil

	if closer, ok := val.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}

func (c *KeyCounterCollector) indexOf(ctx context.Context, target *KV) (runtime.Value, error) {
	length, err := c.Value.Length(ctx)
	if err != nil {
		return nil, err
	}

	for idx := runtime.ZeroInt; idx < length; idx++ {
		value, err := c.Value.At(ctx, idx)
		if err != nil {
			return nil, err
		}

		if value == target {
			return idx, nil
		}
	}

	return runtime.None, runtime.Error(runtime.ErrUnexpected, "counter entry is missing from its backing list")
}
