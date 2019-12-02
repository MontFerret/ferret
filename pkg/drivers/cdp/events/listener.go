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

// Many is a helper function that tells event loop to always call the function
func Always(fn func(ctx context.Context, message interface{})) Handler {
	return func(ctx context.Context, message interface{}) bool {
		fn(ctx, message)

		return true
	}
}

// Many is a helper function that tells event loop to call the function only once
func Once(fn func(ctx context.Context, message interface{})) Handler {
	return func(ctx context.Context, message interface{}) bool {
		fn(ctx, message)

		return false
	}
}
