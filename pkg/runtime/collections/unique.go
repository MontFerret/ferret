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

func (iterator *UniqueIterator) Next(ctx context.Context, scope *core.Scope) (*core.Scope, error) {
	for {
		nextScope, err := iterator.values.Next(ctx, scope.Fork())

		if err != nil {
			return nil, err
		}

		v, err := nextScope.GetVariable(iterator.hashKey)

		if err != nil {
			return nil, err
		}

		h := v.Hash()

		_, exists := iterator.hashes[h]

		if exists {
			continue
		}

		iterator.hashes[h] = true

		return nextScope, nil
	}
}
