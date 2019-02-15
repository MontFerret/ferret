package drivers

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	ctxKey struct{}

	ctxValue struct {
		opts    *options
		drivers map[string]Driver
	}

	Driver interface {
		io.Closer
		Name() string
		GetDocument(ctx context.Context, url values.String) (HTMLDocument, error)
	}
)

func WithContext(ctx context.Context, drv Driver) context.Context {
	ctx, value := resolveValue(ctx)

	value.drivers[drv.Name()] = drv

	return ctx
}

func FromContext(ctx context.Context, name string) (Driver, error) {
	_, value := resolveValue(ctx)
	drv, exists := value.drivers[name]

	if !exists {
		return nil, core.Error(core.ErrNotFound, name)
	}

	return drv, nil
}

func resolveValue(ctx context.Context) (context.Context, *ctxValue) {
	key := ctxKey{}
	v := ctx.Value(key)
	value, ok := v.(*ctxValue)

	if !ok {
		value = &ctxValue{
			opts:    &options{},
			drivers: make(map[string]Driver),
		}

		return context.WithValue(ctx, key, value), value
	}

	return ctx, value
}
