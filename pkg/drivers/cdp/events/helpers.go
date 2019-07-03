package events

import (
	"context"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/pkg/errors"
)

func WaitForLoadEvent(ctx context.Context, client *cdp.Client) error {
	loadEventFired, err := client.Page.LoadEventFired(ctx)

	if err != nil {
		return errors.Wrap(err, "failed to create load event hook")
	}

	_, err = loadEventFired.Recv()

	if err != nil {
		return err
	}

	return loadEventFired.Close()
}

func CreateEventBroker(client *cdp.Client) (*EventBroker, error) {
	var err error
	var onLoad page.LoadEventFiredClient
	var onReload dom.DocumentUpdatedClient
	var onAttrModified dom.AttributeModifiedClient
	var onAttrRemoved dom.AttributeRemovedClient
	var onChildCountUpdated dom.ChildNodeCountUpdatedClient
	var onChildNodeInserted dom.ChildNodeInsertedClient
	var onChildNodeRemoved dom.ChildNodeRemovedClient
	ctx := context.Background()

	onLoad, err = client.Page.LoadEventFired(ctx)

	if err != nil {
		return nil, err
	}

	onReload, err = client.DOM.DocumentUpdated(ctx)

	if err != nil {
		onLoad.Close()
		return nil, err
	}

	onAttrModified, err = client.DOM.AttributeModified(ctx)

	if err != nil {
		onLoad.Close()
		onReload.Close()
		return nil, err
	}

	onAttrRemoved, err = client.DOM.AttributeRemoved(ctx)

	if err != nil {
		onLoad.Close()
		onReload.Close()
		onAttrModified.Close()
		return nil, err
	}

	onChildCountUpdated, err = client.DOM.ChildNodeCountUpdated(ctx)

	if err != nil {
		onLoad.Close()
		onReload.Close()
		onAttrModified.Close()
		onAttrRemoved.Close()
		return nil, err
	}

	onChildNodeInserted, err = client.DOM.ChildNodeInserted(ctx)

	if err != nil {
		onLoad.Close()
		onReload.Close()
		onAttrModified.Close()
		onAttrRemoved.Close()
		onChildCountUpdated.Close()
		return nil, err
	}

	onChildNodeRemoved, err = client.DOM.ChildNodeRemoved(ctx)

	if err != nil {
		onLoad.Close()
		onReload.Close()
		onAttrModified.Close()
		onAttrRemoved.Close()
		onChildCountUpdated.Close()
		onChildNodeInserted.Close()
		return nil, err
	}

	broker := NewEventBroker(
		onLoad,
		onReload,
		onAttrModified,
		onAttrRemoved,
		onChildCountUpdated,
		onChildNodeInserted,
		onChildNodeRemoved,
	)

	err = broker.Start()

	if err != nil {
		onLoad.Close()
		onReload.Close()
		onAttrModified.Close()
		onAttrRemoved.Close()
		onChildCountUpdated.Close()
		onChildNodeInserted.Close()
		onChildNodeRemoved.Close()
		return nil, err
	}

	return broker, nil
}
