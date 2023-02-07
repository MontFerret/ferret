package collections

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type KeyedIterator struct {
	valVar string
	keyVar string
	values KeyedCollection
	keys   []values.String
	pos    int
}

func NewKeyedIterator(
	valVar,
	keyVar string,
	values KeyedCollection,
) (Iterator, error) {
	if valVar == "" {
		return nil, core.Error(core.ErrMissedArgument, "value variable")
	}

	if values == nil {
		return nil, core.Error(core.ErrMissedArgument, "result")
	}

	return &KeyedIterator{valVar, keyVar, values, nil, 0}, nil
}

func NewDefaultKeyedIterator(input KeyedCollection) (Iterator, error) {
	return NewKeyedIterator(DefaultValueVar, DefaultKeyVar, input)
}

func (iterator *KeyedIterator) Next(_ context.Context, scope *core.Scope) (*core.Scope, error) {
	if iterator.keys == nil {
		iterator.keys = iterator.values.Keys()
	}

	if len(iterator.keys) > iterator.pos {
		key := iterator.keys[iterator.pos]
		val, _ := iterator.values.Get(key)

		iterator.pos++

		nextScope := scope.Fork()

		if err := nextScope.SetVariable(iterator.valVar, val); err != nil {
			return nil, err
		}

		if iterator.keyVar != "" {
			if err := nextScope.SetVariable(iterator.keyVar, key); err != nil {
				return nil, err
			}
		}

		return nextScope, nil
	}

	return nil, core.ErrNoMoreData
}
