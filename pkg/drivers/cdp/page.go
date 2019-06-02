package cdp

import (
	"context"
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/mafredri/cdp/protocol/page"
	"sync"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/mafredri/cdp/rpcc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type HTMLPage struct {
	mu       sync.Mutex
	logger   *zerolog.Logger
	conn     *rpcc.Conn
	client   *cdp.Client
	events   *events.EventBroker
	document *HTMLDocument
}

func handleLoadError(logger *zerolog.Logger, client *cdp.Client) {
	err := client.Page.Close(context.Background())

	if err != nil {
		logger.Warn().Timestamp().Err(err).Msg("unabled to close document on load error")
	}
}

func LoadHTMLPage(
	ctx context.Context,
	conn *rpcc.Conn,
	client *cdp.Client,
	params drivers.OpenPageParams,
) (*HTMLPage, error) {
	logger := logging.FromContext(ctx)

	if conn == nil {
		return nil, core.Error(core.ErrMissedArgument, "connection")
	}

	if params.URL == "" {
		return nil, core.Error(core.ErrMissedArgument, "url")
	}

	if params.Cookies != nil {
		cookies := make([]network.CookieParam, 0, len(params.Cookies))

		for _, c := range params.Cookies {
			cookies = append(cookies, fromDriverCookie(params.URL, c))

			logger.
				Debug().
				Timestamp().
				Str("cookie", c.Name).
				Msg("set cookie")
		}

		err := client.Network.SetCookies(
			ctx,
			network.NewSetCookiesArgs(cookies),
		)

		if err != nil {
			return nil, err
		}
	}

	if params.Header != nil {
		j, err := json.Marshal(params.Header)

		if err != nil {
			return nil, err
		}

		for k := range params.Header {
			logger.
				Debug().
				Timestamp().
				Str("header", k).
				Msg("set header")
		}

		err = client.Network.SetExtraHTTPHeaders(
			ctx,
			network.NewSetExtraHTTPHeadersArgs(network.Headers(j)),
		)

		if err != nil {
			return nil, err
		}
	}

	var err error

	if params.URL != BlankPageURL {
		err = events.WaitForLoadEvent(ctx, client)

		if err != nil {
			handleLoadError(logger, client)

			return nil, err
		}
	}

	broker, err := events.CreateEventBroker(client)

	if err != nil {
		handleLoadError(logger, client)
		return nil, errors.Wrap(err, "failed to create event events")
	}

	doc, err := LoadRootHTMLDocument(ctx, logger, client, broker)

	if err != nil {
		broker.Stop()
		broker.Close()
		handleLoadError(logger, client)

		return nil, errors.Wrap(err, "failed to load root element")
	}

	return NewHTMLPage(
		logger,
		conn,
		client,
		broker,
		doc,
	), nil
}

func NewHTMLPage(
	logger *zerolog.Logger,
	conn *rpcc.Conn,
	client *cdp.Client,
	broker *events.EventBroker,
	document *HTMLDocument,
) *HTMLPage {
	p := new(HTMLPage)
	p.logger = logger
	p.conn = conn
	p.client = client
	p.events = broker
	p.document = document

	broker.AddEventListener(events.EventLoad, p.handlePageLoad)
	broker.AddEventListener(events.EventError, p.handleError)

	return p
}

func (p *HTMLPage) MarshalJSON() ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.document.MarshalJSON()
}

func (p *HTMLPage) Type() core.Type {
	return drivers.HTMLPageType
}

func (p *HTMLPage) String() string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.document.GetURL().String()
}

func (p *HTMLPage) Compare(other core.Value) int64 {
	panic("implement me")
}

func (p *HTMLPage) Unwrap() interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.client
}

func (p *HTMLPage) Hash() uint64 {
	panic("implement me")
}

func (p *HTMLPage) Copy() core.Value {
	panic("implement me")
}

func (p *HTMLPage) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	return common.GetInPage(ctx, p, path)
}

func (p *HTMLPage) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	return common.SetInPage(ctx, p, path, value)
}

func (p *HTMLPage) Iterate(ctx context.Context) (core.Iterator, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.document.Iterate(ctx)
}

func (p *HTMLPage) Length() values.Int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.document.Length()
}

func (p *HTMLPage) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	err := p.events.Stop()

	if err != nil {
		p.logger.Warn().
			Timestamp().
			Str("url", p.document.GetURL().String()).
			Err(err).
			Msg("failed to stop event events")
	}

	err = p.events.Close()

	if err != nil {
		p.logger.Warn().
			Timestamp().
			Str("url", p.document.GetURL().String()).
			Err(err).
			Msg("failed to close event events")
	}

	err = p.document.Close()

	if err != nil {
		p.logger.Warn().
			Timestamp().
			Str("url", p.document.GetURL().String()).
			Err(err).
			Msg("failed to close root document")
	}

	err = p.client.Page.Close(context.Background())

	if err != nil {
		p.logger.Warn().
			Timestamp().
			Str("url", p.document.GetURL().String()).
			Err(err).
			Msg("failed to close browser page")
	}

	return p.conn.Close()
}

func (p *HTMLPage) Document() drivers.HTMLDocument {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.document
}

func (p *HTMLPage) GetCookies(ctx context.Context) (*values.Array, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	repl, err := p.client.Network.GetAllCookies(ctx)

	if err != nil {
		return values.NewArray(0), err
	}

	if repl.Cookies == nil {
		return values.NewArray(0), nil
	}

	cookies := values.NewArray(len(repl.Cookies))

	for _, c := range repl.Cookies {
		cookies.Push(toDriverCookie(c))
	}

	return cookies, nil
}

func (p *HTMLPage) SetCookies(ctx context.Context, cookies ...drivers.HTTPCookie) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(cookies) == 0 {
		return nil
	}

	params := make([]network.CookieParam, 0, len(cookies))

	for _, c := range cookies {
		params = append(params, fromDriverCookie(p.document.GetURL().String(), c))
	}

	return p.client.Network.SetCookies(ctx, network.NewSetCookiesArgs(params))
}

func (p *HTMLPage) DeleteCookies(ctx context.Context, cookies ...drivers.HTTPCookie) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(cookies) == 0 {
		return nil
	}

	var err error

	for _, c := range cookies {
		err = p.client.Network.DeleteCookies(ctx, fromDriverCookieDelete(p.document.GetURL().String(), c))

		if err != nil {
			break
		}
	}

	return err
}

func (p *HTMLPage) PrintToPDF(ctx context.Context, params drivers.PDFParams) (values.Binary, error) {
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
	if url == "" {
		url = BlankPageURL
	}

	repl, err := p.client.Page.Navigate(ctx, page.NewNavigateArgs(url.String()))

	if err != nil {
		return err
	}

	if repl.ErrorText != nil {
		return errors.New(*repl.ErrorText)
	}

	return p.WaitForNavigation(ctx)
}

func (p *HTMLPage) NavigateBack(ctx context.Context, skip values.Int) (values.Boolean, error) {
	history, err := p.client.Page.GetNavigationHistory(ctx)

	if err != nil {
		return values.False, err
	}

	// we are in the beginning
	if history.CurrentIndex == 0 {
		return values.False, nil
	}

	if skip < 1 {
		skip = 1
	}

	to := history.CurrentIndex - int(skip)

	if to < 0 {
		// TODO: Return error?
		return values.False, nil
	}

	prev := history.Entries[to]
	err = p.client.Page.NavigateToHistoryEntry(ctx, page.NewNavigateToHistoryEntryArgs(prev.ID))

	if err != nil {
		return values.False, err
	}

	err = p.WaitForNavigation(ctx)

	if err != nil {
		return values.False, err
	}

	return values.True, nil
}

func (p *HTMLPage) NavigateForward(ctx context.Context, skip values.Int) (values.Boolean, error) {
	history, err := p.client.Page.GetNavigationHistory(ctx)

	if err != nil {
		return values.False, err
	}

	length := len(history.Entries)
	lastIndex := length - 1

	// nowhere to go forward
	if history.CurrentIndex == lastIndex {
		return values.False, nil
	}

	if skip < 1 {
		skip = 1
	}

	to := int(skip) + history.CurrentIndex

	if to > lastIndex {
		// TODO: Return error?
		return values.False, nil
	}

	next := history.Entries[to]
	err = p.client.Page.NavigateToHistoryEntry(ctx, page.NewNavigateToHistoryEntryArgs(next.ID))

	if err != nil {
		return values.False, err
	}

	err = p.WaitForNavigation(ctx)

	if err != nil {
		return values.False, err
	}

	return values.True, nil
}

func (p *HTMLPage) WaitForNavigation(ctx context.Context) error {
	onEvent := make(chan struct{})
	var once sync.Once
	listener := func(_ context.Context, _ interface{}) {
		once.Do(func() {
			close(onEvent)
		})
	}

	defer p.events.RemoveEventListener(events.EventLoad, listener)

	p.events.AddEventListener(events.EventLoad, listener)

	select {
	case <-onEvent:
		return nil
	case <-ctx.Done():
		return core.ErrTimeout
	}
}

func (p *HTMLPage) handlePageLoad(ctx context.Context, _ interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	errs := make([]error, 0, 5)

	nextDoc, err := LoadRootHTMLDocument(ctx, p.logger, p.client, p.events)

	if err != nil {
		p.logger.Error().
			Timestamp().
			Err(err).
			Msg("failed to load new root document after page load")

		return
	}

	// close the prev document
	err = p.document.Close()

	if err != nil {
		errs = append(errs, errors.Wrapf(err, "failed to close root document: %s", p.document.GetURL()))
	}

	// set the new root document
	p.document = nextDoc
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
