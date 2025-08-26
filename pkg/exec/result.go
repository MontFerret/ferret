package exec

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type (
	Result interface {
		io.Closer

		HasNext(ctx context.Context) (bool, error)
		Next(ctx context.Context) (runtime.Value, error)
	}

	scalarResult struct {
		value    runtime.Value
		consumed bool
	}

	rowsResult struct {
		value runtime.Iterable
		iter  runtime.Iterator
		err   error
	}
)

func newResult(val runtime.Value) Result {
	iterable, ok := val.(runtime.Iterable)

	if !ok {
		return &scalarResult{value: val}
	}

	_, ok = val.(runtime.Map)

	if ok {
		return &scalarResult{value: val}
	}

	return &rowsResult{value: iterable}
}

func (r *scalarResult) Close() error {
	closable, ok := r.value.(io.Closer)

	if !ok {
		return nil
	}

	return closable.Close()
}

func (r *scalarResult) HasNext(_ context.Context) (bool, error) {
	return !r.consumed, nil
}

func (r *scalarResult) Next(_ context.Context) (runtime.Value, error) {
	if r.consumed {
		return runtime.None, io.EOF
	}

	r.consumed = true
	return r.value, nil
}

func (r *rowsResult) Close() error {
	closable, ok := r.value.(io.Closer)

	if !ok {
		return nil
	}

	return closable.Close()
}

func (r *rowsResult) HasNext(ctx context.Context) (bool, error) {
	if r.err != nil {
		return false, r.err
	}

	if r.iter == nil {
		r.iter, r.err = r.value.Iterate(ctx)

		if r.err != nil {
			return false, r.err
		}
	}

	return r.iter.HasNext(ctx)
}

func (r *rowsResult) Next(ctx context.Context) (runtime.Value, error) {
	if r.err != nil {
		return runtime.None, r.err
	}

	if r.iter == nil {
		return runtime.None, io.ErrUnexpectedEOF
	}

	val, _, err := r.iter.Next(ctx)

	if err != nil {
		r.err = err
		return runtime.None, err
	}

	return val, nil
}
