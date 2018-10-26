package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type LimitIterator struct {
	values    Iterator
	count     int
	offset    int
	currCount int
}

func NewLimitIterator(values Iterator, count, offset int) (*LimitIterator, error) {
	if values == nil {
		return nil, core.Error(core.ErrMissedArgument, "result")
	}

	return &LimitIterator{values, count, offset, 0}, nil
}

func (iterator *LimitIterator) Next(ctx context.Context, scope *core.Scope) (DataSet, error) {
	if err := iterator.verifyOffset(ctx, scope); err != nil {
		return nil, err
	}

	iterator.currCount++

	if iterator.counter() <= iterator.count {
		return iterator.values.Next(ctx, scope)
	}

	return nil, nil
}

func (iterator *LimitIterator) counter() int {
	return iterator.currCount - iterator.offset
}

func (iterator *LimitIterator) verifyOffset(ctx context.Context, scope *core.Scope) error {
	if iterator.offset == 0 {
		return nil
	}

	for iterator.offset > iterator.currCount {
		ds, err := iterator.values.Next(ctx, scope)

		if err != nil {
			return err
		}

		if ds == nil {
			iterator.currCount = iterator.offset
			return nil
		}

		iterator.currCount++
	}

	return nil
}
