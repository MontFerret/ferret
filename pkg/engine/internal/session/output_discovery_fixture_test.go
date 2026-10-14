package session

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	abandonmentMap struct {
		*runtime.Object
		value, key   runtime.Value
		forEachCalls int
	}

	abandonmentKey struct {
		runtime.Value
		mutate func()
	}

	abandonmentResource struct {
		runtime.Value
		closeErr error
		closes   int
	}

	abandonmentWriter struct {
		failure error
		writes  int
	}

	abandonmentIterable struct {
		runtime.Value
		item     runtime.Value
		iterates int
		nexts    int
	}
)

func (v *abandonmentMap) ForEach(_ context.Context, predicate runtime.KeyReadablePredicate) error {
	v.forEachCalls++
	_, err := predicate(context.Background(), v.value, v.key)

	return err
}

func (v *abandonmentKey) String() string {
	v.mutate()

	return "key"
}

func (v *abandonmentResource) Close() error {
	v.closes++

	return v.closeErr
}

func (w *abandonmentWriter) Write(p []byte) (int, error) {
	w.writes++
	if w.writes > 1 {
		return 0, w.failure
	}

	return len(p), nil
}

func (v *abandonmentIterable) Iterate(context.Context) (runtime.Iterator, error) {
	v.iterates++

	return v, nil
}

func (v *abandonmentIterable) Next(context.Context) (runtime.Value, runtime.Value, error) {
	v.nexts++
	if v.nexts > 1 {
		return runtime.None, runtime.None, io.EOF
	}

	return v.item, runtime.None, nil
}
