package dynamic

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/logging"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/dynamic/eval"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/dynamic/events"
	"github.com/corpix/uarand"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/emulation"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"hash/fnv"
	"sync"
	"time"
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

func LoadHTMLDocument(
	ctx context.Context,
	conn *rpcc.Conn,
	url string,
) (*HTMLDocument, error) {
	if conn == nil {
		return nil, core.Error(core.ErrMissedArgument, "connection")
	}

	if url == "" {
		return nil, core.Error(core.ErrMissedArgument, "url")
	}

	client := cdp.NewClient(conn)

	err := runBatch(
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
			return client.Emulation.SetUserAgentOverride(
				ctx,
				emulation.NewSetUserAgentOverrideArgs(uarand.GetRandom()),
			)
		},
	)

	if err != nil {
		return nil, err
	}

	if url != BlankPageURL {
		err = waitForLoadEvent(ctx, client)

		if err != nil {
			return nil, err
		}
	}

	root, innerHTML, err := getRootElement(client)

	if err != nil {
		return nil, err
	}

	broker, err := createEventBroker(client)

	if err != nil {
		return nil, err
	}

	return NewHTMLDocument(
		logging.FromContext(ctx),
		conn,
		client,
		broker,
		root,
		innerHTML,
	), nil
}

func getRootElement(client *cdp.Client) (dom.Node, values.String, error) {
	args := dom.NewGetDocumentArgs()
	args.Depth = pointerInt(1) // lets load the entire document
	ctx := context.Background()

	d, err := client.DOM.GetDocument(ctx, args)

	if err != nil {
		return dom.Node{}, values.EmptyString, err
	}

	innerHTML, err := client.DOM.GetOuterHTML(ctx, dom.NewGetOuterHTMLArgs().SetNodeID(d.Root.NodeID))

	if err != nil {
		return dom.Node{}, values.EmptyString, err
	}

	return d.Root, values.NewString(innerHTML.OuterHTML), nil
}

func NewHTMLDocument(
	logger *zerolog.Logger,
	conn *rpcc.Conn,
	client *cdp.Client,
	broker *events.EventBroker,
	root dom.Node,
	innerHTML values.String,
) *HTMLDocument {
	doc := new(HTMLDocument)
	doc.logger = logger
	doc.conn = conn
	doc.client = client
	doc.events = broker
	doc.element = NewHTMLElement(doc.logger, client, broker, root.NodeID, root, innerHTML)
	doc.url = ""

	if root.BaseURL != nil {
		doc.url = values.NewString(*root.BaseURL)
	}

	broker.AddEventListener("load", func(_ interface{}) {
		doc.Lock()
		defer doc.Unlock()

		updated, innerHTML, err := getRootElement(client)

		if err != nil {
			doc.logger.Error().
				Timestamp().
				Err(err).
				Msg("failed to get root node after page load")

			return
		}

		// close the prev element
		doc.element.Close()

		// create a new root element wrapper
		doc.element = NewHTMLElement(doc.logger, client, broker, updated.NodeID, updated, innerHTML)
		doc.url = ""

		if updated.BaseURL != nil {
			doc.url = values.NewString(*updated.BaseURL)
		}
	})

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
			Msg("failed to stop event broker")
	}

	err = doc.events.Close()

	if err != nil {
		doc.logger.Warn().
			Timestamp().
			Str("url", doc.url.String()).
			Err(err).
			Msg("failed to close event broker")
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
	return doc.url
}

func (doc *HTMLDocument) InnerHTMLBySelector(selector values.String) values.String {
	res, err := eval.Eval(
		doc.client,
		fmt.Sprintf(`
			var el = document.querySelector(%s);

			if (el == null) {
				return "";
			}

			return el.innerHTML;
		`, eval.ParamString(selector.String())),
		true,
		false,
	)

	if err != nil {
		doc.logger.Error().
			Timestamp().
			Err(err).
			Str("selector", selector.String()).
			Msg("failed to get inner HTML by selector")

		return values.EmptyString
	}

	if res.Type() == core.StringType {
		return res.(values.String)
	}

	return values.EmptyString
}

func (doc *HTMLDocument) InnerHTMLBySelectorAll(selector values.String) *values.Array {
	res, err := eval.Eval(
		doc.client,
		fmt.Sprintf(`
			var result = [];
			var elements = document.querySelectorAll(%s);

			if (elements == null) {
				return result;
			}

			elements.forEach((i) => {
				result.push(i.innerHTML);
			});

			return result;
		`, eval.ParamString(selector.String())),
		true,
		false,
	)

	if err != nil {
		doc.logger.Error().
			Timestamp().
			Err(err).
			Str("selector", selector.String()).
			Msg("failed to get an array of inner HTML by selector")

		return values.NewArray(0)
	}

	if res.Type() == core.ArrayType {
		return res.(*values.Array)
	}

	return values.NewArray(0)
}

func (doc *HTMLDocument) InnerTextBySelector(selector values.String) values.String {
	res, err := eval.Eval(
		doc.client,
		fmt.Sprintf(`
			var el = document.querySelector(%s);

			if (el == null) {
				return "";
			}

			return el.innerText;
		`, eval.ParamString(selector.String())),
		true,
		false,
	)

	if err != nil {
		doc.logger.Error().
			Timestamp().
			Err(err).
			Str("selector", selector.String()).
			Msg("failed to get inner text by selector")

		return values.EmptyString
	}

	if res.Type() == core.StringType {
		return res.(values.String)
	}

	return values.EmptyString
}

func (doc *HTMLDocument) InnerTextBySelectorAll(selector values.String) *values.Array {
	res, err := eval.Eval(
		doc.client,
		fmt.Sprintf(`
			var result = [];
			var elements = document.querySelectorAll(%s);

			if (elements == null) {
				return result;
			}

			elements.forEach((i) => {
				result.push(i.innerText);
			});

			return result;
		`, eval.ParamString(selector.String())),
		true,
		false,
	)

	if err != nil {
		doc.logger.Error().
			Timestamp().
			Err(err).
			Str("selector", selector.String()).
			Msg("failed to get an array inner text by selector")

		return values.NewArray(0)
	}

	if res.Type() == core.ArrayType {
		return res.(*values.Array)
	}

	return values.NewArray(0)
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

func (doc *HTMLDocument) InputBySelector(selector values.String, value core.Value) (values.Boolean, error) {
	res, err := eval.Eval(
		doc.client,
		fmt.Sprintf(
			`
			var el = document.querySelector(%s);

			if (el == null) {
				return false;
			}

			var evt = new window.Event('input', { bubbles: true });

			el.value = %s
			el.dispatchEvent(evt);

			return true;
		`,
			eval.ParamString(selector.String()),
			eval.ParamString(value.String()),
		),
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
	timer := time.NewTimer(time.Millisecond * time.Duration(timeout))
	onEvent := make(chan bool)
	listener := func(_ interface{}) {
		onEvent <- true
	}

	defer doc.events.RemoveEventListener("load", listener)
	defer close(onEvent)

	doc.events.AddEventListener("load", listener)

	for {
		select {
		case <-onEvent:
			timer.Stop()

			return nil
		case <-timer.C:
			return core.ErrTimeout
		}
	}
}

func (doc *HTMLDocument) Navigate(url values.String) error {
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

	return doc.WaitForNavigation(5000)
}
