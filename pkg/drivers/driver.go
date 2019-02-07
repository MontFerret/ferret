package drivers

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	ctxKey struct{}

	Driver interface {
		io.Closer
		Name() string
		GetDocument(ctx context.Context, url values.String) (HTMLDocument, error)
	}
)

func WithContext(ctx context.Context, drv Driver) context.Context {
	val := ctx.Value(ctxKey{})
	col, ok := val.(map[string]Driver)

	if !ok {
		col = make(map[string]Driver)
	}

	col[drv.Name()] = drv

	return context.WithValue(
		ctx,
		ctxKey{},
		col,
	)
}

func FromContext(ctx context.Context, name string) (Driver, error) {
	val := ctx.Value(ctxKey{})
	col, ok := val.(map[string]Driver)

	if !ok {
		return nil, core.Error(core.ErrNotFound, name)
	}

	drv, exists := col[name]

	if !exists {
		return nil, core.Error(core.ErrNotFound, name)
	}

	return drv, nil
}

func FromContextAny(ctx context.Context) (Driver, error) {
	val := ctx.Value(ctxKey{})
	col, ok := val.(map[string]Driver)

	if !ok {
		return nil, core.Error(core.ErrNotFound, "html drivers")
	}

	var name string

	for k := range col {
		name = k
		break
	}

	drv := col[name]

	return drv, nil
}
