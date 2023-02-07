package collections

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

const (
	DefaultValueVar = "value"
	DefaultKeyVar   = "key"
)

type IndexedIterator struct {
	valVar string
	keyVar string
	values IndexedCollection
	pos    int
}

func NewIndexedIterator(
	valVar,
	keyVar string,
	values IndexedCollection,
) (Iterator, error) {
	if valVar == "" {
		return nil, core.Error(core.ErrMissedArgument, "value variable")
	}

	if values == nil {
		return nil, core.Error(core.ErrMissedArgument, "result")
	}

	return &IndexedIterator{valVar, keyVar, values, 0}, nil
}

func NewDefaultIndexedIterator(
	values IndexedCollection,
) (Iterator, error) {
	return NewIndexedIterator(DefaultValueVar, DefaultKeyVar, values)
}

func (iterator *IndexedIterator) Next(_ context.Context, scope *core.Scope) (*core.Scope, error) {
	if int(iterator.values.Length()) > iterator.pos {
		idx := values.NewInt(iterator.pos)
		val := iterator.values.Get(idx)

		iterator.pos++

		nextScope := scope.Fork()

		if err := nextScope.SetVariable(iterator.valVar, val); err != nil {
			return nil, err
		}

		if iterator.keyVar != "" {
			if err := nextScope.SetVariable(iterator.keyVar, idx); err != nil {
				return nil, err
			}
		}

		return nextScope, nil
	}

	return nil, core.ErrNoMoreData
}
