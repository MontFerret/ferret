package browser

import (
	"context"
	"fmt"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/rpcc"
	"strings"
	"time"
)

type HtmlDocument struct {
	*HtmlElement
	conn   *rpcc.Conn
	client *cdp.Client
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
			return client.DOM.Enable(ctx)
		},

		func() error {
			return client.Runtime.Enable(ctx)
		},
	)

	if err != nil {
		return nil, err
	}

	loadEventFired, err := client.Page.LoadEventFired(ctx)

	if err != nil {
		return nil, err
	}

	_, err = loadEventFired.Recv()

	if err != nil {
		return nil, err
	}

	loadEventFired.Close()

	args := dom.NewGetDocumentArgs()
	args.Depth = PointerInt(-1) // lets load the entire document

	d, err := client.DOM.GetDocument(ctx, args)

	if err != nil {
		return nil, err
	}

	return &HtmlDocument{
		&HtmlElement{client, d.Root.NodeID, d.Root, nil},
		conn,
		client,
		url,
	}, nil
}

func (doc *HtmlDocument) Close() error {
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
	)

	_, err := task.Run()

	return err
}
