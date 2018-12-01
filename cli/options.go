package cli

import (
	"context"
	"github.com/MontFerret/ferret/pkg/html"
	"github.com/MontFerret/ferret/pkg/html/dynamic"
	"github.com/MontFerret/ferret/pkg/html/static"
)

type Options struct {
	Cdp       string
	Params    map[string]interface{}
	Proxy     string
	UserAgent string
	ShowTime  bool
}

func (opts Options) WithContext(ctx context.Context) context.Context {
	ctx = html.WithDynamicDriver(
		ctx,
		dynamic.WithCDP(opts.Cdp),
		dynamic.WithProxy(opts.Proxy),
		dynamic.WithUserAgent(opts.UserAgent),
	)

	ctx = html.WithStaticDriver(
		ctx,
		static.WithProxy(opts.Proxy),
		static.WithUserAgent(opts.UserAgent),
	)

	return ctx
}
