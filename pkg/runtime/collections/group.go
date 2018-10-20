package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	GroupKey      func(value core.Value) (core.Value, error)
	GroupIterator struct {
		src    Iterator
		keys   []GroupKey
		ready  bool
		values Iterator
	}
)

func NewGroupIterator(
	src Iterator,
	keys ...GroupKey,
) (*GroupIterator, error) {
	if core.IsNil(src) {
		return nil, core.Error(core.ErrMissedArgument, "source")
	}

	if len(keys) == 0 {
		return nil, core.Error(core.ErrMissedArgument, "key(s)")
	}

	return &GroupIterator{src, keys, false, nil}, nil
}

func (iterator *GroupIterator) HasNext() bool {
	if !iterator.ready {
		iterator.ready = true
		groups, err := iterator.group()

		if err != nil {
			iterator.values = NoopIterator

			return false
		}

		iterator.values = groups
	}

	return iterator.values.HasNext()
}

func (iterator *GroupIterator) Next() (core.Value, core.Value, error) {
	return iterator.values.Next()
}

func (iterator *GroupIterator) group() (Iterator, error) {
	groups := make(map[string]core.Value)

	for iterator.src.HasNext() {
		for _, keyFn := range iterator.keys {
			val, _, err := iterator.src.Next()

			if err != nil {
				return nil, err
			}

			keyVal, err := keyFn(val)

			if err != nil {
				return nil, err
			}

			key := keyVal.String()

			group, exists := groups[key]

			if !exists {
				group = values.NewArray(10)
				groups[key] = group
			}

			group.(*values.Array).Push(val)
		}
	}

	return NewMapIterator(groups), nil
}
