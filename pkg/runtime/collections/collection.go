package collections

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	// Measurable represents an interface of a value that can has length.
	Measurable interface {
		Length() values.Int
	}

	IndexedCollection interface {
		core.Value
		Measurable
		Get(idx values.Int) core.Value
		Set(idx values.Int, value core.Value) error
	}

	KeyedCollection interface {
		core.Value
		Measurable
		Keys() []values.String
		Get(key values.String) (core.Value, values.Boolean)
		Set(key values.String, value core.Value)
	}

	coreIterator struct {
		valVar string
		keyVar string
		values core.Iterator
	}
)

func NewCoreIterator(valVar, keyVar string, values core.Iterator) (Iterator, error) {
	return &coreIterator{valVar, keyVar, values}, nil
}

func (iterator *coreIterator) Next(ctx context.Context, scope *core.Scope) (*core.Scope, error) {
	val, key, err := iterator.values.Next(ctx)

	if err != nil {
		return nil, err
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
