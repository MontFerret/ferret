package vm_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	collectionIterable struct {
		values *runtime.Array
	}
	ownedCollectionList struct {
		*runtime.Array
		created *ownedCollectionList
		backend string
		closes  int
	}
)

func (s *collectionIterable) String() string      { return "read-only collection source" }
func (s *collectionIterable) Hash() uint64        { return 1 }
func (s *collectionIterable) Copy() runtime.Value { return s }
func (s *collectionIterable) Iterate(ctx context.Context) (runtime.Iterator, error) {
	return s.values.Iterate(ctx)
}

func (l *ownedCollectionList) New(ctx context.Context) (runtime.List, error) {
	base, err := l.Array.New(ctx)
	if err != nil {
		return nil, err
	}

	l.created = &ownedCollectionList{Array: base.(*runtime.Array), backend: l.backend}

	return l.created, nil
}

func (l *ownedCollectionList) Close() error {
	l.closes++

	return nil
}
