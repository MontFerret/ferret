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
	"github.com/MontFerret/ferret/pkg/drivers/cdp/input"
	net "github.com/MontFerret/ferret/pkg/drivers/cdp/network"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/utils"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/events"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type (
	HTMLPageEvent string

	HTMLPage struct {
		mu      sync.Mutex
		closed  values.Boolean
		logger  zerolog.Logger
		conn    *rpcc.Conn
		client  *cdp.Client
		network *net.Manager
		dom     *dom.Manager
	}
)

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

	closers := make([]io.Closer, 0, 4)

	defer func() {
		if err != nil {
			if err := client.Page.Close(context.Background()); err != nil {
				logger.Error().Err(err)
			}

			if err := conn.Close(); err != nil {
				logger.Error().Err(err)
			}

			common.CloseAll(logger, closers, "failed to close a Page resource")
		}
	}()

	netOpts := net.Options{
		Headers: params.Headers,
	}

	if params.Cookies != nil && params.Cookies.Length() > 0 {
		netOpts.Cookies = make(map[string]*drivers.HTTPCookies)
		netOpts.Cookies[params.URL] = params.Cookies
	}

	if params.Ignore != nil && len(params.Ignore.Resources) > 0 {
		netOpts.Filter = &net.Filter{
			Patterns: params.Ignore.Resources,
		}
	}

	netManager, err := net.New(
		logger,
		client,
		netOpts,
	)

	if err != nil {
		return nil, err
	}

	mouse := input.NewMouse(client)
	keyboard := input.NewKeyboard(client)

	domManager, err := dom.New(
		logger,
		client,
		mouse,
		keyboard,
	)

	if err != nil {
		return nil, err
	}

	p = NewHTMLPage(
		logger,
		conn,
		client,
		netManager,
		domManager,
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
	logger zerolog.Logger,
	conn *rpcc.Conn,
	client *cdp.Client,
	netManager *net.Manager,
	domManager *dom.Manager,
) *HTMLPage {
	p := new(HTMLPage)
	p.closed = values.False
	p.logger = logging.WithName(logger.With(), "cdp_page").Logger()
	p.conn = conn
	p.client = client
	p.network = netManager
	p.dom = domManager

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

func (p *HTMLPage) GetIn(ctx context.Context, path []core.Value) (core.Value, core.PathError) {
	return common.GetInPage(ctx, path, p)
}

func (p *HTMLPage) SetIn(ctx context.Context, path []core.Value, value core.Value) core.PathError {
	return common.SetInPage(ctx, path, p, value)
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

	var url string
	frame := p.dom.GetMainFrame()

	if frame != nil {
		url = frame.GetURL().String()
	}

	p.closed = values.True

	err := p.dom.Close()

	if err != nil {
		p.logger.Warn().
			Str("url", url).
			Err(err).
			Msg("failed to close dom manager")
	}

	err = p.network.Close()

	if err != nil {
		p.logger.Warn().
			Str("url", url).
			Err(err).
			Msg("failed to close network manager")
	}

	err = p.client.Page.Close(context.Background())

	if err != nil {
		p.logger.Warn().
			Str("url", url).
			Err(err).
			Msg("failed to close browser page")
	}

	// Ignore errors from the connection object
	p.conn.Close()

	return nil
}

func (p *HTMLPage) IsClosed() values.Boolean {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.closed
}

func (p *HTMLPage) GetURL() values.String {
	res, err := p.getCurrentDocument().Eval().EvalValue(context.Background(), templates.GetURL())

	if err == nil {
		return values.ToString(res)
	}

	p.logger.Warn().
		Err(err).
		Msg("failed to retrieve URL")

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

func (p *HTMLPage) GetCookies(ctx context.Context) (*drivers.HTTPCookies, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.network.GetCookies(ctx)
}

func (p *HTMLPage) SetCookies(ctx context.Context, cookies *drivers.HTTPCookies) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.network.SetCookies(ctx, p.getCurrentDocument().GetURL().String(), cookies)
}

func (p *HTMLPage) DeleteCookies(ctx context.Context, cookies *drivers.HTTPCookies) error {
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

	clientWidth, clientHeight := utils.GetLayoutViewportWH(metrics)

	if params.Width <= 0 {
		params.Width = values.Float(clientWidth) - params.X
	}

	if params.Height <= 0 {
		params.Height = values.Float(clientHeight) - params.Y
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
	p.mu.Lock()
	defer p.mu.Unlock()

	pattern, err := p.urlToRegexp(targetURL)

	if err != nil {
		return err
	}

	if err := p.network.WaitForNavigation(ctx, net.WaitEventOptions{URL: pattern}); err != nil {
		return err
	}

	return p.reloadMainFrame(ctx)
}

func (p *HTMLPage) WaitForFrameNavigation(ctx context.Context, frame drivers.HTMLDocument, targetURL values.String) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	current := p.dom.GetMainFrame()
	doc, ok := frame.(*dom.HTMLDocument)

	if !ok {
		return errors.New("invalid frame type")
	}

	pattern, err := p.urlToRegexp(targetURL)

	if err != nil {
		return err
	}

	frameID := doc.Frame().Frame.ID
	isMain := current.Frame().Frame.ID == frameID

	opts := net.WaitEventOptions{
		URL: pattern,
	}

	// if it's the current document
	if !isMain {
		opts.FrameID = frameID
	}

	if err = p.network.WaitForNavigation(ctx, opts); err != nil {
		return err
	}

	return p.reloadMainFrame(ctx)
}

func (p *HTMLPage) Subscribe(ctx context.Context, subscription events.Subscription) (events.Stream, error) {
	switch subscription.EventName {
	case drivers.NavigationEvent:
		p.mu.Lock()
		defer p.mu.Unlock()

		stream, err := p.network.OnNavigation(ctx)

		if err != nil {
			return nil, err
		}

		return newPageNavigationEventStream(stream, func(ctx context.Context) error {
			return p.reloadMainFrame(ctx)
		}), nil
	case drivers.RequestEvent:
		return p.network.OnRequest(ctx)
	case drivers.ResponseEvent:
		return p.network.OnResponse(ctx)
	default:
		return nil, core.Errorf(core.ErrInvalidOperation, "unknown event name: %s", subscription.EventName)
	}
}

func (p *HTMLPage) urlToRegexp(targetURL values.String) (*regexp.Regexp, error) {
	if targetURL == "" {
		return nil, nil
	}

	r, err := regexp.Compile(targetURL.String())

	if err != nil {
		return nil, errors.Wrap(err, "invalid URL pattern")
	}

	return r, nil
}

func (p *HTMLPage) reloadMainFrame(ctx context.Context) error {
	prev := p.dom.GetMainFrame()

	if prev != nil {
		if err := p.dom.RemoveFrameRecursively(prev.Frame().Frame.ID); err != nil {
			p.logger.Error().Err(err).Msg("failed to remove main frame")
		}
	}

	next, err := p.dom.LoadRootDocument(ctx)

	if err != nil {
		p.logger.Error().Err(err).Msg("failed to load a new root document")

		return err
	}

	p.dom.SetMainFrame(next)

	return nil
}

func (p *HTMLPage) loadMainFrame(ctx context.Context) error {
	next, err := p.dom.LoadRootDocument(ctx)

	if err != nil {
		return err
	}

	p.dom.SetMainFrame(next)

	return nil
}

func (p *HTMLPage) getCurrentDocument() *dom.HTMLDocument {
	return p.dom.GetMainFrame()
}
