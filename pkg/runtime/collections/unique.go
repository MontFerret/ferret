package collections

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	UniqueIterator struct {
		values  Iterator
		hashes  map[uint64]bool
		hashKey string
	}
)

func NewUniqueIterator(values Iterator, hashKey string) (*UniqueIterator, error) {
	if values == nil {
		return nil, core.Error(core.ErrMissedArgument, "source")
	}

	return &UniqueIterator{
		values:  values,
		hashes:  make(map[uint64]bool),
		hashKey: hashKey,
	}, nil
}

func (iterator *UniqueIterator) Next(ctx context.Context, scope *core.Scope) (DataSet, error) {
	for {
		ds, err := iterator.values.Next(ctx, scope)

		if err != nil {
			return nil, err
		}

		if ds == nil {
			return nil, nil
		}

		h := ds.Get(iterator.hashKey).Hash()

		_, exists := iterator.hashes[h]

		if exists {
			continue
		}

		iterator.hashes[h] = true

		return ds, nil
	}
}
