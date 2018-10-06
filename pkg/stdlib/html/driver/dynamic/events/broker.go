package events

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/mafredri/cdp/rpcc"
	"reflect"
	"sync"
	"time"
)

type (
	MessageFactory func() interface{}
	EventStream    struct {
		stream  rpcc.Stream
		message MessageFactory
	}
	EventListener func(message interface{})

	EventBroker struct {
		sync.Mutex
		events    map[string]*EventStream
		listeners map[string][]EventListener
		cancel    context.CancelFunc
	}
)

func NewEventBroker() *EventBroker {
	broker := new(EventBroker)
	broker.events = make(map[string]*EventStream)
	broker.listeners = make(map[string][]EventListener)

	return broker
}

func (broker *EventBroker) AddEventStream(name string, stream rpcc.Stream, msg MessageFactory) error {
	broker.Lock()
	defer broker.Unlock()

	_, exists := broker.events[name]

	if exists {
		return core.Error(core.ErrNotUnique, name)
	}

	broker.events[name] = &EventStream{stream, msg}

	return nil
}

func (broker *EventBroker) AddEventListener(event string, listener EventListener) {
	broker.Lock()
	defer broker.Unlock()

	listeners, ok := broker.listeners[event]

	if !ok {
		listeners = make([]EventListener, 0, 5)
	}

	broker.listeners[event] = append(listeners, listener)
}

func (broker *EventBroker) RemoveEventListener(event string, listener EventListener) {
	broker.Lock()
	defer broker.Unlock()

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
	}

	broker.listeners[event] = modifiedListeners
}

func (broker *EventBroker) Start() error {
	broker.Lock()
	defer broker.Unlock()

	if broker.cancel != nil {
		return core.Error(core.ErrInvalidOperation, "broker is already started")
	}

	ctx, cancel := context.WithCancel(context.Background())

	broker.cancel = cancel

	go func() {
		counter := 0
		eventsCount := len(broker.events)

		for {
			for name, event := range broker.events {
				counter++

				select {
				case <-ctx.Done():
					return
				case <-event.stream.Ready():
					msg := event.message()
					err := event.stream.RecvMsg(msg)

					if err != nil {
						broker.emit("error", err)

						return
					}

					broker.emit(name, msg)
				default:
					// we have iterated over all events
					// lets pause
					if counter == eventsCount {
						counter = 0
						time.Sleep(DefaultPolling)
					}

					continue
				}
			}
		}
	}()

	return nil
}

func (broker *EventBroker) Stop() error {
	broker.Lock()
	defer broker.Unlock()

	if broker.cancel == nil {
		return core.Error(core.ErrInvalidOperation, "broker is already stopped")
	}

	broker.cancel()

	return nil
}

func (broker *EventBroker) Close() error {
	broker.Lock()
	defer broker.Unlock()

	if broker.cancel != nil {
		broker.cancel()
	}

	for _, event := range broker.events {
		event.stream.Close()
	}

	return nil
}

func (broker *EventBroker) emit(name string, message interface{}) {
	broker.Lock()

	listeners, ok := broker.listeners[name]

	if !ok {
		broker.Unlock()
		return
	}

	// we copy the list of listeners and unlock the broker before the execution.
	// we do it in order to avoid deadlocks during calls of event listeners
	snapshot := listeners[:]
	broker.Unlock()

	for _, listener := range snapshot {
		listener(message)
	}
}
