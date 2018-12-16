package cli

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/drivers/http"
)

type Options struct {
	Cdp       string
	Params    map[string]interface{}
	Proxy     string
	UserAgent string
	ShowTime  bool
}

func (opts Options) WithContext(ctx context.Context) (context.Context, error) {
	var err error

	ctx = drivers.WithDynamic(
		ctx,
		cdp.NewDriver(
			cdp.WithAddress(opts.Cdp),
			cdp.WithProxy(opts.Proxy),
			cdp.WithUserAgent(opts.UserAgent),
		),
	)

	if err != nil {
		return ctx, err
	}

	ctx = drivers.WithStatic(
		ctx,
		http.NewDriver(
			http.WithProxy(opts.Proxy),
			http.WithUserAgent(opts.UserAgent),
		),
	)

	return ctx, err
}
