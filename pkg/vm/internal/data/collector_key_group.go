package data

import (
	"context"
	"errors"
	"hash/fnv"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type KeyGroupCollector struct {
	singleGroup    runtime.List
	singleKey      runtime.Value
	grouping       groupIndex[runtime.List]
	entries        []*KV
	hasSingleGroup bool
	sorted         bool
}

func NewKeyGroupCollector() Transformer {
	return &KeyGroupCollector{
		entries: make([]*KV, 0, 8),
	}
}

func (c *KeyGroupCollector) String() string {
	return "[KeyGroupCollector]"
}

func (c *KeyGroupCollector) Hash() uint64 {
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte("vm.key_group_collector"))

	return hasher.Sum64()
}

func (c *KeyGroupCollector) Copy() runtime.Value {
	return c
}

func (c *KeyGroupCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
	if !c.sorted {
		if err := c.sort(ctx); err != nil {
			return nil, err
		}

		c.sorted = true
	}

	return newKVEntriesIterator(c.entries), nil
}

func (c *KeyGroupCollector) Set(ctx context.Context, key, value runtime.Value) error {
	// Fast path: first key stays in singleKey/singleGroup to avoid map allocation.
	if c.grouping.len() == 0 && !c.hasSingleGroup {
		group := runtime.NewArray(4)
		c.singleKey = key
		c.singleGroup = group
		c.hasSingleGroup = true
		c.entries = append(c.entries, NewKV(key, group))
		c.sorted = false

		return group.Append(ctx, value)
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
					return c.singleGroup.Append(ctx, value)
				}
			}

			if err := c.grouping.insertUnique(ctx, c.singleKey, c.singleGroup); err != nil {
				return err
			}
		}

		c.hasSingleGroup = false
		c.singleKey = nil
		c.singleGroup = nil

		group := runtime.NewArray(4)
		if err := c.grouping.insertUnique(ctx, key, group); err != nil {
			return err
		}

		c.entries = append(c.entries, NewKV(key, group))
		c.sorted = false

		return group.Append(ctx, value)
	}

	group, exists, err := c.grouping.loadOrCreate(ctx, key, func() runtime.List {
		return runtime.NewArray(4)
	})
	if err != nil {
		return err
	}

	if !exists {
		c.entries = append(c.entries, NewKV(key, group))
	}

	c.sorted = false

	return group.Append(ctx, value)
}

func (c *KeyGroupCollector) sort(ctx context.Context) error {
	return sortKVEntries(ctx, c.entries)
}

func (c *KeyGroupCollector) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
	if c.grouping.len() == 0 {
		if c.hasSingleGroup && key.Hash() == c.singleKey.Hash() {
			equal, err := runtime.EqualValues(ctx, key, c.singleKey)
			if err != nil {
				return nil, err
			}

			if equal {
				return c.singleGroup, nil
			}
		}

		return runtime.None, collectorKeyNotFoundValue(ctx, key)
	}

	v, ok, err := c.grouping.get(ctx, key)
	if err != nil {
		return nil, err
	}

	if !ok {
		return runtime.None, collectorKeyNotFoundValue(ctx, key)
	}

	return v, nil
}

func (c *KeyGroupCollector) Length(_ context.Context) (runtime.Int, error) {
	return runtime.Int(len(c.entries)), nil
}

func (c *KeyGroupCollector) Close() error {
	entries := c.entries
	c.entries = nil
	c.grouping = groupIndex[runtime.List]{}
	c.hasSingleGroup = false
	c.singleKey = nil
	c.singleGroup = nil
	c.sorted = false

	errs := make([]error, 0, len(entries))
	for _, entry := range entries {
		if entry == nil {
			continue
		}

		if closer, ok := entry.Value.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}
