package events

import (
	"context"
	"reflect"
	"sync"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Loop struct {
	mu             sync.Mutex
	cancel         context.CancelFunc
	listeners      map[Type][]EventHandler
	sources        []Source
	addSource      chan Source
	removeSource   chan Source
	addListener    chan Listener
	removeListener chan Listener
}

func NewLoop() *Loop {
	loop := new(Loop)
	loop.listeners = make(map[Type][]EventHandler)
	loop.sources = make([]Source, 0, 10)
	loop.addListener = make(chan Listener, 10)
	loop.removeListener = make(chan Listener, 10)
	loop.addSource = make(chan Source, 10)
	loop.removeSource = make(chan Source, 10)

	return loop
}

func (loop *Loop) Start() error {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	if loop.cancel != nil {
		return core.Error(core.ErrInvalidOperation, "event loop is already started")
	}

	ctx, cancel := context.WithCancel(context.Background())

	loop.cancel = cancel

	go loop.run(ctx)

	return nil
}

func (loop *Loop) Stop() error {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	if loop.cancel == nil {
		return core.Error(core.ErrInvalidOperation, "event loops is already stopped")
	}

	loop.cancel()
	loop.cancel = nil

	return nil
}

func (loop *Loop) Close() error {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	if loop.cancel != nil {
		loop.cancel()
		loop.cancel = nil
	}

	errs := make([]error, 0, len(loop.sources))

	for _, e := range loop.sources {
		if err := e.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return core.Errors(errs...)
	}

	return nil
}

func (loop *Loop) ListenerCount(eventType Type) int {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	result := 0

	if eventType == EventTypeAny {
		for _, listeners := range loop.listeners {
			result += len(listeners)
		}
	} else {
		listeners, exists := loop.listeners[eventType]

		if !exists {
			return result
		}

		result = len(listeners)
	}

	return result
}

func (loop *Loop) EventCount() int {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	return len(loop.sources)
}

func (loop *Loop) AddSource(source Source) *Loop {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	if loop.cancel == nil {
		loop.addEventInternal(source)

		return loop
	}

	loop.addSource <- source

	return loop
}

func (loop *Loop) RemoveSource(event Source) *Loop {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	if loop.cancel == nil {
		loop.removeEventInternal(event)

		return loop
	}

	loop.removeSource <- event

	return loop
}

func (loop *Loop) AddListener(eventType Type, handler EventHandler) *Loop {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	listener := Listener{
		Event:   eventType,
		Handler: handler,
	}

	if loop.cancel == nil {
		loop.addListenerInternal(listener)

		return loop
	}

	loop.addListener <- listener

	return loop
}

func (loop *Loop) RemoveListener(eventType Type, handler EventHandler) *Loop {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	listener := Listener{
		Event:   eventType,
		Handler: handler,
	}

	if loop.cancel == nil {
		loop.removeListenerInternal(listener)

		return loop
	}

	loop.removeListener <- listener

	return loop
}

// run starts running an event loop.
// It constantly iterates over each event stream.
// Additionally to that, on each iteration it checks add/remove listener/event channels.
func (loop *Loop) run(ctx context.Context) {
	size := len(loop.sources)
	counter := -1

	// in case event array is empty
	// we use this mock noop event source to simplify the logic
	noop := newNoopSource()

	for {
		counter++

		if counter >= size {
			// reset the counter
			size = len(loop.sources)
			counter = 0
		}

		var source Source

		if size > 0 {
			source = loop.sources[counter]
		} else {
			source = noop
		}

		select {
		case <-ctx.Done():
			return
		case listener := <-loop.addListener:
			loop.addListenerInternal(listener)
		case listener := <-loop.removeListener:
			loop.removeListenerInternal(listener)
		case event := <-loop.addSource:
			loop.addEventInternal(event)
			// update size
			size += 1
		case event := <-loop.removeSource:
			if loop.removeEventInternal(event) {
				size -= 1
			}
		case <-source.Ready():
			if ctxDone(ctx) {
				return
			}

			event, err := source.Recv()

			loop.emit(ctx, event.Type, event.Data, err)
		default:
			continue
		}
	}
}

func (loop *Loop) addEventInternal(src Source) {
	loop.sources = append(loop.sources, src)
}

func (loop *Loop) removeEventInternal(event Source) bool {
	idx := -1

	for i, c := range loop.sources {
		if c == event {
			idx = i
			break
		}
	}

	if idx > -1 {
		loop.sources = append(loop.sources[:idx], loop.sources[idx+1:]...)
	}

	return idx > -1
}

func (loop *Loop) addListenerInternal(listener Listener) {
	bucket, exists := loop.listeners[listener.Event]

	if !exists {
		bucket = make([]EventHandler, 0, 10)
	}

	loop.listeners[listener.Event] = append(bucket, listener.Handler)
}

func (loop *Loop) removeListenerInternal(listener Listener) {
	bucket, exists := loop.listeners[listener.Event]

	if !exists {
		return
	}

	idx := -1

	listenerPointer := reflect.ValueOf(listener.Handler).Pointer()

	for i, l := range bucket {
		itemPointer := reflect.ValueOf(l).Pointer()

		if itemPointer == listenerPointer {
			idx = i
			break
		}
	}

	if idx < 0 {
		return
	}

	var modifiedBucket []EventHandler

	if len(bucket) > 1 {
		modifiedBucket = append(bucket[:idx], bucket[idx+1:]...)
	} else {
		modifiedBucket = make([]EventHandler, 0, 5)
	}

	loop.listeners[listener.Event] = modifiedBucket
}

func (loop *Loop) emit(ctx context.Context, eventType Type, message interface{}, err error) {
	if err != nil {
		eventType = EventTypeError
		message = err
	}

	handlers, ok := loop.listeners[eventType]

	if !ok {
		return
	}

	for _, handler := range handlers {
		select {
		case <-ctx.Done():
			return
		default:
			handler(ctx, message)
		}
	}
}
