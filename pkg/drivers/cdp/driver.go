package cdp

import (
	"context"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/target"
	"github.com/mafredri/cdp/rpcc"
	"github.com/mafredri/cdp/session"
	"github.com/pkg/errors"
	"sync"

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

func (drv *Driver) Open(ctx context.Context, params drivers.Params) (drivers.HTMLPage, error) {
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

	// Args for a new target belonging to the browser context
	createTargetArgs := target.NewCreateTargetArgs(BlankPageURL)

	if !drv.options.KeepCookies && !params.KeepCookies {
		// Set it to an incognito mode
		createTargetArgs.SetBrowserContextID(drv.contextID)
	}

	// New target
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

	if params.UserAgent == "" {
		params.UserAgent = drv.options.UserAgent
	}

	if params.Viewport == nil {
		params.Viewport = defaultViewport
	}

	if drv.options.Headers != nil && params.Headers == nil {
		params.Headers = make(drivers.HTTPHeaders)
	}

	// set default headers
	for k, v := range drv.options.Headers {
		_, exists := params.Headers[k]

		// do not override user's set values
		if !exists {
			params.Headers[k] = v
		}
	}

	if drv.options.Cookies != nil && params.Cookies == nil {
		params.Cookies = make(drivers.HTTPCookies)
	}

	// set default cookies
	for k, v := range drv.options.Cookies {
		_, exists := params.Cookies[k]

		// do not override user's set values
		if !exists {
			params.Cookies[k] = v
		}
	}

	return LoadHTMLPage(ctx, conn, params)
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

func (drv *Driver) init(ctx context.Context) error {
	drv.mu.Lock()
	defer drv.mu.Unlock()

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
