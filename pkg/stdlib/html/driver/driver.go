package driver

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/browser"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/http"
)

const Cdp = "cdp"
const Http = "http"

type Driver interface {
	GetDocument(ctx context.Context, url string) (values.HtmlNode, error)
	Close() error
}

func ToContext(ctx context.Context, name string, drv Driver) context.Context {
	return context.WithValue(ctx, name, drv)
}

func FromContext(ctx context.Context, name string) (Driver, error) {
	val := ctx.Value(name)

	drv, ok := val.(Driver)

	if ok {
		return drv, nil
	}

	return nil, core.Error(core.ErrNotFound, fmt.Sprintf("%s driver", name))
}

func WithCdpDriver(ctx context.Context, addr string) context.Context {
	return context.WithValue(ctx, Cdp, browser.NewDriver(addr))
}

func WithHttpDriver(ctx context.Context, opts ...http.Option) context.Context {
	return context.WithValue(ctx, Http, http.NewDriver(opts...))
}
