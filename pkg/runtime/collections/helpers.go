package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	IterableFn struct {
		fn func(ctx context.Context, scope *core.Scope) (Iterator, error)
	}

	IteratorFn struct {
		fn func(ctx context.Context, scope *core.Scope) (*core.Scope, error)
	}
)

func AsIterable(fn func(ctx context.Context, scope *core.Scope) (Iterator, error)) Iterable {
	return &IterableFn{fn}
}

func (i *IterableFn) Iterate(ctx context.Context, scope *core.Scope) (Iterator, error) {
	return i.fn(ctx, scope)
}

func AsIterator(fn func(ctx context.Context, scope *core.Scope) (*core.Scope, error)) Iterator {
	return &IteratorFn{fn}
}

func (i *IteratorFn) Next(ctx context.Context, scope *core.Scope) (*core.Scope, error) {
	return i.fn(ctx, scope)
}

func ToSlice(ctx context.Context, scope *core.Scope, iterator Iterator) ([]*core.Scope, error) {
	res := make([]*core.Scope, 0, 10)

	for {
		nextScope, err := iterator.Next(ctx, scope.Fork())

		if err != nil {
			if core.IsNoMoreData(err) {
				return res, nil
			}

			return nil, err
		}

		res = append(res, nextScope)
	}
}

func ForEach(ctx context.Context, scope *core.Scope, iter Iterator, predicate func(ctx context.Context, scope *core.Scope) bool) error {
	for {
		nextScope, err := iter.Next(ctx, scope)

		if err != nil {
			if core.IsNoMoreData(err) {
				return nil
			}

			return err
		}

		if !predicate(ctx, nextScope) {
			return nil
		}
	}
}
