package core

import (
	"context"
	"encoding/json"
)

<<<<<<< HEAD
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
=======
//revive:disable-next-line redefines-builtin-id
type Type int64

const (
	NoneType         Type = 0
	BooleanType      Type = 1
	IntType          Type = 2
	FloatType        Type = 3
	StringType       Type = 4
	DateTimeType     Type = 5
	ArrayType        Type = 6
	ObjectType       Type = 7
	HTMLElementType  Type = 8
	HTMLDocumentType Type = 9
	BinaryType       Type = 10
)

var typestr = map[Type]string{
	NoneType:         "none",
	BooleanType:      "boolean",
	IntType:          "int",
	FloatType:        "float",
	StringType:       "string",
	DateTimeType:     "datetime",
	ArrayType:        "array",
	ObjectType:       "object",
	HTMLElementType:  "HTMLElement",
	HTMLDocumentType: "HTMLDocument",
	BinaryType:       "BinaryType",
}

func (t Type) String() string {
	return typestr[t]
}

type Value interface {
	json.Marshaler
	Type() Type
	String() string
	Compare(other Value) int
	Unwrap() interface{}
	Hash() uint64
	Copy() Value
}
>>>>>>> 1c32d2a... rename method Clone to Copy

	// Iterable represents an interface of a value that can be iterated by using an iterator.
	Iterable interface {
		Iterate(ctx context.Context) (Iterator, error)
	}

	// Iterator represents an interface of a value iterator.
	// When iterator is exhausted it must return None as a value.
	Iterator interface {
		Next(ctx context.Context) (value Value, key Value, err error)
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

<<<<<<< HEAD
	// PairValueType is a supporting
	// structure that used in validateValueTypePairs.
	PairValueType struct {
		Value Value
		Types []Type
	}
)
=======
	return nil
}

func IsCloneable(value Value) bool {
	_, ok := cloneableTypes[value.Type()]
	return ok
}
>>>>>>> 6a1f475... move core.IsCloneable to value.go
