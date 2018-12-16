package drivers

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"io"
)

type (
	staticCtxKey struct{}

	dynamicCtxKey struct{}

	Static interface {
		io.Closer
		GetDocument(ctx context.Context, url values.String) (values.HTMLDocument, error)
		ParseDocument(ctx context.Context, str values.String) (values.HTMLDocument, error)
	}

	Dynamic interface {
		io.Closer
		GetDocument(ctx context.Context, url values.String) (values.DHTMLDocument, error)
	}
)

func StaticFrom(ctx context.Context) (Static, error) {
	val := ctx.Value(staticCtxKey{})

	drv, ok := val.(Static)

	if !ok {
		return nil, core.Error(core.ErrNotFound, "HTML Driver")
	}

	return drv, nil
}

func DynamicFrom(ctx context.Context) (Dynamic, error) {
	val := ctx.Value(dynamicCtxKey{})

	drv, ok := val.(Dynamic)

	if !ok {
		return nil, core.Error(core.ErrNotFound, "DHTML Driver")
	}

	return drv, nil
}

func WithStatic(ctx context.Context, drv Static) context.Context {
	return context.WithValue(
		ctx,
		staticCtxKey{},
		drv,
	)
}

func WithDynamic(ctx context.Context, drv Dynamic) context.Context {
	return context.WithValue(
		ctx,
		dynamicCtxKey{},
		drv,
	)
}
