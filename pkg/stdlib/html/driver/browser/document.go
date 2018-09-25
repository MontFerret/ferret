package browser

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/corpix/uarand"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/emulation"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"strings"
	"time"
)

type HtmlDocument struct {
	*HtmlElement
	conn   *rpcc.Conn
	client *cdp.Client
	events *EventBroker
	url    string
}

func NewHtmlDocument(
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

	root, err := getRootElement(ctx, client)

	if err != nil {
		return nil, err
	}

	events, err := createEventBroker(ctx, client)

	if err != nil {
		return nil, err
	}

	doc := &HtmlDocument{
		NewHtmlElement(client, root.NodeID, root),
		conn,
		client,
		events,
		url,
	}

	doc.init()

	return doc, nil
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

func getRootElement(ctx context.Context, client *cdp.Client) (dom.Node, error) {
	args := dom.NewGetDocumentArgs()
	args.Depth = PointerInt(-1) // lets load the entire document

	d, err := client.DOM.GetDocument(ctx, args)

	if err != nil {
		return dom.Node{}, err
	}

	return d.Root, nil
}

func createEventBroker(ctx context.Context, client *cdp.Client) (*EventBroker, error) {
	lfc, err := client.Page.LifecycleEvent(ctx)

	if err != nil {
		return nil, err
	}

	return NewEventBroker(lfc), nil
}

func (doc *HtmlDocument) Close() error {
	doc.events.Stop()
	doc.events.Close()

	doc.client.Page.Close(context.Background())

	return doc.conn.Close()
}

func (doc *HtmlDocument) Type() core.Type {
	return core.HtmlDocumentType
}

func (doc *HtmlDocument) String() string {
	return doc.url
}

func (doc *HtmlDocument) Compare(other core.Value) int {
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

func (doc *HtmlDocument) ClickBySelector(selector values.String) (values.Boolean, error) {
	res, err := Eval(
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
	task := NewWaitTask(
		doc.client,
		fmt.Sprintf(`
			el = document.querySelector("%s");

			if (el != null) {
				return true;
			}

			return null;
		`, selector),
		time.Millisecond*time.Duration(timeout),
		DefaultPolling,
	)

	_, err := task.Run()

	return err
}

func (doc *HtmlDocument) init() {
	// doc.events.AddListener("")
}

func (doc *HtmlDocument) reload() error {
	root, err := getRootElement(context.Background(), doc.client)

	if err != nil {
		return err
	}

	doc.url = *root.BaseURL
	doc.id = root.NodeID

	return nil
}
