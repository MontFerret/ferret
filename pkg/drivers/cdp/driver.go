package cdp

import (
	"context"
	"encoding/json"
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

	url := params.Url

	if url == "" {
		url = BlankPageURL
	}

	// Create a new target belonging to the browser context
	createTargetArgs := target.NewCreateTargetArgs(url)

	if drv.options.KeepCookies == false && params.KeepCookies == false {
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
			if params.Cookies != nil {
				cookies := make([]network.CookieParam, 0, len(params.Cookies))

				for i, c := range params.Cookies {
					cookies[i] = fromDriverCookie(c)

					logger.
						Debug().
						Timestamp().
						Str("cookie", c.Name).
						Msg("set cookie")
				}

				return client.Network.SetCookies(
					ctx,
					network.NewSetCookiesArgs(cookies),
				)
			}

			return nil
		},

		func() error {
			if params.Header != nil {
				j, err := json.Marshal(params.Header)

				if err != nil {
					return err
				}

				for k := range params.Header {
					logger.
						Debug().
						Timestamp().
						Str("header", k).
						Msg("set header")
				}

				return client.Network.SetExtraHTTPHeaders(
					ctx,
					network.NewSetExtraHTTPHeadersArgs(network.Headers(j)),
				)
			}

			return nil
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
