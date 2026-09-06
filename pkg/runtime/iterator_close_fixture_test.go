package runtime_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	traversalSource struct {
		iter runtime.Iterator
		err  error
	}
	traversalIterator struct {
		nextErr, closeErr error
		closes            int
	}
)

func (s *traversalSource) Iterate(context.Context) (runtime.Iterator, error) {
	return s.iter, s.err
}

func (i *traversalIterator) Next(context.Context) (runtime.Value, runtime.Value, error) {
	return runtime.True, runtime.ZeroInt, i.nextErr
}

func (i *traversalIterator) Close() error {
	i.closes++

	return i.closeErr
}
