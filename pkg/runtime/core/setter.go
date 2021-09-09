package core

import "context"

type (
	// Setter represents an interface of
	// complex types that needs to be used to write values by path.
	// The interface is created to let user-defined types be used in dot notation assignment.
	Setter interface {
		SetIn(ctx context.Context, path []Value, value Value) PathError
	}
)
