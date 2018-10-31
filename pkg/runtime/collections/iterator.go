package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	Iterator interface {
		Next(ctx context.Context, scope *core.Scope) (*core.Scope, error)
	}

	Iterable interface {
		Iterate(ctx context.Context, scope *core.Scope) (Iterator, error)
	}
)

func ToSlice(ctx context.Context, scope *core.Scope, iterator Iterator) ([]*core.Scope, error) {
	res := make([]*core.Scope, 0, 10)

	for {
		nextScope, err := iterator.Next(ctx, scope.Fork())

		if err != nil {
			return nil, err
		}

		if nextScope == nil {
			return res, nil
		}

		res = append(res, nextScope)
	}
}
