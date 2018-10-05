package dynamic

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/target"
	"github.com/mafredri/cdp/rpcc"
	"github.com/mafredri/cdp/session"
	"github.com/pkg/errors"
	"sync"
)

type Driver struct {
	sync.Mutex
	dev       *devtool.DevTools
	conn      *rpcc.Conn
	client    *cdp.Client
	session   *session.Manager
	contextID target.BrowserContextID
}

func NewDriver(address string) *Driver {
	drv := new(Driver)
	drv.dev = devtool.New(address)

	return drv
}

func (drv *Driver) GetDocument(ctx context.Context, url string) (values.HTMLNode, error) {
	err := drv.init(ctx)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	if url == "" {
		url = BlankPageUrl
	}

	// Create a new target belonging to the browser context, similar
	// to opening a new tab in an incognito window.
	createTargetArgs := target.NewCreateTargetArgs(url).SetBrowserContextID(drv.contextID)
	createTarget, err := drv.client.Target.CreateTarget(ctx, createTargetArgs)

	if err != nil {
		return nil, err
	}

	// Connect to target using the existing websocket connection.
	conn, err := drv.session.Dial(ctx, createTarget.TargetID)

	if err != nil {
		return nil, err
	}

	return LoadHTMLDocument(ctx, conn, url)
}

func (drv *Driver) Close() error {
	drv.Lock()
	defer drv.Unlock()

	if drv.session != nil {
		drv.session.Close()

		return drv.conn.Close()
	}

	return nil
}

func (drv *Driver) init(ctx context.Context) error {
	drv.Lock()
	defer drv.Unlock()

	if drv.session == nil {
		ver, err := drv.dev.Version(ctx)

		if err != nil {
			return errors.Wrap(err, "failed to initialize driver")
		}

		bconn, err := rpcc.DialContext(ctx, ver.WebSocketDebuggerURL)

		if err != nil {
			return errors.Wrap(err, "failed to initialize driver")
		}

		bc := cdp.NewClient(bconn)

		sess, err := session.NewManager(bc)

		if err != nil {
			bconn.Close()

			return errors.Wrap(err, "failed to initialize driver")
		}

		createCtx, err := bc.Target.CreateBrowserContext(ctx)

		if err != nil {
			bconn.Close()
			sess.Close()

			return err
		}

		drv.conn = bconn
		drv.client = bc
		drv.session = sess
		drv.contextID = createCtx.BrowserContextID
	}

	return nil
}
