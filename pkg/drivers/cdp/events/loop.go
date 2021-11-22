package events

import (
	"context"
	"math/rand"
	"sync"
)

type Loop struct {
	mu        sync.RWMutex
	sources   []SourceFactory
	listeners map[ID]map[ListenerID]Listener
}

func NewLoop(sources ...SourceFactory) *Loop {
	loop := new(Loop)
	loop.listeners = make(map[ID]map[ListenerID]Listener)
	loop.sources = sources

	return loop
}

func (loop *Loop) Run(ctx context.Context) error {
	var err error
	sources := make([]Source, 0, len(loop.sources))

	// create new sources
	for _, factory := range loop.sources {
		src, e := factory(ctx)

		if e != nil {
			err = e

			break
		}

		sources = append(sources, src)
	}

	// if error occurred
	if err != nil {
		// clean up the open ones
		for _, src := range sources {
			src.Close()
		}

		return err
	}

	for _, src := range sources {
		loop.consume(ctx, src)
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
				if ctx.Err() != nil {
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

	var snapshot []Listener
	listeners, exist := loop.listeners[eventID]

	if exist {
		snapshot = make([]Listener, 0, len(listeners))

		for _, listener := range listeners {
			snapshot = append(snapshot, listener)
		}
	}

	loop.mu.Unlock()

	for _, listener := range snapshot {
		if ctx.Err() != nil {
			return
		}

		// if returned false,
		// the handler must be removed after the call
		if !listener.Handler(ctx, message) {
			loop.mu.Lock()
			delete(loop.listeners[eventID], listener.ID)
			loop.mu.Unlock()
		}
	}
}
