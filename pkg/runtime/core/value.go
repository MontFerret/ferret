package core

import (
	"context"
	"encoding/json"
)

type (
	// Value represents an interface of
	// any type that needs to be used during runtime
	Value interface {
		json.Marshaler
		Type() Type
		String() string
		Compare(other Value) int64
		Unwrap() interface{}
		Hash() uint64
		Copy() Value
	}

	// Iterable represents an interface of a value that can be iterated by using an iterator.
	Iterable interface {
		Iterate(ctx context.Context) (Iterator, error)
	}

	// Iterator represents an interface of a value iterator.
	// When iterator is exhausted it must return None as a value.
	Iterator interface {
		Next(ctx context.Context) (value Value, key Value, err error)
	}
)
