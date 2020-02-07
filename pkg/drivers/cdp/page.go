package cdp

import (
	"context"
	"hash/fnv"
	"io"
	"regexp"
	"sync"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/dom"
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
}

func LoadHTMLPage(
	ctx context.Context,
	conn *rpcc.Conn,
	params drivers.Params,
) (p *HTMLPage, err error) {
	logger := logging.FromContext(ctx)

	if conn == nil {
		return nil, core.Error(core.ErrMissedArgument, "connection")
	}

	client := cdp.NewClient(conn)

	if err := enableFeatures(ctx, client, params); err != nil {
		return nil, err
	}

	closers := make([]io.Closer, 0, 2)

	defer func() {
		if err != nil {
			common.CloseAll(logger, closers, "failed to close a Page resource")
		}
	}()

	eventLoop := events.NewLoop()
	closers = append(closers, eventLoop)

	netManager, err := net.New(logger, client, eventLoop)

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

	mouse := input.NewMouse(client)
	keyboard := input.NewKeyboard(client)

	domManager, err := dom.New(
		logger,
		client,
		eventLoop,
		mouse,
		keyboard,
	)

	if err != nil {
		return nil, err
	}

	closers = append(closers, domManager)

	p = NewHTMLPage(
		logger,
		conn,
		client,
		eventLoop,
		netManager,
		domManager,
		mouse,
		keyboard,
	)

	if params.URL != BlankPageURL && params.URL != "" {
		err = p.Navigate(ctx, values.NewString(params.URL))
	} else {
		err = p.loadMainFrame(ctx)
	}

	if err != nil {
		return p, err
	}

	return p, nil
}

func LoadHTMLPageWithContent(
	ctx context.Context,
	conn *rpcc.Conn,
	params drivers.Params,
	content []byte,
) (p *HTMLPage, err error) {
	logger := logging.FromContext(ctx)
	p, err = LoadHTMLPage(ctx, conn, params)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if e := p.Close(); e != nil {
				logger.Error().Err(e).Msg("failed to close page")
			}
		}
	}()

	frameID := p.getCurrentDocument().Frame().Frame.ID
	err = p.client.Page.SetDocumentContent(ctx, page.NewSetDocumentContentArgs(frameID, string(content)))

	if err != nil {
		return nil, errors.Wrap(err, "set document content")
	}

	// Remove prev frames (from a blank page)
	prev := p.dom.GetMainFrame()
	err = p.dom.RemoveFrameRecursively(prev.Frame().Frame.ID)

	if err != nil {
		return nil, err
	}

	err = p.loadMainFrame(ctx)

	if err != nil {
		return nil, err
	}

	return p, nil
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

	eventLoop.AddListener(events.Error, events.Always(p.handleError))

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

	doc := p.getCurrentDocument()
	err := p.events.Stop().Close()

	if err != nil {
		p.logger.Warn().
			Timestamp().
			Str("url", doc.GetURL().String()).
			Err(err).
			Msg("failed to stop event loop")
	}

	err = p.dom.Close()

	if err != nil {
		p.logger.Warn().
			Timestamp().
			Str("url", doc.GetURL().String()).
			Err(err).
			Msg("failed to close dom manager")
	}

	err = p.network.Close()

	if err != nil {
		p.logger.Warn().
			Timestamp().
			Str("url", doc.GetURL().String()).
			Err(err).
			Msg("failed to close network manager")
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
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.dom.GetFrameNodes(ctx)
}

func (p *HTMLPage) GetFrame(ctx context.Context, idx values.Int) (core.Value, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	frames, err := p.dom.GetFrameNodes(ctx)

	if err != nil {
		return values.None, err
	}

	return frames.Get(idx), nil
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

func (p *HTMLPage) GetResponse(ctx context.Context) (drivers.HTTPResponse, error) {
	doc := p.getCurrentDocument()

	if doc == nil {
		return drivers.HTTPResponse{}, nil
	}

	return p.network.GetResponse(ctx, doc.Frame().Frame.ID)
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

	if err := p.network.Navigate(ctx, url); err != nil {
		return err
	}

	return p.reloadMainFrame(ctx)
}

func (p *HTMLPage) NavigateBack(ctx context.Context, skip values.Int) (values.Boolean, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	ret, err := p.network.NavigateBack(ctx, skip)

	if err != nil {
		return values.False, err
	}

	return ret, p.reloadMainFrame(ctx)
}

func (p *HTMLPage) NavigateForward(ctx context.Context, skip values.Int) (values.Boolean, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	ret, err := p.network.NavigateForward(ctx, skip)

	if err != nil {
		return values.False, err
	}

	return ret, p.reloadMainFrame(ctx)
}

func (p *HTMLPage) WaitForNavigation(ctx context.Context, targetURL values.String) error {
	var pattern *regexp.Regexp

	if targetURL != "" {
		r, err := regexp.Compile(targetURL.String())

		if err != nil {
			return errors.Wrap(err, "invalid URL pattern")
		}

		pattern = r
	}

	if err := p.network.WaitForNavigation(ctx, pattern); err != nil {
		return err
	}

	return p.reloadMainFrame(ctx)
}

func (p *HTMLPage) reloadMainFrame(ctx context.Context) error {
	if err := p.dom.WaitForDOMReady(ctx); err != nil {
		return err
	}

	prev := p.dom.GetMainFrame()

	next, err := dom.LoadRootHTMLDocument(
		ctx,
		p.logger,
		p.client,
		p.dom,
		p.mouse,
		p.keyboard,
	)

	if err != nil {
		return err
	}

	if prev != nil {
		if err := p.dom.RemoveFrameRecursively(prev.Frame().Frame.ID); err != nil {
			p.logger.Error().Err(err).Msg("failed to remove main frame")
		}
	}

	p.dom.SetMainFrame(next)

	return nil
}

func (p *HTMLPage) loadMainFrame(ctx context.Context) error {
	next, err := dom.LoadRootHTMLDocument(
		ctx,
		p.logger,
		p.client,
		p.dom,
		p.mouse,
		p.keyboard,
	)

	if err != nil {
		return err
	}

	p.dom.SetMainFrame(next)

	return nil
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
	return p.dom.GetMainFrame()
}
