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

	DriverHTML interface {
		io.Closer
		GetDocument(ctx context.Context, url values.String) (values.HTMLDocument, error)
		ParseDocument(ctx context.Context, str values.String) (values.HTMLDocument, error)
	}

	DriverDHTML interface {
		io.Closer
		GetDocument(ctx context.Context, url values.String) (values.DHTMLDocument, error)
	}
)

func FromContextHTML(ctx context.Context) (DriverHTML, error) {
	val := ctx.Value(staticCtxKey{})

	drv, ok := val.(DriverHTML)

	if !ok {
		return nil, core.Error(core.ErrNotFound, "HTML Driver")
	}

	return drv, nil
}

func FromContextDHTML(ctx context.Context) (DriverDHTML, error) {
	val := ctx.Value(dynamicCtxKey{})

	drv, ok := val.(DriverDHTML)

	if !ok {
		return nil, core.Error(core.ErrNotFound, "DHTML Driver")
	}

	return drv, nil
}

func WithContextHTML(ctx context.Context, drv DriverHTML) context.Context {
	return context.WithValue(
		ctx,
		staticCtxKey{},
		drv,
	)
}

func WithContextDHTML(ctx context.Context, drv DriverDHTML) context.Context {
	return context.WithValue(
		ctx,
		dynamicCtxKey{},
		drv,
	)
}
