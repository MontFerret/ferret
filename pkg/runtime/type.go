package runtime

import (
	"fmt"
	"reflect"
)

type Type string

func NewType(pkg, name string) Type {
	if pkg == "" {
		return Type(pkg)
	}

	if name == "" {
		panic("name is empty")
	}

	return Type(pkg + "." + name)
}

func (t Type) String() string {
	return string(t)
}

const (
	// Actual types
	TypeNone     = Type("none")
	TypeBoolean  = Type("boolean")
	TypeInt      = Type("int")
	TypeFloat    = Type("float")
	TypeString   = Type("string")
	TypeDateTime = Type("date_time")
	TypeArray    = Type("array")
	TypeObject   = Type("object")
	TypeBinary   = Type("binary")

	// Interfaces
	TypeCollection = Type("collection")
	TypeList       = Type("list")
	TypeMap        = Type("map")
	TypeIndexed    = Type("indexed")
	TypeKeyed      = Type("keyed")
	TypeIterable   = Type("iterable")
	TypeIterator   = Type("iterator")
	TypeMeasurable = Type("measurable")
	TypeComparable = Type("comparable")
	TypeCloneable  = Type("cloneable")
	TypeSortable   = Type("sortable")
	TypeObservable = Type("observable")
)

func typeRank(value Value) int64 {
	if value == None {
		return 0
	}

	switch value.(type) {
	case Boolean:
		return 1
	case Int:
		return 2
	case Float:
		return 3
	case String:
		return 4
	case DateTime:
		return 5
	case List:
		return 6
	case Map:
		return 7
	case Binary:
		return 8
	default:
		return -1
	}
}

func Reflect(input Value) Type {
	if input == None || input == nil {
		return TypeNone
	}

	switch input.(type) {
	case Boolean:
		return TypeBoolean
	case Int:
		return TypeInt
	case Float:
		return TypeFloat
	case String:
		return TypeString
	case DateTime:
		return TypeDateTime
	case *Array:
		return TypeArray
	case List:
		return TypeList
	case *Object:
		return TypeObject
	case Map:
		return TypeMap
	case Binary:
		return TypeBinary
	case Iterable:
		return TypeIterable
	case Iterator:
		return TypeIterator
	case Measurable:
		return TypeMeasurable
	case Observable:
		return TypeObservable
	default:
		t := reflect.TypeOf(input)

		return NewType(t.PkgPath(), t.Name())
	}
}

func CompareTypes(a, b Value) int64 {
	aRank := typeRank(a)
	bRank := typeRank(b)

	if aRank == bRank {
		return 0
	}

	if aRank < bRank {
		return -1
	}

	return 1
}

func ValidateType(value Value, required ...Type) error {
	var valid bool
	tid := Reflect(value)

	for _, t := range required {
		if tid == t {
			valid = true
			break
		}
	}

	if !valid {
		return TypeError(tid, required...)
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
			return fmt.Errorf("pair %d: %w", idx, err)
		}
	}

	return nil
}

// PairValueType is a supporting
// structure that used in validateValueTypePairs.
type PairValueType struct {
	Value Value
	Types []Type
}

// NewPairValueType it's a shortcut for creating a new PairValueType.
//
// The common pattern of using PairValueType is:
// ```
//
//	pairs := []core.PairValueType{
//	    core.PairValueType{args[0], []core.Type{types.String}},               // go vet warning
//	    core.PairValueType{Value: args[1], Types: []core.Type{types.Binary}}, // too long
//	}
//
// ```
// With NewPairValueType there is no need to type `[]core.Type{...}` and code becomes
// more readable and maintainable.
//
// That is how the code above looks like with NewPairValueType:
// ```
//
//	pairs := []core.PairValueType{
//	    core.NewPairValueType(args[0], types.String),
//	    core.NewPairValueType(args[1], types.Binary),
//	}
//
// ```
func NewPairValueType(value Value, types ...Type) PairValueType {
	return PairValueType{
		Value: value,
		Types: types,
	}
}
