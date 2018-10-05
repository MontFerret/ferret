package core

import (
	"encoding/json"
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
	Hash() int
	Clone() Value
}

func IsTypeOf(value Value, check Type) bool {
	return value.Type() == check
}

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
