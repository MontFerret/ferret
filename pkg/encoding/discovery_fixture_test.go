package encoding_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	discoveryList struct {
		*runtime.Array
		acquire func()
	}
	discoveryMap struct {
		*runtime.Object
		value, key runtime.Value
		acquire    func()
	}
	discoveryContext struct {
		context.Context
		lookups *int
	}
)

func (v *discoveryList) At(ctx context.Context, index runtime.Int) (runtime.Value, error) {
	if v.acquire != nil {
		v.acquire()
	}
	return v.Array.At(ctx, index)
}
func (v *discoveryMap) ForEach(_ context.Context, predicate runtime.KeyReadablePredicate) error {
	if v.acquire != nil {
		v.acquire()
	}
	_, err := predicate(context.Background(), v.value, v.key)
	return err
}
func (ctx discoveryContext) Value(key any) any { *ctx.lookups++; return ctx.Context.Value(key) }
