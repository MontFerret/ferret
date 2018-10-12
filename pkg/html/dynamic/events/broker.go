package events

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
	"reflect"
	"sync"
)

type (
	Event int

	EventListener func(message interface{})

	EventBroker struct {
		mu                     sync.Mutex
		listeners              map[Event][]EventListener
		cancel                 context.CancelFunc
		onLoad                 rpcc.Stream
		onReload               rpcc.Stream
		onAttrModified         rpcc.Stream
		onAttrRemoved          rpcc.Stream
		onChildrenCountUpdated rpcc.Stream
		onChildNodeInserted    rpcc.Stream
		onChildNodeRemoved     rpcc.Stream
	}
)

var (
	EventError                Event = 0
	EventLoad                 Event = 1
	EventReload               Event = 2
	EventAttrModified         Event = 3
	EventAttrRemoved          Event = 4
	EventChildrenCountUpdated Event = 5
	EventChildNodeInserted    Event = 6
	EventChildNodeRemoved     Event = 7
)

func NewEventBroker(
	onLoad rpcc.Stream,
	onReload rpcc.Stream,
	onAttrModified rpcc.Stream,
	onAttrRemoved rpcc.Stream,
	onChildrenCountUpdated rpcc.Stream,
	onChildNodeInserted rpcc.Stream,
	onChildNodeRemoved rpcc.Stream,
) *EventBroker {
	broker := new(EventBroker)
	broker.listeners = make(map[Event][]EventListener)
	broker.onLoad = onLoad
	broker.onReload = onReload
	broker.onAttrModified = onAttrModified
	broker.onAttrRemoved = onAttrRemoved
	broker.onChildrenCountUpdated = onChildrenCountUpdated
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
	broker.onChildrenCountUpdated.Close()
	broker.onChildNodeInserted.Close()
	broker.onChildNodeRemoved.Close()

	return nil
}

func (broker *EventBroker) runLoop(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	case <-broker.onLoad.Ready():
		msg := new(page.LoadEventFiredReply)
		err := broker.onLoad.RecvMsg(msg)

		broker.emit(EventLoad, msg, err)
	case <-broker.onReload.Ready():
		msg := new(dom.DocumentUpdatedReply)
		err := broker.onReload.RecvMsg(msg)

		broker.emit(EventReload, msg, err)
	case <-broker.onAttrModified.Ready():
		msg := new(dom.AttributeModifiedReply)
		err := broker.onAttrModified.RecvMsg(msg)

		broker.emit(EventAttrModified, msg, err)
	case <-broker.onAttrRemoved.Ready():
		msg := new(dom.AttributeRemovedReply)
		err := broker.onAttrRemoved.RecvMsg(msg)

		broker.emit(EventAttrRemoved, msg, err)
	case <-broker.onChildrenCountUpdated.Ready():
		msg := new(dom.ChildNodeCountUpdatedReply)
		err := broker.onChildrenCountUpdated.RecvMsg(msg)

		broker.emit(EventChildrenCountUpdated, msg, err)
	case <-broker.onChildNodeInserted.Ready():
		msg := new(dom.ChildNodeInsertedReply)
		err := broker.onChildNodeInserted.RecvMsg(msg)

		broker.emit(EventChildNodeInserted, msg, err)
	case <-broker.onChildNodeRemoved.Ready():
		msg := new(dom.ChildNodeRemovedReply)
		err := broker.onChildNodeRemoved.RecvMsg(msg)

		broker.emit(EventChildNodeRemoved, msg, err)
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
