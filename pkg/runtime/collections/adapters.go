package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	coreIterableAdapter struct {
		valVar   string
		keyVar   string
		iterable core.Iterable
	}

	coreIteratorAdapter struct {
		valVar   string
		keyVar   string
		iterator core.Iterator
	}
)

func FromCoreIterable(valVar, keyVar string, iterable core.Iterable) (Iterable, error) {
	if valVar == "" {
		return nil, core.Error(core.ErrMissedArgument, "value variable")
	}

	if iterable == nil {
		return nil, core.Error(core.ErrMissedArgument, "iterable")
	}

	return &coreIterableAdapter{valVar, keyVar, iterable}, nil
}

func (c *coreIterableAdapter) Iterate(ctx context.Context, _ *core.Scope) (Iterator, error) {
	iter, err := c.iterable.Iterate(ctx)

	if err != nil {
		return nil, err
	}

	return FromCoreIterator(c.valVar, c.keyVar, iter)
}

func FromCoreIterator(valVar, keyVar string, iterator core.Iterator) (Iterator, error) {
	if valVar == "" {
		return nil, core.Error(core.ErrMissedArgument, "value variable")
	}

	if iterator == nil {
		return nil, core.Error(core.ErrMissedArgument, "iterator")
	}

	return &coreIteratorAdapter{valVar, keyVar, iterator}, nil
}

func (iterator *coreIteratorAdapter) Next(ctx context.Context, scope *core.Scope) (*core.Scope, error) {
	val, key, err := iterator.iterator.Next(ctx)

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
