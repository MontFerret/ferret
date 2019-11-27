package dom

import (
	"context"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/rpcc"

	"github.com/MontFerret/ferret/pkg/drivers/cdp/events"
)

var (
	eventReload                = events.ToType("dom_reload")
	eventAttrModified          = events.ToType("attr_modified")
	eventAttrRemoved           = events.ToType("attr_removed")
	eventChildNodeCountUpdated = events.ToType("child_count_updated")
	eventChildNodeInserted     = events.ToType("child_inserted")
	eventChildNodeRemoved      = events.ToType("child_removed")
)

type (
	ReloadListener func(ctx context.Context)

	AttrModifiedListener func(ctx context.Context, nodeID dom.NodeID, name, value string)

	AttrRemovedListener func(ctx context.Context, nodeID dom.NodeID, name string)

	ChildNodeCountUpdatedListener func(ctx context.Context, nodeID dom.NodeID, count int)

	ChildNodeInsertedListener func(ctx context.Context, nodeID, previousNodeID dom.NodeID, node dom.Node)

	ChildNodeRemovedListener func(ctx context.Context, nodeID, previousNodeID dom.NodeID)

	Manager struct {
		client *cdp.Client
		events *events.Loop
	}
)

func NewManager(
	ctx context.Context,
	client *cdp.Client,
	eventLoop *events.Loop,
) (*Manager, error) {
	onAttrModified, err := client.DOM.AttributeModified(ctx)

	if err != nil {
		return nil, err
	}

	onAttrRemoved, err := client.DOM.AttributeRemoved(ctx)

	if err != nil {
		onAttrModified.Close()
		return nil, err
	}

	onChildCountUpdated, err := client.DOM.ChildNodeCountUpdated(ctx)

	if err != nil {
		onAttrModified.Close()
		onAttrRemoved.Close()
		return nil, err
	}

	onChildNodeInserted, err := client.DOM.ChildNodeInserted(ctx)

	if err != nil {
		onAttrModified.Close()
		onAttrRemoved.Close()
		onChildCountUpdated.Close()
		return nil, err
	}

	onChildNodeRemoved, err := client.DOM.ChildNodeRemoved(ctx)

	if err != nil {
		onAttrModified.Close()
		onAttrRemoved.Close()
		onChildCountUpdated.Close()
		onChildNodeInserted.Close()
		return nil, err
	}

	eventLoop.AddSource(events.NewSource(eventAttrModified, onAttrModified, func(stream rpcc.Stream) (i interface{}, e error) {
		return stream.(dom.AttributeModifiedClient).Recv()
	}))

	eventLoop.AddSource(events.NewSource(eventAttrRemoved, onAttrRemoved, func(stream rpcc.Stream) (i interface{}, e error) {
		return stream.(dom.AttributeRemovedClient).Recv()
	}))

	eventLoop.AddSource(events.NewSource(eventChildNodeCountUpdated, onChildCountUpdated, func(stream rpcc.Stream) (i interface{}, e error) {
		return stream.(dom.ChildNodeCountUpdatedClient).Recv()
	}))

	eventLoop.AddSource(events.NewSource(eventChildNodeInserted, onChildNodeInserted, func(stream rpcc.Stream) (i interface{}, e error) {
		return stream.(dom.ChildNodeInsertedClient).Recv()
	}))

	eventLoop.AddSource(events.NewSource(eventChildNodeRemoved, onChildNodeRemoved, func(stream rpcc.Stream) (i interface{}, e error) {
		return stream.(dom.ChildNodeRemovedClient).Recv()
	}))

	return &Manager{
		client: client,
		events: eventLoop,
	}, nil
}

func (m *Manager) AddReloadListener(listener ReloadListener) events.ListenerID {
	return m.events.AddListener(eventReload, func(ctx context.Context, _ interface{}) {
		listener(ctx)
	})
}

func (m *Manager) RemoveReloadListener(listenerID events.ListenerID) {
	m.events.RemoveListener(eventReload, listenerID)
}

func (m *Manager) AddAttrModifiedListener(listener AttrModifiedListener) events.ListenerID {
	return m.events.AddListener(eventAttrModified, func(ctx context.Context, message interface{}) {
		reply := message.(dom.AttributeModifiedReply)

		listener(ctx, reply.NodeID, reply.Name, reply.Value)
	})
}

func (m *Manager) RemoveAttrModifiedListener(listenerID events.ListenerID) {
	m.events.RemoveListener(eventAttrModified, listenerID)
}

func (m *Manager) AddAttrRemovedListener(listener AttrRemovedListener) events.ListenerID {
	return m.events.AddListener(eventAttrRemoved, func(ctx context.Context, message interface{}) {
		reply := message.(dom.AttributeRemovedReply)

		listener(ctx, reply.NodeID, reply.Name)
	})
}

func (m *Manager) RemoveAttrRemovedListener(listenerID events.ListenerID) {
	m.events.RemoveListener(eventAttrRemoved, listenerID)
}

func (m *Manager) AddChildNodeCountUpdatedListener(listener ChildNodeCountUpdatedListener) events.ListenerID {
	return m.events.AddListener(eventChildNodeCountUpdated, func(ctx context.Context, message interface{}) {
		reply := message.(dom.ChildNodeCountUpdatedReply)

		listener(ctx, reply.NodeID, reply.ChildNodeCount)
	})
}

func (m *Manager) RemoveChildNodeCountUpdatedListener(listenerID events.ListenerID) {
	m.events.RemoveListener(eventChildNodeCountUpdated, listenerID)
}

func (m *Manager) AddChildNodeInsertedListener(listener ChildNodeInsertedListener) events.ListenerID {
	return m.events.AddListener(eventChildNodeInserted, func(ctx context.Context, message interface{}) {
		reply := message.(dom.ChildNodeInsertedReply)

		listener(ctx, reply.ParentNodeID, reply.PreviousNodeID, reply.Node)
	})
}

func (m *Manager) RemoveChildNodeInsertedListener(listenerID events.ListenerID) {
	m.events.RemoveListener(eventChildNodeInserted, listenerID)
}

func (m *Manager) AddChildNodeRemovedListener(listener ChildNodeRemovedListener) events.ListenerID {
	return m.events.AddListener(eventChildNodeRemoved, func(ctx context.Context, message interface{}) {
		reply := message.(dom.ChildNodeRemovedReply)

		listener(ctx, reply.ParentNodeID, reply.NodeID)
	})
}

func (m *Manager) RemoveChildNodeRemovedListener(listenerID events.ListenerID) {
	m.events.RemoveListener(eventChildNodeRemoved, listenerID)
}
