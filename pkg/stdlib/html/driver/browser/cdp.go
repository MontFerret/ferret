package browser

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/rpcc"
)

type (
	CdpDriver struct {
		address string
	}

	CdpConnection struct {
		target *devtool.Target
		core   *rpcc.Conn
	}
)

func NewDriver(conn string) *CdpDriver {
	return &CdpDriver{
		address: conn,
	}
}

func (drv *CdpDriver) GetDocument(ctx context.Context, url string) (values.HtmlNode, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	devt := devtool.New(drv.address)

	target, err := devt.CreateURL(ctx, url)

	if err != nil {
		return nil, err
	}

	conn, err := rpcc.DialContext(ctx, target.WebSocketDebuggerURL)

	if err != nil {
		return nil, err
	}

	return NewHtmlDocument(ctx, conn, url)
}

func (drv *CdpDriver) Close() error {
	// TODO: Do we need this method?
	return nil
}
