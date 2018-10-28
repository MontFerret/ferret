package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	FilterPredicate func(ctx context.Context, scope *core.Scope) (bool, error)

	FilterIterator struct {
		values    Iterator
		predicate FilterPredicate
	}
)

func NewFilterIterator(values Iterator, predicate FilterPredicate) (*FilterIterator, error) {
	if values == nil {
		return nil, core.Error(core.ErrMissedArgument, "result")
	}

	if predicate == nil {
		return nil, core.Error(core.ErrMissedArgument, "predicate")
	}

	return &FilterIterator{values: values, predicate: predicate}, nil
}

func (iterator *FilterIterator) Next(ctx context.Context, scope *core.Scope) (*core.Scope, error) {
	for {
		nextScope, err := iterator.values.Next(ctx, scope.Fork())

		if err != nil {
			return nil, err
		}

		if nextScope == nil {
			return nil, nil
		}

		take, err := iterator.predicate(ctx, nextScope)

		if take == true {
			return nextScope, nil
		}
	}

	return nil, nil
}
