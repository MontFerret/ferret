package events

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"reflect"
	"sync"
)

type (
	Event int

	EventListener func(message interface{})

	EventBroker struct {
		mu                      sync.Mutex
		listeners               map[Event][]EventListener
		cancel                  context.CancelFunc
		onLoad                  page.LoadEventFiredClient
		onReload                dom.DocumentUpdatedClient
		onAttrModified          dom.AttributeModifiedClient
		onAttrRemoved           dom.AttributeRemovedClient
		onChildNodeCountUpdated dom.ChildNodeCountUpdatedClient
		onChildNodeInserted     dom.ChildNodeInsertedClient
		onChildNodeRemoved      dom.ChildNodeRemovedClient
	}
)

var (
	//revive:disable-next-line var-declaration
	EventError                 Event = 0
	EventLoad                  Event = 1
	EventReload                Event = 2
	EventAttrModified          Event = 3
	EventAttrRemoved           Event = 4
	EventChildNodeCountUpdated Event = 5
	EventChildNodeInserted     Event = 6
	EventChildNodeRemoved      Event = 7
)

func NewEventBroker(
	onLoad page.LoadEventFiredClient,
	onReload dom.DocumentUpdatedClient,
	onAttrModified dom.AttributeModifiedClient,
	onAttrRemoved dom.AttributeRemovedClient,
	onChildNodeCountUpdated dom.ChildNodeCountUpdatedClient,
	onChildNodeInserted dom.ChildNodeInsertedClient,
	onChildNodeRemoved dom.ChildNodeRemovedClient,
) *EventBroker {
	broker := new(EventBroker)
	broker.listeners = make(map[Event][]EventListener)
	broker.onLoad = onLoad
	broker.onReload = onReload
	broker.onAttrModified = onAttrModified
	broker.onAttrRemoved = onAttrRemoved
	broker.onChildNodeCountUpdated = onChildNodeCountUpdated
	broker.onChildNodeInserted = onChildNodeInserted
	broker.onChildNodeRemoved = onChildNodeRemoved

	return broker
}

func (broker *EventBroker) AddEventListener(event Event, listener EventListener) {
	broker.mu.Lock()
	defer broker.mu.Unlock()

	listeners, ok := broker.listeners[event]

	if !ok {
		listeners = make([]EventListener, 0, 5)
	}

	broker.listeners[event] = append(listeners, listener)
}

func (broker *EventBroker) RemoveEventListener(event Event, listener EventListener) {
	broker.mu.Lock()
	defer broker.mu.Unlock()

	idx := -1

	listeners, ok := broker.listeners[event]

	if !ok {
		return
	}

	listenerPointer := reflect.ValueOf(listener).Pointer()

	for i, l := range listeners {
		itemPointer := reflect.ValueOf(l).Pointer()
		if itemPointer == listenerPointer {
			idx = i
			break
		}
	}

	if idx < 0 {
		return
	}

	var modifiedListeners []EventListener

	if len(listeners) > 1 {
		modifiedListeners = append(listeners[:idx], listeners[idx+1:]...)
	} else {
		modifiedListeners = make([]EventListener, 0, 5)
	}

	broker.listeners[event] = modifiedListeners
}

func (broker *EventBroker) Start() error {
	broker.mu.Lock()
	defer broker.mu.Unlock()

	if broker.cancel != nil {
		return core.Error(core.ErrInvalidOperation, "broker is already started")
	}

	ctx, cancel := context.WithCancel(context.Background())

	broker.cancel = cancel

	go broker.runLoop(ctx)

	return nil
}

func (broker *EventBroker) Stop() error {
	broker.mu.Lock()
	defer broker.mu.Unlock()

	if broker.cancel == nil {
		return core.Error(core.ErrInvalidOperation, "broker is already stopped")
	}

	broker.cancel()
	broker.cancel = nil

	return nil
}

func (broker *EventBroker) Close() error {
	broker.mu.Lock()
	defer broker.mu.Unlock()

	if broker.cancel != nil {
		broker.cancel()
		broker.cancel = nil
	}

	broker.onLoad.Close()
	broker.onReload.Close()
	broker.onAttrModified.Close()
	broker.onAttrRemoved.Close()
	broker.onChildNodeCountUpdated.Close()
	broker.onChildNodeInserted.Close()
	broker.onChildNodeRemoved.Close()

	return nil
}

func (broker *EventBroker) runLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-broker.onLoad.Ready():
			reply, err := broker.onLoad.Recv()

			broker.emit(EventLoad, reply, err)
		case <-broker.onReload.Ready():
			reply, err := broker.onReload.Recv()

			broker.emit(EventReload, reply, err)
		case <-broker.onAttrModified.Ready():
			reply, err := broker.onAttrModified.Recv()

			broker.emit(EventAttrModified, reply, err)
		case <-broker.onAttrRemoved.Ready():
			reply, err := broker.onAttrRemoved.Recv()

			broker.emit(EventAttrRemoved, reply, err)
		case <-broker.onChildNodeCountUpdated.Ready():
			reply, err := broker.onChildNodeCountUpdated.Recv()

			broker.emit(EventChildNodeCountUpdated, reply, err)
		case <-broker.onChildNodeInserted.Ready():
			reply, err := broker.onChildNodeInserted.Recv()

			broker.emit(EventChildNodeInserted, reply, err)
		case <-broker.onChildNodeRemoved.Ready():
			reply, err := broker.onChildNodeRemoved.Recv()

			broker.emit(EventChildNodeRemoved, reply, err)
		}
	}
}

func (broker *EventBroker) emit(event Event, message interface{}, err error) {
	if err != nil {
		event = EventError
		message = err
	}

	broker.mu.Lock()
	defer broker.mu.Unlock()

	listeners, ok := broker.listeners[event]

	if !ok {
		return
	}

	snapshot := make([]EventListener, len(listeners))
	copy(snapshot, listeners)

	go func() {
		for _, listener := range snapshot {
			listener(message)
		}
	}()
}
