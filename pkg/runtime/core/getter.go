package core

import "context"

type (
	GetterPathIterator interface {
		Path() []Value
		Current() Value
		CurrentIndex() int
	}

	// Getter represents an interface of
	// complex types that needs to be used to read values by path.
	// The interface is created to let user-defined types be used in dot notation data access.
	Getter interface {
		GetIn(ctx context.Context, path []Value) (Value, PathError)
	}

	// GetterFn represents a type of helper functions that implement complex path resolutions.
	GetterFn func(ctx context.Context, path []Value, src Getter) (Value, PathError)
)
