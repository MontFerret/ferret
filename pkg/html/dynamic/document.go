package dynamic

import (
	"context"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/MontFerret/ferret/pkg/html/dynamic/eval"
	"github.com/MontFerret/ferret/pkg/html/dynamic/events"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/input"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

const BlankPageURL = "about:blank"

type (
	ScreenshotFormat string
	ScreenshotArgs   struct {
		X       float64
		Y       float64
		Width   float64
		Height  float64
		Format  ScreenshotFormat
		Quality int
	}

	HTMLDocument struct {
		sync.Mutex
		logger  *zerolog.Logger
		conn    *rpcc.Conn
		client  *cdp.Client
		events  *events.EventBroker
		url     values.String
		element *HTMLElement
	}
)

const (
	ScreenshotFormatPNG  ScreenshotFormat = "png"
	ScreenshotFormatJPEG ScreenshotFormat = "jpeg"
)

func IsScreenshotFormatValid(format string) bool {
	value := ScreenshotFormat(format)

	return value == ScreenshotFormatPNG || value == ScreenshotFormatJPEG
}

func LoadHTMLDocument(
	ctx context.Context,
	conn *rpcc.Conn,
	client *cdp.Client,
	url string,
) (*HTMLDocument, error) {
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
			return nil, err
		}
	}

	node, err := getRootElement(ctx, client)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get root element")
	}

	broker, err := createEventBroker(client)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create event events")
	}

	logger := logging.FromContext(ctx)

	rootElement, err := LoadElement(
		ctx,
		logger,
		client,
		broker,
		node.Root.NodeID,
		node.Root.BackendNodeID,
	)

	if err != nil {
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

	broker.AddEventListener("load", doc.handlePageLoad)
	broker.AddEventListener("error", doc.handleError)

	return doc
}

func (doc *HTMLDocument) MarshalJSON() ([]byte, error) {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.MarshalJSON()
}

func (doc *HTMLDocument) Type() core.Type {
	return core.HTMLDocumentType
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

func (doc *HTMLDocument) Clone() core.Value {
	return values.None
}

func (doc *HTMLDocument) Compare(other core.Value) int {
	doc.Lock()
	defer doc.Unlock()

	switch other.Type() {
	case core.HTMLDocumentType:
		other := other.(*HTMLDocument)

		return doc.url.Compare(other.url)
	default:
		if other.Type() > core.HTMLDocumentType {
			return -1
		}

		return 1
	}
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

func (doc *HTMLDocument) InnerText() values.String {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.InnerText()
}

func (doc *HTMLDocument) InnerHTML() values.String {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.InnerHTML()
}

func (doc *HTMLDocument) Value() core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.Value()
}

func (doc *HTMLDocument) GetAttributes() core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.GetAttributes()
}

func (doc *HTMLDocument) GetAttribute(name values.String) core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.GetAttribute(name)
}

func (doc *HTMLDocument) GetChildNodes() core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.GetChildNodes()
}

func (doc *HTMLDocument) GetChildNode(idx values.Int) core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.GetChildNode(idx)
}

func (doc *HTMLDocument) QuerySelector(selector values.String) core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.QuerySelector(selector)
}

func (doc *HTMLDocument) QuerySelectorAll(selector values.String) core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.QuerySelectorAll(selector)
}

func (doc *HTMLDocument) URL() core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.url
}

func (doc *HTMLDocument) InnerHTMLBySelector(selector values.String) values.String {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.InnerHTMLBySelector(selector)
}

func (doc *HTMLDocument) InnerHTMLBySelectorAll(selector values.String) *values.Array {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.InnerHTMLBySelectorAll(selector)
}

func (doc *HTMLDocument) InnerTextBySelector(selector values.String) values.String {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.InnerHTMLBySelector(selector)
}

func (doc *HTMLDocument) InnerTextBySelectorAll(selector values.String) *values.Array {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.InnerTextBySelectorAll(selector)
}

func (doc *HTMLDocument) ClickBySelector(selector values.String) (values.Boolean, error) {
	res, err := eval.Eval(
		doc.client,
		fmt.Sprintf(`
			var el = document.querySelector(%s);

			if (el == null) {
				return false;
			}

			var evt = new window.MouseEvent('click', { bubbles: true });
			el.dispatchEvent(evt);

			return true;
		`, eval.ParamString(selector.String())),
		true,
		false,
	)

	if err != nil {
		return values.False, err
	}

	if res.Type() == core.BooleanType {
		return res.(values.Boolean), nil
	}

	return values.False, nil
}

func (doc *HTMLDocument) ClickBySelectorAll(selector values.String) (values.Boolean, error) {
	res, err := eval.Eval(
		doc.client,
		fmt.Sprintf(`
			var elements = document.querySelectorAll(%s);

			if (elements == null) {
				return false;
			}

			elements.forEach((el) => {
				var evt = new window.MouseEvent('click', { bubbles: true });
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

	if res.Type() == core.BooleanType {
		return res.(values.Boolean), nil
	}

	return values.False, nil
}

func (doc *HTMLDocument) InputBySelector(selector values.String, value core.Value, delay values.Int) (values.Boolean, error) {
	ctx := context.Background()

	valStr := value.String()

	res, err := eval.Eval(
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

	if res.Type() == core.BooleanType && res.(values.Boolean) == values.False {
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

func (doc *HTMLDocument) WaitForSelector(selector values.String, timeout values.Int) error {
	task := events.NewEvalWaitTask(
		doc.client,
		fmt.Sprintf(`
			var el = document.querySelector(%s);

			if (el != null) {
				return true;
			}

			// null means we need to repeat
			return null;
		`, eval.ParamString(selector.String())),
		time.Millisecond*time.Duration(timeout),
		events.DefaultPolling,
	)

	_, err := task.Run()

	return err
}

func (doc *HTMLDocument) WaitForClass(selector, class values.String, timeout values.Int) error {
	task := events.NewEvalWaitTask(
		doc.client,
		fmt.Sprintf(`
			var el = document.querySelector(%s);

			if (el == null) {
				return false;
			}

			var className = %s;
			var found = el.className.split(' ').find(i => i === className);

			if (found != null) {
				return true;
			}
			
			// null means we need to repeat
			return null;
		`,
			eval.ParamString(selector.String()),
			eval.ParamString(class.String()),
		),
		time.Millisecond*time.Duration(timeout),
		events.DefaultPolling,
	)

	_, err := task.Run()

	return err
}

func (doc *HTMLDocument) WaitForClassAll(selector, class values.String, timeout values.Int) error {
	task := events.NewEvalWaitTask(
		doc.client,
		fmt.Sprintf(`
			var elements = document.querySelectorAll(%s);

			if (elements == null || elements.length === 0) {
				return false;
			}

			var className = %s;
			var foundCount = 0;

			elements.forEach((el) => {
				var found = el.className.split(' ').find(i => i === className);

				if (found != null) {
					foundCount++;
				}
			});

			if (foundCount === elements.length) {
				return true;
			}
			
			// null means we need to repeat
			return null;
		`,
			eval.ParamString(selector.String()),
			eval.ParamString(class.String()),
		),
		time.Millisecond*time.Duration(timeout),
		events.DefaultPolling,
	)

	_, err := task.Run()

	return err
}

func (doc *HTMLDocument) WaitForNavigation(timeout values.Int) error {
	onEvent := make(chan struct{})
	listener := func(_ interface{}) {
		close(onEvent)
	}

	defer doc.events.RemoveEventListener("load", listener)

	doc.events.AddEventListener("load", listener)

	select {
	case <-onEvent:
		return nil
	case <-time.After(time.Millisecond * time.Duration(timeout)):
		return core.ErrTimeout
	}
}

func (doc *HTMLDocument) Navigate(url values.String, timeout values.Int) error {
	if url == "" {
		url = BlankPageURL
	}

	ctx := context.Background()
	repl, err := doc.client.Page.Navigate(ctx, page.NewNavigateArgs(url.String()))

	if err != nil {
		return err
	}

	if repl.ErrorText != nil {
		return errors.New(*repl.ErrorText)
	}

	return doc.WaitForNavigation(timeout)
}

func (doc *HTMLDocument) CaptureScreenshot(params *ScreenshotArgs) (core.Value, error) {
	ctx := context.Background()
	metrics, err := doc.client.Page.GetLayoutMetrics(ctx)

	if params.Format == ScreenshotFormatJPEG && params.Quality < 0 && params.Quality > 100 {
		params.Quality = 100
	}

	if params.X < 0 {
		params.X = 0
	}

	if params.Y < 0 {
		params.Y = 0
	}

	if params.Width <= 0 {
		params.Width = float64(metrics.LayoutViewport.ClientWidth) - params.X
	}

	if params.Height <= 0 {
		params.Height = float64(metrics.LayoutViewport.ClientHeight) - params.Y
	}

	clip := page.Viewport{
		X:      params.X,
		Y:      params.Y,
		Width:  params.Width,
		Height: params.Height,
		Scale:  1.0,
	}

	format := string(params.Format)
	screenshotArgs := page.CaptureScreenshotArgs{
		Format:  &format,
		Quality: &params.Quality,
		Clip:    &clip,
	}

	reply, err := doc.client.Page.CaptureScreenshot(ctx, &screenshotArgs)

	if err != nil {
		return values.None, err
	}

	return values.NewBinary(reply.Data), nil
}

func (doc *HTMLDocument) handlePageLoad(_ interface{}) {
	doc.Lock()
	defer doc.Unlock()

	ctx, cancel := contextWithTimeout()
	defer cancel()

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

func (doc *HTMLDocument) handleError(val interface{}) {
	err, ok := val.(error)

	if !ok {
		return
	}

	doc.logger.Error().
		Timestamp().
		Err(err).
		Msg("unexpected error")
}
