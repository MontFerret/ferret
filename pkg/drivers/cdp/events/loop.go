package events

import (
	"context"
	"math/rand"
	"sync"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type (
	operation int

	command struct {
		op      operation
		payload interface{}
	}

	Loop struct {
		mu        sync.Mutex
		cancel    context.CancelFunc
		listeners map[ID]map[ListenerID]Listener
		sources   []Source
		commands  chan command
	}
)

const (
	opAddListener operation = iota
	opRemoveListener
	opAddSource
	opRemoveSource
)

func NewLoop() *Loop {
	loop := new(Loop)
	loop.listeners = make(map[ID]map[ListenerID]Listener)
	loop.sources = make([]Source, 0, 10)
	loop.commands = make(chan command, 10)

	return loop
}

func (loop *Loop) Start() *Loop {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	if loop.cancel != nil {
		return loop
	}

	loopCtx, cancel := context.WithCancel(context.Background())

	loop.cancel = cancel

	go loop.run(loopCtx)

	return loop
}

func (loop *Loop) Stop() *Loop {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	if loop.cancel == nil {
		return loop
	}

	loop.cancel()
	loop.cancel = nil

	return loop
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

func (loop *Loop) ListenerCount(eventID ID) int {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	result := 0

	if eventID == Any {
		for _, listeners := range loop.listeners {
			result += len(listeners)
		}
	} else {
		listeners, exists := loop.listeners[eventID]

		if !exists {
			return result
		}

		result = len(listeners)
	}

	return result
}

func (loop *Loop) SourceCount() int {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	return len(loop.sources)
}

func (loop *Loop) AddSource(source Source) {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	if loop.cancel == nil {
		loop.addSourceInternal(source)

		return
	}

	loop.commands <- command{
		op:      opAddSource,
		payload: source,
	}
}

func (loop *Loop) RemoveSource(source Source) {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	if loop.cancel == nil {
		loop.removeSourceInternal(source)

		return
	}

	loop.commands <- command{
		op:      opRemoveSource,
		payload: source,
	}
}

func (loop *Loop) AddListener(eventID ID, handler Handler) ListenerID {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	id := ListenerID(rand.Int())

	listener := Listener{
		ID:      id,
		EventID: eventID,
		Handler: handler,
	}

	if loop.cancel == nil {
		loop.addListenerInternal(listener)

		return id
	}

	loop.commands <- command{
		op:      opAddListener,
		payload: listener,
	}

	return id
}

func (loop *Loop) RemoveListener(eventID ID, listenerID ListenerID) {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	if loop.cancel == nil {
		loop.removeListenerInternal(eventID, listenerID)

		return
	}

	listener := Listener{
		ID:      listenerID,
		EventID: eventID,
		Handler: nil,
	}

	loop.commands <- command{
		op:      opRemoveListener,
		payload: listener,
	}
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
		case cmd := <-loop.commands:
			if isCtxDone(ctx) {
				return
			}

			switch cmd.op {
			case opAddSource:
				loop.addSourceInternal(cmd.payload.(Source))
				// update size
				size += 1
			case opRemoveSource:
				if loop.removeSourceInternal(cmd.payload.(Source)) {
					size -= 1
				}
			case opAddListener:
				loop.addListenerInternal(cmd.payload.(Listener))
			case opRemoveListener:
				listener := cmd.payload.(Listener)
				loop.removeListenerInternal(listener.EventID, listener.ID)
			}
		case <-source.Ready():
			if isCtxDone(ctx) {
				return
			}

			event, err := source.Recv()

			loop.emit(ctx, event.ID, event.Data, err)
		default:
			continue
		}
	}
}

func (loop *Loop) addSourceInternal(src Source) {
	loop.sources = append(loop.sources, src)
}

func (loop *Loop) removeSourceInternal(event Source) bool {
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
	bucket, exists := loop.listeners[listener.EventID]

	if !exists {
		bucket = make(map[ListenerID]Listener)
	}

	bucket[listener.ID] = listener
}

func (loop *Loop) removeListenerInternal(eventID ID, listenerID ListenerID) {
	bucket, exists := loop.listeners[eventID]

	if !exists {
		return
	}

	delete(bucket, listenerID)
}

func (loop *Loop) emit(ctx context.Context, eventID ID, message interface{}, err error) {
	if err != nil {
		eventID = Error
		message = err
	}

	listeners, ok := loop.listeners[eventID]

	if !ok {
		return
	}

	for _, listener := range listeners {
		select {
		case <-ctx.Done():
			return
		default:
			listener.Handler(ctx, message)
		}
	}
}
