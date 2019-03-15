package cdp

import (
	"context"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/eval"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
	"github.com/MontFerret/ferret/pkg/drivers/cdp/templates"
	"github.com/MontFerret/ferret/pkg/drivers/common"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/input"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const BlankPageURL = "about:blank"

type HTMLDocument struct {
	sync.Mutex
	logger  *zerolog.Logger
	conn    *rpcc.Conn
	client  *cdp.Client
	events  *events.EventBroker
	url     values.String
	element *HTMLElement
}

func handleLoadError(logger *zerolog.Logger, client *cdp.Client) {
	err := client.Page.Close(context.Background())

	if err != nil {
		logger.Warn().Timestamp().Err(err).Msg("unabled to close document on load error")
	}
}

func LoadHTMLDocument(
	ctx context.Context,
	conn *rpcc.Conn,
	client *cdp.Client,
	url string,
) (drivers.HTMLDocument, error) {
	logger := logging.FromContext(ctx)

	if conn == nil {
		return nil, core.Error(core.ErrMissedArgument, "connection")
	}

	if url == "" {
		return nil, core.Error(core.ErrMissedArgument, "url")
	}

	var err error

	if url != BlankPageURL {
		err = waitForLoadEvent(ctx, client)

		if err != nil {
			handleLoadError(logger, client)

			return nil, err
		}
	}

	node, err := getRootElement(ctx, client)

	if err != nil {
		handleLoadError(logger, client)
		return nil, errors.Wrap(err, "failed to get root element")
	}

	broker, err := createEventBroker(client)

	if err != nil {
		handleLoadError(logger, client)
		return nil, errors.Wrap(err, "failed to create event events")
	}

	rootElement, err := LoadElement(
		ctx,
		logger,
		client,
		broker,
		node.Root.NodeID,
		node.Root.BackendNodeID,
	)

	if err != nil {
		broker.Stop()
		broker.Close()
		handleLoadError(logger, client)

		return nil, errors.Wrap(err, "failed to load root element")
	}

	return NewHTMLDocument(
		logger,
		conn,
		client,
		broker,
		values.NewString(url),
		rootElement,
	), nil
}

func NewHTMLDocument(
	logger *zerolog.Logger,
	conn *rpcc.Conn,
	client *cdp.Client,
	broker *events.EventBroker,
	url values.String,
	rootElement *HTMLElement,
) *HTMLDocument {
	doc := new(HTMLDocument)
	doc.logger = logger
	doc.conn = conn
	doc.client = client
	doc.events = broker
	doc.url = url
	doc.element = rootElement

	broker.AddEventListener(events.EventLoad, doc.handlePageLoad)
	broker.AddEventListener(events.EventError, doc.handleError)

	return doc
}

func (doc *HTMLDocument) MarshalJSON() ([]byte, error) {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.MarshalJSON()
}

func (doc *HTMLDocument) Type() core.Type {
	return drivers.HTMLDocumentType
}

func (doc *HTMLDocument) String() string {
	doc.Lock()
	defer doc.Unlock()

	return doc.url.String()
}

func (doc *HTMLDocument) Unwrap() interface{} {
	doc.Lock()
	defer doc.Unlock()

	return doc.element
}

func (doc *HTMLDocument) Hash() uint64 {
	doc.Lock()
	defer doc.Unlock()

	h := fnv.New64a()

	h.Write([]byte(doc.Type().String()))
	h.Write([]byte(":"))
	h.Write([]byte(doc.url))

	return h.Sum64()
}

func (doc *HTMLDocument) Copy() core.Value {
	return values.None
}

func (doc *HTMLDocument) Compare(other core.Value) int64 {
	doc.Lock()
	defer doc.Unlock()

	switch other.Type() {
	case drivers.HTMLDocumentType:
		other := other.(drivers.HTMLDocument)

		return doc.url.Compare(other.GetURL())
	default:
		return drivers.Compare(doc.Type(), other.Type())
	}
}

func (doc *HTMLDocument) Iterate(ctx context.Context) (core.Iterator, error) {
	return doc.element.Iterate(ctx)
}

func (doc *HTMLDocument) GetIn(ctx context.Context, path []core.Value) (core.Value, error) {
	return common.GetInDocument(ctx, doc, path)
}

func (doc *HTMLDocument) SetIn(ctx context.Context, path []core.Value, value core.Value) error {
	return common.SetInDocument(ctx, doc, path, value)
}

func (doc *HTMLDocument) Close() error {
	doc.Lock()
	defer doc.Unlock()

	var err error

	err = doc.events.Stop()

	if err != nil {
		doc.logger.Warn().
			Timestamp().
			Str("url", doc.url.String()).
			Err(err).
			Msg("failed to stop event events")
	}

	err = doc.events.Close()

	if err != nil {
		doc.logger.Warn().
			Timestamp().
			Str("url", doc.url.String()).
			Err(err).
			Msg("failed to close event events")
	}

	err = doc.element.Close()

	if err != nil {
		doc.logger.Warn().
			Timestamp().
			Str("url", doc.url.String()).
			Err(err).
			Msg("failed to close root element")
	}

	err = doc.client.Page.Close(context.Background())

	if err != nil {
		doc.logger.Warn().
			Timestamp().
			Str("url", doc.url.String()).
			Err(err).
			Msg("failed to close browser page")
	}

	return doc.conn.Close()
}

func (doc *HTMLDocument) NodeType() values.Int {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.NodeType()
}

func (doc *HTMLDocument) NodeName() values.String {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.NodeName()
}

func (doc *HTMLDocument) Length() values.Int {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.Length()
}

func (doc *HTMLDocument) GetChildNodes(ctx context.Context) core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.GetChildNodes(ctx)
}

func (doc *HTMLDocument) GetChildNode(ctx context.Context, idx values.Int) core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.GetChildNode(ctx, idx)
}

func (doc *HTMLDocument) QuerySelector(ctx context.Context, selector values.String) core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.QuerySelector(ctx, selector)
}

func (doc *HTMLDocument) QuerySelectorAll(ctx context.Context, selector values.String) core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.QuerySelectorAll(ctx, selector)
}

func (doc *HTMLDocument) DocumentElement() drivers.HTMLElement {
	doc.Lock()
	defer doc.Unlock()

	return doc.element
}

func (doc *HTMLDocument) GetURL() core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.url
}

func (doc *HTMLDocument) SetURL(ctx context.Context, url values.String) error {
	return doc.Navigate(ctx, url)
}

func (doc *HTMLDocument) CountBySelector(ctx context.Context, selector values.String) values.Int {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.CountBySelector(ctx, selector)
}

func (doc *HTMLDocument) ExistsBySelector(ctx context.Context, selector values.String) values.Boolean {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.ExistsBySelector(ctx, selector)
}

func (doc *HTMLDocument) ClickBySelector(ctx context.Context, selector values.String) (values.Boolean, error) {
	res, err := eval.Eval(
		ctx,
		doc.client,
		fmt.Sprintf(`
			var el = document.querySelector(%s);
			if (el == null) {
				return false;
			}
			var evt = new window.MouseEvent('click', { bubbles: true, cancelable: true });
			el.dispatchEvent(evt);
			return true;
		`, eval.ParamString(selector.String())),
		true,
		false,
	)

	if err != nil {
		return values.False, err
	}

	if res.Type() == types.Boolean {
		return res.(values.Boolean), nil
	}

	return values.False, nil
}

func (doc *HTMLDocument) ClickBySelectorAll(ctx context.Context, selector values.String) (values.Boolean, error) {
	res, err := eval.Eval(
		ctx,
		doc.client,
		fmt.Sprintf(`
			var elements = document.querySelectorAll(%s);
			if (elements == null) {
				return false;
			}
			elements.forEach((el) => {
				var evt = new window.MouseEvent('click', { bubbles: true, cancelable: true });
				el.dispatchEvent(evt);	
			});
			return true;
		`, eval.ParamString(selector.String())),
		true,
		false,
	)

	if err != nil {
		return values.False, err
	}

	if res.Type() == types.Boolean {
		return res.(values.Boolean), nil
	}

	return values.False, nil
}

func (doc *HTMLDocument) InputBySelector(ctx context.Context, selector values.String, value core.Value, delay values.Int) (values.Boolean, error) {
	valStr := value.String()

	res, err := eval.Eval(
		ctx,
		doc.client,
		fmt.Sprintf(`
			var el = document.querySelector(%s);
			if (el == null) {
				return false;
			}
			el.focus();
			return true;
		`, eval.ParamString(selector.String())),
		true,
		false,
	)

	if err != nil {
		return values.False, err
	}

	if res.Type() == types.Boolean && res.(values.Boolean) == values.False {
		return values.False, nil
	}

	delayMs := time.Duration(delay)

	time.Sleep(delayMs * time.Millisecond)

	for _, ch := range valStr {
		for _, ev := range []string{"keyDown", "keyUp"} {
			ke := input.NewDispatchKeyEventArgs(ev).SetText(string(ch))

			if err := doc.client.Input.DispatchKeyEvent(ctx, ke); err != nil {
				return values.False, err
			}

			time.Sleep(delayMs * time.Millisecond)
		}
	}

	return values.True, nil
}

func (doc *HTMLDocument) SelectBySelector(ctx context.Context, selector values.String, value *values.Array) (*values.Array, error) {
	res, err := eval.Eval(
		ctx,
		doc.client,
		fmt.Sprintf(`
			var element = document.querySelector(%s);
			if (element == null) {
				return [];
			}
			var values = %s;
			if (element.nodeName.toLowerCase() !== 'select') {
				throw new Error('Element is not a <select> element.');
			}
			var options = Array.from(element.options);
      		element.value = undefined;
			for (var option of options) {
        		option.selected = values.includes(option.value);
        	
				if (option.selected && !element.multiple) {
          			break;
				}
      		}
      		element.dispatchEvent(new Event('input', { 'bubbles': true, cancelable: true }));
      		element.dispatchEvent(new Event('change', { 'bubbles': true, cancelable: true }));
      		
			return options.filter(option => option.selected).map(option => option.value);
		`,
			eval.ParamString(selector.String()),
			value.String(),
		),
		true,
		false,
	)

	if err != nil {
		return nil, err
	}

	arr, ok := res.(*values.Array)

	if ok {
		return arr, nil
	}

	return nil, core.TypeError(types.Array, res.Type())
}

func (doc *HTMLDocument) MoveMouseBySelector(ctx context.Context, selector values.String) error {
	err := doc.ScrollBySelector(ctx, selector)

	if err != nil {
		return err
	}

	selectorArgs := dom.NewQuerySelectorArgs(doc.element.id.nodeID, selector.String())
	found, err := doc.client.DOM.QuerySelector(ctx, selectorArgs)

	if err != nil {
		doc.element.logError(err).
			Str("selector", selector.String()).
			Msg("failed to retrieve a node by selector")

		return err
	}

	if found.NodeID <= 0 {
		return errors.New("element not found")
	}

	q, err := getClickablePoint(ctx, doc.client, &HTMLElementIdentity{
		nodeID: found.NodeID,
	})

	if err != nil {
		return err
	}

	return doc.client.Input.DispatchMouseEvent(
		ctx,
		input.NewDispatchMouseEventArgs("mouseMoved", q.X, q.Y),
	)
}

func (doc *HTMLDocument) MoveMouseByXY(ctx context.Context, x, y values.Float) error {
	return doc.client.Input.DispatchMouseEvent(
		ctx,
		input.NewDispatchMouseEventArgs("mouseMoved", float64(x), float64(y)),
	)
}

func (doc *HTMLDocument) WaitForElement(ctx context.Context, selector values.String, when drivers.WaitEvent) error {
	var operator string

	if when == drivers.WaitEventPresence {
		operator = "!="
	} else {
		operator = "=="
	}

	task := events.NewEvalWaitTask(
		doc.client,
		fmt.Sprintf(
			`
				var el = document.querySelector(%s);
				
				if (el %s null) {
					return true;
				}
				
				// null means we need to repeat
				return null;
			`,
			eval.ParamString(selector.String()),
			operator,
		),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForClassBySelector(ctx context.Context, selector, class values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		doc.client,
		templates.WaitBySelector(
			selector,
			when,
			class,
			fmt.Sprintf("el.className.split(' ').find(i => i === %s)", eval.ParamString(class.String())),
		),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForClassBySelectorAll(ctx context.Context, selector, class values.String, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		doc.client,
		templates.WaitBySelectorAll(
			selector,
			when,
			class,
			fmt.Sprintf("el.className.split(' ').find(i => i === %s)", eval.ParamString(class.String())),
		),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForNavigation(ctx context.Context) error {
	onEvent := make(chan struct{})
	listener := func(_ context.Context, _ interface{}) {
		close(onEvent)
	}

	defer doc.events.RemoveEventListener(events.EventLoad, listener)

	doc.events.AddEventListener(events.EventLoad, listener)

	select {
	case <-onEvent:
		return nil
	case <-ctx.Done():
		return core.ErrTimeout
	}
}

func (doc *HTMLDocument) WaitForAttributeBySelector(
	ctx context.Context,
	selector,
	name values.String,
	value core.Value,
	when drivers.WaitEvent,
) error {
	task := events.NewEvalWaitTask(
		doc.client,
		templates.WaitBySelector(
			selector,
			when,
			value,
			templates.AttributeRead(name),
		),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForAttributeBySelectorAll(
	ctx context.Context,
	selector,
	name values.String,
	value core.Value,
	when drivers.WaitEvent,
) error {
	task := events.NewEvalWaitTask(
		doc.client,
		templates.WaitBySelectorAll(
			selector,
			when,
			value,
			templates.AttributeRead(name),
		),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForStyleBySelector(ctx context.Context, selector, name values.String, value core.Value, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		doc.client,
		templates.WaitBySelector(
			selector,
			when,
			value,
			templates.StyleRead(name),
		),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) WaitForStyleBySelectorAll(ctx context.Context, selector, name values.String, value core.Value, when drivers.WaitEvent) error {
	task := events.NewEvalWaitTask(
		doc.client,
		templates.WaitBySelectorAll(
			selector,
			when,
			value,
			templates.StyleRead(name),
		),
		events.DefaultPolling,
	)

	_, err := task.Run(ctx)

	return err
}

func (doc *HTMLDocument) Navigate(ctx context.Context, url values.String) error {
	if url == "" {
		url = BlankPageURL
	}

	repl, err := doc.client.Page.Navigate(ctx, page.NewNavigateArgs(url.String()))

	if err != nil {
		return err
	}

	if repl.ErrorText != nil {
		return errors.New(*repl.ErrorText)
	}

	return doc.WaitForNavigation(ctx)
}

func (doc *HTMLDocument) NavigateBack(ctx context.Context, skip values.Int) (values.Boolean, error) {
	history, err := doc.client.Page.GetNavigationHistory(ctx)

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
	err = doc.client.Page.NavigateToHistoryEntry(ctx, page.NewNavigateToHistoryEntryArgs(prev.ID))

	if err != nil {
		return values.False, err
	}

	err = doc.WaitForNavigation(ctx)

	if err != nil {
		return values.False, err
	}

	return values.True, nil
}

func (doc *HTMLDocument) NavigateForward(ctx context.Context, skip values.Int) (values.Boolean, error) {
	history, err := doc.client.Page.GetNavigationHistory(ctx)

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
	err = doc.client.Page.NavigateToHistoryEntry(ctx, page.NewNavigateToHistoryEntryArgs(next.ID))

	if err != nil {
		return values.False, err
	}

	err = doc.WaitForNavigation(ctx)

	if err != nil {
		return values.False, err
	}

	return values.True, nil
}

func (doc *HTMLDocument) PrintToPDF(ctx context.Context, params drivers.PDFParams) (values.Binary, error) {
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

	reply, err := doc.client.Page.PrintToPDF(ctx, args)

	if err != nil {
		return values.NewBinary([]byte{}), err
	}

	return values.NewBinary(reply.Data), nil
}

func (doc *HTMLDocument) CaptureScreenshot(ctx context.Context, params drivers.ScreenshotParams) (values.Binary, error) {
	metrics, err := doc.client.Page.GetLayoutMetrics(ctx)

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

	reply, err := doc.client.Page.CaptureScreenshot(ctx, &args)

	if err != nil {
		return values.NewBinary([]byte{}), err
	}

	return values.NewBinary(reply.Data), nil
}

func (doc *HTMLDocument) ScrollTop(ctx context.Context) error {
	_, err := eval.Eval(ctx, doc.client, `
		window.scrollTo({
			left: 0,
			top: 0,
    		behavior: 'instant'
  		});
	`, false, false)

	return err
}

func (doc *HTMLDocument) ScrollBottom(ctx context.Context) error {
	_, err := eval.Eval(ctx, doc.client, `
		window.scrollTo({
			left: 0,
			top: window.document.body.scrollHeight,
    		behavior: 'instant'
  		});
	`, false, false)

	return err
}

func (doc *HTMLDocument) ScrollBySelector(ctx context.Context, selector values.String) error {
	_, err := eval.Eval(ctx, doc.client, fmt.Sprintf(`
		var el = document.querySelector(%s);
		if (el == null) {
			throw new Error("element not found");
		}
		el.scrollIntoView({
    		behavior: 'instant'
  		});
		return true;
	`, eval.ParamString(selector.String()),
	), false, false)

	return err
}

func (doc *HTMLDocument) ScrollByXY(ctx context.Context, x, y values.Float) error {
	_, err := eval.Eval(ctx, doc.client, fmt.Sprintf(`
		window.scrollBy({
  			top: %s,
  			left: %s,
  			behavior: 'instant'
		});
	`,
		eval.ParamFloat(float64(x)),
		eval.ParamFloat(float64(y)),
	), false, false)

	return err
}

func (doc *HTMLDocument) handlePageLoad(ctx context.Context, _ interface{}) {
	doc.Lock()
	defer doc.Unlock()

	node, err := getRootElement(ctx, doc.client)

	if err != nil {
		doc.logger.Error().
			Timestamp().
			Err(err).
			Msg("failed to get root node after page load")

		return
	}

	updated, err := LoadElement(
		ctx,
		doc.logger,
		doc.client,
		doc.events,
		node.Root.NodeID,
		node.Root.BackendNodeID,
	)

	if err != nil {
		doc.logger.Error().
			Timestamp().
			Err(err).
			Msg("failed to load root node after page load")

		return
	}

	// close the prev element
	doc.element.Close()

	// create a new root element wrapper
	doc.element = updated
	doc.url = ""

	if node.Root.BaseURL != nil {
		doc.url = values.NewString(*node.Root.BaseURL)
	}
}

func (doc *HTMLDocument) handleError(_ context.Context, val interface{}) {
	err, ok := val.(error)

	if !ok {
		return
	}

	doc.logger.Error().
		Timestamp().
		Err(err).
		Msg("unexpected error")
}
