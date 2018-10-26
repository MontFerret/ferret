package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type TapIterator struct {
	values    Iterator
	predicate core.Expression
}

func NewTapIterator(values Iterator, predicate core.Expression) (Iterator, error) {
	if values == nil {
		return nil, core.Error(core.ErrMissedArgument, "values")
	}

	if predicate == nil {
		return nil, core.Error(core.ErrMissedArgument, "predicate")
	}

	return &TapIterator{values, predicate}, nil
}

func (iterator *TapIterator) Next(ctx context.Context, scope *core.Scope) (DataSet, error) {
	for {
		ds, err := iterator.values.Next(ctx, scope)

		if err != nil {
			return nil, err
		}

		if ds == nil {
			break
		}

		_, err = iterator.predicate.Exec(ctx, scope)

		if err != nil {
			return nil, err
		}

		return ds, nil
	}

	return nil, nil
}
