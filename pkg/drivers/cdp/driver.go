package cdp

import (
	"context"
	"sync"

	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/emulation"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/target"
	"github.com/mafredri/cdp/rpcc"
	"github.com/mafredri/cdp/session"
	"github.com/pkg/errors"
)

type Driver struct {
	sync.Mutex
	dev       *devtool.DevTools
	conn      *rpcc.Conn
	client    *cdp.Client
	session   *session.Manager
	contextID target.BrowserContextID
	options   *Options
}

func NewDriver(opts ...Option) *Driver {
	drv := new(Driver)
	drv.options = newOptions(opts)
	drv.dev = devtool.New(drv.options.address)

	return drv
}

func (drv *Driver) GetDocument(ctx context.Context, targetURL values.String) (values.DHTMLDocument, error) {
	logger := logging.FromContext(ctx)

	err := drv.init(ctx)

	if err != nil {
		logger.
			Error().
			Err(err).
			Str("driver", "dynamic").
			Msg("failed to initialize the driver")

		return nil, err
	}

	url := targetURL.String()

	if url == "" {
		url = BlankPageURL
	}

	// Create a new target belonging to the browser context, similar
	// to opening a new tab in an incognito window.
	createTargetArgs := target.NewCreateTargetArgs(url).SetBrowserContextID(drv.contextID)
	createTarget, err := drv.client.Target.CreateTarget(ctx, createTargetArgs)

	if err != nil {
		logger.
			Error().
			Err(err).
			Str("driver", "dynamic").
			Msg("failed to create a browser target")

		return nil, err
	}

	// Connect to target using the existing websocket connection.
	conn, err := drv.session.Dial(ctx, createTarget.TargetID)

	if err != nil {
		logger.
			Error().
			Err(err).
			Str("driver", "dynamic").
			Msg("failed to establish a connection")

		return nil, err
	}

	client := cdp.NewClient(conn)

	err = runBatch(
		func() error {
			return client.Page.Enable(ctx)
		},

		func() error {
			return client.Page.SetLifecycleEventsEnabled(
				ctx,
				page.NewSetLifecycleEventsEnabledArgs(true),
			)
		},

		func() error {
			return client.DOM.Enable(ctx)
		},

		func() error {
			return client.Runtime.Enable(ctx)
		},

		func() error {
			ua := common.GetUserAgent(drv.options.userAgent)

			logger.
				Debug().
				Str("user-agent", ua).
				Msg("using User-Agent")

			// do not use custom user agent
			if ua == "" {
				return nil
			}

			return client.Emulation.SetUserAgentOverride(
				ctx,
				emulation.NewSetUserAgentOverrideArgs(ua),
			)
		},
	)

	if err != nil {
		return nil, err
	}

	return LoadHTMLDocument(ctx, conn, client, url)
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
