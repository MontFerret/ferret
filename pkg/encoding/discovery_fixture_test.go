package encoding_test

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	discoveryList struct {
		*runtime.Array
		acquire func()
		length  func()
	}
	discoveryMap struct {
		*runtime.Object
		value, key runtime.Value
		acquire    func()
		length     func()
	}
	discoveryContext struct {
		context.Context
		lookups *int
	}

	discoveryMarshaler struct {
		runtime.Value
		marshal func()
	}
)

func (v *discoveryList) Length(ctx context.Context) (runtime.Int, error) {
	if v.length != nil {
		v.length()
	}

	return v.Array.Length(ctx)
}

func (v *discoveryList) At(ctx context.Context, index runtime.Int) (runtime.Value, error) {
	if v.acquire != nil {
		v.acquire()
	}

	return v.Array.At(ctx, index)
}

func (v *discoveryMap) Length(ctx context.Context) (runtime.Int, error) {
	if v.length != nil {
		v.length()
	}

	return v.Object.Length(ctx)
}

func (v *discoveryMap) ForEach(_ context.Context, predicate runtime.KeyReadablePredicate) error {
	if v.acquire != nil {
		v.acquire()
	}

	_, err := predicate(context.Background(), v.value, v.key)

	return err
}

func (ctx discoveryContext) Value(key any) any {
	*ctx.lookups++

	return ctx.Context.Value(key)
}

func (v *discoveryMarshaler) MarshalJSON() ([]byte, error) {
	v.marshal()

	return []byte("null"), nil
}

func (v *discoveryMarshaler) MarshalMsgpack() ([]byte, error) {
	v.marshal()

	return []byte{0xc0}, nil
}
