package core

import "context"

// Cloneable represents an interface of a value that can be cloned.
// The difference between Copy and Clone is that Copy returns a shallow copy of the value
// and Clone returns a deep copy of the value.
type Cloneable interface {
	Value
	Clone(ctx context.Context) (Cloneable, error)
}

// SafeClone creates a deep copy of the given value.
// If the value does not support cloning, it returns None.
func SafeClone(ctx context.Context, origin Cloneable) Cloneable {
	cloned, err := origin.Clone(ctx)

	if err != nil {
		return None
	}

	return cloned
}
