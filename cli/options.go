package cli

import (
	"context"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/drivers/http"
)

type Options struct {
	Cdp         string
	Params      map[string]interface{}
	Proxy       string
	UserAgent   string
	ShowTime    bool
	KeepCookies bool
}

func (opts Options) WithContext(ctx context.Context) (context.Context, context.CancelFunc) {
	httpDriver := http.NewDriver(
		http.WithProxy(opts.Proxy),
		http.WithUserAgent(opts.UserAgent),
	)

	ctx = drivers.WithContext(
		ctx,
		httpDriver,
		drivers.AsDefault(),
	)

	cdpOpts := []cdp.Option{
		cdp.WithAddress(opts.Cdp),
		cdp.WithProxy(opts.Proxy),
		cdp.WithUserAgent(opts.UserAgent),
	}

	if opts.KeepCookies {
		cdpOpts = append(cdpOpts, cdp.WithKeepCookies())
	}

	cdpDriver := cdp.NewDriver(cdpOpts...)

	ctx = drivers.WithContext(
		ctx,
		cdpDriver,
	)

	return context.WithCancel(ctx)
}
