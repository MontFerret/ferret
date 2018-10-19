package collections

import (
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	GroupIterator struct {
		src       Iterator
		selectors []*GroupSelector
		ready     bool
		values    *MapIterator
	}
)

func NewGroupIterator(
	src Iterator,
	selectors ...*GroupSelector,
) (*GroupIterator, error) {
	if core.IsNil(src) {
		return nil, core.Error(core.ErrMissedArgument, "source")
	}

	if len(selectors) == 0 {
		return nil, core.Error(core.ErrMissedArgument, "key(s)")
	}

	return &GroupIterator{src, selectors, false, nil}, nil
}

func (iterator *GroupIterator) HasNext() bool {
	if !iterator.ready {
		iterator.ready = true
		groups, err := iterator.group()

		if err != nil {
			iterator.values = NewMapIterator(map[string]core.Value{})

			return false
		}

		iterator.values = groups
	}

	return iterator.values.HasNext()
}

func (iterator *GroupIterator) Next() (ResultSet, error) {
	return iterator.values.Next()
}

func (iterator *GroupIterator) group() (*MapIterator, error) {
	groups := make(map[string]core.Value)

	for iterator.src.HasNext() {
		for _, selector := range iterator.selectors {
			set, err := iterator.src.Next()

			if err != nil {
				return nil, err
			}

			if len(set) == 0 {
				continue
			}

			keyVal, err := selector.Key(set)

			if err != nil {
				return nil, err
			}

			key := keyVal.String()

			group, exists := groups[key]

			if !exists {
				group = values.NewArray(10)
				groups[key] = group
			}

			val, err := selector.Value(set)

			if err != nil {
				return nil, err
			}

			group.(*values.Array).Push(val)
		}
	}

	return NewMapIterator(groups), nil
}
