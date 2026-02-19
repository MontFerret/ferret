package ferret

import (
	"context"
	"errors"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
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
		peek  runtime.Value
		has   bool
		done  bool
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

	if r.done {
		return false, nil
	}

	if r.iter == nil {
		r.iter, r.err = r.value.Iterate(ctx)

		if r.err != nil {
			return false, r.err
		}
	}

	if r.has {
		return true, nil
	}

	val, _, err := r.iter.Next(ctx)
	if errors.Is(err, io.EOF) {
		r.done = true
		return false, nil
	}
	if err != nil {
		r.err = err
		return false, err
	}

	r.peek = val
	r.has = true
	return true, nil
}

func (r *rowsResult) Next(ctx context.Context) (runtime.Value, error) {
	if r.err != nil {
		return runtime.None, r.err
	}

	if r.done {
		return runtime.None, io.EOF
	}

	if r.iter == nil {
		return runtime.None, io.ErrUnexpectedEOF
	}

	if r.has {
		r.has = false
		return r.peek, nil
	}

	val, _, err := r.iter.Next(ctx)

	if errors.Is(err, io.EOF) {
		r.done = true
		return runtime.None, io.EOF
	}

	if err != nil {
		r.err = err
		return runtime.None, err
	}

	return val, nil
}
