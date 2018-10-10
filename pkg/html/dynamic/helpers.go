package dynamic

import (
	"bytes"
	"context"
	"github.com/MontFerret/ferret/pkg/html/common"
	"github.com/MontFerret/ferret/pkg/html/dynamic/events"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/PuerkitoBio/goquery"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"strings"
)

type batchFunc = func() error

func runBatch(funcs ...batchFunc) error {
	eg := errgroup.Group{}

	for _, f := range funcs {
		eg.Go(f)
	}

	return eg.Wait()
}

func getRootElement(ctx context.Context, client *cdp.Client) (*dom.GetDocumentReply, error) {
	d, err := client.DOM.GetDocument(ctx, dom.NewGetDocumentArgs().SetDepth(1))

	if err != nil {
		return nil, err
	}

	return d, nil
}

func parseAttrs(attrs []string) *values.Object {
	var attr values.String

	res := values.NewObject()

	for _, el := range attrs {
		el = strings.TrimSpace(el)
		str := values.NewString(el)

		if common.IsAttribute(el) {
			attr = str
			res.Set(str, values.EmptyString)
		} else {
			current, ok := res.Get(attr)

			if ok {
				if current.String() != "" {
					res.Set(attr, current.(values.String).Concat(values.SpaceString).Concat(str))
				} else {
					res.Set(attr, str)
				}
			}
		}
	}

	return res
}

func loadInnerHTML(client *cdp.Client, id *HTMLElementIdentity) (values.String, error) {
	var args *dom.GetOuterHTMLArgs

	if id.objectID != "" {
		args = dom.NewGetOuterHTMLArgs().SetObjectID(id.objectID)
	} else if id.backendID > 0 {
		args = dom.NewGetOuterHTMLArgs().SetBackendNodeID(id.backendID)
	} else {
		args = dom.NewGetOuterHTMLArgs().SetNodeID(id.nodeID)
	}

	res, err := client.DOM.GetOuterHTML(context.Background(), args)

	if err != nil {
		return "", err
	}

	return values.NewString(res.OuterHTML), err
}

func loadInnerText(client *cdp.Client, id *HTMLElementIdentity) (values.String, error) {
	h, err := loadInnerHTML(client, id)

	if err != nil {
		return values.EmptyString, err
	}

	if h == values.EmptyString {
		return h, nil
	}

	return parseInnerText(h.String())
}

func parseInnerText(innerHTML string) (values.String, error) {
	buff := bytes.NewBuffer([]byte(innerHTML))

	parsed, err := goquery.NewDocumentFromReader(buff)

	if err != nil {
		return values.EmptyString, err
	}

	return values.NewString(parsed.Text()), nil
}

func createChildrenArray(nodes []dom.Node) []*HTMLElementIdentity {
	children := make([]*HTMLElementIdentity, len(nodes))

	for idx, child := range nodes {
		children[idx] = &HTMLElementIdentity{
			nodeID:    child.NodeID,
			backendID: child.BackendNodeID,
		}
	}

	return children
}

func loadNodes(
	ctx context.Context,
	logger *zerolog.Logger,
	client *cdp.Client,
	broker *events.EventBroker,
	nodes []*HTMLElementIdentity,
) (*values.Array, error) {
	arr := values.NewArray(len(nodes))

	for _, id := range nodes {
		child, err := LoadElement(ctx, logger, client, broker, id.nodeID, id.backendID)

		if err != nil {
			return nil, err
		}

		arr.Push(child)
	}

	return arr, nil
}

func contextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultTimeout)
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

func createEventBroker(client *cdp.Client) (*events.EventBroker, error) {
	ctx := context.Background()
	load, err := client.Page.LoadEventFired(ctx)

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

	destroy, err := client.DOM.DocumentUpdated(ctx)

	if err != nil {
		broker.Close()
		return nil, err
	}

	broker.AddEventStream("reload", destroy, func() interface{} {
		return new(dom.DocumentUpdatedReply)
	})

	attrModified, err := client.DOM.AttributeModified(ctx)

	if err != nil {
		broker.Close()

		return nil, err
	}

	broker.AddEventStream("attr:modified", attrModified, func() interface{} {
		return new(dom.AttributeModifiedReply)
	})

	attrRemoved, err := client.DOM.AttributeRemoved(ctx)

	if err != nil {
		broker.Close()

		return nil, err
	}

	broker.AddEventStream("attr:removed", attrRemoved, func() interface{} {
		return new(dom.AttributeRemovedReply)
	})

	childrenCount, err := client.DOM.ChildNodeCountUpdated(ctx)

	if err != nil {
		broker.Close()

		return nil, err
	}

	broker.AddEventStream("children:count", childrenCount, func() interface{} {
		return new(dom.ChildNodeCountUpdatedReply)
	})

	childrenInsert, err := client.DOM.ChildNodeInserted(ctx)

	if err != nil {
		broker.Close()

		return nil, err
	}

	broker.AddEventStream("children:inserted", childrenInsert, func() interface{} {
		return new(dom.ChildNodeInsertedReply)
	})

	childDeleted, err := client.DOM.ChildNodeRemoved(ctx)

	if err != nil {
		broker.Close()

		return nil, err
	}

	broker.AddEventStream("children:deleted", childDeleted, func() interface{} {
		return new(dom.ChildNodeRemovedReply)
	})

	return broker, nil
}
