package math_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

// Only traversal is implemented. Any accidental dependency on Length, At,
// Copy, sorting, or mutation fails through the nil embedded interface.
type aggregateList struct {
	runtime.List
	context context.Context
	visit   func(context.Context, int, int) error
	values  []runtime.Value
	calls   int
}

func (list *aggregateList) ForEach(ctx context.Context, predicate runtime.IndexReadablePredicate) error {
	list.context = ctx
	list.calls++

	for index, value := range list.values {
		if list.visit != nil {
			if err := list.visit(ctx, list.calls, index); err != nil {
				return err
			}
		}

		more, err := predicate(ctx, value, runtime.Int(index))
		if err != nil {
			return err
		}

		if !more {
			return nil
		}
	}

	if list.visit != nil {
		return list.visit(ctx, list.calls, len(list.values))
	}

	return nil
}
