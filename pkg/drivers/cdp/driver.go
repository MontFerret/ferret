package cdp

import (
	"context"
	"sync"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/emulation"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/protocol/target"
	"github.com/mafredri/cdp/rpcc"
	"github.com/mafredri/cdp/session"
	"github.com/pkg/errors"
)

const DriverName = "cdp"

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
	drv.dev = devtool.New(drv.options.Address)

	return drv
}

func (drv *Driver) Name() string {
	return drv.options.Name
}

func (drv *Driver) LoadDocument(ctx context.Context, params drivers.LoadDocumentParams) (drivers.HTMLDocument, error) {
	logger := logging.FromContext(ctx)

	err := drv.init(ctx)

	if err != nil {
		logger.
			Error().
			Timestamp().
			Err(err).
			Str("driver", drv.options.Name).
			Msg("failed to initialize the driver")

		return nil, err
	}

	url := params.URL

	if url == "" {
		url = BlankPageURL
	}

	// Create a new target belonging to the browser context
	createTargetArgs := target.NewCreateTargetArgs(url)

	if !drv.options.KeepCookies && !params.KeepCookies {
		// Set it to an incognito mode
		createTargetArgs.SetBrowserContextID(drv.contextID)
	}

	createTarget, err := drv.client.Target.CreateTarget(ctx, createTargetArgs)

	if err != nil {
		logger.
			Error().
			Timestamp().
			Err(err).
			Str("driver", drv.options.Name).
			Msg("failed to create a browser target")

		return nil, err
	}

	// Connect to target using the existing websocket connection.
	conn, err := drv.session.Dial(ctx, createTarget.TargetID)

	if err != nil {
		logger.
			Error().
			Timestamp().
			Err(err).
			Str("driver", drv.options.Name).
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
			var ua string

			if params.UserAgent != "" {
				ua = common.GetUserAgent(params.UserAgent)
			} else {
				ua = common.GetUserAgent(drv.options.UserAgent)
			}

			logger.
				Debug().
				Timestamp().
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

		func() error {
			return client.Network.Enable(ctx, network.NewEnableArgs())
		},
	)

	if err != nil {
		return nil, err
	}

	return LoadHTMLDocument(ctx, conn, client, params)
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

		drv.conn = bconn
		drv.client = bc
		drv.session = sess

		if drv.options.KeepCookies {
			return nil
		}

		createCtx, err := bc.Target.CreateBrowserContext(ctx)

		if err != nil {
			bconn.Close()
			sess.Close()

			return err
		}

		drv.contextID = createCtx.BrowserContextID
	}

	return nil
}
