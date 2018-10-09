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

type (
	DriverName    string
	dynamicCtxKey struct{}
	staticCtxKey  struct{}
)

const (
	Dynamic DriverName = "dynamic"
	Static  DriverName = "static"
)

type Driver interface {
	GetDocument(ctx context.Context, url values.String) (values.HTMLNode, error)
	Close() error
}

func ToContext(ctx context.Context, name DriverName, drv Driver) context.Context {
	var key interface{}

	switch name {
	case Dynamic:
		key = dynamicCtxKey{}
	case Static:
		key = staticCtxKey{}
	default:
		return ctx
	}

	return context.WithValue(ctx, key, drv)
}

func FromContext(ctx context.Context, name DriverName) (Driver, error) {
	var key interface{}

	switch name {
	case Dynamic:
		key = dynamicCtxKey{}
	case Static:
		key = staticCtxKey{}
	default:
		return nil, core.Error(core.ErrInvalidArgument, fmt.Sprintf("%s driver", name))
	}

	val := ctx.Value(key)

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
		dynamicCtxKey{},
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
		staticCtxKey{},
		static.NewDriver(
			static.WithProxy(e.ProxyAddress),
			static.WithUserAgent(e.UserAgent),
		),
	)
}
