package session

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	abandonmentMap struct {
		*runtime.Object
		value, key runtime.Value
	}
	abandonmentKey struct {
		runtime.Value
		mutate func()
	}
	abandonmentResource struct {
		runtime.Value
		closes int
	}
	abandonmentWriter struct {
		failure error
		writes  int
	}
)

func (v *abandonmentMap) ForEach(_ context.Context, predicate runtime.KeyReadablePredicate) error {
	_, err := predicate(context.Background(), v.value, v.key)
	return err
}
func (v *abandonmentKey) String() string    { v.mutate(); return "key" }
func (v *abandonmentResource) Close() error { v.closes++; return nil }
func (w *abandonmentWriter) Write(p []byte) (int, error) {
	w.writes++
	if w.writes > 1 {
		return 0, w.failure
	}
	return len(p), nil
}
