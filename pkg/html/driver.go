package html

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/html/dynamic"
	"github.com/MontFerret/ferret/pkg/html/static"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/env"
	"github.com/MontFerret/ferret/pkg/runtime/values"
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

func WithDynamicDriver(ctx context.Context) context.Context {
	e := env.FromContext(ctx)

	return context.WithValue(
		ctx,
		Dynamic,
		dynamic.NewDriver(
			e.CDPAddress,
			dynamic.WithProxy(e.ProxyAddress),
			dynamic.WithUserAgent(e.UserAgent),
		),
	)
}

func WithStaticDriver(ctx context.Context) context.Context {
	e := env.FromContext(ctx)

	return context.WithValue(
		ctx,
		Static,
		static.NewDriver(
			static.WithProxy(e.ProxyAddress),
			static.WithUserAgent(e.UserAgent),
		),
	)
}
