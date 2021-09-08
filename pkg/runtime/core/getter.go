package core

import "context"

type (

	// Getter represents an interface of
	// complex types that needs to be used to read values by path.
	// The interface is created to let user-defined types be used in dot notation data access.
	Getter interface {
		GetIn(ctx context.Context, path []Value) (Value, PathError)
	}
)
