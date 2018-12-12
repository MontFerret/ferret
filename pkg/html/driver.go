package html

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"io"
)

type (
	staticCtxKey struct{}

	dynamicCtxKey struct{}

	DriverStatic interface {
		io.Closer
		GetDocument(ctx context.Context, url values.String) (values.HTMLDocument, error)
		ParseDocument(ctx context.Context, str values.String) (values.HTMLDocument, error)
	}

	DriverDynamic interface {
		io.Closer
		GetDocument(ctx context.Context, url values.String) (values.DHTMLDocument, error)
	}
)

func StaticFrom(ctx context.Context) (DriverStatic, error) {
	val := ctx.Value(staticCtxKey{})

	drv, ok := val.(DriverStatic)

	if !ok {
		return nil, core.Error(core.ErrNotFound, "HTML Driver")
	}

	return drv, nil
}

func DynamicFrom(ctx context.Context) (DriverDynamic, error) {
	val := ctx.Value(dynamicCtxKey{})

	drv, ok := val.(DriverDynamic)

	if !ok {
		return nil, core.Error(core.ErrNotFound, "DHTML Driver")
	}

	return drv, nil
}

func WithStatic(ctx context.Context, drv DriverStatic) context.Context {
	return context.WithValue(
		ctx,
		staticCtxKey{},
		drv,
	)
}

func WithDynamic(ctx context.Context, drv DriverDynamic) context.Context {
	return context.WithValue(
		ctx,
		dynamicCtxKey{},
		drv,
	)
}
