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

	cdpDriver := cdp.NewDriver(
		cdp.WithAddress(opts.Cdp),
		cdp.WithProxy(opts.Proxy),
		cdp.WithUserAgent(opts.UserAgent),
	)

	ctx = drivers.WithContext(
		ctx,
		cdpDriver,
	)

	return context.WithCancel(ctx)
}
