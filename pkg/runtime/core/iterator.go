package core

import "context"

type (
	// Iterable represents an interface of a value that can be iterated by using an iterator.
	Iterable interface {
		Iterate(ctx context.Context) (Iterator, error)
	}

	// Iterator represents an interface of an iterator.
	Iterator interface {
		HasNext(ctx context.Context) (bool, error)
		Next(ctx context.Context) (value Value, key Value, err error)
	}
)
