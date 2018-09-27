package dynamic

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/common"
	"github.com/MontFerret/ferret/pkg/stdlib/html/driver/dynamic/events"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"golang.org/x/sync/errgroup"
)

func pointerInt(input int) *int {
	return &input
}

type batchFunc = func() error

func runBatch(funcs ...batchFunc) error {
	eg := errgroup.Group{}

	for _, f := range funcs {
		eg.Go(f)
	}

	return eg.Wait()
}

func parseAttrs(attrs []string) *values.Object {
	var attr values.String

	res := values.NewObject()

	for _, el := range attrs {
		str := values.NewString(el)

		if common.IsAttribute(el) {
			attr = str
			res.Set(str, values.EmptyString)
		} else {
			current, ok := res.Get(attr)

			if ok {
				res.Set(attr, current.(values.String).Concat(values.SpaceString).Concat(str))
			}
		}
	}

	return res
}

func loadInnerHtml(client *cdp.Client, id dom.NodeID) (values.String, error) {
	res, err := client.DOM.GetOuterHTML(context.Background(), dom.NewGetOuterHTMLArgs().SetNodeID(id))

	if err != nil {
		return "", err
	}

	return values.NewString(res.OuterHTML), err
}

func createChildrenArray(nodes []dom.Node) []dom.NodeID {
	children := make([]dom.NodeID, len(nodes))

	for idx, child := range nodes {
		children[idx] = child.NodeID
	}

	return children
}

func loadNodes(client *cdp.Client, broker *events.EventBroker, nodes []dom.NodeID) (*values.Array, error) {
	arr := values.NewArray(len(nodes))

	for _, id := range nodes {
		child, err := LoadElement(client, broker, id)

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
