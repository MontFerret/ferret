package driver

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/dynamic"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/static"
)

type Name string

const (
	Dynamic Name = "dynamic"
	Static  Name = "static"
)

type Driver interface {
	GetDocument(ctx context.Context, url values.String) (values.HTMLNode, error)
	Close() error
}

func ToContext(ctx context.Context, name Name, drv Driver) context.Context {
	return context.WithValue(ctx, name, drv)
}

func FromContext(ctx context.Context, name Name) (Driver, error) {
	val := ctx.Value(name)

	drv, ok := val.(Driver)

	if ok {
		return drv, nil
	}

	return nil, core.Error(core.ErrNotFound, fmt.Sprintf("%s driver", name))
}

func WithDynamicDriver(ctx context.Context, addr string) context.Context {
	return context.WithValue(ctx, Dynamic, dynamic.NewDriver(addr))
}

func WithStaticDriver(ctx context.Context, opts ...static.Option) context.Context {
	return context.WithValue(ctx, Static, static.NewDriver(opts...))
}
