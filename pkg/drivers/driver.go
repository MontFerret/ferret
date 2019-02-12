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
		opts    *Options
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

func WithContextDefault(ctx context.Context, name string) context.Context {
	ctx, value := resolveValue(ctx)

	value.opts.defaultDriver = name

	return ctx
}

func FromContextDefault(ctx context.Context) (Driver, error) {
	_, value := resolveValue(ctx)

	drv, found := value.drivers[value.opts.defaultDriver]

	if !found {
		return nil, core.Error(core.ErrNotFound, value.opts.defaultDriver)
	}

	return drv, nil
}

func resolveValue(ctx context.Context) (context.Context, *ctxValue) {
	key := ctxKey{}
	v := ctx.Value(key)
	value, ok := v.(*ctxValue)

	if !ok {
		value = &ctxValue{
			opts:    &Options{},
			drivers: make(map[string]Driver),
		}

		return context.WithValue(ctx, key, value), value
	}

	return ctx, value
}
