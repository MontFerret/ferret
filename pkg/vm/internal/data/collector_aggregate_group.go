package data

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type GroupedAggregateCollector struct {
	plan *runtime.AggregatePlan
	// Map key is "<group>::<selectorIdx>" to avoid allocating composite list keys per row.
	entries map[string]*groupedAggregateEntry
}

type groupedAggregateEntry struct {
	kind  runtime.AggregateKind
	state aggregateState
}

func NewGroupedAggregateCollector(plan *runtime.AggregatePlan) Transformer {
	if plan == nil {
		panic("aggregate plan is nil")
	}

	return &GroupedAggregateCollector{
		plan:    plan,
		entries: make(map[string]*groupedAggregateEntry),
	}
}

func (c *GroupedAggregateCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
	keys := make([]string, 0, len(c.entries))
	for key := range c.entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	out := runtime.NewArray(len(keys))
	for _, key := range keys {
		entry := c.entries[key]
		if err := out.Append(ctx, NewKV(runtime.NewString(key), c.valueFor(entry.state, entry.kind))); err != nil {
			return nil, err
		}
	}

	iter, err := out.Iterate(ctx)
	if err != nil {
		return nil, err
	}

	return NewKVIterator(iter), nil
}

func (c *GroupedAggregateCollector) Set(ctx context.Context, key, value runtime.Value) error {
	groupKey, err := c.keyString(ctx, key)
	if err != nil {
		return err
	}

	entry, ok := c.entries[groupKey]
	if !ok {
		idx, err := c.parseSelectorIndex(groupKey)
		if err != nil {
			return err
		}
		if idx < 0 || idx >= c.plan.Size() {
			return runtime.Errorf(runtime.ErrInvalidArgument, "aggregate selector index out of range")
		}
		entry = &groupedAggregateEntry{
			kind: c.plan.KindAt(idx),
		}
		c.entries[groupKey] = entry
	}

	c.update(&entry.state, entry.kind, value)

	return nil
}

func (c *GroupedAggregateCollector) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
	groupKey, err := c.keyString(ctx, key)
	if err != nil {
		return nil, err
	}

	entry, ok := c.entries[groupKey]
	if !ok {
		return runtime.None, runtime.Errorf(runtime.ErrNotFound, "collector key: %s", groupKey)
	}

	return c.valueFor(entry.state, entry.kind), nil
}

func (c *GroupedAggregateCollector) Length(_ context.Context) (runtime.Int, error) {
	return runtime.Int(len(c.entries)), nil
}

func (c *GroupedAggregateCollector) MarshalJSON() ([]byte, error) {
	obj := runtime.NewObject()

	for key, entry := range c.entries {
		_ = obj.Set(context.Background(), runtime.NewString(key), c.valueFor(entry.state, entry.kind))
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
	c.plan = nil
	c.entries = nil
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

func (c *GroupedAggregateCollector) parseSelectorIndex(keyStr string) (int, error) {
	// Expect "<group>::<selectorIdx>" encoding; split on the last separator.
	idx := strings.LastIndex(keyStr, runtime.NamespaceSeparator)
	if idx < 0 {
		return 0, runtime.Errorf(runtime.ErrInvalidArgument, "aggregate key must contain selector index")
	}

	selectorStr := keyStr[idx+len(runtime.NamespaceSeparator):]
	selectorIdx, err := strconv.Atoi(selectorStr)
	if err != nil {
		return 0, runtime.Errorf(runtime.ErrInvalidArgument, "aggregate selector index invalid")
	}

	return selectorIdx, nil
}
