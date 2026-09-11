package vm_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	ownedCollectionMap struct {
		writeErr error
		closeErr error
		*runtime.Object
		created *ownedCollectionMap
		backend string
		closes  int
	}

	measuredCollectionIterable struct {
		*collectionIterable
		length       runtime.Int
		measurements int
		iterations   int
	}
)

func (m *ownedCollectionMap) New(ctx context.Context) (runtime.Map, error) {
	base, err := m.Object.New(ctx)
	if err != nil {
		return nil, err
	}

	m.created = &ownedCollectionMap{Object: base.(*runtime.Object), backend: m.backend, writeErr: m.writeErr, closeErr: m.closeErr}

	return m.created, nil
}

func (m *ownedCollectionMap) Set(ctx context.Context, key, value runtime.Value) error {
	if m.writeErr != nil {
		return m.writeErr
	}

	return m.Object.Set(ctx, key, value)
}

func (m *ownedCollectionMap) Close() error {
	m.closes++

	return m.closeErr
}

func (m *measuredCollectionIterable) Length(context.Context) (runtime.Int, error) {
	m.measurements++

	return m.length, nil
}

func (m *measuredCollectionIterable) Iterate(ctx context.Context) (runtime.Iterator, error) {
	m.iterations++

	return m.collectionIterable.Iterate(ctx)
}
