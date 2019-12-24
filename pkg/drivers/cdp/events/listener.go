package events

import "context"

type (
	// Handler represents a function that is called when a particular event occurs
	// Returned boolean value indicates whether the handler needs to be called again
	// False value indicated that it needs to be removed and never called again
	Handler func(ctx context.Context, message interface{}) bool

	// ListenerID is an internal listener ID that can be used to unsubscribe from a particular event
	ListenerID int

	// Listener is an internal listener representation
	Listener struct {
		ID      ListenerID
		EventID ID
		Handler Handler
	}
)

// Always returns a handler wrapper that always gets executed by an event loop
func Always(fn func(ctx context.Context, message interface{})) Handler {
	return func(ctx context.Context, message interface{}) bool {
		fn(ctx, message)

		return true
	}
}

// Once returns a handler wrapper that gets executed only once by an event loop
func Once(fn func(ctx context.Context, message interface{})) Handler {
	return func(ctx context.Context, message interface{}) bool {
		fn(ctx, message)

		return false
	}
}
