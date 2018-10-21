package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	Iterator interface {
		HasNext() bool
		Next() (DataSet, error)
	}

	Iterable interface {
		Variables() Variables
		Iterate(ctx context.Context, scope *core.Scope) (Iterator, error)
	}
)

func ToSlice(iterator Iterator) ([]DataSet, error) {
	res := make([]DataSet, 0, 10)

	for iterator.HasNext() {
		ds, err := iterator.Next()

		if err != nil {
			return nil, err
		}

		res = append(res, ds)
	}

	return res, nil
}
