package cdp

import (
	"context"
	"sync"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/browser"
	"github.com/mafredri/cdp/protocol/target"
	"github.com/mafredri/cdp/rpcc"
	"github.com/mafredri/cdp/session"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
)

const DriverName = "cdp"
const BlankPageURL = "about:blank"

var defaultViewport = &drivers.Viewport{
	Width:  1600,
	Height: 900,
}

type Driver struct {
	mu        sync.Mutex
	dev       *devtool.DevTools
	conn      *rpcc.Conn
	client    *cdp.Client
	session   *session.Manager
	contextID browser.ContextID
	options   *Options
}

func NewDriver(opts ...Option) *Driver {
	drv := new(Driver)
	drv.options = NewOptions(opts)
	drv.dev = devtool.New(drv.options.Address)

	return drv
}

func (drv *Driver) Name() string {
	return drv.options.Name
}

func (drv *Driver) Open(ctx context.Context, params drivers.Params) (drivers.HTMLPage, error) {
	logger := logging.FromContext(ctx)

	conn, err := drv.createConnection(ctx, params.KeepCookies)

	if err != nil {
		logger.Error().
			Err(err).
			Str("driver", drv.options.Name).
			Msg("failed to create a new connection")

		return nil, err
	}

	return LoadHTMLPage(ctx, conn, drv.setDefaultParams(params))
}

func (drv *Driver) Parse(ctx context.Context, params drivers.ParseParams) (drivers.HTMLPage, error) {
	logger := logging.FromContext(ctx)

	conn, err := drv.createConnection(ctx, true)

	if err != nil {
		logger.Error().
			Err(err).
			Str("driver", drv.options.Name).
			Msg("failed to create a new connection")

		return nil, err
	}

	return LoadHTMLPageWithContent(ctx, conn, drv.setDefaultParams(drivers.Params{
		URL:         BlankPageURL,
		UserAgent:   "",
		KeepCookies: params.KeepCookies,
		Cookies:     params.Cookies,
		Headers:     params.Headers,
		Viewport:    params.Viewport,
	}), params.Content)
}

func (drv *Driver) Close() error {
	drv.mu.Lock()
	defer drv.mu.Unlock()

	if drv.session != nil {
		drv.session.Close()

		return drv.conn.Close()
	}

	return nil
}

func (drv *Driver) createConnection(ctx context.Context, keepCookies bool) (*rpcc.Conn, error) {
	err := drv.init(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "initialize driver")
	}

	// Args for a new target belonging to the browser context
	createTargetArgs := target.NewCreateTargetArgs(BlankPageURL)

	if !drv.options.KeepCookies && !keepCookies {
		// Set it to an incognito mode
		createTargetArgs.SetBrowserContextID(drv.contextID)
	}

	// New target
	createTarget, err := drv.client.Target.CreateTarget(ctx, createTargetArgs)

	if err != nil {
		return nil, errors.Wrap(err, "create a browser target")
	}

	// Connect to target using the existing websocket connection.
	conn, err := drv.session.Dial(ctx, createTarget.TargetID)

	if err != nil {
		return nil, errors.Wrap(err, "establish a new connection")
	}

	return conn, nil
}

func (drv *Driver) setDefaultParams(params drivers.Params) drivers.Params {
	if params.Viewport == nil {
		params.Viewport = defaultViewport
	}

	return drivers.SetDefaultParams(drv.options.Options, params)
}

func (drv *Driver) init(ctx context.Context) error {
	drv.mu.Lock()
	defer drv.mu.Unlock()

	if drv.session == nil {
		ver, err := drv.dev.Version(ctx)

		if err != nil {
			return errors.Wrap(err, "failed to initialize driver")
		}

		dialOpts := make([]rpcc.DialOption, 0, 2)

		if drv.options.Connection != nil {
			if drv.options.Connection.BufferSize > 0 {
				dialOpts = append(dialOpts, rpcc.WithWriteBufferSize(drv.options.Connection.BufferSize))
			}

			if drv.options.Connection.Compression {
				dialOpts = append(dialOpts, rpcc.WithCompression())
			}
		}

		bconn, err := rpcc.DialContext(
			ctx,
			ver.WebSocketDebuggerURL,
			dialOpts...,
		)

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

		createCtx, err := bc.Target.CreateBrowserContext(ctx, &target.CreateBrowserContextArgs{})

		if err != nil {
			bconn.Close()
			sess.Close()

			return err
		}

		drv.contextID = createCtx.BrowserContextID
	}

	return nil
}
