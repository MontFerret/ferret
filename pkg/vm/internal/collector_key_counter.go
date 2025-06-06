package internal

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type KeyCounterCollector struct {
	*BaseCollector
	values   runtime.List
	grouping map[string]runtime.Int
	sorted   bool
}

func NewKeyCounterCollector() Collector {
	return &KeyCounterCollector{
		BaseCollector: &BaseCollector{},
		values:        runtime.NewArray(8),
		grouping:      make(map[string]runtime.Int),
	}
}

func (c *KeyCounterCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
	if !c.sorted {
		if err := c.sort(ctx); err != nil {
			return nil, err
		}

		c.sorted = true
	}

	iter, err := c.values.Iterate(ctx)

	if err != nil {
		return nil, err
	}

	return NewKVIterator(iter), nil
}

func (c *KeyCounterCollector) sort(ctx context.Context) error {
	return runtime.SortListWith(ctx, c.values, func(first, second runtime.Value) int64 {
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

func (c *KeyCounterCollector) Collect(ctx context.Context, key, _ runtime.Value) error {
	k, err := Stringify(ctx, key)

	if err != nil {
		return err
	}

	idx, exists := c.grouping[k]

	var kv *KV

	if !exists {
		size, err := c.values.Length(ctx)

		if err != nil {
			return err
		}

		idx = size
		kv = NewKV(key, runtime.ZeroInt)

		if err := c.values.Add(ctx, kv); err != nil {
			return err
		}

		c.grouping[k] = idx
	} else {
		value, err := c.values.Get(ctx, idx)

		if err != nil {
			return err
		}

		kv = value.(*KV)
	}

	if count, ok := kv.Value.(runtime.Int); ok {
		sum := count + 1
		kv.Value = sum
	} else {
		kv.Value = runtime.NewInt(1)
	}

	return nil
}
