package data

import (
	"context"
	"errors"
	"io"
	"strconv"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	GroupedAggregateCollector struct {
		entries        []*groupedAggregateEntry
		grouping       groupIndex[*groupedAggregateEntry]
		singleEntry    *groupedAggregateEntry
		singleKey      runtime.Value
		plan           bytecode.AggregatePlan
		hasSingleGroup bool
		sorted         bool
	}

	groupedAggregateEntry struct {
		key    runtime.Value
		group  runtime.List
		states []aggregateState
	}
)

func NewGroupedAggregateCollector(plan bytecode.AggregatePlan) Transformer {
	return &GroupedAggregateCollector{
		plan:    plan,
		entries: make([]*groupedAggregateEntry, 0, 8),
	}
}

func (c *GroupedAggregateCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
	if !c.sorted {
		if err := c.sort(ctx); err != nil {
			return nil, err
		}

		c.sorted = true
	}

	return newGroupedAggregateIterator(c.entries, c.plan.TrackGroupValues), nil
}

func (c *GroupedAggregateCollector) Set(ctx context.Context, key, value runtime.Value) error {
	if groupKey, idx, ok, err := c.aggregateKey(ctx, key); err != nil {
		return err
	} else if ok {
		return c.UpdateAggregate(ctx, groupKey, value, idx)
	}

	entry, err := c.entryFor(ctx, key)
	if err != nil {
		return err
	}

	if !c.plan.TrackGroupValues || entry.group == nil {
		return nil
	}

	return entry.group.Append(ctx, value)
}

func (c *GroupedAggregateCollector) UpdateAggregate(ctx context.Context, groupKey, value runtime.Value, idx int) error {
	if err := validateAggregateSelectorIndex(idx, len(c.plan.Keys)); err != nil {
		return err
	}

	entry, err := c.entryFor(ctx, groupKey)
	if err != nil {
		return err
	}

	updateAggregateState(&entry.states[idx], c.plan.Kinds[idx], value)

	return nil
}

func (c *GroupedAggregateCollector) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
	if groupKey, idx, ok, err := c.aggregateKey(ctx, key); err != nil {
		return nil, err
	} else if ok {
		entry, ok, err := c.lookupEntry(ctx, groupKey)
		if err != nil {
			return nil, err
		}

		if !ok {
			return runtime.None, collectorKeyNotFoundValue(ctx, groupKey)
		}

		if err := validateAggregateSelectorIndex(idx, len(c.plan.Keys)); err != nil {
			return nil, err
		}

		return aggregateValueFor(entry.states[idx], c.plan.Kinds[idx]), nil
	}

	entry, ok, err := c.lookupEntry(ctx, key)
	if err != nil {
		return nil, err
	}

	if !ok {
		return runtime.None, collectorKeyNotFoundValue(ctx, key)
	}

	if !c.plan.TrackGroupValues || entry.group == nil {
		return runtime.None, nil
	}

	return entry.group, nil
}

func (c *GroupedAggregateCollector) Length(_ context.Context) (runtime.Int, error) {
	return runtime.Int(len(c.entries)), nil
}

func (c *GroupedAggregateCollector) MarshalJSON() ([]byte, error) {
	obj := runtime.NewObject()

	addEntry := func(entry *groupedAggregateEntry) {
		keyStr, err := collectorKeyString(context.Background(), entry.key)
		if err != nil {
			keyStr = entry.key.String()
		}

		for idx := 0; idx < len(c.plan.Keys); idx++ {
			aggKey := runtime.NewString(keyStr + runtime.NamespaceSeparator + strconv.Itoa(idx))
			_ = obj.Set(context.Background(), aggKey, aggregateValueFor(entry.states[idx], c.plan.Kinds[idx]))
		}
	}

	for _, entry := range c.entries {
		addEntry(entry)
	}

	return json.Default.Encode(obj)
}

func (c *GroupedAggregateCollector) String() string {
	encoded, err := c.MarshalJSON()

	if err != nil {
		return "[GroupedAggregateCollector]"
	}

	return string(encoded)
}

func (c *GroupedAggregateCollector) Hash() uint64 {
	return 0
}

func (c *GroupedAggregateCollector) Copy() runtime.Value {
	return c
}

func (c *GroupedAggregateCollector) Close() error {
	entries := c.entries
	c.entries = nil
	c.grouping = groupIndex[*groupedAggregateEntry]{}
	c.hasSingleGroup = false
	c.singleKey = nil
	c.singleEntry = nil
	c.sorted = false

	if !c.plan.TrackGroupValues {
		return nil
	}

	errs := make([]error, 0, len(entries))
	for _, entry := range entries {
		if entry == nil || entry.group == nil {
			continue
		}

		if closer, ok := entry.group.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}

func (c *GroupedAggregateCollector) sort(ctx context.Context) error {
	order := make([]runtime.Value, len(c.entries))
	for idx := range c.entries {
		order[idx] = runtime.Int(idx)
	}

	err := runtime.SortSliceWith(ctx, order, func(ctx context.Context, first, second runtime.Value) (runtime.Ordering, error) {
		firstIdx := int(first.(runtime.Int))
		secondIdx := int(second.(runtime.Int))

		return runtime.CompareValues(ctx, c.entries[firstIdx].key, c.entries[secondIdx].key)
	})
	if err != nil {
		return err
	}

	sorted := make([]*groupedAggregateEntry, len(c.entries))

	for idx, value := range order {
		sorted[idx] = c.entries[int(value.(runtime.Int))]
	}

	c.entries = sorted

	return nil
}

func (c *GroupedAggregateCollector) entryFor(ctx context.Context, key runtime.Value) (*groupedAggregateEntry, error) {
	// Fast path: first key stays in singleKey/singleEntry to avoid map allocation.
	if c.grouping.len() == 0 && !c.hasSingleGroup {
		entry := c.newEntry(key)
		c.singleKey = key
		c.singleEntry = entry
		c.hasSingleGroup = true
		c.entries = append(c.entries, entry)
		c.sorted = false

		return entry, nil
	}

	// Promote to map when a second distinct key appears.
	if c.grouping.len() == 0 {
		if c.hasSingleGroup {
			if key.Hash() == c.singleKey.Hash() {
				equal, err := runtime.EqualValues(ctx, key, c.singleKey)
				if err != nil {
					return nil, err
				}

				if equal {
					return c.singleEntry, nil
				}
			}

			if err := c.grouping.insertUnique(ctx, c.singleKey, c.singleEntry); err != nil {
				return nil, err
			}
		}

		c.hasSingleGroup = false
		c.singleKey = nil
		c.singleEntry = nil

		entry := c.newEntry(key)
		if err := c.grouping.insertUnique(ctx, key, entry); err != nil {
			return nil, err
		}

		c.entries = append(c.entries, entry)
		c.sorted = false

		return entry, nil
	}

	entry, loaded, err := c.grouping.loadOrCreate(ctx, key, func() *groupedAggregateEntry {
		return c.newEntry(key)
	})
	if err != nil {
		return nil, err
	}

	if loaded {
		return entry, nil
	}

	c.entries = append(c.entries, entry)
	c.sorted = false

	return entry, nil
}

func (c *GroupedAggregateCollector) lookupEntry(ctx context.Context, key runtime.Value) (*groupedAggregateEntry, bool, error) {
	if c.grouping.len() == 0 {
		if c.hasSingleGroup && c.singleEntry != nil && key.Hash() == c.singleKey.Hash() {
			equal, err := runtime.EqualValues(ctx, c.singleKey, key)
			if err != nil {
				return nil, false, err
			}

			if equal {
				return c.singleEntry, true, nil
			}
		}

		return nil, false, nil
	}

	return c.grouping.get(ctx, key)
}

func (c *GroupedAggregateCollector) newEntry(key runtime.Value) *groupedAggregateEntry {
	var group runtime.List
	if c.plan.TrackGroupValues {
		group = runtime.NewArray(4)
	}

	return &groupedAggregateEntry{
		key:    key,
		group:  group,
		states: make([]aggregateState, len(c.plan.Keys)),
	}
}

func (c *GroupedAggregateCollector) aggregateKey(_ context.Context, key runtime.Value) (runtime.Value, int, bool, error) {
	groupKey, idx, ok := DecodeAggregateKey(key)

	return groupKey, idx, ok, nil
}
