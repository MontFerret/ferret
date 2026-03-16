package vm

import (
	"errors"
	"io"
	"reflect"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

type (
	Result struct {
		root         runtime.Value
		seen         map[io.Closer]struct{}
		closers      []io.Closer
		closed       bool
		materialized bool
	}

	Materialized[T any] struct {
		Value   T
		Closers []io.Closer
	}

	Materializer[T any] func(runtime.Value) (Materialized[T], error)
)

func newResult(root runtime.Value) *Result {
	return &Result{
		root: normalizeValue(root),
	}
}

func (r *Result) reset(root runtime.Value) {
	r.root = normalizeValue(root)
	r.closers = nil
	r.seen = nil
	r.closed = false
	r.materialized = false
}

func Materialize[T any](r *Result, materializer Materializer[T]) (T, error) {
	var zero T

	if r == nil {
		return zero, runtime.Error(runtime.ErrInvalidOperation, "result is closed")
	}

	if r.closed {
		return zero, runtime.Error(runtime.ErrInvalidOperation, "result is closed")
	}

	if r.materialized {
		return zero, runtime.Error(runtime.ErrInvalidOperation, "result is already materialized")
	}

	r.materialized = true

	m, err := materializer(r.root)

	if err != nil {
		return zero, err
	}

	for _, closer := range m.Closers {
		r.AdoptCloser(closer)
	}

	return m.Value, nil
}

func (r *Result) Root() runtime.Value {
	if r == nil || r.closed {
		return runtime.None
	}

	return r.root
}

func (r *Result) AdoptValue(val runtime.Value) {
	if r == nil || r.closed {
		return
	}

	closer, ok := comparableCloserOf(val)
	if !ok {
		return
	}

	r.AdoptCloser(closer)
}

func (r *Result) AdoptCloser(closer io.Closer) {
	if r == nil || r.closed {
		return
	}

	comparable, ok := comparableCloser(closer)
	if !ok {
		return
	}

	if r.seen == nil {
		r.seen = make(map[io.Closer]struct{})
	}

	if _, exists := r.seen[comparable]; exists {
		return
	}

	r.seen[comparable] = struct{}{}
	r.closers = append(r.closers, comparable)
}

func (r *Result) adoptOwned(owned *mem.OwnedResources) {
	if owned == nil {
		return
	}

	owned.ForEach(func(closer io.Closer) {
		r.AdoptCloser(closer)
	})
}

func (r *Result) adoptDeferred(deferred *mem.DeferredClosers) {
	if deferred == nil {
		return
	}

	deferred.ForEach(func(closer io.Closer) {
		r.AdoptCloser(closer)
	})
}

func (r *Result) Close() error {
	if r == nil || r.closed {
		return nil
	}

	r.closed = true

	var err error

	for _, closer := range r.closers {
		if e := closer.Close(); e != nil {
			err = errors.Join(err, e)
		}
	}

	r.root = runtime.None
	r.closers = nil
	r.seen = nil

	return err
}

func comparableCloserOf(val runtime.Value) (io.Closer, bool) {
	closer, ok := val.(io.Closer)
	if !ok {
		return nil, false
	}

	return comparableCloser(closer)
}

func comparableCloser(closer io.Closer) (io.Closer, bool) {
	if closer == nil {
		return nil, false
	}

	typ := reflect.TypeOf(closer)
	if typ == nil || !typ.Comparable() {
		return nil, false
	}

	return closer, true
}
