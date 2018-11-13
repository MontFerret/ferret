package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	Collection interface {
		core.Value
		Length() values.Int
	}

	IndexedCollection interface {
		Collection
		Get(idx values.Int) core.Value
		Set(idx values.Int, value core.Value) error
	}

	KeyedCollection interface {
		Collection
		Keys() []string
		Get(key values.String) (core.Value, values.Boolean)
		Set(key values.String, value core.Value)
	}

	IterableCollection interface {
		core.Value
		Iterate(ctx context.Context) (CollectionIterator, error)
	}

	CollectionIterator interface {
		Next(ctx context.Context) (value core.Value, key core.Value, err error)
	}

	collectionIteratorWrapper struct {
		valVar string
		keyVar string
		values CollectionIterator
	}
)

func NewCollectionIterator(
	valVar,
	keyVar string,
	values CollectionIterator,
) (Iterator, error) {
	return &collectionIteratorWrapper{valVar, keyVar, values}, nil
}

func (iterator *collectionIteratorWrapper) Next(ctx context.Context, scope *core.Scope) (*core.Scope, error) {
	val, key, err := iterator.values.Next(ctx)

	if err != nil {
		return nil, err
	}

	// end of iteration
	if val == values.None {
		return nil, nil
	}

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
