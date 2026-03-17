package vm

import (
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

type (
	// Result wraps a VM run result together with any closers whose lifetime is
	// tied to that value. Callers must Close the result when they are done with
	// it. Use Root for low-level inspection while the result is open, or
	// Materialize for a final ownership-aware conversion.
	Result struct {
		root         runtime.Value
		set          mem.CloserSet
		closed       bool
		materialized bool
	}

	// Materialized is a typed value returned by a Materializer together with any
	// additional closers that should be released when the owning Result is
	// closed.
	Materialized[T any] struct {
		Value   T
		Closers []io.Closer
	}

	// Materializer converts a Result root into a typed value and may return
	// additional closers to be adopted by the owning Result.
	Materializer[T any] func(runtime.Value) (Materialized[T], error)
)

func newResult(root runtime.Value) *Result {
	return &Result{
		root: unwrapManaged(normalizeValue(root)),
	}
}

func (r *Result) reset(root runtime.Value) {
	r.root = unwrapManaged(normalizeValue(root))
	r.set.Reset()
	r.closed = false
	r.materialized = false
}

// Materialize converts an open Result into a typed value using materializer.
// It is a terminal, ownership-aware conversion: each Result may be
// materialized at most once, even if materializer returns an error. Callers
// remain responsible for closing the Result after materialization.
func Materialize[T any](r *Result, materializer Materializer[T]) (T, error) {
	var zero T

	if r == nil {
		return zero, runtime.Error(runtime.ErrInvalidOperation, "result is nil")
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

// Root returns the raw runtime value while the result is open. It is intended
// for low-level inspection; once the Result is closed, Root returns
// runtime.None.
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

	closer, ok := val.(io.Closer)
	if !ok {
		return
	}

	r.set.Add(closer)
}

func (r *Result) AdoptCloser(closer io.Closer) {
	if r == nil || r.closed {
		return
	}

	r.set.Add(closer)
}

func (r *Result) adoptOwned(owned *mem.OwnedResources) {
	if owned == nil {
		return
	}

	owned.ForEach(func(closer io.Closer) {
		r.set.Add(closer)
	})
}

func (r *Result) adoptDeferred(deferred *mem.DeferredClosers) {
	if deferred == nil {
		return
	}

	deferred.ForEach(func(closer io.Closer) {
		r.set.Add(closer)
	})
}

func (r *Result) Close() error {
	if r == nil || r.closed {
		return nil
	}

	r.closed = true

	err := r.set.CloseAll()

	r.root = runtime.None

	return err
}
