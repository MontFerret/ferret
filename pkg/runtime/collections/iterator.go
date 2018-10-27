package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	Iterator interface {
		Next(ctx context.Context, scope *core.Scope) (DataSet, error)
	}

	Iterable interface {
		Iterate(ctx context.Context, scope *core.Scope) (Iterator, error)
	}
)

func ToSlice(ctx context.Context, scope *core.Scope, iterator Iterator) ([]DataSet, error) {
	res := make([]DataSet, 0, 10)

	for {
		innerScope := scope.Fork()
		ds, err := iterator.Next(ctx, innerScope)

		if err != nil {
			return nil, err
		}

		if ds == nil {
			return res, nil
		}

		res = append(res, ds)
	}

	return res, nil
}
