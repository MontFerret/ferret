package collections

import (
	"context"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	WhileMode int

	WhilePredicate func(ctx context.Context, scope *core.Scope) (bool, error)

	WhileIterator struct {
		mode      WhileMode
		predicate WhilePredicate
		valVar    string
		pos       int
	}
)

var (
	WhileModePost WhileMode
	WhileModePre  WhileMode = 1
)

func NewWhileIterator(
	mode WhileMode,
	valVar string,
	predicate WhilePredicate,
) (Iterator, error) {
	if valVar == "" {
		return nil, core.Error(core.ErrMissedArgument, "value variable")
	}

	if predicate == nil {
		return nil, core.Error(core.ErrMissedArgument, "predicate")
	}

	return &WhileIterator{mode, predicate, valVar, 0}, nil
}

func NewDefaultWhileIterator(mode WhileMode, predicate WhilePredicate) (Iterator, error) {
	return NewWhileIterator(mode, DefaultValueVar, predicate)
}

func (iterator *WhileIterator) Next(ctx context.Context, scope *core.Scope) (*core.Scope, error) {
	// if it's Post conditional execution, step in always
	// Otherwise, it's not the first iteration
	if iterator.mode == WhileModePost || iterator.pos > 0 {
		doNext, err := iterator.predicate(ctx, scope)

		if err != nil {
			return nil, err
		}

		if !doNext {
			return nil, core.ErrNoMoreData
		}
	}

	counter := values.NewInt(iterator.pos)
	iterator.pos++

	nextScope := scope.Fork()

	if err := nextScope.SetVariable(iterator.valVar, counter); err != nil {
		return nil, err
	}

	return nextScope, nil
}
