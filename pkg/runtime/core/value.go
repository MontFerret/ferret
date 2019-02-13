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

	// Getter represents an interface of
	// complex types that needs to be used to read values by path.
	// The interface is created to let user-defined types be used in dot notation data access.
	Getter interface {
		GetIn(ctx context.Context, path []Value) (Value, error)
	}

	// Setter represents an interface of
	// complex types that needs to be used to write values by path.
	// The interface is created to let user-defined types be used in dot notation assignment.
	Setter interface {
		SetIn(ctx context.Context, path []Value, value Value) error
	}

	// PairValueType is a supporting
	// structure that used in validateValueTypePairs.
	PairValueType struct {
		Value Value
		Types []Type
	}
)
