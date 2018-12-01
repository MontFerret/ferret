package html

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/html/dynamic"
	"github.com/MontFerret/ferret/pkg/html/static"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type DriverName string

const (
	Dynamic DriverName = "dynamic"
	Static  DriverName = "static"
)

type Driver interface {
	GetDocument(ctx context.Context, url values.String) (values.HTMLNode, error)
	Close() error
}

func FromContext(ctx context.Context, name DriverName) (Driver, error) {
	switch name {
	case Dynamic:
		return dynamic.FromContext(ctx)
	case Static:
		return static.FromContext(ctx)
	default:
		return nil, core.Error(core.ErrInvalidArgument, fmt.Sprintf("%s driver", name))
	}
}

func WithDynamicDriver(ctx context.Context, opts ...dynamic.Option) context.Context {
	return dynamic.WithContext(ctx, dynamic.NewDriver(opts...))
}

func WithStaticDriver(ctx context.Context, opts ...static.Option) context.Context {
	return static.WithContext(ctx, static.NewDriver(opts...))
}
