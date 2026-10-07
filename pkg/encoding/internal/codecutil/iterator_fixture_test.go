package codecutil_test

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	acquiredIterator struct {
		ctx      context.Context
		iterator runtime.Iterator
		err      error
		marker   any
		closes   *int
	}

	identifiedIterator struct {
		acquiredIterator
		id uint64
	}
)

func (v acquiredIterator) Iterate(ctx context.Context) (runtime.Iterator, error) {
	if v.ctx != nil && v.ctx != ctx {
		panic("lost iterator context")
	}

	if v.iterator == nil {
		return v, v.err
	}

	return v.iterator, v.err
}

func (v acquiredIterator) Next(context.Context) (runtime.Value, runtime.Value, error) {
	return runtime.None, runtime.None, io.EOF
}

func (v acquiredIterator) Close() error {
	*v.closes += 1

	return nil
}

func (v identifiedIterator) ResourceID() uint64 {
	return v.id
}
