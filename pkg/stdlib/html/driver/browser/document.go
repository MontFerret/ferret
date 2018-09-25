package browser

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/browser/eval"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/browser/events"
	"github.com/corpix/uarand"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/emulation"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"strings"
	"sync"
	"time"
)

type HtmlDocument struct {
	sync.Mutex
	conn    *rpcc.Conn
	client  *cdp.Client
	events  *events.EventBroker
	url     string
	element *HtmlElement
}

func LoadHtmlDocument(
	ctx context.Context,
	conn *rpcc.Conn,
	url string,
) (*HtmlDocument, error) {
	if conn == nil {
		return nil, core.Error(core.ErrMissedArgument, "connection")
	}

	if url == "" {
		return nil, core.Error(core.ErrMissedArgument, "url")
	}

	client := cdp.NewClient(conn)

	err := RunBatch(
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

	err = waitForLoadEvent(ctx, client)

	if err != nil {
		return nil, err
	}

	root, err := getRootElement(client)

	if err != nil {
		return nil, err
	}

	broker, err := createEventBroker(client)

	if err != nil {
		return nil, err
	}

	return NewHtmlDocument(conn, client, root, broker), nil
}

func waitForLoadEvent(ctx context.Context, client *cdp.Client) error {
	loadEventFired, err := client.Page.LoadEventFired(ctx)

	if err != nil {
		return err
	}

	_, err = loadEventFired.Recv()

	if err != nil {
		return err
	}

	return loadEventFired.Close()
}

func getRootElement(client *cdp.Client) (dom.Node, error) {
	args := dom.NewGetDocumentArgs()
	args.Depth = PointerInt(1) // lets load the entire document

	d, err := client.DOM.GetDocument(context.Background(), args)

	if err != nil {
		return dom.Node{}, err
	}

	return d.Root, nil
}

func createEventBroker(client *cdp.Client) (*events.EventBroker, error) {
	load, err := client.Page.LoadEventFired(context.Background())

	if err != nil {
		return nil, err
	}

	broker := events.NewEventBroker()
	broker.AddEventStream("load", load, func() interface{} {
		return new(page.LoadEventFiredReply)
	})

	err = broker.Start()

	if err != nil {
		broker.Close()

		return nil, err
	}

	return broker, nil
}

func NewHtmlDocument(
	conn *rpcc.Conn,
	client *cdp.Client,
	root dom.Node,
	broker *events.EventBroker,
) *HtmlDocument {
	doc := new(HtmlDocument)
	doc.conn = conn
	doc.client = client
	doc.events = broker
	doc.element = NewHtmlElement(client, root.NodeID, root)
	doc.url = ""

	if root.BaseURL != nil {
		doc.url = *root.BaseURL
	}

	broker.AddEventListener("load", func(_ interface{}) {
		doc.Lock()
		defer doc.Unlock()

		updated, err := getRootElement(client)

		if err != nil {
			// TODO: We need somehow log all errors outside of stdout
			return
		}

		// close an old root element
		doc.element.Close()

		// create a new root element wrapper
		doc.element = NewHtmlElement(client, updated.NodeID, updated)
		doc.url = ""

		if updated.BaseURL != nil {
			doc.url = *updated.BaseURL
		}
	})

	return doc
}

func (doc *HtmlDocument) MarshalJSON() ([]byte, error) {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.MarshalJSON()
}

func (doc *HtmlDocument) Type() core.Type {
	return core.HtmlDocumentType
}

func (doc *HtmlDocument) String() string {
	doc.Lock()
	defer doc.Unlock()

	return doc.url
}

func (doc *HtmlDocument) Unwrap() interface{} {
	doc.Lock()
	defer doc.Unlock()

	return doc.element
}

func (doc *HtmlDocument) Hash() int {
	doc.Lock()
	defer doc.Unlock()

	h := sha512.New()

	out, err := h.Write([]byte(doc.url))

	if err != nil {
		return 0
	}

	return out
}

func (doc *HtmlDocument) Compare(other core.Value) int {
	doc.Lock()
	defer doc.Unlock()

	switch other.Type() {
	case core.HtmlDocumentType:
		other := other.(*HtmlDocument)

		return strings.Compare(doc.url, other.url)
	default:
		if other.Type() > core.HtmlDocumentType {
			return -1
		}

		return 1
	}
}

func (doc *HtmlDocument) Close() error {
	doc.Lock()
	defer doc.Unlock()

	doc.events.Stop()
	doc.events.Close()

	doc.client.Page.Close(context.Background())

	return doc.conn.Close()
}

func (doc *HtmlDocument) NodeType() values.Int {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.NodeType()
}

func (doc *HtmlDocument) NodeName() values.String {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.NodeName()
}

func (doc *HtmlDocument) Length() values.Int {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.Length()
}

func (doc *HtmlDocument) InnerText() values.String {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.InnerText()
}

func (doc *HtmlDocument) InnerHtml() values.String {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.InnerHtml()
}

func (doc *HtmlDocument) Value() core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.Value()
}

func (doc *HtmlDocument) GetAttributes() core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.GetAttributes()
}

func (doc *HtmlDocument) GetAttribute(name values.String) core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.GetAttribute(name)
}

func (doc *HtmlDocument) GetChildNodes() core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.GetChildNodes()
}

func (doc *HtmlDocument) GetChildNode(idx values.Int) core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.GetChildNode(idx)
}

func (doc *HtmlDocument) QuerySelector(selector values.String) core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.QuerySelector(selector)
}

func (doc *HtmlDocument) QuerySelectorAll(selector values.String) core.Value {
	doc.Lock()
	defer doc.Unlock()

	return doc.element.QuerySelectorAll(selector)
}

func (doc *HtmlDocument) ClickBySelector(selector values.String) (values.Boolean, error) {
	res, err := eval.Eval(
		doc.client,
		fmt.Sprintf(`
			var el = document.querySelector("%s");

			if (el == null) {
				return false;
			}

			var evt = new window.MouseEvent('click', { bubbles: true });
			el.dispatchEvent(evt);

			return true;
		`, selector),
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

func (doc *HtmlDocument) WaitForSelector(selector values.String, timeout values.Int) error {
	task := events.NewWaitTask(
		doc.client,
		fmt.Sprintf(`
			el = document.querySelector("%s");

			if (el != null) {
				return true;
			}

			return null;
		`, selector),
		time.Millisecond*time.Duration(timeout),
		events.DefaultPolling,
	)

	_, err := task.Run()

	return err
}

func (doc *HtmlDocument) WaitForNavigation(timeout values.Int) error {
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
