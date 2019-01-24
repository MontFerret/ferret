package core

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

type (
	// Type represents runtime type id for quick type check
	//revive:disable-next-line redefines-builtin-id
	Type int64

	// Value represents an interface of
	// any type that needs to be used during runtime
	Value interface {
		json.Marshaler
		Type() Type
		String() string
		Compare(other Value) int
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

const (
	NoneType     Type = 0
	BooleanType  Type = 1
	IntType      Type = 2
	FloatType    Type = 3
	StringType   Type = 4
	DateTimeType Type = 5
	ArrayType    Type = 6
	ObjectType   Type = 7
	BinaryType   Type = 8
	CustomType   Type = 9
)

var typestr = map[Type]string{
	NoneType:     "none",
	BooleanType:  "boolean",
	IntType:      "int",
	FloatType:    "float",
	StringType:   "string",
	DateTimeType: "datetime",
	ArrayType:    "array",
	ObjectType:   "object",
	BinaryType:   "binary",
	CustomType:   "custom",
}

func (t Type) String() string {
	return typestr[t]
}

// IsTypeOf return true when value's type
// is equal to check type.
// Returns false, otherwise.
func IsTypeOf(value Value, check Type) bool {
	return value.Type() == check
}

// ValidateType checks the match of
// value's type and required types.
func ValidateType(value Value, required ...Type) error {
	var valid bool
	ct := value.Type()

	for _, t := range required {
		if ct == t {
			valid = true
			break
		}
	}

	if !valid {
		return TypeError(value.Type(), required...)
	}

	return nil
}

// ValidateValueTypePairs validate pairs of
// Values and Types.
// Returns error when type didn't match
func ValidateValueTypePairs(pairs ...PairValueType) error {
	var err error

	for idx, pair := range pairs {
		err = ValidateType(pair.Value, pair.Types...)
		if err != nil {
			return errors.Errorf("pair %d: %v", idx, err)
		}
	}

	return nil
}
