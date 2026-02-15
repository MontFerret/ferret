package data

import (
	"context"
	"io"
	"strconv"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type GroupedAggregateCollector struct {
	plan *runtime.AggregatePlan
	*runtime.Box[runtime.List]
	grouping map[string]*groupedAggregateEntry
	// Fast path for the common single-key case: keep first group without a map.
	singleKey   string
	singleEntry *groupedAggregateEntry
	sorted      bool
}

type groupedAggregateEntry struct {
	key    runtime.Value
	group  runtime.List
	states []aggregateState
}

func NewGroupedAggregateCollector(plan *runtime.AggregatePlan) Transformer {
	if plan == nil {
		panic("aggregate plan is nil")
	}

	return &GroupedAggregateCollector{
		plan: plan,
		Box: &runtime.Box[runtime.List]{
			Value: runtime.NewArray(8),
		},
	}
}

func (c *GroupedAggregateCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
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

func (c *GroupedAggregateCollector) Set(ctx context.Context, key, value runtime.Value) error {
	if groupKey, idx, ok, err := c.aggregateKey(ctx, key); err != nil {
		return err
	} else if ok {
		entry, err := c.entryFor(ctx, groupKey)
		if err != nil {
			return err
		}

		if idx < 0 || idx >= c.plan.Size() {
			return runtime.Errorf(runtime.ErrInvalidArgument, "aggregate selector index out of range")
		}

		c.update(&entry.states[idx], c.plan.KindAt(idx), value)
		return nil
	}

	entry, err := c.entryFor(ctx, key)
	if err != nil {
		return err
	}

	return entry.group.Append(ctx, value)
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
			return runtime.None, runtime.Errorf(runtime.ErrNotFound, "collector key: %s", groupKey.String())
		}

		if idx < 0 || idx >= c.plan.Size() {
			return runtime.None, runtime.Errorf(runtime.ErrInvalidArgument, "aggregate selector index out of range")
		}

		return c.valueFor(entry.states[idx], c.plan.KindAt(idx)), nil
	}

	keyStr, err := c.keyString(ctx, key)
	if err != nil {
		return nil, err
	}

	entry, ok := c.lookupEntryByString(keyStr)
	if !ok {
		return runtime.None, runtime.Errorf(runtime.ErrNotFound, "collector key: %s", keyStr)
	}

	return entry.group, nil
}

func (c *GroupedAggregateCollector) Length(ctx context.Context) (runtime.Int, error) {
	return c.Value.Length(ctx)
}

func (c *GroupedAggregateCollector) MarshalJSON() ([]byte, error) {
	obj := runtime.NewObject()

	addEntry := func(keyStr string, entry *groupedAggregateEntry) {
		for idx := 0; idx < c.plan.Size(); idx++ {
			aggKey := runtime.NewString(keyStr + runtime.NamespaceSeparator + strconv.Itoa(idx))
			_ = obj.Set(context.Background(), aggKey, c.valueFor(entry.states[idx], c.plan.KindAt(idx)))
		}
	}

	if c.singleKey != "" && c.singleEntry != nil {
		addEntry(c.singleKey, c.singleEntry)
	}

	for key, entry := range c.grouping {
		addEntry(key, entry)
	}

	return obj.MarshalJSON()
}

func (c *GroupedAggregateCollector) String() string {
	data, err := c.MarshalJSON()
	if err != nil {
		return "[GroupedAggregateCollector]"
	}

	return string(data)
}

func (c *GroupedAggregateCollector) Unwrap() interface{} {
	return c
}

func (c *GroupedAggregateCollector) Hash() uint64 {
	return 0
}

func (c *GroupedAggregateCollector) Copy() runtime.Value {
	return c
}

func (c *GroupedAggregateCollector) Close() error {
	val := c.Value
	c.plan = nil
	c.Value = nil
	c.grouping = nil
	c.singleKey = ""
	c.singleEntry = nil
	c.sorted = false

	if closer, ok := val.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}

func (c *GroupedAggregateCollector) update(state *aggregateState, kind runtime.AggregateKind, value runtime.Value) {
	switch kind {
	case runtime.AggregateCount:
		state.count++
	case runtime.AggregateSum:
		if runtime.IsNumber(value) {
			state.sum += toFloat(value)
		}
	case runtime.AggregateAverage:
		if runtime.IsNumber(value) {
			state.sum += toFloat(value)
			state.count++
		}
	case runtime.AggregateMin:
		if runtime.IsNumber(value) {
			v := toFloat(value)
			if !state.hasNumber || v < state.min {
				state.min = v
			}
			state.hasNumber = true
		}
	case runtime.AggregateMax:
		if runtime.IsNumber(value) {
			v := toFloat(value)
			if !state.hasNumber || v > state.max {
				state.max = v
			}
			state.hasNumber = true
		}
	}
}

func (c *GroupedAggregateCollector) valueFor(state aggregateState, kind runtime.AggregateKind) runtime.Value {
	switch kind {
	case runtime.AggregateCount:
		return state.count
	case runtime.AggregateSum:
		return runtime.NewFloat(state.sum)
	case runtime.AggregateAverage:
		if state.count == 0 {
			return runtime.ZeroFloat
		}
		return runtime.NewFloat(state.sum / float64(state.count))
	case runtime.AggregateMin:
		if !state.hasNumber {
			return runtime.None
		}
		return runtime.NewFloat(state.min)
	case runtime.AggregateMax:
		if !state.hasNumber {
			return runtime.None
		}
		return runtime.NewFloat(state.max)
	default:
		return runtime.None
	}
}

func (c *GroupedAggregateCollector) sort(ctx context.Context) error {
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

func (c *GroupedAggregateCollector) entryFor(ctx context.Context, key runtime.Value) (*groupedAggregateEntry, error) {
	keyStr, err := c.keyString(ctx, key)
	if err != nil {
		return nil, err
	}

	// Fast path: first key stays in singleKey/singleEntry to avoid map allocation.
	if c.grouping == nil && c.singleKey == "" {
		entry := c.newEntry(key)
		c.singleKey = keyStr
		c.singleEntry = entry

		if err := c.Value.Append(ctx, NewKV(key, entry.group)); err != nil {
			return nil, err
		}

		return entry, nil
	}

	// Promote to map when a second distinct key appears.
	if c.grouping == nil {
		if keyStr == c.singleKey {
			return c.singleEntry, nil
		}

		c.grouping = map[string]*groupedAggregateEntry{
			c.singleKey: c.singleEntry,
		}
		c.singleKey = ""
		c.singleEntry = nil
	}

	if entry, ok := c.grouping[keyStr]; ok {
		return entry, nil
	}

	entry := c.newEntry(key)
	c.grouping[keyStr] = entry

	if err := c.Value.Append(ctx, NewKV(key, entry.group)); err != nil {
		return nil, err
	}

	return entry, nil
}

func (c *GroupedAggregateCollector) lookupEntry(ctx context.Context, key runtime.Value) (*groupedAggregateEntry, bool, error) {
	keyStr, err := c.keyString(ctx, key)
	if err != nil {
		return nil, false, err
	}

	entry, ok := c.lookupEntryByString(keyStr)
	return entry, ok, nil
}

func (c *GroupedAggregateCollector) lookupEntryByString(keyStr string) (*groupedAggregateEntry, bool) {
	if c.grouping == nil {
		if c.singleKey == keyStr && c.singleEntry != nil {
			return c.singleEntry, true
		}

		return nil, false
	}

	entry, ok := c.grouping[keyStr]
	return entry, ok
}

func (c *GroupedAggregateCollector) newEntry(key runtime.Value) *groupedAggregateEntry {
	return &groupedAggregateEntry{
		key:    key,
		group:  runtime.NewArray(4),
		states: make([]aggregateState, c.plan.Size()),
	}
}

func (c *GroupedAggregateCollector) keyString(ctx context.Context, key runtime.Value) (string, error) {
	var keyStr string

	switch k := key.(type) {
	case runtime.String:
		keyStr = k.String()
	default:
		var err error
		keyStr, err = Stringify(ctx, key)
		if err != nil {
			return "", err
		}
	}

	return keyStr, nil
}

func (c *GroupedAggregateCollector) aggregateKey(ctx context.Context, key runtime.Value) (runtime.Value, int, bool, error) {
	// Aggregation updates use a tagged array key:
	// [AggregateKeyMarker, <groupKey>, <selectorIdx>].
	list, ok := key.(runtime.List)
	if !ok {
		return nil, 0, false, nil
	}

	length, err := list.Length(ctx)
	if err != nil {
		return nil, 0, false, err
	}

	if length != 3 {
		return nil, 0, false, nil
	}

	marker, err := list.Get(ctx, 0)
	if err != nil {
		return nil, 0, false, err
	}

	if marker != runtime.AggregateKeyMarker {
		return nil, 0, false, nil
	}

	groupKey, err := list.Get(ctx, 1)
	if err != nil {
		return nil, 0, false, err
	}

	idxVal, err := list.Get(ctx, 2)
	if err != nil {
		return nil, 0, false, err
	}

	idx, ok := idxVal.(runtime.Int)
	if !ok {
		return nil, 0, false, runtime.Errorf(runtime.ErrInvalidArgument, "aggregate selector index invalid")
	}

	return groupKey, int(idx), true, nil
}
