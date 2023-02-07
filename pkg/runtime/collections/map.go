package collections

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type MapIterator struct {
	valVar string
	keyVar string
	values map[string]core.Value
	keys   []string
	pos    int
}

func NewMapIterator(
	valVar,
	keyVar string,
	values map[string]core.Value,
) (Iterator, error) {
	if valVar == "" {
		return nil, core.Error(core.ErrMissedArgument, "value variable")
	}

	if values == nil {
		return nil, core.Error(core.ErrMissedArgument, "result")
	}

	return &MapIterator{valVar, keyVar, values, nil, 0}, nil
}

func NewDefaultMapIterator(values map[string]core.Value) (Iterator, error) {
	return NewMapIterator(DefaultValueVar, DefaultKeyVar, values)
}

func (iterator *MapIterator) Next(_ context.Context, scope *core.Scope) (*core.Scope, error) {
	// lazy initialization
	if iterator.keys == nil {
		keys := make([]string, len(iterator.values))

		i := 0
		for k := range iterator.values {
			keys[i] = k
			i++
		}

		iterator.keys = keys
	}

	if len(iterator.keys) > iterator.pos {
		key := iterator.keys[iterator.pos]
		val := iterator.values[key]

		iterator.pos++

		nextScope := scope.Fork()

		if err := nextScope.SetVariable(iterator.valVar, val); err != nil {
			return nil, err
		}

		if iterator.keyVar != "" {
			if err := nextScope.SetVariable(iterator.keyVar, values.NewString(key)); err != nil {
				return nil, err
			}
		}

		return nextScope, nil
	}

	return nil, core.ErrNoMoreData
}
