package internal

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type KeyGroupCollector struct {
	*BaseCollector
	values   runtime.List
	grouping map[string]runtime.List
	sorted   bool
}

func NewKeyGroupCollector() Collector {
	return &KeyGroupCollector{
		BaseCollector: &BaseCollector{},
		values:        runtime.NewArray(8),
		grouping:      make(map[string]runtime.List),
	}
}

func (c *KeyGroupCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
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

func (c *KeyGroupCollector) sort(ctx context.Context) error {
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

func (c *KeyGroupCollector) Collect(ctx context.Context, key, value runtime.Value) error {
	k, err := Stringify(ctx, key)

	if err != nil {
		return err
	}

	group, exists := c.grouping[k]

	if !exists {
		group = runtime.NewArray(4)

		c.grouping[k] = group

		err = c.values.Add(ctx, NewKV(key, group))

		if err != nil {
			return err
		}
	}

	return group.Add(ctx, value)
}
