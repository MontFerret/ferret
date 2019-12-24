package events

import (
	"context"
	"math/rand"
	"sync"
)

type Loop struct {
	mu        sync.Mutex
	cancel    context.CancelFunc
	sources   *SourceCollection
	listeners *ListenerCollection
}

func NewLoop() *Loop {
	loop := new(Loop)
	loop.sources = NewSourceCollection()
	loop.listeners = NewListenerCollection()

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

	return loop.sources.Close()
}

func (loop *Loop) AddSource(source Source) {
	loop.sources.Add(source)
}

func (loop *Loop) RemoveSource(source Source) {
	loop.sources.Remove(source)
}

func (loop *Loop) AddListener(eventID ID, handler Handler) ListenerID {
	listener := Listener{
		ID:      ListenerID(rand.Int()),
		EventID: eventID,
		Handler: handler,
	}

	loop.listeners.Add(listener)

	return listener.ID
}

func (loop *Loop) RemoveListener(eventID ID, listenerID ListenerID) {
	loop.listeners.Remove(eventID, listenerID)
}

// run starts running an event loop.
// It constantly iterates over each event source.
// Additionally to that, on each iteration it checks the command channel in order to perform add/remove listener/source operations.
func (loop *Loop) run(ctx context.Context) {
	size := loop.sources.Size()
	counter := -1

	// in case event array is empty
	// we use this mock noop event source to simplify the logic
	noop := newNoopSource()

	for {
		counter++

		if counter >= size {
			// reset the counter
			size = loop.sources.Size()
			counter = 0
		}

		var source Source

		if size > 0 {
			found, err := loop.sources.Get(counter)

			if err == nil {
				source = found
			} else {
				// might be removed
				source = noop
				// force to reset counter
				counter = size
			}
		} else {
			source = noop
		}

		// commands have higher priority
		select {
		case <-ctx.Done():
			return
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

func (loop *Loop) emit(ctx context.Context, eventID ID, message interface{}, err error) {
	if err != nil {
		eventID = Error
		message = err
	}

	snapshot := loop.listeners.Values(eventID)

	for _, listener := range snapshot {
		select {
		case <-ctx.Done():
			return
		default:
			// if returned false, it means the loops should call the handler anymore
			if !listener.Handler(ctx, message) {
				loop.listeners.Remove(eventID, listener.ID)
			}
		}
	}
}
