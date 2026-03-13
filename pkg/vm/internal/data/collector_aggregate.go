package data

import (
	"context"
	"errors"
	"io"
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/encoding/json"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type AggregateCollector struct {
	plan             bytecode.AggregatePlan
	singleGroupValue runtime.List
	groups           map[string]runtime.List
	singleGroupKey   string
	states           []aggregateState
	hasSingleGroup   bool
	hasData          bool
}

func NewAggregateCollector(plan bytecode.AggregatePlan) Transformer {
	return &AggregateCollector{
		plan:   plan,
		states: make([]aggregateState, len(plan.Keys)),
	}
}

func (c *AggregateCollector) Iterate(ctx context.Context) (runtime.Iterator, error) {
	size := 0

	if c.hasData {
		groupCount := len(c.groups)

		if c.hasSingleGroup {
			groupCount++
		}

		size = len(c.plan.Keys) + groupCount
	}

	values := runtime.NewArray(size)

	if c.hasData {
		for i, key := range c.plan.Keys {
			if err := values.Append(ctx, NewKV(key, c.valueFor(i))); err != nil {
				return nil, err
			}
		}

		if c.hasSingleGroup {
			if err := values.Append(ctx, NewKV(runtime.NewString(c.singleGroupKey), c.singleGroupValue)); err != nil {
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
	keyStr, err := normalizeCollectorKey(ctx, key)
	if err != nil {
		return err
	}

	if idx, ok := c.plan.Index[keyStr]; ok {
		c.hasData = true
		c.update(idx, value)

		return nil
	}

	c.hasData = true

	// Fast path: first non-aggregate key is stored in the single-group slot.
	if !c.hasSingleGroup && len(c.groups) == 0 {
		c.singleGroupKey = keyStr
		c.singleGroupValue = runtime.NewArray(4)
		c.hasSingleGroup = true

		return c.singleGroupValue.Append(ctx, value)
	}

	// If a second distinct key appears, promote to the groups map.
	if c.hasSingleGroup {
		if keyStr == c.singleGroupKey {
			return c.singleGroupValue.Append(ctx, value)
		}

		c.groups = promoteSingleGroup(c.groups, c.singleGroupKey, c.singleGroupValue)

		c.hasSingleGroup = false
		c.singleGroupKey = ""
		c.singleGroupValue = nil
	}

	if c.groups == nil {
		c.groups = make(map[string]runtime.List)
	}

	group, exists := c.groups[keyStr]
	if !exists {
		group = runtime.NewArray(4)
		c.groups[keyStr] = group
	}

	return group.Append(ctx, value)
}

func (c *AggregateCollector) Get(ctx context.Context, key runtime.Value) (runtime.Value, error) {
	keyStr, err := normalizeCollectorKey(ctx, key)
	if err != nil {
		return nil, err
	}

	if idx, ok := c.plan.Index[keyStr]; ok {
		return c.valueFor(idx), nil
	}

	if c.hasSingleGroup && keyStr == c.singleGroupKey {
		return c.singleGroupValue, nil
	}

	if group, ok := c.groups[keyStr]; ok {
		return group, nil
	}

	return runtime.None, collectorKeyNotFound(keyStr)
}

func (c *AggregateCollector) Length(_ context.Context) (runtime.Int, error) {
	if !c.hasData {
		return 0, nil
	}

	groupCount := len(c.groups)

	if c.hasSingleGroup {
		groupCount++
	}

	return runtime.Int(len(c.plan.Keys) + groupCount), nil
}

func (c *AggregateCollector) MarshalJSON() ([]byte, error) {
	obj := runtime.NewObject()

	if c.hasData {
		for i, key := range c.plan.Keys {
			_ = obj.Set(context.Background(), key, c.valueFor(i))
		}

		if c.hasSingleGroup {
			_ = obj.Set(context.Background(), runtime.NewString(c.singleGroupKey), c.singleGroupValue)
		}

		for key, value := range c.groups {
			_ = obj.Set(context.Background(), runtime.NewString(key), value)
		}
	}

	return json.Default.Encode(obj)
}

func (c *AggregateCollector) String() string {
	encoded, err := c.MarshalJSON()

	if err != nil {
		return "[AggregateCollector]"
	}

	return string(encoded)
}

func (c *AggregateCollector) Hash() uint64 {
	return 0
}

func (c *AggregateCollector) Copy() runtime.Value {
	return c
}

func (c *AggregateCollector) Close() error {
	var err error
	if c.hasSingleGroup {
		if closer, ok := c.singleGroupValue.(io.Closer); ok {
			err = closer.Close()
		}
	} else {
		errs := make([]error, 0, len(c.groups))
		for _, group := range c.groups {
			if closer, ok := group.(io.Closer); ok {
				if closeErr := closer.Close(); closeErr != nil {
					errs = append(errs, closeErr)
				}
			}
		}

		err = errors.Join(errs...)
	}

	c.states = nil
	c.hasSingleGroup = false
	c.singleGroupKey = ""
	c.singleGroupValue = nil
	c.groups = nil
	c.hasData = false

	return err
}

func (c *AggregateCollector) update(idx int, value runtime.Value) {
	state := &c.states[idx]
	updateAggregateState(state, c.plan.Kinds[idx], value)
}

func (c *AggregateCollector) valueFor(idx int) runtime.Value {
	return aggregateValueFor(c.states[idx], c.plan.Kinds[idx])
}
