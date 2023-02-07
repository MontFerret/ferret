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

func (iterator *LimitIterator) Next(ctx context.Context, scope *core.Scope) (*core.Scope, error) {
	if err := iterator.verifyOffset(ctx, scope); err != nil {
		return nil, err
	}

	iterator.currCount++

	if iterator.counter() <= iterator.count {
		return iterator.values.Next(ctx, scope)
	}

	return nil, core.ErrNoMoreData
}

func (iterator *LimitIterator) counter() int {
	return iterator.currCount - iterator.offset
}

func (iterator *LimitIterator) verifyOffset(ctx context.Context, scope *core.Scope) error {
	if iterator.offset == 0 {
		return nil
	}

	for iterator.offset > iterator.currCount {
		_, err := iterator.values.Next(ctx, scope.Fork())

		if err != nil {
			if core.IsNoMoreData(err) {
				iterator.currCount = iterator.offset
			}

			return err
		}

		iterator.currCount++
	}

	return nil
}
