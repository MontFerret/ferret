package cdp

import (
	"context"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/dom"
	"hash/fnv"
	"sync"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/emulation"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/input"
	net "github.com/MontFerret/ferret/pkg/drivers/cdp/network"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type HTMLPage struct {
	mu       sync.Mutex
	closed   values.Boolean
	logger   *zerolog.Logger
	conn     *rpcc.Conn
	client   *cdp.Client
	events   *events.Loop
	network  *net.Manager
	dom      *dom.Manager
	mouse    *input.Mouse
	keyboard *input.Keyboard
	document *common.AtomicValue
	frames   *common.LazyValue
}

func handleLoadError(logger *zerolog.Logger, client *cdp.Client) {
	err := client.Page.Close(context.Background())

	if err != nil {
		logger.Warn().Timestamp().Err(err).Msg("failed to close document on load error")
	}
}

func LoadHTMLPage(
	ctx context.Context,
	conn *rpcc.Conn,
	params drivers.Params,
) (*HTMLPage, error) {
	logger := logging.FromContext(ctx)

	if conn == nil {
		return nil, core.Error(core.ErrMissedArgument, "connection")
	}

	client := cdp.NewClient(conn)

	if err := client.Page.Enable(ctx); err != nil {
		return nil, err
	}

	err := runBatch(
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
			ua := common.GetUserAgent(params.UserAgent)

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

		func() error {
			return client.Page.SetBypassCSP(ctx, page.NewSetBypassCSPArgs(true))
		},

		func() error {
			if params.Viewport == nil {
				return nil
			}

			orientation := emulation.ScreenOrientation{}

			if !params.Viewport.Landscape {
				orientation.Type = "portraitPrimary"
				orientation.Angle = 0
			} else {
				orientation.Type = "landscapePrimary"
				orientation.Angle = 90
			}

			scaleFactor := params.Viewport.ScaleFactor

			if scaleFactor <= 0 {
				scaleFactor = 1
			}

			deviceArgs := emulation.NewSetDeviceMetricsOverrideArgs(
				params.Viewport.Width,
				params.Viewport.Height,
				scaleFactor,
				params.Viewport.Mobile,
			).SetScreenOrientation(orientation)

			return client.Emulation.SetDeviceMetricsOverride(
				ctx,
				deviceArgs,
			)
		},
	)

	if err != nil {
		return nil, err
	}

	eventLoop := events.NewLoop()
	netManager, err := net.New(ctx, logger, client, eventLoop)

	if err != nil {
		return nil, err
	}

	err = netManager.SetCookies(ctx, params.URL, params.Cookies)

	if err != nil {
		return nil, err
	}

	err = netManager.SetHeaders(ctx, params.Headers)

	if err != nil {
		return nil, err
	}

	eventLoop.Start()

	if params.URL != BlankPageURL && params.URL != "" {
		repl, err := client.Page.Navigate(ctx, page.NewNavigateArgs(params.URL))

		if err != nil {
			handleLoadError(logger, client)

			return nil, errors.Wrap(err, "failed to load the page")
		}

		if repl.ErrorText != nil {
			handleLoadError(logger, client)

			return nil, errors.Wrapf(errors.New(*repl.ErrorText), "failed to load the page: %s", params.URL)
		}

		err = netManager.WaitForNavigation(ctx, values.NewString(params.URL))

		if err != nil {
			handleLoadError(logger, client)

			return nil, errors.Wrap(err, "failed to load the page")
		}
	}

	if err != nil {
		handleLoadError(logger, client)
		return nil, errors.Wrap(err, "failed to create event events")
	}

	mouse := input.NewMouse(client)
	keyboard := input.NewKeyboard(client)

	domManager, err := dom.NewManager(ctx, client, eventLoop)

	if err != nil {
		eventLoop.Stop().Close()

		return nil, errors.Wrap(err, "failed to initialize dom manager")
	}

	doc, err := dom.LoadRootHTMLDocument(
		ctx,
		logger,
		client,
		netManager,
		domManager,
		mouse,
		keyboard,
	)

	if err != nil {
		eventLoop.Stop().Close()
		handleLoadError(logger, client)

		return nil, errors.Wrap(err, "failed to load root element")
	}

	return NewHTMLPage(
		logger,
		conn,
		client,
		eventLoop,
		netManager,
		domManager,
		mouse,
		keyboard,
		doc,
	), nil
}

func NewHTMLPage(
	logger *zerolog.Logger,
	conn *rpcc.Conn,
	client *cdp.Client,
	eventLoop *events.Loop,
	netManager *net.Manager,
	domManager *dom.Manager,
	mouse *input.Mouse,
	keyboard *input.Keyboard,
	document *dom.HTMLDocument,
) *HTMLPage {
	p := new(HTMLPage)
	p.closed = values.False
	p.logger = logger
	p.conn = conn
	p.client = client
	p.events = eventLoop
	p.network = netManager
	p.dom = domManager
	p.mouse = mouse
	p.keyboard = keyboard
	p.document = common.NewAtomicValue(document)
	p.frames = common.NewLazyValue(p.unfoldFrames)

	netManager.AddFrameLoadedListener(p.handlePageLoad)
	eventLoop.AddListener(events.Error, p.handleError)

	return p
}

func (p *HTMLPage) MarshalJSON() ([]byte, error) {
	return p.getCurrentDocument().MarshalJSON()
}

func (p *HTMLPage) Type() core.Type {
	return drivers.HTMLPageType
}

func (p *HTMLPage) String() string {
	return p.getCurrentDocument().GetURL().String()
}

func (p *HTMLPage) Compare(other core.Value) int64 {
	tc := drivers.Compare(p.Type(), other.Type())

	if tc != 0 {
		return tc
	}

	cdpPage, ok := other.(*HTMLPage)

	if !ok {
		return 1
	}

	return p.getCurrentDocument().GetURL().Compare(cdpPage.GetURL())
}

func (p *HTMLPage) Unwrap() interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p
}

func (p *HTMLPage) Hash() uint64 {
	h := fnv.New64a()

	h.Write([]byte("CDP"))
	h.Write([]byte(p.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(p.getCurrentDocument().GetURL()))

	return h.Sum64()
}

func (p *HTMLPage) Copy() core.Value {
	return values.None
}

func (p *HTMLPage) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	return common.GetInPage(ctx, p, path)
}

func (p *HTMLPage) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	return common.SetInPage(ctx, p, path, value)
}

func (p *HTMLPage) Iterate(ctx context.Context) (core.Iterator, error) {
	return p.getCurrentDocument().Iterate(ctx)
}

func (p *HTMLPage) Length() values.Int {
	return p.getCurrentDocument().Length()
}

func (p *HTMLPage) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.closed = values.True

	err := p.events.Stop().Close()
	doc := p.getCurrentDocument()

	if err != nil {
		p.logger.Warn().
			Timestamp().
			Str("url", doc.GetURL().String()).
			Err(err).
			Msg("failed to stop event loop")
	}

	if err != nil {
		p.logger.Warn().
			Timestamp().
			Str("url", doc.GetURL().String()).
			Err(err).
			Msg("failed to close event loop")
	}

	err = doc.Close()

	if err != nil {
		p.logger.Warn().
			Timestamp().
			Str("url", doc.GetURL().String()).
			Err(err).
			Msg("failed to close root document")
	}

	err = p.client.Page.Close(context.Background())

	if err != nil {
		p.logger.Warn().
			Timestamp().
			Str("url", doc.GetURL().String()).
			Err(err).
			Msg("failed to close browser page")
	}

	return p.conn.Close()
}

func (p *HTMLPage) IsClosed() values.Boolean {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.closed
}

func (p *HTMLPage) GetURL() values.String {
	return p.getCurrentDocument().GetURL()
}

func (p *HTMLPage) GetMainFrame() drivers.HTMLDocument {
	return p.getCurrentDocument()
}

func (p *HTMLPage) GetFrames(ctx context.Context) (*values.Array, error) {
	res, err := p.frames.Read(ctx)

	if err != nil {
		return nil, err
	}

	return res.(*values.Array).Clone().(*values.Array), nil
}

func (p *HTMLPage) GetFrame(ctx context.Context, idx values.Int) (core.Value, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	res, err := p.frames.Read(ctx)

	if err != nil {
		return nil, err
	}

	return res.(*values.Array).Get(idx), nil
}

func (p *HTMLPage) GetCookies(ctx context.Context) (drivers.HTTPCookies, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.network.GetCookies(ctx)
}

func (p *HTMLPage) SetCookies(ctx context.Context, cookies drivers.HTTPCookies) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.network.SetCookies(ctx, p.getCurrentDocument().GetURL().String(), cookies)
}

func (p *HTMLPage) DeleteCookies(ctx context.Context, cookies drivers.HTTPCookies) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.network.DeleteCookies(ctx, p.getCurrentDocument().GetURL().String(), cookies)
}

func (p *HTMLPage) GetResponse(_ context.Context) (*drivers.HTTPResponse, error) {
	return nil, core.ErrNotSupported
}

func (p *HTMLPage) PrintToPDF(ctx context.Context, params drivers.PDFParams) (values.Binary, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	args := page.NewPrintToPDFArgs()
	args.
		SetLandscape(bool(params.Landscape)).
		SetDisplayHeaderFooter(bool(params.DisplayHeaderFooter)).
		SetPrintBackground(bool(params.PrintBackground)).
		SetIgnoreInvalidPageRanges(bool(params.IgnoreInvalidPageRanges)).
		SetPreferCSSPageSize(bool(params.PreferCSSPageSize))

	if params.Scale > 0 {
		args.SetScale(float64(params.Scale))
	}

	if params.PaperWidth > 0 {
		args.SetPaperWidth(float64(params.PaperWidth))
	}

	if params.PaperHeight > 0 {
		args.SetPaperHeight(float64(params.PaperHeight))
	}

	if params.MarginTop > 0 {
		args.SetMarginTop(float64(params.MarginTop))
	}

	if params.MarginBottom > 0 {
		args.SetMarginBottom(float64(params.MarginBottom))
	}

	if params.MarginRight > 0 {
		args.SetMarginRight(float64(params.MarginRight))
	}

	if params.MarginLeft > 0 {
		args.SetMarginLeft(float64(params.MarginLeft))
	}

	if params.PageRanges != values.EmptyString {
		args.SetPageRanges(string(params.PageRanges))
	}

	if params.HeaderTemplate != values.EmptyString {
		args.SetHeaderTemplate(string(params.HeaderTemplate))
	}

	if params.FooterTemplate != values.EmptyString {
		args.SetFooterTemplate(string(params.FooterTemplate))
	}

	reply, err := p.client.Page.PrintToPDF(ctx, args)

	if err != nil {
		return values.NewBinary([]byte{}), err
	}

	return values.NewBinary(reply.Data), nil
}

func (p *HTMLPage) CaptureScreenshot(ctx context.Context, params drivers.ScreenshotParams) (values.Binary, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	metrics, err := p.client.Page.GetLayoutMetrics(ctx)

	if err != nil {
		return values.NewBinary(nil), err
	}

	if params.Format == drivers.ScreenshotFormatJPEG && params.Quality < 0 && params.Quality > 100 {
		params.Quality = 100
	}

	if params.X < 0 {
		params.X = 0
	}

	if params.Y < 0 {
		params.Y = 0
	}

	if params.Width <= 0 {
		params.Width = values.Float(metrics.LayoutViewport.ClientWidth) - params.X
	}

	if params.Height <= 0 {
		params.Height = values.Float(metrics.LayoutViewport.ClientHeight) - params.Y
	}

	clip := page.Viewport{
		X:      float64(params.X),
		Y:      float64(params.Y),
		Width:  float64(params.Width),
		Height: float64(params.Height),
		Scale:  1.0,
	}

	format := string(params.Format)
	quality := int(params.Quality)
	args := page.CaptureScreenshotArgs{
		Format:  &format,
		Quality: &quality,
		Clip:    &clip,
	}

	reply, err := p.client.Page.CaptureScreenshot(ctx, &args)

	if err != nil {
		return values.NewBinary([]byte{}), err
	}

	return values.NewBinary(reply.Data), nil
}

func (p *HTMLPage) Navigate(ctx context.Context, url values.String) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.network.Navigate(ctx, url)
}

func (p *HTMLPage) NavigateBack(ctx context.Context, skip values.Int) (values.Boolean, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.network.NavigateBack(ctx, skip)
}

func (p *HTMLPage) NavigateForward(ctx context.Context, skip values.Int) (values.Boolean, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.network.NavigateForward(ctx, skip)
}

func (p *HTMLPage) WaitForNavigation(ctx context.Context, params drivers.NavigationParams) error {
	return p.network.WaitForNavigation(ctx, params.TargetURL)
}

func (p *HTMLPage) handlePageLoad(ctx context.Context, _ page.Frame) {
	err := p.document.Write(func(current core.Value) (core.Value, error) {
		nextDoc, err := dom.LoadRootHTMLDocument(
			ctx,
			p.logger,
			p.client,
			p.network,
			p.dom,
			p.mouse,
			p.keyboard,
		)

		if err != nil {
			return values.None, err
		}

		// close the prev document
		currentDoc := current.(*dom.HTMLDocument)
		err = currentDoc.Close()

		if err != nil {
			p.logger.Warn().
				Timestamp().
				Err(err).
				Msgf("failed to close root document: %s", currentDoc.GetURL())
		}

		// reset all loaded frames
		p.frames.Reset()

		return nextDoc, nil
	})

	if err != nil {
		p.logger.Warn().
			Timestamp().
			Err(err).
			Msg("failed to load new root document after page load")
	}
}

func (p *HTMLPage) handleError(_ context.Context, val interface{}) {
	err, ok := val.(error)

	if !ok {
		return
	}

	p.logger.Error().
		Timestamp().
		Err(err).
		Msg("unexpected error")
}

func (p *HTMLPage) getCurrentDocument() *dom.HTMLDocument {
	return p.document.Read().(*dom.HTMLDocument)
}

func (p *HTMLPage) unfoldFrames(ctx context.Context) (core.Value, error) {
	res := values.NewArray(10)

	err := common.CollectFrames(ctx, res, p.getCurrentDocument())

	if err != nil {
		return nil, err
	}

	return res, nil
}
