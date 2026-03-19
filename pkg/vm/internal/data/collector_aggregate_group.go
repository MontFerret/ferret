package data

import (
	"context"
	"errors"
	"io"
	"sort"
	"strconv"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type GroupedAggregateCollector struct {
	entries        []*groupedAggregateEntry
	grouping       groupIndex[*groupedAggregateEntry]
	singleEntry    *groupedAggregateEntry
	singleKey      runtime.Value
	plan           bytecode.AggregatePlan
	hasSingleGroup bool
	sorted         bool
}

type groupedAggregateEntry struct {
	key    runtime.Value
	group  runtime.List
	states []aggregateState
}

func NewGroupedAggregateCollector(plan bytecode.AggregatePlan) Transformer {
	return &GroupedAggregateCollector{
		plan:    plan,
		entries: make([]*groupedAggregateEntry, 0, 8),
	}
}

func (c *GroupedAggregateCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
	if !c.sorted {
		c.sort()
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

	entry := c.entryFor(key)

	if !c.plan.TrackGroupValues || entry.group == nil {
		return nil
	}

	return entry.group.Append(ctx, value)
}

func (c *GroupedAggregateCollector) UpdateAggregate(_ context.Context, groupKey, value runtime.Value, idx int) error {
	if err := validateAggregateSelectorIndex(idx, len(c.plan.Keys)); err != nil {
		return err
	}

	entry := c.entryFor(groupKey)
	updateAggregateState(&entry.states[idx], c.plan.Kinds[idx], value)

	return nil
}

func (c *GroupedAggregateCollector) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
	if groupKey, idx, ok, err := c.aggregateKey(ctx, key); err != nil {
		return nil, err
	} else if ok {
		entry, ok := c.lookupEntry(groupKey)

		if !ok {
			return runtime.None, collectorKeyNotFoundValue(ctx, groupKey)
		}

		if err := validateAggregateSelectorIndex(idx, len(c.plan.Keys)); err != nil {
			return nil, err
		}

		return aggregateValueFor(entry.states[idx], c.plan.Kinds[idx]), nil
	}

	entry, ok := c.lookupEntry(key)
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

func (c *GroupedAggregateCollector) sort() {
	sort.Slice(c.entries, func(i, j int) bool {
		return runtime.CompareValues(c.entries[i].key, c.entries[j].key) < 0
	})
}

func (c *GroupedAggregateCollector) entryFor(key runtime.Value) *groupedAggregateEntry {
	// Fast path: first key stays in singleKey/singleEntry to avoid map allocation.
	if c.grouping.len() == 0 && !c.hasSingleGroup {
		entry := c.newEntry(key)
		c.singleKey = key
		c.singleEntry = entry
		c.hasSingleGroup = true
		c.entries = append(c.entries, entry)
		c.sorted = false

		return entry
	}

	// Promote to map when a second distinct key appears.
	if c.grouping.len() == 0 {
		if c.hasSingleGroup && runtime.CompareValues(key, c.singleKey) == 0 {
			return c.singleEntry
		}

		if c.hasSingleGroup {
			c.grouping.set(c.singleKey, c.singleEntry)
		}
		c.hasSingleGroup = false
		c.singleKey = nil
		c.singleEntry = nil
	}

	if entry, ok := c.grouping.get(key); ok {
		return entry
	}

	entry := c.newEntry(key)
	c.grouping.set(key, entry)
	c.entries = append(c.entries, entry)
	c.sorted = false

	return entry
}

func (c *GroupedAggregateCollector) lookupEntry(key runtime.Value) (*groupedAggregateEntry, bool) {
	if c.grouping.len() == 0 {
		if c.hasSingleGroup && c.singleEntry != nil && runtime.CompareValues(c.singleKey, key) == 0 {
			return c.singleEntry, true
		}

		return nil, false
	}

	return c.grouping.get(key)
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

type groupedAggregateIterator struct {
	entries          []*groupedAggregateEntry
	idx              int
	trackGroupValues bool
}

func newGroupedAggregateIterator(entries []*groupedAggregateEntry, trackGroupValues bool) runtime.Iterator {
	return &groupedAggregateIterator{
		entries:          entries,
		trackGroupValues: trackGroupValues,
	}
}

func (it *groupedAggregateIterator) Next(_ context.Context) (runtime.Value, runtime.Value, error) {
	if it == nil || it.idx >= len(it.entries) {
		return runtime.None, runtime.None, io.EOF
	}

	entry := it.entries[it.idx]
	it.idx++

	if !it.trackGroupValues || entry.group == nil {
		return runtime.None, entry.key, nil
	}

	return entry.group, entry.key, nil
}
