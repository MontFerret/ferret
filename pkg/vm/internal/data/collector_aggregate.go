package data

import (
	"context"
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type aggregateState struct {
	sum       float64
	count     runtime.Int
	min       float64
	max       float64
	hasNumber bool
}

type AggregateCollector struct {
	plan    *runtime.AggregatePlan
	states  []aggregateState
	groups  map[string]runtime.List
	hasData bool
}

func NewAggregateCollector(plan *runtime.AggregatePlan) Transformer {
	if plan == nil {
		panic("aggregate plan is nil")
	}

	return &AggregateCollector{
		plan:   plan,
		states: make([]aggregateState, plan.Size()),
		groups: make(map[string]runtime.List),
	}
}

func (c *AggregateCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
	size := 0
	if c.hasData {
		size = c.plan.Size() + len(c.groups)
	}

	values := runtime.NewArray(size)

	if c.hasData {
		for i, key := range c.plan.Keys() {
			if err := values.Append(ctx, NewKV(key, c.valueFor(i))); err != nil {
				return nil, err
			}
		}

		if len(c.groups) > 0 {
			keys := make([]string, 0, len(c.groups))
			for key := range c.groups {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			for _, key := range keys {
				if err := values.Append(ctx, NewKV(runtime.NewString(key), c.groups[key])); err != nil {
					return nil, err
				}
			}
		}
	}

	iter, err := values.Iterate(ctx)
	if err != nil {
		return nil, err
	}

	return NewKVIterator(iter), nil
}

func (c *AggregateCollector) Set(ctx context.Context, key, value runtime.Value) error {
	k, err := Stringify(ctx, key)
	if err != nil {
		return err
	}

	if idx, ok := c.plan.Index(k); ok {
		c.hasData = true
		c.update(idx, value)
		return nil
	}

	group, exists := c.groups[k]
	if !exists {
		group = runtime.NewArray(4)
		c.groups[k] = group
	}

	c.hasData = true
	return group.Append(ctx, value)
}

func (c *AggregateCollector) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
	k, err := Stringify(ctx, key)
	if err != nil {
		return nil, err
	}

	if idx, ok := c.plan.Index(k); ok {
		return c.valueFor(idx), nil
	}

	if group, ok := c.groups[k]; ok {
		return group, nil
	}

	return runtime.None, runtime.Errorf(runtime.ErrNotFound, "collector key: %s", k)
}

func (c *AggregateCollector) Length(_ context.Context) (runtime.Int, error) {
	if !c.hasData {
		return 0, nil
	}

	return runtime.Int(c.plan.Size() + len(c.groups)), nil
}

func (c *AggregateCollector) MarshalJSON() ([]byte, error) {
	obj := runtime.NewObject()

	if c.hasData {
		for i, key := range c.plan.Keys() {
			_ = obj.Set(context.Background(), key, c.valueFor(i))
		}

		for key, value := range c.groups {
			_ = obj.Set(context.Background(), runtime.NewString(key), value)
		}
	}

	return obj.MarshalJSON()
}

func (c *AggregateCollector) String() string {
	data, err := c.MarshalJSON()
	if err != nil {
		return "[AggregateCollector]"
	}

	return string(data)
}

func (c *AggregateCollector) Unwrap() interface{} {
	return c
}

func (c *AggregateCollector) Hash() uint64 {
	return 0
}

func (c *AggregateCollector) Copy() runtime.Value {
	return c
}

func (c *AggregateCollector) Close() error {
	c.plan = nil
	c.states = nil
	c.groups = nil
	c.hasData = false

	return nil
}

func (c *AggregateCollector) update(idx int, value runtime.Value) {
	state := &c.states[idx]

	switch c.plan.KindAt(idx) {
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

func (c *AggregateCollector) valueFor(idx int) runtime.Value {
	state := c.states[idx]

	switch c.plan.KindAt(idx) {
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

func toFloat(arg runtime.Value) float64 {
	switch v := arg.(type) {
	case runtime.Float:
		return float64(v)
	case runtime.Int:
		return float64(v)
	default:
		return 0
	}
}
