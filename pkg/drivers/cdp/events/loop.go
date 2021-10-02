package events

import (
	"context"
	"math/rand"
	"sync"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type Loop struct {
	mu        sync.RWMutex
	sources   []Source
	listeners map[ID]map[ListenerID]Listener
	cancel    context.CancelFunc
}

func NewLoop(sources ...Source) *Loop {
	loop := new(Loop)
	loop.sources = sources
	loop.listeners = make(map[ID]map[ListenerID]Listener)

	return loop
}

func (loop *Loop) Run(ctx context.Context) error {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	if loop.cancel != nil {
		return core.Error(core.ErrInvalidOperation, "loop is already running")
	}

	ctx, cancel := context.WithCancel(ctx)
	loop.cancel = cancel

	for _, source := range loop.sources {
		loop.consume(ctx, source)
	}

	return nil
}

func (loop *Loop) Close() error {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	if loop.cancel != nil {
		loop.cancel()
		loop.cancel = nil
	}

	return nil
}

func (loop *Loop) Listeners(eventID ID) int {
	loop.mu.RLock()
	defer loop.mu.RUnlock()

	bucket, exists := loop.listeners[eventID]

	if !exists {
		return 0
	}

	return len(bucket)
}

func (loop *Loop) AddListener(eventID ID, handler Handler) ListenerID {
	loop.mu.RLock()
	defer loop.mu.RUnlock()

	listener := Listener{
		ID:      ListenerID(rand.Int()),
		EventID: eventID,
		Handler: handler,
	}

	bucket, exists := loop.listeners[listener.EventID]

	if !exists {
		bucket = make(map[ListenerID]Listener)
		loop.listeners[listener.EventID] = bucket
	}

	bucket[listener.ID] = listener

	return listener.ID
}

func (loop *Loop) RemoveListener(eventID ID, listenerID ListenerID) {
	loop.mu.RLock()
	defer loop.mu.RUnlock()

	bucket, exists := loop.listeners[eventID]

	if !exists {
		return
	}

	delete(bucket, listenerID)
}

func (loop *Loop) consume(ctx context.Context, src Source) {
	go func() {
		defer func() {
			if err := src.Close(); err != nil {
				loop.emit(ctx, Error, err)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-src.Ready():
				if isCtxDone(ctx) {
					return
				}

				event, err := src.Recv()

				if err != nil {
					loop.emit(ctx, Error, err)

					return
				}

				loop.emit(ctx, event.ID, event.Data)
			}
		}
	}()
}

func (loop *Loop) emit(ctx context.Context, eventID ID, message interface{}) {
	loop.mu.Lock()
	defer loop.mu.Unlock()

	listeners, exist := loop.listeners[eventID]

	if !exist {
		return
	}

	for _, listener := range listeners {
		select {
		case <-ctx.Done():
			return
		default:
			// if returned false, it means the loops should not call the handler anymore
			if !listener.Handler(ctx, message) {
				delete(listeners, listener.ID)
			}
		}
	}
}
