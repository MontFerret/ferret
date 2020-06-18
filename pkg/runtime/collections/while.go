package collections

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	WhilePredicate func(ctx context.Context, scope *core.Scope) (bool, error)

	WhileIterator struct {
		valVar    string
		predicate WhilePredicate
		pos       int
	}
)

func NewWhileIterator(
	valVar string,
	predicate WhilePredicate,
) (Iterator, error) {
	if valVar == "" {
		return nil, core.Error(core.ErrMissedArgument, "value variable")
	}

	if predicate == nil {
		return nil, core.Error(core.ErrMissedArgument, "predicate")
	}

	return &WhileIterator{valVar, predicate, 0}, nil
}

func NewDefaultWhileIterator(predicate WhilePredicate) (Iterator, error) {
	return NewWhileIterator(DefaultValueVar, predicate)
}

func (iterator *WhileIterator) Next(ctx context.Context, scope *core.Scope) (*core.Scope, error) {
	doNext, err := iterator.predicate(ctx, scope)

	if err != nil {
		return nil, err
	}

	if doNext {
		counter := values.NewInt(iterator.pos)

		iterator.pos++

		nextScope := scope.Fork()

		if err := nextScope.SetVariable(iterator.valVar, counter); err != nil {
			return nil, err
		}

		return nextScope, nil
	}

	return nil, core.ErrNoMoreData
}
