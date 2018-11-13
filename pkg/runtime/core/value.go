package core

import (
	"encoding/json"

	"github.com/pkg/errors"
)

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
	CustomType       Type = 99
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
	CustomType:       "CustomType",
}

func (t Type) String() string {
	return typestr[t]
}

// Value represents an interface of
// any type that needs to be used during runtime
type Value interface {
	json.Marshaler
	Type() Type
	String() string
	Compare(other Value) int
	Unwrap() interface{}
	Hash() uint64
	Copy() Value
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

// PairValueType is a supporting
// structure that used in validateValueTypePairs.
type PairValueType struct {
	Value Value
	Types []Type
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
